package proxy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rueian/pgbroker/backend"
	"github.com/rueian/pgbroker/message"
)

const (
	InitMessageSizeLength = 4
)

type Server struct {
	PGResolver                    backend.PGResolver
	PGStartupMessageRewriter      backend.PGStartupMessageRewriter
	ConnInfoStore                 backend.ConnInfoStore
	ClientMessageHandlers         *ClientMessageHandlers
	ServerMessageHandlers         *ServerMessageHandlers
	ClientStreamCallbackFactories *StreamCallbackFactories
	ServerStreamCallbackFactories *StreamCallbackFactories
	OnHandleConnError             func(err error, ctx *Ctx, conn net.Conn)

	wg      sync.WaitGroup
	ln      net.Listener
	clients sync.Map
}

func (s *Server) Serve(ln net.Listener) error {
	if s.ServerMessageHandlers != nil {
		s.ServerMessageHandlers.AddHandleBackendKeyData(func(ctx *Ctx, msg *message.BackendKeyData) (data *message.BackendKeyData, e error) {
			ctx.ConnInfo.BackendProcessID = msg.ProcessID
			ctx.ConnInfo.BackendSecretKey = msg.SecretKey
			if err := s.ConnInfoStore.Save(&ctx.ConnInfo); err != nil {
				// TODO: log error
			}
			return msg, nil
		})
	} else if s.ServerStreamCallbackFactories != nil {
		s.ServerStreamCallbackFactories.SetFactory('K', func(ctx *Ctx) StreamCallback {
			return func(slice Slice) Slice {
				if !slice.Head {
					ctx.ConnInfo.BackendProcessID = binary.BigEndian.Uint32(slice.Data[0:4])
					ctx.ConnInfo.BackendSecretKey = binary.BigEndian.Uint32(slice.Data[4:8])
					if err := s.ConnInfoStore.Save(&ctx.ConnInfo); err != nil {
						// TODO: log error
					}
				}
				return slice
			}
		})
	}
	s.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		s.clients.Store(conn.LocalAddr().String(), conn)

		go func() {
			s.wg.Add(1)
			defer s.wg.Done()
			defer conn.Close()
			defer s.clients.Delete(conn.LocalAddr().String())

			c, f := context.WithCancel(context.Background())
			ctx := &Ctx{
				ClientConn: conn,
				Context:    c,
				Cancel:     f,
				AuthPhase:  PhaseStartup,
			}
			if err := s.handleConn(ctx, conn); err != nil {
				if s.OnHandleConnError != nil {
					s.OnHandleConnError(err, ctx, conn)
				}
			}
		}()
	}
}

func (s *Server) Shutdown() {
	s.ln.Close()
	s.clients.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*net.TCPConn); ok {
			conn.CloseRead()
		}
		return true
	})
	s.wg.Wait()
}

func (s *Server) handleConn(ctx *Ctx, client net.Conn) (err error) {
	if cc, ok := client.(*net.TCPConn); ok {
		if err := cc.SetKeepAlivePeriod(30 * time.Second); err != nil {
			return err
		}
		if err := cc.SetKeepAlive(true); err != nil {
			return err
		}
	}

	var server net.Conn
	var startup message.Reader
	for {
		startup, err = s.readStartupMessage(client)
		if err != nil {
			io.Copy(client, errorResp("ERROR", "08P01", err.Error()).Reader())
			return err
		}
		if m, ok := startup.(*message.CancelRequest); ok {
			info, err := s.ConnInfoStore.Find(client.RemoteAddr(), m.ProcessID, m.SecretKey)
			if info != nil {
				if server, err = net.Dial("tcp", info.ServerAddress.String()); err == nil {
					io.Copy(server, m.Reader())
					server.Close()
				}
			}
			if err != nil {
				io.Copy(client, errorResp("ERROR", "08006", err.Error()).Reader())
			}
			return err
		}
		if _, ok := startup.(*message.SSLRequest); ok {
			// TODO: SSLRequest, currently reject
			client.Write([]byte{'N'})
			continue
		}
		if m, ok := startup.(*message.StartupMessage); ok {
			resolved := make(chan bool)
			checking := make(chan error, 1)
			go func() {
				defer close(checking)
				if err := isConnReadable(resolved, client); err != nil {
					checking <- err
					ctx.Cancel()
				}
			}()

			server, err = s.PGResolver.GetPGConn(ctx.Context, client.RemoteAddr(), m.Parameters)
			close(resolved)

			if err != nil {
				select {
				case err2 := <-checking:
					if err2 != nil {
						err = errors.New("client leave before pg conn resolved: " + err2.Error())
					}
				default:
				}
				client.SetWriteDeadline(time.Now().Add(5 * time.Second))
				io.Copy(client, errorResp("ERROR", "08006", err.Error()).Reader())
				return err
			}

			ctx.ServerConn = server
			ctx.ConnInfo.ClientAddress = client.RemoteAddr()
			ctx.ConnInfo.ServerAddress = server.RemoteAddr()
			ctx.ConnInfo.StartupParameters = m.Parameters

			if s.PGStartupMessageRewriter != nil {
				m.Parameters = s.PGStartupMessageRewriter.RewriteParameters(m.Parameters)
			}

			<-checking
			break
		}
	}
	defer server.Close()
	defer s.ConnInfoStore.Delete(&ctx.ConnInfo)

	_, err = io.Copy(server, startup.Reader())
	if err != nil {
		return err
	}

	// cleanup startup
	startup = nil

	clientCh := make(chan error)
	serverCh := make(chan error)

	go func() {
		defer close(clientCh)
		if s.ClientMessageHandlers != nil {
			clientCh <- s.processMessages(ctx, client, server, s.ClientMessageHandlers)
		} else if s.ClientStreamCallbackFactories != nil {
			clientCh <- s.processStreamCallback(ctx, client, server, s.ClientStreamCallbackFactories)
		}
	}()

	go func() {
		defer close(serverCh)
		if s.ServerMessageHandlers != nil {
			serverCh <- s.processMessages(ctx, server, client, s.ServerMessageHandlers)
		} else if s.ServerStreamCallbackFactories != nil {
			serverCh <- s.processStreamCallback(ctx, server, client, s.ServerStreamCallbackFactories)
		}
	}()

	select {
	case err = <-serverCh:
		if err != nil && err != io.EOF {
			io.Copy(client, errorResp("ERROR", "08006", err.Error()).Reader())
		}
		client.Close()
		<-clientCh
		return err
	case err = <-clientCh:
		server.Close()
		<-serverCh
		return err
	}
}

func (s *Server) readStartupMessage(client io.Reader) (message.Reader, error) {
	head := make([]byte, InitMessageSizeLength)
	if _, err := io.ReadFull(client, head); err != nil {
		return nil, err
	}

	data := make([]byte, binary.BigEndian.Uint32(head)-InitMessageSizeLength)
	if _, err := io.ReadFull(client, data); err != nil {
		return nil, err
	}

	return message.ReadStartupMessage(data)
}

func (s *Server) processMessages(ctx *Ctx, r io.Reader, w io.Writer, hg MessageHandlerRegister) (err error) {
	rb := newMsgBuffer(r, 4096)
	wb := bufio.NewWriter(w)

	var ms []byte
	for {
		ms, err = rb.ReadMessage()
		if err != nil {
			return err
		}
		if handler := hg.GetHandler(ms[0]); handler != nil {
			var msg message.Reader
			if msg, err = handler(ctx, ms[5:]); err == nil {
				_, err = io.Copy(wb, msg.Reader())
			}
		} else {
			_, err = io.Copy(wb, bytes.NewReader(ms))
		}
		if err != nil {
			return err
		}
		if rb.End() {
			wb.Flush()
		}
	}
}

type Slice struct {
	Head bool
	Last bool
	Data []byte
}

func (s *Server) processStreamCallback(ctx *Ctx, r io.Reader, w io.Writer, hg *StreamCallbackFactories) (err error) {
	rb := bufio.NewReader(r)
	wb := bufio.NewWriter(w)
	sz := 4096 * 4
	cc := 0

	buf := make([]byte, sz)

	var slice Slice
	for {
		header := buf[:5] // 1 byte for message type + 4 byte for message length

		if _, err = io.ReadFull(rb, header); err != nil {
			return err
		}

		factory := hg.GetFactory(header[0])
		remaining := int(binary.BigEndian.Uint32(header[1:5])) - 4

		slice = Slice{
			Head: true,
			Last: remaining == 0,
			Data: header,
		}

		handler := factory(ctx)
		if _, err = wb.Write(handler(slice).Data); err != nil {
			return err
		}

		for remaining > 0 {
			chunk := buf[:min(sz, remaining)]

			cc, err = io.ReadFull(rb, chunk)
			if err != nil {
				return err
			}

			remaining = remaining - cc

			slice = Slice{
				Head: false,
				Last: remaining == 0,
				Data: chunk[:cc],
			}

			if _, err = wb.Write(handler(slice).Data); err != nil {
				return err
			}
		}

		if rb.Buffered() == 0 {
			if err = wb.Flush(); err != nil {
				return err
			}
		}
	}
}

type msgBuffer struct {
	buf []byte
	rs  int
	we  int
	r   io.Reader
}

func (b *msgBuffer) End() bool {
	return b.we == b.rs
}

func (b *msgBuffer) ReadMessage() (msg []byte, err error) {
	var n int

	// ensure there is a message header in the buf
	for b.we-b.rs < 5 {
		if b.we == len(b.buf) {
			// need to move msg to the beginning
			copy(b.buf, b.buf[b.rs:])
			b.we = b.we - b.rs
			b.rs = 0
		}

		n, err = b.r.Read(b.buf[b.we:])
		if err != nil {
			return
		}
		b.we += n
	}

	msgLen := int(binary.BigEndian.Uint32(b.buf[b.rs+1:b.rs+5])) + 1

	if msgLen > len(b.buf) {
		// need large buf
		buf := make([]byte, msgLen)
		copy(buf, b.buf[b.rs:b.we])
		b.buf = buf
		b.we = b.we - b.rs
		b.rs = 0
	}

	msgEnd := b.rs + msgLen

	if msgEnd > len(b.buf) {
		// need to move msg to the beginning
		copy(b.buf, b.buf[b.rs:b.we])
		b.we = b.we - b.rs
		b.rs = 0
		msgEnd = b.rs + msgLen
	}

	// ensure the full message is in the buf
	for b.we < msgEnd {
		n, err = b.r.Read(b.buf[b.we:])
		if err != nil {
			return
		}
		b.we += n
	}

	msg = b.buf[b.rs:msgEnd]
	b.rs = msgEnd
	// no remaining message in the buffer, we can use full buffer next time
	if b.we == b.rs {
		b.we = 0
		b.rs = 0
	}
	return
}

func newMsgBuffer(r io.Reader, size int) *msgBuffer {
	return &msgBuffer{buf: make([]byte, size), r: r}
}

func errorResp(s, c, m string) *message.ErrorResponse {
	return &message.ErrorResponse{
		Fields: []message.ErrorField{
			{
				Type:  byte('S'),
				Value: s,
			},
			{
				Type:  byte('C'),
				Value: c,
			},
			{
				Type:  byte('M'),
				Value: m,
			},
		},
	}
}

func isConnReadable(stop <-chan bool, conn net.Conn) error {
	timer := time.NewTicker(10 * time.Second)
	defer timer.Stop()
	defer conn.SetReadDeadline(time.Time{})
	buf := make([]byte, 1)
	for {
		select {
		case <-stop:
			return nil
		case <-timer.C:
			if err := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond)); err != nil {
				return err
			}
			if n, err := conn.Read(buf); err != nil {
				if e, ok := err.(*net.OpError); ok && e.Timeout() {
					continue
				}
				return err
			} else if n > 0 {
				return errors.New("unexpected message from client during startup")
			}
		}
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

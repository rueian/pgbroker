package proxy

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rueian/pgbroker/backend"
	"github.com/rueian/pgbroker/message"
)

const (
	MessageTypeLength     = 1
	MessageSizeLength     = 4
	InitMessageSizeLength = 4
)

type Server struct {
	PGResolver            backend.PGResolver
	ConnInfoStore         backend.ConnInfoStore
	ClientMessageHandlers *ClientMessageHandlers
	ServerMessageHandlers *ServerMessageHandlers

	stop bool
	wg   sync.WaitGroup
	ln   net.Listener
}

func (s *Server) Serve(ln net.Listener) error {
	s.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go func() {
			s.wg.Add(1)
			defer s.wg.Done()
			defer conn.Close()
			if err := s.handleConn(conn); err != nil {
				// TODO: log error
			}
		}()
	}
}

func (s *Server) Shutdown() {
	s.stop = true
	s.ln.Close()
	s.wg.Wait()
}

func (s *Server) handleConn(client net.Conn) (err error) {
	s.ServerMessageHandlers.AddHandleBackendKeyData(func(ctx *Metadata, msg *message.BackendKeyData) (data *message.BackendKeyData, e error) {
		ctx.ConnInfo.BackendProcessID = msg.ProcessID
		ctx.ConnInfo.BackendSecretKey = msg.SecretKey
		if err := s.ConnInfoStore.Save(&ctx.ConnInfo); err != nil {
			// TODO: log error
		}
		return msg, nil
	})

	ctx := &Metadata{
		Context:   context.Background(),
		AuthPhase: PhaseStartup,
	}

	var server net.Conn
	var startup message.Reader
	for {
		startup, err = s.readStartupMessage(client)
		if err != nil {
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
				// TODO: log error
			}
			return err
		}
		if _, ok := startup.(*message.SSLRequest); ok {
			// TODO: SSLRequest, currently reject
			client.Write([]byte{'N'})
			continue
		}
		if m, ok := startup.(*message.StartupMessage); ok {
			server, err = s.PGResolver.GetPGConn(client.RemoteAddr(), m.Parameters)
			if err != nil {
				return err
			}

			ctx.ConnInfo.ClientAddress = client.RemoteAddr()
			ctx.ConnInfo.ServerAddress = server.RemoteAddr()
			ctx.ConnInfo.StartupParameters = m.Parameters
			break
		}
	}
	defer server.Close()
	defer s.ConnInfoStore.Delete(&ctx.ConnInfo)

	_, err = io.Copy(server, startup.Reader())
	if err != nil {
		return err
	}

	clientCh := make(chan error)
	serverCh := make(chan error)

	go func() {
		defer close(clientCh)
		clientCh <- s.processMessages(ctx, client, server, s.ClientMessageHandlers)
	}()

	go func() {
		defer close(serverCh)
		serverCh <- s.processMessages(ctx, server, client, s.ServerMessageHandlers)
	}()

	var wait chan error

	select {
	case <-clientCh:
		wait = serverCh
	case <-serverCh:
		wait = clientCh
	}

	select {
	case err = <-wait:
	case <-time.Tick(10 * time.Second):
	}

	return err
}

func (s *Server) readStartupMessage(client io.Reader) (message.Reader, error) {
	c := container{}

	c.head = make([]byte, InitMessageSizeLength)
	if _, err := io.ReadFull(client, c.head); err != nil {
		return c, err
	}

	c.data = make([]byte, c.GetSize()-InitMessageSizeLength)
	if _, err := io.ReadFull(client, c.data); err != nil {
		return c, err
	}

	m := message.ReadStartupMessage(c.data)

	return m, nil
}

func (s *Server) processMessages(ctx *Metadata, in io.Reader, out io.Writer, hg MessageHandlerRegister) error {
	for !s.stop {
		c := container{}

		c.head = make([]byte, MessageTypeLength+MessageSizeLength)
		if _, err := io.ReadFull(in, c.head); err != nil {
			return err
		}

		if handler := hg.GetHandler(c.GetType()); handler != nil {
			c.data = make([]byte, c.GetSize()-MessageSizeLength)
			if _, err := io.ReadFull(in, c.data); err != nil {
				return err
			}

			if msg, err := handler(ctx, c.data); err != nil {
				if _, err := io.Copy(out, c.Reader()); err != nil {
					return err
				}
			} else {
				if _, err := io.Copy(out, msg.Reader()); err != nil {
					return err
				}
			}
		} else {
			c.data = make([]byte, 0)
			if _, err := io.Copy(out, c.Reader()); err != nil {
				return err
			}
			if _, err := io.CopyN(out, in, int64(c.GetSize()-MessageSizeLength)); err != nil {
				return err
			}
		}
	}
	return nil
}

type container struct {
	head []byte
	data []byte
}

func (c container) Reader() io.Reader {
	return io.MultiReader(bytes.NewReader(c.head), bytes.NewReader(c.data))
}

func (c container) GetType() byte {
	return c.head[0]
}

func (c container) GetSize() uint32 {
	return binary.BigEndian.Uint32(c.head[len(c.head)-4:])
}

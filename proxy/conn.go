package proxy

import (
	"io"
	"net"
	"time"

	"github.com/rueian/pgbroker/backend"
	"github.com/rueian/pgbroker/message"
)

func HandleConn(client net.Conn, backendResolver backend.Resolver, frontendHandlers FrontendMessageHandlers, backendHandlers BackendMessageHandlers) (err error) {
	defer client.Close()

	backendHandlers.SetHandleBackendKeyData(func(ctx *Context, msg *message.BackendKeyData) (data *message.BackendKeyData, e error) {
		ctx.BackendKeyData = msg
		return msg, nil
	})

	ctx := &Context{}

	var server net.Conn
	var initMsg message.Reader
	for {
		initMsg, err = ReadInitMessage(client)
		if err != nil {
			panic(err)
		}
		if _, ok := initMsg.(*message.CancelRequest); ok {
			// TODO: CancelRequest, currently ignored
			return
		}
		if _, ok := initMsg.(*message.SSLRequest); ok {
			// TODO: SSLRequest, currently reject
			client.Write([]byte{'N'})
			continue
		}
		if m, ok := initMsg.(*message.StartupMessage); ok {
			server, err = backendResolver.GetBackend(client.RemoteAddr(), m.Parameters)
			if err != nil {
				panic(err)
			}

			ctx.FrontendAddress = client.RemoteAddr()
			ctx.BackendAddress = server.RemoteAddr()
			break
		}
	}
	defer server.Close()

	_, err = io.Copy(server, initMsg.Reader())
	if err != nil {
		panic(err)
	}

	readCh := make(chan error)
	writeCh := make(chan error)

	go func() {
		defer close(readCh)
		readCh <- ProcessMessages(ctx, client, server, frontendHandlers)
	}()

	go func() {
		defer close(writeCh)
		writeCh <- ProcessMessages(ctx, server, client, backendHandlers)
	}()

	var wait chan error

	select {
	case <-readCh:
		wait = writeCh
	case <-writeCh:
		wait = readCh
	}

	select {
	case err = <-wait:
	case <-time.Tick(10 * time.Second):
	}

	return err
}

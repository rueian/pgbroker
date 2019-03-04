package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rueian/pgbroker/backend"
	"github.com/rueian/pgbroker/message"
	"github.com/rueian/pgbroker/proxy"
)

func main() {
	ln, err := net.Listen("tcp", ":54322")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	clientMessageHandlers := proxy.ClientMessageHandlers{}
	serverMessageHandlers := proxy.ServerMessageHandlers{}

	clientMessageHandlers.SetHandleQuery(func(ctx *proxy.Context, msg *message.Query) (query *message.Query, e error) {
		fmt.Println("Query: ", msg.QueryString)
		return msg, nil
	})
	serverMessageHandlers.SetHandleRowDescription(func(ctx *proxy.Context, msg *message.RowDescription) (data *message.RowDescription, e error) {
		ctx.RowDescription = msg
		return msg, nil
	})
	serverMessageHandlers.SetHandleDataRow(func(ctx *proxy.Context, msg *message.DataRow) (data *message.DataRow, e error) {
		fmt.Printf("DataDes\t")
		for _, f := range ctx.RowDescription.Fields {
			fmt.Printf("%s\t", f.Name)
		}
		fmt.Println("")
		fmt.Printf("DataRow\t")
		for _, f := range msg.ColumnValues {
			fmt.Printf("%s\t", string(f.DataBytes()))
		}
		fmt.Println("")
		return msg, nil
	})

	server := proxy.Server{
		PGResolver:            backend.NewStaticPGResolver(":5432"),
		ConnInfoStore:         backend.NewInMemoryConnInfoStore(),
		ServerMessageHandlers: serverMessageHandlers,
		ClientMessageHandlers: clientMessageHandlers,
	}

	go server.Serve(ln)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	server.Shutdown()
}

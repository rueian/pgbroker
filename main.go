package main

import (
	"fmt"
	"net"

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

	frontendHandlers := proxy.FrontendMessageHandlers{}
	backendHandlers := proxy.BackendMessageHandlers{}

	frontendHandlers.SetHandleQuery(func(ctx *proxy.Context, msg *message.Query) (query *message.Query, e error) {
		fmt.Println("Query: ", msg.QueryString)
		return msg, nil
	})
	backendHandlers.SetHandleRowDescription(func(ctx *proxy.Context, msg *message.RowDescription) (data *message.RowDescription, e error) {
		ctx.RowDescription = msg
		return msg, nil
	})
	backendHandlers.SetHandleDataRow(func(ctx *proxy.Context, msg *message.DataRow) (data *message.DataRow, e error) {
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

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go proxy.HandleConn(conn, backend.StaticResolver{Address: ":5432"}, frontendHandlers, backendHandlers)
	}
}

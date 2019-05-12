package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/rueian/pgbroker/backend"
	"github.com/rueian/pgbroker/message"
	"github.com/rueian/pgbroker/proxy"
)

func main() {
	{
		f, err := os.Create("/pprof/cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	ln, err := net.Listen("tcp", ":5432")
	if err != nil {
		panic(err)
	}

	clientMessageHandlers := proxy.NewClientMessageHandlers()
	serverMessageHandlers := proxy.NewServerMessageHandlers()

	clientMessageHandlers.AddHandleQuery(func(ctx *proxy.Ctx, msg *message.Query) (query *message.Query, e error) {
		fmt.Println("Query: ", msg.QueryString)
		return msg, nil
	})
	serverMessageHandlers.AddHandleRowDescription(func(ctx *proxy.Ctx, msg *message.RowDescription) (data *message.RowDescription, e error) {
		ctx.RowDescription = msg
		return msg, nil
	})
	serverMessageHandlers.AddHandleDataRow(func(ctx *proxy.Ctx, msg *message.DataRow) (data *message.DataRow, e error) {
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
		PGResolver:            backend.NewStaticPGResolver("postgres:5432"),
		ConnInfoStore:         backend.NewInMemoryConnInfoStore(),
		ServerMessageHandlers: serverMessageHandlers,
		ClientMessageHandlers: clientMessageHandlers,
	}

	go server.Serve(ln)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	server.Shutdown()

	{
		f, err := os.Create("/pprof/mem.prof")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

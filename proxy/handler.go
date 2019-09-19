package proxy

import (
	"github.com/rueian/pgbroker/message"
	"io"
)

type MessageHandler func(*Ctx, []byte) (message.Reader, error)

type MessageHandlerRegister interface {
	GetHandler(byte) MessageHandler
}

type StreamHandler func(ctx *Ctx, count int, pi io.Reader, po io.Writer)
type StreamHandler2 func(ctx *Ctx, pi *sliceChanReader, po *sliceChanWriter)

type StreamHandlerRegister interface {
	GetHandler(byte) StreamHandler
}

type StreamHandlerRegister2 interface {
	GetHandler(byte) StreamHandler2
}

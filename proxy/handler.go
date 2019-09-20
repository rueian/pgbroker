package proxy

import (
	"github.com/rueian/pgbroker/message"
)

type MessageHandler func(*Ctx, []byte) (message.Reader, error)

type MessageHandlerRegister interface {
	GetHandler(byte) MessageHandler
}

type StreamHandler func(ctx *Ctx, pi *SliceChanReader, po *SliceChanWriter)

type StreamHandlerRegister interface {
	GetHandler(byte) StreamHandler
}

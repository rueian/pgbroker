package proxy

import "github.com/rueian/pgbroker/message"

type MessageHandler func(*Metadata, []byte) (message.Reader, error)

type MessageHandlerRegister interface {
	GetHandler(byte) MessageHandler
}

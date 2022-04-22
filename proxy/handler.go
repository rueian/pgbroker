package proxy

import (
	"github.com/pioneerworks/pgbroker/message"
)

type MessageHandler func(*Ctx, []byte) (message.Reader, error)

type MessageHandlerRegister interface {
	GetHandler(byte) MessageHandler
}

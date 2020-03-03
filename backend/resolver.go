package backend

import (
	"context"
	"github.com/rueian/pgbroker/message"
	"net"
)

type PGResolver interface {
	GetPGConn(ctx context.Context, clientAddr net.Addr, parameters map[string]string) (net.Conn, error)
}

type PGStartupMessageRewriter interface {
	Rewrite(original *message.StartupMessage) *message.StartupMessage
}

type StaticPGResolver struct {
	Address string
}

func (r *StaticPGResolver) GetPGConn(ctx context.Context, clientAddr net.Addr, parameters map[string]string) (net.Conn, error) {
	return net.Dial("tcp", r.Address)
}

func NewStaticPGResolver(addr string) *StaticPGResolver {
	return &StaticPGResolver{Address: addr}
}

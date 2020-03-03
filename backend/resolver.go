package backend

import (
	"context"
	"net"
)

type PGResolver interface {
	GetPGConn(ctx context.Context, clientAddr net.Addr, parameters map[string]string) (net.Conn, error)
}

type PGStartupMessageRewriter interface {
	RewriteParameters(original map[string]string) map[string]string
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

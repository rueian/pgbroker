package backend

import (
	"net"
)

type PGResolver interface {
	GetPGConn(clientAddr net.Addr, parameters map[string]string) (net.Conn, error)
}

type StaticPGResolver struct {
	Address string
}

func (r *StaticPGResolver) GetPGConn(clientAddr net.Addr, parameters map[string]string) (net.Conn, error) {
	return net.Dial("tcp", r.Address)
}

func NewStaticPGResolver(addr string) *StaticPGResolver {
	return &StaticPGResolver{Address: addr}
}

package backend

import (
	"net"
)

type Resolver interface {
	GetBackend(clientAddr net.Addr, parameters map[string]string) (net.Conn, error)
}

type StaticResolver struct {
	Address string
}

func (r StaticResolver) GetBackend(clientAddr net.Addr, parameters map[string]string) (net.Conn, error) {
	return net.Dial("tcp", r.Address)
}

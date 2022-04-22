package proxy

import (
	"context"
	"net"

	"github.com/pioneerworks/pgbroker/backend"
	"github.com/pioneerworks/pgbroker/message"
)

type AuthPhase int

const (
	PhaseStartup AuthPhase = iota
	PhaseGSS
	PhaseSASLInit
	PhaseSASL
	PhaseOK
)

type Ctx struct {
	ClientConn     net.Conn
	ServerConn     net.Conn
	ConnInfo       backend.ConnInfo
	RowDescription *message.RowDescription
	AuthPhase      AuthPhase
	Context        context.Context
	Cancel         context.CancelFunc
}

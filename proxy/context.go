package proxy

import (
	"net"

	"github.com/rueian/pgbroker/message"
)

type AuthPhase int

const (
	PhaseStartup AuthPhase = iota
	PhaseGSS
	PhaseSASLInit
	PhaseSASL
	PhaseOK
)

type Context struct {
	FrontendAddress   net.Addr
	BackendAddress    net.Addr
	StartupParameters map[string]string
	BackendKeyData    *message.BackendKeyData
	RowDescription    *message.RowDescription
	AuthPhase         AuthPhase
}

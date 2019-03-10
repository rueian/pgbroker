package proxy

import (
	"context"

	"github.com/rueian/pgbroker/backend"
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

type Metadata struct {
	ConnInfo       backend.ConnInfo
	RowDescription *message.RowDescription
	AuthPhase      AuthPhase
	Context        context.Context
}

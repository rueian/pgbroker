package message

import "io"

type PortalSuspended struct{}

func (m *PortalSuspended) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('s').Reader()
}

func ReadPortalSuspended(raw []byte) *PortalSuspended {
	return &PortalSuspended{}
}

package message

import "io"

type Execute struct {
	PortalName          string
	MaximumNumberOfRows uint32
}

func (m *Execute) Reader() io.Reader {
	b := NewBase(len(m.PortalName) + 1 + 4)
	b.WriteString(m.PortalName)
	b.WriteUint32(m.MaximumNumberOfRows)
	return b.SetType('E').Reader()
}

func ReadExecute(raw []byte) *Execute {
	b := NewBaseFromBytes(raw)
	e := &Execute{}
	e.PortalName = b.ReadString()
	e.MaximumNumberOfRows = b.ReadUint32()
	return e
}

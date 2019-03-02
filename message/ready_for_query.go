package message

import "io"

type ReadyForQuery struct {
	Status byte
}

func (m *ReadyForQuery) Reader() io.Reader {
	b := NewBase(1)
	b.WriteByte(m.Status)
	return b.SetType('Z').Reader()
}

func ReadReadyForQuery(raw []byte) *ReadyForQuery {
	b := NewBaseFromBytes(raw)
	c := &ReadyForQuery{}
	c.Status = b.ReadByte()
	return c
}

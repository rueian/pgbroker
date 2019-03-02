package message

import "io"

type Close struct {
	TargetType byte
	TargetName string
}

func (m *Close) Reader() io.Reader {
	b := NewBase(1 + len(m.TargetName) + 1)
	b.WriteByte(m.TargetType)
	b.WriteString(m.TargetName)
	return b.SetType('C').Reader()
}

func ReadClose(raw []byte) *Close {
	b := NewBaseFromBytes(raw)
	c := &Close{}
	c.TargetType = b.ReadByte()
	c.TargetName = b.ReadString()
	return c
}

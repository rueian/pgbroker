package message

import "io"

type Describe struct {
	TargetType byte
	TargetName string
}

func (m *Describe) Reader() io.Reader {
	b := NewBase(1 + len(m.TargetName) + 1)
	b.WriteByte(m.TargetType)
	b.WriteString(m.TargetName)
	return b.SetType('D').Reader()
}

func ReadDescribe(raw []byte) *Describe {
	b := NewBaseFromBytes(raw)
	d := &Describe{}
	d.TargetType = b.ReadByte()
	d.TargetName = b.ReadString()
	return d
}

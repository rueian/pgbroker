package message

import "io"

type CommandComplete struct {
	CommandTag string
}

func (m *CommandComplete) Reader() io.Reader {
	b := NewBase(len(m.CommandTag) + 1)
	b.WriteString(m.CommandTag)
	return b.SetType('C').Reader()
}

func ReadCommandComplete(raw []byte) *CommandComplete {
	b := NewBaseFromBytes(raw)
	c := &CommandComplete{}
	c.CommandTag = b.ReadString()
	return c
}

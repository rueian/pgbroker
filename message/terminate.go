package message

import "io"

type Terminate struct{}

func (m *Terminate) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('X').Reader()
}

func ReadTerminate(raw []byte) *Terminate {
	return &Terminate{}
}

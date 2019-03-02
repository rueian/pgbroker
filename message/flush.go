package message

import "io"

type Flush struct{}

func (m *Flush) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('H').Reader()
}

func ReadFlush(raw []byte) *Flush {
	return &Flush{}
}

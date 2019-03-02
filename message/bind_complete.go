package message

import "io"

type BindComplete struct{}

func (m *BindComplete) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('2').Reader()
}

func ReadBindComplete(raw []byte) *BindComplete {
	return &BindComplete{}
}

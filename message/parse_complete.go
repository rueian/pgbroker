package message

import "io"

type ParseComplete struct{}

func (m *ParseComplete) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('1').Reader()
}

func ReadParseComplete(raw []byte) *ParseComplete {
	return &ParseComplete{}
}

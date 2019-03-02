package message

import "io"

type CopyDone struct{}

func (m *CopyDone) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('c').Reader()
}

func ReadCopyDone(raw []byte) *CopyDone {
	return &CopyDone{}
}

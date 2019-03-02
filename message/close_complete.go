package message

import "io"

type CloseComplete struct{}

func (m *CloseComplete) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('3').Reader()
}

func ReadCloseComplete(raw []byte) *CloseComplete {
	return &CloseComplete{}
}

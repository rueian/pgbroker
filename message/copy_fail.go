package message

import "io"

type CopyFail struct {
	ErrorMessage string
}

func (m *CopyFail) Reader() io.Reader {
	b := NewBase(len(m.ErrorMessage) + 1)
	b.WriteString(m.ErrorMessage)
	return b.SetType('f').Reader()
}

func ReadCopyFail(raw []byte) *CopyFail {
	b := NewBaseFromBytes(raw)
	c := &CopyFail{}
	c.ErrorMessage = b.ReadString()
	return c
}

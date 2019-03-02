package message

import "io"

type Sync struct{}

func (m *Sync) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('S').Reader()
}

func ReadSync(raw []byte) *Sync {
	return &Sync{}
}

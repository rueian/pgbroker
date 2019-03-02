package message

import "io"

type NoData struct{}

func (m *NoData) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('n').Reader()
}

func ReadNoData(raw []byte) *NoData {
	return &NoData{}
}

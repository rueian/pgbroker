package message

import "io"

type EmptyQueryResponse struct{}

func (m *EmptyQueryResponse) Reader() io.Reader {
	b := NewBase(0)
	return b.SetType('I').Reader()
}

func ReadEmptyQueryResponse(raw []byte) *EmptyQueryResponse {
	return &EmptyQueryResponse{}
}

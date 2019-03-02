package message

import "io"

type SASLInitialResponse struct {
	Mechanism string
	Response  Value
}

func (m *SASLInitialResponse) Reader() io.Reader {
	b := NewBase(len(m.Mechanism) + 1 + m.Response.Length())
	b.WriteString(m.Mechanism)
	b.WriteValue(m.Response)
	return b.SetType('p').Reader()
}

func ReadSASLInitialResponse(raw []byte) *SASLInitialResponse {
	b := NewBaseFromBytes(raw)
	resp := &SASLInitialResponse{}
	resp.Mechanism = b.ReadString()
	resp.Response = ReadValue(b)
	return resp
}

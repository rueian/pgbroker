package message

import "io"

type GSSResponse struct {
	Data []byte
}

func (m *GSSResponse) Reader() io.Reader {
	b := NewBase(len(m.Data))
	b.WriteByteN(m.Data)
	return b.SetType('p').Reader()
}

func ReadGSSResponse(raw []byte) *GSSResponse {
	b := NewBaseFromBytes(raw)
	resp := &GSSResponse{}
	resp.Data = b.ReadByteN(uint32(len(raw)))
	return resp
}

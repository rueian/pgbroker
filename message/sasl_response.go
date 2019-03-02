package message

import "io"

type SASLResponse struct {
	Data []byte
}

func (m *SASLResponse) Reader() io.Reader {
	b := NewBase(len(m.Data))
	b.WriteByteN(m.Data)
	return b.SetType('p').Reader()
}

func ReadSASLResponse(raw []byte) *SASLResponse {
	b := NewBaseFromBytes(raw)
	resp := &SASLResponse{}
	resp.Data = b.ReadByteN(uint32(len(raw)))
	return resp
}

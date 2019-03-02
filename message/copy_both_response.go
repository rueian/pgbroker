package message

import "io"

type CopyBothResponse struct {
	OverallFormat     uint8
	ColumnFormatCodes []uint16
}

func (m *CopyBothResponse) Reader() io.Reader {
	b := NewBase(1 + 2 + 2*len(m.ColumnFormatCodes))
	b.WriteUint8(m.OverallFormat)
	b.WriteUint16(uint16(len(m.ColumnFormatCodes)))
	for _, c := range m.ColumnFormatCodes {
		b.WriteUint16(c)
	}
	return b.SetType('W').Reader()
}

func ReadCopyBothResponse(raw []byte) *CopyBothResponse {
	b := NewBaseFromBytes(raw)
	resp := &CopyBothResponse{}
	resp.OverallFormat = b.ReadUint8()
	resp.ColumnFormatCodes = make([]uint16, b.ReadUint16())
	for i := 0; i < len(resp.ColumnFormatCodes); i++ {
		resp.ColumnFormatCodes[i] = b.ReadUint16()
	}
	return resp
}

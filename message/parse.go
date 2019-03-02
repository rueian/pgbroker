package message

import "io"

type Parse struct {
	PreparedStatementName string
	QueryString           string
	ParameterIDs          []uint32
}

func (m *Parse) Reader() io.Reader {
	b := NewBase(len(m.PreparedStatementName) + 1 + len(m.QueryString) + 1 + 2 + 4*len(m.ParameterIDs))
	b.WriteString(m.PreparedStatementName)
	b.WriteString(m.QueryString)
	b.WriteUint16(uint16(len(m.ParameterIDs)))
	for _, id := range m.ParameterIDs {
		b.WriteUint32(id)
	}
	return b.SetType('P').Reader()
}

func ReadParse(raw []byte) *Parse {
	b := NewBaseFromBytes(raw)
	p := &Parse{}
	p.PreparedStatementName = b.ReadString()
	p.QueryString = b.ReadString()
	p.ParameterIDs = make([]uint32, b.ReadUint16())
	for i := 0; i < len(p.ParameterIDs); i++ {
		p.ParameterIDs[i] = b.ReadUint32()
	}
	return p
}

package message

import "io"

type ParameterDescription struct {
	ParameterIDs []uint32
}

func (m *ParameterDescription) Reader() io.Reader {
	b := NewBase(2 + 4*len(m.ParameterIDs))
	b.WriteUint16(uint16(len(m.ParameterIDs)))
	for _, id := range m.ParameterIDs {
		b.WriteUint32(id)
	}
	return b.SetType('t').Reader()
}

func ReadParameterDescription(raw []byte) *ParameterDescription {
	b := NewBaseFromBytes(raw)
	desc := &ParameterDescription{}
	desc.ParameterIDs = make([]uint32, b.ReadUint16())
	for i := 0; i < len(desc.ParameterIDs); i++ {
		desc.ParameterIDs[i] = b.ReadUint32()
	}
	return desc
}

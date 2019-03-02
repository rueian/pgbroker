package message

import "io"

type RowDescription struct {
	Fields []Field
}

type Field struct {
	Name            string
	TableID         uint32
	AttributeNumber uint16
	DataType        uint32
	DataSize        uint16
	ModifierType    uint32
	FormatCode      uint16
}

func (m *RowDescription) Reader() io.Reader {
	length := 2
	for _, f := range m.Fields {
		length += len(f.Name) + 1
		length += 4 + 2 + 4 + 2 + 4 + 2
	}
	b := NewBase(length)
	b.WriteUint16(uint16(len(m.Fields)))
	for _, f := range m.Fields {
		b.WriteString(f.Name)
		b.WriteUint32(f.TableID)
		b.WriteUint16(f.AttributeNumber)
		b.WriteUint32(f.DataType)
		b.WriteUint16(f.DataSize)
		b.WriteUint32(f.ModifierType)
		b.WriteUint16(f.FormatCode)
	}
	return b.SetType('T').Reader()
}

func ReadRowDescription(raw []byte) *RowDescription {
	b := NewBaseFromBytes(raw)
	desc := &RowDescription{}
	desc.Fields = make([]Field, b.ReadUint16())
	for i := 0; i < len(desc.Fields); i++ {
		desc.Fields[i].Name = b.ReadString()
		desc.Fields[i].TableID = b.ReadUint32()
		desc.Fields[i].AttributeNumber = b.ReadUint16()
		desc.Fields[i].DataType = b.ReadUint32()
		desc.Fields[i].DataSize = b.ReadUint16()
		desc.Fields[i].ModifierType = b.ReadUint32()
		desc.Fields[i].FormatCode = b.ReadUint16()
	}
	return desc
}

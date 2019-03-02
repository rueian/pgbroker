package message

import "io"

type Bind struct {
	PortalName              string
	PreparedStatementName   string
	ParameterFormatCodes    []uint16
	ParameterValues         []Value
	ResultColumnFormatCodes []uint16
}

func (m *Bind) Reader() io.Reader {
	length := len(m.PortalName) + 1 + len(m.PreparedStatementName) + 1
	length += 2 + len(m.ParameterFormatCodes)*2
	length += 2
	for _, v := range m.ParameterValues {
		length += v.Length()
	}
	length += 2 + len(m.ResultColumnFormatCodes)*2

	b := NewBase(length)
	b.WriteString(m.PortalName)
	b.WriteString(m.PreparedStatementName)
	b.WriteUint16(uint16(len(m.ParameterFormatCodes)))
	for _, c := range m.ParameterFormatCodes {
		b.WriteUint16(c)
	}
	b.WriteUint16(uint16(len(m.ParameterValues)))
	for _, v := range m.ParameterValues {
		b.WriteValue(v)
	}
	b.WriteUint16(uint16(len(m.ResultColumnFormatCodes)))
	for _, c := range m.ResultColumnFormatCodes {
		b.WriteUint16(c)
	}
	return b.SetType('B').Reader()
}

func ReadBind(raw []byte) *Bind {
	b := NewBaseFromBytes(raw)

	bind := &Bind{}
	bind.PortalName = b.ReadString()
	bind.PreparedStatementName = b.ReadString()

	bind.ParameterFormatCodes = make([]uint16, b.ReadUint16())
	for i := 0; i < len(bind.ParameterFormatCodes); i++ {
		bind.ParameterFormatCodes[i] = b.ReadUint16()
	}

	bind.ParameterValues = make([]Value, b.ReadUint16())
	for i := 0; i < len(bind.ParameterValues); i++ {
		bind.ParameterValues[i] = ReadValue(b)
	}

	bind.ResultColumnFormatCodes = make([]uint16, b.ReadUint16())
	for i := 0; i < len(bind.ResultColumnFormatCodes); i++ {
		bind.ResultColumnFormatCodes[i] = b.ReadUint16()
	}

	return bind
}

package message

import "io"

type FunctionCall struct {
	ID                  uint32
	ArgumentFormatCodes []uint16
	Arguments           []Value
	ResultFormatCode    uint16
}

func (m *FunctionCall) Reader() io.Reader {
	length := 4 + 2 + 2*len(m.Arguments) + 2 + 2
	for _, a := range m.Arguments {
		length += a.Length()
	}
	b := NewBase(length)
	b.WriteUint32(m.ID)
	b.WriteUint16(uint16(len(m.ArgumentFormatCodes)))
	for _, c := range m.ArgumentFormatCodes {
		b.WriteUint16(c)
	}
	b.WriteUint16(uint16(len(m.Arguments)))
	for _, a := range m.Arguments {
		b.WriteValue(a)
	}
	b.WriteUint16(m.ResultFormatCode)
	return b.SetType('F').Reader()
}

func ReadFunctionCall(raw []byte) *FunctionCall {
	b := NewBaseFromBytes(raw)
	call := &FunctionCall{}
	call.ID = b.ReadUint32()
	call.ArgumentFormatCodes = make([]uint16, b.ReadUint16())
	for i := 0; i < len(call.ArgumentFormatCodes); i++ {
		call.ArgumentFormatCodes[i] = b.ReadUint16()
	}
	call.Arguments = make([]Value, b.ReadUint16())
	for i := 0; i < len(call.Arguments); i++ {
		call.Arguments[i] = ReadValue(b)
	}
	call.ResultFormatCode = b.ReadUint16()
	return call
}

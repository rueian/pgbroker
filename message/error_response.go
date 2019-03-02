package message

import "io"

type ErrorResponse struct {
	Fields []ErrorField
}

type ErrorField struct {
	Type  byte
	Value string
}

func (m *ErrorResponse) Reader() io.Reader {
	length := 1
	for _, f := range m.Fields {
		length += 1
		length += len(f.Value) + 1
	}

	b := NewBase(length)
	for _, f := range m.Fields {
		b.WriteByte(f.Type)
		b.WriteString(f.Value)
	}
	b.WriteByte(0)

	return b.SetType('E').Reader()
}

func ReadErrorResponse(raw []byte) *ErrorResponse {
	b := NewBaseFromBytes(raw)
	e := &ErrorResponse{Fields: make([]ErrorField, 0)}
	for !b.isEnd() {
		t := b.ReadByte()
		if t != 0 {
			v := b.ReadString()
			e.Fields = append(e.Fields, ErrorField{
				Type:  t,
				Value: v,
			})
		}
	}
	return e
}

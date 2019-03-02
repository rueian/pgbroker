package message

import "io"

type NoticeResponse struct {
	Fields []NoticeField
}

type NoticeField struct {
	Type  byte
	Value string
}

func (m *NoticeResponse) Reader() io.Reader {
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

	return b.SetType('N').Reader()
}

func ReadNoticeResponse(raw []byte) *NoticeResponse {
	b := NewBaseFromBytes(raw)
	e := &NoticeResponse{Fields: make([]NoticeField, 0)}
	for !b.isEnd() {
		t := b.ReadByte()
		if t != 0 {
			v := b.ReadString()
			e.Fields = append(e.Fields, NoticeField{
				Type:  t,
				Value: v,
			})
		}
	}
	return e
}

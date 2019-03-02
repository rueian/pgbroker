package message

import "io"

type Query struct {
	QueryString string
}

func (m *Query) Reader() io.Reader {
	b := NewBase(len(m.QueryString) + 1)
	b.WriteString(m.QueryString)
	return b.SetType('Q').Reader()
}

func ReadQuery(raw []byte) *Query {
	b := NewBaseFromBytes(raw)
	msg := &Query{}
	msg.QueryString = b.ReadString()
	return msg
}

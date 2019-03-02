package message

import "io"

type PasswordMessage struct {
	Password string
}

func (m *PasswordMessage) Reader() io.Reader {
	b := NewBase(len(m.Password) + 1)
	b.WriteString(m.Password)
	return b.SetType('p').Reader()
}

func ReadPasswordMessage(raw []byte) *PasswordMessage {
	b := NewBaseFromBytes(raw)
	msg := &PasswordMessage{}
	msg.Password = b.ReadString()
	return msg
}

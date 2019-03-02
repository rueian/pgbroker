package message

import (
	"io"
)

type AuthenticationOk struct {
	ID uint32
}
type AuthenticationKerberosV5 struct {
	ID uint32
}
type AuthenticationCleartextPassword struct {
	ID uint32
}
type AuthenticationMD5Password struct {
	ID   uint32
	Salt []byte
}
type AuthenticationSCMCredential struct {
	ID uint32
}
type AuthenticationGSS struct {
	ID uint32
}
type AuthenticationSSPI struct {
	ID uint32
}
type AuthenticationGSSContinue struct {
	ID   uint32
	Data []byte
}
type AuthenticationSASL struct {
	ID         uint32
	Mechanisms []string
}
type AuthenticationSASLContinue struct {
	ID   uint32
	Data []byte
}
type AuthenticationSASLFinal struct {
	ID   uint32
	Data []byte
}

func (m *AuthenticationOk) Reader() io.Reader {
	b := NewBase(4)
	b.WriteUint32(m.ID)
	return b.SetType('R').Reader()
}
func (m *AuthenticationKerberosV5) Reader() io.Reader {
	b := NewBase(4)
	b.WriteUint32(m.ID)
	return b.SetType('R').Reader()
}
func (m *AuthenticationCleartextPassword) Reader() io.Reader {
	b := NewBase(4)
	b.WriteUint32(m.ID)
	return b.SetType('R').Reader()
}
func (m *AuthenticationMD5Password) Reader() io.Reader {
	b := NewBase(4 + len(m.Salt))
	b.WriteUint32(m.ID)
	b.WriteByteN(m.Salt)
	return b.SetType('R').Reader()
}
func (m *AuthenticationSCMCredential) Reader() io.Reader {
	b := NewBase(4)
	b.WriteUint32(m.ID)
	return b.SetType('R').Reader()
}
func (m *AuthenticationGSS) Reader() io.Reader {
	b := NewBase(4)
	b.WriteUint32(m.ID)
	return b.SetType('R').Reader()
}
func (m *AuthenticationSSPI) Reader() io.Reader {
	b := NewBase(4)
	b.WriteUint32(m.ID)
	return b.SetType('R').Reader()
}
func (m *AuthenticationGSSContinue) Reader() io.Reader {
	b := NewBase(4 + len(m.Data))
	b.WriteUint32(m.ID)
	b.WriteByteN(m.Data)
	return b.SetType('R').Reader()
}
func (m *AuthenticationSASL) Reader() io.Reader {
	listLength := len(m.Mechanisms) + 1
	for _, str := range m.Mechanisms {
		listLength += len(str)
	}

	b := NewBase(4 + listLength)
	b.WriteUint32(m.ID)
	for _, str := range m.Mechanisms {
		b.WriteString(str)
	}
	b.WriteByte(0)

	return b.SetType('R').Reader()
}
func (m *AuthenticationSASLContinue) Reader() io.Reader {
	b := NewBase(4 + len(m.Data))
	b.WriteUint32(m.ID)
	b.WriteByteN(m.Data)
	return b.SetType('R').Reader()
}
func (m *AuthenticationSASLFinal) Reader() io.Reader {
	b := NewBase(4 + len(m.Data))
	b.WriteUint32(m.ID)
	b.WriteByteN(m.Data)
	return b.SetType('R').Reader()
}

func ReadAuthentication(raw []byte) interface{} {
	b := NewBaseFromBytes(raw)
	id := b.ReadUint32()

	switch id {
	case 0:
		return &AuthenticationOk{ID: id}
	case 2:
		return &AuthenticationKerberosV5{ID: id}
	case 3:
		return &AuthenticationCleartextPassword{ID: id}
	case 5:
		return &AuthenticationMD5Password{ID: id, Salt: b.ReadByteN(4)}
	case 6:
		return &AuthenticationSCMCredential{ID: id}
	case 7:
		return &AuthenticationGSS{ID: id}
	case 8:
		return &AuthenticationGSSContinue{ID: id, Data: b.ReadByteN(uint32(len(raw)) - 4)}
	case 9:
		return &AuthenticationSSPI{ID: id}
	case 10:
		var list []string
		for !b.isEnd() {
			if s := b.ReadString(); s != "" {
				list = append(list, s)
			}
		}
		return &AuthenticationSASL{ID: id, Mechanisms: list}
	case 11:
		return &AuthenticationSASLContinue{ID: id, Data: b.ReadByteN(uint32(len(raw)) - 4)}
	case 12:
		return &AuthenticationSASLFinal{ID: id, Data: b.ReadByteN(uint32(len(raw)) - 4)}
	}

	return nil
}

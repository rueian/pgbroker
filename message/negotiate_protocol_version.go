package message

import "io"

type NegotiateProtocolVersion struct {
	NewestMinorProtocolVersion uint32
	Options                    []string
}

func (m *NegotiateProtocolVersion) Reader() io.Reader {
	length := 4 + 4
	for _, o := range m.Options {
		length += len(o) + 1
	}
	b := NewBase(length)
	b.WriteUint32(m.NewestMinorProtocolVersion)
	b.WriteUint32(uint32(len(m.Options)))
	for _, o := range m.Options {
		b.WriteString(o)
	}
	return b.SetType('v').Reader()
}

func ReadNegotiateProtocolVersion(raw []byte) *NegotiateProtocolVersion {
	b := NewBaseFromBytes(raw)
	msg := &NegotiateProtocolVersion{}
	msg.NewestMinorProtocolVersion = b.ReadUint32()
	msg.Options = make([]string, b.ReadUint32())
	for i := 0; i < len(msg.Options); i++ {
		msg.Options[i] = b.ReadString()
	}
	return msg
}

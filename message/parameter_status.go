package message

import "io"

type ParameterStatus struct {
	Name  string
	Value string
}

func (m *ParameterStatus) Reader() io.Reader {
	b := NewBase(len(m.Name) + 1 + len(m.Value) + 1)
	b.WriteString(m.Name)
	b.WriteString(m.Value)
	return b.SetType('S').Reader()
}

func ReadParameterStatus(raw []byte) *ParameterStatus {
	b := NewBaseFromBytes(raw)
	status := &ParameterStatus{}
	status.Name = b.ReadString()
	status.Value = b.ReadString()
	return status
}

package message

import "io"

type BackendKeyData struct {
	ProcessID uint32
	SecretKey uint32
}

func (m *BackendKeyData) Reader() io.Reader {
	b := NewBase(8)
	b.WriteUint32(m.ProcessID)
	b.WriteUint32(m.SecretKey)
	return b.SetType('K').Reader()
}

func ReadBackendKeyData(raw []byte) *BackendKeyData {
	b := NewBaseFromBytes(raw)
	data := &BackendKeyData{}
	data.ProcessID = b.ReadUint32()
	data.SecretKey = b.ReadUint32()
	return data
}

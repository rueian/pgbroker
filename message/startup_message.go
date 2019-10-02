package message

import (
	"errors"
	"io"
)

type StartupMessage struct {
	ProtocolVersion uint32
	Parameters      map[string]string
}

type CancelRequest struct {
	RequestCode uint32
	ProcessID   uint32
	SecretKey   uint32
}

type SSLRequest struct {
	RequestCode uint32
}

func (m *StartupMessage) Reader() io.Reader {
	length := 4
	for k, v := range m.Parameters {
		length += len(k) + 1
		length += len(v) + 1
	}
	length += 1
	b := NewBase(length)
	b.WriteUint32(m.ProtocolVersion)
	for k, v := range m.Parameters {
		b.WriteString(k)
		b.WriteString(v)
	}
	b.WriteByte(0)
	return b.Reader()
}

func (m *CancelRequest) Reader() io.Reader {
	b := NewBase(12)
	b.WriteUint32(m.RequestCode)
	b.WriteUint32(m.ProcessID)
	b.WriteUint32(m.SecretKey)
	return b.Reader()
}

func (m *SSLRequest) Reader() io.Reader {
	b := NewBase(4)
	b.WriteUint32(m.RequestCode)
	return b.Reader()
}

func ReadStartupMessage(raw []byte) (Reader, error) {
	b := NewBaseFromBytes(raw)
	code := b.ReadUint32()
	if code == 80877103 {
		return &SSLRequest{RequestCode: code}, nil
	}

	if code == 80877102 {
		req := &CancelRequest{RequestCode: code}
		req.ProcessID = b.ReadUint32()
		req.SecretKey = b.ReadUint32()
		return req, nil
	}

	if majorVersion := code >> 16; majorVersion < 3 {
		return nil, errors.New("pg protocol < 3.0 is not supported")
	}

	startup := &StartupMessage{ProtocolVersion: code, Parameters: make(map[string]string)}
	for {
		k := b.ReadString()
		if k == "" {
			break
		}
		startup.Parameters[k] = b.ReadString()
	}
	return startup, nil
}

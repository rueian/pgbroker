package message

import "io"

type NotificationResponse struct {
	ProcessID   uint32
	ChannelName string
	Payload     string
}

func (m *NotificationResponse) Reader() io.Reader {
	b := NewBase(4 + len(m.ChannelName) + 1 + len(m.Payload) + 1)
	b.WriteUint32(m.ProcessID)
	b.WriteString(m.ChannelName)
	b.WriteString(m.Payload)
	return b.SetType('A').Reader()
}

func ReadNotificationResponse(raw []byte) *NotificationResponse {
	b := NewBaseFromBytes(raw)
	resp := &NotificationResponse{}
	resp.ProcessID = b.ReadUint32()
	resp.ChannelName = b.ReadString()
	resp.Payload = b.ReadString()
	return resp
}

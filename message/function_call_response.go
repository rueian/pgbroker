package message

import "io"

type FunctionCallResponse struct {
	Value Value
}

func (m *FunctionCallResponse) Reader() io.Reader {
	b := NewBase(m.Value.Length())
	b.WriteValue(m.Value)
	return b.SetType('V').Reader()
}

func ReadFunctionCallResponse(raw []byte) *FunctionCallResponse {
	b := NewBaseFromBytes(raw)
	resp := &FunctionCallResponse{}
	resp.Value = ReadValue(b)
	return resp
}

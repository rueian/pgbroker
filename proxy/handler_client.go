package proxy

import "github.com/rueian/pgbroker/message"

type HandleBind func(ctx *Metadata, msg *message.Bind) (*message.Bind, error)
type HandleClose func(ctx *Metadata, msg *message.Close) (*message.Close, error)
type HandleCopyFail func(ctx *Metadata, msg *message.CopyFail) (*message.CopyFail, error)
type HandleDescribe func(ctx *Metadata, msg *message.Describe) (*message.Describe, error)
type HandleExecute func(ctx *Metadata, msg *message.Execute) (*message.Execute, error)
type HandleFlush func(ctx *Metadata, msg *message.Flush) (*message.Flush, error)
type HandleFunctionCall func(ctx *Metadata, msg *message.FunctionCall) (*message.FunctionCall, error)
type HandleParse func(ctx *Metadata, msg *message.Parse) (*message.Parse, error)
type HandlePasswordMessage func(ctx *Metadata, msg *message.PasswordMessage) (*message.PasswordMessage, error)
type HandleGSSResponse func(ctx *Metadata, msg *message.GSSResponse) (*message.GSSResponse, error)
type HandleSASLInitialResponse func(ctx *Metadata, msg *message.SASLInitialResponse) (*message.SASLInitialResponse, error)
type HandleSASLResponse func(ctx *Metadata, msg *message.SASLResponse) (*message.SASLResponse, error)
type HandleQuery func(ctx *Metadata, msg *message.Query) (*message.Query, error)
type HandleSync func(ctx *Metadata, msg *message.Sync) (*message.Sync, error)
type HandleTerminate func(ctx *Metadata, msg *message.Terminate) (*message.Terminate, error)
type HandleCopyData func(ctx *Metadata, msg *message.CopyData) (*message.CopyData, error)
type HandleCopyDone func(ctx *Metadata, msg *message.CopyDone) (*message.CopyDone, error)

type HandleAuthenticationResponse struct {
	HandlePasswordMessage     HandlePasswordMessage
	HandleGSSResponse         HandleGSSResponse
	HandleSASLInitialResponse HandleSASLInitialResponse
	HandleSASLResponse        HandleSASLResponse
}

type ClientMessageHandlers map[byte]MessageHandler

func (s *ClientMessageHandlers) SetHandleBind(h HandleBind) {
	if h == nil {
		return
	}
	(*s)['B'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadBind(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleClose(h HandleClose) {
	if h == nil {
		return
	}
	(*s)['C'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadClose(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleCopyFail(h HandleCopyFail) {
	if h == nil {
		return
	}
	(*s)['f'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyFail(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleDescribe(h HandleDescribe) {
	if h == nil {
		return
	}
	(*s)['D'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadDescribe(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleExecute(h HandleExecute) {
	if h == nil {
		return
	}
	(*s)['E'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadExecute(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleFlush(h HandleFlush) {
	if h == nil {
		return
	}
	(*s)['H'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadFlush(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleFunctionCall(h HandleFunctionCall) {
	if h == nil {
		return
	}
	(*s)['F'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadFunctionCall(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleParse(h HandleParse) {
	if h == nil {
		return
	}
	(*s)['P'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadParse(raw))
	}
}
func (s *ClientMessageHandlers) SetHandlePMessage(h HandleAuthenticationResponse) {
	(*s)['p'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		switch ctx.AuthPhase {
		case PhaseSASLInit:
			if h.HandleSASLInitialResponse == nil {
				return message.ReadSASLInitialResponse(raw), nil
			}
			return h.HandleSASLInitialResponse(ctx, message.ReadSASLInitialResponse(raw))
		case PhaseSASL:
			if h.HandleSASLResponse == nil {
				return message.ReadSASLResponse(raw), nil
			}
			return h.HandleSASLResponse(ctx, message.ReadSASLResponse(raw))
		case PhaseGSS:
			if h.HandleGSSResponse == nil {
				return message.ReadGSSResponse(raw), nil
			}
			return h.HandleGSSResponse(ctx, message.ReadGSSResponse(raw))
		default:
			return h.HandlePasswordMessage(ctx, message.ReadPasswordMessage(raw))
		}
	}
}
func (s *ClientMessageHandlers) SetHandleQuery(h HandleQuery) {
	if h == nil {
		return
	}
	(*s)['Q'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadQuery(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleSync(h HandleSync) {
	if h == nil {
		return
	}
	(*s)['S'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadSync(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleTerminate(h HandleTerminate) {
	if h == nil {
		return
	}
	(*s)['X'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadTerminate(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleCopyData(h HandleCopyData) {
	if h == nil {
		return
	}
	(*s)['d'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyData(raw))
	}
}
func (s *ClientMessageHandlers) SetHandleCopyDone(h HandleCopyDone) {
	if h == nil {
		return
	}
	(*s)['c'] = func(ctx *Metadata, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyDone(raw))
	}
}

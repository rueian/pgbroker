package proxy

import "github.com/rueian/pgbroker/message"

type HandleBind func(ctx *Context, msg *message.Bind) (*message.Bind, error)
type HandleClose func(ctx *Context, msg *message.Close) (*message.Close, error)
type HandleCopyFail func(ctx *Context, msg *message.CopyFail) (*message.CopyFail, error)
type HandleDescribe func(ctx *Context, msg *message.Describe) (*message.Describe, error)
type HandleExecute func(ctx *Context, msg *message.Execute) (*message.Execute, error)
type HandleFlush func(ctx *Context, msg *message.Flush) (*message.Flush, error)
type HandleFunctionCall func(ctx *Context, msg *message.FunctionCall) (*message.FunctionCall, error)
type HandleParse func(ctx *Context, msg *message.Parse) (*message.Parse, error)
type HandlePasswordMessage func(ctx *Context, msg *message.PasswordMessage) (*message.PasswordMessage, error)
type HandleGSSResponse func(ctx *Context, msg *message.GSSResponse) (*message.GSSResponse, error)
type HandleSASLInitialResponse func(ctx *Context, msg *message.SASLInitialResponse) (*message.SASLInitialResponse, error)
type HandleSASLResponse func(ctx *Context, msg *message.SASLResponse) (*message.SASLResponse, error)
type HandleQuery func(ctx *Context, msg *message.Query) (*message.Query, error)
type HandleSync func(ctx *Context, msg *message.Sync) (*message.Sync, error)
type HandleTerminate func(ctx *Context, msg *message.Terminate) (*message.Terminate, error)
type HandleCopyData func(ctx *Context, msg *message.CopyData) (*message.CopyData, error)
type HandleCopyDone func(ctx *Context, msg *message.CopyDone) (*message.CopyDone, error)

type HandlePMessage struct {
	HandlePasswordMessage     HandlePasswordMessage
	HandleGSSResponse         HandleGSSResponse
	HandleSASLInitialResponse HandleSASLInitialResponse
	HandleSASLResponse        HandleSASLResponse
}

type FrontendMessageHandlers map[byte]MessageHandler

func (s *FrontendMessageHandlers) SetHandleBind(h HandleBind) {
	if h == nil {
		return
	}
	(*s)['B'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadBind(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleClose(h HandleClose) {
	if h == nil {
		return
	}
	(*s)['C'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadClose(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleCopyFail(h HandleCopyFail) {
	if h == nil {
		return
	}
	(*s)['f'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyFail(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleDescribe(h HandleDescribe) {
	if h == nil {
		return
	}
	(*s)['D'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadDescribe(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleExecute(h HandleExecute) {
	if h == nil {
		return
	}
	(*s)['E'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadExecute(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleFlush(h HandleFlush) {
	if h == nil {
		return
	}
	(*s)['H'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadFlush(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleFunctionCall(h HandleFunctionCall) {
	if h == nil {
		return
	}
	(*s)['F'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadFunctionCall(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleParse(h HandleParse) {
	if h == nil {
		return
	}
	(*s)['P'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadParse(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandlePMessage(h HandlePMessage) {
	(*s)['p'] = func(ctx *Context, raw []byte) (message.Reader, error) {
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
func (s *FrontendMessageHandlers) SetHandleQuery(h HandleQuery) {
	if h == nil {
		return
	}
	(*s)['Q'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadQuery(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleSync(h HandleSync) {
	if h == nil {
		return
	}
	(*s)['S'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadSync(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleTerminate(h HandleTerminate) {
	if h == nil {
		return
	}
	(*s)['X'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadTerminate(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleCopyData(h HandleCopyData) {
	if h == nil {
		return
	}
	(*s)['d'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyData(raw))
	}
}
func (s *FrontendMessageHandlers) SetHandleCopyDone(h HandleCopyDone) {
	if h == nil {
		return
	}
	(*s)['c'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyDone(raw))
	}
}

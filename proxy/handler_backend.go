package proxy

import (
	"errors"

	"github.com/rueian/pgbroker/message"
)

type HandleAuthenticationOk func(ctx *Context, msg *message.AuthenticationOk) (*message.AuthenticationOk, error)
type HandleAuthenticationKerberosV5 func(ctx *Context, msg *message.AuthenticationKerberosV5) (*message.AuthenticationKerberosV5, error)
type HandleAuthenticationCleartextPassword func(ctx *Context, msg *message.AuthenticationCleartextPassword) (*message.AuthenticationCleartextPassword, error)
type HandleAuthenticationMD5Password func(ctx *Context, msg *message.AuthenticationMD5Password) (*message.AuthenticationMD5Password, error)
type HandleAuthenticationSCMCredential func(ctx *Context, msg *message.AuthenticationSCMCredential) (*message.AuthenticationSCMCredential, error)
type HandleAuthenticationGSS func(ctx *Context, msg *message.AuthenticationGSS) (*message.AuthenticationGSS, error)
type HandleAuthenticationSSPI func(ctx *Context, msg *message.AuthenticationSSPI) (*message.AuthenticationSSPI, error)
type HandleAuthenticationGSSContinue func(ctx *Context, msg *message.AuthenticationGSSContinue) (*message.AuthenticationGSSContinue, error)
type HandleAuthenticationSASL func(ctx *Context, msg *message.AuthenticationSASL) (*message.AuthenticationSASL, error)
type HandleAuthenticationSASLContinue func(ctx *Context, msg *message.AuthenticationSASLContinue) (*message.AuthenticationSASLContinue, error)
type HandleAuthenticationSASLFinal func(ctx *Context, msg *message.AuthenticationSASLFinal) (*message.AuthenticationSASLFinal, error)
type HandleBackendKeyData func(ctx *Context, msg *message.BackendKeyData) (*message.BackendKeyData, error)
type HandleBindComplete func(ctx *Context, msg *message.BindComplete) (*message.BindComplete, error)
type HandleCloseComplete func(ctx *Context, msg *message.CloseComplete) (*message.CloseComplete, error)
type HandleCommandComplete func(ctx *Context, msg *message.CommandComplete) (*message.CommandComplete, error)
type HandleCopyInResponse func(ctx *Context, msg *message.CopyInResponse) (*message.CopyInResponse, error)
type HandleCopyOutResponse func(ctx *Context, msg *message.CopyOutResponse) (*message.CopyOutResponse, error)
type HandleCopyBothResponse func(ctx *Context, msg *message.CopyBothResponse) (*message.CopyBothResponse, error)
type HandleDataRow func(ctx *Context, msg *message.DataRow) (*message.DataRow, error)
type HandleEmptyQueryResponse func(ctx *Context, msg *message.EmptyQueryResponse) (*message.EmptyQueryResponse, error)
type HandleErrorResponse func(ctx *Context, msg *message.ErrorResponse) (*message.ErrorResponse, error)
type HandleFunctionCallResponse func(ctx *Context, msg *message.FunctionCallResponse) (*message.FunctionCallResponse, error)
type HandleNegotiateProtocolVersion func(ctx *Context, msg *message.NegotiateProtocolVersion) (*message.NegotiateProtocolVersion, error)
type HandleNoData func(ctx *Context, msg *message.NoData) (*message.NoData, error)
type HandleNoticeResponse func(ctx *Context, msg *message.NoticeResponse) (*message.NoticeResponse, error)
type HandleNotificationResponse func(ctx *Context, msg *message.NotificationResponse) (*message.NotificationResponse, error)
type HandleParameterDescription func(ctx *Context, msg *message.ParameterDescription) (*message.ParameterDescription, error)
type HandleParameterStatus func(ctx *Context, msg *message.ParameterStatus) (*message.ParameterStatus, error)
type HandleParseComplete func(ctx *Context, msg *message.ParseComplete) (*message.ParseComplete, error)
type HandlePortalSuspended func(ctx *Context, msg *message.PortalSuspended) (*message.PortalSuspended, error)
type HandleReadyForQuery func(ctx *Context, msg *message.ReadyForQuery) (*message.ReadyForQuery, error)
type HandleRowDescription func(ctx *Context, msg *message.RowDescription) (*message.RowDescription, error)

type HandleAuthentication struct {
	HandleAuthenticationOk                HandleAuthenticationOk
	HandleAuthenticationKerberosV5        HandleAuthenticationKerberosV5
	HandleAuthenticationCleartextPassword HandleAuthenticationCleartextPassword
	HandleAuthenticationMD5Password       HandleAuthenticationMD5Password
	HandleAuthenticationSCMCredential     HandleAuthenticationSCMCredential
	HandleAuthenticationGSS               HandleAuthenticationGSS
	HandleAuthenticationSSPI              HandleAuthenticationSSPI
	HandleAuthenticationGSSContinue       HandleAuthenticationGSSContinue
	HandleAuthenticationSASL              HandleAuthenticationSASL
	HandleAuthenticationSASLContinue      HandleAuthenticationSASLContinue
	HandleAuthenticationSASLFinal         HandleAuthenticationSASLFinal
}

type BackendMessageHandlers map[byte]MessageHandler

func (s *BackendMessageHandlers) SetHandleAuthentication(h HandleAuthentication) {
	(*s)['R'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		switch m := message.ReadAuthentication(raw).(type) {
		case *message.AuthenticationOk:
			ctx.AuthPhase = PhaseOK
			if h.HandleAuthenticationOk == nil {
				return m, nil
			}
			return h.HandleAuthenticationOk(ctx, m)
		case *message.AuthenticationKerberosV5:
			if h.HandleAuthenticationKerberosV5 == nil {
				return m, nil
			}
			return h.HandleAuthenticationKerberosV5(ctx, m)
		case *message.AuthenticationCleartextPassword:
			if h.HandleAuthenticationCleartextPassword == nil {
				return m, nil
			}
			return h.HandleAuthenticationCleartextPassword(ctx, m)
		case *message.AuthenticationMD5Password:
			if h.HandleAuthenticationMD5Password == nil {
				return m, nil
			}
			return h.HandleAuthenticationMD5Password(ctx, m)
		case *message.AuthenticationSCMCredential:
			if h.HandleAuthenticationSCMCredential == nil {
				return m, nil
			}
			return h.HandleAuthenticationSCMCredential(ctx, m)
		case *message.AuthenticationGSS:
			ctx.AuthPhase = PhaseGSS
			if h.HandleAuthenticationGSS == nil {
				return m, nil
			}
			return h.HandleAuthenticationGSS(ctx, m)
		case *message.AuthenticationSSPI:
			if h.HandleAuthenticationSSPI == nil {
				return m, nil
			}
			return h.HandleAuthenticationSSPI(ctx, m)
		case *message.AuthenticationGSSContinue:
			if h.HandleAuthenticationGSSContinue == nil {
				return m, nil
			}
			return h.HandleAuthenticationGSSContinue(ctx, m)
		case *message.AuthenticationSASL:
			ctx.AuthPhase = PhaseSASLInit
			if h.HandleAuthenticationSASL == nil {
				return m, nil
			}
			return h.HandleAuthenticationSASL(ctx, m)
		case *message.AuthenticationSASLContinue:
			ctx.AuthPhase = PhaseSASL
			if h.HandleAuthenticationSASLContinue == nil {
				return m, nil
			}
			return h.HandleAuthenticationSASLContinue(ctx, m)
		case *message.AuthenticationSASLFinal:
			if h.HandleAuthenticationSASLFinal == nil {
				return m, nil
			}
			return h.HandleAuthenticationSASLFinal(ctx, m)
		}
		return nil, errors.New("fail to cast authentication message")
	}
}
func (s *BackendMessageHandlers) SetHandleBackendKeyData(h HandleBackendKeyData) {
	if h == nil {
		return
	}
	(*s)['K'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadBackendKeyData(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleBindComplete(h HandleBindComplete) {
	if h == nil {
		return
	}
	(*s)['2'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadBindComplete(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleCloseComplete(h HandleCloseComplete) {
	if h == nil {
		return
	}
	(*s)['3'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCloseComplete(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleCommandComplete(h HandleCommandComplete) {
	if h == nil {
		return
	}
	(*s)['C'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCommandComplete(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleCopyInResponse(h HandleCopyInResponse) {
	if h == nil {
		return
	}
	(*s)['G'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyInResponse(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleCopyOutResponse(h HandleCopyOutResponse) {
	if h == nil {
		return
	}
	(*s)['H'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyOutResponse(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleCopyBothResponse(h HandleCopyBothResponse) {
	if h == nil {
		return
	}
	(*s)['W'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyBothResponse(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleDataRow(h HandleDataRow) {
	if h == nil {
		return
	}
	(*s)['D'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadDataRow(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleEmptyQueryResponse(h HandleEmptyQueryResponse) {
	if h == nil {
		return
	}
	(*s)['I'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadEmptyQueryResponse(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleErrorResponse(h HandleErrorResponse) {
	if h == nil {
		return
	}
	(*s)['E'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadErrorResponse(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleFunctionCallResponse(h HandleFunctionCallResponse) {
	if h == nil {
		return
	}
	(*s)['V'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadFunctionCallResponse(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleNegotiateProtocolVersion(h HandleNegotiateProtocolVersion) {
	if h == nil {
		return
	}
	(*s)['v'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadNegotiateProtocolVersion(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleNoData(h HandleNoData) {
	if h == nil {
		return
	}
	(*s)['n'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadNoData(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleNoticeResponse(h HandleNoticeResponse) {
	if h == nil {
		return
	}
	(*s)['N'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadNoticeResponse(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleNotificationResponse(h HandleNotificationResponse) {
	if h == nil {
		return
	}
	(*s)['A'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadNotificationResponse(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleParameterDescription(h HandleParameterDescription) {
	if h == nil {
		return
	}
	(*s)['t'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadParameterDescription(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleParameterStatus(h HandleParameterStatus) {
	if h == nil {
		return
	}
	(*s)['S'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadParameterStatus(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleParseComplete(h HandleParseComplete) {
	if h == nil {
		return
	}
	(*s)['1'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadParseComplete(raw))
	}
}
func (s *BackendMessageHandlers) SetHandlePortalSuspended(h HandlePortalSuspended) {
	if h == nil {
		return
	}
	(*s)['s'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadPortalSuspended(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleReadyForQuery(h HandleReadyForQuery) {
	if h == nil {
		return
	}
	(*s)['Z'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadReadyForQuery(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleRowDescription(h HandleRowDescription) {
	if h == nil {
		return
	}
	(*s)['T'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadRowDescription(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleCopyData(h HandleCopyData) {
	if h == nil {
		return
	}
	(*s)['d'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyData(raw))
	}
}
func (s *BackendMessageHandlers) SetHandleCopyDone(h HandleCopyDone) {
	if h == nil {
		return
	}
	(*s)['c'] = func(ctx *Context, raw []byte) (message.Reader, error) {
		return h(ctx, message.ReadCopyDone(raw))
	}
}

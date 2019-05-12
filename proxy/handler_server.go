package proxy

import (
	"errors"

	"github.com/rueian/pgbroker/message"
)

type HandleAuthenticationOk func(ms *Ctx, msg *message.AuthenticationOk) (*message.AuthenticationOk, error)
type HandleAuthenticationKerberosV5 func(ms *Ctx, msg *message.AuthenticationKerberosV5) (*message.AuthenticationKerberosV5, error)
type HandleAuthenticationCleartextPassword func(ms *Ctx, msg *message.AuthenticationCleartextPassword) (*message.AuthenticationCleartextPassword, error)
type HandleAuthenticationMD5Password func(ms *Ctx, msg *message.AuthenticationMD5Password) (*message.AuthenticationMD5Password, error)
type HandleAuthenticationSCMCredential func(ms *Ctx, msg *message.AuthenticationSCMCredential) (*message.AuthenticationSCMCredential, error)
type HandleAuthenticationGSS func(ms *Ctx, msg *message.AuthenticationGSS) (*message.AuthenticationGSS, error)
type HandleAuthenticationSSPI func(ms *Ctx, msg *message.AuthenticationSSPI) (*message.AuthenticationSSPI, error)
type HandleAuthenticationGSSContinue func(ms *Ctx, msg *message.AuthenticationGSSContinue) (*message.AuthenticationGSSContinue, error)
type HandleAuthenticationSASL func(ms *Ctx, msg *message.AuthenticationSASL) (*message.AuthenticationSASL, error)
type HandleAuthenticationSASLContinue func(ms *Ctx, msg *message.AuthenticationSASLContinue) (*message.AuthenticationSASLContinue, error)
type HandleAuthenticationSASLFinal func(ms *Ctx, msg *message.AuthenticationSASLFinal) (*message.AuthenticationSASLFinal, error)
type HandleBackendKeyData func(ms *Ctx, msg *message.BackendKeyData) (*message.BackendKeyData, error)
type HandleBindComplete func(ms *Ctx, msg *message.BindComplete) (*message.BindComplete, error)
type HandleCloseComplete func(ms *Ctx, msg *message.CloseComplete) (*message.CloseComplete, error)
type HandleCommandComplete func(ms *Ctx, msg *message.CommandComplete) (*message.CommandComplete, error)
type HandleCopyInResponse func(ms *Ctx, msg *message.CopyInResponse) (*message.CopyInResponse, error)
type HandleCopyOutResponse func(ms *Ctx, msg *message.CopyOutResponse) (*message.CopyOutResponse, error)
type HandleCopyBothResponse func(ms *Ctx, msg *message.CopyBothResponse) (*message.CopyBothResponse, error)
type HandleDataRow func(ms *Ctx, msg *message.DataRow) (*message.DataRow, error)
type HandleEmptyQueryResponse func(ms *Ctx, msg *message.EmptyQueryResponse) (*message.EmptyQueryResponse, error)
type HandleErrorResponse func(ms *Ctx, msg *message.ErrorResponse) (*message.ErrorResponse, error)
type HandleFunctionCallResponse func(ms *Ctx, msg *message.FunctionCallResponse) (*message.FunctionCallResponse, error)
type HandleNegotiateProtocolVersion func(ms *Ctx, msg *message.NegotiateProtocolVersion) (*message.NegotiateProtocolVersion, error)
type HandleNoData func(ms *Ctx, msg *message.NoData) (*message.NoData, error)
type HandleNoticeResponse func(ms *Ctx, msg *message.NoticeResponse) (*message.NoticeResponse, error)
type HandleNotificationResponse func(ms *Ctx, msg *message.NotificationResponse) (*message.NotificationResponse, error)
type HandleParameterDescription func(ms *Ctx, msg *message.ParameterDescription) (*message.ParameterDescription, error)
type HandleParameterStatus func(ms *Ctx, msg *message.ParameterStatus) (*message.ParameterStatus, error)
type HandleParseComplete func(ms *Ctx, msg *message.ParseComplete) (*message.ParseComplete, error)
type HandlePortalSuspended func(ms *Ctx, msg *message.PortalSuspended) (*message.PortalSuspended, error)
type HandleReadyForQuery func(ms *Ctx, msg *message.ReadyForQuery) (*message.ReadyForQuery, error)
type HandleRowDescription func(ms *Ctx, msg *message.RowDescription) (*message.RowDescription, error)

type ServerMessageHandlers struct {
	m                                     map[byte]MessageHandler
	handleAuthenticationOk                []HandleAuthenticationOk
	handleAuthenticationKerberosV5        []HandleAuthenticationKerberosV5
	handleAuthenticationCleartextPassword []HandleAuthenticationCleartextPassword
	handleAuthenticationMD5Password       []HandleAuthenticationMD5Password
	handleAuthenticationSCMCredential     []HandleAuthenticationSCMCredential
	handleAuthenticationGSS               []HandleAuthenticationGSS
	handleAuthenticationSSPI              []HandleAuthenticationSSPI
	handleAuthenticationGSSContinue       []HandleAuthenticationGSSContinue
	handleAuthenticationSASL              []HandleAuthenticationSASL
	handleAuthenticationSASLContinue      []HandleAuthenticationSASLContinue
	handleAuthenticationSASLFinal         []HandleAuthenticationSASLFinal
	handleBackendKeyData                  []HandleBackendKeyData
	handleBindComplete                    []HandleBindComplete
	handleCloseComplete                   []HandleCloseComplete
	handleCommandComplete                 []HandleCommandComplete
	handleCopyInResponse                  []HandleCopyInResponse
	handleCopyOutResponse                 []HandleCopyOutResponse
	handleCopyBothResponse                []HandleCopyBothResponse
	handleDataRow                         []HandleDataRow
	handleEmptyQueryResponse              []HandleEmptyQueryResponse
	handleErrorResponse                   []HandleErrorResponse
	handleFunctionCallResponse            []HandleFunctionCallResponse
	handleNegotiateProtocolVersion        []HandleNegotiateProtocolVersion
	handleNoData                          []HandleNoData
	handleNoticeResponse                  []HandleNoticeResponse
	handleNotificationResponse            []HandleNotificationResponse
	handleParameterDescription            []HandleParameterDescription
	handleParameterStatus                 []HandleParameterStatus
	handleParseComplete                   []HandleParseComplete
	handlePortalSuspended                 []HandlePortalSuspended
	handleReadyForQuery                   []HandleReadyForQuery
	handleRowDescription                  []HandleRowDescription
	handleCopyData                        []HandleCopyData
	handleCopyDone                        []HandleCopyDone
}

func NewServerMessageHandlers() *ServerMessageHandlers {
	return &ServerMessageHandlers{m: make(map[byte]MessageHandler)}
}

func (s *ServerMessageHandlers) GetHandler(b byte) MessageHandler {
	if handler, ok := s.m[b]; ok {
		return handler
	}
	switch b {
	case 'R':
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			switch msg := message.ReadAuthentication(raw).(type) {
			case *message.AuthenticationOk:
				md.AuthPhase = PhaseOK
				for _, h := range s.handleAuthenticationOk {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationKerberosV5:
				for _, h := range s.handleAuthenticationKerberosV5 {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationCleartextPassword:
				for _, h := range s.handleAuthenticationCleartextPassword {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationMD5Password:
				for _, h := range s.handleAuthenticationMD5Password {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationSCMCredential:
				for _, h := range s.handleAuthenticationSCMCredential {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationGSS:
				md.AuthPhase = PhaseGSS
				for _, h := range s.handleAuthenticationGSS {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationSSPI:
				for _, h := range s.handleAuthenticationSSPI {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationGSSContinue:
				for _, h := range s.handleAuthenticationGSSContinue {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationSASL:
				md.AuthPhase = PhaseSASLInit
				for _, h := range s.handleAuthenticationSASL {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationSASLContinue:
				md.AuthPhase = PhaseSASL
				for _, h := range s.handleAuthenticationSASLContinue {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case *message.AuthenticationSASLFinal:
				for _, h := range s.handleAuthenticationSASLFinal {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			}
			return nil, errors.New("fail to cast authentication message")
		}
		return s.m[b]
	case 'K':
		if len(s.handleBackendKeyData) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadBackendKeyData(raw)
			for _, h := range s.handleBackendKeyData {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case '2':
		if len(s.handleBindComplete) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadBindComplete(raw)
			for _, h := range s.handleBindComplete {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case '3':
		if len(s.handleCloseComplete) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadCloseComplete(raw)
			for _, h := range s.handleCloseComplete {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'C':
		if len(s.handleCommandComplete) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadCommandComplete(raw)
			for _, h := range s.handleCommandComplete {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'G':
		if len(s.handleCopyInResponse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadCopyInResponse(raw)
			for _, h := range s.handleCopyInResponse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'H':
		if len(s.handleCopyOutResponse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadCopyOutResponse(raw)
			for _, h := range s.handleCopyOutResponse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'W':
		if len(s.handleCopyBothResponse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadCopyBothResponse(raw)
			for _, h := range s.handleCopyBothResponse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'D':
		if len(s.handleDataRow) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadDataRow(raw)
			for _, h := range s.handleDataRow {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'I':
		if len(s.handleEmptyQueryResponse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadEmptyQueryResponse(raw)
			for _, h := range s.handleEmptyQueryResponse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'E':
		if len(s.handleErrorResponse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadErrorResponse(raw)
			for _, h := range s.handleErrorResponse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'V':
		if len(s.handleFunctionCallResponse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadFunctionCallResponse(raw)
			for _, h := range s.handleFunctionCallResponse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'v':
		if len(s.handleNegotiateProtocolVersion) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadNegotiateProtocolVersion(raw)
			for _, h := range s.handleNegotiateProtocolVersion {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'n':
		if len(s.handleNoData) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadNoData(raw)
			for _, h := range s.handleNoData {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'N':
		if len(s.handleNoticeResponse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadNoticeResponse(raw)
			for _, h := range s.handleNoticeResponse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'A':
		if len(s.handleNotificationResponse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadNotificationResponse(raw)
			for _, h := range s.handleNotificationResponse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 't':
		if len(s.handleParameterDescription) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadParameterDescription(raw)
			for _, h := range s.handleParameterDescription {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'S':
		if len(s.handleParameterStatus) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadParameterStatus(raw)
			for _, h := range s.handleParameterStatus {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case '1':
		if len(s.handleParseComplete) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadParseComplete(raw)
			for _, h := range s.handleParseComplete {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 's':
		if len(s.handlePortalSuspended) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadPortalSuspended(raw)
			for _, h := range s.handlePortalSuspended {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'Z':
		if len(s.handleReadyForQuery) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadReadyForQuery(raw)
			for _, h := range s.handleReadyForQuery {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'T':
		if len(s.handleRowDescription) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadRowDescription(raw)
			for _, h := range s.handleRowDescription {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'd':
		if len(s.handleCopyData) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadCopyData(raw)
			for _, h := range s.handleCopyData {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'c':
		if len(s.handleCopyDone) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadCopyDone(raw)
			for _, h := range s.handleCopyDone {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	}
	return nil
}

func (s *ServerMessageHandlers) AddHandleAuthenticationOk(h HandleAuthenticationOk) {
	if h == nil {
		return
	}
	s.handleAuthenticationOk = append(s.handleAuthenticationOk, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationKerberosV5(h HandleAuthenticationKerberosV5) {
	if h == nil {
		return
	}
	s.handleAuthenticationKerberosV5 = append(s.handleAuthenticationKerberosV5, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationCleartextPassword(h HandleAuthenticationCleartextPassword) {
	if h == nil {
		return
	}
	s.handleAuthenticationCleartextPassword = append(s.handleAuthenticationCleartextPassword, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationMD5Password(h HandleAuthenticationMD5Password) {
	if h == nil {
		return
	}
	s.handleAuthenticationMD5Password = append(s.handleAuthenticationMD5Password, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationSCMCredential(h HandleAuthenticationSCMCredential) {
	if h == nil {
		return
	}
	s.handleAuthenticationSCMCredential = append(s.handleAuthenticationSCMCredential, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationGSS(h HandleAuthenticationGSS) {
	if h == nil {
		return
	}
	s.handleAuthenticationGSS = append(s.handleAuthenticationGSS, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationSSPI(h HandleAuthenticationSSPI) {
	if h == nil {
		return
	}
	s.handleAuthenticationSSPI = append(s.handleAuthenticationSSPI, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationGSSContinue(h HandleAuthenticationGSSContinue) {
	if h == nil {
		return
	}
	s.handleAuthenticationGSSContinue = append(s.handleAuthenticationGSSContinue, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationSASL(h HandleAuthenticationSASL) {
	if h == nil {
		return
	}
	s.handleAuthenticationSASL = append(s.handleAuthenticationSASL, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationSASLContinue(h HandleAuthenticationSASLContinue) {
	if h == nil {
		return
	}
	s.handleAuthenticationSASLContinue = append(s.handleAuthenticationSASLContinue, h)
}
func (s *ServerMessageHandlers) AddHandleAuthenticationSASLFinal(h HandleAuthenticationSASLFinal) {
	if h == nil {
		return
	}
	s.handleAuthenticationSASLFinal = append(s.handleAuthenticationSASLFinal, h)
}
func (s *ServerMessageHandlers) AddHandleBackendKeyData(h HandleBackendKeyData) {
	if h == nil {
		return
	}
	s.handleBackendKeyData = append(s.handleBackendKeyData, h)
}
func (s *ServerMessageHandlers) AddHandleBindComplete(h HandleBindComplete) {
	if h == nil {
		return
	}
	s.handleBindComplete = append(s.handleBindComplete, h)
}
func (s *ServerMessageHandlers) AddHandleCloseComplete(h HandleCloseComplete) {
	if h == nil {
		return
	}
	s.handleCloseComplete = append(s.handleCloseComplete, h)
}
func (s *ServerMessageHandlers) AddHandleCommandComplete(h HandleCommandComplete) {
	if h == nil {
		return
	}
	s.handleCommandComplete = append(s.handleCommandComplete, h)
}
func (s *ServerMessageHandlers) AddHandleCopyInResponse(h HandleCopyInResponse) {
	if h == nil {
		return
	}
	s.handleCopyInResponse = append(s.handleCopyInResponse, h)
}
func (s *ServerMessageHandlers) AddHandleCopyOutResponse(h HandleCopyOutResponse) {
	if h == nil {
		return
	}
	s.handleCopyOutResponse = append(s.handleCopyOutResponse, h)
}
func (s *ServerMessageHandlers) AddHandleCopyBothResponse(h HandleCopyBothResponse) {
	if h == nil {
		return
	}
	s.handleCopyBothResponse = append(s.handleCopyBothResponse, h)
}
func (s *ServerMessageHandlers) AddHandleDataRow(h HandleDataRow) {
	if h == nil {
		return
	}
	s.handleDataRow = append(s.handleDataRow, h)
}
func (s *ServerMessageHandlers) AddHandleEmptyQueryResponse(h HandleEmptyQueryResponse) {
	if h == nil {
		return
	}
	s.handleEmptyQueryResponse = append(s.handleEmptyQueryResponse, h)
}
func (s *ServerMessageHandlers) AddHandleErrorResponse(h HandleErrorResponse) {
	if h == nil {
		return
	}
	s.handleErrorResponse = append(s.handleErrorResponse, h)
}
func (s *ServerMessageHandlers) AddHandleFunctionCallResponse(h HandleFunctionCallResponse) {
	if h == nil {
		return
	}
	s.handleFunctionCallResponse = append(s.handleFunctionCallResponse, h)
}
func (s *ServerMessageHandlers) AddHandleNegotiateProtocolVersion(h HandleNegotiateProtocolVersion) {
	if h == nil {
		return
	}
	s.handleNegotiateProtocolVersion = append(s.handleNegotiateProtocolVersion, h)
}
func (s *ServerMessageHandlers) AddHandleNoData(h HandleNoData) {
	if h == nil {
		return
	}
	s.handleNoData = append(s.handleNoData, h)
}
func (s *ServerMessageHandlers) AddHandleNoticeResponse(h HandleNoticeResponse) {
	if h == nil {
		return
	}
	s.handleNoticeResponse = append(s.handleNoticeResponse, h)
}
func (s *ServerMessageHandlers) AddHandleNotificationResponse(h HandleNotificationResponse) {
	if h == nil {
		return
	}
	s.handleNotificationResponse = append(s.handleNotificationResponse, h)
}
func (s *ServerMessageHandlers) AddHandleParameterDescription(h HandleParameterDescription) {
	if h == nil {
		return
	}
	s.handleParameterDescription = append(s.handleParameterDescription, h)
}
func (s *ServerMessageHandlers) AddHandleParameterStatus(h HandleParameterStatus) {
	if h == nil {
		return
	}
	s.handleParameterStatus = append(s.handleParameterStatus, h)
}
func (s *ServerMessageHandlers) AddHandleParseComplete(h HandleParseComplete) {
	if h == nil {
		return
	}
	s.handleParseComplete = append(s.handleParseComplete, h)
}
func (s *ServerMessageHandlers) AddHandlePortalSuspended(h HandlePortalSuspended) {
	if h == nil {
		return
	}
	s.handlePortalSuspended = append(s.handlePortalSuspended, h)
}
func (s *ServerMessageHandlers) AddHandleReadyForQuery(h HandleReadyForQuery) {
	if h == nil {
		return
	}
	s.handleReadyForQuery = append(s.handleReadyForQuery, h)
}
func (s *ServerMessageHandlers) AddHandleRowDescription(h HandleRowDescription) {
	if h == nil {
		return
	}
	s.handleRowDescription = append(s.handleRowDescription, h)
}
func (s *ServerMessageHandlers) AddHandleCopyData(h HandleCopyData) {
	if h == nil {
		return
	}
	s.handleCopyData = append(s.handleCopyData, h)
}
func (s *ServerMessageHandlers) AddHandleCopyDone(h HandleCopyDone) {
	if h == nil {
		return
	}
	s.handleCopyDone = append(s.handleCopyDone, h)
}

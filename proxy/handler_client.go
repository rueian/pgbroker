package proxy

import (
	"github.com/rueian/pgbroker/message"
)

type HandleBind func(md *Ctx, msg *message.Bind) (*message.Bind, error)
type HandleClose func(md *Ctx, msg *message.Close) (*message.Close, error)
type HandleCopyFail func(md *Ctx, msg *message.CopyFail) (*message.CopyFail, error)
type HandleDescribe func(md *Ctx, msg *message.Describe) (*message.Describe, error)
type HandleExecute func(md *Ctx, msg *message.Execute) (*message.Execute, error)
type HandleFlush func(md *Ctx, msg *message.Flush) (*message.Flush, error)
type HandleFunctionCall func(md *Ctx, msg *message.FunctionCall) (*message.FunctionCall, error)
type HandleParse func(md *Ctx, msg *message.Parse) (*message.Parse, error)
type HandlePasswordMessage func(md *Ctx, msg *message.PasswordMessage) (*message.PasswordMessage, error)
type HandleGSSResponse func(md *Ctx, msg *message.GSSResponse) (*message.GSSResponse, error)
type HandleSASLInitialResponse func(md *Ctx, msg *message.SASLInitialResponse) (*message.SASLInitialResponse, error)
type HandleSASLResponse func(md *Ctx, msg *message.SASLResponse) (*message.SASLResponse, error)
type HandleQuery func(md *Ctx, msg *message.Query) (*message.Query, error)
type HandleSync func(md *Ctx, msg *message.Sync) (*message.Sync, error)
type HandleTerminate func(md *Ctx, msg *message.Terminate) (*message.Terminate, error)
type HandleCopyData func(md *Ctx, msg *message.CopyData) (*message.CopyData, error)
type HandleCopyDone func(md *Ctx, msg *message.CopyDone) (*message.CopyDone, error)

type ClientMessageHandlers struct {
	m                         map[byte]MessageHandler
	handleBind                []HandleBind
	handleClose               []HandleClose
	handleCopyFail            []HandleCopyFail
	handleDescribe            []HandleDescribe
	handleExecute             []HandleExecute
	handleFlush               []HandleFlush
	handleFunctionCall        []HandleFunctionCall
	handleParse               []HandleParse
	handlePasswordMessage     []HandlePasswordMessage
	handleGSSResponse         []HandleGSSResponse
	handleSASLInitialResponse []HandleSASLInitialResponse
	handleSASLResponse        []HandleSASLResponse
	handleQuery               []HandleQuery
	handleSync                []HandleSync
	handleTerminate           []HandleTerminate
	handleCopyData            []HandleCopyData
	handleCopyDone            []HandleCopyDone
}

func NewClientMessageHandlers() *ClientMessageHandlers {
	return &ClientMessageHandlers{m: make(map[byte]MessageHandler)}
}

func (s *ClientMessageHandlers) GetHandler(b byte) MessageHandler {
	if handler, ok := s.m[b]; ok {
		return handler
	}
	switch b {
	case 'B':
		if len(s.handleBind) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadBind(raw)
			for _, h := range s.handleBind {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'C':
		if len(s.handleClose) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadClose(raw)
			for _, h := range s.handleClose {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'f':
		if len(s.handleCopyFail) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadCopyFail(raw)
			for _, h := range s.handleCopyFail {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'D':
		if len(s.handleDescribe) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadDescribe(raw)
			for _, h := range s.handleDescribe {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'E':
		if len(s.handleExecute) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadExecute(raw)
			for _, h := range s.handleExecute {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'H':
		if len(s.handleFlush) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadFlush(raw)
			for _, h := range s.handleFlush {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'F':
		if len(s.handleFunctionCall) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadFunctionCall(raw)
			for _, h := range s.handleFunctionCall {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'P':
		if len(s.handleParse) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadParse(raw)
			for _, h := range s.handleParse {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'Q':
		if len(s.handleQuery) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadQuery(raw)
			for _, h := range s.handleQuery {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'S':
		if len(s.handleSync) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadSync(raw)
			for _, h := range s.handleSync {
				if msg, err = h(md, msg); err != nil {
					break
				}
			}
			return msg, err
		}
		return s.m[b]
	case 'X':
		if len(s.handleTerminate) == 0 {
			s.m[b] = nil
			break
		}
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			msg := message.ReadTerminate(raw)
			for _, h := range s.handleTerminate {
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
	case 'p':
		s.m[b] = func(md *Ctx, raw []byte) (message.Reader, error) {
			var err error
			switch md.AuthPhase {
			case PhaseSASLInit:
				msg := message.ReadSASLInitialResponse(raw)
				for _, h := range s.handleSASLInitialResponse {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case PhaseSASL:
				msg := message.ReadSASLResponse(raw)
				for _, h := range s.handleSASLResponse {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			case PhaseGSS:
				msg := message.ReadGSSResponse(raw)
				for _, h := range s.handleGSSResponse {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			default:
				msg := message.ReadPasswordMessage(raw)
				for _, h := range s.handlePasswordMessage {
					if msg, err = h(md, msg); err != nil {
						break
					}
				}
				return msg, err
			}
		}
		return s.m[b]
	}
	return nil
}

func (s *ClientMessageHandlers) AddHandleBind(h HandleBind) {
	if h == nil {
		return
	}
	s.handleBind = append(s.handleBind, h)
}
func (s *ClientMessageHandlers) AddHandleClose(h HandleClose) {
	if h == nil {
		return
	}
	s.handleClose = append(s.handleClose, h)
}
func (s *ClientMessageHandlers) AddHandleCopyFail(h HandleCopyFail) {
	if h == nil {
		return
	}
	s.handleCopyFail = append(s.handleCopyFail, h)
}
func (s *ClientMessageHandlers) AddHandleDescribe(h HandleDescribe) {
	if h == nil {
		return
	}
	s.handleDescribe = append(s.handleDescribe, h)
}
func (s *ClientMessageHandlers) AddHandleExecute(h HandleExecute) {
	if h == nil {
		return
	}
	s.handleExecute = append(s.handleExecute, h)
}
func (s *ClientMessageHandlers) AddHandleFlush(h HandleFlush) {
	if h == nil {
		return
	}
	s.handleFlush = append(s.handleFlush, h)
}
func (s *ClientMessageHandlers) AddHandleFunctionCall(h HandleFunctionCall) {
	if h == nil {
		return
	}
	s.handleFunctionCall = append(s.handleFunctionCall, h)
}
func (s *ClientMessageHandlers) AddHandleParse(h HandleParse) {
	if h == nil {
		return
	}
	s.handleParse = append(s.handleParse, h)
}
func (s *ClientMessageHandlers) AddHandleQuery(h HandleQuery) {
	if h == nil {
		return
	}
	s.handleQuery = append(s.handleQuery, h)
}
func (s *ClientMessageHandlers) AddHandleSync(h HandleSync) {
	if h == nil {
		return
	}
	s.handleSync = append(s.handleSync, h)
}
func (s *ClientMessageHandlers) AddHandleTerminate(h HandleTerminate) {
	if h == nil {
		return
	}
	s.handleTerminate = append(s.handleTerminate, h)
}
func (s *ClientMessageHandlers) AddHandleCopyData(h HandleCopyData) {
	if h == nil {
		return
	}
	s.handleCopyData = append(s.handleCopyData, h)
}
func (s *ClientMessageHandlers) AddHandleCopyDone(h HandleCopyDone) {
	if h == nil {
		return
	}
	s.handleCopyDone = append(s.handleCopyDone, h)
}
func (s *ClientMessageHandlers) AddHandleSASLInitialResponse(h HandleSASLInitialResponse) {
	if h == nil {
		return
	}
	s.handleSASLInitialResponse = append(s.handleSASLInitialResponse, h)
}
func (s *ClientMessageHandlers) AddHandleSASLResponse(h HandleSASLResponse) {
	if h == nil {
		return
	}
	s.handleSASLResponse = append(s.handleSASLResponse, h)
}
func (s *ClientMessageHandlers) AddHandleGSSResponse(h HandleGSSResponse) {
	if h == nil {
		return
	}
	s.handleGSSResponse = append(s.handleGSSResponse, h)
}
func (s *ClientMessageHandlers) AddHandlePasswordMessage(h HandlePasswordMessage) {
	if h == nil {
		return
	}
	s.handlePasswordMessage = append(s.handlePasswordMessage, h)
}

package proxy

import "io"

func NewStreamMessageHandler() *StreamMessageHandlers {
	return &StreamMessageHandlers{m: make(map[byte]StreamHandler)}
}

func NewStreamMessageHandler2() *StreamMessageHandlers2 {
	return &StreamMessageHandlers2{m: make(map[byte]StreamHandler2)}
}

type StreamMessageHandlers struct {
	m map[byte]StreamHandler
}

func (s *StreamMessageHandlers) GetHandler(b byte) StreamHandler {
	return s.m[b]
}

func (s *StreamMessageHandlers) AddHandler(b byte, handler StreamHandler) {
	s.m[b] = handler
}

func DefaultStreamHandler(ctx *Ctx, count int, pi io.Reader, po io.Writer) {
	io.Copy(po, pi)
}

type StreamMessageHandlers2 struct {
	m map[byte]StreamHandler2
}

func (s *StreamMessageHandlers2) GetHandler(b byte) StreamHandler2 {
	if s.m[b] == nil {
		return DefaultStreamHandler2
	}
	return s.m[b]
}

func (s *StreamMessageHandlers2) AddHandler(b byte, handler StreamHandler2) {
	s.m[b] = handler
}

func DefaultStreamHandler2(ctx *Ctx, pi *sliceChanReader, po *sliceChanWriter) {
	for {
		s, err := pi.NextSlice()
		if err != nil {
			return
		}
		po.Write(s)
	}
}

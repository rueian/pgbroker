package proxy

type StreamCallback func(slice Slice) Slice

type StreamCallbackFactory func(ctx *Ctx) StreamCallback

func NewStreamCallbackFactories() *StreamCallbackFactories {
	return &StreamCallbackFactories{m: make(map[byte]StreamCallbackFactory)}
}

type StreamCallbackFactories struct {
	m map[byte]StreamCallbackFactory
}

func (s *StreamCallbackFactories) GetFactory(b byte) StreamCallbackFactory {
	if s.m[b] == nil {
		return DefaultStreamCallbackFactory
	}
	return s.m[b]
}

func (s *StreamCallbackFactories) SetFactory(b byte, factory StreamCallbackFactory) {
	s.m[b] = factory
}

func DefaultStreamCallbackFactory(ctx *Ctx) StreamCallback {
	return DefaultStreamCallback
}

func DefaultStreamCallback(s Slice) Slice {
	return s
}

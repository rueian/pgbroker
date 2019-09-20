package proxy

func NewStreamMessageHandler() *StreamMessageHandlers {
	return &StreamMessageHandlers{m: make(map[byte]StreamHandler)}
}

type StreamMessageHandlers struct {
	m map[byte]StreamHandler
}

func (s *StreamMessageHandlers) GetHandler(b byte) StreamHandler {
	if s.m[b] == nil {
		return DefaultStreamHandler
	}
	return s.m[b]
}

func (s *StreamMessageHandlers) AddHandler(b byte, handler StreamHandler) {
	s.m[b] = handler
}

func DefaultStreamHandler(ctx *Ctx, pi *SliceChanReader, po *SliceChanWriter) {
	for {
		s, err := pi.NextSlice()
		if err != nil {
			return
		}
		po.Write(s)
	}
}

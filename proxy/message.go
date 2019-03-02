package proxy

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/rueian/pgbroker/message"
)

const (
	MessageTypeLength     = 1
	MessageSizeLength     = 4
	InitMessageSizeLength = 4
)

type MessageHandler func(*Context, []byte) (message.Reader, error)

type container struct {
	head []byte
	data []byte
}

func (c container) Reader() io.Reader {
	return io.MultiReader(bytes.NewReader(c.head), bytes.NewReader(c.data))
}

func (c container) GetType() byte {
	return c.head[0]
}

func (c container) GetSize() uint32 {
	return binary.BigEndian.Uint32(c.head[len(c.head)-4:])
}

func ReadInitMessage(client io.Reader) (message.Reader, error) {
	c := container{}

	c.head = make([]byte, InitMessageSizeLength)
	if _, err := io.ReadFull(client, c.head); err != nil {
		return c, err
	}

	c.data = make([]byte, c.GetSize()-InitMessageSizeLength)
	if _, err := io.ReadFull(client, c.data); err != nil {
		return c, err
	}

	m := message.ReadStartupMessage(c.data)

	return m, nil
}

func ProcessMessages(ctx *Context, in io.Reader, out io.Writer, handlers map[byte]MessageHandler) error {
	for {
		c := container{}

		c.head = make([]byte, MessageTypeLength+MessageSizeLength)
		if _, err := io.ReadFull(in, c.head); err != nil {
			return err
		}

		if handler, ok := handlers[c.GetType()]; ok {
			c.data = make([]byte, c.GetSize()-MessageSizeLength)
			if _, err := io.ReadFull(in, c.data); err != nil {
				return err
			}

			if msg, err := handler(ctx, c.data); err != nil {
				if _, err := io.Copy(out, c.Reader()); err != nil {
					return err
				}
			} else {
				if _, err := io.Copy(out, msg.Reader()); err != nil {
					return err
				}
			}
		} else {
			c.data = make([]byte, 0)
			if _, err := io.Copy(out, c.Reader()); err != nil {
				return err
			}
			if _, err := io.CopyN(out, in, int64(c.GetSize()-MessageSizeLength)); err != nil {
				return err
			}
		}
	}
}

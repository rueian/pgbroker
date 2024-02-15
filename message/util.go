package message

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Reader interface {
	Reader() io.Reader
}

type Message struct {
	data   []byte
	offset uint32
	t      []byte
}

func (b *Message) ReadByte() (r byte) {
	r = b.data[b.offset]
	b.offset++
	return r
}

func (b *Message) ReadByteN(n uint32) (r []byte) {
	r = b.data[b.offset : b.offset+n]
	b.offset += n
	return r
}

func (b *Message) ReadUint8() (r uint8) {
	r = uint8(b.data[b.offset])
	b.offset++
	return r
}

func (b *Message) ReadUint16() (r uint16) {
	r = binary.BigEndian.Uint16(b.data[b.offset : b.offset+2])
	b.offset += 2
	return r
}

func (b *Message) ReadUint32() (r uint32) {
	r = binary.BigEndian.Uint32(b.data[b.offset : b.offset+4])
	b.offset += 4
	return r
}

func (b *Message) ReadString() (r string) {
	end := b.offset
	max := uint32(len(b.data))
	for ; end != max && b.data[end] != 0; end++ {
	}
	r = string(b.data[b.offset:end])
	b.offset = end + 1
	return r
}

func (b *Message) WriteByte(i byte) {
	b.data[b.offset] = i
	b.offset++
}

func (b *Message) WriteByteN(i []byte) {
	for _, s := range i {
		b.WriteByte(s)
	}
}

func (b *Message) WriteUint8(i uint8) {
	b.WriteByte(i)
}

func (b *Message) WriteUint16(i uint16) {
	binary.BigEndian.PutUint16(b.data[b.offset:b.offset+2], i)
	b.offset += 2
}

func (b *Message) WriteUint32(i uint32) {
	binary.BigEndian.PutUint32(b.data[b.offset:b.offset+4], i)
	b.offset += 4
}

func (b *Message) WriteString(i string) {
	b.WriteByteN([]byte(i))
	b.WriteByte(0)
}

func (b *Message) WriteValue(i Value) {
	b.WriteUint32(i.DataLength())
	b.WriteByteN(i.DataBytes())
}

func (b *Message) Length() int {
	return len(b.data)
}

func (b *Message) isEnd() bool {
	return b.offset >= uint32(len(b.data))
}

func (b *Message) Reader() io.Reader {
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(b.Length()+4))
	return io.MultiReader(
		bytes.NewReader(b.t),
		bytes.NewReader(length),
		bytes.NewReader(b.data),
	)
}

func (b *Message) SetType(t byte) *Message {
	b.t = []byte{t}
	return b
}

func NewBaseFromBytes(b []byte) *Message {
	return &Message{data: b}
}

func NewBase(len int) *Message {
	return &Message{data: make([]byte, len)}
}

type Value struct {
	data []byte
	null bool
}

func (v *Value) Length() int {
	return 4 + len(v.data)
}

func (v *Value) DataLength() uint32 {
	if v.null {
		return uint32(4294967295)
	}
	return uint32(len(v.data))
}

func (v *Value) DataBytes() []byte {
	return v.data
}

func NewValue(data []byte) Value {
	return Value{data: data}
}

func NewNullValue() Value {
	return Value{null: true}
}

func ReadValue(b *Message) Value {
	length := b.ReadUint32()
	if length == 4294967295 {
		return NewNullValue()
	}
	return Value{data: b.ReadByteN(length)}
}

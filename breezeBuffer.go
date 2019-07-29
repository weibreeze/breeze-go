package breeze

import (
	"encoding/binary"
	"io"
)

// Buffer is A variable-sized buffer of bytes with Read and Write methods.
// Buffer is not thread safe for multi goroutine operation.
type Buffer struct {
	buf     []byte // contents are the bytes buf[0 : wpos] in write, are the bytes buf[rpos: len(buf)] in read
	rpos    int    // read position
	wpos    int    // write position
	order   binary.ByteOrder
	temp    []byte
	context *Context
}

// NewBuffer create A empty Buffer with initial size
func NewBuffer(initSize int) *Buffer {
	return NewBufferWithOrder(initSize, binary.BigEndian)
}

// NewBufferWithOrder create A empty Buffer with initial size and byte order
func NewBufferWithOrder(initSize int, order binary.ByteOrder) *Buffer {
	return &Buffer{buf: make([]byte, initSize),
		order: order,
		temp:  make([]byte, 8),
	}
}

// CreateBuffer create A Buffer from data bytes
func CreateBuffer(data []byte) *Buffer {
	return CreateBufferWithOrder(data, binary.BigEndian)
}

// CreateBufferWithOrder create A Buffer from data bytes with bytes order
func CreateBufferWithOrder(data []byte, order binary.ByteOrder) *Buffer {
	return &Buffer{buf: data,
		order: order,
		temp:  make([]byte, 8),
		wpos:  len(data),
	}
}

// SetWPos set the write position of Buffer
func (b *Buffer) SetWPos(pos int) {
	if len(b.buf) < pos {
		b.grow(pos - len(b.buf))
	}
	b.wpos = pos
}

// GetWPos get the write position of Buffer
func (b *Buffer) GetWPos() int {
	return b.wpos
}

// SetRPos get the read position of Buffer
func (b *Buffer) SetRPos(pos int) {
	b.rpos = pos
}

// GetRPos get the read position of Buffer
func (b *Buffer) GetRPos() int {
	return b.rpos
}

// WriteByte write A byte append the Buffer, the wpos will increase one
func (b *Buffer) WriteByte(c byte) {
	if len(b.buf) < b.wpos+1 {
		b.grow(1)
	}
	b.buf[b.wpos] = c
	b.wpos++
}

// Write write A byte array append the Buffer, and the wpos will increase len(bytes)
func (b *Buffer) Write(bytes []byte) {
	l := len(bytes)
	if l > 0 {
		if len(b.buf) < b.wpos+l {
			b.grow(l)
		}
		copy(b.buf[b.wpos:], bytes)
		b.wpos += l
	}
}

// WriteUint16 write A uint16 append the Buffer according to buffer's order
func (b *Buffer) WriteUint16(u uint16) {
	if len(b.buf) < b.wpos+2 {
		b.grow(2)
	}
	b.order.PutUint16(b.temp, u)
	copy(b.buf[b.wpos:], b.temp[:2])
	b.wpos += 2
}

// WriteUint32 write a uint32 append to the Buffer
func (b *Buffer) WriteUint32(u uint32) {
	if len(b.buf) < b.wpos+4 {
		b.grow(4)
	}
	b.order.PutUint32(b.temp, u)
	copy(b.buf[b.wpos:], b.temp[:4])
	b.wpos += 4
}

// WriteUint64 write a uint64 append to the Buffer
func (b *Buffer) WriteUint64(u uint64) {
	if len(b.buf) < b.wpos+8 {
		b.grow(8)
	}
	b.order.PutUint64(b.temp, u)
	copy(b.buf[b.wpos:], b.temp[:8])
	b.wpos += 8
}

// WriteZigzag32 write a uint32 append to the Buffer with zigzag algorithm
func (b *Buffer) WriteZigzag32(u uint32) int {
	return b.WriteVarInt(uint64((u << 1) ^ uint32(int32(u)>>31)))
}

// WriteZigzag64 write a uint64 append to the Buffer with zigzag algorithm
func (b *Buffer) WriteZigzag64(u uint64) int {
	return b.WriteVarInt(uint64((u << 1) ^ uint64(int64(u)>>63)))
}

// WriteVarInt write a uint64 into buffer with variable length
func (b *Buffer) WriteVarInt(u uint64) int {
	l := 0
	for u >= 1<<7 {
		b.WriteByte(uint8(u&0x7f | 0x80))
		u >>= 7
		l++
	}
	b.WriteByte(uint8(u))
	l++
	return l
}

func (b *Buffer) grow(n int) {
	buf := make([]byte, 2*len(b.buf)+n)
	copy(buf, b.buf[:b.wpos])
	b.buf = buf
}

// Bytes return a bytes slice of under byte buffer.
func (b *Buffer) Bytes() []byte { return b.buf[:b.wpos] }

// Read read buffer's byte to byte array. return value n is read size.
func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.rpos >= len(b.buf) {
		return 0, io.EOF
	}

	n = copy(p, b.buf[b.rpos:])
	b.rpos += n
	return n, nil
}

// ReadFull read buffer's byte to byte array. if read size not equals len(p), will return error ErrNotEnough
func (b *Buffer) ReadFull(p []byte) error {
	if b.Remain() < len(p) {
		return ErrNotEnough
	}
	n := copy(p, b.buf[b.rpos:])
	if n < len(p) {
		return ErrNotEnough
	}
	b.rpos += n
	return nil
}

// ReadUint16 read a uint16 from buffer.
func (b *Buffer) ReadUint16() (n uint16, err error) {
	if b.Remain() < 2 {
		return 0, ErrNotEnough
	}
	n = b.order.Uint16(b.buf[b.rpos : b.rpos+2])
	b.rpos += 2
	return n, nil
}

// ReadInt read next int32
func (b *Buffer) ReadInt() (int, error) {
	n, err := b.ReadUint32()
	return int(n), err
}

// ReadUint32 read a uint32 from buffer.
func (b *Buffer) ReadUint32() (n uint32, err error) {
	if b.Remain() < 4 {
		return 0, ErrNotEnough
	}
	n = b.order.Uint32(b.buf[b.rpos : b.rpos+4])
	b.rpos += 4
	return n, nil
}

// ReadUint64 read a uint64 from buffer.
func (b *Buffer) ReadUint64() (n uint64, err error) {
	if b.Remain() < 8 {
		return 0, ErrNotEnough
	}
	n = b.order.Uint64(b.buf[b.rpos : b.rpos+8])
	b.rpos += 8
	return n, nil
}

// ReadZigzag64 read a zigzag uint64 from buffer.
func (b *Buffer) ReadZigzag64() (x uint64, err error) {
	x, err = b.ReadVarInt()
	if err != nil {
		return
	}
	x = (x >> 1) ^ uint64(-int64(x&1))
	return
}

// ReadZigzag32 read a zigzag uint32 from buffer.
func (b *Buffer) ReadZigzag32() (x uint64, err error) {
	x, err = b.ReadVarInt()
	if err != nil {
		return
	}
	x = uint64((uint32(x) >> 1) ^ uint32(-int32(x&1)))
	return
}

// ReadVarInt read a variable length uint64 form buffer
func (b *Buffer) ReadVarInt() (x uint64, err error) {
	var temp byte
	for offset := uint(0); offset < 64; offset += 7 {
		temp, err = b.ReadByte()
		if err != nil {
			return 0, err
		}
		if (temp & 0x80) != 0x80 {
			x |= uint64(temp) << offset
			return x, nil
		}
		x |= uint64(temp&0x7f) << offset
	}
	return 0, ErrOverflow
}

/*
Next get next n bytes from the buffer.
notice that return bytes is A slice of under byte array, hold the return value means hold all under byte array.
so , this method only for short-lived use
*/
func (b *Buffer) Next(n int) ([]byte, error) {
	m := b.Remain()
	if n > m {
		return nil, ErrNotEnough
	}
	data := b.buf[b.rpos : b.rpos+n]
	b.rpos += n
	return data, nil
}

// ReadByte read a byte form buffer
func (b *Buffer) ReadByte() (byte, error) {
	if b.rpos >= len(b.buf) {
		return 0, io.EOF
	}
	c := b.buf[b.rpos]
	b.rpos++
	return c, nil
}

// Reset reset the read position and write position to zero
func (b *Buffer) Reset() {
	b.rpos = 0
	b.wpos = 0
}

// Remain is used in buffer read, it return a size of bytes the buffer remained
func (b *Buffer) Remain() int { return b.wpos - b.rpos }

// Len return the len of buffer' bytes,
func (b *Buffer) Len() int { return b.wpos - 0 }

// Cap return the capacity of the under byte buffer
func (b *Buffer) Cap() int { return cap(b.buf) }

// GetContext get breeze context
func (b *Buffer) GetContext() *Context {
	if b.context == nil {
		b.context = &Context{}
	}
	return b.context
}

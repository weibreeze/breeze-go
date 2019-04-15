package breeze

import (
	"encoding/binary"
	"io"
)

// Buffer is A variable-sized buffer of bytes with Read and Write methods.
// Buffer is not thread safe for multi goroutine operation.
type Buffer struct {
	buf   []byte // contents are the bytes buf[0 : woff] in write, are the bytes buf[roff: len(buf)] in read
	rpos  int    // read position
	wpos  int    // write position
	order binary.ByteOrder
	temp  []byte
}

// NewBuffer create A empty Buffer with initial size
func NewBuffer(initsize int) *Buffer {
	return NewBufferWithOrder(initsize, binary.BigEndian)
}

// NewBufferWithOrder create A empty Buffer with initial size and byte order
func NewBufferWithOrder(initsize int, order binary.ByteOrder) *Buffer {
	return &Buffer{buf: make([]byte, initsize),
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
	if len(b.buf) < b.wpos+l {
		b.grow(l)
	}
	copy(b.buf[b.wpos:], bytes)
	b.wpos += l
}

// WriteUint16 write A uint16 append the Buffer acording to buffer's order
func (b *Buffer) WriteUint16(u uint16) {
	if len(b.buf) < b.wpos+2 {
		b.grow(2)
	}
	b.order.PutUint16(b.temp, u)
	copy(b.buf[b.wpos:], b.temp[:2])
	b.wpos += 2
}

func (b *Buffer) WriteUint32(u uint32) {
	if len(b.buf) < b.wpos+4 {
		b.grow(4)
	}
	b.order.PutUint32(b.temp, u)
	copy(b.buf[b.wpos:], b.temp[:4])
	b.wpos += 4
}

func (b *Buffer) WriteUint64(u uint64) {
	if len(b.buf) < b.wpos+8 {
		b.grow(8)
	}
	b.order.PutUint64(b.temp, u)
	copy(b.buf[b.wpos:], b.temp[:8])
	b.wpos += 8
}

func (b *Buffer) WriteZigzag32(u uint32) int {
	return b.WriteVarint(uint64((u << 1) ^ uint32(int32(u)>>31)))
}

func (b *Buffer) WriteZigzag64(u uint64) int {
	return b.WriteVarint(uint64((u << 1) ^ uint64(int64(u)>>63)))
}

func (b *Buffer) WriteVarint(u uint64) int {
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

func (b *Buffer) Bytes() []byte { return b.buf[:b.wpos] }

func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.rpos >= len(b.buf) {
		return 0, io.EOF
	}

	n = copy(p, b.buf[b.rpos:])
	b.rpos += n
	return n, nil
}

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

func (b *Buffer) ReadUint32() (n uint32, err error) {
	if b.Remain() < 4 {
		return 0, ErrNotEnough
	}
	n = b.order.Uint32(b.buf[b.rpos : b.rpos+4])
	b.rpos += 4
	return n, nil
}

func (b *Buffer) ReadUint64() (n uint64, err error) {
	if b.Remain() < 8 {
		return 0, ErrNotEnough
	}
	n = b.order.Uint64(b.buf[b.rpos : b.rpos+8])
	b.rpos += 8
	return n, nil
}

func (b *Buffer) ReadZigzag64() (x uint64, err error) {
	x, err = b.ReadVarint()
	if err != nil {
		return
	}
	x = (x >> 1) ^ uint64(-int64(x&1))
	return
}

func (b *Buffer) ReadZigzag32() (x uint64, err error) {
	x, err = b.ReadVarint()
	if err != nil {
		return
	}
	x = uint64((uint32(x) >> 1) ^ uint32(-int32(x&1)))
	return
}

func (b *Buffer) ReadVarint() (x uint64, err error) {
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

/**
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

func (b *Buffer) ReadByte() (byte, error) {
	if b.rpos >= len(b.buf) {
		return 0, io.EOF
	}
	c := b.buf[b.rpos]
	b.rpos++
	return c, nil
}

func (b *Buffer) Reset() {
	b.rpos = 0
	b.wpos = 0
}

func (b *Buffer) Remain() int { return b.wpos - b.rpos }

func (b *Buffer) Len() int { return b.wpos - 0 }

func (b *Buffer) Cap() int { return cap(b.buf) }
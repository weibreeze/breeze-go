package breeze

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	initsize := 1
	buf := NewBufferWithOrder(initsize, binary.LittleEndian)
	if buf.order != binary.LittleEndian {
		t.Errorf("order not correct. real:%s, expect:%s\n", buf.order, binary.LittleEndian)
	}

	buf = NewBuffer(initsize)
	if len(buf.buf) != initsize {
		t.Errorf("wrong initsize. real size:%d, expect size:%d", len(buf.buf), initsize)
	}
	if buf.order != binary.BigEndian {
		t.Errorf("default order not bigendian.")
	}
	if buf.Len() != 0 {
		t.Errorf("new buf length not zero.")
	}
	if buf.Cap() != initsize {
		t.Errorf("buf cap not correct.real:%d, expect:%d\n", buf.Cap(), initsize)
	}
	if buf.wpos != 0 || buf.rpos != 0 {
		t.Errorf("buf wpos or rpos init value not correct.wpos:%d, rpos:%d\n", buf.wpos, buf.rpos)
	}

	buf.SetWPos(3)
	if buf.Cap() < 3 || buf.wpos != 3 {
		t.Errorf("buf SetWPos expand buffer failed: %d", buf.Cap())
	}
}

func TestWrite(t *testing.T) {
	buf := NewBuffer(32)
	// write byte
	buf.WriteByte('A')
	if buf.wpos != 1 || buf.Len() != 1 {
		t.Errorf("buf wpos not correct.buf:%+v\n", buf)
	}

	// write []byte
	oldpos := buf.GetWPos()
	size := 20
	bytes := make([]byte, 0, size)
	for i := 0; i < size; i++ {
		bytes = append(bytes, 'b')
	}
	buf.Write(bytes)
	if buf.wpos != size+oldpos || buf.Len() != size+oldpos {
		t.Errorf("buf wpos not correct.buf:%+v\n", buf)
	}

	// write uint
	buf.Reset()
	buf.WriteUint32(uint32(123))
	buf.WriteUint64(uint64(789))
	tempbytes := buf.Bytes()
	if 123 != int(binary.BigEndian.Uint32(tempbytes[:4])) {
		t.Errorf("write uint32 not correct.buf:%+v\n", buf)
	}
	if 789 != int(binary.BigEndian.Uint64(tempbytes[4:12])) {
		t.Errorf("write uint32 not correct.buf:%+v\n", buf)
	}
}

func TestWriteWithGrow(t *testing.T) {
	// write with grow
	buf := NewBuffer(32)
	size := 107
	bytes := make([]byte, 0, size)
	for i := 0; i < size; i++ {
		bytes = append(bytes, 'c')
	}
	buf.Write(bytes)
	if buf.wpos != size || buf.Len() != size {
		t.Errorf("buf wpos not correct.buf:%+v\n", buf)
	}
	// set wpos
	buf.Reset()
	if buf.GetWPos() != 0 {
		t.Errorf("buf reset wpos not correct.buf:%+v\n", buf)
	}
	buf.SetWPos(4)
	buf.Write(bytes)
	buf.SetWPos(0)
	buf.WriteUint32(uint32(len(bytes)))
	buf.SetWPos(buf.GetWPos() + len(bytes))
	tempbytes := buf.Bytes()
	if int(binary.BigEndian.Uint32(tempbytes[:4])) != len(bytes) {
		t.Errorf("write uint32 not correct.buf:%+v\n", buf)
	}
	if len(tempbytes) != 4+len(bytes) {
		t.Errorf("set wpos test not correct.buf:%+v\n", buf)
	}
}

func TestRead(t *testing.T) {
	buf := NewBuffer(128)
	b := byte(45)
	s := "jlk>E&(*L#?>"
	i16 := uint16(34)
	i32 := uint32(56)
	i64 := uint64(56578)

	buf.WriteByte(b)
	buf.Write([]byte(s))
	buf.WriteUint16(i16)
	buf.WriteUint32(i32)
	buf.WriteUint64(i64)

	bytes := buf.Bytes()
	buf2 := CreateBuffer(bytes)
	// read full
	rbytes := make([]byte, len(bytes))
	buf2.ReadFull(rbytes)
	if !reflect.DeepEqual(rbytes, bytes) {
		t.Errorf("wrong buf bytes. expect:%v, real:%v\n", bytes, rbytes)
	}

	//read next
	buf2.rpos = 0
	rbytes, _ = buf2.Next(len(bytes))
	if !reflect.DeepEqual(rbytes, bytes) {
		t.Errorf("wrong buf bytes. expect:%v, real:%v\n", bytes, rbytes)
	}

	//read
	buf2.rpos = 0
	rb, _ := buf2.ReadByte()
	if rb != b {
		t.Errorf("wrong buf value. expect:%v, real:%v\n", b, rb)
	}
	rbytes = make([]byte, len(s))
	buf2.Read(rbytes)
	if string(rbytes) != s {
		t.Errorf("wrong buf value. expect:%v, real:%v\n", s, string(rbytes))
	}
	ri16, _ := buf2.ReadUint16()
	if ri16 != i16 {
		t.Errorf("wrong buf value. expect:%v, real:%v\n", i16, ri16)
	}
	ri32, _ := buf2.ReadUint32()
	if ri32 != i32 {
		t.Errorf("wrong buf value. expect:%v, real:%v\n", i32, ri32)
	}
	ri64, _ := buf2.ReadUint64()
	if ri64 != i64 {
		t.Errorf("wrong buf value. expect:%v, real:%v\n", i64, ri64)
	}
}

func TestZigzag(t *testing.T) {
	times := 128
	f1 := 1678
	buf := NewBuffer(times * 8)
	// zigzag32
	for i := 0; i < times; i++ {
		buf.WriteZigzag32(uint32(i * f1))
	}
	bytes := buf.Bytes()
	//fmt.Printf("bytes:%v\n", bytes)
	buf2 := CreateBuffer(bytes)
	for i := 0; i < times; i++ {
		ni, err := buf2.ReadZigzag32()
		if err != nil || int(ni) != i*f1 {
			t.Errorf("zigzag32 not correct. ni: %d, i:%d, err :%v, buf:%v\n", ni, i, err, buf2)
		}
	}

	//zigzag64
	buf.Reset()
	f2 := 7289374928
	for i := 0; i < times; i++ {
		buf.WriteZigzag64(uint64(i * f2))
	}
	bytes = buf.Bytes()
	//fmt.Printf("bytes:%v\n", bytes)
	buf2 = CreateBuffer(bytes)
	for i := 0; i < times; i++ {
		ni, err := buf2.ReadZigzag64()
		if err != nil || int(ni) != i*f2 {
			t.Errorf("zigzag64 not correct. ni: %d, i:%d, err :%v, buf:%v\n", ni, i, err, buf2)
		}
	}
}

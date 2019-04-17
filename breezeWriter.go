package breeze

import (
	"github.com/pkg/errors"
	"math"
	"reflect"
)

// WriteFieldsFunc is a func interface of how to write all fields of a breeze message to the buffer.
type WriteFieldsFunc func(buf *Buffer)

// WriteBool write a bool value into the buffer
func WriteBool(buf *Buffer, b bool) {
	if b {
		buf.WriteByte(TRUE)
	} else {
		buf.WriteByte(FALSE)
	}
}

// WriteString write a string value into the buffer
func WriteString(buf *Buffer, s string) {
	buf.WriteByte(STRING)
	buf.WriteZigzag32(uint32(len(s)))
	buf.Write([]byte(s))
}

// WriteByte write a byte value into the buffer
func WriteByte(buf *Buffer, b byte) {
	buf.WriteByte(BYTE)
	buf.WriteByte(b)
}

// WriteBytes write a byte slice into the buffer
func WriteBytes(buf *Buffer, bytes []byte) {
	buf.WriteByte(BYTES)
	buf.WriteUint32(uint32(len(bytes)))
	buf.Write(bytes)
}

// WriteInt16 write a uint16 value into the buffer
func WriteInt16(buf *Buffer, ui uint16) {
	buf.WriteByte(INT16)
	buf.WriteUint16(ui)
}

// WriteInt32 write a uint32 value into the buffer
func WriteInt32(buf *Buffer, ui uint32) {
	buf.WriteByte(INT32)
	buf.WriteZigzag32(ui)
}

// WriteInt64 write a uint64 value into the buffer
func WriteInt64(buf *Buffer, ui uint64) {
	buf.WriteByte(INT64)
	buf.WriteZigzag64(ui)
}

// WriteFloat32 write a float32 value into the buffer
func WriteFloat32(buf *Buffer, f float32) {
	buf.WriteByte(FLOAT32)
	buf.WriteUint32(math.Float32bits(float32(f)))
}

// WriteFloat64 write a float64 value into the buffer
func WriteFloat64(buf *Buffer, f float64) {
	buf.WriteByte(FLOAT64)
	buf.WriteUint64(math.Float64bits(f))
}

// WriteValue can write primtive type and ptr of primtive type, and breeze message.
func WriteValue(buf *Buffer, v interface{}) error {
	if v == nil {
		buf.WriteByte(NULL)
		return nil
	}
	if msg, ok := v.(Message); ok {
		return msg.WriteTo(buf)
	}

	var rv reflect.Value
	if nrv, ok := v.(reflect.Value); ok {
		rv = nrv
	} else {
		rv = reflect.ValueOf(v)
	}
	k := rv.Kind()
	if k == reflect.Ptr {
		if rv.CanInterface() { //message
			realV := rv.Interface()
			if msg, ok := realV.(Message); ok {
				return msg.WriteTo(buf)
			}
		}
		//TODO extension for custom process
		rv = rv.Elem()
		k = rv.Kind()
	}
	return writeByKind(buf, k, rv)
}

func writeByKind(buf *Buffer, k reflect.Kind, rv reflect.Value) error {
	if k == reflect.Interface {
		rv = reflect.ValueOf(rv.Interface())
		k = rv.Kind()
	}
	switch k {
	case reflect.String:
		WriteString(buf, rv.String())
	case reflect.Bool:
		WriteBool(buf, rv.Bool())
	case reflect.Uint8:
		WriteByte(buf, byte(rv.Uint()))
	case reflect.Int16, reflect.Uint16:
		WriteInt16(buf, uint16(rv.Int()))
	case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
		WriteInt32(buf, uint32(rv.Int()))
	case reflect.Int64, reflect.Uint64:
		WriteInt64(buf, uint64(rv.Int()))
	case reflect.Float32:
		WriteFloat32(buf, float32(rv.Float()))
	case reflect.Float64:
		WriteFloat64(buf, rv.Float())
	case reflect.Slice:
		t := rv.Type().Elem().Kind()
		if t == reflect.Uint8 {
			WriteBytes(buf, rv.Bytes())
		} else {
			return writeArray(buf, rv)
		}
	case reflect.Map:
		return writeMap(buf, rv)
	default:
		return errors.New("breeze: unsupport type " + k.String())
	}
	return nil
}

// WriteMessage write a breeze message according to WriteFieldsFunc
func WriteMessage(buf *Buffer, name string, fieldsFunc WriteFieldsFunc) (err error) {
	defer func() {
		if inner := recover(); inner != nil {
			err = inner.(error)
		}
	}()
	buf.WriteByte(MESSAGE)
	WriteString(buf, name)
	pos := skipLength(buf)
	fieldsFunc(buf)
	writeLength(buf, pos)
	return err
}

// WriteMessageField write a message field into buffer
func WriteMessageField(buf *Buffer, index int, v interface{}) {
	if v != nil {
		buf.WriteZigzag32(uint32(index))
		err := WriteValue(buf, v)
		if err != nil {
			panic(err)
		}
	}
}

func writeArray(buf *Buffer, v reflect.Value) (err error) {
	buf.WriteByte(ARRAY)
	pos := skipLength(buf)
	for i := 0; i < v.Len(); i++ {
		err = WriteValue(buf, v.Index(i))
		if err != nil {
			return err
		}
	}
	writeLength(buf, pos)
	return nil
}

func writeMap(buf *Buffer, v reflect.Value) (err error) {
	buf.WriteByte(MAP)
	pos := skipLength(buf)
	err = rangeMap(buf, v)
	writeLength(buf, pos)
	return err
}

// keep 4 bytes for write length later
func skipLength(buf *Buffer) int {
	pos := buf.GetWPos()
	buf.SetWPos(pos + 4)
	return pos
}

// write length into keep position
func writeLength(buf *Buffer, keepPos int) {
	curPos := buf.GetWPos()
	buf.SetWPos(keepPos)
	buf.WriteUint32(uint32(curPos - keepPos - 4))
	buf.SetWPos(curPos)
}

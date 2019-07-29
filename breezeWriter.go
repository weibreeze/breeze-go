package breeze

import (
	"github.com/pkg/errors"
	"math"
	"reflect"
)

// WriteFieldsFunc is a func interface of how to write all fields of a breeze message to the buffer.
type WriteFieldsFunc func(buf *Buffer)

// WriteElemFunc write map or array elements
type WriteElemFunc func(buf *Buffer)

// WriteBool write a bool value into the buffer
func WriteBool(buf *Buffer, b bool, withType bool) {
	if b {
		buf.WriteByte(TrueType)
	} else {
		buf.WriteByte(FalseType)
	}
}

// WriteString write a string value into the buffer
func WriteString(buf *Buffer, s string, withType bool) {
	if withType {
		l := len(s)
		if l <= DirectStringMaxLength { // direct string
			buf.WriteByte(byte(l))
			buf.Write([]byte(s))
			return
		}
		WriteStringType(buf)
	}
	buf.WriteVarInt(uint64(len(s)))
	buf.Write([]byte(s))
}

// WriteByte write a byte value into the buffer
func WriteByte(buf *Buffer, b byte, withType bool) {
	if withType {
		WriteByteType(buf)
	}
	buf.WriteByte(b)
}

// WriteBytes write a byte slice into the buffer
func WriteBytes(buf *Buffer, bytes []byte, withType bool) {
	if withType {
		WriteBytesType(buf)
	}
	buf.WriteUint32(uint32(len(bytes)))
	buf.Write(bytes)
}

// WriteInt16 write a int16 value into the buffer
func WriteInt16(buf *Buffer, i int16, withType bool) {
	if withType {
		WriteInt16Type(buf)
	}
	buf.WriteUint16(uint16(i))
}

// WriteInt32 write a int32 value into the buffer
func WriteInt32(buf *Buffer, i int32, withType bool) {
	if withType {
		if i >= DirectInt32MinValue && i <= DirectInt32MaxValue {
			buf.WriteByte(byte(i + Int32Zero))
			return
		}
		WriteInt32Type(buf)
	}
	buf.WriteZigzag32(uint32(i))
}

// WriteInt64 write a int64 value into the buffer
func WriteInt64(buf *Buffer, i int64, withType bool) {
	if withType {
		if i >= DirectInt64MinValue && i <= DirectInt64MaxValue {
			buf.WriteByte(byte(i + Int64Zero))
			return
		}
		WriteInt64Type(buf)
	}
	buf.WriteZigzag64(uint64(i))
}

// WriteFloat32 write a float32 value into the buffer
func WriteFloat32(buf *Buffer, f float32, withType bool) {
	if withType {
		WriteFloat32Type(buf)
	}
	buf.WriteUint32(math.Float32bits(float32(f)))
}

// WriteFloat64 write a float64 value into the buffer
func WriteFloat64(buf *Buffer, f float64, withType bool) {
	if withType {
		WriteFloat64Type(buf)
	}
	buf.WriteUint64(math.Float64bits(f))
}

// WritePackedMap write packed map by WriteElemFunc
func WritePackedMap(buf *Buffer, withType bool, size int, f WriteElemFunc) {
	if withType {
		WritePackedMapType(buf)
	}
	buf.WriteVarInt(uint64(size))
	f(buf)
}

// WritePackedArray write packed array by WriteElemFunc
func WritePackedArray(buf *Buffer, withType bool, size int, f WriteElemFunc) {
	if withType {
		WritePackedArrayType(buf)
	}
	buf.WriteVarInt(uint64(size))
	f(buf)
}

// WriteStringStringMapEntries write map[string]string directly
func WriteStringStringMapEntries(buf *Buffer, m map[string]string) {
	WriteStringType(buf)
	WriteStringType(buf)
	for k, v := range m {
		WriteString(buf, k, false)
		WriteString(buf, v, false)
	}
}

// WriteStringInt32MapEntries write map[string]int32 directly
func WriteStringInt32MapEntries(buf *Buffer, m map[string]int32) {
	WriteStringType(buf)
	WriteInt32Type(buf)
	for k, v := range m {
		WriteString(buf, k, false)
		WriteInt32(buf, v, false)
	}
}

// WriteStringInt64MapEntries write map[string]int64 directly
func WriteStringInt64MapEntries(buf *Buffer, m map[string]int64) {
	WriteStringType(buf)
	WriteInt64Type(buf)
	for k, v := range m {
		WriteString(buf, k, false)
		WriteInt64(buf, v, false)
	}
}

// WriteStringArrayElems write []string directly
func WriteStringArrayElems(buf *Buffer, a []string) {
	WriteStringType(buf)
	for _, v := range a {
		WriteString(buf, v, false)
	}
}

// WriteInt32ArrayElems write []int32 directly
func WriteInt32ArrayElems(buf *Buffer, a []int32) {
	WriteInt32Type(buf)
	for _, v := range a {
		WriteInt32(buf, v, false)
	}
}

// WriteInt64ArrayElems write []int64 directly
func WriteInt64ArrayElems(buf *Buffer, a []int64) {
	WriteInt64Type(buf)
	for _, v := range a {
		WriteInt64(buf, v, false)
	}
}

// WriteMessageWithoutType write a breeze message according to WriteFieldsFunc. without message type
func WriteMessageWithoutType(buf *Buffer, fieldsFunc WriteFieldsFunc) (err error) {
	defer func() {
		if inner := recover(); inner != nil {
			err = inner.(error)
		}
	}()
	pos := skipLength(buf)
	fieldsFunc(buf)
	writeLength(buf, pos)
	return err
}

//========== write BreezeType to buffer, only for packed model(packed map and packed array) =====================

// WriteBoolType write bool type
func WriteBoolType(buf *Buffer) {
	buf.WriteByte(TrueType)
}

// WriteStringType write string type
func WriteStringType(buf *Buffer) {
	buf.WriteByte(StringType)
}

// WriteByteType write byte type
func WriteByteType(buf *Buffer) {
	buf.WriteByte(ByteType)
}

// WriteBytesType write byte array type
func WriteBytesType(buf *Buffer) {
	buf.WriteByte(BytesType)
}

// WriteInt16Type write int16 type
func WriteInt16Type(buf *Buffer) {
	buf.WriteByte(Int16Type)
}

// WriteInt32Type write int32 type
func WriteInt32Type(buf *Buffer) {
	buf.WriteByte(Int32Type)
}

// WriteInt64Type write int64 type
func WriteInt64Type(buf *Buffer) {
	buf.WriteByte(Int64Type)
}

// WriteFloat32Type write float32 type
func WriteFloat32Type(buf *Buffer) {
	buf.WriteByte(Float32Type)
}

// WriteFloat64Type write float64 type
func WriteFloat64Type(buf *Buffer) {
	buf.WriteByte(Float64Type)
}

// WritePackedMapType write packed map type
func WritePackedMapType(buf *Buffer) {
	buf.WriteByte(PackedMapType)
}

// WritePackedArrayType write packed array type
func WritePackedArrayType(buf *Buffer) {
	buf.WriteByte(PackedArrayType)
}

// WriteMessageType write message type. it can be a ref index or message with name
func WriteMessageType(buf *Buffer, name string) {
	index := buf.GetContext().getMessageTypeIndex(name)
	if index < 0 { // first write
		buf.WriteByte(MessageType)
		WriteString(buf, name, false)
		buf.GetContext().putMessageType(name)
	} else {
		if index > DirectRefMessageMaxValue {
			buf.WriteByte(RefMessageType)
			buf.WriteVarInt(uint64(index))
		} else {
			buf.WriteByte(byte(RefMessageType + index))
		}
	}
}

//========== write message field by type. it will not write if the value is default =====================

// WriteBoolField write field with index
func WriteBoolField(buf *Buffer, index int, b bool) {
	if b {
		buf.WriteVarInt(uint64(index))
		WriteBool(buf, b, true)
	}
}

// WriteStringField write field with index
func WriteStringField(buf *Buffer, index int, s string) {
	if s != "" {
		buf.WriteVarInt(uint64(index))
		WriteString(buf, s, true)
	}
}

// WriteByteField write field with index
func WriteByteField(buf *Buffer, index int, b byte) {
	buf.WriteVarInt(uint64(index))
	WriteByte(buf, b, true)
}

// WriteBytesField write field with index
func WriteBytesField(buf *Buffer, index int, b []byte) {
	if len(b) > 0 {
		buf.WriteVarInt(uint64(index))
		WriteBytes(buf, b, true)
	}
}

// WriteInt16Field write field with index
func WriteInt16Field(buf *Buffer, index int, i int16) {
	if i != 0 {
		buf.WriteVarInt(uint64(index))
		WriteInt16(buf, i, true)
	}
}

// WriteInt32Field write field with index
func WriteInt32Field(buf *Buffer, index int, i int32) {
	if i != 0 {
		buf.WriteVarInt(uint64(index))
		WriteInt32(buf, i, true)
	}
}

// WriteInt64Field write field with index
func WriteInt64Field(buf *Buffer, index int, i int64) {
	if i != 0 {
		buf.WriteVarInt(uint64(index))
		WriteInt64(buf, i, true)
	}
}

// WriteFloat32Field write field with index
func WriteFloat32Field(buf *Buffer, index int, f float32) {
	if f != 0 {
		buf.WriteVarInt(uint64(index))
		WriteFloat32(buf, f, true)
	}
}

// WriteFloat64Field write field with index
func WriteFloat64Field(buf *Buffer, index int, f float64) {
	if f != 0 {
		buf.WriteVarInt(uint64(index))
		WriteFloat64(buf, f, true)
	}
}

// WriteMapField write field with index
func WriteMapField(buf *Buffer, index int, size int, f WriteElemFunc) {
	buf.WriteVarInt(uint64(index))
	WritePackedMap(buf, true, size, f)
}

// WriteArrayField write field with index
func WriteArrayField(buf *Buffer, index int, size int, f WriteElemFunc) {
	buf.WriteVarInt(uint64(index))
	WritePackedArray(buf, true, size, f)
}

// WriteMessageField write field with index
func WriteMessageField(buf *Buffer, index int, m Message) {
	buf.WriteVarInt(uint64(index))
	WriteMessageType(buf, m.GetName())
	m.WriteTo(buf)
}

// WriteField write an any type field into buffer.
func WriteField(buf *Buffer, index int, v interface{}) {
	if v != nil {
		buf.WriteVarInt(uint64(index))
		err := WriteValue(buf, v)
		if err != nil {
			panic(err)
		}
	}
}

// WriteValue can write primitive type and ptr of primitive type, and breeze message.
func WriteValue(buf *Buffer, v interface{}) error {
	if v == nil {
		buf.WriteByte(NullType)
		return nil
	}
	if msg, ok := v.(Message); ok {
		return writeMessage(buf, msg)
	}

	var rv reflect.Value
	if nrv, ok := v.(reflect.Value); ok {
		rv = nrv
	} else {
		rv = reflect.ValueOf(v)
	}
	return writeReflectValue(buf, rv, true)

}

func writeReflectValue(buf *Buffer, rv reflect.Value, withType bool) error {
	k := rv.Kind()
	if k == reflect.Ptr {
		if rv.CanInterface() { //message
			realV := rv.Interface()
			if msg, ok := realV.(Message); ok {
				return writeMessage(buf, msg)
			}
		}
		//TODO extension for custom process
		rv = rv.Elem()
		k = rv.Kind()
	}
	if k == reflect.Interface {
		rv = reflect.ValueOf(rv.Interface())
		k = rv.Kind()
	}
	switch k {
	case reflect.String:
		WriteString(buf, rv.String(), withType)
	case reflect.Bool:
		WriteBool(buf, rv.Bool(), withType)
	case reflect.Int, reflect.Int32:
		WriteInt32(buf, int32(rv.Int()), withType)
	case reflect.Int64:
		WriteInt64(buf, rv.Int(), withType)
	case reflect.Map:
		return writeMap(buf, rv, withType)
	case reflect.Slice:
		t := rv.Type().Elem().Kind()
		if t == reflect.Uint8 {
			WriteBytes(buf, rv.Bytes(), withType)
		} else {
			return writeArray(buf, rv, withType)
		}
	case reflect.Uint, reflect.Uint32:
		WriteInt32(buf, int32(rv.Uint()), withType)
	case reflect.Uint64:
		WriteInt64(buf, int64(rv.Uint()), withType)
	case reflect.Uint8:
		WriteByte(buf, byte(rv.Uint()), withType)
	case reflect.Int16:
		WriteInt16(buf, int16(rv.Int()), withType)
	case reflect.Uint16:
		WriteInt16(buf, int16(rv.Uint()), withType)
	case reflect.Float32:
		WriteFloat32(buf, float32(rv.Float()), withType)
	case reflect.Float64:
		WriteFloat64(buf, rv.Float(), withType)
	default:
		return errors.New("breeze: unsupported type " + k.String())
	}
	return nil
}

func writeType(buf *Buffer, rv reflect.Value) {
	k := rv.Kind()
	if k == reflect.Ptr {
		if rv.CanInterface() { //message
			realV := rv.Interface()
			if msg, ok := realV.(Message); ok {
				WriteMessageType(buf, msg.GetName())
				return
			}
		}
		rv = rv.Elem()
		k = rv.Kind()
	}
	switch k {
	case reflect.String:
		WriteStringType(buf)
	case reflect.Bool:
		WriteBoolType(buf)
	case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
		WriteInt32Type(buf)
	case reflect.Int64, reflect.Uint64:
		WriteInt64Type(buf)
	case reflect.Map:
		if canPackMap(rv.Type()) {
			WritePackedMapType(buf)
		} else {
			buf.WriteByte(MapType)
		}
	case reflect.Slice:
		tp := rv.Type()
		if tp.Elem().Kind() == reflect.Uint8 {
			WriteBytesType(buf)
		} else {
			if canPackArray(rv.Type()) {
				WritePackedArrayType(buf)
			} else {
				buf.WriteByte(ArrayType)
			}
		}
	case reflect.Uint8:
		WriteByteType(buf)
	case reflect.Int16, reflect.Uint16:
		WriteInt16Type(buf)
	case reflect.Float32:
		WriteFloat32Type(buf)
	case reflect.Float64:
		WriteFloat64Type(buf)
	default:
		panic(errors.New("breeze: unsupported type " + k.String()))
	}
}

func writeArray(buf *Buffer, v reflect.Value, withType bool) (err error) {
	if canPackArray(v.Type()) {
		defer func() {
			if inner := recover(); inner != nil {
				err = inner.(error)
			}
		}()
		WritePackedArray(buf, withType, v.Len(), func(buf *Buffer) {
			for i := 0; i < v.Len(); i++ {
				elem := v.Index(i)
				if i == 0 {
					writeType(buf, elem)
				}
				err = writeReflectValue(buf, elem, false)
				if err != nil {
					panic(err)
				}
			}
		})
	} else {
		if withType {
			buf.WriteByte(ArrayType)
		}
		buf.WriteVarInt(uint64(v.Len()))
		for i := 0; i < v.Len(); i++ {
			err = writeReflectValue(buf, v.Index(i), true)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func writeMap(buf *Buffer, v reflect.Value, withType bool) (err error) {
	if canPackMap(v.Type()) {
		defer func() {
			if inner := recover(); inner != nil {
				err = inner.(error)
			}
		}()
		WritePackedMap(buf, withType, v.Len(), func(buf *Buffer) {
			rangePackedMap(buf, v)
		})
	} else {
		if withType {
			buf.WriteByte(MapType)
		}
		buf.WriteVarInt(uint64(v.Len()))
		err = rangeMap(buf, v)
	}
	return err
}

func canPackArray(t reflect.Type) bool {
	return t.Elem().Kind() != reflect.Interface
}

func canPackMap(t reflect.Type) bool {
	return (t.Key().Kind() != reflect.Interface) && (t.Elem().Kind() != reflect.Interface)
}

func writeMessage(buf *Buffer, message Message) error {
	WriteMessageType(buf, message.GetName())
	return message.WriteTo(buf)
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

package breeze

import (
	"github.com/pkg/errors"
	"math"
	"reflect"
	"strconv"
)

// default value of breeze reader
const (
	DefaultSize = 16
)

// ReadFieldsFunc is a func interface of how to read a breeze message field by field index.
type ReadFieldsFunc func(buf *Buffer, index int) error

// ReadElemFunc read one element of map or array
type ReadElemFunc func(buf *Buffer) error

// ReadBool read a bool value into the bool pointer
func ReadBool(buf *Buffer, b *bool) (err error) {
	*b, err = ReadBoolWithoutType(buf)
	return err
}

// ReadBoolWithoutType read without type
func ReadBoolWithoutType(buf *Buffer) (bool, error) {
	tp, err := buf.ReadByte()
	if err != nil {
		return false, err
	}
	if tp == TrueType {
		return true, nil
	} else if tp == FalseType {
		return false, nil
	} else {
		return false, errors.New("BreezeRead: wrong type expect bool, real " + strconv.Itoa(int(tp)))
	}
}

// ReadString read a string value into a pointer
func ReadString(buf *Buffer, s *string) (err error) {
	tp, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if tp <= DirectStringMaxType {
		bytes, err := buf.Next(int(tp))
		if err != nil {
			return err
		}
		*s = string(bytes)
		return nil
	}
	switch int(tp) {
	case StringType:
		*s, err = ReadStringWithoutType(buf)
	case Int32Type:
		i32, err := ReadInt32WithoutType(buf)
		if err != nil {
			return err
		}
		*s = strconv.Itoa(int(i32))
	case Int64Type:
		i64, err := ReadInt64WithoutType(buf)
		if err != nil {
			return err
		}
		*s = strconv.FormatInt(i64, 10)
	case Int16Type:
		i16, err := ReadInt16WithoutType(buf)
		if err != nil {
			return err
		}
		*s = strconv.Itoa(int(i16))
	case Float32Type:
		f32, err := ReadFloat32WithoutType(buf)
		if err != nil {
			return err
		}
		*s = strconv.FormatFloat(float64(f32), 'f', -1, 32)
	case Float64Type:
		f64, err := ReadFloat64WithoutType(buf)
		if err != nil {
			return err
		}
		*s = strconv.FormatFloat(f64, 'f', -1, 64)
	default:
		err = errors.New("Breeze: not convert to string, type " + strconv.Itoa(int(tp)))
	}
	return err
}

// ReadByte read a byte value into a pointer
func ReadByte(buf *Buffer, b *byte) error {
	err := checkType(buf, ByteType)
	if err != nil {
		return err
	}
	*b, err = buf.ReadByte()
	return err
}

// ReadBytes read a byte slice value into a pointer
func ReadBytes(buf *Buffer, bytes *[]byte) error {
	err := checkType(buf, BytesType)
	if err != nil {
		return err
	}
	*bytes, err = ReadBytesWithoutType(buf)
	return err
}

// ReadInt16 read a int16 value into a pointer
func ReadInt16(buf *Buffer, i *int16) error {
	tp, err := buf.ReadByte()
	if err != nil {
		return err
	}
	switch int(tp) {
	case Int16Type:
		*i, err = ReadInt16WithoutType(buf)
	case StringType:
		s, err := ReadStringWithoutType(buf)
		if err != nil {
			return err
		}
		si, err := strconv.Atoi(s)
		*i = int16(si)
		return err
	case Int32Type:
		i32, err := ReadInt32WithoutType(buf)
		if err != nil {
			return err
		}
		*i = int16(i32)
	case Int64Type:
		i64, err := ReadInt64WithoutType(buf)
		if err != nil {
			return err
		}
		*i = int16(i64)
	default:
		err = errors.New("Breeze: not convert to int16, type " + strconv.Itoa(int(tp)))
	}
	return err
}

// ReadInt read int32 as int
func ReadInt(buf *Buffer, i *int) error {
	var i32 int32
	err := ReadInt32(buf, &i32)
	*i = int(i32)
	return err
}

// ReadInt32 read a int32 value into a pointer
func ReadInt32(buf *Buffer, i *int32) error {
	tp, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if tp >= DirectInt32MinType && tp <= DirectInt32MaxType {
		*i = int32(tp) - Int32Zero
		return nil
	}
	switch int(tp) {
	case Int32Type:
		*i, err = ReadInt32WithoutType(buf)
	case StringType:
		s, err := ReadStringWithoutType(buf)
		if err != nil {
			return err
		}
		si, err := strconv.Atoi(s)
		*i = int32(si)
		return err
	case Int64Type:
		i64, err := ReadInt64WithoutType(buf)
		if err != nil {
			return err
		}
		*i = int32(i64)
	case Int16Type:
		i16, err := ReadInt16WithoutType(buf)
		if err != nil {
			return err
		}
		*i = int32(i16)
	default:
		err = errors.New("Breeze: not convert to int32, type " + strconv.Itoa(int(tp)))
	}
	return err
}

// ReadInt64 read a int64 value into a pointer
func ReadInt64(buf *Buffer, i *int64) error {
	tp, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if tp >= DirectInt64MinType && tp <= DirectInt64MaxType {
		*i = int64(tp) - Int64Zero
		return nil
	}
	switch int(tp) {
	case Int64Type:
		*i, err = ReadInt64WithoutType(buf)
	case StringType:
		s, err := ReadStringWithoutType(buf)
		if err != nil {
			return err
		}
		*i, err = strconv.ParseInt(s, 10, 64)
		return err
	case Int32Type:
		i32, err := ReadInt32WithoutType(buf)
		if err != nil {
			return err
		}
		*i = int64(i32)
	case Int16Type:
		i16, err := ReadInt16WithoutType(buf)
		if err != nil {
			return err
		}
		*i = int64(i16)
	default:
		err = errors.New("Breeze: not convert to int64, type " + strconv.Itoa(int(tp)))
	}
	return err
}

// ReadFloat32 read a float32 value into a pointer
func ReadFloat32(buf *Buffer, f *float32) error {
	tp, err := buf.ReadByte()
	if err != nil {
		return err
	}
	switch int(tp) {
	case Float32Type:
		*f, err = ReadFloat32WithoutType(buf)
	case Float64Type:
		f64, err := ReadFloat64WithoutType(buf)
		if err != nil {
			return err
		}
		*f = float32(f64)
	case StringType:
		s, err := ReadStringWithoutType(buf)
		if err != nil {
			return err
		}
		f64, err := strconv.ParseFloat(s, 64)
		*f = float32(f64)
		return err
	default:
		err = errors.New("Breeze: not convert to float32, type " + strconv.Itoa(int(tp)))
	}
	return err
}

// ReadFloat64 read a float64 value into a pointer
func ReadFloat64(buf *Buffer, f *float64) error {
	tp, err := buf.ReadByte()
	if err != nil {
		return err
	}
	switch int(tp) {
	case Float64Type:
		*f, err = ReadFloat64WithoutType(buf)
	case Float32Type:
		f32, err := ReadFloat32WithoutType(buf)
		if err != nil {
			return err
		}
		*f = float64(f32)
	case StringType:
		s, err := ReadStringWithoutType(buf)
		if err != nil {
			return err
		}
		*f, err = strconv.ParseFloat(s, 64)
	default:
		err = errors.New("Breeze: not convert to float64 type " + strconv.Itoa(int(tp)))
	}
	return err
}

// ReadPackedSize read size of packed map or packed array
func ReadPackedSize(buf *Buffer, withType bool) (int, error) {
	if withType {
		_, err := buf.ReadByte() //ignored
		if err != nil {
			return 0, err
		}
	}
	i, err := buf.ReadVarInt()
	if err != nil {
		return 0, err
	}
	size := int(i)
	if size > MaxElemSize {
		return 0, errors.New("collection elem size overflow. size:" + strconv.Itoa(size))
	}
	return size, nil
}

// ReadPacked read packed map or packed array without type
func ReadPacked(buf *Buffer, size int, isMap bool, f ReadElemFunc) (err error) {
	if size <= 0 {
		return nil
	}
	readType(buf) // skip key type
	if isMap {
		readType(buf) // skip value type
	}
	for i := 0; i < size; i++ {
		err = f(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadByEnum read enum with type
func ReadByEnum(buf *Buffer, enum Enum, asAddr bool) (interface{}, error) {
	tp, _, err := readType(buf)
	if err != nil {
		return nil, err
	}
	if tp != MessageType {
		return nil, errors.New("ReadByEnum fail, type not message, tp " + strconv.Itoa(int(tp)))
	}
	return enum.ReadEnum(buf, asAddr)
}

// ReadByMessage read message with type
func ReadByMessage(buf *Buffer, msg Message) (err error) {
	tp, _, err := readType(buf)
	if err != nil {
		return err
	}
	if tp != MessageType {
		return errors.New("ReadByEnum fail, type not message, tp " + strconv.Itoa(int(tp)))
	}
	return msg.ReadFrom(buf)
}

// ReadStringStringMap read map[string]string
func ReadStringStringMap(buf *Buffer, withType bool) (m map[string]string, err error) {
	size, err := ReadPackedSize(buf, withType)
	if err != nil {
		return nil, err
	}
	m = make(map[string]string, size)
	err = ReadPacked(buf, size, true, func(buf *Buffer) (err error) {
		var k, v string
		k, err = ReadStringWithoutType(buf)
		if err != nil {
			return err
		}
		v, err = ReadStringWithoutType(buf)
		if err == nil {
			m[k] = v
		}
		return err
	})
	return m, err
}

// ReadStringInt32Map read map[string]int32
func ReadStringInt32Map(buf *Buffer, withType bool) (m map[string]int32, err error) {
	size, err := ReadPackedSize(buf, withType)
	if err != nil {
		return nil, err
	}
	m = make(map[string]int32, size)
	err = ReadPacked(buf, size, true, func(buf *Buffer) error {
		k, err := ReadStringWithoutType(buf)
		if err != nil {
			return err
		}
		v, err := ReadInt32WithoutType(buf)
		if err == nil {
			m[k] = v
		}
		return err
	})
	return m, err
}

// ReadStringInt64Map read map[string]int64
func ReadStringInt64Map(buf *Buffer, withType bool) (m map[string]int64, err error) {
	size, err := ReadPackedSize(buf, withType)
	if err != nil {
		return nil, err
	}
	m = make(map[string]int64, size)
	err = ReadPacked(buf, size, true, func(buf *Buffer) error {
		k, err := ReadStringWithoutType(buf)
		if err != nil {
			return err
		}
		v, err := ReadInt64WithoutType(buf)
		if err == nil {
			m[k] = v
		}
		return err
	})
	return m, err
}

// ReadStringArray read []string
func ReadStringArray(buf *Buffer, withType bool) ([]string, error) {
	size, err := ReadPackedSize(buf, withType)
	if err != nil {
		return nil, err
	}
	a := make([]string, 0, size)
	if size > 0 {
		readType(buf)
		for i := 0; i < size; i++ {
			s, err := ReadStringWithoutType(buf)
			if err != nil {
				return nil, err
			}
			a = append(a, s)
		}
	}
	return a, nil
}

// ReadInt32Array read []int32
func ReadInt32Array(buf *Buffer, withType bool) ([]int32, error) {
	size, err := ReadPackedSize(buf, withType)
	if err != nil {
		return nil, err
	}
	a := make([]int32, 0, size)
	if size > 0 {
		readType(buf)
		for i := 0; i < size; i++ {
			i32, err := ReadInt32WithoutType(buf)
			if err != nil {
				return nil, err
			}
			a = append(a, i32)
		}
	}
	return a, nil
}

// ReadInt64Array read []int64
func ReadInt64Array(buf *Buffer, withType bool) ([]int64, error) {
	size, err := ReadPackedSize(buf, withType)
	if err != nil {
		return nil, err
	}
	a := make([]int64, 0, size)
	if size > 0 {
		readType(buf)
		for i := 0; i < size; i++ {
			i64, err := ReadInt64WithoutType(buf)
			if err != nil {
				return nil, err
			}
			a = append(a, i64)
		}
	}
	return a, nil
}

func readType(buf *Buffer) (tp byte, name string, err error) {
	tp, err = buf.ReadByte()
	if err != nil {
		return tp, name, err
	}
	if tp >= MessageType { //message
		name, err = readMessageType(buf, tp)
		tp = MessageType
	}
	return tp, name, err
}

func readMessageType(buf *Buffer, tp byte) (name string, err error) {
	if tp == MessageType {
		name, err = ReadStringWithoutType(buf)
		if err == nil {
			buf.GetContext().putMessageType(name)
		}
	} else if tp == RefMessageType {
		index, err := buf.ReadVarInt()
		if err != nil {
			return name, err
		}
		name = buf.GetContext().getMessageTypeName(int(index))
	} else {
		name = buf.GetContext().getMessageTypeName(int(tp - RefMessageType))
	}
	return name, err
}

// ReadMessageField can read all message fields according to the ReadFieldsFunc
func ReadMessageField(buf *Buffer, readField ReadFieldsFunc) (err error) {
	total, err := buf.ReadInt()
	if err != nil {
		return err
	}
	if total > 0 {
		pos := buf.GetRPos()
		endPos := pos + total
		var index uint64
		for buf.GetRPos() < endPos {
			index, err = buf.ReadVarInt()
			if err != nil {
				return err
			}
			err = readField(buf, int(index))
			if err != nil {
				return err
			}
		}
		if buf.GetRPos() != endPos {
			return ErrWrongSize
		}
	}
	return nil
}

func checkType(buf *Buffer, expect byte) error {
	tp, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if tp != expect {
		return errors.New("BreezeRead: wrong type expect " + strconv.Itoa(int(expect)) + ", real " + strconv.Itoa(int(tp)))
	}
	return nil
}

/*
ReadValue read A value from Buffer based v.
v can be A reflect type or an address which receive the deserialize value.
*/
func ReadValue(buf *Buffer, v interface{}) (interface{}, error) {
	return readValueByType(buf, v, true, 0, "")
}

func readValueByType(buf *Buffer, v interface{}, withType bool, t byte, msgName string) (interface{}, error) {
	var err error
	if withType {
		t, msgName, err = readType(buf)
		if err != nil {
			return nil, err
		}
	}

	// string
	if t <= StringType {
		var s string
		if t == StringType {
			s, err = ReadStringWithoutType(buf)
		} else {
			bytes, err := buf.Next(int(t))
			if err != nil {
				return nil, err
			}
			s = string(bytes)
		}
		if err != nil {
			return nil, err
		}
		return adaptString(s, v)
	}
	// int32
	if t >= DirectInt32MinType && t <= Int32Type {
		var i32 int32
		if t == Int32Type {
			i32, err = ReadInt32WithoutType(buf)
		} else {
			i32 = int32(t) - Int32Zero
		}
		if err != nil {
			return nil, err
		}
		return adaptInt32(i32, v)
	}
	//int64
	if t >= DirectInt64MinType && t <= Int64Type {
		var i64 int64
		if t == Int64Type {
			i64, err = ReadInt64WithoutType(buf)
		} else {
			i64 = int64(t) - Int64Zero
		}
		if err != nil {
			return nil, err
		}
		return adaptInt64(i64, v)
	}

	switch t {
	case NullType:
		return nil, nil
	case MessageType:
		return readMessage(buf, v, msgName)
	case MapType, PackedMapType:
		return readMap(buf, v, t == PackedMapType)
	case ArrayType, PackedArrayType:
		return readArray(buf, v, t == PackedArrayType)
	case TrueType:
		if castV, ok := v.(*bool); ok {
			*castV = true
		}
		return true, nil
	case FalseType:
		if castV, ok := v.(*bool); ok {
			*castV = false
		}
		return false, nil
	case Float32Type:
		f, err := ReadFloat32WithoutType(buf)
		if err != nil {
			return nil, err
		}
		return adaptFloat32(f, v)
	case Float64Type:
		f, err := ReadFloat64WithoutType(buf)
		if err != nil {
			return nil, err
		}
		return adaptFloat64(f, v)
	case BytesType:
		bytes, err := ReadBytesWithoutType(buf)
		if err != nil {
			return nil, err
		}
		if castV, ok := v.(*[]byte); ok {
			*castV = bytes
			return *castV, nil
		}
		return bytes, nil
	case ByteType:
		return adaptByte(buf, v)
	case Int16Type:
		i16, err := ReadInt16WithoutType(buf)
		if err != nil {
			return nil, err
		}
		return adaptInt16(i16, v)
	}
	return nil, errors.New("BreezeRead: unsupported type " + strconv.Itoa(int(t)))
}

func readMessage(buf *Buffer, v interface{}, name string) (interface{}, error) {
	if enum, ok := v.(Enum); ok {
		result, err := enum.ReadEnum(buf, false)
		if err == nil {
			rv := reflect.ValueOf(v)
			if rv.Type().Kind() == reflect.Ptr {
				rv.Elem().Set(reflect.ValueOf(result))
			}
		}
		return result, err
	}
	message, ok := v.(Message)
	if ok {
		if message.GetName() != name && message.GetAlias() != name {
			return nil, errors.New("BreezeRead: wrong message type. expect " + message.GetName() + ", real " + name)
		}
	} else if v == nil || reflect.TypeOf(v).Kind() == reflect.Interface {
		message = &GenericMessage{Name: name}
	} else if rt, isType := v.(reflect.Type); isType {
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		if rt.Kind() == reflect.Interface {
			message = &GenericMessage{Name: name}
		} else {
			newValue := reflect.New(rt).Interface()
			if enum, ok := newValue.(Enum); ok {
				return enum.ReadEnum(buf, true)
			}
			message, ok = newValue.(Message)
		}
	}
	if message != nil {
		err := message.ReadFrom(buf)
		if err != nil {
			return nil, err
		}
		return message, nil
	}
	return nil, errors.New("BreezeRead: can not read breeze message to type" + reflect.TypeOf(v).String())
}

func readArray(buf *Buffer, v interface{}, isPacked bool) (interface{}, error) {
	total, err := buf.ReadVarInt()
	if err != nil {
		return nil, err
	}
	if total <= 0 {
		return nil, nil
	}
	size := int(total)
	if size > MaxElemSize {
		return nil, errors.New("collection elem size overflow. size:" + strconv.Itoa(size))
	}
	var orgRv reflect.Value
	var rv reflect.Value
	rt, isType := v.(reflect.Type)
	if !isType {
		if v == nil {
			tmp := make([]interface{}, 0, size)
			rv = reflect.ValueOf(&tmp)
		} else {
			rv = reflect.ValueOf(v)
		}
		rt = rv.Type()
		if rt.Kind() != reflect.Ptr {
			return nil, errors.New("BreezeRead: can not read slice to type " + rt.String())
		}
		orgRv = rv
		rv = rv.Elem()
		rt = rv.Type()
	}
	if rt.Kind() != reflect.Slice && rt.Kind() != reflect.Interface {
		return nil, errors.New("BreezeRead: can not read slice to type " + rt.String())
	}
	if isType {
		if rt.Kind() == reflect.Interface {
			rv = reflect.ValueOf(make([]interface{}, 0, size))
			rt = rv.Type()
		} else {
			rv = reflect.MakeSlice(rt, 0, size)
		}
	}

	var tp byte
	var name string
	if isPacked {
		tp, name, err = readType(buf)
		if err != nil {
			return nil, err
		}
	}

	var sv interface{}
	for i := 0; i < size; i++ {
		if isPacked {
			sv, err = readValueByType(buf, rt.Elem(), false, tp, name)
		} else {
			sv, err = ReadValue(buf, rt.Elem())
		}
		if err != nil {
			return nil, err
		}
		rv = reflect.Append(rv, reflect.ValueOf(sv))
	}
	if !isType {
		orgRv.Elem().Set(rv)
	}
	return rv.Interface(), nil
}

func readMap(buf *Buffer, v interface{}, isPacked bool) (interface{}, error) {
	total, err := buf.ReadVarInt()
	if err != nil {
		return nil, err
	}
	if total <= 0 {
		return nil, nil
	}
	size := int(total)
	if size > MaxElemSize {
		return nil, errors.New("collection elem size overflow. size:" + strconv.Itoa(size))
	}
	var orgRv reflect.Value
	var rv reflect.Value
	rt, isType := v.(reflect.Type)
	if !isType {
		if v == nil {
			tmp := make(map[interface{}]interface{}, size)
			rv = reflect.ValueOf(&tmp)
		} else {
			rv = reflect.ValueOf(v)
		}
		rt = rv.Type()
		if rt.Kind() != reflect.Ptr {
			return nil, errors.New("BreezeRead: can not read map to type " + rt.String())
		}
		orgRv = rv
		rv = rv.Elem()
		rt = rv.Type()
	}
	if rt.Kind() != reflect.Map && rt.Kind() != reflect.Interface {
		return nil, errors.New("BreezeRead: can not read map to type " + rt.String())
	}
	if isType {
		if rt.Kind() == reflect.Interface {
			rv = reflect.ValueOf(make(map[interface{}]interface{}, size))
			rt = rv.Type()
		} else {
			rv = reflect.MakeMapWithSize(rt, size)
		}
	}
	var ktp, vtp byte
	var kn, vn string
	if isPacked {
		ktp, kn, err = readType(buf)
		vtp, vn, err = readType(buf)
		if err != nil {
			return nil, err
		}
	}

	var mk, mv interface{}
	for i := 0; i < size; i++ {
		if isPacked {
			mk, err = readValueByType(buf, rt.Key(), false, ktp, kn)
			if err != nil {
				return nil, err
			}
			mv, err = readValueByType(buf, rt.Elem(), false, vtp, vn)
			if err != nil {
				return nil, err
			}
		} else {
			mk, err = ReadValue(buf, rt.Key())
			if err != nil {
				return nil, err
			}
			mv, err = ReadValue(buf, rt.Elem())
			if err != nil {
				return nil, err
			}
		}
		rv.SetMapIndex(reflect.ValueOf(mk), reflect.ValueOf(mv))
	}
	if !isType {
		orgRv.Elem().Set(rv)
	}
	return rv.Interface(), nil
}

// ReadFloat64WithoutType read without type
func ReadFloat64WithoutType(buf *Buffer) (float64, error) {
	i, err := buf.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(i), nil

}

func adaptFloat64(f float64, v interface{}) (interface{}, error) {
	if v == nil {
		return f, nil
	}
	if castV, ok := v.(*float64); ok {
		*castV = f
		return *castV, nil
	}
	rt, isType := v.(reflect.Type)
	if isType && (rt.Kind() == reflect.Float64 || rt.Kind() == reflect.Interface) {
		return f, nil
	}
	if !isType {
		rt = reflect.TypeOf(v)
	}
	return nil, errors.New("BreezeRead: can not read float64 to type " + rt.String())
}

// ReadFloat32WithoutType read without type
func ReadFloat32WithoutType(buf *Buffer) (float32, error) {
	i, err := buf.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(i), nil
}

func adaptFloat32(f float32, v interface{}) (interface{}, error) {
	if v == nil {
		return f, nil
	}
	if castV, ok := v.(*float32); ok {
		*castV = f
		return *castV, nil
	}
	rt, isType := v.(reflect.Type)
	if isType && (rt.Kind() == reflect.Float32 || rt.Kind() == reflect.Interface) {
		return f, nil
	}
	if !isType {
		rt = reflect.TypeOf(v)
	}
	return nil, errors.New("BreezeRead: can not read float32 to type " + rt.String())
}

// ReadInt64WithoutType read without type
func ReadInt64WithoutType(buf *Buffer) (int64, error) {
	i, err := buf.ReadZigzag64()
	return int64(i), err
}

func adaptInt64(i int64, v interface{}) (interface{}, error) {
	if v == nil {
		return i, nil
	}
	if castV, ok := v.(*int64); ok {
		*castV = i
		return *castV, nil
	}
	rt, isType := v.(reflect.Type)
	if isType && (rt.Kind() == reflect.Int64 || rt.Kind() == reflect.Interface) {
		return i, nil
	}
	return adaptToInt(i, v)
}

// ReadInt32WithoutType read without type
func ReadInt32WithoutType(buf *Buffer) (int32, error) {
	i, err := buf.ReadZigzag32()
	return int32(i), err
}

func adaptInt32(i int32, v interface{}) (interface{}, error) {
	if v == nil {
		return i, nil
	}
	if castV, ok := v.(*int); ok {
		*castV = int(i)
		return *castV, nil
	}
	if castV, ok := v.(*int32); ok {
		*castV = i
		return *castV, nil
	}
	rt, isType := v.(reflect.Type)
	if isType {
		if rt.Kind() == reflect.Int || rt.Kind() == reflect.Interface {
			return int(i), nil
		} else if rt.Kind() == reflect.Int32 {
			return i, nil
		}
	}
	return adaptToInt(int64(i), v)
}

// ReadInt16WithoutType read without type
func ReadInt16WithoutType(buf *Buffer) (int16, error) {
	i, err := buf.ReadUint16()
	return int16(i), err
}

func adaptInt16(i int16, v interface{}) (interface{}, error) {
	if v == nil {
		return i, nil
	}
	if castV, ok := v.(*int16); ok {
		*castV = i
		return *castV, nil
	}
	rt, isType := v.(reflect.Type)
	if isType && (rt.Kind() == reflect.Int16 || rt.Kind() == reflect.Interface) {
		return i, nil
	}
	return adaptToInt(int64(i), v)
}

// ReadBytesWithoutType read without type
func ReadBytesWithoutType(buf *Buffer) ([]byte, error) {
	size, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}
	ret := make([]byte, size)
	err = buf.ReadFull(ret)
	return ret, err
}

func adaptByte(buf *Buffer, v interface{}) (interface{}, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if castV, ok := v.(*byte); ok {
		*castV = b
	}
	return b, nil
}

// ReadStringWithoutType read without type
func ReadStringWithoutType(buf *Buffer) (string, error) {
	size, err := buf.ReadVarInt()
	if err != nil {
		return "", err
	}
	bytes, err := buf.Next(int(size))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func adaptString(s string, v interface{}) (interface{}, error) {
	if v == nil {
		return s, nil
	}
	if castV, ok := v.(*string); ok {
		*castV = s
		return s, nil
	}
	rt, isType := v.(reflect.Type)
	if isType && (rt.Kind() == reflect.String || rt.Kind() == reflect.Interface) {
		return s, nil
	}
	//compatible with other type
	if isType {
		return parseStringByType(s, rt)
	}
	rv := reflect.ValueOf(v)
	if rv.CanSet() && rv.Type().Kind() == reflect.Ptr {
		tmp, err := parseStringByType(s, rv.Type().Elem())
		if err != nil {
			return nil, err
		}
		rv.Set(reflect.ValueOf(tmp))
		return tmp, nil
	}
	return nil, errors.New("can not read string to type " + rv.Type().String())
}

func parseStringByType(s string, rt reflect.Type) (interface{}, error) {
	switch rt.Kind() {
	case reflect.Slice:
		if rt.Elem().Kind() == reflect.Uint8 {
			return []byte(s), nil
		}
	case reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int, reflect.Uint, reflect.Int64, reflect.Uint64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return getIntByKind(i, rt.Kind())
	case reflect.Float32:
		return strconv.ParseFloat(s, 32)
	case reflect.Float64:
		return strconv.ParseFloat(s, 64)
	}
	return nil, errors.New("BreezeRead: can not read string to type " + rt.String())
}

func adaptToInt(i int64, v interface{}) (interface{}, error) {
	if rt, isType := v.(reflect.Type); isType {
		return getIntByKind(int64(i), rt.Kind())
	}
	rv := reflect.ValueOf(v)
	if rv.CanSet() && rv.Type().Kind() == reflect.Ptr {
		tmp, err := getIntByKind(int64(i), rv.Type().Elem().Kind())
		if err != nil {
			return nil, err
		}
		rv.Set(reflect.ValueOf(tmp))
		return tmp, nil
	}
	return nil, errors.New("BreezeRead: can not read int to type " + rv.Type().String())
}

func getIntByKind(i int64, k reflect.Kind) (interface{}, error) {
	switch k {
	case reflect.Int16:
		return int16(i), nil
	case reflect.Uint16:
		return uint16(i), nil
	case reflect.Int32:
		return int32(i), nil
	case reflect.Uint32:
		return uint32(i), nil
	case reflect.Int: // for compatible with other language, int is regarded as int32 in breeze. u should use int64 if value over int32
		return int(int32(i)), nil
	case reflect.Uint:
		return uint(i), nil
	case reflect.Int64:
		return i, nil
	case reflect.Uint64:
		return uint64(i), nil
	default:
		return nil, errors.New("BreezeRead: can not convert to type " + k.String())
	}
}

package breeze

import (
	"github.com/pkg/errors"
	"math"
	"reflect"
	"strconv"
)

const (
	DefaultSize = 16
)

type ReadFieldsFunc func(buf *Buffer, index int) error

func ReadBool(buf *Buffer, b *bool) error {
	tp, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if tp == TRUE {
		*b = true
		return nil
	} else if tp == FALSE {
		*b = false
		return nil
	} else {
		return errors.New("BreezeRead: wrong type expect bool, real " + strconv.Itoa(int(tp)))
	}
}

func ReadString(buf *Buffer, s *string) error {
	err := checkType(buf, STRING)
	if err != nil {
		return err
	}
	_, err = readStringWithoutType(buf, s)
	return err
}

func ReadByte(buf *Buffer, b *byte) error {
	err := checkType(buf, BYTE)
	if err != nil {
		return err
	}
	_, err = readByteWithoutType(buf, b)
	return err
}

func ReadBytes(buf *Buffer, bytesAddr *[]byte) error {
	err := checkType(buf, BYTES)
	if err != nil {
		return err
	}
	_, err = readBytesWithoutType(buf, bytesAddr)
	return err
}

func ReadInt16(buf *Buffer, i *int16) error {
	err := checkType(buf, INT16)
	if err != nil {
		return err
	}
	_, err = readInt16WithoutType(buf, i)
	return err
}

func ReadInt(buf *Buffer, i *int) error {
	err := checkType(buf, INT32)
	if err != nil {
		return err
	}
	_, err = readInt32WithoutType(buf, i)
	return err
}

func ReadInt64(buf *Buffer, i *int64) error {
	err := checkType(buf, INT64)
	if err != nil {
		return err
	}
	_, err = readInt64WithoutType(buf, i)
	return err
}

func ReadFloat32(buf *Buffer, f *float32) error {
	err := checkType(buf, FLOAT32)
	if err != nil {
		return err
	}
	_, err = readFloat32WithoutType(buf, f)
	return err
}

func ReadFloat64(buf *Buffer, f *float64) error {
	err := checkType(buf, FLOAT64)
	if err != nil {
		return err
	}
	_, err = readFloat64WithoutType(buf, f)
	return err
}

func ReadMessageByField(buf *Buffer, readField ReadFieldsFunc) error {
	total, err := buf.ReadInt()
	if err != nil {
		return err
	}
	if total > 0 {
		pos := buf.GetRPos()
		endPos := pos + total
		var index uint64
		for buf.GetRPos() < endPos {
			index, err = buf.ReadZigzag32()
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
v can be A reflect type or an address which recieve the deserialized value.
*/
func ReadValue(buf *Buffer, v interface{}) (interface{}, error) {
	tp, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	switch int(tp) {
	case NULL:
		return nil, nil
	case TRUE:
		if castV, ok := v.(*bool); ok {
			*castV = true
		}
		return true, nil
	case FALSE:
		if castV, ok := v.(*bool); ok {
			*castV = false
		}
		return false, nil
	case STRING:
		return readStringWithoutType(buf, v)
	case BYTE:
		return readByteWithoutType(buf, v)
	case BYTES:
		return readBytesWithoutType(buf, v)
	case INT16:
		return readInt16WithoutType(buf, v)
	case INT32:
		return readInt32WithoutType(buf, v)
	case INT64:
		return readInt64WithoutType(buf, v)
	case FLOAT32:
		return readFloat32WithoutType(buf, v)
	case FLOAT64:
		return readFloat64WithoutType(buf, v)
	case MAP:
		return readMap(buf, v)
	case ARRAY:
		return readArray(buf, v)
	case MESSAGE:
		return readMessage(buf, v)
	}
	return nil, errors.New("BreezeRead: unsupport type " + strconv.Itoa(int(tp)))
}

func readMessage(buf *Buffer, v interface{}) (interface{}, error) {
	var name string
	err := ReadString(buf, &name)
	if err != nil {
		return nil, err
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
			message, _ = reflect.New(rt).Interface().(Message)
		}
	}
	if message != nil {
		err = message.ReadFrom(buf)
		if err != nil {
			return nil, err
		}
		return message, nil
	}
	return nil, errors.New("BreezeRead: can not read breeze message to type" + reflect.TypeOf(v).String())
}

func readArray(buf *Buffer, v interface{}) (interface{}, error) {
	total, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}
	if total <= 0 {
		return nil, nil
	}
	var orgRv reflect.Value
	var rv reflect.Value
	rt, isType := v.(reflect.Type)
	if !isType {
		if v == nil {
			tmp := make([]interface{}, 0, DefaultSize)
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
			rv = reflect.ValueOf(make([]interface{}, 0, DefaultSize))
			rt = rv.Type()
		} else {
			rv = reflect.MakeSlice(rt, 0, DefaultSize)
		}
	}
	pos := buf.GetRPos()
	endPos := pos + total
	var sv interface{}
	for buf.GetRPos() < endPos {
		sv, err = ReadValue(buf, rt.Elem())
		if err != nil {
			return nil, err
		}
		rv = reflect.Append(rv, reflect.ValueOf(sv))
	}
	if buf.GetRPos() != endPos {
		return nil, ErrWrongSize
	}
	if !isType {
		orgRv.Elem().Set(rv)
	}
	return rv.Interface(), nil
}

func readMap(buf *Buffer, v interface{}) (interface{}, error) {
	total, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}
	if total <= 0 {
		return nil, nil
	}
	var rv reflect.Value
	rt, isType := v.(reflect.Type)
	if !isType {
		if v == nil {
			tmp := make(map[interface{}]interface{}, DefaultSize)
			rv = reflect.ValueOf(&tmp)
		} else {
			rv = reflect.ValueOf(v)
		}
		rt = rv.Type()
		if rt.Kind() != reflect.Ptr {
			return nil, errors.New("BreezeRead: can not read map to type " + rt.String())
		}
		rv = rv.Elem()
		rt = rv.Type()
	}
	if rt.Kind() != reflect.Map && rt.Kind() != reflect.Interface {
		return nil, errors.New("BreezeRead: can not read map to type " + rt.String())
	}
	if isType {
		if rt.Kind() == reflect.Interface {
			rv = reflect.ValueOf(make(map[interface{}]interface{}, DefaultSize))
			rt = rv.Type()
		} else {
			rv = reflect.MakeMapWithSize(rt, DefaultSize)
		}
	}
	pos := buf.GetRPos()
	endPos := pos + total
	var mk, mv interface{}
	for buf.GetRPos() < endPos {
		mk, err = ReadValue(buf, rt.Key())
		if err != nil {
			return nil, err
		}
		mv, err = ReadValue(buf, rt.Elem())
		if err != nil {
			return nil, err
		}
		rv.SetMapIndex(reflect.ValueOf(mk), reflect.ValueOf(mv))
	}
	if buf.GetRPos() != endPos {
		return nil, ErrWrongSize
	}
	return rv.Interface(), nil
}

func readFloat64WithoutType(buf *Buffer, v interface{}) (interface{}, error) {
	i, err := buf.ReadUint64()
	if err != nil {
		return 0, err
	}
	f := math.Float64frombits(i)
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

func readFloat32WithoutType(buf *Buffer, v interface{}) (interface{}, error) {
	i, err := buf.ReadUint32()
	if err != nil {
		return 0, err
	}
	f := math.Float32frombits(i)
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

func readInt64WithoutType(buf *Buffer, v interface{}) (interface{}, error) {
	i, err := buf.ReadZigzag64()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return int64(i), nil
	}
	if castV, ok := v.(*int64); ok {
		*castV = int64(i)
		return *castV, nil
	}
	rt, isType := v.(reflect.Type)
	if isType && (rt.Kind() == reflect.Int64 || rt.Kind() == reflect.Interface) {
		return int64(i), nil
	}
	return compatInt(int64(i), v)
}

func readInt32WithoutType(buf *Buffer, v interface{}) (interface{}, error) {
	i, err := buf.ReadZigzag32()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return int32(i), nil
	}
	if castV, ok := v.(*int); ok {
		*castV = int(i)
		return *castV, nil
	}
	if castV, ok := v.(*int32); ok {
		*castV = int32(i)
		return *castV, nil
	}
	rt, isType := v.(reflect.Type)
	if isType {
		if rt.Kind() == reflect.Int || rt.Kind() == reflect.Interface {
			return int(i), nil
		} else if rt.Kind() == reflect.Int32 {
			return int32(i), nil
		}
	}
	return compatInt(int64(i), v)
}

func readInt16WithoutType(buf *Buffer, v interface{}) (interface{}, error) {
	i, err := buf.ReadUint16()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return int16(i), nil
	}
	if castV, ok := v.(*int16); ok {
		*castV = int16(i)
		return *castV, nil
	}
	rt, isType := v.(reflect.Type)
	if isType && (rt.Kind() == reflect.Int16 || rt.Kind() == reflect.Interface) {
		return int16(i), nil
	}
	return compatInt(int64(i), v)
}

func readBytesWithoutType(buf *Buffer, v interface{}) (interface{}, error) {
	size, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}
	ret := make([]byte, size)
	err = buf.ReadFull(ret)
	if err != nil {
		return nil, err
	}
	if castV, ok := v.(*[]byte); ok {
		*castV = ret
		return *castV, nil
	}
	return ret, nil
}

func readByteWithoutType(buf *Buffer, v interface{}) (interface{}, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if castV, ok := v.(*byte); ok {
		*castV = b
	}
	return b, nil
}

func readStringWithoutType(buf *Buffer, v interface{}) (interface{}, error) {
	size, err := buf.ReadZigzag32()
	if err != nil {
		return nil, err
	}
	bytes, err := buf.Next(int(size))
	if err != nil {
		return nil, err
	}
	s := string(bytes)
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

func compatInt(i int64, v interface{}) (interface{}, error) {
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
	case reflect.Int:
		return int(i), nil
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

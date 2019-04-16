package breeze

import (
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestWriteBool(t *testing.T) {
	type args struct {
		buf *Buffer
		b   bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"true", args{NewBuffer(32), true}},
		{"false", args{NewBuffer(32), false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteBool(tt.args.buf, tt.args.b)
			if tt.args.buf.Len() != 1 {
				t.Errorf("wrong write size. expect %d, real %d", 1, tt.args.buf.Len())
			}
			bytes := tt.args.buf.Bytes()
			var b bool
			err := ReadBool(CreateBuffer(bytes), &b)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if b != tt.args.b {
				t.Errorf("wrong result. expect %v, real %v", tt.args.b, b)
			}
		})
	}
}

func TestWriteString(t *testing.T) {
	type args struct {
		buf *Buffer
		s   string
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty", args{NewBuffer(32), ""}},
		{"string", args{NewBuffer(32), "uwoerj8093lsd#!@#$%^^&&*()lkd"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteString(tt.args.buf, tt.args.s)
			bytes := tt.args.buf.Bytes()
			var s string
			err := ReadString(CreateBuffer(bytes), &s)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if s != tt.args.s {
				t.Errorf("wrong result. expect %v, real %v", tt.args.s, s)
			}
		})
	}
}

func TestWriteByte(t *testing.T) {
	type args struct {
		buf *Buffer
		b   byte
	}
	tests := []struct {
		name string
		args args
	}{
		{"zero", args{NewBuffer(32), 0}},
		{"maxint", args{NewBuffer(32), math.MaxInt8}},
		{"maxuint", args{NewBuffer(32), math.MaxUint8}},
		{"normal", args{NewBuffer(32), 36}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteByte(tt.args.buf, tt.args.b)
			if tt.args.buf.Len() != 2 {
				t.Errorf("wrong write size. expect %d, real %d", 2, tt.args.buf.Len())
			}
			bytes := tt.args.buf.Bytes()
			var b byte
			err := ReadByte(CreateBuffer(bytes), &b)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if b != tt.args.b {
				t.Errorf("wrong result. expect %v, real %v", tt.args.b, b)
			}
		})
	}
}

func TestWriteBytes(t *testing.T) {
	type args struct {
		buf   *Buffer
		bytes []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty", args{NewBuffer(32), make([]byte, 0)}},
		{"normal", args{NewBuffer(32), []byte("jlkw!@#%$#%#$%hjsde23kd\\n\\t")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteBytes(tt.args.buf, tt.args.bytes)
			bytes := tt.args.buf.Bytes()
			var newBytes []byte
			err := ReadBytes(CreateBuffer(bytes), &newBytes)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if len(newBytes) != len(tt.args.bytes) {
				t.Errorf("wrong result. expect %v, real %v", tt.args.bytes, newBytes)
			}
			if !reflect.DeepEqual(newBytes, tt.args.bytes) {
				t.Errorf("wrong result. expect %v, real %v", tt.args.bytes, newBytes)
			}
		})
	}
}

func TestWriteInt16(t *testing.T) {
	type args struct {
		buf *Buffer
		i   int16
	}
	tests := []struct {
		name string
		args args
	}{
		{"zero", args{NewBuffer(32), 0}},
		{"negative", args{NewBuffer(32), -123}},
		{"max", args{NewBuffer(32), math.MaxInt16}},
		{"min", args{NewBuffer(32), math.MinInt16}},
		{"normal", args{NewBuffer(32), 234}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteInt16(tt.args.buf, uint16(tt.args.i))
			if tt.args.buf.Len() != 3 {
				t.Errorf("wrong write size. expect %d, real %d", 3, tt.args.buf.Len())
			}
			bytes := tt.args.buf.Bytes()
			var i int16
			err := ReadInt16(CreateBuffer(bytes), &i)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if i != tt.args.i {
				t.Errorf("wrong result. expect %v, real %v", tt.args.i, i)
			}
		})
	}
}

func TestWriteInt32(t *testing.T) {
	type args struct {
		buf *Buffer
		i   int32
	}
	tests := []struct {
		name string
		args args
	}{
		{"zero", args{NewBuffer(32), 0}},
		{"negative", args{NewBuffer(32), -12345}},
		{"max", args{NewBuffer(32), math.MaxInt32}},
		{"min", args{NewBuffer(32), math.MinInt32}},
		{"normal", args{NewBuffer(32), 2332454}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteInt32(tt.args.buf, uint32(tt.args.i))
			bytes := tt.args.buf.Bytes()
			var i int
			err := ReadInt(CreateBuffer(bytes), &i)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if int32(i) != tt.args.i {
				t.Errorf("wrong result. expect %v, real %v", tt.args.i, i)
			}
		})
	}
}

func TestWriteInt64(t *testing.T) {
	type args struct {
		buf *Buffer
		i   int64
	}
	tests := []struct {
		name string
		args args
	}{
		{"zero", args{NewBuffer(32), 0}},
		{"negative", args{NewBuffer(32), -123234092342345}},
		{"max", args{NewBuffer(32), math.MaxInt64}},
		{"min", args{NewBuffer(32), math.MinInt64}},
		{"normal", args{NewBuffer(32), 23324542384092384}},
		{"normal2", args{NewBuffer(32), 234}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteInt64(tt.args.buf, uint64(tt.args.i))
			bytes := tt.args.buf.Bytes()
			var i int64
			err := ReadInt64(CreateBuffer(bytes), &i)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if i != tt.args.i {
				t.Errorf("wrong result. expect %v, real %v", tt.args.i, i)
			}
		})
	}
}

func TestWriteFloat32(t *testing.T) {
	type args struct {
		buf *Buffer
		f   float32
	}
	tests := []struct {
		name string
		args args
	}{
		{"zero", args{NewBuffer(32), 0}},
		{"negative", args{NewBuffer(32), -123.234}},
		{"max", args{NewBuffer(32), math.MaxFloat32}},
		{"normal", args{NewBuffer(32), 23.324542384092384}},
		{"normal2", args{NewBuffer(32), 23423749823749.45}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteFloat32(tt.args.buf, tt.args.f)
			if tt.args.buf.Len() != 5 {
				t.Errorf("wrong write size. expect %d, real %d", 5, tt.args.buf.Len())
			}
			bytes := tt.args.buf.Bytes()
			var f float32
			err := ReadFloat32(CreateBuffer(bytes), &f)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if f != tt.args.f {
				t.Errorf("wrong result. expect %v, real %v", tt.args.f, f)
			}
		})
	}
}

func TestWriteFloat64(t *testing.T) {
	type args struct {
		buf *Buffer
		f   float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"zero", args{NewBuffer(32), 0}},
		{"negative", args{NewBuffer(32), -122343.234}},
		{"max", args{NewBuffer(32), math.MaxFloat64}},
		{"normal", args{NewBuffer(32), 233480.324542384092384}},
		{"normal2", args{NewBuffer(32), 23423749823749.452343}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteFloat64(tt.args.buf, tt.args.f)
			if tt.args.buf.Len() != 9 {
				t.Errorf("wrong write size. expect %d, real %d", 9, tt.args.buf.Len())
			}
			bytes := tt.args.buf.Bytes()
			var f float64
			err := ReadFloat64(CreateBuffer(bytes), &f)
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if f != tt.args.f {
				t.Errorf("wrong result. expect %v, real %v", tt.args.f, f)
			}
		})
	}
}

func TestWriteValueBasic(t *testing.T) {
	type args struct {
		buf *Buffer
		v   interface{}
	}
	var b = true
	var s = "ewkleruc8738(&^9?//n"
	var by = byte(16)
	var i16 = int16(234)
	var i32 = int32(2389473)
	var i = -2342
	var i64 = int64(2903402374328432983)
	var f32 = float32(3.1415)
	var f64 = float64(23487924.234823904)
	var byArray = []byte("wioejfn//n?><#@)$%(")
	var sArray = []string{"sjie", "erowir23<&*^", "", "23j8"}
	var m = make(map[string]int, 16)
	m["wjeriew"] = 234
	m[">@#D3"] = 234234
	m["@#>$:P:"] = 98023

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"bool", args{NewBuffer(32), b}, false},
		{"bool_addr", args{NewBuffer(32), &b}, false},
		{"string", args{NewBuffer(32), s}, false},
		{"string_addr", args{NewBuffer(32), &s}, false},
		{"byte", args{NewBuffer(32), by}, false},
		{"byte_addr", args{NewBuffer(32), &by}, false},
		{"int16", args{NewBuffer(32), i16}, false},
		{"int16_addr", args{NewBuffer(32), &i16}, false},
		{"int32", args{NewBuffer(32), i32}, false},
		{"int32_addr", args{NewBuffer(32), &i32}, false},
		{"int", args{NewBuffer(32), i}, false},
		{"int_addr", args{NewBuffer(32), &i}, false},
		{"int64", args{NewBuffer(32), i64}, false},
		{"int64_addr", args{NewBuffer(32), &i64}, false},
		{"float32", args{NewBuffer(32), f32}, false},
		{"float32_addr", args{NewBuffer(32), &f32}, false},
		{"float64", args{NewBuffer(32), f64}, false},
		{"float64_addr", args{NewBuffer(32), &f64}, false},
		{"bytes", args{NewBuffer(32), byArray}, false},
		{"bytes_addr", args{NewBuffer(32), &byArray}, false},
		{"slice", args{NewBuffer(32), sArray}, false},
		{"slice_addr", args{NewBuffer(32), &sArray}, false},
		{"map", args{NewBuffer(32), m}, false},
		{"map_addr", args{NewBuffer(32), &m}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteValue(tt.args.buf, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("WriteValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			bytes := tt.args.buf.Bytes()
			rv := reflect.ValueOf(tt.args.v)
			if rv.Type().Kind() == reflect.Ptr { // basic type do not support read pointer type
				rv = rv.Elem()
			}
			ret, err := ReadValue(CreateBuffer(bytes), rv.Type())
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if reflect.TypeOf(ret) != rv.Type() {
				t.Errorf("wrong result. expect %v, real %v", rv.Type().String(), ret)
			}
			if !reflect.DeepEqual(ret, rv.Interface()) {
				t.Errorf("wrong result. expect %v, real %v", tt.args.v, ret)
			}
		})
	}
}

func TestWriteValueComplex(t *testing.T) {
	type args struct {
		buf *Buffer
		v   interface{}
	}
	m := make(map[string][]map[int]float32, 16)
	mSize := 5
	aSize := 6
	imSize := 8
	for i := 0; i < mSize; i++ {
		a := make([]map[int]float32, 0, 16)
		for j := 0; j < aSize; j++ {
			im := make(map[int]float32, 16)
			for k := 0; k < imSize; k++ {
				im[k] = float32(k+i*j*k) * 0.2
			}
			a = append(a, im)
		}
		m[strconv.Itoa(i)] = a
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"complex map", args{NewBuffer(32), m}, false},
		{"complex map", args{NewBuffer(32), &m}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteValue(tt.args.buf, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("WriteValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			bytes := tt.args.buf.Bytes()
			rv := reflect.ValueOf(tt.args.v)
			if rv.Type().Kind() == reflect.Ptr { // basic type do not support read pointer type
				rv = rv.Elem()
			}
			ret, err := ReadValue(CreateBuffer(bytes), rv.Type())
			if err != nil {
				t.Errorf("err :%s", err.Error())
			}
			if !reflect.DeepEqual(ret, rv.Interface()) {
				t.Errorf("wrong result. expect %v, real %v", tt.args.v, ret)
			}
		})
	}
}

func TestWriteMessage(t *testing.T) {
	type args struct {
		buf        *Buffer
		name       string
		fieldsFunc WriteFieldsFunc
	}
	name := "test message"
	m := getTestSubMsg()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"message", args{NewBuffer(32), name, func(buf *Buffer) {
			WriteMessageField(buf, 1, m.S)
			WriteMessageField(buf, 2, m.I)
			WriteMessageField(buf, 3, m.I64)
			WriteMessageField(buf, 4, m.F32)
			WriteMessageField(buf, 5, m.F64)
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteMessage(tt.args.buf, tt.args.name, tt.args.fieldsFunc); (err != nil) != tt.wantErr {
				t.Errorf("WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			bytes := tt.args.buf.Bytes()
			newBuf := CreateBuffer(bytes)
			b, err := newBuf.ReadByte()
			if err != nil {
				t.Errorf("read Message error = %v", err)
			}
			if b != MESSAGE {
				t.Errorf("read wrong message type.expect:%v, real:%v", MESSAGE, b)
			}
			var rname string
			ReadString(newBuf, &rname)
			if rname != name {
				t.Errorf("read wrong message name. expect:%v, real:%v", name, rname)
			}
			ReadMessageByField(newBuf, func(buf *Buffer, index int) error {
				switch index {
				case 1:
					var v string
					ReadString(buf, &v)
					if v != m.S {
						t.Errorf("read wrong message name. expect:%v, real:%v", m.S, v)
					}
				case 2:
					var v int
					ReadInt(buf, &v)
					if v != m.I {
						t.Errorf("read wrong message name. expect:%v, real:%v", m.I, v)
					}
				case 3:
					var v int64
					ReadInt64(buf, &v)
					if v != m.I64 {
						t.Errorf("read wrong message name. expect:%v, real:%v", m.I64, v)
					}
				case 4:
					var v float32
					ReadFloat32(buf, &v)
					if v != m.F32 {
						t.Errorf("read wrong message name. expect:%v, real:%v", m.F32, v)
					}
				case 5:
					var v float64
					ReadFloat64(buf, &v)
					if v != m.F64 {
						t.Errorf("read wrong message name. expect:%v, real:%v", m.F64, v)
					}
				default:
					t.Errorf("read wrong message index :%v", index)
				}
				return nil
			})
		})
	}
}

func TestWriteValueMessage(t *testing.T) {
	msg := getTestMsg()
	buf := NewBuffer(32)
	err := WriteValue(buf, msg)
	if err != nil {
		t.Errorf("write message err:%v", err)
	}
	bytes := buf.Bytes()
	var result TestMsg
	// read by value
	newBuf := CreateBuffer(bytes)
	ReadValue(newBuf, &result)
	if !reflect.DeepEqual(&result, msg) {
		t.Errorf("wrong result. expect %v, real %v", msg, result)
	}

	// read by type
	newBuf = CreateBuffer(bytes)
	r, _ := ReadValue(newBuf, reflect.TypeOf(&result))
	if !reflect.DeepEqual(&result, msg) {
		t.Errorf("wrong result. expect %v, real %v", msg, result)
	}

	// test GenericMessage (read by nil)
	newBuf = CreateBuffer(bytes)
	r, _ = ReadValue(newBuf, nil)
	gm := r.(*GenericMessage)
	sgm := gm.GetFieldByIndex(3).(map[interface{}]interface{})["m1"].(*GenericMessage)
	if len(sgm.GetFieldByIndex(10).([]interface{})) != len(msg.M["m1"].List) {
		t.Errorf("read wrong message. expect:%v, real:%v", gm, msg)
	}
}

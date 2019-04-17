package breeze

// TestMsg is a breeze message for test
type TestMsg struct {
	I int
	S string
	M map[string]*TestSubMsg
	A []*TestSubMsg
}

// GetAlias return breeze message alias for multi language compat
func (t *TestMsg) GetAlias() string {
	return "TestMsg"
}

// WriteTo write breeze message to breeze buffer.
func (t *TestMsg) WriteTo(buf *Buffer) error {
	return WriteMessage(buf, "TestMsg", func(buffer *Buffer) {
		WriteMessageField(buf, 1, t.I)
		WriteMessageField(buf, 2, t.S)
		WriteMessageField(buf, 3, t.M)
		WriteMessageField(buf, 4, t.A)
	})
}

// ReadFrom read a breeze message from breeze buffer
func (t *TestMsg) ReadFrom(buf *Buffer) error {
	return ReadMessageByField(buf, func(buf *Buffer, index int) (err error) {
		switch index {
		case 1:
			return ReadInt(buf, &t.I)
		case 2:
			return ReadString(buf, &t.S)
		case 3:
			t.M = make(map[string]*TestSubMsg, 16)
			_, err = ReadValue(buf, &t.M)
		case 4:
			t.A = make([]*TestSubMsg, 0, 16)
			_, err = ReadValue(buf, &t.A)
		default: // for compat
			_, err = ReadValue(buf, nil)
		}
		return err
	})
}

// GetName get the name of breeze message
func (t *TestMsg) GetName() string {
	return "TestMsg"
}

// GetSchema get breeze message's schema
func (t *TestMsg) GetSchema() *Schema {
	panic("implement me")
}

// TestSubMsg is a breeze message for test
type TestSubMsg struct {
	S     string
	I     int
	I64   int64
	F32   float32
	F64   float64
	Byte  byte
	Bytes []byte
	Map1  map[string][]byte
	Map2  map[int][]interface{}
	List  []int
	B     bool
}

// GetAlias return breeze message alias for multi language compat
func (t *TestSubMsg) GetAlias() string {
	return "TestSubMsg"
}

// WriteTo write breeze message to breeze buffer.
func (t *TestSubMsg) WriteTo(buf *Buffer) error {
	return WriteMessage(buf, "TestSubMsg", func(buffer *Buffer) {
		WriteMessageField(buf, 1, t.S)
		WriteMessageField(buf, 2, t.I)
		WriteMessageField(buf, 3, t.I64)
		WriteMessageField(buf, 4, t.F32)
		WriteMessageField(buf, 5, t.F64)
		WriteMessageField(buf, 6, t.Byte)
		WriteMessageField(buf, 7, t.Bytes)
		WriteMessageField(buf, 8, t.Map1)
		WriteMessageField(buf, 9, &t.Map2)
		WriteMessageField(buf, 10, t.List)
		WriteMessageField(buf, 11, t.B)
	})
}

// ReadFrom read a breeze message from breeze buffer
func (t *TestSubMsg) ReadFrom(buf *Buffer) error {
	return ReadMessageByField(buf, func(buf *Buffer, index int) (err error) {
		switch index {
		case 1:
			return ReadString(buf, &t.S)
		case 2:
			return ReadInt(buf, &t.I)
		case 3:
			return ReadInt64(buf, &t.I64)
		case 4:
			return ReadFloat32(buf, &t.F32)
		case 5:
			return ReadFloat64(buf, &t.F64)
		case 6:
			return ReadByte(buf, &t.Byte)
		case 7:
			return ReadBytes(buf, &t.Bytes)
		case 8:
			t.Map1 = make(map[string][]byte, 16)
			_, err = ReadValue(buf, &t.Map1)
		case 9:
			t.Map2 = make(map[int][]interface{}, 16)
			_, err = ReadValue(buf, &t.Map2)
		case 10:
			t.List = make([]int, 0, 16)
			_, err = ReadValue(buf, &t.List)
		case 11:
			return ReadBool(buf, &t.B)
		default: // for compat
			_, err = ReadValue(buf, nil)
		}
		return err
	})
}

// GetName get the name of breeze message
func (t *TestSubMsg) GetName() string {
	return "TestSubMsg"
}

// GetSchema get breeze message's schema
func (t *TestSubMsg) GetSchema() *Schema {
	panic("implement me")
}

func getTestMsg() *TestMsg {
	tsm := getTestSubMsg()
	t := &TestMsg{I: 123, S: "jiernoce"}
	t.M = make(map[string]*TestSubMsg)
	t.M["m1"] = tsm
	t.A = make([]*TestSubMsg, 0, 12)
	t.A = append(t.A, tsm)
	return t
}

func getTestSubMsg() *TestSubMsg {
	tsm := &TestSubMsg{S: "uoiwer", I: 2134, I64: 234, F32: 23.434, F64: 8923.234234, Byte: 5, Bytes: []byte("ipower"), B: true}
	im1 := make(map[string][]byte, 16)
	im1["jdie"] = []byte("ierjkkkd")
	im1["jddfwwie"] = []byte("ieere9943rjkkkd")
	tsm.Map1 = im1
	il := make([]interface{}, 0, 12)
	il = append(il, 34)
	il = append(il, 56)
	im2 := make(map[int][]interface{}, 16)
	im2[12] = il
	im2[3] = []interface{}{34, 45, 657}
	im2[6] = []interface{}{23, 66}
	tsm.Map2 = im2
	tsm.List = []int{234, 6456, 234, 6859}
	return tsm
}

package breeze

import (
	"errors"
)

// breeze type
const (
	StringType              = 0x3f
	DirectStringMinType     = 0x00
	DirectStringMaxType     = 0x3e
	Int32Type               = 0x7f
	DirectInt32MinType      = 0x40
	DirectInt32MaxType      = 0x7e
	Int64Type               = 0x98
	DirectInt64MinType      = 0x80
	DirectInt64MaxType      = 0x97
	NullType                = 0x99
	TrueType                = 0x9a
	FalseType               = 0x9b
	ByteType                = 0x9c
	BytesType               = 0x9d
	Int16Type               = 0x9e
	Float32Type             = 0x9f
	Float64Type             = 0xa0
	MapType                 = 0xd9
	ArrayType               = 0xda
	PackedMapType           = 0xdb
	PackedArrayType         = 0xdc
	SchemaType              = 0xdd
	MessageType             = 0xde
	RefMessageType          = 0xdf
	DirectRefMessageMaxType = 0xff
)

// direct value limit
const (
	Int32Zero                int32 = 0x50
	Int64Zero                int64 = 0x88
	DirectStringMaxLength          = DirectStringMaxType
	DirectInt32MinValue            = DirectInt32MinType - Int32Zero
	DirectInt32MaxValue            = DirectInt32MaxType - Int32Zero
	DirectInt64MinValue            = DirectInt64MinType - Int64Zero
	DirectInt64MaxValue            = DirectInt64MaxType - Int64Zero
	DirectRefMessageMaxValue       = DirectRefMessageMaxType - RefMessageType
)

// configure
var (
	MaxWriteCount = 0 // TODO check circular reference
	MaxElemSize   = 100000
)

// common errors
var (
	ErrNoSchema  = errors.New("breeze: not have breeze schema")
	ErrWrongSize = errors.New("breeze: read byte size not correct")
	ErrNotEnough = errors.New("breeze: not enough bytes")
	ErrOverflow  = errors.New("breeze: integer overflow")
)

// Message is a interface of breeze message. all breeze message must implement Message
type Message interface {
	WriteTo(buf *Buffer) error
	ReadFrom(buf *Buffer) error
	GetName() string
	GetAlias() string
	GetSchema() *Schema
}

// Enum is a special breeze message, it can read a `new` enum value from breeze buffer. Enum should be declared as pointer not value, thus the default value of Enum will be null in serialization.
type Enum interface {
	Message
	ReadEnum(buf *Buffer, asAddr bool) (Enum, error)
}

// Context is a context of breeze encode/decode in BreezeBuffer
type Context struct {
	messageTypeRefCount int
	messageTypeRefName  map[int]string
	messageTypeRefIndex map[string]int
}

func (c *Context) getMessageTypeName(index int) (name string) {
	if c.messageTypeRefName != nil {
		name = c.messageTypeRefName[index]
	}
	return name
}

func (c *Context) getMessageTypeIndex(name string) int {
	if c.messageTypeRefIndex != nil {
		index, ok := c.messageTypeRefIndex[name]
		if ok {
			return index
		}
	}
	return -1
}

func (c *Context) putMessageType(name string) {
	if c.messageTypeRefName == nil {
		c.messageTypeRefName = make(map[int]string, 16)
		c.messageTypeRefIndex = make(map[string]int, 16)
	}
	c.messageTypeRefCount++
	c.messageTypeRefName[c.messageTypeRefCount] = name
	c.messageTypeRefIndex[name] = c.messageTypeRefCount
}

// GenericMessage is a generic breeze message. it can receive any breeze message
type GenericMessage struct {
	Name   string
	Alias  string
	schema *Schema
	fields map[int]interface{}
}

// GetAlias return breeze message alias for multi language compatible
func (g *GenericMessage) GetAlias() string {
	return g.Alias
}

// WriteTo write breeze message to breeze buffer.
func (g *GenericMessage) WriteTo(buf *Buffer) error {
	return WriteMessageWithoutType(buf, func(buf *Buffer) {
		for k, v := range g.fields {
			WriteField(buf, k, v)
		}
	})
}

// ReadFrom read a breeze message from breeze buffer
func (g *GenericMessage) ReadFrom(buf *Buffer) error {
	return ReadMessageField(buf, func(buf *Buffer, index int) (err error) {
		v, err := ReadValue(buf, nil)
		if err != nil {
			return err
		}
		if g.fields == nil {
			g.fields = make(map[int]interface{}, DefaultSize)
		}
		g.fields[index] = v
		return nil
	})
}

// GetName get the name of breeze message
func (g *GenericMessage) GetName() string {
	return g.Name
}

// GetSchema get breeze message's schema
func (g *GenericMessage) GetSchema() *Schema {
	return g.schema
}

// GetFieldByIndex get a GenericMessage's field by field index
func (g *GenericMessage) GetFieldByIndex(index int) interface{} {
	if g.fields == nil {
		return nil
	}
	return g.fields[index]
}

// GetFieldByName get a GenericMessage's field by field name
func (g *GenericMessage) GetFieldByName(name string) (interface{}, error) {
	if g.schema == nil {
		return nil, ErrNoSchema
	}
	if g.fields == nil {
		return nil, nil
	}
	field := g.schema.GetFieldByName(name)
	if field != nil {
		return g.fields[field.Index], nil
	}
	return nil, nil
}

// PutField put a field into a GenericMessage
func (g *GenericMessage) PutField(index int, field interface{}) {
	if index > -1 && field != nil {
		if g.fields == nil {
			g.fields = make(map[int]interface{}, 16)
		}
		g.fields[index] = field
	}
}

// Schema describes a breeze message, include name, alias, all fields of message
type Schema struct {
	Name          string
	Alias         string
	indexFieldMap map[int]*Field
	nameFieldMap  map[string]*Field
}

// PutFields put a field into a schema
func (s *Schema) PutFields(fields ...*Field) {
	if s.indexFieldMap == nil {
		s.indexFieldMap = make(map[int]*Field, DefaultSize)
	}
	if s.nameFieldMap == nil {
		s.nameFieldMap = make(map[string]*Field, DefaultSize)
	}
	for _, value := range fields {
		if value != nil && value.Index > -1 {
			s.indexFieldMap[value.Index] = value
			s.nameFieldMap[value.Name] = value
		}
	}
}

// GetFieldByIndex get a message's field from schema by field index
func (s *Schema) GetFieldByIndex(index int) *Field {
	if s.indexFieldMap != nil {
		return s.indexFieldMap[index]
	}
	return nil
}

// GetFieldByName get a message's field from schema by field name
func (s *Schema) GetFieldByName(name string) *Field {
	if s.nameFieldMap != nil {
		return s.nameFieldMap[name]
	}
	return nil
}

// Field describes a message field, include field index, field name and field type
type Field struct {
	Index int
	Name  string
	Type  string
}

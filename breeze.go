package breeze

import (
	"errors"
)

// breeze type
const (
	NULL = iota
	TRUE
	FALSE
	STRING
	BYTE
	BYTES
	INT16
	INT32
	INT64
	FLOAT32
	FLOAT64

	MAP     = 20
	ARRAY   = 21
	MESSAGE = 22
	SCHEMA  = 23
)

// configure
var (
	MaxWriteCount = 0 // TODO check circular reference
)

// common errors
var (
	ErrNoSchema  = errors.New("Breeze: not have breeze schema")
	ErrWrongSize = errors.New("Breeze: read byte size not correct")
	ErrNotEnough = errors.New("Breeze: not enough bytes")
	ErrOverflow  = errors.New("Breeze: integer overflow")
)

type Serializer interface {
	WriteTo(buf *Buffer, v interface{}) (bool, error)
	ReadFrom(buf *Buffer, v interface{}) error
}

type Message interface {
	WriteTo(buf *Buffer) error
	ReadFrom(buf *Buffer) error
	GetName() string
	GetAlias() string
	GetSchema() *Schema
}

type GenericMessage struct {
	Name   string
	Alias  string
	schema *Schema
	fields map[int]interface{}
}

func (g *GenericMessage) GetAlias() string {
	return g.Alias
}

func (g *GenericMessage) WriteTo(buf *Buffer) error {
	return WriteMessage(buf, g.Name, func(buf *Buffer) {
		for k, v := range g.fields {
			WriteMessageField(buf, k, v)
		}
	})
}

func (g *GenericMessage) ReadFrom(buf *Buffer) error {
	return ReadMessageByField(buf, func(buf *Buffer, index int) (err error) {
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

func (g *GenericMessage) GetName() string {
	return g.Name
}

func (g *GenericMessage) GetSchema() *Schema {
	return g.schema
}

func (g *GenericMessage) GetFieldByIndex(index int) interface{} {
	if g.fields == nil {
		return nil
	}
	return g.fields[index]
}

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

func (g *GenericMessage) PutField(index int, field interface{}) {
	if index > -1 && field != nil {
		if g.fields == nil {
			g.fields = make(map[int]interface{}, 16)
		}
		g.fields[index] = field
	}
}

type Schema struct {
	Name          string
	Alias         string
	indexFieldMap map[int]*Field
	nameFieldMap  map[string]*Field
}

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

func (s *Schema) GetFieldByIndex(index int) *Field {
	if s.indexFieldMap != nil {
		return s.indexFieldMap[index]
	}
	return nil
}

func (s *Schema) GetFieldByName(name string) *Field {
	if s.nameFieldMap != nil {
		return s.nameFieldMap[name]
	}
	return nil
}

type Field struct {
	Index int
	Name  string
	Type  string
}

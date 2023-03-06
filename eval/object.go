package eval

import (
	"fmt"
	"strconv"
)

type ObjectType = string

const (
	IntegerType = "INTEGER"
	BooleanType = "BOOLEAN"
	NullType    = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type IntegerObject struct {
	Value int64
}

func (integer *IntegerObject) Type() ObjectType {
	return IntegerType
}

func (integer *IntegerObject) Inspect() string {
	return strconv.FormatInt(integer.Value, 10)
}

type BooleanObject struct {
	Value bool
}

func (boolean *BooleanObject) Type() ObjectType {
	return BooleanType
}

func (boolean *BooleanObject) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}

type NullObject struct {
}

func (null *NullObject) Type() ObjectType {
	return NullType
}

func (null *NullObject) Inspect() string {
	return "null"
}

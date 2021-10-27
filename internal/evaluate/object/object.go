package object

import (
	"container/list"
	"fmt"
	"strings"
)

// Type is type of object which represented with string.
type Type string

// Object interface.
type Object interface {
	Type() Type
	Inspect() string
}

// Word is type for string expressions.
type Word struct {
	Value string
}

// Inspect returns string representation of object.
func (s *Word) Inspect() string {
	return s.Value
}

// Type returns type of object.
func (s *Word) Type() Type {
	return WordObj
}

type DocIDs struct {
	Value *list.List
}

// Inspect returns string representation of object.
func (s *DocIDs) Inspect() string {
	var buffer strings.Builder

	buffer.WriteString("[")

	for e := s.Value.Front(); e != nil; e = e.Next() {
		buffer.WriteString(fmt.Sprintf("%d ", e.Value))
	}

	buffer.WriteString("]")

	return buffer.String()
}

// Type returns type of object.
func (s *DocIDs) Type() Type {
	return DocIDsObj
}

// Error is type for errors handling.
type Error struct {
	Message string
}

// Inspect returns string representation of object.
func (e *Error) Inspect() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

// Type returns type of object.
func (e *Error) Type() Type {
	return ErrorObj
}

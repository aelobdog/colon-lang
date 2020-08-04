package colobj

import "fmt"

// At runtime, every node from the ast will be wrapped into an equivalent "object" for the evaluator

// ObjectType : type of "object" that a particular node in the ast produces at runtime
type ObjectType string

const (
	INTEGER  = "INTEGER"
	BOOLEAN  = "BOOLEAN"
	FLOATING = "FLOATING"
	STRING   = "STRING"
)

// Object : an interface that wraps values of all types which can be fed into the evaluator
type Object interface {
	ObType() ObjectType
	ObValue() string
}

// ----------------------------------------------------------------------------

// Integer : A wrapper for integer values
type Integer struct {
	Value int64
}

// ObValue : Integer
func (i *Integer) ObValue() string {
	return fmt.Sprintf("%v", i.Value)
}

// ObType : Integer
func (i *Integer) ObType() ObjectType {
	return INTEGER
}

// ----------------------------------------------------------------------------

// Boolean : A wrapper for boolean values
type Boolean struct {
	Value bool
}

// ObValue : Boolean
func (b *Boolean) ObValue() string {
	return fmt.Sprintf("%v", b.Value)
}

// ObType : Boolean
func (b *Boolean) ObType() ObjectType {
	return BOOLEAN
}

// ----------------------------------------------------------------------------

// Floating : A wrapper for floating-point values
type Floating struct {
	Value float64
}

// ObValue : Floating
func (f *Floating) ObValue() string {
	return fmt.Sprintf("%v", f.Value)
}

// ObType : Floating
func (f *Floating) ObType() ObjectType {
	return FLOATING
}

// ----------------------------------------------------------------------------

// String : A wrapper for string values
type String struct {
	Value string
}

// ObValue : String
func (s *String) ObValue() string {
	return fmt.Sprintf("%v", s.Value)
}

// ObType : String
func (s *String) ObType() ObjectType {
	return STRING
}

// ----------------------------------------------------------------------------

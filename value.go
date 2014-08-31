//  value.go -- the parent of all Goaldi types

package goaldi

import (
	"fmt"
)

type Value interface {
	String() string // return (Go) string image for printing
	Deref() Value   // return dereferenced value

	AsString() *VString // convert self to VString
	AsNumber() *VNumber // convert self to VNumber

	Add(v2 Value) (Value, *Closure)  // e1 + e2
	Mult(v2 Value) (Value, *Closure) // e1 * e2
}

// verify that the interface is implemented as expected
var _ Value = NewNil()
var _ Value = NewNumber(2.7183)
var _ Value = NewString("cowabunga")
var _ Value = NewException("boom")

//  V(x) builds a value of appropriate Goaldi type
func V(x interface{}) Value {
	switch v := x.(type) {
	case nil:
		return NewNil()
	case int:
		return NewNumber(float64(v))
	case float64:
		return NewNumber(v)
	case string:
		return NewString(v)
	case []byte:
		return NewString(string(v))
	default:
		panic(fmt.Sprintf("V(%T:%v)", x, x))
	}
}

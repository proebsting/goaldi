//  value.go

package goaldi

import (
	"fmt"
)

//  V(x) builds a value of appropriate Goaldi type
func V(x interface{}) Value {
	switch v := x.(type) {
	case nil:
		return NilValue
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

//  NewStatic() creates a new static variable, initialized to nil.
func NewStatic() *Value {
	v := new(Value)
	*v = NilValue
	return v
}

//  value.go

package goaldi

import (
	"fmt"
)

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

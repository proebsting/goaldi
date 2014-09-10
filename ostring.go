//  ostring.go -- string operations

package goaldi

import "fmt"

//  extract string value from arbitrary Value, or panic
func sval(v Value) string {
	if n, ok := v.(Stringable); ok {
		return string(*(n.ToString()))
	} else {
		panic("not a string: " + fmt.Sprintf("%v", v))
	}
}

//  e1 || e2

func (v1 *VNumber) Concat(v2 Value) Value {
	return v1.ToString().Concat(v2)
}

func (v1 *VString) Concat(v2 Value) Value {
	return NewString(string(*v1) + sval(v2))
}

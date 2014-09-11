//  ostring.go -- string operations

package goaldi

//  extract string value from arbitrary Value, or panic
func sval(v Value) string {
	if n, ok := v.(Stringable); ok {
		return string(*(n.ToString()))
	} else {
		panic(&RunErr{"Not a string", v})
	}
}

//  Concat:  e1 || e2

type IConcat interface { // e1 || e2
	Concat(Value) Value
}

var _ IConcat = NewNumber(1)
var _ IConcat = NewString("a")

func (v1 *VNumber) Concat(v2 Value) Value {
	return v1.ToString().Concat(v2)
}

func (v1 *VString) Concat(v2 Value) Value {
	return NewString(string(*v1) + sval(v2))
}

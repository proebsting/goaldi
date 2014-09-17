//  ostring.go -- string operations

package goaldi

//  extract VString value from arbitrary Value, or panic
func sval(v Value) *VString {
	if n, ok := v.(Stringable); ok {
		return n.ToString()
	} else {
		panic(&RunErr{"Not a string", v})
	}
}

//------------------------------------  Concat:  e1 || e2

type IConcat interface {
	Concat(Value) Value
}

func (v1 *VNumber) Concat(v2 Value) Value {
	return v1.ToString().Concat(v2)
}

func (v1 *VString) Concat(v2 Value) Value {
	return NewString(v1.data + sval(v2).data)
}

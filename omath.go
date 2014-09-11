//  omath.go -- math operations

package goaldi

//  extract float64 value from arbitrary Value, or panic
func fval(v Value) float64 {
	if n, ok := v.(Numerable); ok {
		return float64(*(n.ToNumber()))
	} else {
		panic(&RunErr{"Not a number", v})
	}
}

//  Negate:  -e

type INegate interface { // -e
	Negate() Value
}

var _ INegate = NewNumber(1)
var _ INegate = NewString("1")

func (v1 *VString) Negate() Value {
	return v1.ToNumber().Negate()
}

func (v1 *VNumber) Negate() Value {
	return NewNumber(-float64(*v1))
}

//  Add:  e1 + e2

type IAdd interface { // e1 + e2
	Add(Value) Value
}

var _ IAdd = NewNumber(1)
var _ IAdd = NewString("1")

func (v1 *VString) Add(v2 Value) Value {
	return v1.ToNumber().Add(v2)
}

func (v1 *VNumber) Add(v2 Value) Value {
	return NewNumber(float64(*v1) + fval(v2))
}

//  Mult:  e1 * e2

type IMult interface { // e1 * e2
	Mult(Value) Value
}

var _ IMult = NewNumber(1)
var _ IMult = NewString("1")

func (v1 *VString) Mult(v2 Value) Value {
	return v1.ToNumber().Mult(v2)
}

func (v1 *VNumber) Mult(v2 Value) Value {
	return NewNumber(float64(*v1) * fval(v2))
}

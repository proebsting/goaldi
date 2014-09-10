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

//  -e

func (v1 *VString) Negate() Value {
	return v1.ToNumber().Negate()
}

func (v1 *VNumber) Negate() Value {
	return NewNumber(-float64(*v1))
}

//  e1 + e2

func (v1 *VString) Add(v2 Value) Value {
	return v1.ToNumber().Add(v2)
}

func (v1 *VNumber) Add(v2 Value) Value {
	return NewNumber(float64(*v1) + fval(v2))
}

//  e1 * e2

func (v1 *VString) Mult(v2 Value) Value {
	return v1.ToNumber().Mult(v2)
}

func (v1 *VNumber) Mult(v2 Value) Value {
	return NewNumber(float64(*v1) * fval(v2))
}

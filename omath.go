//  omath.go -- math operations

//#%#%  not quite right yet.  not guaranteed that e2 is a VNumber.

package goaldi

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
	n1 := v1
	n2 := v2.(Numerable).ToNumber()
	return NewNumber(float64(*n1) + float64(*n2))
}

//  e1 * e2

func (v1 *VString) Mult(v2 Value) Value {
	return v1.ToNumber().Mult(v2)
}

func (v1 *VNumber) Mult(v2 Value) Value {
	n1 := v1
	n2 := v2.(Numerable).ToNumber()
	return NewNumber(float64(*n1) * float64(*n2))
}

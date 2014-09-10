//  omath.go -- math operations

//#%#%  not quite right yet.  not guaranteed that e2 is a VNumber.

package goaldi

//  -e

func (v1 *VString) Negate() IMath {
	return v1.ToNumber().Negate()
}

func (v1 *VNumber) Negate() IMath {
	return NewNumber(-float64(*v1))
}

//  e1 + e2

func (v1 *VString) Add(v2 IMath) IMath {
	return v1.ToNumber().Add(v2)
}

func (v1 *VNumber) Add(v2 IMath) IMath {
	return NewNumber(float64(*v1) + float64(*(v2.(*VNumber))))
}

//  e1 * e2

func (v1 *VString) Mult(v2 IMath) IMath {
	return v1.ToNumber().Mult(v2)
}

func (v1 *VNumber) Mult(v2 IMath) IMath {
	return NewNumber(float64(*v1) * float64(*(v2.(*VNumber))))
}

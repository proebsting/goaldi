//  math.go -- math operations

package goaldi

//  e1 + e2

func (v1 *VString) Add(v2 Value) (Value, *Closure) {
	return v1.AsNumber().Add(v2)
}

func (v1 *VNumber) Add(v2 Value) (Value, *Closure) {
	return v2.AsNumber().AddInto(v1)
}

func (v2 *VNumber) AddInto(v1 *VNumber) (Value, *Closure) {
	return Return(NewNumber(float64(*v1) + float64(*v2)))
}

//  e1 * e2

func (v1 *VString) Mult(v2 Value) (Value, *Closure) {
	return v1.AsNumber().Mult(v2)
}

func (v1 *VNumber) Mult(v2 Value) (Value, *Closure) {
	return v2.AsNumber().MultInto(v1)
}

func (v2 *VNumber) MultInto(v1 *VNumber) (Value, *Closure) {
	return Return(NewNumber(float64(*v1) * float64(*v2)))
}

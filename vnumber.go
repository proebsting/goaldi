//  vnumber.go -- VNumber, the Goaldi type "number"

package goaldi

import (
	"fmt"
)

type VNumber float64

//  NewNumber -- construct a Goaldi number
func NewNumber(n float64) *VNumber {
	vn := VNumber(n)
	return &vn
}

// predefined constants
var (
	ZERO      = NewNumber(0)
	ONE       = NewNumber(1)
	MAX_EXACT = 1 << 53 // beyond 9e15, integers are noncontiguous
)

//  VNumber.Val -- return underlying float64 value
func (v *VNumber) Val() float64 {
	return float64(*v)
}

//  VNumber.String -- convert to Go string
func (v *VNumber) String() string {
	return fmt.Sprintf("%g", float64(*v))
}

//  VNumber.ToString -- convert to Goaldi string
func (v *VNumber) ToString() *VString {
	return NewString(v.String())
}

//  VNumber.Number -- return self
func (v *VNumber) ToNumber() *VNumber {
	return v
}

//  VNumber.Type -- return "number"
func (v *VNumber) Type() Value {
	return type_number
}

var type_number = NewString("number")

//  VNumber.Identical -- check equality for === operator
func (a *VNumber) Identical(x Value) Value {
	b, ok := x.(*VNumber)
	if ok && a.Val() == b.Val() {
		return x
	} else {
		return nil
	}
}

//  VNumber.Export returns a float64
func (v *VNumber) Export() interface{} {
	return float64(*v)
}

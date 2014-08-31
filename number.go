//  number.go -- the Goaldi type "number"

package goaldi

import (
	"fmt"
)

type VNumber float64

func NewNumber(n float64) *VNumber {
	vn := VNumber(n)
	return &vn
}

func (v *VNumber) String() string {
	return fmt.Sprintf("%g", float64(*v))
}

func (v *VNumber) Deref() Value {
	return v
}

func (v *VNumber) AsString() *VString {
	return NewString(v.String())
}

func (v *VNumber) AsNumber() *VNumber {
	return v
}

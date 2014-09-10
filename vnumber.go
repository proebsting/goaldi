//  vnumber.go -- VNumber, the Goaldi type "number"

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

func (v *VNumber) ToString() *VString {
	return NewString(v.String())
}

func (v *VNumber) ToNumber() *VNumber {
	return v
}

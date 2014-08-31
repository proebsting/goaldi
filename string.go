//  string.go -- the Goaldi type "string"

package goaldi

import (
	"fmt"
)

type VString string

func NewString(s string) *VString {
	vs := VString(s)
	return &vs
}

func (v *VString) String() string {
	return string(*v)
}
func (v *VString) Deref() Value {
	return v
}

func (v *VString) AsString() *VString {
	return v
}

func (v *VString) AsNumber() *VNumber {
	var f float64
	var b byte
	n, _ := fmt.Sscanf(string(*v), "%f%c", &f, &b)
	if n == 1 {
		return NewNumber(f)
	} else {
		return nil
	}
}

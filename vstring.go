//  vstring.go -- VString, the Goaldi type "string"

package goaldi

import (
	"fmt"
)

type VString string

//  NewString -- construct a Goaldi string
func NewString(s string) *VString {
	vs := VString(s)
	return &vs
}

//  VString.String -- return image of string, quoted, for debugging
func (v *VString) String() string {
	return `"` + string(*v) + `"`
}

//  VString.ToString -- for a Goaldi string, this just returns self
func (v *VString) ToString() *VString {
	return v
}

//  VString.Number -- return conversion to VNumber, or nil for failure
func (v *VString) ToNumber() *VNumber {
	var f float64
	var b byte
	n, _ := fmt.Sscanf(string(*v), "%f%c", &f, &b)
	if n == 1 {
		return NewNumber(f)
	} else {
		return nil
	}
}

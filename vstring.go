//  vstring.go -- VString, the Goaldi type "string"
//
//  Strings encode sequences of Unicode characters (Code Points or Runes)

package goaldi

import (
	"fmt"
)

type VString struct {
	data  string
	hints interface{} // TBD
}

//  NewString -- construct a Goaldi string from a Go UTF8 string
func NewString(s string) *VString {
	vs := &VString{s, nil}
	return vs
}

//  VString.ToUTF8 -- convert Goaldi Unicode string to Go UTF8 string
func (v *VString) ToUTF8() string {
	return v.data
}

//  VString.String -- return image of string, quoted, as a Go string
func (v *VString) String() string {
	return `"` + v.ToUTF8() + `"`
}

//  VString.ToString -- for a Goaldi string, this just returns self
func (v *VString) ToString() *VString {
	return v
}

//  VString.Number -- return conversion to VNumber, or nil for failure
func (v *VString) ToNumber() *VNumber {
	var f float64
	var b byte
	n, _ := fmt.Sscanf(v.data, "%f%c", &f, &b)
	if n == 1 {
		return NewNumber(f)
	} else {
		return nil
	}
}

//  VString.Type -- return "string"
func (v *VString) Type() Value {
	return type_string
}

var type_string = NewString("string")

//  VString.Identical -- check equality for === operator
func (s *VString) Identical(x Value) Value {
	t, ok := x.(*VString)
	if ok && s.data == t.data {
		return x
	} else {
		return nil
	}
}

//  VString.Export returns a Go string
func (v *VString) Export() interface{} {
	return v.ToUTF8()
}

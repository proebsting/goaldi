//  vstring.go -- VString, the Goaldi type "string"
//
//  Strings contain sequences of Unicode characters (Code Points or Runes)
//

package goaldi

import (
	"fmt"
	"unicode/utf8"
)

//  A string is encoded by one (usually) or two parallel slices
type VString struct {
	low  []uint8  // required: low-order 8 bits of each rune
	high []uint16 // optional: high-order 13 bits of each rune
}

//  NewString -- construct a Goaldi string from a Go UTF8 string
func NewString(s string) *VString {
	n := len(s)
	low := make([]uint8, n, n)
	high := make([]uint16, n, n)
	h := '\000'
	i := 0
	for _, c := range s {
		h |= c
		low[i] = uint8(c)
		high[i] = uint16(c >> 8)
		i++
	}
	// #%#% could copy now to smaller underlying arrays if warranted
	if (h >> 8) == 0 {
		return &VString{low[:i], nil}
	} else {
		return &VString{low[:i], high[:i]}
	}
}

//  EasyString -- construct a Goaldi string from ASCII input (no byte > 0x7F)
func EasyString(s string) *VString {
	return &VString{[]uint8(s), nil}
}

//  BinaryString -- construct a Goaldi string from Go Latin1 bytes
func BinaryString(s []byte) *VString {
	b := make([]uint8, len(s), len(s))
	copy(b, s)
	return &VString{b, nil}
}

//  VString.ToUTF8 -- convert Goaldi Unicode string to Go UTF8 string
func (v *VString) ToUTF8() string {
	b := make([]byte, 0, len(v.low))
	p := make([]byte, 8, 8)
	for i, c := range v.low {
		r := rune(c)
		if v.high != nil {
			r |= rune(v.high[i]) << 8
		}
		n := utf8.EncodeRune(p, r)
		b = append(b, p[:n]...)
	}
	return string(b)
}

//  VString.ToBinary -- convert Goaldi Unicode to 8-bit bytes by truncation
func (v *VString) ToBinary() []byte {
	return []byte(v.low)
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
	if v.high != nil { // if has exotic characters //#%#% bogus test?
		return nil // it can't be valid
	}
	n, _ := fmt.Sscanf(string(v.low), "%f%c", &f, &b)
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

var type_string = EasyString("string")

//  VString.Identical -- check equality for === operator
func (s *VString) Identical(x Value) Value {
	t, ok := x.(*VString)
	if !ok {
		return nil
	} else if s == t {
		return t
	} else {
		return s.LEqual(t)
	}
}

//  VString.Export returns a Go string
func (v *VString) Export() interface{} {
	return v.ToUTF8()
}

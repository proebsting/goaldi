//  vnumber.go -- VNumber, the Goaldi type "number"

package goaldi

import (
	"fmt"
	"math"
	"strconv"
)

type VNumber float64

//  NewNumber -- construct a Goaldi number from a float value
func NewNumber(n float64) *VNumber {
	vn := VNumber(n)
	return &vn
}

//  ParseNumber -- standard string-to-number conversion for Goaldi
//  Currently allows only Go standard format, plus leading and trailing spaces.
func ParseNumber(s string) (float64, error) {
	var f float64
	var b byte
	n, _ := fmt.Sscanf(s, "%f %c", &f, &b)
	if n == 1 {
		return f, nil
	} else {
		return math.NaN(), &RunErr{"Not a number", s}
	}

}

//  MustParseNum -- make a float from a string, or throw a RunErr
func MustParseNum(s string) float64 {
	f, e := ParseNumber(s)
	if e != nil {
		panic(e)
	} else {
		return f
	}
}

// predefined constants
var (
	ZERO = NewNumber(0)
	ONE  = NewNumber(1)
	E    = NewNumber(math.E)
	PI   = NewNumber(math.Pi)
	PHI  = NewNumber(math.Phi)
)

//  VNumber.Val -- return underlying float64 value
func (v *VNumber) Val() float64 {
	return float64(*v)
}

//  VNumber.String -- default conversion to Go string
func (v *VNumber) String() string {
	i := int64(*v)
	if v.IsExactInt(i) {
		return strconv.FormatInt(i, 10) // if exact integer
	} else {
		return fmt.Sprintf("%.4g", float64(*v)) // if has fractional bits
	}
}

//  VNumber.GoString -- convert to Go string for image() and printf("%#v")
//  The difference vs String() is that all significant digits are returned
func (v *VNumber) GoString() string {
	i := int64(*v)
	if v.IsExactInt(i) {
		return strconv.FormatInt(i, 10) // if exact integer
	} else {
		return fmt.Sprintf("%g", float64(*v)) // if has fractional bits
	}
}

//  VNumber.ToString -- convert to Goaldi string
func (v *VNumber) ToString() *VString {
	return NewString(v.String())
}

//  VNumber.Number -- return self
func (v *VNumber) ToNumber() *VNumber {
	return v
}

//  VNumber.Rank returns rNumber
func (v *VNumber) Rank() int {
	return rNumber
}

//  VNumber.Type -- return "number"
func (v *VNumber) Type() Value {
	return type_number
}

var type_number = NewString("number")

//  VNumber.Copy returns itself
func (v *VNumber) Copy() Value {
	return v
}

//  VNumber.Identical -- check equality for === operator
func (a *VNumber) Identical(x Value) Value {
	b, ok := x.(*VNumber)
	if ok && a.Val() == b.Val() {
		return x
	} else {
		return nil
	}
}

//  VNumber.Import returns itself
func (v *VNumber) Import() Value {
	return v
}

//  VNumber.Export returns a float64
func (v *VNumber) Export() interface{} {
	return float64(*v)
}

//  VNumber.IsExactInt returns true if this VNumber represents int i exactly
func (v *VNumber) IsExactInt(i int64) bool {
	return float64(i) == float64(*v) && i <= MAX_EXACT && i >= -MAX_EXACT
}

const MAX_EXACT = 1 << 53 // beyond 9e15, integers are noncontiguous

//  vnumber.go -- VNumber, the Goaldi type "number"

package goaldi

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type VNumber float64

//  NewNumber -- construct a Goaldi number from a float value
func NewNumber(n float64) *VNumber {
	vn := VNumber(n)
	return &vn
}

const rNumber = 10              // declare sort ranking
var _ ICore = NewNumber(1)      // validate implementation
var _ Numerable = NewNumber(1)  // validate implementation
var _ Stringable = NewNumber(1) // validate implementation

//  NumberType is the number instance of type type.
var NumberType = NewType("number", "n", rNumber, Number, nil,
	"number", "x", "convert to number")

//  ParseNumber -- standard string-to-number conversion for Goaldi.
//  Trims leading spaces and tabs, then allows either stadard Go
//  "ParseFloat" form or any Goaldi radix form (nnb, nno, nnx, nnrxxxx).
func ParseNumber(s string) (float64, error) {
	// trim leading and trailing strings; must have something left
	s = strings.Trim(s, " \t")
	if len(s) == 0 {
		return 0, mtyerr
	}
	// try first to interpret as a decimal number (fixed or floating)
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f, nil
	}
	// check next for old Icon nnRxxx form
	parts := nnrxx.FindStringSubmatch(s)
	if parts != nil {
		radix, _ := strconv.Atoi(parts[1])
		value, err := strconv.ParseInt(parts[2], radix, 64)
		return float64(value), err
	}
	// the only other possibility is radix suffix form:  nnnb, nnno, nnnx
	radix := 0
	switch s[len(s)-1] {
	case 'b':
		radix = 2
	case 'o':
		radix = 8
	case 'x':
		radix = 16
	default:
		return 0, numerr
	}
	value, err := strconv.ParseInt(s[0:len(s)-1], radix, 64)
	return float64(value), err
}

var nnrxx = regexp.MustCompile("^([0-9]+)[rR]([0-9a-zA-Z]+)$")
var mtyerr = errors.New("empty string for numeric conversion")
var numerr = errors.New("malformed number")

// predefined constants
var (
	ZERO = NewNumber(0)
	ONE  = NewNumber(1)
	INF  = NewNumber(math.Inf(+1))
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

//  VNumber.Type -- return the number type
func (v *VNumber) Type() IRank {
	return NumberType
}

//  VNumber.Copy returns itself
func (v *VNumber) Copy() Value {
	return v
}

//  VNumber.Before compares two numbers for sorting
func (a *VNumber) Before(b Value, i int) bool {
	return *a < *(b.(*VNumber))
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

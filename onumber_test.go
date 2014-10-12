//  onumber_test.go -- test numeric conversions and onumber.go operations

package goaldi

import (
	"fmt"
	"testing"
)

func TestMath(t *testing.T) {
	i6, s6 := nspair(t, 6)
	i7, s7 := nspair(t, 7)
	i8, s8 := nspair(t, 8)
	i11, s11 := nspair(t, 11)
	i12, s12 := nspair(t, 12)
	i30, s30 := nspair(t, 30)
	f25, s25 := nspair(t, 2.5)
	ck4n(t, "Numerate", 7,
		i7.Numerate(), i7.Numerate(), s7.Numerate(), s7.Numerate())
	ck4n(t, "Negate", -11,
		i11.Negate(), i11.Negate(), s11.Negate(), s11.Negate())
	ck4n(t, "Add", 18,
		i7.Add(i11), i7.Add(s11), s7.Add(i11), s7.Add(s11))
	ck4n(t, "Sub", -4,
		i7.Sub(i11), i7.Sub(s11), s7.Sub(i11), s7.Sub(s11))
	ck4n(t, "Mul", 77,
		i7.Mul(i11), i7.Mul(s11), s7.Mul(i11), s7.Mul(s11))
	ck4n(t, "Div1", 5,
		i30.Div(i6), i30.Div(s6), s30.Div(i6), s30.Div(s6))
	ck4n(t, "Div2", 2.5,
		i30.Div(i12), i30.Div(s12), s30.Div(i12), s30.Div(s12))
	ck4n(t, "Divt", 2,
		i30.Divt(i11), i30.Divt(s11), s30.Divt(i11), s30.Divt(s11))
	ck4n(t, "Mod1", 8,
		i30.Mod(i11), i30.Mod(s11), s30.Mod(i11), s30.Mod(s11))
	ck4n(t, "Mod2", 2,
		i12.Mod(f25), i12.Mod(s25), s12.Mod(f25), s12.Mod(s25))
	ck4n(t, "Mod3", 0.5,
		i8.Mod(f25), i8.Mod(s25), s8.Mod(f25), s8.Mod(s25))
	ck4n(t, "Power", 117649,
		i7.Power(i6), i7.Power(s6), s7.Power(i6), s7.Power(s6))
}

//  nspair -- return number as a pair (number, string), checking conversions
func nspair(t *testing.T, v float64) (*VNumber, *VString) {
	n1 := NewNumber(v)
	s1 := NewString(fmt.Sprintf("%g", v))
	n2 := s1.ToNumber()
	s2 := n2.ToString()
	if n1.Val() != n2.Val() {
		t.Errorf("numbers %v != %v", n1, n2)
	}
	if s1.String() != s2.String() {
		t.Errorf("strings %v != %v", s1, s2)
	}
	return n1, s1
}

//  ck4n -- check four numeric values for equality with expected value
func ck4n(t *testing.T, label string, n0 float64, v1, v2, v3, v4 Value) {
	t.Log("testing", label)
	n1 := v1.(*VNumber).Val()
	n2 := v2.(*VNumber).Val()
	n3 := v3.(*VNumber).Val()
	n4 := v4.(*VNumber).Val()
	if n0 != n1 || n1 != n2 || n2 != n3 || n3 != n4 {
		t.Errorf("Expected %g, got %g %g %g %g", n0, n1, n2, n3, n4)
	}
}

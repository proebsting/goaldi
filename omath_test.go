//  omath_test.go -- test numeric conversions and omath.go operations

package goaldi

import (
	"fmt"
	"testing"
)

func TestMath(t *testing.T) {
	i7, s7 := nspair(t, 7)
	i11, s11 := nspair(t, 11)
	t.Log("values:", i7, s7, i11, s11)
	ck4n(t, "Negate", -7,
		i7.Negate(), i7.Negate(), s7.Negate(), s7.Negate())
	ck4n(t, "Add", 18,
		i7.Add(i11), i7.Add(s11), s7.Add(i11), s7.Add(s11))
	ck4n(t, "Mult", 77,
		i7.Mult(i11), i7.Mult(s11), s7.Mult(i11), s7.Mult(s11))
}

//  nspair -- return number as a pair (number, string), checking conversions
func nspair(t *testing.T, v float64) (*VNumber, *VString) {
	n1 := NewNumber(v)
	s1 := NewString(fmt.Sprintf("%g", v))
	n2 := s1.ToNumber()
	s2 := n2.ToString()
	if n1.val() != n2.val() {
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
	n1 := v1.(*VNumber).val()
	n2 := v2.(*VNumber).val()
	n3 := v3.(*VNumber).val()
	n4 := v4.(*VNumber).val()
	if n0 != n1 || n1 != n2 || n2 != n3 || n3 != n4 {
		t.Errorf("Expected %g: %g %g %g %g", n0, n1, n2, n3, n4)
	}
}

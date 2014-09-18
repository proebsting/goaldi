//  unicode_test.go -- test Unicode string construction and comparison

package goaldi

import (
	"math/rand"
	"testing"
)

func TestUnicode(t *testing.T) {
	testUniRange(t, '0', '9')          // digits
	testUniRange(t, ' ', '_')          // SIXBIT
	testUniRange(t, ' ', '~')          // ASCII printable
	testUniRange(t, ' ', 'ÿ')          // Latin1 printable
	testUniRange(t, ' ', '\u06FF')     // 9-bit mostly printable
	testUniRange(t, ' ', '\u0FFF')     // 12-bit mostly printable
	testUniRange(t, ' ', '\uD78F')     // 16-bit mostly printable
	testUniRange(t, ' ', '\U0002FFFF') // 17-bit most chars in use
	testUniRange(t, ' ', '\U0010FFFF') // full range
}

func testUniRange(t *testing.T, low, high rune) {
	for len := 1; len <= 6; len++ {
		r1 := make([]rune, 0, len)
		r2 := make([]rune, len, len)
		r3 := make([]rune, len, len)
		spanm2 := int(high - low - 2)
		for i := 0; i < len; i++ {
			r1 = append(r1, low+1+rune(rand.Intn(spanm2)))
		}
		copy(r2, r1)
		copy(r3, r1)
		r1[rand.Intn(len)]--
		r3[rand.Intn(len)]++
		u1 := string(r1)
		u2 := string(r2)
		u3 := string(r3)
		testUniString(t, len, u1, u2, u3)
	}
}

func testUniString(t *testing.T, len int, u1, u2, u3 string) {
	s1 := NewString(u1)
	s2 := NewString(u2)
	s3 := NewString(u3)
	t.Log(len, u1, s1, u2, s2, u3, s3)
	uu1 := s1.ToUTF8()
	uu2 := s2.ToUTF8()
	uu3 := s3.ToUTF8()
	if u1 != uu1 || u2 != uu2 || u3 != uu3 {
		t.Error("there-and-back error")
	}
	lenf := float64(len)
	z1 := s1.Size().(*VNumber).val()
	z2 := s2.Size().(*VNumber).val()
	z3 := s3.Size().(*VNumber).val()
	if z1 != lenf || z2 != lenf || z3 != lenf {
		t.Error("expected length", lenf, " but got ", z1, z2, z3)
	}
	c11 := s1.compare(s1)
	c12 := s1.compare(s2)
	c13 := s1.compare(s3)
	c21 := s2.compare(s1)
	c22 := s2.compare(s2)
	c23 := s2.compare(s3)
	c31 := s3.compare(s1)
	c32 := s3.compare(s2)
	c33 := s3.compare(s3)
	if c11 != 0 || c12 >= 0 || c13 >= 0 ||
		c21 <= 0 || c22 != 0 || c23 >= 0 ||
		c31 <= 0 || c32 <= 0 || c33 != 0 {
		t.Error("comparisons", c11, c12, c13, c23, c22, c23, c31, c32, c33)
	}
}

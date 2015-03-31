//  unicode_test.go -- test Unicode string construction and comparison

package runtime

import (
	"math/rand"
	"testing"
)

const maxUniLen = 7

func TestUnicode(t *testing.T) {
	testUniConcat(t, EMPTY, EMPTY, EMPTY)
	testUniConcat(t, NewString("a"), NewString("b"), NewString("c"))
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

func testUniRange(t *testing.T, low rune, high rune) {
	slist := make([]string, 0)
	for len := 1; len <= maxUniLen; len++ {
		slist = append(slist, randomUniString(len, low+1, high-1))
	}
	t.Log("maxrune", int(high), slist)
	for len := 1; len <= maxUniLen; len++ {
		s := slist[len-1]
		r1 := []rune(s)
		r2 := []rune(s)
		r3 := []rune(s)
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
	uu1 := s1.ToUTF8()
	uu2 := s2.ToUTF8()
	uu3 := s3.ToUTF8()
	if u1 != uu1 || u2 != uu2 || u3 != uu3 {
		t.Error("there-and-back error", u1, uu1, u2, uu2, u3, uu3)
	}
	lenf := float64(len)
	z1 := s1.Size().(*VNumber).Val()
	z2 := s2.Size().(*VNumber).Val()
	z3 := s3.Size().(*VNumber).Val()
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
	abc := NewString("abc")
	testUniConcat(t, s1, EMPTY, EMPTY)
	testUniConcat(t, EMPTY, s1, EMPTY)
	testUniConcat(t, EMPTY, EMPTY, s1)
	testUniConcat(t, s1, EMPTY, abc)
	testUniConcat(t, abc, s1, EMPTY)
	testUniConcat(t, abc, s1, abc)
}

func testUniConcat(t *testing.T, s1, s2, s3 *VString) {
	n1 := s1.length()
	n2 := s2.length()
	n3 := s3.length()
	u1 := s1.ToUTF8()
	u2 := s2.ToUTF8()
	u3 := s3.ToUTF8()
	u123 := u1 + u2 + u3
	s123 := scat(s1, 0, n1, s2, 0, n2, s3, 0, n3)
	us123 := s123.ToUTF8()
	expect(t, "concatenation", u123, us123)
}

func randomUniString(len int, low rune, high rune) string {
	s := ""
	n := int(high - low + 1)
	for i := 0; i < len; i++ {
		s = s + string(low+rune(rand.Intn(n)))
	}
	return s
}

//  ostring_test.go -- test string conversions and ostring.go operations

package goaldi

import (
	"testing"
)

func TestStringOps(t *testing.T) {
	i123, s123 := nspair(t, 123)
	i456, s456 := nspair(t, 456)
	t.Log("values:", i123, s123, i456, s456)
	ck4s(t, "Concat", "123456", i123.Concat(i456), i123.Concat(s456),
		s123.Concat(i456), s123.Concat(s456))
	sh := NewString("♡") // heart
	sd := NewString("♢") // diamond
	sc := NewString("♣") // club
	ss := NewString("♠") // spade
	hd := sh.Concat(sd)
	cs := sc.Concat(ss)
	hdcs := hd.(*VString).Concat(cs)
	ck4s(t, "Concat", "♡♢♣♠", hdcs, hdcs, hdcs, hdcs)
	sz := hdcs.(ISize).Size().(*VNumber).Val()
	if sz != 4.0 {
		t.Errorf("String %s length %d, expected 4", hdcs, sz)
	}
}

// ck4s -- check four string values for equality with expected value
func ck4s(t *testing.T, label string, s0 string, v1, v2, v3, v4 Value) {
	t.Log("testing", label)
	s1 := v1.(*VString).String()
	s2 := v2.(*VString).String()
	s3 := v3.(*VString).String()
	s4 := v4.(*VString).String()
	if s0 != s1 || s1 != s2 || s2 != s3 || s3 != s4 {
		t.Errorf("Expected %s: %s %s %s %s", s0, s1, s2, s3, s4)
	}
}

//  for nspair() see onumber_test.go

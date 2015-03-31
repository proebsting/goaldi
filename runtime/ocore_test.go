//  ocore_test.go -- test core functions Identical and NotIdentical

package runtime

import (
	"testing"
)

func TestCore(t *testing.T) {
	ab := NewString("ab")
	cd := NewString("cd")
	abcd1 := NewString("abcd")
	abcd2 := ab.Concat(cd)
	expect(t, "1s=", abcd1, Identical(abcd1, abcd1))
	expect(t, "2s=", abcd2, Identical(abcd1, abcd2))
	expect(t, "3s=", abcd1, Identical(abcd2, abcd1))
	expect(t, "4s=", nil, Identical(ab, cd))
	expect(t, "5s~", cd, NotIdentical(ab, cd))
	expect(t, "6s~", nil, NotIdentical(ab, ab))
	expect(t, "7s~", nil, NotIdentical(abcd1, abcd2))

	n2 := NewNumber(2)
	n3 := NewNumber(3)
	n6a := NewNumber(6)
	n6b := n2.Mul(n3)
	expect(t, "1n=", n6a, Identical(n6a, n6a))
	expect(t, "2n=", n6b, Identical(n6a, n6b))
	expect(t, "3n=", n6a, Identical(n6b, n6a))
	expect(t, "4n=", nil, Identical(n6b, n3))
	expect(t, "5n~", n3, NotIdentical(n2, n3))
	expect(t, "6n~", nil, NotIdentical(n3, n3))
	expect(t, "7n~", nil, NotIdentical(n6a, n6b))

	expect(t, "1x=", nil, Identical(ab, n2))
	expect(t, "2x=", nil, Identical(n3, cd))
	expect(t, "3x~", n3, NotIdentical(ab, n3))
	expect(t, "4x~", cd, NotIdentical(n2, cd))

	expect(t, "1z=", NilValue, Identical(NilValue, NilValue))
	expect(t, "2z~", ab, NotIdentical(NilValue, ab))
	expect(t, "3z~", n3, NotIdentical(NilValue, n3))
}

//  expect -- check result against expected value
//  n.b. uses Go comparison not Goaldi (does not look inside String or Number)
func expect(t *testing.T, label string, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("%s: expected %v, found %v\n", label, expected, actual)
	}
}

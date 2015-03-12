//	vset.go -- VSet, the Goaldi type "set"
//
//	Implementation:
//	A Goaldi set is just a type name VSet attached to a Go map[Value]bool.
//	This distinguishes it from an external Go map and allows attaching methods.
//	Goaldi strings and numbers are converted to Go string and float64 values.
//  Much of this is very similar to the implementation of the table type.

package goaldi

import (
	"bytes"
	"fmt"
)

//  VSet implements a Goaldi set.
//  Strings and numbers are converted before use as keys;
//  otherwise, unconverted "identical" values would appear distinct.
//  All map values are "true"; deletions remove non-member keys.
type VSet map[Value]bool

//  NewSet -- construct a new Goaldi set from a Goaldi list.
func NewSet(L *VList) VSet {
	S := VSet(make(map[Value]bool))
	for _, v := range L.data {
		S[GoKey(v)] = true
	}
	return S
}

//  GoKey(v) turns a Goaldi value into something usable as a Go map key.
//  Strings and numbers must be converted becuase otherwise two
//  identical values in different objects would be seen as distinct keys.
func GoKey(v Value) interface{} {
	switch t := v.(type) {
	case *VString:
		return t.ToUTF8()
	case *VNumber:
		return t.Val()
	default:
		return v // use key as is
	}
}

//  SetType is the set instance of type type.
var SetType = NewType("set", "S", rSet, Set, SetMethods,
	"set", "L", "create a new set from list L")

//  VSet.String -- default conversion to Go string returns "S:size"
func (S VSet) String() string {
	return fmt.Sprintf("S:%d", len(S))
}

//  VSet.GoString -- convert to Go string for image() and printf("%#v")
//
//  For utility and reproducibility, we accept the cost of sorting the set.
func (S VSet) GoString() string {
	if len(S) == 0 {
		return "set{}"
	}
	l, _ := S.Sort(ONE) // sort on key values
	var b bytes.Buffer
	fmt.Fprintf(&b, "set{")
	for _, e := range l.(*VList).data {
		fmt.Fprintf(&b, "%v,", e)
	}
	s := b.Bytes()
	s[len(s)-1] = '}'
	return string(s)
}

//  VSet.Type -- return the set type
func (S VSet) Type() IRank {
	return SetType
}

//  VSet.Copy returns a duplicate of itself
func (S VSet) Copy() Value {
	r := NewSet(EMPTYLIST)
	for k := range S {
		r[k] = true
	}
	return r
}

//  VSet.Before compares two sets for sorting
func (a VSet) Before(b Value, i int) bool {
	return false // no ordering defined
}

//  VSet.Import returns itself
func (v VSet) Import() Value {
	return v
}

//  VSet.Export returns itself.
//  Go extensions may wish to use v.Member(), v.Delete(), etc.
//  to ensure proper conversion of keys.
func (v VSet) Export() interface{} {
	return v
}

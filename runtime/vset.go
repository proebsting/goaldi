//	vset.go -- VSet, the Goaldi type "set"
//
//	Implementation:
//	A Goaldi set is just a type name VSet attached to a Go map[Value]bool.
//	This distinguishes it from an external Go map and allows attaching methods.
//	Goaldi strings and numbers are converted to Go string and float64 values.
//  Much of this is very similar to the implementation of the table type.

package runtime

import (
	"bytes"
	"fmt"
)

// VSet implements a Goaldi set.
// Strings and numbers are converted before use as keys;
// otherwise, unconverted "identical" values would appear distinct.
// All map values are "true"; deletions remove non-member keys.
type VSet map[Value]bool

const rSet = 65       // declare sort ranking
var _ ICore = &VSet{} // validate implementation

// NewSet -- construct a new Goaldi set from a Goaldi list.
func NewSet(L *VList) *VSet {
	S := VSet(make(map[Value]bool))
	for _, v := range L.data {
		S[GoKey(v)] = true
	}
	return &S
}

// GoKey(v) turns a Goaldi value into something usable as a Go map key.
// Strings and numbers must be converted because otherwise two
// identical values in different objects would be seen as distinct keys.
// The inverse of GoKey() is Import().
func GoKey(v Value) interface{} {
	switch t := v.(type) {
	case *VString:
		return t.ToUTF8()
	case *VNumber:
		return t.Val()
	default:
		return v
	}
}

// SetType is the set instance of type type.
var SetType = NewType("set", "S", rSet, Set, SetMethods,
	"set", "L", "create a new set from list L")

// SetVal(x) return x as a Set, or throws an exception.
func SetVal(x Value) *VSet {
	if S, ok := x.(*VSet); ok {
		return S
	} else {
		panic(NewExn("Not a set", x))
	}
}

// VSet.String -- default conversion to Go string returns "S:size"
func (S *VSet) String() string {
	return fmt.Sprintf("S:%d", len(*S))
}

// VSet.GoString -- convert to Go string for image() and printf("%#v")
//
// For utility and reproducibility, we accept the cost of sorting the set.
func (S *VSet) GoString() string {
	if len(*S) == 0 {
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

// VSet.Type -- return the set type
func (S *VSet) Type() IRank {
	return SetType
}

// VSet.Copy returns a duplicate of itself
func (S *VSet) Copy() Value {
	r := NewSet(EMPTYLIST)
	for k := range *S {
		(*r)[k] = true
	}
	return r
}

// VSet.Before compares two sets for sorting
func (a *VSet) Before(b Value, i int) bool {
	return false // no ordering defined
}

// VSet.Import returns itself
func (S *VSet) Import() Value {
	return S
}

// VSet.Export returns its underlying map[Value]bool.
// Go extensions may wish to use GoKey() for proper conversion of keys.
func (S *VSet) Export() interface{} {
	return map[Value]bool(*S)
}

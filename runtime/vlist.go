//	vlist.go -- VList, the Goaldi type "list"
//
//	Implementation:
//	A list is a slice of values accompanied by a "reversal" flag.
//	Lookups are simple.  Stack and queue operations append or trim the slice.
//	Switching between put and push reverses the list in place for appending.
//	All other modifications, and lookups, are O(1) in amortized cost.

package runtime

import (
	"bytes"
	"fmt"
)

//  A Goaldi list is a Go slice of values.
//  The values are reversed if the last addition was via "push".
type VList struct {
	data []Value // current list contents
	rev  bool    // true if list is reversed
}

const rList = 60       // declare sort ranking
var _ ICore = &VList{} // validate implementation

//  ListType is the list instance of type type.
var ListType = NewType("list", "L", rList, List, ListMethods,
	"list", "size,x", "create list of copies of x")

//  NewList(n, x) -- make a new list of n elements each initialized to copy(x)
func NewList(n int, x Value) *VList {
	v := &VList{make([]Value, n), false}
	for i := range v.data {
		v.data[i], _ = Copy(nil, x)
	}
	return v
}

//  InitList(v []Value) -- construct a new VList containing the given list
//  (directly, without copying)
func InitList(v []Value) *VList {
	return &VList{v, false}
}

//  VList.Elem(i) -- return the ith element, allowing for a reversed list
func (v *VList) Elem(i int) Value {
	n := len(v.data)
	if i < 0 || i >= n {
		return nil
	} else if v.rev {
		return v.data[n-i-1]
	} else {
		return v.data[i]
	}
}

//  VList.String -- default conversion to Go string produces L:size
func (v *VList) String() string {
	return fmt.Sprintf("L:%d", len(v.data))
}

//  VList.GoString -- convert to Go string for image() and printf("%#v")
func (v *VList) GoString() string {
	if len(v.data) == 0 {
		return "[]"
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "[")
	for i := 0; i < len(v.data); i++ {
		fmt.Fprintf(&b, "%v,", v.Elem(i))
	}
	s := b.Bytes()
	s[len(s)-1] = ']'
	return string(s)
}

//  VList.Type -- return the list type
func (v *VList) Type() IRank {
	return ListType
}

//  VList.Copy returns a new list with identical contents
func (v *VList) Copy() Value {
	return InitList(v.Export().([]Value))
}

//  VList.Before compares two lists for sorting on field i
func (a *VList) Before(x Value, i int) bool {
	b := x.(*VList)
	if i >= 0 && len(a.data) > i && len(b.data) > i {
		aref := &vListRef{a, i}
		bref := &vListRef{b, i}
		return LT(aref.Deref(), bref.Deref(), -1)
	} else {
		// put missing one first; otherwise #%#% we don't care
		return len(a.data) < len(b.data)
	}
}

//  VList.Import returns itself
func (v *VList) Import() Value {
	return v
}

//  VList.Export returns a copy of the data slice.
func (v *VList) Export() interface{} {
	n := len(v.data)
	r := make([]Value, n)
	copy(r, v.data)
	if v.rev {
		ReverseValues(r)
	}
	return r
}

//  -------------------------- trapped references ---------------------

type vListRef struct {
	list *VList // Goaldi list
	ix   int    // zero-based nonnegative Go index
}

//  vListRef.String() -- show string representation: produces (list[k])
func (lr *vListRef) String() string {
	return fmt.Sprintf("(list[%v])", lr.ix)
}

//  vListRef.Deref() -- extract value from list
func (lr *vListRef) Deref() Value {
	return lr.list.Elem(lr.ix)
}

//  vListRef.Assign -- store value in list
func (lr *vListRef) Assign(v Value) IVariable {
	list := lr.list
	n := len(list.data)
	if lr.ix >= n {
		return nil
	} else if list.rev {
		list.data[n-1-lr.ix] = v
	} else {
		list.data[lr.ix] = v
	}
	return lr
}

//  -------------------------- internal functions ---------------------

//  vList.Grow(front, name, args) -- put or push onto list
func (v *VList) Grow(front bool, name string, args ...Value) (Value, *Closure) {
	defer Traceback(name, args)
	if front != v.rev { // if wrong way around for appending
		ReverseValues(v.data)
		v.rev = !v.rev
	}
	v.data = append(v.data, args...)
	return Return(v)
}

//  vList.Snip(front, name, args) -- remove / return value from front or back
func (v *VList) Snip(front bool, name string, args ...Value) (Value, *Closure) {
	defer Traceback(name, args)
	if len(v.data) == 0 {
		return Fail()
	}
	if front != v.rev {
		r := v.data[0]
		v.data = v.data[1:]
		return Return(r)
	} else {
		n := len(v.data) - 1
		r := v.data[n]
		v.data = v.data[:n]
		return Return(r)
	}
}

//  ReverseValues(v) reverses a slice of values in place.
func ReverseValues(v []Value) {
	n := len(v)
	for i := n/2 - 1; i >= 0; i-- {
		v[i], v[n-i-1] = v[n-i-1], v[i]
	}
}

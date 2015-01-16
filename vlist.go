//	vlist.go -- VList, the Goaldi type "list"
//
//	Implementation:
//	A list is a slice of values accompanied by a "reversal" flag.
//	Lookups are simple.  Stack and queue operations append or trim the slice.
//	Switching between put and push reverses the list in place for appending.
//	All other modifications, and lookups, are O(1) in amortized cost.

package goaldi

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

//  ListType is the list instance of type type.
var ListType = NewType("list",
	DefProc(List, "list", "n,x", "create list of n copies of x"))

//  Install ListType in standard library.
func init() {
	StdLib["list"] = ListType
}

//  NewList(n, x) -- make a new list of n elements each initialized to copy(x)
//  (Unless x is a Go external that does not implement Copy().)
func NewList(n int, x Value) *VList {
	v := &VList{make([]Value, n), false}
	c, ok := x.(ICopy)
	for i := range v.data {
		if ok {
			v.data[i] = c.Copy()
		} else {
			v.data[i] = x
		}
	}
	return v
}

//  InitList(v []Value) -- construct a new VList containing the given list
//  (directly, without copying)
func InitList(v []Value) *VList {
	return &VList{v, false}
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
	for _, x := range v.data {
		fmt.Fprintf(&b, "%v,", x)
	}
	s := b.Bytes()
	s[len(s)-1] = ']'
	return string(s)
}

//  VList.Rank returns rList
func (v *VList) Rank() int {
	return rList
}

//  VList.Type -- return the list type
func (v *VList) Type() Value {
	return ListType
}

//  VList.Copy returns a new list with identical contents
func (v *VList) Copy() Value {
	return InitList(v.Export().([]Value))
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
	list := lr.list
	n := len(list.data)
	if lr.ix >= n {
		return nil
	} else if list.rev {
		return list.data[n-1-lr.ix]
	} else {
		return list.data[lr.ix]
	}
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

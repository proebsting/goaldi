//  flist.go -- list functions

package goaldi

import (
	"fmt"
)

var _ = fmt.Printf // enable debugging

//  Declare methods
var ListMethods = MethodTable([]*VProcedure{
	DefMeth((*VList).Push, "push", "x[]", "add to front"),
	DefMeth((*VList).Pop, "pop", "", "remove from front"),
	DefMeth((*VList).Get, "get", "", "remove from front"),
	DefMeth((*VList).Put, "put", "x[]", "add to end"),
	DefMeth((*VList).Pull, "pull", "", "remove from end"),
	DefMeth((*VList).Sort, "sort", "i", "return sorted copy"),
})

//  List(n, x) -- return a new list of n elements initialized to copy(x)
func List(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("list", args)
	n := int(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val())
	x := ProcArg(args, 1, NilValue)
	return Return(NewList(n, x))
}

//------------------------------------  Push:  L.push(x...)

func (v *VList) Push(args ...Value) (Value, *Closure) {
	return v.Grow(true, "L.push", args...)
}

//------------------------------------  Pop:  L.pop()

func (v *VList) Pop(args ...Value) (Value, *Closure) {
	return v.Snip(true, "L.pop", args...)
}

//------------------------------------  Get:  L.get()

func (v *VList) Get(args ...Value) (Value, *Closure) {
	return v.Snip(true, "L.get", args...)
}

//------------------------------------  Put:  L.put(x...)

func (v *VList) Put(args ...Value) (Value, *Closure) {
	return v.Grow(false, "L.put", args...)
}

//------------------------------------  Pull:  L.pull()

func (v *VList) Pull(args ...Value) (Value, *Closure) {
	return v.Snip(false, "L.pull", args...)
}

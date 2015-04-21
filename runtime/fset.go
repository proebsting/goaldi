//  fset.go -- set functions and methods

package runtime

import (
	"fmt"
)

var _ = fmt.Printf // enable debugging

//  Declare methods
var SetMethods = MethodTable([]*VProcedure{
	DefMeth((*VSet).Put, "put", "x[]", "add members"),
	DefMeth((*VSet).Insert, "insert", "x[]", "add members"),
	DefMeth((*VSet).Delete, "delete", "x[]", "remove members"),
	DefMeth((*VSet).Member, "member", "x", "test membership"),
	DefMeth((*VSet).Sort, "sort", "i", "produce sorted list"),
})

//  set(L) creates a set initialized by the values of list L.
func Set(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("set", args)
	L := ProcArg(args, 0, EMPTYLIST).(*VList)
	return Return(NewSet(L))
}

var EMPTYLIST = NewList(0, nil)

//  S.member(x) returns x if x is a member of set S;
//  otherwise it fails.
func (S *VSet) Member(args ...Value) (Value, *Closure) {
	defer Traceback("S.member", args)
	x := ProcArg(args, 0, NilValue)
	if (*S)[GoKey(x)] {
		return Return(x)
	} else {
		return Fail()
	}
}

//  S.put(x...) adds all its arguments to set S.
//  It returns S.
func (S *VSet) Put(args ...Value) (Value, *Closure) {
	defer Traceback("S.put", args)
	for _, x := range args {
		(*S)[GoKey(x)] = true
	}
	return Return(S)
}

//  S.insert(x...) adds all its arguments to set S.
//  It returns S.
func (S *VSet) Insert(args ...Value) (Value, *Closure) {
	defer Traceback("S.insert", args)
	for _, x := range args {
		(*S)[GoKey(x)] = true
	}
	return Return(S)
}

//  S.delete(x...) removes all of its arguments from set S.
//  It returns S.
func (S *VSet) Delete(args ...Value) (Value, *Closure) {
	defer Traceback("S.delete", args)
	for _, x := range args {
		delete(*S, GoKey(x))
	}
	return Return(S)
}

//  S.sort(i) returns a sorted list of the members of set S.
//  This is equivalent to [:!S:].sort(i).
func (S *VSet) Sort(args ...Value) (Value, *Closure) {
	defer Traceback("S.sort", args)
	i := ProcArg(args, 0, ONE).(Numerable).ToNumber()
	members := make([]Value, 0, len(*S))
	for k := range *S {
		members = append(members, Import(k)) // convert back from GoKey form
	}
	return InitList(members).Sort(i)
}

//  ftable.go -- table functions and methods

package goaldi

import (
	"fmt"
	"reflect"
)

var _ = fmt.Printf // enable debugging

//  Declare methods
var TableMethods = MethodTable([]*VProcedure{
	DefMeth((*VTable).Member, "member", "x", "test membership"),
	DefMeth((*VTable).Delete, "delete", "x[]", "remove entries"),
	DefMeth((*VTable).Sort, "sort", "i", "produce sorted list"),
})

//  Declare methods on Go Tables
var GoMapMethods = MethodTable([]*VProcedure{
	DefMeth(GoMapMember, "member", "x", "test membership"),
	DefMeth(GoMapDelete, "delete", "x[]", "remove entries"),
	DefMeth(GoMapSort, "sort", "i", "produce sorted list"),
})

//  Declare elemtype record for generating table values
var ElemType = NewCtor("elemtype", nil, []string{"key", "value"})

func init() {
	StdLib["elemtype"] = ElemType
}

//  table(dfval) creates a new, empty table.
func Table(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("table", args)
	dfval := ProcArg(args, 0, NilValue)
	return Return(NewTable(dfval))
}

//  T.member(k) returns k if k is an existing key in table T;
//  otherwise it fails.
func (T *VTable) Member(args ...Value) (Value, *Closure) {
	return GoMapMember(T, args...)
}

//  GoMapMember(T, k) succeeds if k is an existing key in T
func GoMapMember(T Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.member", args)
	key := ProcArg(args, 0, NilValue)
	if TrapMap(T, key).Exists() {
		return Return(key)
	} else {
		return Fail()
	}
}

//  T.delete(k...) deletes the entries with the given keys from the table T.
//  It returns T.
func (T *VTable) Delete(args ...Value) (Value, *Closure) {
	return GoMapDelete(T, args...)
}

//  GoMapDelete(T, k...) deletes the entries with the given keys from Go map T
func GoMapDelete(T Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.delete", args)
	for i := 0; i < len(args); i++ {
		key := ProcArg(args, i, NilValue)
		TrapMap(T, key).Delete()
	}
	return Return(T)
}

//  T.sort(i) returns a sorted list of elemtype(key,value) records
//  holding the contents of table T.
//  Sorting is by key if i=1 and by value if i=2.
//  T.sort(i) is equivalent to [:!T:].sort(i).
func (T *VTable) Sort(args ...Value) (Value, *Closure) {
	return GoMapSort(T.data, args...)
}

//  GoMapSort(T, i) produces [:!T:].sort(i)
func GoMapSort(T Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.sort", args)
	i := ProcArg(args, 0, ONE).(Numerable).ToNumber()
	mv := reflect.ValueOf(T)
	klist := mv.MapKeys()
	vlist := make([]Value, mv.Len())
	for i, kv := range klist {
		k := Import(kv.Interface())
		v := Import(mv.MapIndex(kv).Interface())
		vlist[i] = ElemType.New([]Value{k, v})
	}
	return InitList(vlist).Sort(i)
}

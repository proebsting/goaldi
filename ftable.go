//  ftable.go -- table functions and methods

package goaldi

import (
	"fmt"
	"math/rand"
	"reflect"
)

var _ = fmt.Printf // enable debugging

//  Declare methods
var TableMethods = MethodTable([]*VProcedure{
	DefMeth(VTable.Member, "member", "x", "test membership"),
	DefMeth(VTable.Delete, "delete", "x", "remove entry"),
	DefMeth(VTable.Sort, "sort", "", "produce sorted list"),
})

//  Declare methods on Go Tables
var GoMapMethods = MethodTable([]*VProcedure{
	DefMeth(GoMapMember, "member", "x", "test membership"),
	DefMeth(GoMapDelete, "delete", "x", "remove entry"),
	DefMeth(GoMapSort, "sort", "", "produce sorted list"),
})

//  Declare elemtype record for generating table values
var ElemType = NewCtor("elemtype", nil, []string{"key", "value"})

func init() {
	StdLib["elemtype"] = ElemType
}

//  table() creates a new, empty table.
func Table(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("table", args)
	return Return(NewTable())
}

//  T.member(k) returns k if k is an existing key in table T;
//  otherwise it fails.
func (m VTable) Member(args ...Value) (Value, *Closure) {
	return GoMapMember(m, args...)
}

//  GoMapMember(m, k) succeeds if k is an existing key in m
func GoMapMember(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.member", args)
	key := ProcArg(args, 0, NilValue)
	if TrapMap(m, key).Exists() {
		return Return(key)
	} else {
		return Fail()
	}
}

//  T.delete(k) deletes the entry with key k, if any, from the table T.
//  It returns T.
func (m VTable) Delete(args ...Value) (Value, *Closure) {
	return GoMapDelete(m, args...)
}

//  GoMapDelete(m, k) deletes the entry in Go map m, if any, with key k
func GoMapDelete(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.delete", args)
	key := ProcArg(args, 0, NilValue)
	TrapMap(m, key).Delete()
	return Return(m)
}

//  T.sort(i) returns a sorted list of elemtype(key,value) records
//  holding the contents of table T.
//  Sorting is by key if i=1 and by value if i=2.
//  T.sort(i) is equivalent to [:!T:].sort(i).
func (m VTable) Sort(args ...Value) (Value, *Closure) {
	return GoMapSort(m, args...)
}

//  GoMapSort(m, i) produces [:!T:].sort(i)
func GoMapSort(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.sort", args)
	i := ProcArg(args, 0, ONE).(Numerable).ToNumber()
	mv := reflect.ValueOf(m)
	klist := mv.MapKeys()
	vlist := make([]Value, mv.Len())
	for i, kv := range klist {
		k := Import(kv.Interface())
		v := Import(mv.MapIndex(kv).Interface())
		vlist[i] = ElemType.New([]Value{k, v})
	}
	return InitList(vlist).Sort(i)
}

//  -------------------------- key/value pairs ---------------------

//  ChooseMap returns a key/value pair from any Goaldi table or Go map
func ChooseMap(m interface{} /*anymap*/) Value {
	mv := reflect.ValueOf(m)
	klist := mv.MapKeys()
	n := len(klist)
	if n == 0 { // if map empty
		return nil // fail
	}
	i := rand.Intn(n)
	k := Import(klist[i].Interface())
	v := Import(mv.MapIndex(klist[i]).Interface())
	return ElemType.New([]Value{k, v})
}

//  DispenseMap generates key/value pairs for any Goaldi table or Go map
func DispenseMap(m interface{} /*anymap*/) (Value, *Closure) {
	mv := reflect.ValueOf(m)
	klist := mv.MapKeys()
	i := -1
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		i++
		if i < len(klist) {
			k := Import(klist[i].Interface())
			v := Import(mv.MapIndex(klist[i]).Interface())
			return ElemType.New([]Value{k, v}), c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

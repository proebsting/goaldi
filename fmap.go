//  fmap.go -- map functions and methods

package goaldi

import (
	"math/rand"
	"reflect"
)

//  Declare methods
var MapMethods = map[string]interface{}{
	"type":   VMap.Type,
	"copy":   VMap.Copy,
	"string": VMap.String,
	"image":  VMap.GoString,
	"member": VMap.Member,
	"delete": VMap.Delete,
	"sort":   VMap.Sort,
}

//  VMap.Field implements method calls
func (m VMap) Field(f string) Value {
	return GetMethod(MapMethods, m, f)
}

//  init() declares the constructor function
func init() {
	// Goaldi procedures
	LibProcedure("map", Map)
}

//  Declare methods on Go Maps
var GoMapMethods = map[string]interface{}{
	"member": GoMapMember,
	"delete": GoMapDelete,
	"sort":   GoMapSort,
}

//  Map() returns a new map
func Map(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("map", args)
	return Return(NewMap())
}

//  VMap.Member(k) succeeds if k is an existing key
func (m VMap) Member(args ...Value) (Value, *Closure) {
	return GoMapMember(m, args...)
}

//  GoMapMember(M, k) succeeds if k is an existing key in m
func GoMapMember(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("M.member", args)
	key := ProcArg(args, 0, NilValue)
	if TrapMap(m, key).Exists() {
		return Return(key)
	} else {
		return Fail()
	}
}

//  VMap.Delete(k) deletes the entry, if any, with key k
func (m VMap) Delete(args ...Value) (Value, *Closure) {
	return GoMapDelete(m, args...)
}

//  GoMapDelete(M, k) deletes the entry in map m, if any, with key k
func GoMapDelete(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("M.delete", args)
	key := ProcArg(args, 0, NilValue)
	TrapMap(m, key).Delete()
	return Return(m)
}

//  VMap.Sort(i) produces [:!M:].sort(i)
func (m VMap) Sort(args ...Value) (Value, *Closure) {
	return GoMapSort(m, args...)
}

//  GoMapSort(M, i) produces [:!M:].sort(i)
func GoMapSort(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("M.sort", args)
	i := ProcArg(args, 0, ONE).(Numerable).ToNumber()
	mv := reflect.ValueOf(m)
	klist := mv.MapKeys()
	vlist := make([]Value, mv.Len())
	for i, kv := range klist {
		k := Import(kv.Interface())
		v := Import(mv.MapIndex(kv).Interface())
		vlist[i] = kvstruct.New([]Value{k, v})
	}
	return InitList(vlist).Sort(i)
}

//  -------------------------- key/value pairs ---------------------

//  kvstruct defines the {key,value} struct returned by ?M and !M
var kvstruct = NewDefn("mapElem", []string{"key", "value"})

//  ChooseMap returns a key/value pair from any Go (or Goaldi) map
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
	return kvstruct.New([]Value{k, v})
}

//  DispenseMap generates key/value pairs for any Go (or Goaldi) map
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
			return kvstruct.New([]Value{k, v}), c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

//  TakeMap returns a key/value pair after removing the entry from the map
func TakeMap(m interface{} /*anymap*/) Value {
	mv := reflect.ValueOf(m)
	klist := mv.MapKeys()
	n := len(klist)
	if n == 0 { // if map empty
		return nil // fail
	}
	i := rand.Intn(n)
	kv := klist[i]
	k := Import(kv.Interface())
	v := Import(mv.MapIndex(klist[i]).Interface())
	mv.SetMapIndex(kv, reflect.Value{}) // delete
	return kvstruct.New([]Value{k, v})
}

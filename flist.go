//  flist.go -- list functions

package goaldi

import (
	"sort"
)

//  Declare methods
var ListMethods = map[string]interface{}{
	"type":  (*VList).Type,
	"copy":  (*VList).Copy,
	"image": Image,
	"push":  (*VList).Push,
	"pop":   (*VList).Pop,
	"get":   (*VList).Get,
	"put":   (*VList).Put,
	"pull":  (*VList).Pull,
	"sort":  (*VList).Sort,
}

//  VList.Field implements methods
func (v *VList) Field(f string) Value {
	return GetMethod(ListMethods, v, f)
}

//  Declare constructor function
func init() {
	LibProcedure("list", List)
}

//  List(n, x) -- return a new list of n elements initialize to x
func List(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("list", a)
	n := int(ProcArg(a, 0, ZERO).(Numerable).ToNumber().Val())
	x := ProcArg(a, 1, NilValue)
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

//------------------------------------  Sort:  L.sort(i)

//  L.Sort(i) -- sort list L on field i (default i=1; use i=0 for no field)
func (v *VList) Sort(args ...Value) (Value, *Closure) {
	defer Traceback("sort", args)
	i := int(ProcArg(args, 0, ONE).(Numerable).ToNumber().Val()) - 1
	if i < 0 {
		panic(&RunErr{"Nonpositive field index", args[0]})
	}
	d := &lsort{make([]Value, len(v.data)), i}
	copy(d.v, v.data)
	sort.Stable(d)
	return Return(InitList(d.v))
}

//  a list to be sorted, with field index
type lsort struct {
	v []Value // Goaldi values
	f int     // zero-based field index, or -1
}

//  sort interface functions
func (a *lsort) Len() int           { return len(a.v) }
func (a *lsort) Swap(i, j int)      { a.v[i], a.v[j] = a.v[j], a.v[i] }
func (a *lsort) Less(i, j int) bool { return LT(a.v[i], a.v[j], a.f) }

//  LT(x, y, i) -- return x < y on field i
func LT(x Value, y Value, i int) bool {
	rx := rank(x)
	ry := rank(y)
	if rx != ry { // if different types
		return rx < ry // order by type rank
	}
	// both values have the same type
	switch ry {
	case rNumber:
		return x.(*VNumber).Val() < y.(*VNumber).Val()
	case rString:
		return x.(*VString).compare(y.(*VString)) < 0
	case rFile:
		return x.(*VFile).Name < y.(*VFile).Name
	case rDefn:
		return x.(*VDefn).Name < y.(*VDefn).Name
	case rMethB:
		return x.(*VMethB).Name < y.(*VMethB).Name
	case rProc:
		return x.(*VProcedure).Name < y.(*VProcedure).Name
	case rStruct:
		xs := x.(*VStruct)
		ys := y.(*VStruct)
		if xs.Defn != ys.Defn {
			// different struct types; order by struct name
			return xs.Defn.Name < ys.Defn.Name
		}
		if i >= 0 && len(xs.Data) > i && len(ys.Data) > i {
			// both sides have an item i
			return LT(xs.Data[i], ys.Data[i], -1)
		} else {
			// put missing one first; otherwise #%#% we don't care
			return len(xs.Data) < len(ys.Data)
		}
	case rList:
		xl := x.(*VList)
		yl := y.(*VList)
		if i >= 0 && len(xl.data) > i && len(yl.data) > i {
			xr := &vListRef{xl, i}
			yr := &vListRef{yl, i}
			return LT(xr.Deref(), yr.Deref(), -1)
		} else {
			// put missing one first; otherwise #%#% we don't care
			return len(xl.data) < len(yl.data)
		}
	case rMap:
		return len(x.(VMap)) < len(y.(VMap)) //#%#% got anything better?
	default:
		return false //#%#% not comparable?
	}
}

//  ranking of types for sorting
const (
	rNil = iota
	rNumber
	rString
	rFile
	rDefn
	rMethB
	rProc
	rList
	rMap
	rStruct
	rExternal
)

//  rank(x) -- return sort ranking for the type of x
func rank(x Value) int {
	if t, ok := x.(IRank); ok {
		return t.Rank()
	} else {
		return rExternal
	}
}

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

//  ranking of types for sorting
const (
	rNil = iota
	rNumber
	rString
	rFile
	rDefn
	rProc
	rList
	rMap
	rStruct
	rExternal
)

//  a list to be sorted, with field index
type sortdata struct {
	data []Value // Goaldi values
	f    int     // field index or 0
}

//  L.Sort(i) -- sort list L on field i		#%#% ignores i
func (v *VList) Sort(args ...Value) (Value, *Closure) {
	defer Traceback("sort", args)
	i := int(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val())
	d := &sortdata{make([]Value, len(v.data)), i}
	copy(d.data, v.data)
	sort.Sort(d)
	return Return(InitList(d.data))
}

//  sort interface functions
func (a *sortdata) Len() int      { return len(a.data) }
func (a *sortdata) Swap(i, j int) { a.data[i], a.data[j] = a.data[j], a.data[i] }
func (a *sortdata) Less(i, j int) bool {
	ri := rank(a.data[i])
	rj := rank(a.data[j])
	if ri != rj {
		return ri < rj
	}
	switch ri {
	case rNumber:
		return a.data[i].(*VNumber).Val() < a.data[j].(*VNumber).Val()
	case rString:
		return a.data[i].(*VString).String() < a.data[j].(*VString).String()
	default:
		return false //#%#% not comparable?
	}
}

//  rank(x) -- return sort ranking for the type of x
func rank(x Value) int {
	if t, ok := x.(IRank); ok {
		return t.Rank()
	} else {
		return rExternal
	}
}

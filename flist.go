//  flist.go -- list functions

package goaldi

import (
	"fmt"
	"math/rand"
	"sort"
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
	DefMeth((*VList).Shuffle, "shuffle", "", "return randomized copy"),
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

//------------------------------------  Shuffle:  L.shffle()

func (v *VList) Shuffle(args ...Value) (Value, *Closure) {
	defer Traceback("shuffle", args)
	n := len(v.data)
	d := make([]Value, n, n)
	copy(d, v.data)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
	return Return(InitList(d))
}

//------------------------------------  Sort: L.Sort(i)

//  L.Sort(i) -- sort list L on field i (default i=1)
func (v *VList) Sort(args ...Value) (Value, *Closure) {
	defer Traceback("sort", args)
	i := int(ProcArg(args, 0, ONE).(Numerable).ToNumber().Val()) - 1
	if i < 0 {
		panic(NewExn("Nonpositive field index", args[0]))
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
	if v, ok := x.(ICore); ok { // if standard types (not external)
		return v.Before(y, i) // use the type's comparison function
	}
	return false // otherwise no ordering defined
}

//  rank(x) -- return sort ranking for the type of x
func rank(x Value) int {
	if t, ok := x.(IType); ok {
		return t.Type().Rank()
	} else {
		return rExternal
	}
}

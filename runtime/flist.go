//  flist.go -- list functions

package runtime

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

//  list(size, x) builds and returns a new list of the given size
//  with each element initialized to a copy of x.
func List(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("list", args)
	n := int(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val())
	x := ProcArg(args, 1, NilValue)
	return Return(NewList(n, x))
}

//  L.push(x...) adds its arguments, in order, to the beginning of list L.
//  The last argument thus ends up as the first element of L.
func (v *VList) Push(args ...Value) (Value, *Closure) {
	return v.Grow(true, "L.push", args...)
}

//  L.pop() removes the first element from list L
//  and returns the element's value.
func (v *VList) Pop(args ...Value) (Value, *Closure) {
	return v.Snip(true, "L.pop", args...)
}

//  L.get() removes the first element from list L
//  and returns the element's value.
func (v *VList) Get(args ...Value) (Value, *Closure) {
	return v.Snip(true, "L.get", args...)
}

//  L.put(x...) adds its arguments, in order, to the end of list L.
//  The last argument becomes the final element of L.
func (v *VList) Put(args ...Value) (Value, *Closure) {
	return v.Grow(false, "L.put", args...)
}

//  L.pull() removes the final element from list L
//  and returns the element's value.
func (v *VList) Pull(args ...Value) (Value, *Closure) {
	return v.Snip(false, "L.pull", args...)
}

//  L.shuffle() returns a copy of list L in which the elements
//  have been randomly reordered.
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

//  L.sort(i) returns a copy of list L in which the elements have been sorted.
//  Values are ordered first by type, then within types by their values.
//  Among lists and among records of the same type,
//  ordering is based on field i.
//  Lists with no element i are sorted ahead of lists that have one.
//  The value i defaults to 1 and must be strictly positive.
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

//  rank(x) -- return sort ranking for the type of x.
//  Every type declares its own type independently, rather like the
//  atomic weights of the elements.  There is no central coordinator.
func rank(x Value) int {
	if t, ok := x.(IType); ok {
		return t.Type().Rank()
	} else {
		return rExternal
	}
}

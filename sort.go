//  sort.go -- code dealing with ordering and sorting

package goaldi

import (
	"fmt"
	"sort"
)

var _ = fmt.Printf // enable debugging

//  ranking of types for sorting
const (
	rNil = iota
	rTrapped
	rType
	rNumber
	rString
	rFile
	rChannel
	rDefn
	rMethVal
	rProc
	rList
	rTable
	rRecord
	rExternal
)

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
	// both values have the same type
	switch ry {
	case rType:
		return x.(IRank).Rank() < y.(IRank).Rank()
	case rNumber:
		return x.(*VNumber).Val() < y.(*VNumber).Val()
	case rString:
		return x.(*VString).compare(y.(*VString)) < 0
	case rFile:
		return x.(*VFile).Name < y.(*VFile).Name
	case rDefn:
		return x.(*VDefn).Name < y.(*VDefn).Name
	case rMethVal:
		return x.(*VMethVal).Proc.Name < y.(*VMethVal).Proc.Name
	case rProc:
		return x.(*VProcedure).Name < y.(*VProcedure).Name
	case rRecord:
		xs := x.(*VRecord)
		ys := y.(*VRecord)
		if xs.Defn != ys.Defn {
			// different record types; order by type name
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
	case rTable:
		return len(x.(VTable)) < len(y.(VTable)) //#%#% got anything better?
	default:
		return false //#%#% not comparable?
	}
}

//  rank(x) -- return sort ranking for the type of x
func rank(x Value) int {
	if t, ok := x.(IType); ok {
		return t.Type().Rank()
	} else {
		return rExternal
	}
}

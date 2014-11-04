//  oelem.go -- operations for accessing elements of structures

package goaldi

import (
	"reflect"
)

//  Index(lval,x,y) calls x.Index(lval, y) or falls back to reflection.
func Index(lval IVariable, x Value, y Value) Value {
	if t, ok := x.(IIndex); ok {
		return t.Index(lval, y)
	}
	xv := reflect.ValueOf(x)
	if xv.Kind() != reflect.Slice && xv.Kind() != reflect.Array {
		panic(&RunErr{"Wrong type for indexing", x})
	}
	n := xv.Len()
	i := GoIndex(int(y.(Numerable).ToNumber().Val()), n)
	if i >= n {
		return nil // out of bounds
	}
	if lval == nil {
		return Import(xv.Index(i).Interface())
	} else {
		return TrapValue(xv.Index(i))
	}
}

//  IField -- interface for x.Field(id), used by x.id
type IField interface {
	Field(string) Value
}

//  Field(x,s) calls x.Field(s) or (#%#%TBD) falls back to reflection.
func Field(x Value, s string) Value {
	if t, ok := x.(IField); ok {
		return t.Field(s)
	}
	//#%#% try looking up field in Go struct or map using reflection.
	return nil
}

//  GoIndex(i, n) -- return Go index for i into length n, or n+1 if out of range
//  i is a 1-based index and may be nonpositive.  i==n or i==0 is in range.
//  The caller may want to check for result<n or result<=n depending on context.
func GoIndex(i int, n int) int {
	if i > 0 {
		i-- // convert to zero-based
	} else {
		i = n + i // count backwards from end
	}
	if i >= 0 && i <= n {
		return i // index is valid
	} else {
		return n + 1 // index is out of range
	}
}

//  oelem.go -- operations for accessing elements of structures

package goaldi

import (
	"math/rand"
	"reflect"
)

//  IField -- interface for x.Field(id), used by x.id
type IField interface {
	Field(string) Value
}

//  Field(x,s) calls x.Field(s) or falls back to reflection.
func Field(x Value, s string) Value {
	// first check to see if this value implements Field()
	if t, ok := x.(IField); ok {
		return t.Field(s)
	}
	// using reflection, peek inside interface and/or pointer to actual value
	xv := reflect.ValueOf(x)
	if xv.Kind() == reflect.Interface {
		xv = xv.Elem()
	}
	// look for an explicitly implemented method
	if m, ok := xv.Type().MethodByName(s); ok {
		return GoMethod(x, s, m)
	}
	if xv.Kind() == reflect.Ptr {
		xv = xv.Elem()
	}
	// what kind of a Go value is this?
	switch xv.Kind() {
	case reflect.Struct:
		// check for field of arbitrary struct type
		if f := xv.FieldByName(s); f.IsValid() {
			return TrapValue(f)
		}
	case reflect.Chan:
		// we have implicit methods for any kind of map
		return GetMethod(GoChanMethods, x, s)
	case reflect.Map:
		// we have implicit methods for any kind of map
		return GetMethod(GoMapMethods, x, s)
	}
	// nothing found
	panic(&Exception{"Field not found: " + s, x})
}

//  Index(lval,x,y) calls x.Index(lval, y) or falls back to reflection.
func Index(lval Value, x Value, y Value) Value {
	if t, ok := x.(IIndex); ok {
		return t.Index(lval, y)
	}
	xv := reflect.ValueOf(x)
	if xv.Kind() == reflect.Map {
		return TrapMap(x, y)
	}
	if xv.Kind() != reflect.Slice && xv.Kind() != reflect.Array {
		panic(&Exception{"Wrong type for indexing", x})
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

//  Choose(lval, x) calls x.Choose(lval) or uses Size() and Index().
func Choose(lval Value, x Value) Value {
	if t, ok := x.(IChoose); ok {
		return t.Choose(lval)
	}
	if reflect.ValueOf(x).Kind() == reflect.Map {
		return ChooseMap(x)
	}
	n := int(Size(x).(*VNumber).Val())
	i := rand.Intn(n) + 1 // +1 for 1-based indexing
	return Index(lval, x, NewNumber(float64(i)))
}

//  Dispense(lval, x) calls x.Dispense(lval) or steps through Go values.
func Dispense(lval Value, x Value) (Value, *Closure) {
	if t, ok := x.(IDispense); ok {
		return t.Dispense(lval)
	}
	k := reflect.ValueOf(x).Kind()
	if k == reflect.Chan {
		return DispenseChan(x)
	} else if k == reflect.Map {
		return DispenseMap(x)
	}
	i := 0.0
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		i++
		return Index(lval, x, NewNumber(i)), c
	}}
	return c.Resume()
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

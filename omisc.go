//  omisc.go -- miscellaneous runtime operations

package goaldi

import (
	"reflect"
)

//  IIdentical -- interface for a.Identical(b), used by a===b
//  Must be implemented for types where === is not just a pointer match
type IIdentical interface {
	Identical(Value) Value
}

var _ IIdentical = NewNumber(1)   // confirm implementation by VNumber
var _ IIdentical = NewString("a") // confirm implementation by VString

//  Identical(a,b) implements the === operator.
//  NotIdentical(a,b) implements the ~=== operator.
//  Both call a.Identical(b) if implemented (interface IIdentical).
func Identical(a, b Value) Value {
	if _, ok := a.(IIdentical); ok {
		return a.(IIdentical).Identical(b)
	} else if a == b {
		return b
	} else {
		return nil
	}
}

func NotIdentical(a, b Value) Value {
	if Identical(b, a) != nil {
		return nil
	} else {
		return b
	}
}

//  ISize -- interface for x.Size(), used by *x
type ISize interface {
	Size() Value
}

//  Size(x) calls x.Size() or falls back to reflection.
//  It panics on an inappropriate argument type.
func Size(x Value) Value {
	if t, ok := x.(ISize); ok {
		return t.Size()
	} else {
		return NewNumber(float64(reflect.ValueOf(x).Len()))
	}
}

//  ITake -- interface for x.Take(), used by @x
type ITake interface { // @x
	Take() Value
}

//  Take(x) calls x.Take() or uses reflection for a map of any type.
//  It panics on an inappropriate argument type.
func Take(x Value) Value {
	if t, ok := x.(ITake); ok {
		return t.Take()
	} else if reflect.ValueOf(x).Kind() == reflect.Map {
		return TakeMap(x)
	} else {
		return x.(ITake).Take() // provoke panic
	}
}

//  VNumber.ICall -- implement i(e1, e2, e3...)
func (v *VNumber) Call(env *Env, args ...Value) (Value, *Closure) {
	i := GoIndex(int(v.Val()), len(args))
	if i < len(args) {
		return Return(args[i])
	} else {
		return Fail()
	}
}

//  ToBy -- implement "e1 to e2 by e3"
func ToBy(e1 Value, e2 Value, e3 Value) (Value, *Closure) {
	n1 := e1.(Numerable).ToNumber()
	if n1 == nil {
		panic(&RunErr{"ToBy: e1 bad", e1})
	}
	n2 := e2.(Numerable).ToNumber()
	if n2 == nil {
		panic(&RunErr{"ToBy: e2 bad", e2})
	}
	n3 := e3.(Numerable).ToNumber()
	if n3 == nil {
		panic(&RunErr{"ToBy: e3 bad", e3})
	}
	if *n3 == 0 {
		panic(&RunErr{"ToBy: by 0", nil})
	}
	v1 := *n1
	v2 := *n2
	v3 := *n3
	v1 -= v3
	var f *Closure
	f = &Closure{func() (Value, *Closure) {
		v1 += v3
		if (v3 > 0 && v1 <= v2) || (v3 < 0 && v1 >= v2) {
			return NewNumber(float64(v1)), f
		} else {
			return Fail()
		}
	}}
	return f.Resume()
}

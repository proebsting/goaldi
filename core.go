//  core.go -- some fundamental functions

package goaldi

import (
	"fmt"
)

//  Image() makes a Goaldi string image of any value
//  It uses the same String() function as Sprintf("%v").

type IImage interface {
	String() string
}

func Image(v Value) Value {
	return NewString(fmt.Sprintf("%v", v))
}

//  Type() returns the type of any value

type IType interface {
	Type() Value
}

func Type(v Value) Value {
	switch t := v.(type) {
	case IExternal:
		return NewString(t.ExternalType())
	case IType:
		return t.Type()
	default:
		return type_external
	}
}

var type_external = NewString("external")

//  Identical(a,b) implements the === operator.
//  NotIdentical(a,b) implements the ~=== operator.
//  Both call a.Identical(b) if implemented (interface IIdentical).

type IIdentical interface {
	Identical(Value) Value
}

var _ IIdentical = NewNumber(1)   // confirm implementation by VNumber
var _ IIdentical = NewString("a") // confirm implementation by VString

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

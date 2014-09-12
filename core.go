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

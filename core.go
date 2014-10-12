//  core.go -- some fundamental functions

package goaldi

import (
	"fmt"
)

func Image(v Value) Value {
	return NewString(fmt.Sprintf("%v", v))
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

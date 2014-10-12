//  fmisc.go -- miscellaneous functions

package goaldi

import (
	"fmt"
)

func init() {
	LibGoFunc("exit", os.Exit)
	LibGoFunc("remove", os.Remove)
	LibProcedure("image", Image)
	LibProcedure("type", Type)
}

//  Image(v) -- return string image of value v
func Image(env *Env, a ...Value) (Value, *Closure) {
	v := NilValue
	if len(a) > 0 {
		v = a[0]
	}
	return Return(NewString(fmt.Sprintf("%#v", v)))
}

//  Type(v) -- return the name of v's type, as a string
func Type(env *Env, a ...Value) (Value, *Closure) {
	v := NilValue
	if len(a) > 0 {
		v = a[0]
	}
	switch t := v.(type) {
	case IExternal:
		return Return(NewString(t.ExternalType()))
	case IType:
		return Return(t.Type())
	default:
		return Return(type_external)
	}
}

var type_external = NewString("external")

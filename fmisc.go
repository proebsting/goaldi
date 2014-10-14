//  fmisc.go -- standard library setup and miscellaneous functions

package goaldi

import (
	"fmt"
	"os"
)

//  StdLib is the set of procedures available at link time
var StdLib = make(map[string]*VProcedure)

//  LibProcedure registers a standard library procedure taking Goaldi arguments.
//  This must be done before linking (e.g. via init func) to be effective.
func LibProcedure(name string, p Procedure) {
	StdLib[name] = NewProcedure(name, p)
}

//  LibGoFunc registers a Go function as a standard library procedure.
//  This must be done before linking (e.g. via init func) to be effective.
func LibGoFunc(name string, f interface{}) {
	StdLib[name] = GoProcedure(name, f)
}

//  Miscellaneous standard library procedures
func init() {
	LibGoFunc("remove", os.Remove)
	LibProcedure("image", Image)
	LibProcedure("type", Type)
	LibProcedure("exit", Exit)
}

//  ProcArg(a,i,d) -- return procedure argument a[i], defaulting to d
func ProcArg(a []Value, i int, d Value) Value {
	if i < len(a) && a[i] != NilValue {
		return a[i]
	} else {
		return d
	}
}

//  Image(v) -- return string image of value v
func Image(env *Env, a ...Value) (Value, *Closure) {
	return Return(NewString(fmt.Sprintf("%#v", ProcArg(a, 0, NilValue))))
}

//  Type(v) -- return the name of v's type, as a string
func Type(env *Env, a ...Value) (Value, *Closure) {
	v := ProcArg(a, 0, NilValue)
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

//  Exit(n) -- terminate program
func Exit(env *Env, a ...Value) (Value, *Closure) {
	Shutdown(int(ProcArg(a, 0, ZERO).(Numerable).ToNumber().Val()))
	return Fail() // NOTREACHED
}

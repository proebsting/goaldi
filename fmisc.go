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

//  Miscellaneous library procedures
func init() {
	// Goaldi procedures
	LibProcedure("image", Image)
	LibProcedure("type", Type)
	LibProcedure("exit", Exit)
	// Go library functions
	LibGoFunc("getenv", os.Getenv)
	LibGoFunc("setenv", os.Setenv)
	LibGoFunc("expandenv", os.ExpandEnv)
	LibGoFunc("clearenv", os.Clearenv)
	LibGoFunc("hostname", os.Hostname)
	LibGoFunc("getpid", os.Getpid)
	LibGoFunc("getppid", os.Getppid)
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
	defer Traceback("image", a)
	return Return(NewString(fmt.Sprintf("%#v", ProcArg(a, 0, NilValue))))
}

//  Type(v) -- return the name of v's type, as a string
func Type(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("type", a)
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
	defer Traceback("exit", a)
	Shutdown(int(ProcArg(a, 0, ZERO).(Numerable).ToNumber().Val()))
	return Fail() // NOTREACHED
}

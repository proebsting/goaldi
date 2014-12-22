//  fmisc.go -- standard library setup and miscellaneous functions

package goaldi

import (
	"archive/zip"
	"fmt"
	"os"
	"time"
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
	LibProcedure("type", Type)
	LibProcedure("copy", Copy)
	LibProcedure("image", Image)
	LibProcedure("noresult", NoResult)
	LibProcedure("nilresult", NilResult)
	LibProcedure("errresult", ErrResult)
	LibProcedure("exit", Exit)
	LibProcedure("runerr", Runerr)
	LibProcedure("sleep", Sleep)
	// Go library functions
	LibGoFunc("getenv", os.Getenv)
	LibGoFunc("setenv", os.Setenv)
	LibGoFunc("environ", os.Environ)
	LibGoFunc("expandenv", os.ExpandEnv)
	LibGoFunc("clearenv", os.Clearenv)
	LibGoFunc("hostname", os.Hostname)
	LibGoFunc("getpid", os.Getpid)
	LibGoFunc("getppid", os.Getppid)
	// Heavy-duty package interfaces
	LibGoFunc("zipreader", zip.OpenReader)
}

//  Type(v) -- return the name of v's type, as a string
func Type(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("type", args)
	v := ProcArg(args, 0, NilValue)
	if t, ok := v.(IType); ok {
		return Return(t.Type())
	} else {
		return Return(type_external)
	}
}

var type_external = NewString("external")

//  Copy(v) -- return a copy of v (or just v if a simple value).
//  The type of v *must* implement ICopy.
func Copy(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("copy", args)
	v := ProcArg(args, 0, NilValue)
	return Return(v.(ICopy).Copy())
}

//  Image(v) -- return string image of value v
func Image(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("image", args)
	return Return(NewString(fmt.Sprintf("%#v", ProcArg(args, 0, NilValue))))
}

//  NoResult() -- fail immediately
func NoResult(env *Env, args ...Value) (Value, *Closure) {
	return Fail()
}

//  NilResult() -- return nilresult
func NilResult(env *Env, args ...Value) (Value, *Closure) {
	return Return(NilValue)
}

//  ErrResult() -- return &error
func ErrResult(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("errresult", args)
	return Return(env.VarMap["error"])
}

//  Sleep(n) -- delay execution for n seconds (may be fractional)
func Sleep(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("sleep", args)
	v := ProcArg(args, 0, ONE).(Numerable).ToNumber()
	n := v.Val()
	d := time.Duration(n * float64(time.Second))
	time.Sleep(d)
	return Return(v)
}

//  Exit(n) -- terminate program
func Exit(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("exit", args)
	Shutdown(int(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val()))
	return Fail() // NOTREACHED
}

//  Runerr(x, v) -- terminate with error x and offending value v
func Runerr(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("runerr", args)
	x := ProcArg(args, 0, err_fatal)
	v := ProcArg(args, 1, nil)
	if n, ok := x.(*VNumber); ok {
		x = NewString(fmt.Sprintf("Fatal error %v", n))
	}
	panic(&RunErr{fmt.Sprintf("%v", x), v})
}

var err_fatal = NewString("Unspecified fatal error")

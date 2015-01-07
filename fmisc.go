//  fmisc.go -- standard library setup and miscellaneous functions

package goaldi

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
)

//  StdLib is the set of procedures available at link time
var StdLib = make(map[string]*VProcedure)

//  LibProcedure registers a standard library procedure taking Goaldi arguments.
//  This must be done before linking (e.g. via init func) to be effective.
func LibProcedure(name string, p Procedure) {
	StdLib[name] = NewProcedure(name, nil, p, p)
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
	LibProcedure("throw", Throw)
	LibProcedure("sleep", Sleep)
	LibProcedure("date", Date)
	LibProcedure("time", Time)
	LibProcedure("now", Now)
	LibProcedure("duration", Duration)
	LibProcedure("cputime", CPUtime)
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

//  ShowLibrary(f) -- list all library functions on file f
func ShowLibrary(f io.Writer) {
	fmt.Fprintln(f)
	fmt.Fprintln(f, "Standard Library")
	fmt.Fprintln(f, "------------------------------")
	for k := range SortedKeys(StdLib) {
		v := StdLib[k]
		fmt.Fprintf(f, "%-12s %s\n", k, v.ImplBy())
	}
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
	return Return(env.VarMap["error"].(IVariable).Deref())
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

//  Date() -- return current date in the form yyyy/mm/dd
func Date(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("date", args)
	return Return(NewString(time.Now().Format("2006/01/02")))
}

//  Time() -- return current time in the form hh:mm:ss
func Time(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("time", args)
	return Return(NewString(time.Now().Format("15:04:05")))
}

//  Now() -- return current time as a Go.Time struct for user formatting
func Now(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("now", args)
	return Return(time.Now())
}

//  Duration(x) -- convert x to a duration, or fail
//  If x is a string, it is passed directly to time.ParseDuration.
//  If x is a number, "s" is appended to interpret it as an interval in seconds.
func Duration(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("duration", args)
	v := ProcArg(args, 0, ZERO)
	s := ""
	if n, ok := v.(*VNumber); ok {
		s = n.String() + "s"
	} else {
		s = v.(Stringable).ToString().String()
	}
	d, err := time.ParseDuration(s)
	if err == nil {
		return Return(d)
	} else {
		return Fail()
	}
}

//  CPUtime() -- return u+s CPU usage in seconds (may be fractional)
func CPUtime(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("cputime", args)
	var ustruct syscall.Rusage
	err := syscall.Getrusage(0, &ustruct)
	if err != nil {
		panic(err)
	}
	user := time.Duration(syscall.TimevalToNsec(ustruct.Utime))
	sys := time.Duration(syscall.TimevalToNsec(ustruct.Stime))
	total := user + sys
	return Return(NewNumber(total.Seconds()))
}

//  Exit(n) -- terminate program
func Exit(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("exit", args)
	Shutdown(int(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val()))
	return Fail() // NOTREACHED
}

//  Throw(x, v) -- terminate with error x and offending value v
func Throw(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("throw", args)
	x := ProcArg(args, 0, err_fatal)
	v := ProcArg(args, 1, NilValue)
	if len(args) < 2 {
		v = nil // distingish no argument from explicit %nil
	}
	if n, ok := x.(*VNumber); ok {
		x = NewString(fmt.Sprintf("Fatal error %v", n))
	}
	panic(&Exception{fmt.Sprintf("%v", x), v})
}

var err_fatal = NewString("Unspecified fatal error")

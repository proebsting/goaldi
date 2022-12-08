//  fmisc.go -- standard library setup and miscellaneous functions

package runtime

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"reflect"
	"syscall"
	"time"
)

// Miscellaneous library procedures
func init() {
	// Goaldi procedures
	DefLib(Copy, "copy", "x", "copy value")
	DefLib(Image, "image", "x", "return detailed string image")
	DefLib(NoResult, "noresult", "e", "fail immediately")
	DefLib(NilResult, "nilresult", "e", "return nil")
	DefLib(ErrResult, "errresult", "e", "return e")
	DefLib(Exit, "exit", "i", "terminate program with exit status")
	DefLib(Throw, "throw", "e,x[]", "terminate with error and offending values")
	DefLib(Sleep, "sleep", "n", "pause execution momentarily")
	DefLib(Date, "date", "", "return the current date")
	DefLib(Time, "time", "", "return the current time")
	DefLib(Now, "now", "", "return the current instant as a Go Time struct")
	DefLib(Duration, "duration", "x", "convert value to a Go Duration struct")
	DefLib(CPUtime, "cputime", "", "return total processor time used")
	// Go library functions
	GoLib(os.Getenv, "getenv", "key", "read environment variable")
	GoLib(os.Setenv, "setenv", "key,value", "set environment variable")
	GoLib(os.Environ, "environ", "", "get list of environment variables")
	GoLib(os.Clearenv, "clearenv", "", "delete all environment variables")
	GoLib(os.Hostname, "hostname", "", "get host machine name")
	GoLib(os.Getpid, "getpid", "", "get process ID")
	GoLib(os.Getppid, "getppid", "", "get parent process ID")
	GoLib(exec.Command, "command", "name,args[]", "build struct to run command")
}

// copy(x) returns a copy of x if x is a structure,
// or just x itself if x is a simple value.
// This is a shallow copy; nested structures are not duplicated.
func Copy(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("copy", args)
	x := ProcArg(args, 0, NilValue)
	if v, ok := x.(ICopy); ok {
		return Return(v.Copy())
	}
	// doesn't implement Copy(); must be an external
	y := reflect.Indirect(reflect.New(reflect.TypeOf(x)))
	y.Set(reflect.ValueOf(x))
	return Return(y.Interface())
}

// image(x) returns a string image of x.
// This is the same conversion applied by sprintf("%#v",x)
// and is typically more verbose and detailed than the result of string(x).
func Image(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("image", args)
	return Return(NewString(fmt.Sprintf("%#v", ProcArg(args, 0, NilValue))))
}

// noresult(e) fails immediately, ignoring e.
// It is suitable for use as a catch handler.
func NoResult(env *Env, args ...Value) (Value, *Closure) {
	return Fail()
}

// nilresult(e) returns nil, ignoring e.
// It is suitable for use as a catch handler.
func NilResult(env *Env, args ...Value) (Value, *Closure) {
	return Return(NilValue)
}

// errresult(e) returns its argument e.
// It is suitable for use as a catch handler.
func ErrResult(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("errresult", args)
	return Return(ProcArg(args, 0, NilValue))
}

// sleep(n) delays execution for n seconds, which may be a fractional value.
// If n is nil, sleep() blocks indefinitely.
func Sleep(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("sleep", args)
	a := ProcArg(args, 0, NilValue)
	if a == NilValue {
		time.Sleep(time.Duration(math.MaxInt64)) // approx 290 years
		return nil, nil                          // not reached
	} else {
		n := FloatVal(a)
		d := time.Duration(n * float64(time.Second))
		time.Sleep(d)
		return Return(d)
	}
}

// date() returns the current date in the form "yyyy/mm/dd".
func Date(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("date", args)
	return Return(NewString(time.Now().Format("2006/01/02")))
}

// time() returns the current time of day in the form "hh:mm:ss".
func Time(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("time", args)
	return Return(NewString(time.Now().Format("15:04:05")))
}

// now() returns the current time as an external Go
// http://golang.org/pkg/time#Time[time.Time] value,
// which can then be formatted or otherwise manipulated by calling
// http://golang.org/pkg/time/#Time.Format[tval.Format()]
// or other associated methods.
func Now(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("now", args)
	return Return(time.Now())
}

// duration(x) converts x to an external Go
// http://golang.org/pkg/time#Duration[time.Duration] value.
// If x is a string, it is passed directly to
// http://golang.org/pkg/time#ParseDuration[time.ParseDuration()].
// If x is a number, "s" is appended to interpret it as an interval in seconds.
// If the conversion is unsuccessful, duration() fails.
func Duration(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("duration", args)
	v := ProcArg(args, 0, ZERO)
	s := ""
	if n, ok := v.(*VNumber); ok {
		s = n.String() + "s"
	} else {
		s = ToString(v).ToUTF8()
	}
	d, err := time.ParseDuration(s)
	if err == nil {
		return Return(d)
	} else {
		return Fail()
	}
}

// cputime() returns processor usage in seconds, likely a fractional value.
// The result includes both "user" and "system" time.
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

// exit(i) terminates execution and returns exit status i,
// truncated to integer, to the system.
// A status of 0 signifies normal termination.
func Exit(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("exit", args)
	Shutdown(IntVal(ProcArg(args, 0, ZERO)))
	return Fail() // NOTREACHED
}

// throw(e, x...) raises an exception
// with error value e and zero or more offending values.
// If not caught, the exception terminates execution.
//
// If e is a number or string, a Goaldi exception is created using e.
// Otherwise, the value e is thrown directly, without interpretation.
func Throw(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("throw", args)
	x := ProcArg(args, 0, err_fatal)
	switch v := x.(type) {
	case *VString:
		panic(NewExn(v.String(), args[1:]...))
	case *VNumber:
		panic(NewExn(fmt.Sprintf("Fatal error %v", v), args[1:]...))
	default:
		panic(x)
	}
}

var err_fatal = NewString("Unspecified fatal error")

//  fmisc.go -- standard library setup and miscellaneous functions

package goaldi

import (
	"archive/zip"
	"fmt"
	"os"
	"reflect"
	"syscall"
	"time"
)

//  Miscellaneous library procedures
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
	DefLib(Duration, "duration", "n", "convert value to a Go Duration struct")
	DefLib(CPUtime, "cputime", "", "return total processor time used")
	// Go library functions
	GoLib(os.Getenv, "getenv", "key", "read environment variable")
	GoLib(os.Setenv, "setenv", "key,value", "set environment variable")
	GoLib(os.Environ, "environ", "", "get list of environment variables")
	GoLib(os.Clearenv, "clearenv", "", "delete all environment variables")
	GoLib(os.Hostname, "hostname", "", "get host machine name")
	GoLib(os.Getpid, "getpid", "", "get process ID")
	GoLib(os.Getppid, "getppid", "", "get parent process ID")
	// Heavy-duty package interfaces
	GoLib(zip.OpenReader, "zipreader", "name", "open a Zip file")
}

//  Copy(v) -- return a copy of v (or just v if a simple value).
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

//  Image(v) -- return string image of value v
func Image(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("image", args)
	return Return(NewString(fmt.Sprintf("%#v", ProcArg(args, 0, NilValue))))
}

//  NoResult(e) -- fail immediately -- can be used as a catch handler
func NoResult(env *Env, args ...Value) (Value, *Closure) {
	return Fail()
}

//  NilResult(e) -- return nilresult -- can be used as a catch handler
func NilResult(env *Env, args ...Value) (Value, *Closure) {
	return Return(NilValue)
}

//  ErrResult(e) -- return e -- can be used as a catch handler
func ErrResult(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("errresult", args)
	return Return(ProcArg(args, 0, NilValue))
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

//  Throw(x, v...) -- terminate with error x and offending values v
//  If x is a number or string, a Goaldi exception is created using v.
//  Otherwise, the value x is thrown directly.
func Throw(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("throw", args)
	x := ProcArg(args, 0, err_fatal)
	switch v := x.(type) {
	case *VString:
		panic(&Exception{v.String(), args[1:]})
	case *VNumber:
		panic(&Exception{fmt.Sprintf("Fatal error %v", v), args[1:]})
	default:
		panic(x)
	}
}

var err_fatal = NewString("Unspecified fatal error")

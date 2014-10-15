//  ffile.go -- I/O functions
//
//  #%#% TO DO:
//  implement flags for open(), including new ones:
//		f   fail, don't panic, on open failure
//		m	memory file, implies r/w, buffer in memory (not on disk)
//		s	scratch file, implies r/w, alter name randomly, delete after open
//  add
//	    reads(), readb(), writeb()
//	    implement methods???
//
//  In general:  Files can be passed to Go I/O functions.
//  Goaldi I/O functions panic on error; Go functions return a status code.

package goaldi

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func init() {
	LibProcedure("open", Open)
	LibProcedure("flush", Flush)
	LibProcedure("close", Close)
	LibProcedure("read", Read)
	LibProcedure("write", Write)
	LibProcedure("writes", Writes)
	LibProcedure("print", Print)
	LibProcedure("println", Println)
	LibProcedure("stop", Stop)
	LibGoFunc("printf", fmt.Printf)   // Go library function
	LibGoFunc("fprintf", fmt.Fprintf) // Go library function
}

var noBytes = []byte("")
var spByte = []byte(" ")
var nlByte = []byte("\n")

//  Open(name,flags) -- open a file, or fail
//  #%#%#% currently ignores all flags and opens for sequential buffered reading
func Open(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("open", a)
	name := ProcArg(a, 0, NilValue).(Stringable).ToString().String()
	flags := ProcArg(a, 1, EMPTY).(Stringable).ToString().String()
	f, e := os.Open(name)
	if e != nil {
		return Fail()
	}
	return Return(NewFile(name, flags, f, bufio.NewReader(f), nil))
}

//  Flush(f) -- flush output on a Goaldi file
func Flush(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("flush", a)
	f := ProcArg(a, 0, STDOUT)
	Ock(0, f.(*VFile).Flush())
	return Return(f)
}

//  Close(f) -- close a Goaldi file
func Close(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("close", a)
	f := ProcArg(a, 0, NilValue)
	Ock(0, f.(*VFile).Close())
	return Return(f)
}

//  Read(f) -- return next line from file
func Read(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("read", a)
	r := ProcArg(a, 0, STDIN).(*VFile).Reader
	s, e := r.ReadString('\n')
	if e == io.EOF {
		if s != "" {
			return Return(NewString(s)) // unterminated by \n at EOF
		} else {
			return Fail() // read EOF
		}
	}
	if e != nil {
		panic(e) // other error
	}
	return Return(NewString(s[:len(s)-1])) // trim \n and return
}

//  Write(x,...)
func Write(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("write", a)
	return Wrt(STDOUT, noBytes, nlByte, a)
}

//  Writes(x,...)
func Writes(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("writes", a)
	return Wrt(STDOUT, noBytes, noBytes, a)
}

//  Print(x,...)
func Print(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("print", a)
	return Wrt(STDOUT, spByte, noBytes, a)
}

//  Println(x,...)
func Println(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("println", a)
	return Wrt(STDOUT, spByte, nlByte, a)
}

//  Stop(x,...):
func Stop(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("stop", a)
	Wrt(STDERR, spByte, nlByte, a)
	Shutdown(1) // does not return
	return Fail()
}

//  Wrt(file, between, atEnd, x[]) -- implement write/writes/print/println/stop
func Wrt(f *VFile, between []byte, atEnd []byte, a []Value) (Value, *Closure) {
	if len(a) > 0 { // if there is a first argument
		if altf, ok := a[0].(*VFile); ok { // and it's a file
			f = altf  // use that as the output file
			a = a[1:] // and remove from arglist
		}
	}
	r := NilValue
	for i, v := range a {
		if i > 0 {
			Ock(f.Write(between))
		}
		Ock(fmt.Fprint(f, v))
		r = v
	}
	Ock(f.Write(atEnd))
	Ock(0, f.Flush()) // #%#% seems necessary; should it be?
	return Return(r)
}

//  Ock(n, e) -- output error check: panics if e is non-nil after output call
func Ock(n int, err error) {
	if err != nil {
		panic(err)
	}
}

//  ffile.go -- I/O functions
//
//  In general:  Files can be passed to Go I/O functions.
//  Goaldi I/O functions panic on error; Go functions return a status code.
//
//  #%#% TO DO:
//
//  add random I/O  (same as Icon? including seek/where offsets?)
//  add directory reading?
//
//  implement flags for open(), including new ones:
//		m	memory file, implies r/w, buffer in memory (not on disk)
//		s	scratch file, implies r/w, alter name randomly, delete after open
//
//  add methods???

package goaldi

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func init() {
	LibProcedure("open", Open)
	LibProcedure("flush", Flush)
	LibProcedure("close", Close)
	LibProcedure("read", Read)
	LibProcedure("readb", Readb)
	LibProcedure("write", Write)
	LibProcedure("writes", Writes)
	LibProcedure("writeb", Writeb)
	LibProcedure("print", Print)
	LibProcedure("println", Println)
	LibProcedure("stop", Stop)
	LibGoFunc("printf", fmt.Printf)   // Go library function; don't use %d
	LibGoFunc("fprintf", fmt.Fprintf) // Go library function; don't use %d
}

var noBytes = []byte("")
var spByte = []byte(" ")
var nlByte = []byte("\n")
var dflt_open = NewString("r")

//  Open(name,flags) -- open a file
//	flags:
//		r	open for reading
//		w	open for writing
//      a	open for appending
//		f	fail on error (instead of panicking)  #%#% new
//  #%#% no flag "b": use "rw"
//  #%#% no flag "c": implied by "w"
//  #%#% no flag "t" or "u": done differently (see readb/writeb)
//  #%#% no flag "p": to be considered
func Open(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("open", a)

	name := ProcArg(a, 0, NilValue).(Stringable).ToString().String()
	flags := ProcArg(a, 1, dflt_open).(Stringable).ToString().String()
	fail := false
	read := false
	write := false
	append := false

	// scan flags
	for _, f := range flags {
		switch f {
		case 'r':
			read = true
		case 'w':
			write = true
		case 'a':
			write = true
			append = true
		case 'f':
			fail = true
		default:
			panic(&RunErr{"Unrecognized flag", string([]rune{f})})
		}
	}
	flags = strings.Replace(flags, "f", "", -1) // remove "f" from flags

	// deduce access mode and flags
	amode := 0
	if !write {
		amode = os.O_RDONLY // "r" or unspecified
	} else if read {
		amode = os.O_CREATE | os.O_RDWR // "rw"
	} else {
		amode = os.O_CREATE | os.O_WRONLY // "w" or "a"
	}
	if append {
		amode |= os.O_APPEND
	} else if write {
		amode |= os.O_TRUNC
	}

	// open the file
	f, e := os.OpenFile(name, amode, 0666) // umask modifies 0666
	if e != nil {                          // if error
		if fail {
			return Fail()
		} else {
			panic(e)
		}
	}

	// construct Goaldi file value
	reader := io.Reader(f)
	writer := io.Writer(f)
	if !read {
		reader = nil
	}
	if !write {
		writer = nil
	}
	return Return(NewFile(name, flags, f, reader, writer))
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
//  Fails at EOF when no more data is available.
func Read(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("read", a)
	r := ProcArg(a, 0, STDIN).(*VFile)
	s := r.ReadLine()
	if s == nil {
		return Fail()
	} else {
		return Return(s)
	}
}

//  Readb(f,n) -- read next n binary bytes from file
//  Reads up to n bytes into individual characters without decoding as UTF-8.
//  Useful for reading binary files.
//  Fails at EOF when no more data is available.
func Readb(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("readb", a)
	r := ProcArg(a, 0, STDIN).(*VFile).Reader
	n := int(ProcArg(a, 1, ONE).(Numerable).ToNumber().Val())
	b := make([]byte, n, n)
	n, err := r.Read(b)
	if err == io.EOF {
		return Fail()
	} else if err != nil {
		panic(err)
	} else {
		return Return(BinaryString(b[:n]))
	}
}

//  Writeb(f,s) -- write string s as bytes
//  Writes the low 8 bits of each character of s to file f.
func Writeb(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("writeb", a)
	w := ProcArg(a, 0, STDIN).(*VFile).Writer
	s := ProcArg(a, 1, NilValue).(Stringable).ToString()
	Ock(w.Write(s.ToBinary()))
	return Return(s)
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
	Wrt(STDERR, noBytes, nlByte, a)
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

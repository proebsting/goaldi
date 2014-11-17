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

//  Declare methods
//  Method names begin with an extra F to distinguish from those in vfile.go
//  (whose names are fixed by the need to implement io.ReadWriteCloser).
var FileMethods = map[string]interface{}{
	"type":   (*VFile).Type,
	"copy":   (*VFile).Copy,
	"string": (*VFile).String,
	"image":  (*VFile).GoString,
	"flush":  (*VFile).FFlush,
	"close":  (*VFile).FClose,
	"read":   (*VFile).FRead,
	"readb":  (*VFile).FReadb,
	"writeb": (*VFile).FWriteb,
}

//  VFile.Field implements methods
func (v *VFile) Field(f string) Value {
	return GetMethod(FileMethods, v, f)
}

func init() {
	// Goaldi procedures
	LibProcedure("open", Open)
	LibProcedure("read", Read)
	LibProcedure("write", Write)
	LibProcedure("writes", Writes)
	LibProcedure("print", Print)
	LibProcedure("println", Println)
	LibProcedure("stop", Stop)
	// Go library functions
	LibGoFunc("chdir", os.Chdir)
	LibGoFunc("getwd", os.Getwd)
	LibGoFunc("chmod", os.Chmod)
	LibGoFunc("remove", os.Remove)
	LibGoFunc("mkdir", os.Mkdir)
	LibGoFunc("mkdirall", os.MkdirAll)
	LibGoFunc("rename", os.Rename)
	LibGoFunc("truncate", os.Truncate)
	LibGoFunc("printf", fmt.Printf)   // use %.0f to format as integer
	LibGoFunc("fprintf", fmt.Fprintf) // use %.0f to format as integer
	LibGoFunc("sprintf", fmt.Sprintf) // use %.0f to format as integer
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
func Open(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("open", args)

	name := ProcArg(args, 0, NilValue).(Stringable).ToString().String()
	flags := ProcArg(args, 1, dflt_open).(Stringable).ToString().String()
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
	return Return(NewFile(name, reader, writer, f))
}

//  VFile.FFlush() -- flush output
func (f *VFile) FFlush(args ...Value) (Value, *Closure) {
	defer Traceback("f.flush", args)
	f.Flush()
	return Return(f)
}

//  VFile.FClose(f) -- close a Goaldi file
func (f *VFile) FClose(args ...Value) (Value, *Closure) {
	defer Traceback("f.close", args)
	f.Flush()
	f.Close()
	return Return(f)
}

//  Read(f) -- return next line from file
//  Fails at EOF when no more data is available.
func Read(env *Env, args ...Value) (Value, *Closure) {
	return ProcArg(args, 0, STDIN).(*VFile).FRead(args)
}

//  VFile.FRead() -- read next line from file, failing at EOF.
func (f *VFile) FRead(args ...Value) (Value, *Closure) {
	defer Traceback("f.read", args)
	s := f.ReadLine()
	if s == nil {
		return Fail()
	} else {
		return Return(s)
	}
}

//  VFile.FReadb(n) -- read next n binary bytes from file
//  Reads up to n bytes into individual characters without decoding as UTF-8.
//  Useful for reading binary files.
//  Fails at EOF when no more data is available.
func (f *VFile) FReadb(args ...Value) (Value, *Closure) {
	n := int(ProcArg(args, 0, ONE).(Numerable).ToNumber().Val())
	b := make([]byte, n)
	n, err := f.Reader.Read(b)
	if err == io.EOF {
		return Fail()
	} else if err != nil {
		panic(err)
	} else {
		return Return(BinaryString(b[:n]))
	}
}

//  VFile.FWriteb(s) -- write string s as bytes
//  Writes the low 8 bits of each character of s to file f.
func (f *VFile) FWriteb(args ...Value) (Value, *Closure) {
	defer Traceback("f.writeb", args)
	s := ProcArg(args, 0, NilValue).(Stringable).ToString()
	Ock(f.Writer.Write(s.ToBinary()))
	return Return(f)
}

//  Write(x,...) -- write values and newline to stdout
func Write(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("write", args)
	return Wrt(STDOUT, noBytes, nlByte, args)
}

//  VFile.FWrite(x,...) -- write values and newline to file
func (f *VFile) FWrite(args ...Value) (Value, *Closure) {
	defer Traceback("f.write", args)
	return Wrt(f, noBytes, nlByte, args)
}

//  Writes(x,...) -- write values without newline to stdout
func Writes(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("writes", args)
	return Wrt(STDOUT, noBytes, noBytes, args)
}

//  VFile.FWrites(x,...) -- write values without newline to file
func (f *VFile) FWrites(args ...Value) (Value, *Closure) {
	defer Traceback("f.writes", args)
	return Wrt(f, noBytes, noBytes, args)
}

//  Print(x,...) -- write values with separating whitespace to stdout
func Print(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("print", args)
	return Wrt(STDOUT, spByte, noBytes, args)
}

//  VFile.FPrint(x,...)  -- write values with separating whitespace to file
func (f *VFile) FPrint(args ...Value) (Value, *Closure) {
	defer Traceback("f.print", args)
	return Wrt(f, spByte, noBytes, args)
}

//  Println(x,...)  -- write values with whitespace and newline to stdout
func Println(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("println", args)
	return Wrt(STDOUT, spByte, nlByte, args)
}

//  VFile.FPrintln(x,...) -- write values with whitespace and newline to file
func (f *VFile) FPrintln(args ...Value) (Value, *Closure) {
	defer Traceback("f.println", args)
	return Wrt(f, spByte, nlByte, args)
}

//  Stop(x,...): -- write values to stderr and terminate program
func Stop(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("stop", args)
	Wrt(STDERR, noBytes, nlByte, args)
	Shutdown(1) // does not return
	return Fail()
}

//  Wrt(file, between, atEnd, x[]) -- implement write/writes/print/println/stop
func Wrt(v Value, between []byte, atEnd []byte, args []Value) (Value, *Closure) {
	f := v.(*VFile)
	for i, v := range args {
		if i > 0 {
			Ock(f.Write(between))
		}
		Ock(fmt.Fprint(f, v))
	}
	Ock(f.Write(atEnd))
	Ock(0, f.Flush()) // #%#% seems necessary; should it be?
	return Return(f)
}

//  Ock(n, e) -- output error check: panics if e is non-nil after output call
func Ock(n int, err error) {
	if err != nil {
		panic(err)
	}
}

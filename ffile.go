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
var FileMethods = MethodTable([]*VProcedure{
	DefMeth((*VFile).FFlush, "flush", "", "flush file"),
	DefMeth((*VFile).FClose, "close", "", "close file"),
	DefMeth((*VFile).FRead, "read", "", "read one line"),
	DefMeth((*VFile).FReadb, "readb", "size", "read binary bytes"),
	DefMeth((*VFile).FWriteb, "writeb", "s", "write binary bytes"),
	DefMeth((*VFile).FWrite, "write", "x[]", "write values and newline"),
	DefMeth((*VFile).FWrites, "writes", "x[]", "write values"),
	DefMeth((*VFile).FPrint, "print", "x[]", "write values with spacing"),
	DefMeth((*VFile).FPrintln, "println", "x[]", "write line of values"),
})

//  Declare procedures
func init() {
	// Goaldi procedures
	DefLib(Open, "open", "name,flags", "open a file")
	DefLib(Read, "read", "f", "read one line from a file")
	DefLib(Write, "write", "x[]", "write values and newline")
	DefLib(Writes, "writes", "x[]", "write values")
	DefLib(Print, "print", "x[]", "write values with spacing")
	DefLib(Println, "println", "x[]", "write line of values")
	DefLib(Stop, "stop", "x[]", "write values and abort program")
	// Go library functions
	GoLib(os.Chdir, "chdir", "dir", "change working directory")
	GoLib(os.Getwd, "getwd", "", "get working directory")
	GoLib(os.Chmod, "chmod", "name,mode", "change file mode")
	GoLib(os.Remove, "remove", "name", "delete file")
	GoLib(os.Mkdir, "mkdir", "name,perm", "create directory")
	GoLib(os.MkdirAll, "mkdirall", "path,perm", "create directory tree")
	GoLib(os.Rename, "rename", "old,new", "change file name")
	GoLib(os.Truncate, "truncate", "name,size", "change file size")
	GoLib(fmt.Printf, "printf", "fmt,x[]", "write with formatting")
	GoLib(fmt.Fprintf, "fprintf", "f,fmt,x[]", "write to file with formatting")
	GoLib(fmt.Sprintf, "sprintf", "fmt,x[]", "make string by formatting values")
}

var noBytes = []byte("")
var spByte = []byte(" ")
var nlByte = []byte("\n")
var dflt_open = NewString("r")

//  open(name,flags) opens a file and returns a file value.
//
//  Each character of the optional flags argument selects an option:
//		"r"	open for reading
//		"w"	open for writing
//		"a"	open for appending
//		"f"	fail on error (instead of panicking)
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
			panic(NewExn("Unrecognized flag", string([]rune{f})))
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

//  VFile.FClose() -- close a Goaldi file
func (f *VFile) FClose(args ...Value) (Value, *Closure) {
	defer Traceback("f.close", args)
	f.Flush()
	f.Close()
	return Return(f)
}

//  read(f) consumes and returns next line of text from file f.
//  The trailing linefeed or CRLF is removed from the returned value.
//  read() fails at EOF when no more data is available.
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
	defer Traceback("f.readb", args)
	if f.Reader == nil {
		panic(NewExn("Not open for reading", f))
	}
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
	if f.Writer == nil {
		panic(NewExn("Not open for writing", f))
	}
	s := ProcArg(args, 0, NilValue).(Stringable).ToString()
	Ock(f.Writer.Write(s.ToBinary()))
	return Return(f)
}

//  write(x,...) writes its arguments to %stdout followed by a newline.
func Write(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("write", args)
	return Wrt(STDOUT, noBytes, nlByte, args)
}

//  VFile.FWrite(x,...) -- write values and newline to file
func (f *VFile) FWrite(args ...Value) (Value, *Closure) {
	defer Traceback("f.write", args)
	return Wrt(f, noBytes, nlByte, args)
}

//  writes(x,...) write its arguments to %stdout with no following newline.
func Writes(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("writes", args)
	return Wrt(STDOUT, noBytes, noBytes, args)
}

//  VFile.FWrites(x,...) -- write values without newline to file
func (f *VFile) FWrites(args ...Value) (Value, *Closure) {
	defer Traceback("f.writes", args)
	return Wrt(f, noBytes, noBytes, args)
}

//  print(x,...) writes its arguments to %stdout, separated by spaces.
func Print(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("print", args)
	return Wrt(STDOUT, spByte, noBytes, args)
}

//  VFile.FPrint(x,...)  -- write values with separating whitespace to file
func (f *VFile) FPrint(args ...Value) (Value, *Closure) {
	defer Traceback("f.print", args)
	return Wrt(f, spByte, noBytes, args)
}

//  println(x,...) writes its arguments to %stdout,
//  separated by spaces and terminated by a newline character.
func Println(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("println", args)
	return Wrt(STDOUT, spByte, nlByte, args)
}

//  VFile.FPrintln(x,...) -- write values with whitespace and newline to file
func (f *VFile) FPrintln(args ...Value) (Value, *Closure) {
	defer Traceback("f.println", args)
	return Wrt(f, spByte, nlByte, args)
}

//  stop(x,...) writes its arguments to %stderr and terminates execution
//  with an exit code of 1 (indicating an error).
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

//  ffile.go -- I/O functions
//
//  In general:  Files can be passed to Go I/O functions.
//  Goaldi I/O functions panic on error; Go functions return a status code.

package runtime

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Declare methods
// Method names begin with an extra F to distinguish from those in vfile.go
// (whose names are fixed by the need to implement io.ReadWriteCloser).
var FileMethods = MethodTable([]*VProcedure{
	DefMeth((*VFile).FUnbuffer, "unbuffer", "", "stop file buffering"),
	DefMeth((*VFile).FFlush, "flush", "", "flush file"),
	DefMeth((*VFile).FClose, "close", "", "close file"),
	DefMeth((*VFile).FGet, "get", "", "read one line"),
	DefMeth((*VFile).FRead, "read", "", "read one line"),
	DefMeth((*VFile).FReadb, "readb", "size", "read binary bytes"),
	DefMeth((*VFile).FWriteb, "writeb", "s", "write binary bytes"),
	DefMeth((*VFile).FPut, "put", "x[]", "write values as lines"),
	DefMeth((*VFile).FWrite, "write", "x[]", "write values and newline"),
	DefMeth((*VFile).FWrites, "writes", "x[]", "write values"),
	DefMeth((*VFile).FPrint, "print", "x[]", "write values with spacing"),
	DefMeth((*VFile).FPrintln, "println", "x[]", "write line of values"),
	DefMeth((*VFile).FSeek, "seek", "n", "set file position"),
	DefMeth((*VFile).FWhere, "where", "", "report current file position"),
})

// Declare procedures
func init() {
	// Goaldi procedures
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

// file(name,flags) opens a file and returns a file value.
//
// Each character of the optional flags argument selects an option:
//
//	"r"   open for reading
//	"w"   open for writing
//	"a"   open for appending
//	"c"   create and open for writing
//	"n"   no buffering
//	"f"   fail on error (instead of panicking)
//
// If none of "w", "a", or "c" are specified, then "r" is implied.
// "w" implies "c" unless "r" is also specified.
// Buffering is used if "n" is absent and the file is opened
// exclusively for reading or writing but not both.
//
// In the absence of "f", any error throws an exception.
func File(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("file", args)

	name := ToString(ProcArg(args, 0, NilValue)).ToUTF8()
	flags := ToString(ProcArg(args, 1, EMPTY)).ToUTF8()
	fail := false
	read := false
	write := false
	append := false
	create := false
	buffer := true

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
		case 'c':
			write = true
			create = true
		case 'n':
			buffer = false
		case 'f':
			fail = true
		default:
			panic(NewExn("Unrecognized flag", string([]rune{f})))
		}
	}

	// deduce access mode
	var amode int
	if !write {
		read = true
		amode = os.O_RDONLY // "r" or unspecified
	} else if read {
		amode = os.O_RDWR // "rw"
	} else {
		amode = os.O_CREATE | os.O_WRONLY // "w" or "a" or "c"
	}
	if create {
		amode |= os.O_CREATE
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
	if read && !write { // if read only
		writer = nil
		if buffer {
			reader = bufio.NewReader(reader)
		}
	} else if write && !read { // if write only
		reader = nil
		if buffer {
			writer = bufio.NewWriter(writer)
		}
	}
	return Return(NewFile(name, f, reader, writer, f))
}

// f.unbuffer() removes any buffering from file f.
// Any buffered output is flushed; any buffered input is lost.
func (f *VFile) FUnbuffer(args ...Value) (Value, *Closure) {
	defer Traceback("f.unbuffer", args)
	if _, ok := f.Reader.(*bufio.Reader); ok {
		f.Reader = f.Original.(io.Reader)
	}
	if _, ok := f.Writer.(*bufio.Writer); ok {
		f.Flush()
		f.Writer = f.Original.(io.Writer)
	}
	return Return(f)
}

// f.flush() flushes output on file f.
func (f *VFile) FFlush(args ...Value) (Value, *Closure) {
	defer Traceback("f.flush", args)
	f.Flush()
	return Return(f)
}

// f.close() closes file f.
func (f *VFile) FClose(args ...Value) (Value, *Closure) {
	defer Traceback("f.close", args)
	f.Flush()
	f.Close()
	return Return(f)
}

// f.seek(n) sets the position for the next read or write on file f.
// File positions are measured in bytes, not characters, counting the
// first byte as 1.  A value of 0 seeks to end of file, and a negative
// value is an offset from the end.
func (f *VFile) FSeek(args ...Value) (Value, *Closure) {
	defer Traceback("f.seek", args)
	posn := int64(FloatVal(ProcArg(args, 0, ONE)))
	whence := os.SEEK_SET
	if posn > 0 {
		posn = posn - 1
	} else {
		whence = os.SEEK_END
	}
	_, err := f.Seek(posn, whence)
	if err != nil {
		panic(err)
	}
	return Return(f)
}

// f.where() reports the current position of file f.
// File positions are measured in bytes, counting the first byte as 1.
func (f *VFile) FWhere(args ...Value) (Value, *Closure) {
	defer Traceback("f.where", args)
	offset, err := f.Seek(0, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	return Return(NewNumber(float64(offset + 1)))
}

// read(f) consumes and returns next line of text from file f.
// The trailing linefeed or CRLF is removed from the returned value.
// read() fails at EOF when no more data is available.
func Read(env *Env, args ...Value) (Value, *Closure) {
	if len(args) > 0 {
		return ProcArg(args, 0, NilValue).(*VFile).FRead(args)
	} else {
		return env.Lookup("stdin", true).(*VFile).FRead(args)
	}
}

// f.get() consumes and returns next line of text from file f.
// The trailing linefeed or CRLF is removed from the returned value.
// f.get() fails at EOF when no more data is available.
func (f *VFile) FGet(args ...Value) (Value, *Closure) {
	defer Traceback("f.get", args)
	s := f.ReadLine()
	if s == nil {
		return Fail()
	} else {
		return Return(s)
	}
}

// f.read() consumes and returns next line of text from file f.
// The trailing linefeed or CRLF is removed from the returned value.
// f.read() fails at EOF when no more data is available.
func (f *VFile) FRead(args ...Value) (Value, *Closure) {
	defer Traceback("f.read", args)
	s := f.ReadLine()
	if s == nil {
		return Fail()
	} else {
		return Return(s)
	}
}

// f.readb(n) reads up to n bytes into individual characters
// without attempting any UTF-8 decoding.
// This is useful for reading binary files.
// f.readb() fails at EOF when no more data is available.
func (f *VFile) FReadb(args ...Value) (Value, *Closure) {
	defer Traceback("f.readb", args)
	if f.Reader == nil {
		panic(NewExn("Not open for reading", f))
	}
	n := IntVal(ProcArg(args, 0, ONE))
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

// f.writeb(s) writes the string s to file f without any UTF-8 encoding.
// Instead, the low 8 bits of each character are written as a single byte,
// ignoring all other bits.
// This is useful for writing binary files.
func (f *VFile) FWriteb(args ...Value) (Value, *Closure) {
	defer Traceback("f.writeb", args)
	if f.Writer == nil {
		panic(NewExn("Not open for writing", f))
	}
	s := ToString(ProcArg(args, 0, NilValue))
	Ock(f.Writer.Write(s.ToBinary()))
	return Return(f)
}

// write(x,...) writes its arguments to %stdout followed by a newline.
func Write(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("write", args)
	return Wrt(env.Lookup("stdout", true), noBytes, nlByte, args)
}

// f.put(x,...) writes its arguments to file f, each followed by a newline.
// This treats a file as as a container of text values separated by newlines,
// which is consistent with the interpretation used by f.get().
func (f *VFile) FPut(args ...Value) (Value, *Closure) {
	defer Traceback("f.put", args)
	return Wrt(f, nlByte, nlByte, args)
}

// f.write(x,...) writes its arguments to file f followed by a single newline.
func (f *VFile) FWrite(args ...Value) (Value, *Closure) {
	defer Traceback("f.write", args)
	return Wrt(f, noBytes, nlByte, args)
}

// writes(x,...) write its arguments to %stdout with no following newline.
func Writes(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("writes", args)
	return Wrt(env.Lookup("stdout", true), noBytes, noBytes, args)
}

// f.writes(x,...) write its arguments to file f with no following newline.
func (f *VFile) FWrites(args ...Value) (Value, *Closure) {
	defer Traceback("f.writes", args)
	return Wrt(f, noBytes, noBytes, args)
}

// print(x,...) writes its arguments to %stdout, separated by spaces.
func Print(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("print", args)
	return Wrt(env.Lookup("stdout", true), spByte, noBytes, args)
}

// f.print(x,...) writes its arguments to file f, separated by spaces.
func (f *VFile) FPrint(args ...Value) (Value, *Closure) {
	defer Traceback("f.print", args)
	return Wrt(f, spByte, noBytes, args)
}

// println(x,...) writes its arguments to %stdout,
// separated by spaces and terminated by a newline character.
func Println(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("println", args)
	return Wrt(env.Lookup("stdout", true), spByte, nlByte, args)
}

// f.println(x,...) writes its arguments to file f,
// separated by spaces and terminated by a newline character.
func (f *VFile) FPrintln(args ...Value) (Value, *Closure) {
	defer Traceback("f.println", args)
	return Wrt(f, spByte, nlByte, args)
}

// stop(x,...) writes its arguments to %stderr and terminates execution
// with an exit code of 1 (indicating an error).
func Stop(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("stop", args)
	Wrt(env.Lookup("stderr", true), noBytes, nlByte, args)
	Shutdown(1) // does not return
	return Fail()
}

// Wrt(file, between, atEnd, x[]) -- implement write/writes/print/println/stop
func Wrt(v Value, between []byte, atEnd []byte, args []Value) (Value, *Closure) {
	f := v.(*VFile)
	for i, v := range args {
		if i > 0 {
			Ock(f.Write(between))
		}
		Ock(fmt.Fprint(f, v))
	}
	Ock(f.Write(atEnd))
	Ock(0, f.Flush()) // flush at end of operation -- file may be buffered
	return Return(f)
}

// Ock(n, e) -- output error check: panics if e is non-nil after output call
func Ock(n int, err error) {
	if err != nil {
		panic(err)
	}
}

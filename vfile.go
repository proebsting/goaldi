//  vfile.go -- implementation of a file type
//
//  A Goaldi file is produced by open().
//  It implements io.ReadWriteCloser so it can be passed to Go funcs.
//
//  NOTE:  Read and Write are Go methods.  read and write are Goaldi methods.
//  #%#%#% once Goaldi has methods, that is.

//  #%#% Do we need to register all files in order to flush them on exit?

package goaldi

import (
	"bufio"
	"io"
	"os"
)

// confirm implementation of promised interfaces
var _ ICore = &VFile{}
var _ io.ReadWriteCloser = &VFile{}

//  standard files
var (
	STDIN  = NewFile("%stdin", "", bufio.NewReader(os.Stdin), os.Stdin)
	STDOUT = NewFile("%stdout", "", bufio.NewWriter(os.Stdout), os.Stdout)
	STDERR = NewFile("%stderr", "", io.Writer(os.Stderr), os.Stderr)
)

type VFile struct {
	Name  string      // name when opened
	Flags string      // attributes when opened
	Actor interface{} // an io.Reader or io.Writer (or both)
	File  *os.File    // underlying file (needed for close etc.)
}

//  NewFile(name, flags, actor) -- construct new Goaldi file
//  flags are from "open", EXCLUDING "r" and "w"
//  actor is an io.Reader or io.Writer
func NewFile(name string, flags string, actor interface{}, f *os.File) *VFile {
	a := ""
	if _, ok := actor.(io.Reader); ok {
		a = a + "r"
	}
	if _, ok := actor.(io.Writer); ok {
		a = a + "w"
	}
	if a == "" {
		panic(&RunErr{"Neither reader nor writer", actor})
	}
	return &VFile{name, a + flags, actor, f}
}

//  VFile.String -- conversion to Go string returns "file(name)"
func (v *VFile) String() string {
	return "file(" + v.Name + ")"
}

//  VFile.GoString -- image returns "file(name,flags)"
func (v *VFile) GoString() string {
	return "file(" + v.Name + "," + v.Flags + ")"
}

//  VFile.Type returns "file"
func (v *VFile) Type() Value {
	return type_file
}

var type_file = NewString("file")

//  VFile.Export returns the underlying io.Reader / io.Writer
func (v *VFile) Export() interface{} {
	return v.Actor
}

//  VFile.Read() calls io.Read().
func (v *VFile) Read(p []byte) (int, error) {
	if r, ok := v.Actor.(io.Reader); ok {
		return r.Read(p)
	} else {
		panic(&RunErr{"Not open for reading", v})
	}
}

//  VFile.Write() calls io.Write().
func (v *VFile) Write(p []byte) (int, error) {
	if w, ok := v.Actor.(io.Writer); ok {
		return w.Write(p)
	} else {
		panic(&RunErr{"Not open for writing", v})
	}
}

//  VFile.Close() marks the file as closed and calls io.Close().
func (v *VFile) Close() error {
	a := v.Actor
	if b, ok := a.(bufio.Writer); ok {
		b.Flush()
	}
	v.Actor = nil
	v.Flags = "-"
	return v.File.Close()
}

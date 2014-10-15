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
	STDIN  = NewFile("%stdin", "r", os.Stdin, bufio.NewReader(os.Stdin), nil)
	STDOUT = NewFile("%stdout", "w", os.Stdout, nil, bufio.NewWriter(os.Stdout))
	STDERR = NewFile("%stderr", "w", os.Stderr, nil, io.Writer(os.Stderr))
)

type VFile struct {
	Name   string        // name when opened
	Flags  string        // attributes when opened
	File   *os.File      // underlying file (needed for close etc.)
	Reader *bufio.Reader // underlying reader, if open for read
	Writer io.Writer     // underlying writer, if open for write
}

//  NewFile(name, flags, file, reader, writer) -- construct new Goaldi file
func NewFile(name string, flags string, file *os.File,
	reader *bufio.Reader, writer io.Writer) *VFile {
	return &VFile{name, flags, file, reader, writer}
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

//  VFile.Export returns itself, which implements the ReadWriteCloser interface
func (v *VFile) Export() interface{} {
	return v
}

//  VFile.Read() calls io.Read().
func (v *VFile) Read(p []byte) (int, error) {
	if v.Reader != nil {
		return v.Reader.Read(p)
	} else {
		panic(&RunErr{"Not open for reading", v})
	}
}

//  VFile.Write() calls io.Write().
func (v *VFile) Write(p []byte) (int, error) {
	if v.Writer != nil {
		return v.Writer.Write(p)
	} else {
		panic(&RunErr{"Not open for writing", v})
	}
}

//  VFile.Flush() flushes the output stream if possible.
func (v *VFile) Flush() error {
	if b, ok := v.Writer.(*bufio.Writer); ok {
		return b.Flush()
	} else {
		return nil
	}
}

//  VFile.Close() marks the file as closed and calls io.Close().
func (v *VFile) Close() error {
	if v.File == nil {
		panic(&RunErr{"File not open", v})
	}
	err := v.Flush()
	if err != nil {
		return err
	}
	f := v.File
	v.Flags = "-"
	v.File = nil
	v.Reader = nil
	v.Writer = nil
	return f.Close()
}

//  vfile.go -- implementation of a file type
//
//  A Goaldi file is produced by open().
//  It implements io.ReadWriteCloser so that it can be passed to Go funcs.
//
//  NOTE:  Read and Write in here are Go methods, not Goaldi procedures.

//  #%#% Do we need to register all files in order to flush them on exit?

package goaldi

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

// confirm implementation of promised interfaces
var _ ICore = &VFile{}
var _ io.ReadWriteCloser = &VFile{}

//  standard files, referenced (and changeable) by keyword / dynamic variables
var (
	STDIN  Value = NewFile("%stdin", bufio.NewReader(os.Stdin), nil, os.Stdin)
	STDOUT Value = NewFile("%stdout", nil, bufio.NewWriter(os.Stdout), os.Stdout)
	STDERR Value = NewFile("%stderr", nil, io.Writer(os.Stderr), os.Stderr)
)

type VFile struct {
	Name   string        // name when opened
	Reader *bufio.Reader // reader, if open for read
	Writer io.Writer     // writer, if open for write
	Closer io.Closer     // closer
}

//  NewFile(name, reader, writer, closer) -- construct new Goaldi file
func NewFile(name string,
	reader io.Reader, writer io.Writer, closer io.Closer) *VFile {
	// if no closer, add one, because nil means file is already closed
	if closer == nil {
		closer = ioutil.NopCloser(reader)
	}
	// if file is not for reading, nothing much to do
	if reader == nil {
		return &VFile{name, nil, writer, closer}
	}
	// if reader is not bufio.Reader, wrap it (needed for readline)
	if _, ok := reader.(*bufio.Reader); !ok {
		reader = bufio.NewReader(reader)
	}
	// create file with correctly typed buffered reader
	return &VFile{name, reader.(*bufio.Reader), writer, closer}
}

//  VFile.String -- conversion to Go string returns "F:name"
func (v *VFile) String() string {
	return "F:" + v.Name
}

//  VFile.GoString -- image returns "file(name,[r][w])"
func (v *VFile) GoString() string {
	s := "file(" + v.Name + ","
	if v.Reader != nil {
		s = s + "r"
	}
	if v.Writer != nil {
		s = s + "w"
	}
	return s + ")"
}

//  VFile.Rank returns rFile
func (v *VFile) Rank() int {
	return rFile
}

//  VFile.Type returns "file"
func (v *VFile) Type() Value {
	return type_file
}

var type_file = NewString("file")

//  VFile.Copy returns itself
func (v *VFile) Copy() Value {
	return v
}

//  VFile.Import returns itself
func (v *VFile) Import() Value {
	return v
}

//  VFile.Export returns itself, which implements the ReadWriteCloser interface
func (v *VFile) Export() interface{} {
	return v
}

//  VFile.Dispense() implements the !f operator
func (f *VFile) Dispense(unused IVariable) (Value, *Closure) {
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		s := f.ReadLine()
		if s != nil {
			return s, c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

//  VFile.Take() implements the @f operator
func (f *VFile) Take() Value {
	s := f.ReadLine()
	if s != nil {
		return s
	} else {
		return nil
	}
}

//  VFile.Read() calls io.Read().  This implements the Go io.Reader interface.
func (v *VFile) Read(p []byte) (int, error) {
	if v.Reader != nil {
		return v.Reader.Read(p)
	} else {
		panic(&RunErr{"Not open for reading", v})
	}
}

//  VFile.ReadLine() returns the next line from this file, or nil at EOF.
func (v *VFile) ReadLine() *VString {
	if v.Reader == nil {
		panic(&RunErr{"Not open for reading", v})
	}
	s, e := v.Reader.ReadString('\n')
	if e == nil {
		n := len(s) - 1              // position of trailing \n
		if n > 0 && s[n-1] == '\r' { // if preceded by \r
			return NewString(s[:n-1]) // trim CRLF and return
		} else {
			return NewString(s[:n]) // trim \n and return
		}
	} else if e != io.EOF {
		panic(e) // actual error
	} else if s != "" {
		return NewString(s) // unterminated by \n at EOF
	} else {
		return nil // hit EOF with no more data
	}
}

//  VFile.Write() calls io.Write().  This implements the Go io.Writer interface.
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

//  VFile.Close() closes a file.  This implements the Go io.Closer interface.
//  It marks the file as closed and calls io.Close() on the underlying Closer.
func (v *VFile) Close() error {
	if v.Closer == nil {
		panic(&RunErr{"File not open", v})
	}
	err := v.Flush()
	if err != nil {
		return err
	}
	c := v.Closer
	v.Reader = nil
	v.Writer = nil
	v.Closer = nil
	return c.Close()
}

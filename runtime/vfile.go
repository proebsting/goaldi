//  vfile.go -- VFile, the Goaldi type "file"
//
//  A Goaldi file is produced by file(), which replaces open().
//  A file implements io.ReadWriteCloser so that it can be passed to Go funcs.

package runtime

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

const rFile = 30                    // declare sort ranking
var _ ICore = &VFile{}              // validate implementation
var _ io.ReadWriteCloser = &VFile{} // ensure promise is kept

// standard files, referenced (and changeable) by keyword / dynamic variables
var (
	STDIN Value = NewFile(
		"%stdin", os.Stdin, bufio.NewReader(os.Stdin), nil, os.Stdin)
	STDOUT Value = NewFile(
		"%stdout", os.Stdout, nil, bufio.NewWriter(os.Stdout), os.Stdout)
	STDERR Value = NewFile(
		"%stderr", os.Stderr, nil, io.Writer(os.Stderr), os.Stderr)
)

type VFile struct {
	Name     string      // name when opened
	Original interface{} // underlying object (os.File? other Reader or Writer?)
	Reader   io.Reader   // reader, if open for read
	Writer   io.Writer   // writer, if open for write
	Closer   io.Closer   // closer; the underlying file, if buffered
}

var FileType = NewType("file", "f", rFile, File, FileMethods,
	"file", "name,flags", "open a file")

// NewFile(name, file, reader, writer, closer) -- construct new Goaldi file
func NewFile(name string, file interface{},
	reader io.Reader, writer io.Writer, closer io.Closer) *VFile {
	if closer == nil { // need a closer; nil means already closed
		closer = ioutil.NopCloser(reader)
	}
	return &VFile{name, file, reader, writer, closer}
}

// VFile.String -- conversion to Go string returns "f:name"
func (v *VFile) String() string {
	return "f:" + v.Name
}

// VFile.GoString -- image returns "file(name,[r][w][n])"
func (v *VFile) GoString() string {
	s := "file(" + v.Name + ","
	if v.Reader != nil {
		s = s + "r"
	}
	if v.Writer != nil {
		s = s + "w"
	}
	if v.Closer != nil && !v.IsBuffered() {
		s = s + "n"
	}
	return s + ")"
}

// VFile.Type returns the file type
func (v *VFile) Type() IRank {
	return FileType
}

// VFile.Copy returns itself
func (v *VFile) Copy() Value {
	return v
}

// VFile.Before compares two files for sorting
func (a *VFile) Before(b Value, i int) bool {
	return a.Name < b.(*VFile).Name
}

// VFile.Import returns itself
func (v *VFile) Import() Value {
	return v
}

// VFile.Export exports a Goaldi file
// If the file is buffered, it returns the Goaldi file,
// which implements the ReadWriteCloser interface.
// If not, it returns the underlying file.
func (v *VFile) Export() interface{} {
	if v.IsBuffered() {
		return v // Goaldi file is buffered; return that
	} else {
		return v.Original // no; return underlying file
	}
}

// VFile.IsBuffered returns true if f is buffered.
func (v *VFile) IsBuffered() bool {
	if _, ok := v.Reader.(*bufio.Reader); ok { // if input is buffered
		return true
	}
	if _, ok := v.Writer.(*bufio.Writer); ok { // if output is buffered
		return true
	}
	return false // neither is buffered
}

// VFile.Read() calls io.Read().  This implements the Go io.Reader interface.
func (v *VFile) Read(p []byte) (int, error) {
	if v.Reader != nil {
		return v.Reader.Read(p)
	} else {
		panic(NewExn("Not open for reading", v))
	}
}

// VFile.ReadLine() returns the next line from this file, or nil at EOF.
func (v *VFile) ReadLine() *VString {
	if v.Reader == nil {
		panic(NewExn("Not open for reading", v))
	}
	var s string
	var e error
	if r, ok := v.Reader.(*bufio.Reader); ok {
		// use library func to read a line
		s, e = r.ReadString('\n')
	} else {
		// not buffered; read a char at a time up through newline
		var b bytes.Buffer
		p := make([]byte, 1)
		for p[0] != '\n' {
			n, e := v.Reader.Read(p)
			if n > 0 {
				b.Write(p)
			}
			if e != nil {
				break
			}
		}
		s = b.String()
	}
	// interpret and return results of reading
	if e == nil {
		n := len(s)
		if n == 0 {
			return nil // EOF
		}
		if s[n-1] == '\n' { // if ends with \n, remove it
			n--
			if n > 0 && s[n-1] == '\r' { // if preceded by \r, remove that
				n--
			}
		}
		return NewString(s[:n])
	} else if e != io.EOF {
		panic(e) // actual error
	} else if s != "" {
		return NewString(s) // unterminated by \n at EOF
	} else {
		return nil // hit EOF with no more data
	}
}

// VFile.Write() calls io.Write().  This implements the Go io.Writer interface.
func (v *VFile) Write(p []byte) (int, error) {
	if v.Writer != nil {
		return v.Writer.Write(p)
	} else {
		panic(NewExn("Not open for writing", v))
	}
}

// VFile.Flush() flushes the output stream if possible.
func (v *VFile) Flush() error {
	if b, ok := v.Writer.(*bufio.Writer); ok {
		return b.Flush()
	} else {
		return nil
	}
}

// VFile.Seek(offset, whence) implements the Go io.Seeker interface.
func (v *VFile) Seek(offset int64, whence int) (int64, error) {
	var f io.Seeker
	var ok bool
	// Either the Reader or Writer must be seekable.
	// (If both are non-nil then they are identical and either will do.)
	f, ok = v.Reader.(io.Seeker)
	if !ok {
		f, ok = v.Writer.(io.Seeker)
	}
	if !ok {
		panic(NewExn("Not seekable", v))
	}
	return f.Seek(offset, whence)
}

// VFile.Close() closes a file.  This implements the Go io.Closer interface.
// It marks the file as closed and calls io.Close() on the underlying Closer.
func (v *VFile) Close() error {
	if v.Closer == nil {
		panic(NewExn("File not open", v))
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

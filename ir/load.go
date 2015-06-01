//  load.go -- read intermediate representation from JSON file

package ir

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"encoding/json"
	"fmt"
	g "goaldi/runtime"
	"io"
	"reflect"
	"runtime"
	"unicode"
	"unicode/utf8"
)

var fileNumber = 1

//  Load(reader) -- read a JSON-encoded IR file.
//
//  Each section of the input file is a JSON list value corresponding to a
//  single source file.  (It is typical, then, to find just one section.)
//  A list of sections -- a list of lists of IR structs -- is returned.
//
//  A per-section distinguishing integer is prepended to each procedure name
//  that begins with "$".  No other changes are made during input.
//
func Load(rdr io.Reader) (comments []string, ircode [][]interface{}) {

	//  collect initial comment lines (e.g. #!/usr/bin/env goaldi ...)
	buffi := bufio.NewReader(rdr)
	comments = []string{}
	for {
		b, e := buffi.Peek(1)
		if e != nil || b[0] != '#' {
			break
		}
		cmt, _ := buffi.ReadBytes('\n')
		comments = append(comments, string(cmt[:len(cmt)-1]))
	}

	//  check for bzip2-encoded file
	gcode := io.Reader(buffi)
	bzheader := []byte("BZh91AY&SY")
	b, _ := buffi.Peek(10)
	if len(b) == 10 && bytes.Compare(b, bzheader) == 0 {
		gcode = bzip2.NewReader(buffi) // use decompressing reader
		if runtime.NumCPU() > 1 {
			// with multiple CPUs, read in parallel
			// (buffer size tuned empirically on Perigrine.cs 22-May-2015)
			gcode = ParallelReader(gcode, 16*1024)
			defer func() { gcode.(io.Closer).Close() }()
		}
	}

	//  load the JSON-encoded program
	jd := json.NewDecoder(gcode)
	ircode = make([][]interface{}, 0)

	for {
		var jtree []interface{}
		err := jd.Decode(&jtree)
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				panic(err)
			}
		}
		jtree = jstructs(jtree).([]interface{})
		ircode = append(ircode, jtree)
		fileNumber++
	}
	return comments, ircode
}

//  jstructs -- replace maps by IR structs in Json tree
func jstructs(jtree interface{}) interface{} {
	switch x := jtree.(type) {
	case []interface{}:
		for i, v := range x {
			x[i] = jstructs(v)
		}
		return x
	case map[string]interface{}:
		for k, v := range x {
			x[k] = jstructs(v)
		}
		return structFor(x)
	default:
		return jtree
	}
}

//  structFor -- return IR struct equivalent to map
func structFor(m map[string]interface{}) interface{} {
	tag := m["tag"].(string)
	if tag == "" {
		panic(g.Malfunction(fmt.Sprintf("No tag in %v", m)))
	}
	rtype := irtable[tag]
	if rtype == nil {
		panic(g.Malfunction(fmt.Sprintf("Unrecognized IR tag %s", tag)))
	}
	resultp := reflect.New(rtype)
	result := resultp.Elem()
	for key, val := range m {
		key = Capitalize(key)
		f := result.FieldByName(key)
		setField(f, key, val)
	}
	return result.Interface()
}

//  setField -- set field in struct
func setField(f reflect.Value, key string, val interface{}) {

	defer func() {
		if x := recover(); x != nil {
			fmt.Println("SETFIELD PANIC: ", f, f.Interface(), key, val)
			panic(x)
		}
	}()

	if key == "Tag" || val == nil {
		return // nothing to do
	}
	if !f.CanSet() {
		panic(g.Malfunction("Cannot set key " + key))
	}

	t := f.Type()
	if t.Kind() != reflect.Slice || t.Elem().Kind() == reflect.Interface {
		// set a simple value
		if t.Kind() == reflect.Int { // ints are read as floats
			val = int(val.(float64))
		} else if t.Kind() == reflect.String {
			s := val.(string)
			if len(s) > 0 && s[0] == '$' {
				// prefix file number to Name/Parent/Fn beginning with '$'
				switch key {
				case "Name", "Fn", "Parent":
					val = fmt.Sprintf("%d%s", fileNumber, val)
				}
			}
		}
		v := reflect.ValueOf(val)
		if f.Kind() == reflect.Ptr && v.Kind() != reflect.Ptr {
			// we have a value but need a pointer;
			// copy the value to get an assignable pointer
			p := reflect.New(v.Type())
			p.Elem().Set(v)
			v = p
		}
		f.Set(v)
		return
	}

	// we have to make a typed slice and copy in the elements
	resultp := reflect.New(t)
	result := resultp.Elem()
	for _, v := range val.([]interface{}) {
		if f, ok := v.(float64); ok { // if float64
			v = int(f) // convert to int
		}
		result = reflect.Append(result, reflect.ValueOf(v))
	}
	f.Set(result)
}

//  Capitalize(s) -- convert first character of string to upper case
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

//  DeCapit(s) -- convert first character of string to lower case
func DeCapit(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

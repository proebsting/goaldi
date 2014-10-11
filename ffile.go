//  ffile.go -- I/O functions

package goaldi

import (
	"fmt"
	"os"
)

func init() {
	LibProcedure("write", write)
	LibProcedure("writes", writes)
	LibProcedure("print", print)
	LibProcedure("println", println)
	LibGoFunc("printf", fmt.Printf) // Go library function
}

var noBytes = []byte("")
var spByte = []byte(" ")
var nlByte = []byte("\n")

func write(env *Env, a ...Value) (Value, *Closure) {
	return wrt(noBytes, nlByte, a)
}

func writes(env *Env, a ...Value) (Value, *Closure) {
	return wrt(noBytes, noBytes, a)
}

func print(env *Env, a ...Value) (Value, *Closure) {
	return wrt(spByte, noBytes, a)
}

func println(env *Env, a ...Value) (Value, *Closure) {
	return wrt(spByte, nlByte, a)
}

func wrt(between []byte, atEnd []byte, a []Value) (Value, *Closure) {
	f := os.Stdout //#%#% should eventually use a buffered version?
	//#%#% if a[0] is a file, switch files, adjust "between" handling
	r := NilValue
	for i, v := range a {
		if i > 0 {
			f.Write(between)
		}
		fmt.Fprint(f, v)
		r = v
	}
	f.Write(atEnd)
	return Return(r)
}

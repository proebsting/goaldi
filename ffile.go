//  ffile.go -- I/O functions

package goaldi

import (
	"fmt"
	"os"
)

func init() {
	LibProcedure("write", Write)
	LibProcedure("writes", Writes)
	LibProcedure("print", Print)
	LibProcedure("println", Println)
	LibGoFunc("printf", fmt.Printf) // Go library function
}

var noBytes = []byte("")
var spByte = []byte(" ")
var nlByte = []byte("\n")

func Write(env *Env, a ...Value) (Value, *Closure) {
	return Wrt(noBytes, nlByte, a)
}

func Writes(env *Env, a ...Value) (Value, *Closure) {
	return Wrt(noBytes, noBytes, a)
}

func Print(env *Env, a ...Value) (Value, *Closure) {
	return Wrt(spByte, noBytes, a)
}

func Println(env *Env, a ...Value) (Value, *Closure) {
	return Wrt(spByte, nlByte, a)
}

func Wrt(between []byte, atEnd []byte, a []Value) (Value, *Closure) {
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

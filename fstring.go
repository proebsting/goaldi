//  fstring.go -- string functions

package goaldi

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

//  Declare methods
var StringMethods = MethodTable([]*VProcedure{
	DefMeth("type", (*VString).Type, []string{}, "return string type"),
	DefMeth("copy", (*VString).Copy, []string{}, "return string value"),
	DefMeth("string", (*VString).String, []string{}, "return string value"),
	DefMeth("image", (*VString).GoString, []string{}, "return string image"),
})

//  VString.Field implements methods
func (v *VString) Field(f string) Value {
	return GetMethod(StringMethods, v, f)
}

//  This init function adds a set of Go functions to the standard library
func init() {
	// Goaldi procedures
	LibProcedure("string", String)
	LibProcedure("char", Char)
	LibProcedure("ord", Ord)
	LibProcedure("reverse", Reverse)
	// Go library functions
	LibGoFunc("contains", strings.Contains)
	LibGoFunc("containsany", strings.ContainsAny)
	LibGoFunc("equalfold", strings.EqualFold)
	LibGoFunc("fields", strings.Fields)
	LibGoFunc("regex", regexp.Compile)
	LibGoFunc("regexp", regexp.CompilePOSIX)
	LibGoFunc("replace", strings.Replace)
	LibGoFunc("repl", strings.Repeat)
	LibGoFunc("split", strings.Split)
	LibGoFunc("toupper", strings.ToUpper)
	LibGoFunc("tolower", strings.ToLower)
	LibGoFunc("trim", strings.Trim)
}

//  String(x) -- return argument converted to string (always)
//  The result is exactly the same value used by write(x) etc.
func String(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("string", args)
	v := ProcArg(args, 0, NilValue)
	return Return(NewString(fmt.Sprint(v)))
}

//  Char(i) -- return one-character string with Unicode value i
func Char(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("char", args)
	var r [1]rune
	i := int(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	if i < 0 || i > int(unicode.MaxRune) {
		panic(NewExn("character code out of range", args[0]))
	}
	r[0] = rune(i)
	return Return(RuneString(r[:]))
}

//  Ord(c) -- return Unicode value of one-character string
func Ord(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("ord", args)
	r := ProcArg(args, 0, NilValue).(Stringable).ToString().ToRunes()
	if len(r) != 1 {
		panic(NewExn("string length not 1", args[0]))
	}
	return Return(NewNumber(float64(r[0])))
}

//  Reverse(s) -- return mirror image of string
func Reverse(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("reverse", args)
	r := ProcArg(args, 0, NilValue).(Stringable).ToString().ToRunes()
	n := len(r)
	for i := 0; i < n/2; i++ {
		r[i], r[n-1-i] = r[n-1-i], r[i]
	}
	return Return(RuneString(r))
}

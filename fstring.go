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
	DefMeth((*VString).Type, "type", "", "return string type"),
	DefMeth((*VString).Copy, "copy", "", "return string value"),
	DefMeth((*VString).String, "string", "", "return string value"),
	DefMeth((*VString).GoString, "image", "", "return string image"),
})

//  VString.Field implements methods
func (v *VString) Field(f string) Value {
	return GetMethod(StringMethods, v, f)
}

//  This init function adds a set of Go functions to the standard library
func init() {
	// Goaldi procedures
	DefLib(Char, "char", "n", "interpret rune as one-character string")
	DefLib(Ord, "ord", "s", "return Unicode ordinal of one-character string")
	DefLib(Reverse, "reverse", "s", "return mirror image of string")
	// Go library functions
	GoLib(strings.Contains, "contains", "s,substr", "return 1 if substr is in s")
	GoLib(strings.ContainsAny, "containsany", "s,chars", "return 1 if any char is in s")
	GoLib(strings.EqualFold, "equalfold", "s,t", "return 1 if s==t with case folding")
	GoLib(strings.Fields, "fields", "s", "return fields of s broken by whitespace")
	GoLib(regexp.Compile, "regex", "expr", "compile Go regular expression")
	GoLib(regexp.CompilePOSIX, "regexp", "expr", "compile POSIX regular expression")
	GoLib(strings.Replace, "replace", "s,old,new", "return s with new replacing old")
	GoLib(strings.Repeat, "repl", "s,count", "return s repeated n times")
	GoLib(strings.Split, "split", "s,sep", "return fields separted by sep")
	GoLib(strings.ToUpper, "toupper", "s", "convert to upper case")
	GoLib(strings.ToLower, "tolower", "s", "convert to lower case")
	GoLib(strings.Trim, "trim", "s,cutset", "remove leading and trailing characters")
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

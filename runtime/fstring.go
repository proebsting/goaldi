//  fstring.go -- string functions

package runtime

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

//  This init function adds a set of Go functions to the standard library
func init() {
	// Goaldi procedures
	DefLib(Char, "char", "n", "return single character for Unicode value")
	DefLib(Ord, "ord", "s", "return Unicode ordinal of single character")
	DefLib(Reverse, "reverse", "s", "return mirror image of string")
	DefLib(Left, "left", "s,w,p", "left-justify with padding p to width w")
	DefLib(Center, "center", "s,w,p", "center with padding p to width w")
	DefLib(Right, "right", "s,w,p", "right-justify with padding p to width w")
	DefLib(Unquote, "unquote", "s", "remove delimiters and escapes from s")
	// Go library functions
	GoLib(strings.Contains, "contains", "s,substr", "return 1 if substr is in s")
	GoLib(strings.ContainsAny, "containsany", "s,chars", "return 1 if any char is in s")
	GoLib(strings.EqualFold, "equalfold", "s,t", "return 1 if s==t with case folding")
	GoLib(strings.Fields, "fields", "s", "return fields of s delimited by whitespace")
	GoLib(regexp.Compile, "regex", "expr", "compile Go regular expression")
	GoLib(regexp.CompilePOSIX, "regexp", "expr", "compile POSIX regular expression")
	GoLib(strconv.Quote, "quote", "s", "add quotation marks and escapes to s")
	GoLib(strings.Replace, "replace", "s,old,new", "return s with new replacing old")
	GoLib(strings.Repeat, "repl", "s,count", "concatenate copies of s")
	GoLib(strings.Split, "split", "s,sep", "return fields delimited by sep")
	GoLib(strings.ToUpper, "toupper", "s", "convert to upper case")
	GoLib(strings.ToLower, "tolower", "s", "convert to lower case")
	GoLib(strings.Trim, "trim", "s,cutset", "remove leading and trailing characters")
}

//  string(x) returns a string representation of x.
//  The result is identical to the value used by write(x) or sprintf("%v",x).
func String(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("string", args)
	v := ProcArg(args, 0, NilValue)
	return Return(NewString(fmt.Sprint(v)))
}

//  char(n) returns the one-character string corresponding to the
//  Unicode value of n truncated to integer.
func Char(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("char", args)
	var r [1]rune
	i := int(FloatVal(ProcArg(args, 0, NilValue)))
	if i < 0 || i > int(unicode.MaxRune) {
		panic(NewExn("Character code out of range", args[0]))
	}
	r[0] = rune(i)
	return Return(RuneString(r[:]))
}

//  ord(s) returns the Unicode value corresponding to the one-character string s.
func Ord(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("ord", args)
	r := ProcArg(args, 0, NilValue).(Stringable).ToString().ToRunes()
	if len(r) != 1 {
		panic(NewExn("String length not 1", args[0]))
	}
	return Return(NewNumber(float64(r[0])))
}

//  left(s,w,p) left-justifies s in a string of width w, padding with p.
func Left(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("left", args)
	s := ProcArg(args, 0, NilValue).(Stringable).ToString().ToRunes()
	w := int(FloatVal(ProcArg(args, 1, ONE)))
	p := ProcArg(args, 2, SPACE).(Stringable).ToString().ToRunes()
	if len(p) == 0 {
		panic(NewExn("Empty padding string", args[2]))
	}
	r := make([]rune, w)
	copy(r, s)
	n := w - len(s)
	for i := 0; i < n; i++ {
		j := i % len(p)
		r[w-i-1] = p[len(p)-j-1]
	}
	return Return(RuneString(r))
}

//  right(s,w,p) right-justifies s in a string of width w, padding with p.
func Right(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("right", args)
	s := ProcArg(args, 0, NilValue).(Stringable).ToString().ToRunes()
	w := int(FloatVal(ProcArg(args, 1, ONE)))
	p := ProcArg(args, 2, SPACE).(Stringable).ToString().ToRunes()
	if len(p) == 0 {
		panic(NewExn("Empty padding string", args[2]))
	}
	n := w - len(s)
	if n > 0 {
		r := make([]rune, w)
		copy(r[n:], s)
		for i := 0; i < n; i++ {
			j := i % len(p)
			r[i] = p[j]
		}
		return Return(RuneString(r))
	} else {
		return Return(RuneString(s[-n:]))
	}
}

//  center(s,w,p) centers s in a string of width w, padding with p.
func Center(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("center", args)
	s := ProcArg(args, 0, NilValue).(Stringable).ToString().ToRunes()
	w := int(FloatVal(ProcArg(args, 1, ONE)))
	p := ProcArg(args, 2, SPACE).(Stringable).ToString().ToRunes()
	if len(p) == 0 {
		panic(NewExn("Empty padding string", args[2]))
	}
	n := w - len(s) // amount of padding needed
	if n > 0 {      // if any
		r := make([]rune, w)      // result
		nl := n / 2               // left-side padding count
		nr := n - nl              // right-side padding count
		copy(r[nl:], s)           // original string
		for i := 0; i < nl; i++ { // pad left
			j := i % len(p)
			r[i] = p[j]
		}
		for i := 0; i < nr; i++ { // pad right
			j := i % len(p)
			r[w-i-1] = p[len(p)-j-1]
		}
		return Return(RuneString(r))
	} else { // no padding needed
		i := (-n + 1) / 2
		return Return(RuneString(s[i : i+w]))
	}
}

//  reverse(s) returns the end-for-end reversal of the string s.
func Reverse(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("reverse", args)
	r := ProcArg(args, 0, NilValue).(Stringable).ToString().ToRunes()
	n := len(r)
	for i := 0; i < n/2; i++ {
		r[i], r[n-1-i] = r[n-1-i], r[i]
	}
	return Return(RuneString(r))
}

//  unquote(s) removes delimiters and escapes from a quoted string.
//  The argument s must begin and end with explicit "double quotes" or
//  \`backticks`.  unquote() fails if s is not properly quoted or if it
//  contains an invalid (by Go rules) escape sequence.
func Unquote(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("unquote", args)
	s := ProcArg(args, 0, NilValue).(Stringable).ToString().ToUTF8()
	s, err := strconv.Unquote(s)
	if err != nil {
		return Fail()
	} else {
		return Return(NewString(s))
	}
}

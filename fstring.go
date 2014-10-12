//  fstring.go -- string functions

package goaldi

import (
	"strings"
)

//  This init function adds a set of Go functions to the standard library
func init() {
	LibGoFunc("equalfold", strings.EqualFold)
	LibGoFunc("replace", strings.Replace)
	LibGoFunc("toupper", strings.ToUpper)
	LibGoFunc("tolower", strings.ToLower)
	LibGoFunc("trim", strings.Trim)
	LibProcedure("reverse", Reverse)
}

//  Reverse(s) -- return mirror image of string
func Reverse(env *Env, a ...Value) (Value, *Closure) {
	r := ProcArg(a, 0, NilValue).(Stringable).ToString().ToRunes()
	n := len(r)
	for i := 0; i < n/2; i++ {
		r[i], r[n-1-i] = r[n-1-i], r[i]
	}
	return Return(RuneString(r))
}

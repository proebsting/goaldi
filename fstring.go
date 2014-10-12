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
}

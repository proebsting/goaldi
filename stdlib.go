//  stdlib.go -- standard library and miscellaneous functions

//  #%#% this initial set is for testing and illustration; it is NOT final!

package goaldi

import (
	"fmt"
	"os"
	"strings"
)

//  StdLib is the set of procedures available at link time
var StdLib = make(map[string]*VProcedure)

//  LibProc registers a Go function as a standard library procedure.
//  This must be done before linking (e.g. via init func) to be effective.
func LibProc(name string, f interface{}) {
	StdLib[name] = GoProcedure(name, f)
}

//  This init function adds a set of Go functions to the standard library
func init() {

	LibProc("equalfold", strings.EqualFold)
	LibProc("replace", strings.Replace)
	LibProc("toupper", strings.ToUpper)
	LibProc("tolower", strings.ToLower)
	LibProc("trim", strings.Trim)

	LibProc("print", fmt.Print)
	LibProc("println", fmt.Println)
	LibProc("printf", fmt.Printf)
	LibProc("fprint", fmt.Fprint)
	LibProc("fprintln", fmt.Fprintln)
	LibProc("fprintf", fmt.Fprintf)
	LibProc("write", fmt.Println) // not like Icon: no file, spacing
	LibProc("writes", fmt.Print)  // not like Icon: no file, spacing

	LibProc("exit", os.Exit)
	LibProc("remove", os.Remove)
}

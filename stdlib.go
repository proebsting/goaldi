//  stdlib.go -- standard library support
//
//  Most library procedure definitions are grouped by type and defined in f*.go.

package goaldi

import (
	"fmt"
	"io"
)

//  StdLib is the set of procedures (including types) initially available
var StdLib = make(map[string]ICall)

//  ShowLibrary(f) lists all library functions on file f
func ShowLibrary(f io.Writer) {
	fmt.Fprintln(f)
	fmt.Fprintln(f, "Standard Library")
	fmt.Fprintln(f, "------------------------------")
	columns := "%-12s %s\n"
	for k := range SortedKeys(StdLib) {
		x := StdLib[k]
		switch v := x.(type) {
		case *VProcedure:
			fmt.Fprintf(f, columns, k, v.ImplBy())
		case *VType:
			fmt.Fprintf(f, columns, k, "[standard type]")
		}
	}
}

//  DefLib constructs and registers a standard library procedure.
func DefLib(entry Procedure, name string, pspec string, descr string) *VProcedure {
	p := DefProc(entry, name, pspec, descr)
	StdLib[name] = p
	return p
}

//#%#% to be replaced by DefLib(above)
//  LibProcedure registers a standard library procedure taking Goaldi arguments.
//  This must be done before linking (e.g. via init func) to be effective.
func LibProcedure(name string, p Procedure) {
	StdLib[name] = NewProcedure(name, nil /*#%#% TIGHTEN */, true, p, p, "")
}

//  LibGoFunc registers a Go function as a standard library procedure.
//  This must be done before linking (e.g. via init func) to be effective.
func LibGoFunc(name string, f interface{}) {
	StdLib[name] = GoProcedure(name, f)
}

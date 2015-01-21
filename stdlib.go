//  stdlib.go -- standard library support
//
//  Most library procedure definitions are grouped by type and defined in f*.go.

package goaldi

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
)

//  StdLib is the set of procedures (including types) initially available
var StdLib = make(map[string]ICall)

//  DefLib constructs and registers a standard library procedure.
func DefLib(entry Procedure, name string, pspec string, descr string) *VProcedure {
	p := DefProc(entry, name, pspec, descr)
	StdLib[name] = p
	return p
}

//  GoLib registers a Go function as a standard library procedure
func GoLib(entry interface{}, name string, pspec string, descr string) *VProcedure {
	pnames, isvar := ParmsFromSpec(pspec)
	p := NewProcedure(name, pnames, isvar, GoShim(name, entry), entry, descr)
	StdLib[name] = p
	return p
}

//  ShowLibrary(f) lists all library functions and standard types on file f
func ShowLibrary(f io.Writer) {
	typelist := make([]*VType, 0)
	linelen := 79
	fmt.Fprintln(f)
	fmt.Fprintln(f, "Standard Library")
	fmt.Fprintln(f, strings.Repeat("-", linelen))
	for k := range SortedKeys(StdLib) {
		x := StdLib[k]
		switch v := x.(type) {
		case *VProcedure:
			s1 := v.GoString()[10:] + " -- " + v.Descr
			s3 := v.ImplBy()
			s2 := strings.Repeat(" ", linelen-len(s1)-len(s3)-2)
			fmt.Fprintln(f, s1, s2, s3)
		case *VType:
			typelist = append(typelist, v)
		default:
			fmt.Fprintf(f, "%x : UNRECOGNIZED : %T\n", k, x)
		}
	}

	fmt.Fprintln(f)
	fmt.Fprintln(f, "Standard Types")
	fmt.Fprintln(f, "-------------------------------------------")
	for _, t := range typelist {
		ctor := t.Ctor
		fmt.Fprintln(f, ctor.GoString()[10:]+" -- "+ctor.Descr)
		if t.Methods != nil {
			for _, meth := range t.Methods {
				fmt.Fprintf(f, "    %s.%s -- %s\n",
					t.Abbr, meth.GoString()[10:], meth.Descr)
			}
		}
	}
}

//  VProcedure.ImplBy -- return name of implementing underlying function
func (v *VProcedure) ImplBy() string {
	if v.GoFunc == nil {
		return v.Name // no further information available
	} else {
		return runtime.FuncForPC(reflect.ValueOf(v.GoFunc).Pointer()).Name()
	}
}

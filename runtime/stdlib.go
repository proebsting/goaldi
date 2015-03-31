//  stdlib.go -- standard library support
//
//  Most library procedure definitions are grouped by type and defined in f*.go.

package runtime

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

//  GoLib registers a Go function as a standard library procedure.
//  The ETOSS option is used to make regex, printf, remove, etc. throw
//  an exceptions when an error occurs.
func GoLib(entry interface{}, name string, pspec string, descr string) *VProcedure {
	pnames, isvar := ParmsFromSpec(pspec)
	p := NewProcedure(name, pnames, isvar,
		GoShim(name, entry, ETOSS), entry, descr)
	StdLib[name] = p
	return p
}

const showLen = 79

//  ShowLibrary(f) lists all library functions and standard types on file f
func ShowLibrary(f io.Writer) {
	typelist := make([]*VType, 0)
	hrule := strings.Repeat("-", showLen)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "Standard Library")
	fmt.Fprintln(f, hrule)
	for k := range SortedKeys(StdLib) {
		x := StdLib[k]
		switch v := x.(type) {
		case *VProcedure:
			showProc(f, "", v)
		case *VType:
			typelist = append(typelist, v)
			showProc(f, "", v.Ctor)
		case *VCtor:
			// ignore (e.g. elemtype)
		default:
			fmt.Fprintf(f, "%x : UNRECOGNIZED : %T\n", k, x)
		}
	}

	fmt.Fprintln(f, hrule)
	for k := range SortedKeys(UniMethods) {
		showProc(f, "x.", UniMethods[k])
	}
	for _, t := range typelist {
		if t.Methods != nil && len(t.Methods) > 0 {
			fmt.Fprintln(f, hrule)
			for k := range SortedKeys(t.Methods) {
				showProc(f, t.Abbr+".", t.Methods[k])
			}
		}
	}
	fmt.Fprintln(f, hrule)
}

//  showProc(f, c, p) -- format and print one-line procedure reference
//  f is the output file
//  c is a prefix (e.g. "x." or nothing)
//  p is the procedure
func showProc(f io.Writer, c string, p *VProcedure) {
	l := fmt.Sprintf("%s%s -- %s", c, p.GoString()[10:], p.Descr)
	r := p.ImplBy()
	n := showLen - len(l) - len(r)
	if n < 2 {
		n = 2
	}
	s := strings.Repeat(" ", n)
	fmt.Fprintln(f, l+s+r)
}

//  VProcedure.ImplBy -- return name of implementing underlying function
func (v *VProcedure) ImplBy() string {
	if v.GoFunc == nil {
		return v.Name // no further information available
	} else {
		return runtime.FuncForPC(reflect.ValueOf(v.GoFunc).Pointer()).Name()
	}
}

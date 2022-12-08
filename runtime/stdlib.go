//  stdlib.go -- standard library support
//
//  Most library procedure definitions are grouped by type
//  and defined in source files named f*.go.

package runtime

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

// StdLib is the set of procedures (including types) initially available
var StdLib = make(map[string]ICall)

// DefLib constructs and registers a standard library procedure.
func DefLib(entry Procedure, name string, pspec string, descr string) *VProcedure {
	p := DefProc(entry, name, pspec, descr)
	StdLib[name] = p
	return p
}

// GoLib registers a Go function as a standard library procedure.
// The ETOSS option is used to make regex, printf, remove, etc. throw
// an exceptions when an error occurs.
func GoLib(entry interface{}, name string, pspec string, descr string) *VProcedure {
	pnames, isvar := ParmsFromSpec(pspec)
	p := NewProcedure(name, pnames, isvar,
		GoShim(name, entry, ETOSS), entry, descr)
	StdLib[name] = p
	return p
}

const showLen = 79

// libentry registers a library entry for sorting and output
type libentry struct {
	tpfx string      // type prefix for method (e.g. "L.")
	name string      // procedure, constructor, or method name
	rank int         // type rank (secondary key)
	proc *VProcedure // underlying procedure
}

// bykey implements sort.Interface for []*libentry
type bykey []*libentry

func (a bykey) Len() int      { return len(a) }
func (a bykey) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a bykey) Less(i, j int) bool {
	if a[i].name != a[j].name {
		return a[i].name < a[j].name
	} else {
		return a[i].rank < a[j].rank
	}
}

// ShowLibrary(f) lists all library functions and standard types on file f
func ShowLibrary(f io.Writer) {

	types := make([]*VType, 0)    // types seen
	procs := make([]*libentry, 0) // procedures and methods to show

	// first enter simple procedure and type constructors (noting the types)
	for name, x := range StdLib {
		switch v := x.(type) {
		case *VProcedure:
			procs = append(procs, &libentry{"  ", name, -2, v})
		case *VType:
			types = append(types, v)
			procs = append(procs, &libentry{v.Abbr + " ", name, -3, v.Ctor})
		case *VCtor:
			// ignore (e.g. elemtype)
		default:
			fmt.Fprintf(f, "%x : UNRECOGNIZED : %T\n", name, x)
		}
	}

	// add the "universal" methods that can be applied to any value
	for name, proc := range UniMethods {
		procs = append(procs, &libentry{"x.", name, -1, proc})
	}

	// add the methods of the types we've seen, sorting in rank order
	for _, t := range types {
		tpfx := t.Abbr + "."
		for name, proc := range t.Methods {
			procs = append(procs, &libentry{tpfx, name, t.SortRank, proc})
		}
	}

	// sort the procedures and methods
	sort.Sort(bykey(procs))

	// output them
	hrule := strings.Repeat("-", showLen)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "Standard Library")
	fmt.Fprintln(f, hrule)
	for _, e := range procs {
		p := e.proc
		l := fmt.Sprintf("%s%s -- %s", e.tpfx, p.GoString()[10:], p.Descr)
		r := p.ImplBy()
		n := showLen - len(l) - len(r)
		if n < 2 {
			n = 2
		}
		s := strings.Repeat(" ", n)
		fmt.Fprintln(f, l+s+r)
	}
	fmt.Fprintln(f, hrule)
}

// VProcedure.ImplBy -- return name of implementing underlying function
func (v *VProcedure) ImplBy() string {
	if v.GoFunc == nil {
		return v.Name // no further information available
	} else {
		return runtime.FuncForPC(reflect.ValueOf(v.GoFunc).Pointer()).Name()
	}
}

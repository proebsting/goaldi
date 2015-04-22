//  main.go -- overall control of interpreter
//
//  If the first command line argument is "-x", then additional arguments
//  direct the loading and execution of IR code (gcode) from input files.
//
//  If not, the embedded translator app receives all arguments.

package main

import (
	"bytes"
	"fmt"
	_ "goaldi/extensions"
	"goaldi/ir"
	g "goaldi/runtime"
	"goaldi/translator"
	"io"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

//  globals

var PubSpace = g.GetSpace("")          // the public (unnamed) namespace
var Undeclared = make(map[string]bool) // is var x undeclared?

var GlobInit = make([]*ir.Ir_Global, 0)  // globals with initialization
var InitList = make([]*ir.Ir_Initial, 0) // sequential initialization blocks

var nFatals = 0   // count of fatal errors
var nWarnings = 0 // count of nonfatal errors

//  main is the overall supervisor.
func main() {

	// use all available processors
	runtime.GOMAXPROCS(runtime.NumCPU())

	// handle command line
	files, args := options()

	// start profiling if requested
	if opt_profile {
		pfile, err := os.Create("PROFILE")
		checkError(err)
		pprof.StartCPUProfile(pfile)
		defer pprof.StopCPUProfile()
	}

	// show library environment
	if opt_envmt {
		g.ShowLibrary(os.Stdout)
		g.ShowEnvironment(os.Stdout)
		fmt.Println()
	}

	// load the IR code
	parts := make([][]interface{}, 0)
	if files == nil {
		bbuf := bytes.NewBuffer(translator.GCode)
		parts = append(parts, loadfile("[embedded]", bbuf)...)
	} else if len(files) == 0 {
		parts = append(parts, loadfile("[stdin]", os.Stdin)...)
	} else {
		for _, fname := range files {
			f, err := os.Open(fname)
			checkError(err)
			parts = append(parts, loadfile(fname, f)...)
			if opt_delete {
				os.Remove(fname)
			}
		}
	}
	showInterval("loading")

	// quit now if this was just a run to get an assembly listing
	if opt_noexec && opt_adump {
		quit(0)
	}

	// link everything together
	link(parts)
	showInterval("linking")
	if nFatals > 0 {
		quit(1)
	}

	// quit now if -c was given
	if opt_noexec {
		quit(0)
	}

	// set environment flag if to dump Go stack on panic
	if opt_debug {
		g.EnvInit("gostack", g.ONE)
	}

	// make a list for dependency-based global initialization
	dlist := &g.DependencyList{}
	// put procedures at the front of the list for proper dependency checking
	// (excluding procedures associated with global:= and initial{})
	for _, proc := range ProcTable {
		if !strings.Contains(proc.name, "$global$") &&
			!strings.Contains(proc.name, "$initial$") {
			ulist := proc.ir.UnboundList
			if ulist != nil && len(ulist) > 0 {
				dlist.Add(proc.qname, nil, ulist)
			}
		}
	}
	// enter all globals that initialize
	for _, gi := range GlobInit {
		p := ProcTable[gi.Fn].vproc
		uses := ProcTable[gi.Fn].ir.UnboundList
		q := g.GetSpace(gi.Namespace).GetQual()
		dlist.Add(q+gi.Name, p, uses)
	}
	// reorder the list for dependencies
	err := dlist.Reorder(opt_trace)
	if err != nil {
		abort(fmt.Sprintf("fatal   %v\n", err))
	}

	// before running any initialization code, make sure main() exists
	gmain := PubSpace.Get("main")
	if gmain == nil {
		abort("no main procedure")
	}
	if gv, ok := gmain.(g.IVariable); ok {
		gmain = gv.Deref()
	}

	// run the sequence of initialization procedures
	dlist.RunAll()                // global initializers as reordered
	for _, ip := range InitList { // initial{} blocks in lexical order
		g.Run(ProcTable[ip.Fn].vproc, []g.Value{})
	}
	showInterval("initialization")

	// execute main()
	arglist := make([]g.Value, 0)
	for _, s := range args {
		arglist = append(arglist, g.NewString(s))
	}
	g.Run(gmain, arglist)

	// exit
	showInterval("execution")
	g.Shutdown(0)
}

//  loadfile(label, reader) -- load and possibly print one file
func loadfile(label string, rdr io.Reader) [][]interface{} {
	_, parts := ir.Load(rdr)
	if opt_adump {
		for _, p := range parts {
			ir.Print(label, p)
		}
	}
	return parts
}

//  warning -- report nonfatal error and continue
func warning(s string) {
	nWarnings++
	fmt.Fprintf(os.Stderr, "Warning: %s\n", s)
}

//  fatal -- report fatal error (but continue)
func fatal(s string) {
	nFatals++
	fmt.Fprintf(os.Stderr, "Fatal:   %s\n", s)
}

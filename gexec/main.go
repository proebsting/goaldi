//  main.go -- overall control of interpreter

package main

import (
	"fmt"
	g "goaldi"
	_ "goaldi/extensions"
	"os"
	"runtime/pprof"
	"strings"
)

//  globals

var PubSpace = g.GetSpace("")          // the public (unnamed) namespace
var Undeclared = make(map[string]bool) // is var x undeclared?

var GlobInit = make([]*ir_Global, 0)  // globals with initialization
var InitList = make([]*ir_Initial, 0) // sequential initialization blocks

var nFatals = 0   // count of fatal errors
var nWarnings = 0 // count of nonfatal errors

//  main is the overall supervisor.
func main() {

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
	if len(files) == 0 {
		parts = append(parts, load("-"))
	} else {
		for _, f := range files {
			parts = append(parts, load(f))
		}
	}
	showInterval("loading")

	// link everything together
	link(parts)
	showInterval("linking")
	if nFatals > 0 {
		pprof.StopCPUProfile()
		os.Exit(1)
	}

	// list the globals
	if opt_verbose {
		for nsname := range g.AllSpaces() {
			ns := g.GetSpace(nsname)
			if nsname == "" {
				fmt.Printf("\nGLOBALS:")
			} else {
				fmt.Printf("\n%s::", nsname)
			}
			for name := range ns.All() {
				fmt.Printf(" %s", name)
				if _, ok := ns.Get(name).(*g.VProcedure); ok {
					fmt.Print("()")
				}
			}
			fmt.Printf("\n")
		}
	}

	// quit now if -c was given
	if opt_noexec {
		pprof.StopCPUProfile()
		os.Exit(0)
	}

	// set execution flag
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
	for _, ir := range GlobInit {
		p := ProcTable[ir.Fn].vproc
		uses := ProcTable[ir.Fn].ir.UnboundList
		q := g.GetSpace(ir.Namespace).GetQual()
		dlist.Add(q+ir.Name, p, uses)
	}
	// reorder the list for dependencies
	err := dlist.Reorder(opt_trace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal:   %v\n", err)
		pprof.StopCPUProfile()
		os.Exit(1)
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
	//#%#% each call to Run resets a clean environment. is that valid?
	dlist.RunAll()                // global initializers as reordered
	for _, ir := range InitList { // initial{} blocks in lexical order
		g.Run(ProcTable[ir.Fn].vproc, []g.Value{})
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

//  warning -- report nonfatal error and continue
func warning(s string) {
	nWarnings++
	fmt.Fprintf(os.Stderr, "warning: %s\n", s)
}

//  fatal -- report fatal error (but continue)
func fatal(s string) {
	nFatals++
	fmt.Fprintf(os.Stderr, "fatal:   %s\n", s)
}

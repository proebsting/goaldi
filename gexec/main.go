//  main.go -- overall control of interpreter

package main

import (
	"fmt"
	g "goaldi"
	"os"
	"runtime/pprof"
	"unicode"
)

//  globals

var GlobalDict = make(map[string]g.Value) // global dictionary
var Undeclared = make(map[string]bool)    // is var x undeclared?

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
		fmt.Printf("\nGLOBALS:")
		for name := range g.SortedKeys(GlobalDict) {
			fmt.Printf(" %s", name)
			if _, ok := GlobalDict[name].(*g.VProcedure); ok {
				fmt.Print("()")
			}
		}
		fmt.Printf("\n")
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

	// run the interdependent global initialization procedures
	ilist := make([]*g.InitItem, 0)
	for _, ir := range GlobInit { // enter all globals that initialize
		p := GlobalDict[ir.Fn].(*g.VProcedure)
		uses := ProcTable[ir.Fn].ir.UnboundList
		ilist = append(ilist, g.NewInit(p, uses, ir.NameList[0]))
	}
	// need to factor in the dependencies of called procedures, too
	for _, proc := range ProcTable { // enter real procedures that ref globals
		if !unicode.IsDigit(rune(proc.name[0])) {
			// this is a top-level user-declared procedure
			ulist := proc.ir.UnboundList
			if ulist != nil && len(ulist) > 0 {
				ilist = append(ilist, g.NewInit(nil, ulist, proc.name))
			}
		}
	}
	err := g.RunDep(ilist, opt_trace) // init globals in dependency order
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal:   %v\n", err)
		pprof.StopCPUProfile()
		os.Exit(1)
	}

	// run the sequence of initialization procedures
	//#%#% each call to Run resets a clean environment. is that valid?
	for _, ir := range InitList {
		g.Run(GlobalDict[ir.Fn].(*g.VProcedure), []g.Value{})
	}
	showInterval("initialization")

	// find and execute main()
	arglist := make([]g.Value, 0)
	for _, s := range args {
		arglist = append(arglist, g.NewString(s))
	}
	gmain := GlobalDict["main"]
	if gmain == nil {
		abort("no main procedure")
	}
	if gv, ok := gmain.(g.IVariable); ok {
		gmain = gv.Deref()
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

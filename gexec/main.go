//  main.go -- overall control of interpreter

package main

import (
	"fmt"
	g "goaldi"
	"os"
)

type UNKNOWN interface{} // temporary designation for type TBD

//  globals

var GlobalDict = make(map[string]g.Value)
var Undeclared = make(map[string]bool)

var nFatals = 0   // count of fatal errors
var nWarnings = 0 // count of nonfatal errors

//  main is the overall supervisor.
func main() {

	// handle command line
	files, args := options()

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
		os.Exit(1)
	}

	// list the globals
	if opt_verbose {
		fmt.Printf("\nGLOBALS:")
		for name := range sortedKeys(GlobalDict) {
			fmt.Printf(" %s", name)
			if _, ok := GlobalDict[name].(*g.VProcedure); ok {
				fmt.Print("()")
			}
		}
		fmt.Printf("\n")
	}

	// quit now if -c was given
	if opt_noexec {
		os.Exit(0)
	}

	// set up recovery code
	//#%#% TBD

	// find and execute main()
	env := &g.Env{}
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
	gmain.(g.ICall).Call(env, arglist...)

	// exit
	showInterval("execution")
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

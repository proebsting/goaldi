//  main.go -- overall control of interpreter

package main

import (
	"fmt"
	g "goaldi"
	"os"
)

type UNKNOWN interface{} // temporary designation for type TBD

var GlobalDict = make(map[string]g.Value)
var Undeclared = make(map[string]bool)

var nFatals = 0   // count of fatal errors
var nWarnings = 0 // count of nonfatal errors

//  main is the overall supervisor.
func main() {
	files, args := options()
	parts := make([][]interface{}, 0)
	if len(files) == 0 {
		parts = append(parts, load("-"))
	} else {
		for _, f := range files {
			parts = append(parts, load(f))
		}
	}
	showInterval("loading")
	prog := link(parts)
	showInterval("linking")

	if opt_verbose {
		fmt.Printf("\nGLOBALS:")
		for name, value := range GlobalDict {
			fmt.Printf(" %s", name)
			if _, ok := value.(*g.VProcedure); ok {
				fmt.Print("()")
			}
		}
		fmt.Printf("\n")
	}

	if nFatals > 0 {
		os.Exit(1)
	}

	//#%#%#% to be continued...
	_ = prog
	_ = args
	return
	// run(prog, args)
	// showInterval("execution")
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

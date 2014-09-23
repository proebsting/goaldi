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

	fmt.Printf("\nGLOBALS:")
	for k, _ := range GlobalDict {
		fmt.Printf(" %s", k)
	}
	fmt.Printf("\nUNDECLARED:")
	for k, _ := range Undeclared {
		fmt.Printf(" %s", k)
	}
	fmt.Printf("\n")

	os.Exit(0) //#%#%#%#%#%#%#%
	run(prog, args)
	showInterval("execution")
}

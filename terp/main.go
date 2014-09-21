//  main.go -- overall control of interpreter

package main

import (
	"os"
)

type UNKNOWN interface{} // temporary designation for type TBD

//  main is the overall supervisor.
func main() {
	files, args := options()
	parts := make([]UNKNOWN, 0)
	if len(files) == 0 {
		parts = append(parts, load("-"))
	} else {
		for _, f := range files {
			parts = append(parts, load(f))
		}
	}
	os.Exit(0) //#%#%#%#%#%#%#%
	prog := link(parts)
	showInterval("loading")
	run(prog, args)
	showInterval("execution")
}

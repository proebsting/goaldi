//  main.go -- overall control of interpreter

package main

import ()

type UNKNOWN interface{} // temporary designation for type TBD

//  main is the overall supervisor.
func main() {
	files, args := options()
	parts := make([]UNKNOWN, 0)
	for _, f := range files {
		parts = append(parts, load(f))
	}
	prog := link(parts)
	showInterval("loading")
	run(prog, args)
	showInterval("execution")
}

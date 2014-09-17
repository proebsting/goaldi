//  options.go -- declaration and processing of command line arguments

package main

import (
	"flag"
	"fmt"
	"os"
)

//  command-line options
var opt_timings bool // show CPU timings
var opt_verbose bool // issue verbose commentary

//  usage prints a usage message (with option descriptions) and aborts.
func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [options] file [args]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

//  options sets global flags and returns file names and execution arguments.
func options() (files []string, args []string) {

	flag.BoolVar(&opt_timings, "t", true, "show CPU timings")
	flag.BoolVar(&opt_verbose, "v", true, "issue verbose commentary")
	flag.Usage = usage
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		usage()
	}
	files = append(files, args[0])
	args = args[1:]
	return files, args
}

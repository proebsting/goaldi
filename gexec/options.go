//  options.go -- declaration and processing of command line arguments

package main

import (
	"flag"
	"fmt"
	"os"
)

//  command-line options
var opt_noexec bool  // -l: load and link only; don't execute
var opt_timings bool // -t: show CPU timings
var opt_verbose bool // -v: issue verbose commentary
var opt_adump bool   // -A: dump assembly-style IR code
var opt_debug bool   // -D: set debug flag (dump Go stack on panic)
var opt_tally bool   // -F: tally (static) IR field usage
var opt_jdump bool   // -J: dump JSON in outline form
var opt_profile bool // -P: produce CPU profile on ./PROFILE
var opt_trace bool   // -T: trace IR instruction execution

//  usage prints a usage message (with option descriptions) and aborts.
func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [options] [file [args]]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

//  options sets global flags and returns file names and execution arguments.
func options() (files []string, args []string) {

	flag.BoolVar(&opt_noexec, "l", false, "load and link only")
	flag.BoolVar(&opt_timings, "t", false, "show CPU timings")
	flag.BoolVar(&opt_verbose, "v", false, "issue verbose commentary")
	flag.BoolVar(&opt_adump, "A", false, "dump assembly-style IR code")
	flag.BoolVar(&opt_debug, "D", false, "dump Go stack on panic")
	flag.BoolVar(&opt_tally, "F", false, "tally IR field usage")
	flag.BoolVar(&opt_jdump, "J", false, "dump JSON IR in outline form")
	flag.BoolVar(&opt_profile, "P", false, "produce CPU profile on ./PROFILE")
	flag.BoolVar(&opt_trace, "T", false, "trace IR instruction execution")
	flag.Usage = usage
	flag.Parse()
	args = flag.Args()
	if len(args) > 0 {
		files = append(files, args[0])
		args = args[1:]
	}
	return files, args
}

//  options.go -- declaration and processing of command line arguments
//
//  NOTE:  If the first command line argument is not "-x", then
//  no argument processing is done under the assumption that all
//  options and arguments will be passed to the embedded app.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// command-line options
var opt_noexec bool  // -l: load and link only; don't execute
var opt_timings bool // -t: show CPU timings
var opt_adump bool   // -A: dump assembly-style IR code
var opt_debug bool   // -D: set debug flag (dump Go stack on panic)
var opt_init bool    // -I: trace initialization ordering
var opt_envmt bool   // -E: show initial environment before loading
var opt_profile bool // -P: produce CPU profile on ./PROFILE
var opt_trace bool   // -T: trace IR instruction execution
var opt_delete bool  // -#: delete IR files after loading

// usage prints a usage message (with option descriptions) and aborts.
func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s -x [options] file.gir... [--] [arg...]]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

// options sets global flags and returns file names and execution arguments.
func options() (files []string, args []string) {

	// check for enabling magic flag
	// (if not set, return files=nil as an indicator)
	if len(os.Args) < 2 || os.Args[1] != "-x" {
		return nil, os.Args[1:]
	}

	flag.Bool("x", false, "process command line as described here")
	flag.BoolVar(&opt_noexec, "l", false, "load and link only")
	flag.BoolVar(&opt_timings, "t", false, "show CPU timings")
	flag.BoolVar(&opt_adump, "A", false, "dump assembly-style IR code")
	flag.BoolVar(&opt_debug, "D", false, "dump Go stack on panic")
	flag.BoolVar(&opt_init, "I", false, "trace initialization ordering")
	flag.BoolVar(&opt_envmt, "E", false, "show initial environment")
	flag.BoolVar(&opt_profile, "P", false, "produce ./PROFILE file (Linux)")
	flag.BoolVar(&opt_trace, "T", false, "trace IR instruction execution")
	flag.BoolVar(&opt_delete, "#", false, "delete IR files after loading")
	flag.Usage = usage
	flag.Parse()

	// get remaining (positional) command arguments
	args = flag.Args()
	if len(args) == 0 { // must have at least one
		usage()
	}
	files = append(files, args[0]) // first argument is always a file
	args = args[1:]

	// any immediately following args that end in ".gir" are also files to load
	for len(args) > 0 && strings.HasSuffix(args[0], ".gir") {
		files = append(files, args[0])
		args = args[1:]
	}

	// a "--" argument is a separator to be removed
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	return files, args
}

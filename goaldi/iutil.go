//  iutil.go -- interpreter utility routines

package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"syscall"
	"time"
)

//  checkError aborts if error value e is not nil.
func checkError(e error) {
	if e != nil {
		abort(e)
	}
}

//  abort issues an error message and aborts.
func abort(e interface{}) {
	fmt.Fprintln(os.Stderr, e)
	quit(1)
}

//  quit exits with a given code after stopping profiling.
func quit(xc int) {
	pprof.StopCPUProfile()
	os.Exit(xc)
}

//  babble prints commentary on Stderr if opt_verbose is set.
//  The first argument is a printf format.  A newline is added automatically.
func babble(format string, values ...interface{}) {
	if opt_verbose {
		fmt.Fprintf(os.Stderr, format, values...)
		fmt.Fprintln(os.Stderr)
	}
}

//  showInterval prints timing for the latest interval if opt_timings is set.
func showInterval(label string) {
	dt := cpuInterval().Seconds()
	if label != "" && opt_timings {
		fmt.Fprintf(os.Stderr, "%7.3f %s\n", dt, label)
	}
}

//  cpuInterval returns the CPU time (user + system) since the preceding call.
func cpuInterval() time.Duration {
	total := cpuTime()
	delta := total - prevCPU
	prevCPU = total
	return delta
}

var prevCPU time.Duration // total time at list check

//  cpuTime returns the current CPU usage (user time + system time).
func cpuTime() time.Duration {
	var ustruct syscall.Rusage
	checkError(syscall.Getrusage(0, &ustruct))
	user := time.Duration(syscall.TimevalToNsec(ustruct.Utime))
	sys := time.Duration(syscall.TimevalToNsec(ustruct.Stime))
	return user + sys
}

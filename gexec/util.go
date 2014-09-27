//  util.go -- general-purpose utility routines

package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"syscall"
	"time"
	"unicode"
	"unicode/utf8"
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
	os.Exit(1)
}

//  sortedKeys generates in order (over a channel) the keys of a map[string].
//  usage:  for k := range sortedKeys(mymap) { ... }
func sortedKeys(m interface{}) chan string {
	vlist := reflect.ValueOf(m).MapKeys()
	n := len(vlist)
	slist := make([]string, n, n)
	for i, k := range vlist {
		slist[i] = k.String()
	}
	sort.Strings(slist)
	ch := make(chan string, n)
	go func() {
		for _, k := range slist {
			ch <- k
		}
		close(ch)
	}()
	return ch
}

//  assert panics if the test argument is false
func assert(test bool, err string) {
	if !test {
		panic("assertion failed: " + err)
	}
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
	return 0
}

//  Capitalize -- convert first character of string to upper case
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}
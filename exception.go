//  exception.go -- things dealing with exceptions and panics

package goaldi

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

//  RunErr records a Goaldi runtime error
type RunErr struct {
	msg  string // explanatory message
	offv Value  // offending value
}

//  RunErr.String() returns a string form of a RunErr
func (e *RunErr) String() string {
	return fmt.Sprintf("RunErr: %s (%v)", e.msg, e.offv)
}

//  CallFrame records one frame of traceback information
type CallFrame struct {
	cause interface{} // underlying panic call
	offv  Value       // offending value
	fname string      // source filename
	ln    string      // source line number
	pname string      // procedure name
	args  []Value     // procedure arguments
}

//  Run wraps a Goaldi procedure in an environment and an exception catcher,
//  and calls it from Go
func Run(p Value, arglist []Value) {
	env := &Env{}
	defer func() {
		if x := recover(); x != nil {
			r := Diagnose(os.Stderr, x) // write Goaldi stack trace
			if !r {                     // if unrecognized
				fmt.Fprintf(os.Stderr, "Go stack:\n%s\n",
					debug.Stack()) // write Go stack trace
			}
			os.Exit(1)
		}
	}()
	p.(ICall).Call(env, arglist...)
}

//  Catch annotates a caught panic value with traceback information
func Catch(p interface{}, ev Value, fname string, ln string,
	procname string, arglist []Value) *CallFrame {
	return &CallFrame{p, ev, fname, ln, procname, arglist}
}

//  Diagnose handles traceback for a panic caught by Run()
//  It returns true for an "expected" (recognized) error.
func Diagnose(f io.Writer, v Value) bool {
	switch x := v.(type) {
	case *CallFrame:
		rv := Diagnose(f, x.cause)
		if _, ok := x.cause.(*runtime.TypeAssertionError); ok {
			fmt.Fprintf(f, "Offending value: %v\n", x.offv)
		}
		fmt.Fprintf(f, "Called by %s(%v) at %s line %s\n",
			x.pname, x.args, x.fname, x.ln)
		return rv
	case *RunErr:
		fmt.Fprintln(f, x.msg)
		if x.offv != nil {
			fmt.Fprintf(f, "Offending value: %v\n", x.offv)
		}
		return true
	case *runtime.TypeAssertionError:
		s := fmt.Sprintf("%#v", x)
		conc := extract(s, "concreteString")
		asst := extract(s, "assertedString")
		fmt.Fprintf(f, "Type %s does not implement %s\n",
			conc, asst)
		return true
	default:
		fmt.Fprintf(f, "%#v\n", x)
		return false
	}
}

//  extract finds a field in the %#v image of a struct
func extract(s string, label string) string {
	label = label + `:"`
	i := strings.Index(s, label)
	if i >= 0 {
		s = s[i+len(label) : len(s)]
		j := strings.Index(s, `"`)
		return s[0:j]
	} else {
		return "[?]"
	}
}

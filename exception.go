//  exception.go -- things dealing with exceptions and panics

package goaldi

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
)

//  RunErr records a Goaldi runtime error
type RunErr struct {
	Msg  string      // explanatory message
	Offv interface{} // offending value (Goaldi or Go value)
}

//  RunErr.String() returns a string form of a RunErr
func (e *RunErr) String() string {
	return fmt.Sprintf("RunErr: %s (%v)", e.Msg, e.Offv)
}

//  RunErr.Error() implements the interface that makes a RunErr a Go "error"
func (e *RunErr) Error() string {
	return e.String()
}

//  CallFrame records one frame of traceback information
type CallFrame struct {
	cause interface{} // underlying panic call
	offv  Value       // offending value
	coord string      // source coords (file:line:colm)
	pname string      // procedure name
	args  []Value     // procedure arguments
}

//  Cause(x) returns the original panic underlying a chain of CallFrame structs.
func Cause(x interface{}) interface{} {
	for {
		if f, ok := x.(*CallFrame); ok {
			x = f.cause
		} else {
			return x
		}
	}
}

//  Run wraps a Goaldi procedure in an environment and an exception catcher,
//  and calls it from Go
func Run(p Value, arglist []Value) {
	env := NewEnv(nil)
	defer Catcher(env)
	p.(ICall).Call(env, arglist...)
}

//  Catcher(env) tries to recover from a panic and print a traceback.
func Catcher(env *Env) {
	if x := recover(); x != nil {
		Diagnose(os.Stderr, x)            // write Goaldi stack trace
		if env.VarMap["gostack"] != nil { // if interpreter set %gostack
			fmt.Fprintf(os.Stderr, "Go stack:\n%s\n",
				debug.Stack()) // write Go stack trace
		}
		Shutdown(1)
		panic(x)
	}
}

//  Shutdown terminates execution with the given exit code.
func Shutdown(e int) {
	if f, ok := STDOUT.(*VFile); ok {
		f.Flush()
	}
	if f, ok := STDERR.(*VFile); ok {
		f.Flush()
	}
	pprof.StopCPUProfile()
	os.Exit(e)
}

//  Traceback is called as a deferred function to catch and annotate a panic
func Traceback(procname string, arglist []Value) {
	if p := recover(); p != nil {
		panic(Catch(p, nil, "", procname, arglist))
	}
}

//  Catch annotates a caught panic value with traceback information
func Catch(p interface{}, ev Value, coord string,
	procname string, arglist []Value) *CallFrame {
	return &CallFrame{p, ev, coord, procname, arglist}
}

//  Diagnose handles traceback for a panic caught by Run()
//  It returns true for an "expected" (recognized) error.
func Diagnose(f io.Writer, v interface{}) bool {
	switch x := v.(type) {
	case *CallFrame:
		rv := Diagnose(f, x.cause)
		if _, ok := x.cause.(*runtime.TypeAssertionError); ok {
			fmt.Fprintf(f, "Offending value: %#v\n", x.offv)
		}
		fmt.Fprintf(f, "Called by %s(", x.pname)
		for i, a := range x.args {
			if i > 0 {
				fmt.Fprintf(f, ",")
			}
			fmt.Fprintf(f, "%#v", a)
		}
		if x.coord != "" {
			fmt.Fprintf(f, ") at %s\n", x.coord)
		} else {
			fmt.Fprintf(f, ")\n")
		}
		return rv
	case *RunErr:
		fmt.Fprintln(f, x.Msg)
		if x.Offv != nil {
			fmt.Fprintf(f, "Offending value: %#v\n", x.Offv)
		}
		return true
	case *runtime.TypeAssertionError:
		s := fmt.Sprintf("%#v", x)
		conc := extract(s, "concreteString")
		asst := extract(s, "assertedString")
		fmt.Fprintf(f, "Type %s does not implement %s\n", conc, asst)
		return true
	case string:
		fmt.Fprintf(f, "PANIC: %v\n", x)
		return false
	default:
		fmt.Fprintf(f, "%T: %v\n", x, x)
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

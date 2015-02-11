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

//  Exception records a Goaldi panic value
type Exception struct {
	Msg  string  // explanatory message
	Offv []Value // offending values (Goaldi or Go values)
}

//  Exception.Error(), by its existence, makes an Exception a Go "error"
func (e *Exception) Error() string {
	return e.String()
}

//  Exception.String() returns a string form of a Exception
func (e *Exception) String() string {
	s := fmt.Sprintf("Exception(%#v", e.Msg)
	for _, v := range e.Offv {
		s = fmt.Sprintf("%s,%#v", s, v)
	}
	return s + ")"
}

//  Exception.GoString() converts an exception for image() or printf(%#v)
func (e *Exception) GoString() string {
	return e.String()
}

//  NewExn(s,v,...) creates and returns an Exception struct
func NewExn(s string, v ...Value) *Exception {
	return &Exception{s, v}
}

//  A Malfunction indicates an internal Goaldi problem (vs. a user error)
type Malfunction string

//  Malfunction.String() returns the default string representation.
func (e Malfunction) String() string {
	return "Malfunction: " + string(e)
}

//  Malfunction.Error() makes a Malfunction a Go "error"
func (e Malfunction) Error() string {
	return e.String()
}

//  CallFrame records one frame of traceback information
type CallFrame struct {
	cause interface{} // underlying panic call
	offv  []Value     // offending value
	coord string      // source coords (file:line:colm)
	pname string      // procedure name
	args  []Value     // procedure arguments
}

//  Traceback is called as a deferred function to catch and annotate a panic
func Traceback(procname string, arglist []Value) {
	if p := recover(); p != nil {
		panic(Catch(p, []Value{}, "", procname, arglist))
	}
}

//  Catch annotates a caught panic value with traceback information
func Catch(p interface{}, ev []Value, coord string,
	procname string, arglist []Value) *CallFrame {
	return &CallFrame{p, ev, coord, procname, arglist}
}

//  Cause(x) returns the original panic underlying a chain of CallFrame structs.
//  This is the value passed to an exception catcher.
func Cause(x interface{}) interface{} {
	for {
		if f, ok := x.(*CallFrame); ok {
			x = f.cause
		} else {
			return x
		}
	}
}

//  Catcher(env) prints a tracepback after a panic.
//  This is the recovery procedure at the top of the main (or coexpr) stack.
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

//  Diagnose prints traceback of a panic.
//  It returns true for an "expected" (recognized) error.
func Diagnose(f io.Writer, v interface{}) bool {
	switch x := v.(type) {
	case *CallFrame:
		rv := Diagnose(f, x.cause)
		if _, ok := x.cause.(*runtime.TypeAssertionError); ok {
			for _, v := range x.offv {
				fmt.Fprintf(f, "Offending value: %#v\n", v)
			}
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
	case *Exception:
		fmt.Fprintln(f, x.Msg)
		for _, v := range x.Offv {
			fmt.Fprintf(f, "Offending value: %#v\n", v)
		}
		return true
	case *runtime.TypeAssertionError:
		s := fmt.Sprintf("%#v", x)
		conc := extract(s, "concreteString")
		asst := extract(s, "assertedString")
		fmt.Fprintf(f, "Type %s does not implement %s\n", conc, asst)
		return true
	case Malfunction:
		fmt.Fprintf(f, "Goaldi runtime malfunction: %s\n", string(x))
		return false
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

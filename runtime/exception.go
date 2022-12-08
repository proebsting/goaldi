//  exception.go -- things dealing with exceptions and panics

package runtime

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

// Exception records a Goaldi panic value
type Exception struct {
	Msg  string  // explanatory message
	Offv []Value // offending values (Goaldi or Go values)
}

// Exception.Error(), by its existence, makes an Exception a Go "error"
func (e *Exception) Error() string {
	return e.String()
}

// Exception.String() returns a string form of a Exception
func (e *Exception) String() string {
	s := fmt.Sprintf("Exception(%#v", e.Msg)
	for _, v := range e.Offv {
		s = fmt.Sprintf("%s,%#v", s, v)
	}
	return s + ")"
}

// Exception.GoString() converts an exception for image() or printf(%#v)
func (e *Exception) GoString() string {
	return e.String()
}

// NewExn(s,v,...) creates and returns an Exception struct
func NewExn(s string, v ...Value) *Exception {
	return &Exception{s, v}
}

// A Malfunction indicates an internal Goaldi problem (vs. a user error)
type Malfunction string

// Malfunction.String() returns the default string representation.
func (e Malfunction) String() string {
	return "Malfunction: " + string(e)
}

// Malfunction.Error() makes a Malfunction a Go "error"
func (e Malfunction) Error() string {
	return e.String()
}

// CallFrame records one frame of traceback information
type CallFrame struct {
	cause interface{} // underlying panic call
	offv  []Value     // offending value
	coord string      // source coords (file:line:colm)
	pname string      // procedure name
	args  []Value     // procedure arguments
}

// Traceback is called as a deferred function to catch and annotate a panic
func Traceback(procname string, arglist []Value) {
	if p := recover(); p != nil {
		panic(Catch(p, []Value{}, "", procname, arglist))
	}
}

// Catch annotates a caught panic value with traceback information
func Catch(p interface{}, ev []Value, coord string,
	procname string, arglist []Value) *CallFrame {
	if te, ok := p.(*runtime.TypeAssertionError); ok {
		p = (*TypeError)(te)
	}
	return &CallFrame{p, ev, coord, procname, arglist}
}

// Cause(x) returns the original panic underlying a chain of CallFrame structs.
// This is the value passed to an exception catcher.
func Cause(x interface{}) interface{} {
	for {
		if f, ok := x.(*CallFrame); ok {
			x = f.cause
		} else {
			return x
		}
	}
}

// Catcher(env) prints a traceback after a panic.
// This is the recovery procedure at the top of the main (or coexpr) stack.
func Catcher(env *Env) {
	if x := recover(); x != nil {
		Diagnose(os.Stderr, x)                       // write Goaldi stack trace
		if env.Lookup("gostack", true) != NilValue { // if interpr set %gostack
			fmt.Fprintf(os.Stderr, "Go stack:\n%s\n",
				debug.Stack()) // write Go stack trace
		}
		Shutdown(1)
		panic(x)
	}
}

// Diagnose prints traceback of a panic.
// It returns true for an "expected" (recognized) error.
func Diagnose(f io.Writer, v interface{}) bool {
	switch x := v.(type) {
	case *CallFrame:
		rv := Diagnose(f, x.cause)
		if _, ok := x.cause.(*TypeError); ok {
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
	case *TypeError:
		fmt.Fprintln(f, x.Cleanup())
		return true
	case *runtime.TypeAssertionError:
		fmt.Fprintln(f, (*TypeError)(x).Cleanup())
		return true
	case Malfunction:
		fmt.Fprintf(f, "Goaldi runtime malfunction: %s\n", string(x))
		return false
	case string:
		fmt.Fprintf(f, "PANIC: %v\n", x)
		return false
	default:
		fmt.Fprintf(f, "PANIC(%T): %v\n", x, x)
		return false
	}
}

// A TypeError wraps a Go TypeAssertionError so we can change how it prints.
type TypeError runtime.TypeAssertionError

func (e *TypeError) Error() string {
	return `TypeError("` + e.Cleanup() + `")`
}

// Cleanup() simplifies the underlying Go runtime.TypeAssertionError.
// (This would be a lot easier if the error object fields weren't protected.)
func (e *TypeError) Cleanup() string {
	errstr := ((*runtime.TypeAssertionError)(e)).Error()
	subj := extract(errstr, "conversion: ")
	itis := extract(errstr, " is ")
	isnot := extract(errstr, " not ")
	switch isnot {
	case "IVariable":
		return "Variable expected"
	case "Numerable":
		return "Number expected"
	case "Stringable":
		return "String expected"
	default:
		if itis != "not" { // i.e. "e is t" not "e is not ..."
			return fmt.Sprintf("%s is not %s", itis, isnot)
		} else {
			return fmt.Sprintf("%s does not implement %s", subj, isnot)
		}
	}
}

// extract finds the field following a given indicator prefix
// and cleans it up, removing any further [*][runtime.[V]] prefix
func extract(s string, prefix string) string {
	i := strings.Index(s, prefix)
	if i < 0 {
		return ""
	}
	s = s[i+len(prefix) : len(s)]
	j := strings.IndexAny(s, ", :")
	if j > 0 {
		s = s[0:j]
	}
	if strings.HasPrefix(s, "*") {
		s = s[1:]
	}
	if strings.HasPrefix(s, "runtime.V") {
		s = s[9:]
	} else if strings.HasPrefix(s, "runtime.") {
		s = s[8:]
	}
	return s
}

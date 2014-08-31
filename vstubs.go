//  vstubs.go -- a collection of stub methods that panic
//
//  A struct type that includes a Stubs field, and additionally implements
//  the String() method, effectively implements the Value interface.
//  These stub functions act as defaults if not overridden by
//  implementations in the enclosing struct.  Each one panics.
//
//  Note that a stub can't access its enclosing struct.
//  This means that a default stub can't do anything *useful*.

package goaldi

import (
	"fmt"
	"runtime"
	"strings"
)

type Stubs struct { // the Stubs struct itself is empty
}

//  --------------- stub functions --------------

// The String() method, used by fmt.Printf, is deliberately not implemented.
// Every Value interface *must* supply at least that one method.

func (p *Stubs) Deref() Value       { no(); return nil }
func (p *Stubs) AsString() *VString { no(); return nil }
func (p *Stubs) AsNumber() *VNumber { no(); return nil }

func (p *Stubs) Add(v2 Value) (Value, *Closure)  { return no() }
func (p *Stubs) Mult(v2 Value) (Value, *Closure) { return no() }

//  --------------- validation --------------

type stubsplus struct {
	Stubs
}

func (p *stubsplus) String() string { return "stubsplus"}

var _ Value = &stubsplus{}	// if error, stub collection is incomplete

//  --------------- support functions --------------

//  no() panics with a message about a missing function
//  and a source reference identifying the caller's caller.
//
//  Assuming the information is available, the message takes the form (e.g.)
//	No Divide() function for supplied argument type (sourcefile.go+29)
//
//  It would be really nice if we could include the offending value and type,
//  but that information is not available.
//  #%#% We could however include more traceback information if warranted.
func no() (Value, *Closure) {
	msg := "No "
	if pc, _, _, ok := runtime.Caller(1); ok {
		if f := runtime.FuncForPC(pc); f != nil {
			msg = msg + tail(f.Name(), ".") + "() "
		}
	}
	msg = msg + "function for supplied argument type"
	if _, file, line, ok := runtime.Caller(2); ok {
		msg = fmt.Sprintf("%s (%s+%d)", msg, tail(file, "/"), line)
	}
	panic(msg)
}

//  tail returns the last field after a member of the separator set.
func tail(s string, seps string) string {
	if i := strings.LastIndex(s, seps); i > 0 {
		s = s[i+1 : len(s)]
	}
	return s
}

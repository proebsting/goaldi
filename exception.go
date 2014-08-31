//  exception.go -- things dealing with exceptions

//  NOTE:  Not all of the runtime code currently checks for exceptions.
//  It adds clutter, and we may end up doing this differently via panics.

package goaldi

import (
	"fmt"
)

var CATCHME *Closure = &Closure{} // special flag value for exceptions

//  NewException constructs an exception from a format string and arguments
func NewException(format string, args ...interface{}) Value {
	return NewString(fmt.Sprintf(format, args...))
}

//  Throw returns a simple Goaldi value as an exception
func Throw(v Value) (Value, *Closure) {
	return v, CATCHME
}

//  Throwf formats and throws an exception a la printf.
func Throwf(format string, args ...interface{}) (Value, *Closure) {
	return Throw(NewException(format, args...))
}

//  Run wraps a Goaldi procedure in an exception catcher, and calls it from Go
func Run(p Procedure) {
	fmt.Println("[--------------------- begin ---------------------]")

//	defer func() {
//		// this works, but get more detailed traceback without it
//		if x := recover(); x != nil {
//			fmt.Println("PANIC:", x)
//		}
//	}()

	t1, c1 := p(nil)
	if c1 == CATCHME {
		fmt.Println("UNCAUGHT EXCEPTION: ", t1)
	} else if t1 == nil {
		fmt.Println("[failed]")
	} else {
		fmt.Printf("[returned %s]\n", t1)
	}
}

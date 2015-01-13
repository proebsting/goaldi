//  calling.go -- procedure call / suspension / resumption interface
//
//  In the Go implementation, a Goaldi procedure or operation result
//  is the pair (Value, *Closure) with this meaning:
//
//  Value  *Closure   Interpretation
//  -----  -------   ---------------------------------
//   nil     nil     procedure failed
//  value    nil     procedure returned a value
//  value   resume   procedure suspended and can be resumed

package goaldi

import (
	"fmt"
)

var _ = fmt.Printf // enable debugging

type Closure struct {
	Go Resumer // start or resumption function, depending on context
}

//  Procedure function prototype
type Procedure func(env *Env, args ...Value) (Value, *Closure)

//  Procedure resumption prototype
type Resumer func() (Value, *Closure)

//  ICall interface
type ICall interface {
	Call(env *Env, args []Value, names []string) (Value, *Closure)
}

//  ProcArg(a,i,d)  returns procedure argument a[i], defaulting to d
func ProcArg(a []Value, i int, d Value) Value {
	if i < len(a) && a[i] != NilValue {
		return a[i]
	} else {
		return d
	}
}

//  Resume() executes the entry point in a Closure to produce the next result.
//  If the pointer is nil, failure is produced.
func (c *Closure) Resume() (Value, *Closure) {
	if c == nil {
		return Fail()
	}
	return c.Go()
}

//  Fail returns a failure indicator
func Fail() (Value, *Closure) {
	return nil, nil
}

//  Return returns a simple value as a duo
func Return(v Value) (Value, *Closure) {
	return v, nil
}

//  ArgNames handles named arguments by building a new arglist.
//  offv is given as the "offending value" in case of errors.
//  The pnames value may be nil to indicate no param names are known.
//#%#% DOES NOT HANDLE VARARGS
func ArgNames(args []Value, names []string, offv Value, pnames *[]string) []Value {
	if len(names) == 0 {
		return args
	}
	if pnames == nil {
		panic(NewExn("Named arguments not allowed", offv))
	}

	// make a list of target indexes for storing the named arguments seen
	locs := make([]int, len(names)) // list of indexes
	nslots := 0                     // totall number of parameters to pass
	for i, s := range names {
		j := argIndex(s, pnames) // get index of name i
		locs[i] = j              // save it
		if nslots <= j {
			nslots = j + 1
		}
	}

	// make a new argument list of sufficient size
	// in which a Go nil indicates an unfilled slot
	newargs := make([]Value, nslots, nslots)
	nbase := len(args) - len(names) // base of named arguments
	copy(newargs, args[:nbase])     // copy in the unnamed arguments

	// copy in the named arguments
	for i, j := range locs {
		if newargs[j] != nil {
			panic(NewExn("Duplicate argument", names[i]))
		}
		newargs[j] = args[nbase+i]
	}

	// fill unused slots with Goaldi nils
	for i := nbase; i < len(newargs); i++ {
		if newargs[i] == nil {
			newargs[i] = NilValue
		}
	}
	return newargs
}

//  argIndex finds the index of an argument name in a list of strings
func argIndex(name string, pnames *[]string) int {
	for i, s := range *pnames {
		if s == name {
			return i
		}
	}
	panic(NewExn("No parameter matches name", name))
}

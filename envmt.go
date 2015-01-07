//  envmt.go -- dynamic variables and procedure environment

package goaldi

import (
	"fmt"
	"io"
)

//  execution environment
//
//  #%#% This needs more thought, especially if it is to be dynamic.
//  #%#% Currently it only changes, by copying, on creation of a new thread.
type Env struct {
	ThreadID int              // thread ID
	VarMap   map[string]Value // %variable map
	//#%#% more to be determined
	//#%#% dynamic variables?
}

//  NewEnv(e) returns a new environment with a distinct ThreadID.
func NewEnv(e *Env) *Env {
	if e == nil {
		return &Env{0, StdEnv}
	} else {
		return &Env{<-TID, e.VarMap}
	}
}

//  ThreadID production
var TID = make(chan int)

func init() {
	go func() {
		tid := 0
		for {
			tid++
			TID <- tid
		}
	}()
}

//  StdEnv is the initial environment
var StdEnv = make(map[string]Value)

//  EnvInit registers a standard environment value or variable at init time.
//  (Variables should be registered as trapped values).
func EnvInit(name string, v Value) {
	StdEnv[name] = v
}

//  Initial values and variables
func init() {

	// math constants
	EnvInit("e", E)
	EnvInit("phi", PHI)
	EnvInit("pi", PI)

	// standard files (mutable)
	EnvInit("stdin", Trapped(&STDIN))
	EnvInit("stdout", Trapped(&STDOUT))
	EnvInit("stderr", Trapped(&STDERR))

	// error recovery
	EnvInit("error", Trapped(NewVariable(NilValue)))
}

//	ShowEnvironment(f) -- list standard environment on file f
func ShowEnvironment(f io.Writer) {
	fmt.Fprintln(f)
	fmt.Fprintln(f, "Standard Environment")
	fmt.Fprintln(f, "------------------------------")
	for k := range SortedKeys(StdEnv) {
		cv := "c"
		v := StdEnv[k]
		if t, ok := v.(*VTrapped); ok {
			cv = "v"
			v = t.Deref()
		}
		fmt.Fprintf(f, "%%%-8s %s  %#v\n", k, cv, v)
	}
}

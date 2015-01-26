//  run.go -- overall control of execution

package goaldi

import (
	"os"
	"runtime/pprof"
)

//  An InitItem is a global initialization procedure with dependencies
type InitItem struct {
	Proc *VProcedure // procedure to execute
	Uses []string    // variables used by this procedure
	Sets string      // variable set by running this
}

//  Run wraps a Goaldi procedure in an environment and an exception catcher,
//  and calls it from Go
func Run(p Value, arglist []Value) {
	env := NewEnv(nil)
	defer Catcher(env)
	p.(ICall).Call(env, arglist, []string{})
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

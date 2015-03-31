//  run.go -- overall control of execution

package runtime

import (
	"os"
	"runtime/pprof"
)

//  Run wraps a Goaldi procedure in an environment and an exception catcher,
//  and calls it from Go.
//  This is used first for any initialization blocks and then for main().
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

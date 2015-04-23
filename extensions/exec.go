//  exec.go -- extensions for executing commands
//
//  These functions are used by the Goaldi translator
//  to restart the interpreter after it has compiled a program.

package extensions

import (
	g "goaldi/runtime"
	"os"
)

//  declare procedures
func init() {
	g.GoLib(OSArgs, "osargs", "", "get program argument vector")
	g.DefLib(OSFile, "osfile", "fd,name", "get Go os.File for FD")
}

//  osargs() returns the program argument vector, argv.
//  This is the argument vector as seen by the Goaldi interpreter,
//  and is a superset of the vector passed to the Goaldi program.
func OSArgs() []string {
	return os.Args
}

//  osfile(fd,name) returns the Go os.File struct for a file descriptor.
//  This bypasses the normal import mechanism that would make a Goaldi file.
func OSFile(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
	defer g.Traceback("osfile", args)
	i := g.IntVal(g.ProcArg(args, 0, g.ZERO))
	s := g.ToString(g.ProcArg(args, 1, g.EMPTY)).ToUTF8()
	switch i {
	case 0:
		return os.Stdin, nil
	case 1:
		return os.Stdout, nil
	case 2:
		return os.Stderr, nil
	default:
		return os.NewFile(uintptr(i), s), nil
	}
}

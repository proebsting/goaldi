//  run.go -- overall control of execution

package goaldi

import (
	"errors"
	"fmt"
	"os"
	"runtime/pprof"
)

//  An InitItem is a global initialization procedure with dependencies
type InitItem struct {
	Proc     *VProcedure // procedure to execute
	Uses     []string    // variables used by this procedure
	Sets     string      // variable set by running this procedure
	pending  int         // number of others we are waiting on
	releases []int       // list of others to notify on set
}

//  NewInit creates an InitItem
func NewInit(proc *VProcedure, uses []string, sets string) *InitItem {
	return &InitItem{proc, uses, sets, 0, nil}
}

//  Run wraps a Goaldi procedure in an environment and an exception catcher,
//  and calls it from Go.
//  This is used first for any initial{} blocks and then for main().
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

//  RunDep runs a set of procedures in dependency order.
//  This is used for initializing globals.
//  Execution errors are handled by the usual exception handling.
//  RunDep returns an error if circular dependencies remain at the end.
func RunDep(ilist []*InitItem, trace bool) error {

	// make a table of all the globals of interest
	// (we don't care about any global that is not *set* in this list)
	itable := make(map[string]int)
	for i, item := range ilist {
		itable[item.Sets] = i + 1 // store index+1 so that 0 means not found
		item.pending = 0
		item.releases = make([]int, 0)
	}

	// for each item, count the number of others it depends on,
	// and register the item on the list of each of those others
	// for notification when set.
	for i, item := range ilist { // for each init item
		for _, id := range item.Uses { // for each reference listed
			j := itable[id] - 1   // look up in table
			if j >= 0 && j != i { // if registered, and if not self
				item.pending++ // increment wait count
				guard := ilist[j]
				guard.releases = append(guard.releases, i) // register
			}
		}
	}

	// if tracing, show initial data structures
	if trace {
		for _, item := range ilist {
			fmt.Printf("global %s depends on [", item.Sets)
			for _, s := range item.Uses {
				fmt.Printf("%s,", s)
			}
			fmt.Printf("] used by [")
			for _, j := range item.releases {
				fmt.Printf("%s,", ilist[j].Sets)
			}
			fmt.Println("]")
		}
	}

	todo := len(ilist) // number of procedures to run
	trynext := 0       // next one to run, if ready

	// loop until we reach the end of the list, running procedures
	// -- skip items not yet ready to run
	// -- go back to earliest one when running something makes it ready
	for trynext < len(ilist) {
		item := ilist[trynext] // get next potential candidate
		trynext++              // and bump the pointer
		if item.pending == 0 { // if this one is ready to run
			if trace {
				fmt.Printf("global %s initializing:\n", item.Sets)
			}
			Run(item.Proc, []Value{}) // run the procedure
			item.pending--            // mark it as done
			todo--                    // count it
			// decrement the wait count of each dependent item
			for _, j := range item.releases {
				ilist[j].pending--
				// if this item is now ready to run,
				// make it next if it precedes the currently chosen one
				if ilist[j].pending == 0 && j < trynext {
					trynext = j
				}
			}
		} else if trace {
			fmt.Printf("global %s wait count = %d\n", item.Sets, item.pending)
		}
	}

	if todo == 0 {
		return nil // success
	}

	// there was a circular dependency; report an error
	s := "Circular dependency among:"
	for _, item := range ilist {
		if item.pending > 0 {
			s = s + " " + item.Sets
		}
	}
	return errors.New(s)
}

//  run.go -- overall control of execution

package goaldi

import (
	"errors"
	"fmt"
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

//  An InitItem is a global initialization procedure with dependencies
type InitItem struct {
	Proc     *VProcedure // initialization procedure for globals (only!)
	Uses     []string    // variables used by this global or procedure
	Sets     string      // name of procedure or associated global
	status   int         // current status
	passnum  int         // last visit time
	awaiting *InitItem   // one item that blocks this one
}

//  InitItem status values
//  The order matters here -- see InitItem.setStatus()
const (
	initUnk   = iota // never calculated
	waitGlob         // waiting on at least one global
	waitProc         // waiting only on procedures
	initReady        // ready to initialize
	initDone         // initialization done (or scheduled)
)

//  NewInit creates an InitItem
func NewInit(proc *VProcedure, uses []string, sets string) *InitItem {
	return &InitItem{proc, uses, sets, initUnk, 0, nil}
}

//  table of all items being initialized (in global init pass)
var itemTable = make(map[string]*InitItem)

//  RunDep runs a set of procedures in dependency order.
//  This is used for initializing globals.
//  Execution errors are handled by the usual exception handling.
//  RunDep returns an error if circular dependencies remain at the end,
//  or if an attempt is made to set the same global twice.
func RunDep(ilist []*InitItem, trace bool) error {
	// initialize table of items
	for _, item := range ilist {
		if itemTable[item.Sets] != nil {
			return errors.New("Multiple initializations of global: " + item.Sets)
		}
		if trace {
			fmt.Printf("[-] init %s: depends on %s\n", item.Sets, item.Uses)
		}
		itemTable[item.Sets] = item
	}
	passnum := 0
	// restart from beginning of list each time to preserve lexical ordering
OuterLoop:
	for {
		// remove completed items from front of list
		for len(ilist) > 0 && ilist[0].status == initDone {
			ilist = ilist[1:]
		}
		// if nothing is left, we are done
		if len(ilist) == 0 {
			if trace {
				fmt.Printf("[-] global variable initialization complete\n")
			}
			return nil // success
		}
		// make a pass through all items looking for something runnable
		passnum++
		for _, item := range ilist {
			item.setStatus(passnum)
			if item.status == initReady {
				// found something; run it and mark it
				if trace {
					if item.Proc != nil { // if a global initializer
						fmt.Printf("[-] global %s initializing:\n", item.Sets)
					} else {
						fmt.Printf("[-] procedure %s ready\n", item.Sets)
					}
				}
				if item.Proc != nil { // if a global initializer
					Run(item.Proc, []Value{}) // run the procedure
				}
				item.status = initDone
				continue OuterLoop
			} else if item.status != initDone && trace {
				fmt.Printf("[-] %s waiting on %s\n",
					item.Sets, item.awaiting.Sets)
			}
		}
		// didn't find anything but list is not empty; this is an error
		s := "Circular dependency among:"
		for _, item := range ilist {
			if item.status == waitGlob {
				s = s + " " + item.Sets
			}
		}
		return errors.New(s)
	}
}

//  InitItem.setSatus computes and returns the status for the current pass
func (m *InitItem) setStatus(passnum int) int {
	// if m is nil, this entry isn't even in the table and is considered done
	if m == nil {
		return initDone
	}
	// if already visited this pass, or if already ready, bail out now
	if m.passnum == passnum || m.status >= initReady {
		return m.status
	}
	m.passnum = passnum  // note this visit to break recursion
	m.status = initReady // assume ready unless we learn otherwise
	m.awaiting = nil
	// check all the other items on which this one depends
	for _, u := range m.Uses {
		o := itemTable[u]
		if o == m {
			continue // don't wait on self
		}
		s := o.setStatus(passnum)
		if s != initDone { // if we need to wait for this
			m.awaiting = o
			if o.Proc != nil { // if this is a global
				m.status = waitGlob
			} else if m.status != waitGlob {
				m.status = waitProc
			}
		}
	}
	// if this is a procedure, not a global, don't wait for other procs
	// (this permits circularity within procedure calls only)
	if m.status == waitProc && m.Proc == nil {
		m.status = initReady
	}
	return m.status
}

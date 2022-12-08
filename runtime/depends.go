//  depends.go -- initialization dependency processing

package runtime

import (
	"errors"
	"fmt"
)

// A DependencyList holds a collection of globals and procedures for ordering
type DependencyList struct {
	list    []*InitItem          // ordered list of entries
	table   map[string]*InitItem // entries indexed by name
	passnum int                  // pass number during processing
}

// An InitItem is a global initialization procedure with dependencies
type InitItem struct {
	proc     *VProcedure // initialization procedure for globals (only!)
	uses     []string    // variables used by this global or procedure
	sets     string      // name of procedure or associated global
	status   int         // current status
	passnum  int         // last visit time
	awaiting *InitItem   // one item that blocks this one
}

// InitItem status values
// The order matters here -- see InitItem.setStatus()
const (
	initUnk   = iota // never calculated
	waitGlob         // waiting on at least one global
	waitProc         // waiting only on procedures
	initReady        // ready to initialize
	initDone         // initialization done (or scheduled)
)

// DependencyList.Add inserts an item in the list.
// This ordering is preserved in the absence of actual dependencies.
func (dl *DependencyList) Add(
	name string, initProc *VProcedure, uses []string) {
	item := &InitItem{initProc, uses, name, initUnk, 0, nil}
	dl.list = append(dl.list, item)
}

// DependencyList.Reorder places the procedures in dependency order.
// This is used for initializing globals.
// Reorder returns an error if circular dependencies remain at the end,
// or if an attempt is made to set the same global twice.
func (dl *DependencyList) Reorder(trace bool) error {
	if trace {
		fmt.Printf("[-] begin dependency computation\n")
	}
	ilist := dl.list
	runlist := make([]*InitItem, 0, len(dl.list))
	// initialize table of items
	dl.table = make(map[string]*InitItem)
	for _, item := range ilist {
		if dl.table[item.sets] != nil {
			return errors.New("Multiple initializations of global: " + item.sets)
		}
		if trace {
			fmt.Printf("[-]    %s depends on %s\n", item.sets, item.uses)
		}
		dl.table[item.sets] = item
	}
	// restart from beginning of list each time to preserve lexical ordering
OuterLoop:
	for {
		// remove completed items from front of list
		for len(ilist) > 0 && ilist[0].status == initDone {
			ilist = ilist[1:]
		}
		// if nothing is left, we are done
		if len(ilist) == 0 {
			dl.list = runlist
			if trace {
				fmt.Printf("[-] dependency reordering complete\n")
			}
			return nil // success
		}
		// make a pass through all items looking for something runnable
		dl.passnum++
		for _, item := range ilist {
			item.setStatus(dl)
			if item.status == initReady {
				// found something; run it and mark it
				if trace {
					if item.proc != nil { // if a global initializer
						fmt.Printf("[-] global %s scheduled\n", item.sets)
					} else {
						fmt.Printf("[-] procedure %s ready\n", item.sets)
					}
				}
				runlist = append(runlist, item)
				item.status = initDone
				continue OuterLoop
			} else if item.status != initDone && trace {
				fmt.Printf("[-] %s awaits %s\n",
					item.sets, item.awaiting.sets)
			}
		}
		// didn't find anything but list is not empty; this is an error
		s := "Circular dependency among:"
		for _, item := range ilist {
			if item.status == waitGlob {
				s = s + " " + item.sets
			}
		}
		return errors.New(s)
	}
}

// DependencyList.RunAll runs all the initializers in their current order.
// Execution errors are handled by the usual exception handling.
func (dl *DependencyList) RunAll() {
	for _, item := range dl.list {
		if item.proc != nil { // if a global initializer
			Run(item.proc, []Value{}) // run it
		}
	}
}

// InitItem.setSatus computes and returns the status for the current pass
func (m *InitItem) setStatus(dl *DependencyList) int {
	// if m is nil, this entry isn't even in the table and is considered done
	if m == nil {
		return initDone
	}
	// if already visited this pass, or if already ready, bail out now
	if m.passnum == dl.passnum || m.status >= initReady {
		return m.status
	}
	m.passnum = dl.passnum // note this visit to break recursion
	m.status = initReady   // assume ready unless we learn otherwise
	m.awaiting = nil
	// check all the other items on which this one depends
	for _, u := range m.uses {
		o := dl.table[u]
		if o == m {
			continue // don't wait on self
		}
		s := o.setStatus(dl)
		if s != initDone { // if we need to wait for this
			m.awaiting = o
			if o.proc != nil || o.status == waitGlob {
				m.status = waitGlob // waiting on a global, or indirectly so
			} else if m.status != waitGlob {
				m.status = waitProc // waiting on a procedure
			}
		}
	}
	// if this is a procedure, not a global, don't wait for other procs
	// (this permits circularity within procedure calls only)
	if m.status == waitProc && m.proc == nil {
		m.status = initReady
	}
	return m.status
}

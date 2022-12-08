//  interp.go -- interpret procedure

package main

import (
	"fmt"
	g "github.com/proebsting/goaldi/runtime"
)

// procedure frame
type pr_frame struct {
	env   *g.Env                 // dynamic execution environment
	info  *pr_Info               // static procedure information
	args  []g.Value              // arglist as called
	vars  map[string]interface{} // variables and scopes
	temps []interface{}          // temporaries
	coord string                 // last known source location
	offv  g.Value                // offending value for traceback
	cxout g.VChannel             // co-expression output pipe
	onerr *g.VProcedure          // recovery procedure
}

// newframe(f) -- duplicate a procedure frame for "create e"
func newframe(f *pr_frame) *pr_frame {
	fnew := &pr_frame{} // allocate new frame
	*fnew = *f          // duplicate values
	fnew.onerr = nil    // don't copy recovery procedure
	fnew.temps = make([]interface{}, len(f.temps))
	fnew.vars = make(map[string]interface{})
	for k, v := range f.vars {
		fnew.vars[k] = v
	}
	// make new copies of all parameter values
	for _, name := range f.info.params {
		fnew.vars[name] = g.NewVariable(g.Deref(f.vars[name]))
	}
	// make new copies of all locals (n.b. does not include statics)
	for _, name := range f.info.locals {
		fnew.vars[name] = g.NewVariable(g.Deref(f.vars[name]))
	}
	return fnew
}

// duplvars(a) -- duplicate a list of (trapped) variables or parameters
func duplvars(a []g.Value) []g.Value {
	b := make([]g.Value, len(a))
	for i, x := range a {
		b[i] = g.NewVariable(g.Deref(x))
	}
	return b
}

// interp -- interpret one procedure
func interp(env *g.Env, pr *pr_Info, outer map[string]interface{},
	args ...g.Value) (g.Value, *g.Closure) {

	if opt_trace {
		fmt.Printf("[%d] enter procedure %s\n", env.ThreadID, pr.qname)
	}

	// initialize procedure frame
	var f pr_frame
	f.env = env                                // environment
	f.info = pr                                // static procedure info
	f.args = args                              // argument list
	f.temps = make([]interface{}, 1+pr.ntemps) // temporaries

	// initialize variable dictionary with inherited variables;
	// any of these may be subsequently hidden (replaced)
	f.vars = make(map[string]interface{})
	for k, v := range outer {
		f.vars[k] = v
	}

	// store the initial inherited (blank) scope in the variable map
	f.vars[""] = env

	// our own locals and scopes are not defined here;
	// they are later added dynamically by Ir_EnterScope instructions

	// define static variables
	for k, v := range pr.statics {
		f.vars[k] = v
	}

	// initialize parameters
	for i, name := range pr.params {
		if i < len(args) {
			f.vars[name] = g.NewVariable(args[i])
		} else {
			f.vars[name] = g.NewVariable(g.NilValue)
		}
	}

	//  handle variadic procedure
	if pr.variadic {
		n := len(pr.params) - 1
		vp := new(g.Value)
		if len(args) < n {
			*vp = g.NewList(0, nil)
		} else {
			vals := make([]g.Value, len(args)-n)
			copy(vals, args[n:])
			*vp = g.InitList(vals)
		}
		f.vars[pr.params[n]] = g.Trapped(vp)
	}

	// execute the IR code
	return execute(&f, pr.ir.CodeStart)
}

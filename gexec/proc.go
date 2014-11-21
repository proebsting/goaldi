//  proc.go -- things dealing with procedures at link time

package main

import (
	"fmt"
	g "goaldi"
	"regexp"
)

//  information about a procedure that is shared by all invocations
type pr_Info struct {
	name     string                   // procedure name
	outer    *pr_Info                 // enclosing procedure, if any
	ir       *ir_Function             // intermediate code structure
	insns    map[string][]interface{} // map from labels to IR code chunks
	known    map[string]bool          // set of locally declared identifiers
	statics  map[string]interface{}   // table of statics (including globals)
	locals   []string                 // list of local names
	params   []string                 // list of parameter names
	variadic bool                     // true if last param is []
}

//  global index of procedure information
var ProcTable = make(map[string]*pr_Info)

//  pattern for extracting name of enclosing procedure
var enclpat = regexp.MustCompile(`^(.*)\$nested\$[0-9]*$`)

//  declareProc initializes and returns a procedure info structure
func declareProc(ir *ir_Function) *pr_Info {
	pr := &pr_Info{}
	pr.name = ir.Name
	pr.ir = ir
	pr.variadic = (ir.Accumulate != "")
	pr.params = pr.ir.ParamList
	pr.locals = pr.ir.LocalList
	pr.known = make(map[string]bool)
	for _, name := range pr.params {
		pr.known[name] = true
	}
	for _, name := range pr.locals {
		pr.known[name] = true
	}
	for _, name := range ir.StaticList {
		pr.known[name] = true
	}
	ProcTable[ir.Name] = pr
	// if nested, we also know all idents known to enclosing procedure
	matches := enclpat.FindStringSubmatch(pr.name) // check pattern of proc name
	if matches != nil {                            // if nested
		pr.outer = ProcTable[matches[1]]   // look up parent info
		for name := range pr.outer.known { // every identifier known there
			pr.known[name] = true // is known here, too
		}
	}
	return pr
}

//  irProcedure makes a runtime procedure from a procedure info structure
func irProcedure(pr *pr_Info) *g.VProcedure {
	return g.NewProcedure(pr.name,
		func(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
			return interp(env, pr, args...)
		})
}

//  setupProc finishes procedure setup now that the GlobalDict is set
//	report undeclared identifiers
//	create combined dictionary of global + local variables
//	create chunk table indexed by labels
//	#%#% TODO: handle nested procedures
func setupProc(pr *pr_Info) {
	undeclared(pr)
	pr.statics = makeDict(pr)
	pr.insns = getInsns(pr)
	if opt_verbose {
		fmt.Printf("\n%s()  %d param  %d local  static+global %d\n",
			pr.name, len(pr.params), len(pr.locals), len(pr.statics))
	}
}

//  makeDict creates the initial mapping of identifiers to variables for a proc
//  This initial dictionary contains only globals and statics.
func makeDict(pr *pr_Info) map[string]interface{} {
	dict := make(map[string]interface{})
	//  start with every global that is not hidden by a local/param/static
	for name, value := range GlobalDict {
		if !pr.known[name] {
			dict[name] = value
		}
	}
	// add a trapped variable for every static
	for _, name := range pr.ir.StaticList {
		dict[name] = g.Trapped(g.NewVariable(g.NilValue))
	}
	return dict
}

//  getInsns creates the index of IR code chunks
func getInsns(pr *pr_Info) map[string][]interface{} {
	insns := make(map[string][]interface{})
	for _, ch := range pr.ir.CodeList {
		insns[ch.Label] = ch.InsnList
	}
	return insns
}

//  undeclared reports identifiers never declared anywhere
func undeclared(pr *pr_Info) {
	for _, chunk := range pr.ir.CodeList {
		for _, insn := range chunk.InsnList {
			if i, ok := insn.(ir_Var); ok {
				if !pr.known[i.Name] && Undeclared[i.Name] {
					//%#% warn now, later fatal
					warning(fmt.Sprintf("%v %s undeclared",
						i.Coord, i.Name))
					// inhibit repeated warnings
					delete(Undeclared, i.Name)
				}
			}
		}
	}
}

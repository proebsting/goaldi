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

//  setupProc finishes procedure setup now that all globals are known
func setupProc(pr *pr_Info) {

	// report undeclared identifiers
	undeclared(pr)

	// make a trapped variable for every static
	pr.statics = make(map[string]interface{})
	for _, name := range pr.ir.StaticList {
		pr.statics[name] = g.Trapped(g.NewVariable(g.NilValue))
	}

	// create an index of IR code chunks
	pr.insns = make(map[string][]interface{})
	for _, ch := range pr.ir.CodeList {
		if pr.insns[ch.Label] != nil {
			panic("Duplicate IR label: " + ch.Label)
		}
		pr.insns[ch.Label] = ch.InsnList
	}

	if opt_verbose {
		fmt.Printf("\n%s()  %d param  %d local  %d static\n",
			pr.name, len(pr.params), len(pr.locals), len(pr.statics))
	}
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

//  irProcedure makes a runtime procedure from static info and inherited vars
func irProcedure(pr *pr_Info, outer map[string]interface{}) *g.VProcedure {
	return g.NewProcedure(pr.name,
		func(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
			return interp(env, pr, outer, args...)
		})
}

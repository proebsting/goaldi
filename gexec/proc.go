//  proc.go -- things dealing with procedures at link time

package main

import (
	"fmt"
	g "goaldi"
	"strings"
)

//  information about a procedure that is shared by all invocations
type pr_Info struct {
	name     string                   // procedure name
	outer    *pr_Info                 // enclosing procedure, if any
	ir       *ir_Function             // intermediate code structure
	insns    map[string][]interface{} // map from labels to IR code chunks
	statics  map[string]interface{}   // table of statics (including globals)
	locals   []string                 // list of local names
	params   []string                 // list of parameter names
	variadic bool                     // true if last param is []
}

//  global index of procedure information
var ProcTable = make(map[string]*pr_Info)

//  declareProc initializes and returns a procedure info structure
func declareProc(ir *ir_Function) *pr_Info {
	if ProcTable[ir.Name] != nil {
		fatal("duplicate procedure definition: " + ir.Name)
	}
	pr := &pr_Info{}
	pr.name = ir.Name
	pr.ir = ir
	pr.variadic = (ir.Accumulate != "")
	pr.params = pr.ir.ParamList
	pr.locals = pr.ir.LocalList
	ProcTable[ir.Name] = pr
	return pr
}

//  setupProc finishes procedure setup now that all globals are known
func setupProc(pr *pr_Info) {

	// report undeclared identifiers
	for _, id := range pr.ir.UnboundList {
		if GlobalDict[id] == nil {
			fatal("in " + pr.name + "(): undeclared identifier: " + id)
		}
	}

	// make a trapped variable for every static
	pr.statics = make(map[string]interface{})
	for _, name := range pr.ir.StaticList {
		pr.statics[name] = g.NewVariable(g.NilValue)
	}

	// create an index of IR code chunks
	pr.insns = make(map[string][]interface{})
	for _, ch := range pr.ir.CodeList {
		if pr.insns[ch.Label] != nil {
			panic(g.Malfunction("Duplicate IR label: " + ch.Label))
		}
		pr.insns[ch.Label] = ch.InsnList
	}

	if opt_verbose {
		fmt.Printf("\n%s()  %d param  %d local  %d static\n",
			pr.name, len(pr.params), len(pr.locals), len(pr.statics))
	}
}

//  irProcedure makes a runtime procedure from static info and inherited vars
func irProcedure(pr *pr_Info, outer map[string]interface{}) *g.VProcedure {

	// make a list of unadorned parameter names
	pnames := make([]string, len(pr.params))
	for i, s := range pr.params {
		pnames[i] = s[:strings.Index(s, ":")]
	}

	// copy (references to) any inherited variables
	vars := make(map[string]interface{})
	if outer != nil {
		for k, v := range outer {
			vars[k] = v
		}
	}

	return g.NewProcedure(pr.name, &pnames, pr.variadic,
		func(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
			return interp(env, pr, vars, args...)
		}, nil, "")
}

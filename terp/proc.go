//  proc.go -- things dealing with procedures in the interpreter

package main

import (
	"fmt"
	g "goaldi"
)

//  info about a procedure
type pr_Info struct {
	name  string                 // procedure name
	outer *pr_Info               // immediate parent #%#% NOT YET SET
	ir    *ir_Function           // intermediate code structure
	lset  map[string]bool        // set of locally declared identifiers
	dict  map[string]interface{} // map of identifiers to variables
}

//  global index of procedure information
var ProcTable = make(map[string]*pr_Info)

//  declareProc initializes and returns a procedure info structure
func declareProc(ir *ir_Function) *pr_Info {
	pr := &pr_Info{}
	pr.name = ir.Name
	pr.ir = ir
	pr.lset = make(map[string]bool)
	for _, name := range ir.ParamList {
		pr.lset[name] = true
	}
	for _, name := range ir.LocalList {
		pr.lset[name] = true
	}
	for _, name := range ir.StaticList {
		pr.lset[name] = true
	}
	ProcTable[ir.Name] = pr
	return pr
}

//  irProcedure makes a runtime procedure from a procedure info structure
func irProcedure(pr *pr_Info) *g.VProcedure {
	return g.GoProcedure(pr.name,
		func(args ...g.Value) (g.Value, *g.Closure) {
			//#%#% return interpMe(pr, args, ...) ??
			return nil, nil
		})
}

//  setupProc finishes procedure setup now that the GlobalDict is set
//	#%#% TODO: create LocalDict for finding variables
//	report undeclared identifiers
//	#%#% TODO: handle nested procedures
func setupProc(pr *pr_Info) {
	undeclared(pr)
	makedict(pr)
}

//  makedict creates the mapping from identifiers to variables within the proc
func makedict(pr *pr_Info) {
	pr.dict = make(map[string]interface{})
	// start with the globals; may overwrite some of these with locals
	//#%#% later: start with outer proc dict to grab its statics etc
	for name, value := range GlobalDict {
		pr.dict[name] = value
	}
	// add statics
	for _, name := range pr.ir.StaticList {
		v := g.Value(g.Trapped(g.NewStatic()))
		pr.dict[name] = g.Trapped(&v)
	}
	// add outer locals
	// add locals
	// add params
}

//  undeclared reports such identifiers
func undeclared(pr *pr_Info) {
	for _, chunk := range pr.ir.CodeList {
		for _, insn := range chunk.InsnList {
			if i, ok := insn.(ir_Var); ok {
				if !pr.lset[i.Name] && Undeclared[i.Name] {
					//%#% warn now, later fatal
					warning(fmt.Sprintf("%v %s undeclared",
						i.Coord, i.Name))
					// inhibit repeated warnings
					delete(Undeclared, pr.name)
				}
			}
		}
	}
}

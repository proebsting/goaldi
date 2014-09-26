//  proc.go -- things dealing with procedures in the interpreter

package main

import (
	"fmt"
	g "goaldi"
)

//  info about a procedure
type pr_Info struct {
	name    string                 // procedure name
	outer   *pr_Info               // immediate parent #%#% NOT YET SET
	ir      *ir_Function           // intermediate code structure
	nparams int                    // number of parameters
	nlocals int                    // number of locals
	lset    map[string]bool        // set of locally declared identifiers
	dict    map[string]interface{} // map of identifiers to variables
}

//  global index of procedure information
var ProcTable = make(map[string]*pr_Info)

//  a local variable
type pr_local int // value is index of this particular local

//  a parameter
type pr_param int // value is index of this particular parameter

//  declareProc initializes and returns a procedure info structure
func declareProc(ir *ir_Function) *pr_Info {
	pr := &pr_Info{}
	pr.name = ir.Name
	pr.ir = ir
	pr.nlocals = len(pr.ir.LocalList)
	pr.nparams = len(pr.ir.ParamList)
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
	return g.NewProcedure(pr.name,
		func(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
			//#%#% return interpMe(pr, args, ...) ??
			assert(false, "reached "+pr.name)
			return nil, nil
		})
}

//  setupProc finishes procedure setup now that the GlobalDict is set
//	#%#% TODO: handle nested procedures
//	report undeclared identifiers
//	create combined dictionary of global + local variables
func setupProc(pr *pr_Info) {
	undeclared(pr)
	pr.dict = makedict(pr)
	if opt_verbose {
		fmt.Printf("\n%s()  %d param  %d local  dict %d\n   ",
			pr.name, pr.nparams, pr.nlocals, len(pr.dict))
		for name := range sortedKeys(pr.dict) {
			v := pr.dict[name]
			switch x := v.(type) {
			case pr_param:
				fmt.Printf(" p%d:%s", int(x), name)
			case pr_local:
				fmt.Printf(" l%d:%s", int(x), name)
			default:
				fmt.Printf(" g:%s", name)
			}
		}
		fmt.Println()
	}
}

//  makedict creates the mapping from identifiers to variables within the proc
func makedict(pr *pr_Info) map[string]interface{} {

	dict := make(map[string]interface{})

	// start with the globals; may overwrite some of these with locals
	//#%#% later: start with outer proc dict to grab its statics etc
	for name, value := range GlobalDict {
		dict[name] = value
	}

	// add statics
	for _, name := range pr.ir.StaticList {
		v := g.Value(g.Trapped(g.NewStatic()))
		dict[name] = g.Trapped(&v)
	}
	// add outer locals
	//#%#% TBD

	// add locals
	for i, name := range pr.ir.LocalList {
		dict[name] = pr_local(i)
	}

	// add params
	for i, name := range pr.ir.ParamList {
		dict[name] = pr_param(i)
	}
	return dict
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

//  proc.go -- things dealing with procedures in the interpreter

package main

import (
	"fmt"
	g "goaldi"
)

//  info about a procedure
type pr_Info struct {
	name    string                   // procedure name
	outer   *pr_Info                 // immediate parent #%#% NOT YET SET
	ir      *ir_Function             // intermediate code structure
	accum   bool                     // true if last param is [] #%#% NYET IMPL
	nparams int                      // number of parameters
	nlocals int                      // number of locals
	lset    map[string]bool          // set of locally declared identifiers
	dict    map[string]interface{}   // map from identifiers to variables
	insns   map[string][]interface{} // map from labels to IR code chunks
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
	pr.accum = (ir.Accumulate != "")
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
	pr.dict = makeDict(pr)
	pr.insns = getInsns(pr)
	if opt_verbose {
		showProc(pr)
	}
}

//  makeDict creates the mapping from identifiers to variables within the proc
func makeDict(pr *pr_Info) map[string]interface{} {

	dict := make(map[string]interface{})

	// start with the globals; may overwrite some of these with locals
	//#%#% later: start with outer proc dict to grab its statics etc
	for name, value := range GlobalDict {
		dict[name] = value
	}

	// add statics
	for _, name := range pr.ir.StaticList {
		dict[name] = g.Trapped(g.NewVariable())
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
				if !pr.lset[i.Name] && Undeclared[i.Name] {
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

//  showProc prints information about the procedure in verbose mode
func showProc(pr *pr_Info) {
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

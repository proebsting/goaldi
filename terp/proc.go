//  proc.go -- things dealing with procedures in the interpreter

package main

import (
	"fmt"
	g "goaldi"
)

//  info about a procedure
type pr_Info struct {
	name  string          // procedure name
	outer *pr_Info        // immediate parent, if nested #%#% NOT YET SET
	ir    *ir_Function    // intermediate code structure
	lset  map[string]bool // set of locally declared identifiers
}

//  global index of procedure information
var ProcTable = make(map[string]*pr_Info)

//  declareProc initializes and returns a procedure info structure
func declareProc(ir *ir_Function) *pr_Info {
	info := &pr_Info{}
	info.name = ir.Name
	info.ir = ir
	info.lset = make(map[string]bool)
	for _, name := range ir.ParamList {
		info.lset[name] = true
	}
	for _, name := range ir.LocalList {
		info.lset[name] = true
	}
	for _, name := range ir.StaticList {
		info.lset[name] = true
	}
	ProcTable[ir.Name] = info
	return info
}

//  irProcedure makes a runtime procedure from a procedure info structure
func irProcedure(p *pr_Info) *g.VProcedure {
	return g.GoProcedure(p.name,
		func(args ...g.Value) (g.Value, *g.Closure) {
			//#%#% return interpMe(p, args, ...) ??
			return nil, nil
		})
}

//  setupProc finishes procedure setup now that the GlobalDict is set
//	#%#% TODO: create LocalDict for finding variables
//	report undeclared identifiers
//	#%#% TODO: handle nested procedures
func setupProc(p *pr_Info) {
	for _, chunk := range p.ir.CodeList {
		for _, insn := range chunk.InsnList {
			if i, ok := insn.(ir_Var); ok {
				if !p.lset[i.Name] && Undeclared[i.Name] {
					//%#% warn now, later fatal
					warning(fmt.Sprintf("%v %s undeclared",
						i.Coord, i.Name))
					// inhibit repeated warnings
					delete(Undeclared, p.name)
				}
			}
		}
	}
}

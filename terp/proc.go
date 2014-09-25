//  proc.go -- things dealing with procedures in the interpreter

package main

import (
	g "goaldi"
)

//  info about a procedure
type pr_Info struct {
	name  string          // procedure name
	ir    *ir_Function    // IR struct
	outer *pr_Info        // immediate parent, if nested
	lset  map[string]bool // set of locally declared identifiers
}

//  global index of procedure information
var ProcTable = make(map[string]*pr_Info)

//  pRegister enters a procedure into the info table
func pRegister(p *ir_Function) *pr_Info {
	info := &pr_Info{}
	info.name = p.Name
	info.ir = p
	info.lset = localSet(p)
	ProcTable[p.Name] = info
	return info
}

//  irProcedure makes a runtime procedure from an IR struct
//
//  #%#% this is just a skeleton -- details to be filled in later
func irProcedure(p *ir_Function) *g.VProcedure {
	//#%#% create static variables
	//#%#% set up data structures
	//#%#% handle nested procedures
	info := ProcTable[p.Name]
	assert(info != nil, "lost proc")
	return g.GoProcedure(p.Name,
		func(args ...g.Value) (g.Value, *g.Closure) {
			//#%#% return interpMe(info, args, ...) ??
			return nil, nil
		})
}

//  localSet(p) -- return set of locally declared ids
//  #%#%#% does not handle references to parent from nested procedure
func localSet(p *ir_Function) map[string]bool {
	lset := make(map[string]bool)
	for _, name := range p.ParamList {
		lset[name] = true
	}
	for _, name := range p.LocalList {
		lset[name] = true
	}
	for _, name := range p.StaticList {
		lset[name] = true
	}
	return lset
}

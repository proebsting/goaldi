//  proc.go -- things dealing with procedures in the interpreter

package main

import (
	g "goaldi"
)

//  irProcedure makes a runtime procedure from an IR struct
//
//  #%#% this is just a skeleton -- details to be filled in later
func irProcedure(p ir_Function) *g.VProcedure {
	//#%#% create static variables
	//#%#% set up data structures
	//#%#% handle nested procedures
	return g.GoProcedure(p.Name,
		func(args ...g.Value) (g.Value, *g.Closure) {
			//#%#% interpMe(...) ??
			return nil, nil
		})
}

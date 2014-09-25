//  link.go -- linking together loaded files

package main

import (
	"fmt"
	g "goaldi"
)

//  link combines IR files to make a complete program.
func link(parts [][]interface{}) UNKNOWN {
	babble("linking")

	walkTree(parts, gdecl1) // pass 1: declared globals
	walkTree(parts, gdecl2) // pass 2: declared procedures
	stdProcs()              // standard library
	walkTree(parts, undecl) // report undeclared identifiers
	return nil
}

//  walkTree calls a function for every top-level declaration in every file
func walkTree(parts [][]interface{}, f func(interface{})) {
	for _, file := range parts {
		for _, decl := range file {
			f(decl)
		}
	}
}

//  gdecl1 -- global dictionary, pass1
//	install declared global variables in global dictionary
//	install procedures in proc info table
//	note references to undeclared identifiers
func gdecl1(decl interface{}) {
	switch x := decl.(type) {
	case ir_Global:
		for _, name := range x.NameList {
			if GlobalDict[name] == nil {
				v := g.Value(g.NewNil())
				GlobalDict[name] = g.Trapped(&v)
			}
		}
	case ir_Function:
		pr := pRegister(&x)
		for _, chunk := range x.CodeList {
			for _, insn := range chunk.InsnList {
				if i, ok := insn.(ir_Var); ok {
					if !pr.lset[i.Name] {
						Undeclared[i.Name] = true
					}
				}
			}
		}
	case ir_Record:
		//#%#%#%# TBD
	case ir_Invocable, ir_Link:
		//#%#%#% nothing?
	default:
		panic(fmt.Sprintf("gdecl1: %#v", x))
	}
}

//  gdecl2 -- global dictionary, pass 2
//	satisfy undeclared IDs with declared procedures as constants
func gdecl2(decl interface{}) {
	switch x := decl.(type) {
	case ir_Global:
		// nothing to do
	case ir_Function:
		regProc(&x)
	case ir_Record:
		//#%#%#%# TBD
	case ir_Invocable, ir_Link:
		//#%#%#% nothing?
	default:
		panic(fmt.Sprintf("gdecl2: %#v", x))
	}
}

//  regProc(p) -- register procedure p in globals
func regProc(p *ir_Function) {
	gv := GlobalDict[p.Name]
	if gv == nil {
		// not declared as global, and not seen before:
		// create global with unmodifiable procedure value
		GlobalDict[p.Name] = irProcedure(p)
	} else if t, ok := gv.(*g.VTrapped); ok && t.Target == g.Value(g.NIL) {
		// uninitialized declared global:
		// initialize global trapped variable with procedure value
		*t.Target = irProcedure(p) //#%#% TEST THIS!
	} else {
		// duplicate global: fatal error
		fatal("duplicate global declaration: " + p.Name)
	}
	Undeclared[p.Name] = false
}

//  stdProcs() -- add referenced stdlib procedures to globals
func stdProcs() {
	for _, p := range g.StdLib {
		if Undeclared[p.Name] {
			if GlobalDict[p.Name] != nil {
				panic("undeclared but present: " + p.Name)
			}
			GlobalDict[p.Name] = p
			if Undeclared[p.Name] {
				Undeclared[p.Name] = false
			}
		}
	}
}

//  undecl -- report undeclared identifiers
func undecl(decl interface{}) {
	p, ok := decl.(ir_Function)
	if !ok { // if not a procedure declaration
		return
	}
	info := ProcTable[p.Name]
	for _, chunk := range p.CodeList {
		for _, insn := range chunk.InsnList {
			if i, ok := insn.(ir_Var); ok {
				if !info.lset[i.Name] && Undeclared[i.Name] {
					//%#% warn now, later fatal
					warning(fmt.Sprintf("%v %s undeclared",
						i.Coord, i.Name))
					// inhibit repeated warnings
					Undeclared[i.Name] = false
				}
			}
		}
	}
}

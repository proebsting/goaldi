//  link.go -- linking together loaded files

package main

import (
	"fmt"
	g "goaldi"
)

//  link combines IR files to make a complete program.
func link(parts [][]interface{}) UNKNOWN {
	babble("linking")

	walkTree(parts, gdecl1)
	walkTree(parts, gdecl2)
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
//	note references to undeclared identifiers
//	note declared procedures
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
		regUndecl(x)
	case ir_Record:
		//#%#%#%# TBD
	case ir_Invocable, ir_Link:
		//#%#%#% nothing?
	default:
		panic(fmt.Sprintf("gdecl1: %#v", x))
	}
}

//  regUndecl(p) -- register undeclared variables in procedure p
func regUndecl(p ir_Function) {
	localDict := make(map[string]bool)
	for _, name := range p.ParamList {
		localDict[name] = true
	}
	for _, name := range p.LocalList {
		localDict[name] = true
	}
	for _, name := range p.StaticList {
		localDict[name] = true
	}
	for _, chunk := range p.CodeList {
		for _, insn := range chunk.InsnList {
			if i, ok := insn.(ir_Var); ok {
				if !localDict[i.Name] {
					Undeclared[i.Name] = true
				}
			}
		}
	}
}

//  gdecl2 -- global dictionary, pass 2
//	satisfy undeclared IDs with declared and stdlib procedures as constants
func gdecl2(decl interface{}) {
	switch x := decl.(type) {
	case ir_Global:
		// nothing to do
	case ir_Function:
		regProc(x)
	case ir_Record:
		//#%#%#%# TBD
	case ir_Invocable, ir_Link:
		//#%#%#% nothing?
	default:
		panic(fmt.Sprintf("gdecl2: %#v", x))
	}
}

//  regProc(p) -- register procedure p in globals
func regProc(p ir_Function) {
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

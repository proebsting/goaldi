//  link.go -- linking together loaded files

package main

import (
	"fmt"
	g "goaldi"
)

//  link combines IR files to make a complete program.
func link(parts [][]interface{}) UNKNOWN {

	babble("linking")

	//  process individual declarations (proc, global, etc) from IR
	for _, file := range parts {
		for _, decl := range file {
			irDecl(decl)
		}
	}

	//  register procedures in global namespace
	for _, info := range ProcTable {
		registerProc(info)
	}

	//  add standard library procedures for names not yet found
	stdProcs()

	// set up procedures and report undeclared identifiers
	for _, info := range ProcTable {
		setupProc(info)
	}

	return nil
}

//  irDecl -- process IR file declaration
//	install declared global variables in global dictionary
//	install procedures in proc info table
//	note references to undeclared identifiers
func irDecl(decl interface{}) {
	switch x := decl.(type) {
	case ir_Global:
		for _, name := range x.NameList {
			if GlobalDict[name] == nil {
				v := g.Value(g.NewNil())
				GlobalDict[name] = g.Trapped(&v)
			}
		}
	case ir_Function:
		pr := declareProc(&x)
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

//  registerProc(p) -- register procedure p in globals
func registerProc(p *pr_Info) {
	gv := GlobalDict[p.name]
	if gv == nil {
		// not declared as global, and not seen before:
		// create global with unmodifiable procedure value
		GlobalDict[p.name] = irProcedure(p)
	} else if t, ok := gv.(*g.VTrapped); ok && t.Target == g.Value(g.NIL) {
		// uninitialized declared global:
		// initialize global trapped variable with procedure value
		*t.Target = irProcedure(p) //#%#% TEST THIS!
	} else {
		// duplicate global: fatal error
		fatal("duplicate global declaration: " + p.name)
	}
	delete(Undeclared, p.name)
}

//  stdProcs() -- add referenced stdlib procedures to globals
func stdProcs() {
	for _, p := range g.StdLib {
		if Undeclared[p.Name] {
			if GlobalDict[p.Name] != nil {
				panic("undeclared but present: " + p.Name)
			}
			GlobalDict[p.Name] = p
			delete(Undeclared, p.Name)
		}
	}
}

//  link.go -- linking together loaded files

package main

import (
	"fmt"
	g "goaldi"
)

//  list of struct (record) declarations seen
var StructList = make([]*ir_Record, 0)

//  link combines IR files to make a complete program.
func link(parts [][]interface{}) {

	babble("linking")

	//  process individual declarations (proc, global, etc) from IR
	for _, file := range parts {
		for _, decl := range file {
			irDecl(decl)
		}
	}

	//  register procedures in global namespace
	for _, pr := range ProcTable {
		registerProc(pr)
	}

	//  register struct constructors
	for _, sc := range StructList {
		registerStruct(sc)
	}

	//  add standard library procedures for names not yet found
	stdProcs()

	//  remove globals from "Undeclared" list
	for name := range GlobalDict {
		delete(Undeclared, name)
	}

	// set up procedures and report undeclared identifiers
	for _, pr := range ProcTable {
		setupProc(pr)
	}
}

//  irDecl -- process IR file declaration
//	install declared global variables in global dictionary
//		#%#% as trapped variables; is that really the best way?
//	install procedures in proc info table
//	note references to undeclared identifiers
func irDecl(decl interface{}) {
	switch x := decl.(type) {
	case ir_Global:
		for _, name := range x.NameList {
			if GlobalDict[name] == nil {
				GlobalDict[name] = g.Trapped(g.NewVariable(g.NilValue))
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
		StructList = append(StructList, &x)
	case ir_Invocable, ir_Link:
		//#%#%#% nothing?
	default:
		panic(fmt.Sprintf("gdecl1: %#v", x))
	}
}

//  registerProc(pr) -- register procedure pr in globals
func registerProc(pr *pr_Info) {
	gv := GlobalDict[pr.name]
	if gv == nil {
		// not declared as global, and not seen before:
		// create global with unmodifiable procedure value
		GlobalDict[pr.name] = irProcedure(pr)
	} else if t, ok := gv.(*g.VTrapped); ok && t.Target == &g.NilValue {
		// uninitialized declared global:
		// initialize global trapped variable with procedure value
		*t.Target = irProcedure(pr) //#%#% TEST THIS!
	} else {
		// duplicate global: fatal error
		fatal("duplicate global declaration: " + pr.name)
	}
	delete(Undeclared, pr.name)
}

//  registerStruct(sc) -- register struct constructor in globals
func registerStruct(sc *ir_Record) {
	gv := GlobalDict[sc.Name]
	if gv == nil {
		// not declared as global, and not seen before:
		// create global with unmodifiable procedure value
		GlobalDict[sc.Name] = g.NewDefn(sc.Name, sc.FieldList)
	} else if t, ok := gv.(*g.VTrapped); ok && t.Target == &g.NilValue {
		// uninitialized declared global:
		// initialize global trapped variable with procedure value
		*t.Target = g.NewDefn(sc.Name, sc.FieldList)
	} else {
		// duplicate global: fatal error
		fatal("duplicate global declaration: " + sc.Name)
	}
	delete(Undeclared, sc.Name)
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

//  link.go -- linking together loaded files

package main

import (
	"fmt"
	g "goaldi"
	"strings"
)

//  list of record declarations seen
var RecordList = make([]*ir_Record, 0)

//  link combines IR files to make a complete program.
func link(parts [][]interface{}) {

	babble("linking")

	//  process individual declarations (proc, global, etc) from IR
	for _, file := range parts {
		for _, decl := range file {
			irDecl(decl)
		}
	}

	//  register the record constructors
	for _, sc := range RecordList {
		registerRecord(sc)
	}

	//  register methods in constructors and procedures in global namespace
	for _, pr := range ProcTable {
		a := strings.Split(pr.name, ".") // look for xxx.yyy form
		if len(a) == 1 {                 // if simple procedure name
			registerProc(pr)
		} else { // no, this is typename.methodname
			registerMethod(pr, a[0], a[1])
		}
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
//
//	Install declared global variables as trapped refs in global dictionary.
//	Install procedures in proc info table.
//  Register initial procedures and global initialization procedures.
func irDecl(decl interface{}) {
	switch x := decl.(type) {
	case ir_Global:
		for _, name := range x.NameList {
			gv := GlobalDict[name]
			if gv == nil {
				GlobalDict[name] = g.NewVariable(g.NilValue)
			} else if t, ok := gv.(*g.VTrapped); ok && *t.Target == g.NilValue {
				// okay, previously declared global, no problem
			} else {
				fatal("duplicate global declaration: global " + name)
			}
		}
		if x.Fn != "" {
			GlobInit = append(GlobInit, &x)
		}
	case ir_Initial:
		InitList = append(InitList, &x)
	case ir_Function:
		declareProc(&x)
		for _, id := range x.UnboundList {
			Undeclared[id] = true
		}
	case ir_Record:
		RecordList = append(RecordList, &x)
	default: // including ir_Invocable, ir_Link
		panic(g.Malfunction(fmt.Sprintf("unrecognized: %#v", x)))
	}
}

//  registerMethod(pr, recname, methname) -- register method in record ctor
func registerMethod(pr *pr_Info, recname string, methname string) {
	gv := GlobalDict[recname]
	if gv != nil {
		gv = g.Deref(gv)
	}
	if d, ok := gv.(*g.VCtor); ok && d != nil {
		if !d.AddMethod(methname, irProcedure(pr, nil)) {
			fatal(fmt.Sprintf("Method %s.%s() duplicates field name %s",
				recname, methname, methname))
		}
	} else {
		fatal(fmt.Sprintf("No type %s found for method %s.%s()",
			recname, recname, methname))
	}
}

//  registerProc(pr) -- register procedure pr in globals
func registerProc(pr *pr_Info) {
	gv := GlobalDict[pr.name]
	if gv == nil {
		// create global with unmodifiable procedure value
		GlobalDict[pr.name] = irProcedure(pr, nil)
	} else {
		// duplicate global: fatal error
		fatal("duplicate global declaration: procedure " + pr.name)
	}
	delete(Undeclared, pr.name)
}

//  registerRecord(sc) -- register a record constructor in the globals
func registerRecord(sc *ir_Record) {
	gv := GlobalDict[sc.Name]
	if gv == nil {
		// create global with unmodifiable procedure value
		GlobalDict[sc.Name] = g.NewCtor(sc.Name, sc.FieldList)
	} else {
		// duplicate global: fatal error
		fatal("duplicate global declaration: record " + sc.Name)
	}
	delete(Undeclared, sc.Name)
}

//  stdProcs() -- add referenced stdlib procedures to globals
func stdProcs() {
	for name, p := range g.StdLib {
		if Undeclared[name] {
			if GlobalDict[name] != nil {
				panic(g.Malfunction("undeclared but present: " + name))
			}
			GlobalDict[name] = p
			delete(Undeclared, name)
		}
	}
}

//  link.go -- linking together loaded files

package main

import (
	"fmt"
	"goaldi/ir"
	g "goaldi/runtime"
	"strings"
)

//  A RecordEntry adds info to an Ir_Record
type RecordEntry struct {
	ir.Ir_Record          // ir struct
	ctor         *g.VCtor // constructor
}

//  RecordTable registers all the record declarations that have been seen
var RecordTable = make(map[string]*RecordEntry, 0)

//  link combines IR files to make a complete program.
func link(parts [][]interface{}) {

	//  process individual declarations (proc, global, etc) from IR
	for _, file := range parts {
		for _, decl := range file {
			irDecl(decl)
		}
	}

	//  register the record constructors
	for _, re := range RecordTable {
		registerRecord(re)
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
	for name := range PubSpace.All() {
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
	case ir.Ir_Global:
		name := x.Name
		ns := g.GetSpace(x.Namespace)
		gv := ns.Get(name)
		if gv == nil {
			ns.Declare(name, g.NewVariable(g.NilValue))
		} else if t, ok := gv.(*g.VTrapped); ok && *t.Target == g.NilValue {
			// okay, previously declared global, no problem
		} else {
			fatal("duplicate global declaration: global " + name)
		}
		if x.Fn != "" {
			GlobInit = append(GlobInit, &x)
		}
	case ir.Ir_Initial:
		InitList = append(InitList, &x)
	case ir.Ir_Function:
		declareProc(&x)
		for _, id := range x.UnboundList {
			if !strings.Contains(id, "::") { // if no explicit namespace
				Undeclared[id] = true
			}
		}
	case ir.Ir_Record:
		ns := g.GetSpace(x.Namespace)
		qname := ns.GetQual() + x.Name
		if RecordTable[qname] == nil {
			RecordTable[qname] = &RecordEntry{x, nil}
		} else {
			fatal("duplicate record declaration: record " + qname)
		}
	default:
		panic(g.Malfunction(fmt.Sprintf("unrecognized: %#v", x)))
	}
}

//  registerMethod(pr, recname, methname) -- register method in record ctor
func registerMethod(pr *pr_Info, recname string, methname string) {
	gv := pr.space.Get(recname)
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
	pr.vproc = irProcedure(pr, nil)
	gv := pr.space.Get(pr.name)
	if gv == nil {
		// create global with unmodifiable procedure value
		pr.space.Declare(pr.name, pr.vproc)
	} else {
		// duplicate global: fatal error
		fatal("duplicate global declaration: procedure " + pr.name)
	}
	delete(Undeclared, pr.name)
}

//  registerRecord(re) -- register a record constructor in the globals
func registerRecord(re *RecordEntry) {
	defer func() { // catch "duplicate field name" exception
		if e := recover(); e != nil {
			x := e.(*g.Exception)
			fatal(fmt.Sprintf("in record %s: %s: %v",
				re.Name, x.Msg, x.Offv[0]))
		}
	}()
	if re.ctor == nil { // if not already processed
		re.ctor = regMark // prevent infinite recursion on error
		ns := g.GetSpace(re.Namespace)
		gv := ns.Get(re.Name)
		if gv == nil {
			// this is a new definition
			var ext *g.VCtor
			if re.ExtendsRec != "" {
				pt := RecordTable[re.ExtendsRec]
				if pt == nil {
					fatal("parent type not found: record " +
						re.Name + " extends " + re.ExtendsRec)
				} else if pt.ctor == regMark {
					fatal("recursive definition: record " +
						re.Name + " extends " + re.ExtendsRec + " extends...")
				} else {
					registerRecord(pt) // ensure parent is done first
					ext = pt.ctor
				}
			}
			// create global with unmodifiable procedure value
			re.ctor = g.NewCtor(re.Name, ext, re.FieldList)
			ns.Declare(re.Name, re.ctor)
		} else {
			// duplicate global: fatal error
			fatal("duplicate global declaration: record " + re.Name)
		}
		delete(Undeclared, re.Name)
	}
}

var regMark = &g.VCtor{} // marker for catching recursive definitions

//  stdProcs() -- add referenced stdlib procedures to globals
func stdProcs() {
	for name, p := range g.StdLib {
		if Undeclared[name] {
			if PubSpace.Get(name) != nil {
				panic(g.Malfunction("undeclared but present: " + name))
			}
			PubSpace.Declare(name, p)
			delete(Undeclared, name)
		}
	}
}

//  proc.go -- things dealing with procedures at link time

package main

import (
	"github.com/proebsting/goaldi/ir"
	g "github.com/proebsting/goaldi/runtime"
	"strings"
	"unicode"
)

// information about a procedure that is shared by all invocations
type pr_Info struct {
	space    *g.Namespace             // procedure namespace
	name     string                   // procedure name
	qname    string                   // qualified name (namespace::name)
	ir       *ir.Ir_Function          // intermediate code structure
	insns    map[string][]interface{} // map from labels to IR code chunks
	statics  map[string]interface{}   // table of statics (including globals)
	locals   []string                 // list of local names
	params   []string                 // list of parameter names
	variadic bool                     // true if last param is []
	ntemps   int                      // number of temporaries
	vproc    *g.VProcedure            // execution-time procedure struct
}

// global index of procedure information (indexed by qualified name)
var ProcTable = make(map[string]*pr_Info)

// declareProc initializes and returns a procedure info structure
func declareProc(irf *ir.Ir_Function) *pr_Info {
	pr := &pr_Info{}
	pr.name = irf.Name
	pr.space = g.GetSpace(irf.Namespace)
	if unicode.IsDigit(rune(pr.name[0])) { // if generated procedure
		pr.qname = pr.name // leave the name alone
	} else { // if explicit user procedure
		pr.qname = pr.space.GetQual() + pr.name // add namespace qualifier
	}
	if ProcTable[pr.qname] != nil {
		fatal("Duplicate procedure definition: " + irf.Name)
	}
	pr.ir = irf
	pr.variadic = (irf.Accumulate != "")
	pr.params = pr.ir.ParamList
	pr.locals = pr.ir.LocalList
	pr.ntemps = pr.ir.TempCount
	ProcTable[pr.qname] = pr
	return pr
}

// setupProc finishes procedure setup now that all globals are known
func setupProc(pr *pr_Info) {

	// add qualifiers to unbound identifiers for dependency processing
	// report identifiers not declared anywhere
	for i, id := range pr.ir.UnboundList {
		nsid := strings.Split(id, "::")
		if len(nsid) == 1 { // if unqualified name
			if pr.space.Get(id) != nil {
				// found in current space; make this explicit
				pr.ir.UnboundList[i] = pr.space.GetQual() + id
			} else if PubSpace.Get(id) == nil {
				fatal("In " + pr.qname + "(): Undeclared identifier: " + id)
			}
		} else { // explicitly qualified by namespace
			if g.GetSpace(nsid[0]).Get(nsid[1]) == nil {
				fatal("In " + pr.qname + "(): Undeclared identifier: " + id)
			}
		}
	}

	// add this proc to outer procedure's dependency list
	if pr.ir.Parent != "" {
		pt := ProcTable[pr.space.GetQual()+pr.ir.Parent]
		pt.ir.UnboundList = append(pt.ir.UnboundList, pr.qname)
	}

	// make a trapped variable for every static
	pr.statics = make(map[string]interface{})
	for _, name := range pr.ir.StaticList {
		pr.statics[name] = g.NewVariable(g.NilValue)
	}

	// create an index of IR code chunks
	pr.insns = make(map[string][]interface{})
	for _, ch := range pr.ir.CodeList {
		if pr.insns[ch.Label] != nil {
			panic(g.Malfunction("Duplicate IR label: " + ch.Label))
		}
		pr.insns[ch.Label] = ch.InsnList
	}
}

// irProcedure makes a runtime procedure from static info and inherited vars
func irProcedure(pr *pr_Info, outer map[string]interface{}) *g.VProcedure {

	// make a list of unadorned parameter names
	pnames := make([]string, len(pr.params))
	for i, s := range pr.params {
		pnames[i] = s[:strings.Index(s, ":")]
	}

	// copy (references to) any inherited variables
	vars := make(map[string]interface{})
	if outer != nil {
		for k, v := range outer {
			vars[k] = v
		}
	}

	return g.NewProcedure(pr.qname, &pnames, pr.variadic,
		func(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
			return interp(env, pr, vars, args...)
		}, nil, "")
}

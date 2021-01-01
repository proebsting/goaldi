//  execute.go -- the interpreter main loop

package main

import (
	"fmt"
	"github.com/proebsting/goaldi/ir"
	g "github.com/proebsting/goaldi/runtime"
)

//  iLiteral replaces Ir_NilLit, Ir_IntLit, Ir_RealLit, Ir_StrLit
type iLiteral struct {
	Lhs   int
	Value g.Value
}

//  coexecute wraps an execute call to catch a panic in a co-expression
func coexecute(f *pr_frame, label string) (g.Value, *g.Closure) {
	defer g.Catcher(f.env)
	return execute(f, label)
}

//  execute dispatches and interprets IR instructions for a procedure or coexpr
func execute(f *pr_frame, label string) (rv g.Value, rc *g.Closure) {

	// set up error catcher to call user recovery procedure
	defer func() {
		if p := recover(); p != nil {
			if f.onerr != nil { // if user called recover()
				// find true panic value hiding under traceback info
				arglist := []g.Value{g.Cause(p)}
				if opt_trace {
					fmt.Printf("[%d] panic: %v\n", f.env.ThreadID, arglist[0])
					fmt.Printf("[%d] calling %v\n", f.env.ThreadID, f.onerr)
				}
				// call recovery procedure and return its result
				rv, _ = f.onerr.Call(f.env, arglist, []string{})
				rc = nil
			} else {
				// re-throw the exception
				panic(p)
			}
		}
	}()

	// create re-entrant interpreter
	var self *g.Closure
	self = &g.Closure{func() (g.Value, *g.Closure) {

		// set up traceback recovery
		// (must do that here to include resumed procedures in traceback)
		defer func() {
			if p := recover(); p != nil {
				// add traceback information and re-throw exception
				panic(g.Catch(p,
					[]g.Value{f.offv}, f.coord, f.info.name, f.args))
			}
		}()

		// interpret the IR code (main loop)
		// a "goto" in the IR code sets "label" and jumps to NextChunk
	NextChunk:
		for {
			if opt_trace {
				fmt.Printf("[%d] %s:\n", f.env.ThreadID, label)
			}
			ilist := f.info.insns[label] // look up label
			if len(ilist) == 0 {
				panic(g.Malfunction("No instructions for IR label: " + label))
			}
			inchunk := label
			label = "UNSET"        // should never see this
			for j := range ilist { // execute insns in chunk
			Dispatch: // jump target for reissuing JIT-modified instruction
				insn := ilist[j]
				if opt_trace {
					t := fmt.Sprintf("%T", insn)[6:]
					fmt.Printf("[%d]    %s %v\n", f.env.ThreadID, t, insn)
				}
				f.coord = "" // unnecessary but prudent
				f.offv = nil // unnecessary but prudent
				switch i := insn.(type) {
				default: // incl ScanSwap, Assign, Deref, Unreachable
					panic(g.Malfunction(fmt.Sprintf(
						"Unrecognized interpreter instruction: %#v", i)))
				case ir.Ir_NoOp:
					// nothing to do
				case ir.Ir_Fail:
					return nil, nil
				case ir.Ir_Succeed:
					v := g.Deref(f.temps[i.Expr].(g.Value))
					if i.ResumeLabel == "" {
						return v, nil
					} else {
						label = i.ResumeLabel
						return v, self
					}
				case ir.Ir_Catch:
					f.offv = g.Deref(f.temps[i.Fn])
					if f.offv == g.NilValue {
						f.onerr = nil // clear if nil
					} else {
						f.onerr = f.offv.(*g.VProcedure) // else must be proc
					}
					if i.Lhs != 0 {
						f.temps[i.Lhs] = f.onerr
					}
				case ir.Ir_Create:
					fnew := newframe(f)
					fnew.cxout = g.NewChannel(0)
					e := g.NewEnv(f.env)
					e.ThreadID = <-g.TID
					e.VarMap["current"] = fnew.cxout // set %current
					fnew.env = e
					fnew.vars[i.Scope] = e
					fnew.coord = i.Coord
					if i.Lhs != 0 {
						f.temps[i.Lhs] = fnew.cxout
					}
					go coexecute(fnew, i.CoexpLabel)
				case ir.Ir_Select:
					label = irSelect(f, i)
					continue NextChunk
				case ir.Ir_CoRet:
					f.coord = i.Coord
					if g.CoSend(f.cxout, f.temps[i.Value]) == nil {
						return nil, nil // kill self: channel was closed
					}
					label = i.ResumeLabel
					continue NextChunk
				case ir.Ir_CoFail:
					close(f.cxout)
					return nil, nil // i.e. die
				case ir.Ir_Key: // dynamic variable reference
					f.coord = i.Coord
					e := f.vars[i.Scope].(*g.Env) // get correct environment
					v := e.Lookup(i.Name, i.Rval != "")
					if i.Lhs != 0 {
						f.temps[i.Lhs] = v
					}
				case iLiteral: // replaces ir_{Nil,Int,Real,Str}Lit
					f.temps[i.Lhs] = i.Value
				case ir.Ir_NilLit:
					ilist[j] = iLiteral{i.Lhs, g.NilValue}
					goto Dispatch
				case ir.Ir_IntLit:
					n, _ := g.ParseNumber(i.Val)
					ilist[j] = iLiteral{i.Lhs, g.NewNumber(n)}
					goto Dispatch
				case ir.Ir_RealLit:
					n, _ := g.ParseNumber(i.Val)
					ilist[j] = iLiteral{i.Lhs, g.NewNumber(n)}
					goto Dispatch
				case ir.Ir_StrLit:
					ilist[j] = iLiteral{i.Lhs, g.NewString(i.Val)}
					goto Dispatch
				case ir.Ir_MakeList:
					n := len(i.ValueList)
					a := make([]g.Value, n)
					for j, t := range i.ValueList {
						a[j] = g.Deref(f.temps[t])
					}
					f.temps[i.Lhs] = g.InitList(a)
				case ir.Ir_Var:
					var v g.Value
					if i.Namespace != "" {
						v = g.GetSpace(i.Namespace).Get(i.Name)
					} else {
						v = f.vars[i.Name]
						if v == nil {
							v = f.info.space.Get(i.Name)
							if v == nil {
								v = PubSpace.Get(i.Name)
							}
						}
					}
					if v == nil {
						panic(g.Malfunction("Unbound identifier: " +
							i.Namespace + "::" + i.Name))
					}
					if i.Rval != "" {
						v = g.Deref(v)
					}
					f.temps[i.Lhs] = v
				case ir.Ir_EnterScope:
					e := f.env                 // environment at procedure entry
					p := f.vars[i.ParentScope] // look it up
					if p != nil {              // if known
						e = p.(*g.Env) // now e has our current env
					}
					if len(i.DynamicList) > 0 { // if any dynamic vars declared
						e = g.NewEnv(e)                      // make new env
						for _, name := range i.DynamicList { // install dynamics
							e.VarMap[name] = g.NewVariable(nil) // uninitialized
						}
					}
					f.vars[i.Scope] = e               // save envmt of scope
					for _, name := range i.NameList { // init locals
						f.vars[name] = g.NewVariable(g.NilValue)
					}
				case ir.Ir_ExitScope:
					for _, name := range i.NameList {
						f.vars[name] = nil // allow garbage collection
					}
					for _, name := range i.DynamicList {
						f.env.VarMap[name] = nil
					}
				case ir.Ir_Move:
					f.temps[i.Lhs] = f.temps[i.Rhs]
				case ir.Ir_MoveLabel:
					f.temps[i.Lhs] = i.Label
				case ir.Ir_Goto:
					label = i.TargetLabel
					continue NextChunk
				case ir.Ir_IndirectGoto:
					label = f.temps[i.TargetTmpLabel].(string)
					for _, s := range i.LabelList {
						if s == label {
							continue NextChunk
						}
					}
					panic(g.Malfunction(
						"IndirectGoto: unlisted label: " + label))
				case ir.Ir_MakeClosure:
					// potential future optimization:
					// only pass in *referenced* variables
					// so that the remainder can get garbage collected
					f.temps[i.Lhs] = irProcedure(ProcTable[i.Name], f.vars)
				case ir.Ir_OpFunction:
					// replace by equivalent iOperator and re-dispatch
					e := getOperator(&i)
					ilist[j] = *e
					goto Dispatch
				case iOperator:
					v, c := operate(f.env, f, &i)         // execute operation
					if (i.Flags&rflag) != 0 && v != nil { // if rval needed
						v = g.Deref(v) // then make sure we have one
						// note v can be set nil by failing Deref
					}
					if v != nil {
						if i.Lhs != 0 {
							f.temps[i.Lhs] = v
						}
						if i.Lhsclosure != 0 {
							f.temps[i.Lhsclosure] = c
						}
					} else if i.FailLabel != "" {
						label = i.FailLabel
						continue NextChunk
					}
				case ir.Ir_Field:
					f.coord = i.Coord
					x := g.Deref(f.temps[i.Expr].(g.Value))
					v := g.Field(x, i.Field)
					if v != nil {
						if i.Rval != "" { // if an rval is required
							v = g.Deref(v) // then make sure we have one
						}
						if i.Lhs != 0 {
							f.temps[i.Lhs] = v
						}
					}
				case ir.Ir_Call:
					f.coord = i.Coord
					proc := g.Deref(f.temps[i.Fn].(g.Value))
					n := len(i.ArgList)
					arglist := make([]g.Value, n)
					for j, a := range i.ArgList {
						v := f.temps[a]
						arglist[j] = g.Deref(v.(g.Value))
					}
					f.offv = proc
					e := f.vars[i.Scope].(*g.Env) // get correct environment
					v, c := proc.(g.ICall).Call(e, arglist, i.NameList)
					if v != nil {
						if i.Lhs != 0 {
							f.temps[i.Lhs] = v
						}
						if i.Lhsclosure != 0 {
							f.temps[i.Lhsclosure] = c
						}
					} else if i.FailLabel != "" {
						label = i.FailLabel
						continue NextChunk
					}
				case ir.Ir_ResumeValue:
					f.coord = i.Coord
					var v g.Value
					c := f.temps[i.Closure].(*g.Closure)
					if c != nil {
						v, c = c.Go()
					}
					if v != nil {
						if i.Lhs != 0 {
							f.temps[i.Lhs] = v
						}
						if i.Lhsclosure != 0 {
							f.temps[i.Lhsclosure] = c
						}
					} else if i.FailLabel != "" {
						label = i.FailLabel
						continue NextChunk
					}
				}
			}
			panic(g.Malfunction("Ir_Chunk exhausted: " + inchunk))
		}
	}}

	// start up the interpreter
	return self.Resume()
}

//  irSelect -- execute select statement, returning label of chosen case body
func irSelect(f *pr_frame, irs ir.Ir_Select) string {

	// set up data structures for selection
	s := g.NewSelector(len(irs.CaseList))
	for _, sc := range irs.CaseList {
		f.coord = sc.Coord
		switch sc.Kind {
		case "send":
			s.SendCase(g.Deref(f.temps[sc.Lhs]), g.Deref(f.temps[sc.Rhs]))
		case "receive":
			s.RecvCase(g.Deref(f.temps[sc.Rhs]))
		case "default":
			s.DfltCase()
		default:
			panic(g.Malfunction("Bad SelectCase kind: " + sc.Kind))
		}
	}

	// do the selection
	f.coord = irs.Coord
	i, v := s.Execute()

	if i < 0 {
		return irs.FailLabel // select failed, no default case supplied
	}
	sc := irs.CaseList[i]
	f.coord = sc.Coord
	if sc.Kind == "receive" {
		// assign received value before executing body
		f.temps[sc.Lhs].(g.IVariable).Assign(v)
	}
	return sc.BodyLabel
}

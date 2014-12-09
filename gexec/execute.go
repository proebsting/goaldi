//  execute.go -- the interpreter main loop

package main

import (
	"fmt"
	g "goaldi"
)

//  execute IR instructions for procedure or co-expression
func execute(f *pr_frame, label string) (g.Value, *g.Closure) {

	// set up traceback recovery
	defer func() {
		if p := recover(); p != nil {
			panic(g.Catch(p, f.offv, f.coord, f.info.name, f.args))
		}
	}()

	// create re-entrant interpreter
	f.temps = make(map[string]interface{}) // each cx needs own copy
	var self *g.Closure
	self = &g.Closure{func() (g.Value, *g.Closure) {

		// interpret the IR code
		for {
			if opt_trace {
				fmt.Printf("[%d] L: %s\n", f.env.ThreadID, label)
			}
			ilist := f.info.insns[label] // look up label
		Chunk:
			for _, insn := range ilist { // execute insns in chunk
				if opt_trace {
					fmt.Printf("[%d] I: %T %v\n", f.env.ThreadID, insn, insn)
				}
				f.coord = "" //#%#% prudent, but s/n/b needed
				f.offv = nil //#%#% prudent, but s/n/b needed
				switch i := insn.(type) {
				default: // incl ScanSwap, Assign, Deref, Unreachable
					panic(&g.RunErr{
						"Unrecognized interpreter instruction",
						fmt.Sprintf("%#v", i)})
				case ir_Fail:
					return nil, nil
				case ir_Succeed:
					v := g.Deref(f.temps[i.Expr].(g.Value))
					if i.ResumeLabel == "" {
						return v, nil
					} else {
						label = i.ResumeLabel
						return v, self
					}
				case ir_Create:
					fnew := newframe(f)
					fnew.cxout = g.NewChannel(0)
					fnew.env = g.NewEnv(f.env)
					fnew.coord = i.Coord
					if i.Lhs != "" {
						f.temps[i.Lhs] = fnew.cxout
					}
					go coexecute(fnew, i.CoexpLabel)
				case ir_CoRet:
					f.coord = i.Coord
					if g.CoSend(f.cxout, f.temps[i.Value]) == nil {
						return nil, nil // kill self: channel was closed
					}
					label = i.ResumeLabel
				case ir_CoFail:
					close(f.cxout)
					return nil, nil // i.e. die
				case ir_Key:
					//#%#% keywords are dynamic vars fetched from env
					f.coord = i.Coord
					v := f.env.VarMap[i.Name]
					if v == nil {
						panic(&g.RunErr{"Unrecognized dynamic variable",
							"%" + i.Name})
					}
					//#%#% ignoring failure and FailLabel
					if i.Lhs != "" {
						f.temps[i.Lhs] = v
					}
				case ir_IntLit:
					f.temps[i.Lhs] = g.NewString(i.Val).ToNumber()
				case ir_RealLit:
					f.temps[i.Lhs] = g.NewString(i.Val).ToNumber()
				case ir_StrLit:
					f.temps[i.Lhs] = g.NewString(i.Val)
				case ir_MakeList:
					n := len(i.ValueList)
					a := make([]g.Value, n)
					for j, v := range i.ValueList {
						a[j] = g.Deref(f.temps[v.(string)])
					}
					f.temps[i.Lhs] = g.InitList(a)
				case ir_Var:
					frame := f
					v := frame.vars[i.Name]
					if v == nil {
						v = GlobalDict[i.Name]
					}
					if v == nil {
						//#%#% eventually make a link-time error
						panic(&g.RunErr{"Undeclared identifier", i.Name})
					}
					f.temps[i.Lhs] = v
				case ir_Move:
					f.temps[i.Lhs] = f.temps[i.Rhs]
				case ir_MoveLabel:
					f.temps[i.Lhs] = i.Label
				case ir_Goto:
					label = i.TargetLabel
					break Chunk
				case ir_IndirectGoto:
					label = i.TargetTmpLabel
					label = f.temps[label].(string)
					break Chunk
				case ir_MakeClosure:
					//#%#% potential later optimization:
					//#%#% only pass in *referenced* variables
					//#%#% so that the remainder can get garbage collected
					f.temps[i.Lhs] = irProcedure(ProcTable[i.Name], f.vars)
				case ir_OpFunction:
					f.coord = i.Coord
					v, c := operator(f.env, f, &i)
					if v != nil {
						if i.Lhs != "" {
							f.temps[i.Lhs] = v
						}
						if i.Lhsclosure != "" {
							f.temps[i.Lhsclosure] = c
						}
					} else if i.FailLabel != "" {
						label = i.FailLabel
						break Chunk
					}
				case ir_Field:
					f.coord = i.Coord
					x := g.Deref(f.temps[i.Expr].(g.Value))
					v := g.Field(x, i.Field)
					if v != nil {
						if i.Lhs != "" {
							f.temps[i.Lhs] = v
						}
					} else if i.FailLabel != "" {
						label = i.FailLabel
						break Chunk
					}
				case ir_Call:
					f.coord = i.Coord
					proc := g.Deref(f.temps[i.Fn].(g.Value))
					argl := getArgs(f, 0, i.ArgList)
					f.offv = proc
					v, c := proc.(g.ICall).Call(f.env, argl...)
					if v != nil {
						if i.Lhs != "" {
							f.temps[i.Lhs] = v
						}
						if i.Lhsclosure != "" {
							f.temps[i.Lhsclosure] = c
						}
					} else if i.FailLabel != "" {
						label = i.FailLabel
						break Chunk
					}
				case ir_ResumeValue:
					f.coord = i.Coord
					var v g.Value
					c := f.temps[i.Closure].(*g.Closure)
					if c != nil {
						v, c = c.Go()
					}
					if v != nil {
						if i.Lhs != "" {
							f.temps[i.Lhs] = v
						}
						if i.Lhsclosure != "" {
							f.temps[i.Lhsclosure] = c
						}
					} else if i.FailLabel != "" {
						label = i.FailLabel
						break Chunk
					}
				}
			}
		}
		return nil, nil
	}}

	// start up the interpreter
	return self.Resume()
}

//  getArgs -- load values from heterogeneous ArgList slice field
//  nd is the number of leading arguments that should *not* be dereferenced
func getArgs(f *pr_frame, nd int, arglist []interface{}) []g.Value {
	n := len(arglist)
	argl := make([]g.Value, n)
	for i, a := range arglist {
		switch t := a.(type) {
		case string:
			a = f.temps[t]
		default:
			// nothing to do: use entry as is
		}
		if i < nd {
			argl[i] = a.(g.Value)
		} else {
			argl[i] = g.Deref(a.(g.Value))
		}
	}
	return argl
}

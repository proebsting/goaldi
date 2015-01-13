//  execute.go -- the interpreter main loop

package main

import (
	"fmt"
	g "goaldi"
	"reflect"
)

//  coexecute wraps an execute call to catch a panic in a co-expression
func coexecute(f *pr_frame, label string) (g.Value, *g.Closure) {
	defer g.Catcher(f.env)
	return execute(f, label)
}

//  execute IR instructions for procedure or co-expression
func execute(f *pr_frame, label string) (rv g.Value, rc *g.Closure) {

	// set up traceback recovery
	defer func() {
		if p := recover(); p != nil {
			if f.onerr != nil {
				// find true panic value hiding under traceback info
				arglist := []g.Value{g.Cause(p)}
				// call recovery procedure and return its result
				rv, _ = f.onerr.Call(f.env, arglist, []string{})
				rc = nil
			} else {
				// add traceback information and re-throw exception
				panic(g.Catch(p,
					[]g.Value{f.offv}, f.coord, f.info.name, f.args))
			}
		}
	}()

	// create re-entrant interpreter
	f.temps = make(map[string]interface{}) // each cx needs own copy
	var self *g.Closure
	self = &g.Closure{func() (g.Value, *g.Closure) {

		// interpret the IR code
		for {
			if opt_trace {
				fmt.Printf("[%d] %s:\n", f.env.ThreadID, label)
			}
			ilist := f.info.insns[label] // look up label
			if len(ilist) == 0 {
				panic(g.Malfunction("No instructions for IR label: " + label))
			}
		Chunk:
			for _, insn := range ilist { // execute insns in chunk
				if opt_trace {
					t := fmt.Sprintf("%T", insn)[8:]
					fmt.Printf("[%d]    %s %v\n", f.env.ThreadID, t, insn)
				}
				f.coord = "" //#%#% prudent, but s/n/b needed
				f.offv = nil //#%#% prudent, but s/n/b needed
				switch i := insn.(type) {
				default: // incl ScanSwap, Assign, Deref, Unreachable
					panic(g.Malfunction(fmt.Sprintf(
						"Unrecognized interpreter instruction: %#v", i)))
				case ir_NoOp:
					// nothing to do
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
				case ir_Catch:
					f.offv = g.Deref(f.temps[i.Fn])
					f.onerr = f.offv.(*g.VProcedure)
				case ir_Create:
					fnew := newframe(f)
					fnew.cxout = g.NewChannel(0)
					fnew.env = g.NewEnv(f.env)
					fnew.coord = i.Coord
					if i.Lhs != "" {
						f.temps[i.Lhs] = fnew.cxout
					}
					go coexecute(fnew, i.CoexpLabel)
				case ir_Select:
					label = irSelect(f, i)
					break Chunk
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
						panic(g.NewExn("Unrecognized dynamic variable",
							"%"+i.Name))
					}
					if i.Lhs != "" {
						f.temps[i.Lhs] = v
					}
				case ir_NilLit:
					f.temps[i.Lhs] = g.NilValue
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
						panic(g.NewExn("Undeclared identifier", i.Name))
					}
					f.temps[i.Lhs] = v
				case ir_EnterScope:
					for _, name := range i.NameList {
						f.vars[name] = g.Trapped(g.NewVariable(g.NilValue))
					}
				case ir_ExitScope:
					for _, name := range i.NameList {
						f.vars[name] = nil // allow garbage collection
					}
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
					for _, s := range i.LabelList {
						if s == label {
							break Chunk
						}
					}
					panic(g.Malfunction(
						"ir_IndirectGoto: label not in list: " + label))
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
					}
				case ir_Call:
					f.coord = i.Coord
					proc := g.Deref(f.temps[i.Fn].(g.Value))
					arglist := getArgs(f, 0, i.ArgList)
					f.offv = proc
					v, c := proc.(g.ICall).Call(f.env, arglist, i.NameList)
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

//  irSelect -- execute select statement, returning label of chosen body
//  #%#% most of this should be moved to runtime package
//  #%#% (set up some data structures here and call that)
func irSelect(f *pr_frame, ir ir_Select) string {

	// set up data structures for reflect.Select
	n := len(ir.CaseList)
	cases := make([]reflect.SelectCase, n, n+1)
	seenDefault := false
	for i, sc := range ir.CaseList {
		f.coord = sc.Coord
		switch sc.Kind {
		case "send":
			ch := g.Deref(f.temps[sc.Lhs])
			v := g.Deref(f.temps[sc.Rhs])
			if _, ok := ch.(g.VChannel); !ok {
				// not a Goaldi channel; convert data value to best Go type
				v = g.Export(v)
			}
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: channelValue(ch),
				Send: reflect.ValueOf(v)}
		case "receive":
			ch := g.Deref(f.temps[sc.Rhs])
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: channelValue(ch)}
		case "default":
			cases[i] = reflect.SelectCase{
				Dir: reflect.SelectDefault}
			seenDefault = true
		default:
			panic(g.Malfunction("Bad selectcase kind: " + sc.Kind))
		}
	}
	if !seenDefault {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	}
	// repeat until we get anything other than a read on a closed channel
	for {
		f.coord = ir.Coord
		// call select through the reflection interface
		i, v, recvOK := reflect.Select(cases)
		// select has returned, having chosen case i
		if i == n {
			// this is the default case we added, because there was none
			return ir.FailLabel // so the select expression fails
		}
		sc := ir.CaseList[i]
		f.coord = sc.Coord
		if sc.Kind == "receive" {
			if recvOK {
				// assign received value before executing body
				f.temps[sc.Lhs].(g.IVariable).Assign(g.Import(v.Interface()))
			} else {
				// a closed channel was selected
				cases[i].Chan = hungChannel // disable this case
				continue                    // and retry
			}
		}
		return sc.BodyLabel // all scenarios except receive from closed channel
	}
}

//  get channel value and validate
func channelValue(ch g.Value) reflect.Value {
	cv := reflect.ValueOf(ch)
	if cv.Kind() != reflect.Chan {
		panic(g.NewExn("Not a channel", ch))
	}
	return cv
}

//  used for disabling one branch of a select
var hungChannel = reflect.ValueOf(make(chan interface{}))

//  execute.go -- the interpreter main loop

package main

import (
	"fmt"
	"goaldi/ir"
	g "goaldi/runtime"
	"reflect"
)

//  coexecute wraps an execute call to catch a panic in a co-expression
func coexecute(f *pr_frame, label string) (g.Value, *g.Closure) {
	defer g.Catcher(f.env)
	return execute(f, label)
}

//  execute IR instructions for procedure or co-expression
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
	f.temps = make(map[string]interface{}) // each cx needs own copy
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

		// interpret the IR code
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
			label = "UNSET"              // should never see this
			for _, insn := range ilist { // execute insns in chunk
				if opt_trace {
					t := fmt.Sprintf("%T", insn)[6:]
					fmt.Printf("[%d]    %s %v\n", f.env.ThreadID, t, insn)
				}
				f.coord = "" //#%#% prudent, but s/n/b needed
				f.offv = nil //#%#% prudent, but s/n/b needed
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
					if i.Lhs != "" {
						f.temps[i.Lhs] = f.onerr
					}
				case ir.Ir_Create:
					fnew := newframe(f)
					fnew.cxout = g.NewChannel(0)
					fnew.env = g.NewEnv(f.env)
					fnew.env.ThreadID = <-g.TID
					fnew.coord = i.Coord
					if i.Lhs != "" {
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
					v := e.Lookup(i.Name)         // look up name
					if v == nil {
						panic(g.NewExn("Unrecognized dynamic variable",
							"%"+i.Name))
					}
					if i.Rval != "" { // if an rval is required
						v = g.Deref(v) // then make sure we have one
					}
					if i.Lhs != "" {
						f.temps[i.Lhs] = v
					}
				case ir.Ir_NilLit:
					f.temps[i.Lhs] = g.NilValue
				case ir.Ir_IntLit:
					f.temps[i.Lhs] = g.NewString(i.Val).ToNumber()
				case ir.Ir_RealLit:
					f.temps[i.Lhs] = g.NewString(i.Val).ToNumber()
				case ir.Ir_StrLit:
					f.temps[i.Lhs] = g.NewString(i.Val)
				case ir.Ir_MakeList:
					n := len(i.ValueList)
					a := make([]g.Value, n)
					for j, v := range i.ValueList {
						a[j] = g.Deref(f.temps[v.(string)])
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
						for _, name := range i.DynamicList { // init dynamics
							e.VarMap[name] = g.NewVariable(g.NilValue)
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
					label = i.TargetTmpLabel
					label = f.temps[label].(string)
					for _, s := range i.LabelList {
						if s == label {
							continue NextChunk
						}
					}
					panic(g.Malfunction(
						"IndirectGoto: unlisted label: " + label))
				case ir.Ir_MakeClosure:
					//#%#% potential later optimization:
					//#%#% only pass in *referenced* variables
					//#%#% so that the remainder can get garbage collected
					f.temps[i.Lhs] = irProcedure(ProcTable[i.Name], f.vars)
				case ir.Ir_OpFunction:
					f.coord = i.Coord
					v, c := operator(f.env, f, &i)
					if i.Rval != "" && v != nil { // if an rval is required
						v = g.Deref(v) // then make sure we have one
						// note v can be set nil by failing Deref
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
						if i.Lhs != "" {
							f.temps[i.Lhs] = v
						}
					}
				case ir.Ir_Call:
					f.coord = i.Coord
					proc := g.Deref(f.temps[i.Fn].(g.Value))
					arglist := getArgs(f, 0, i.ArgList)
					f.offv = proc
					e := f.vars[i.Scope].(*g.Env) // get correct environment
					v, c := proc.(g.ICall).Call(e, arglist, i.NameList)
					if v != nil {
						if i.Lhs != "" {
							f.temps[i.Lhs] = v
						}
						if i.Lhsclosure != "" {
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
						if i.Lhs != "" {
							f.temps[i.Lhs] = v
						}
						if i.Lhsclosure != "" {
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
		if a == nil {
			panic(g.Malfunction(
				fmt.Sprintf("Go nil in getArgs(): i=%d %#v", i, arglist[i])))
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
func irSelect(f *pr_frame, irs ir.Ir_Select) string {

	// set up data structures for reflect.Select
	n := len(irs.CaseList)
	cases := make([]reflect.SelectCase, n, n+1)
	seenDefault := false
	for i, sc := range irs.CaseList {
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
			panic(g.Malfunction("Bad SelectCase kind: " + sc.Kind))
		}
	}
	if !seenDefault {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	}
	// repeat until we get anything other than a read on a closed channel
	for {
		f.coord = irs.Coord
		// call select through the reflection interface
		i, v, recvOK := reflect.Select(cases)
		// select has returned, having chosen case i
		if i == n {
			// this is the default case we added, because there was none
			return irs.FailLabel // so the select expression fails
		}
		sc := irs.CaseList[i]
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

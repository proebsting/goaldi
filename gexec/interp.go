//  interp.go -- the interpreter main loop

package main

import (
	"fmt"
	g "goaldi"
)

//  procedure frame
type pr_frame struct {
	env   *g.Env                 // dynamic execution environment
	info  *pr_Info               // static procedure information
	args  []g.Value              // arglist as called
	vars  map[string]interface{} // variables
	temps map[string]interface{} // temporaries
	coord string                 // last known source location
	offv  g.Value                // offending value for traceback
	cxout g.VChannel             // co-expression output pipe
}

//  newframe(f) -- duplicate a procedure frame
func newframe(f *pr_frame) *pr_frame {
	fnew := &pr_frame{}
	*fnew = *f
	fnew.vars = make(map[string]interface{})
	for k, v := range f.vars {
		fnew.vars[k] = v
	}
	// make new copies of all parameter values
	for _, name := range f.info.params {
		fnew.vars[name] = g.Trapped(g.NewVariable(g.Deref(f.vars[name])))
	}
	// make new copies of all locals (n.b. does not include statics)
	for _, name := range f.info.locals {
		fnew.vars[name] = g.Trapped(g.NewVariable(g.Deref(f.vars[name])))
	}
	return fnew
}

//  duplvars(a) -- duplicate a list of (trapped) variables or parameters
func duplvars(a []g.Value) []g.Value {
	b := make([]g.Value, len(a))
	for i, x := range a {
		b[i] = g.Trapped(g.NewVariable(g.Deref(x)))
	}
	return b
}

//  interp -- interpret one procedure
func interp(env *g.Env, pr *pr_Info, outer map[string]interface{},
	args ...g.Value) (g.Value, *g.Closure) {

	if opt_trace {
		fmt.Printf("[%d] P: %s\n", env.ThreadID, pr.name)
	}

	// initialize procedure frame
	var f pr_frame
	f.env = env
	f.info = pr
	f.args = args

	// initialize variable dictionary with inherited variables;
	// any of these may be subsequently hidden (replaced)
	f.vars = make(map[string]interface{})
	for k, v := range outer {
		f.vars[k] = v
	}

	// add static variables
	for k, v := range pr.statics {
		f.vars[k] = v
	}

	// initialize parameters
	for i, name := range pr.params {
		if i < len(args) {
			f.vars[name] = g.Trapped(g.NewVariable(args[i]))
		} else {
			f.vars[name] = g.Trapped(g.NewVariable(g.NilValue))
		}
	}

	//  handle variadic procedure
	if pr.variadic {
		n := len(pr.params) - 1
		vp := new(g.Value)
		if len(args) < n {
			*vp = g.NewList(0, nil)
		} else {
			vals := make([]g.Value, len(args)-n)
			copy(vals, args[n:])
			*vp = g.InitList(vals)
		}
		f.vars[pr.params[n]] = g.Trapped(vp)
	}

	// initialize locals
	for _, name := range pr.locals {
		f.vars[name] = g.Trapped(g.NewVariable(g.NilValue))
	}

	// execute the IR code
	return execute(&f, pr.ir.CodeStart)
}

//  coexecute wraps an execute call to catch a panic in a co-expression
func coexecute(f *pr_frame, label string) (g.Value, *g.Closure) {
	defer g.Catcher(f.env)
	return execute(f, label)
}

//  execute IR code for procedure or co-expression
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
					v, c := opFunc(f.env, f, &i)
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

//  opFunc -- implement operator function
func opFunc(env *g.Env, f *pr_frame, i *ir_OpFunction) (g.Value, *g.Closure) {
	op := string('0'+len(i.ArgList)) + i.Fn
	a := getArgs(f, nonDeref[op], i.ArgList)
	f.offv = a[0]        // save offending value
	var lval g.IVariable // lvalue for some operators
	if i.Rval == "" {    // set it if warranted
		lval, _ = a[0].(g.IVariable) // and if possible
	}

	switch op {
	default:
		panic(&g.RunErr{"Unimplemented operator", g.NewString(op)})

	// fundamental operations
	case "1.":
		return g.Return(a[0]) // was dereferenced by getArgs
	case "1#":
		// means e > 0, used with x \ e
		return g.ZERO.NumLT(a[0])
	case "1/":
		v := g.Deref(a[0])
		if v == g.NilValue {
			return g.Return(a[0]) // NOT dereferenced!
		} else {
			return g.Fail()
		}
	case "1\\":
		v := g.Deref(a[0])
		if v != g.NilValue {
			return g.Return(a[0]) // NOT dereferenced!
		} else {
			return g.Fail()
		}
	case "2===":
		return g.Identical(a[0], a[1]), nil
	case "2~===":
		return g.NotIdentical(a[0], a[1]), nil

	// assignment
	case "2:=":
		return a[0].(g.IVariable).Assign(a[1]), nil
	case "2<-":
		return g.RevAssign(a[0], a[1])
	case "2:=:":
		return g.Return(g.Swap(a[0], a[1]))
	case "2<->":
		return g.RevSwap(a[0], a[1])

	// multi-type operations
	case "1*":
		return g.Size(a[0]), nil
	case "1@", "2@": //#%#% 1@(x) is passed as 2@(x,null)
		return g.Take(a[0]), nil
	case "1?":
		return g.Choose(lval, g.Deref(a[0])), nil
	case "1!":
		return g.Dispense(lval, g.Deref(a[0]))
	case "2[]":
		return g.Index(lval, g.Deref(a[0]), a[1]), nil
	case "3[:]":
		return g.Deref(a[0]).(g.ISlice).Slice(lval, a[1], a[2]), nil
	case "3[+:]":
		return deltaSlice(lval, a, +1)
	case "3[-:]":
		return deltaSlice(lval, a, -1)

	// miscellaneous operations
	case "2@:":
		return g.Send(a[0], a[1]), nil
	case "2!":
		return a[0].(g.ICall).Call(env, a[1].(*g.VList).Export().([]g.Value)...)
	case "2put":
		return a[0].(g.IListPut).ListPut(a[1]), nil
	case "2|||":
		return a[0].(g.IListCat).ListCat(a[1]), nil

	// string operations
	case "2||":
		return a[0].(g.IConcat).Concat(a[1]), nil
	case "2<<":
		return a[0].(g.IStrLT).StrLT(a[1]), nil
	case "2<<=":
		return a[0].(g.IStrLE).StrLE(a[1]), nil
	case "2==":
		return a[0].(g.IStrEQ).StrEQ(a[1]), nil
	case "2~==":
		return a[0].(g.IStrNE).StrNE(a[1]), nil
	case "2>>=":
		return a[0].(g.IStrGE).StrGE(a[1]), nil
	case "2>>":
		return a[0].(g.IStrGT).StrGT(a[1]), nil

	// numeric operations
	case "1+":
		return a[0].(g.INumerate).Numerate(), nil
	case "1-":
		return a[0].(g.INegate).Negate(), nil
	case "2+":
		return a[0].(g.IAdd).Add(a[1]), nil
	case "2-":
		return a[0].(g.ISub).Sub(a[1]), nil
	case "2*":
		return a[0].(g.IMul).Mul(a[1]), nil
	case "2/":
		return a[0].(g.IDiv).Div(a[1]), nil
	case "2//":
		return a[0].(g.IDivt).Divt(a[1]), nil
	case "2%":
		return a[0].(g.IMod).Mod(a[1]), nil
	case "2^":
		return a[0].(g.IPower).Power(a[1]), nil
	case "2<":
		return a[0].(g.INumLT).NumLT(a[1])
	case "2<=":
		return a[0].(g.INumLE).NumLE(a[1])
	case "2=":
		return a[0].(g.INumEQ).NumEQ(a[1])
	case "2~=":
		return a[0].(g.INumNE).NumNE(a[1])
	case "2>=":
		return a[0].(g.INumGE).NumGE(a[1])
	case "2>":
		return a[0].(g.INumGT).NumGT(a[1])
	case "3...":
		return g.ToBy(a[0], a[1], a[2])
	}
}

var nonDeref = make(map[string]int)

func init() {
	nonDeref["1/"] = 1
	nonDeref["1\\"] = 1
	nonDeref["1?"] = 1
	nonDeref["1!"] = 1
	nonDeref["2:="] = 1
	nonDeref["2<-"] = 1
	nonDeref["2:=:"] = 2
	nonDeref["2<->"] = 2
	nonDeref["2[]"] = 1
	nonDeref["3[:]"] = 1
	nonDeref["3[+:]"] = 1
	nonDeref["3[-:]"] = 1
}

//  deltaSlice handles x[i+:k] or x[i-:k] by calling x[i:j]
func deltaSlice(lval g.IVariable, a []g.Value, sign int) (g.Value, *g.Closure) {
	x := g.Deref(a[0]).(g.ISlice)
	i := int(a[1].(g.Numerable).ToNumber().Val())
	j := i + sign*int(a[2].(g.Numerable).ToNumber().Val())
	if (i > 0 && j <= 0) || (i <= 0 && j > 0) { // if wraparound
		return nil, nil // fail
	}
	return x.Slice(lval, g.NewNumber(float64(i)), g.NewNumber(float64(j))), nil
}

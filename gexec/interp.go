//  interp.go -- the interpreter main loop

package main

import (
	"fmt"
	g "goaldi"
)

//  procedure frame
type pr_frame struct {
	env    *g.Env             // dynamic execution enviromment
	info   *pr_Info           // static procedure information
	params []g.Value          // parameters
	locals []g.Value          // locals
	temps  map[string]g.Value // temporaries
	coord  string             // last known source location
	offv   g.Value            // offending value for traceback
}

//  catchf -- annotate a panic value with procedure frame information
func catchf(p interface{}, f *pr_frame, args []g.Value) *g.CallFrame {
	return g.Catch(p, f.offv, f.coord, f.info.name, args)
}

//  interp -- interpret one procedure
func interp(env *g.Env, pr *pr_Info, args ...g.Value) (g.Value, *g.Closure) {

	if opt_trace {
		fmt.Printf("P: %s\n", pr.name)
	}

	// initialize procedure frame: params, locals, temps
	var f pr_frame
	f.env = env
	f.info = pr
	f.temps = make(map[string]g.Value)
	f.params = make([]g.Value, pr.nparams, pr.nparams)
	f.locals = make([]g.Value, pr.nlocals, pr.nlocals)
	for i := 0; i < len(f.params); i++ {
		if i < len(args) {
			f.params[i] = args[i]
		} else {
			f.params[i] = g.NewNil()
		}
	}
	for i := 0; i < len(f.locals); i++ {
		f.locals[i] = g.NewNil()
	}

	// set starting point
	label := pr.ir.CodeStart

	// create re-entrant interpreter
	var self *g.Closure
	self = &g.Closure{func() (g.Value, *g.Closure) {

		// set up tracback recovery
		defer func() {
			if p := recover(); p != nil {
				panic(catchf(p, &f, args))
			}
		}()

		// interpret the IR code
		for {
			if opt_trace {
				fmt.Printf("L: %s\n", label)
			}
			ilist := pr.insns[label] // look up label
		Chunk:
			for _, insn := range ilist { // execute insns in chunk
				if opt_trace {
					fmt.Printf("I: %T %v\n", insn, insn)
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
					v := f.temps[i.Expr]
					if i.ResumeLabel == "" {
						return v, nil
					} else {
						label = i.ResumeLabel
						return v, self
					}
				case ir_Key:
					v := keyword(i.Name)
					//#%#% ignoring failure and FailLabel
					if i.Lhs != "" {
						f.temps[i.Lhs] = v
					}
				case ir_IntLit:
					f.temps[i.Lhs] =
						g.NewString(i.Val).ToNumber()
				case ir_RealLit:
					f.temps[i.Lhs] =
						g.NewString(i.Val).ToNumber()
				case ir_StrLit:
					f.temps[i.Lhs] =
						g.NewString(i.Val)
				case ir_Var:
					v := pr.dict[i.Name]
					switch t := v.(type) {
					case pr_local:
						v = g.Trapped(&f.locals[int(t)])
					case pr_param:
						v = g.Trapped(&f.params[int(t)])
					case nil:
						panic("nil in ir_Var; undeclared?")
					default:
						// global or static: already trapped
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
				case ir_OpFunction:
					f.coord = i.Coord
					v, c := opFunc(&f, &i)
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
				case ir_Call:
					f.coord = i.Coord
					proc := f.temps[i.Fn]
					argl := getArgs(&f, 0, i.ArgList)
					f.offv = proc
					v, c := proc.(g.ICall).Call(env, argl...)
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
	argl := make([]g.Value, n, n)
	for i, a := range arglist {
		switch t := a.(type) {
		case string:
			a = f.temps[t]
		default:
			// nothing to do: use entry as is
		}
		if i < nd {
			argl[i] = a
		} else {
			argl[i] = g.Deref(a)
		}
	}
	return argl
}

//  opFunc -- implement operator function
func opFunc(f *pr_frame, i *ir_OpFunction) (g.Value, *g.Closure) {
	op := string('0'+len(i.ArgList)) + i.Fn
	a := getArgs(f, nonDeref[op], i.ArgList)
	f.offv = a[0]

	switch op {
	default:
		panic(&g.RunErr{"Unimplemented operator", g.NewString(op)})

	// fundamental operations
	case "1/":
		v := g.Deref(a[0])
		if v == g.NilVal {
			return g.Return(a[0]) // NOT dereferenced!
		} else {
			return g.Fail()
		}
	case "1\\":
		v := g.Deref(a[0])
		if v != g.NilVal {
			return g.Return(a[0]) // NOT dereferenced!
		} else {
			return g.Fail()
		}
	case "2===":
		return g.Identical(a[0], a[1]), nil
	case "2~===":
		return g.NotIdentical(a[0], a[1]), nil
	case "2:=":
		return a[0].(g.IVariable).Assign(a[1]), nil

	// multi-type operations
	case "1*":
		return a[0].(g.ISize).Size(), nil
	case "1?":
		v := g.Deref(a[0])
		if i.Rval != "" {
			// produce value
			return v.(g.IChoose).Choose(nil), nil
		} else {
			// produce variable
			return v.(g.IChoose).Choose(a[0].(g.IVariable)), nil
		}
	case "1!":
		v := g.Deref(a[0])
		if i.Rval != "" {
			// generate values
			return v.(g.IDispense).Dispense(nil)
		} else {
			// generate variables
			return v.(g.IDispense).Dispense(a[0].(g.IVariable))
		}
	case "2[]":
		v := g.Deref(a[0])
		if i.Rval != "" {
			// return slice value
			return v.(g.IIndex).Index(nil, a[1]), nil
		} else {
			// return variable slice
			return v.(g.IIndex).Index(a[0].(g.IVariable), a[1]), nil
		}
	case "3[:]":
		v := g.Deref(a[0])
		if i.Rval != "" {
			// return slice value
			return v.(g.ISlice).Slice(nil, a[1], a[2]), nil
		} else {
			// return variable slice
			return v.(g.ISlice).Slice(a[0].(g.IVariable), a[1], a[2]), nil
		}

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
}

//  keyword -- return keyword value
//  #%#% cannot handle generator keywords
//  #%#% currently only handles &null
func keyword(name string) g.Value {
	switch name {
	default:
		panic(&g.RunErr{"Unrecognized keyword", g.NewString(name)})
	case "null":
		return g.NilVal
	}
}

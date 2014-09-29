//  interp.go -- the interpreter main loop

package main

import (
	"fmt"
	g "goaldi"
)

//  procedure frame
type pr_frame struct {
	env    *g.Env
	info   *pr_Info
	params []g.Value
	locals []g.Value
	temps  map[string]g.Value
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

	//#%#%#% defer recover ...

	// interpret the IR code
	label := pr.ir.CodeStart.Value
	for {
		if opt_trace {
			fmt.Printf("L: %s\n", label)
		}
		ilist := pr.insns[label]     // look up label
		for _, insn := range ilist { // execute insns in chunk
			if opt_trace {
				fmt.Printf("I: %T %v\n", insn, insn)
			}
			switch i := insn.(type) {
			default:
				panic(i) // unrecognized or unimplemented
			case ir_Fail:
				return nil, nil
			case ir_IntLit:
				f.temps[i.Lhs.Name] =
					g.NewString(i.Val).ToNumber()
			case ir_RealLit:
				f.temps[i.Lhs.Name] =
					g.NewString(i.Val).ToNumber()
			case ir_StrLit:
				f.temps[i.Lhs.Name] =
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
				f.temps[i.Lhs.Name] = v
			case ir_OpFunction:
				n := len(i.ArgList)
				argl := make([]g.Value, n, n)
				for j, a := range i.ArgList {
					switch t := a.(type) {
					case ir_Tmp:
						argl[j] = f.temps[t.Name]
					default:
						argl[j] = g.Deref(a) //#%#%
					}
				}
				v, c := opFunc(i.Fn, argl)
				if v == nil && i.FailLabel.Value != "" {
					label = i.FailLabel.Value
					break
				}
				f.temps[i.Lhs.Name] = v
				f.temps[i.Lhsclosure.Name] = c
			case ir_Call:
				//#%#% combine shared code with OpFunction
				proc := f.temps[i.Fn.Name]
				n := len(i.ArgList)
				argl := make([]g.Value, n, n)
				for j, a := range i.ArgList {
					switch t := a.(type) {
					case ir_Tmp:
						argl[j] = f.temps[t.Name]
					default:
						argl[j] = g.Deref(a) //#%#%
					}
				}
				f.temps[i.Lhs.Name], f.temps[i.Lhsclosure.Name] =
					proc.(g.ICall).Call(env, argl...)
			}
		}
	}
	return nil, nil
}

//  opFunc -- implement operator function
func opFunc(o *ir_operator, a []g.Value) (g.Value, *g.Closure) {
	op := o.Arity + o.Name
	switch op {
	default:
		panic("unimplemented operator: " + op)
	case "2+":
		return a[0].(g.IAdd).Add(a[1]), nil
	}
}

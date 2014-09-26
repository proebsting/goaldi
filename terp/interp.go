//  interp.go -- the interpreter main loop

package main

import (
	"fmt"
	g "goaldi"
)

//  interp -- interpret one procedure
func interp(env *g.Env, pr *pr_Info, args ...g.Value) (g.Value, *g.Closure) {

	if opt_trace {
		fmt.Printf("P: %s\n", pr.name)
	}

	// initialize params, locals, temps
	temps := make(map[string]g.Value)
	params := make([]g.Value, pr.nparams, pr.nparams)
	locals := make([]g.Value, pr.nlocals, pr.nlocals)
	for i := 0; i < len(params); i++ {
		if i < len(args) {
			params[i] = args[i]
		} else {
			params[i] = g.NewNil()
		}
	}
	for i := 0; i < len(locals); i++ {
		locals[i] = g.NewNil()
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
			case ir_Var:
				v := pr.dict[i.Name]
				switch t := v.(type) {
				case pr_local:
					v = g.Trapped(&locals[int(t)])
				case pr_param:
					v = g.Trapped(&params[int(t)])
				case nil:
					panic("nil in ir_Var; undeclared?")
				default:
					// global or static: already trapped
				}
				temps[i.Lhs.Name] = v
			case ir_IntLit:
				temps[i.Lhs.Name] = g.NewString(i.Val).ToNumber()
			case ir_RealLit:
				temps[i.Lhs.Name] = g.NewString(i.Val).ToNumber()
			case ir_StrLit:
				temps[i.Lhs.Name] = g.NewString(i.Val)
			case ir_Call:
				proc := temps[i.Fn.Name]
				argl := make([]g.Value, 0)
				for _, a := range i.ArgList {
					switch t := a.(type) {
					case ir_Tmp:
						argl = append(argl, temps[t.Name])
					default: //#%#%#%
						argl = append(argl, g.Deref(a))
					}
				}
				temps[i.Lhs.Name], temps[i.Lhsclosure.Name] =
					proc.(g.ICall).Call(env, argl...)
			case ir_Fail:
				return nil, nil
			}
		}
	}
	_ = temps
	return nil, nil
}

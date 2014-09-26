//  interp.go -- the interpreter main loop

package main

import (
	"fmt"
	g "goaldi"
)

func interp(env *g.Env, pr *pr_Info, args ...g.Value) (g.Value, *g.Closure) {
	label := pr.ir.CodeStart.Value

	for {
		ilist := pr.insns[label]
		for _, insn := range ilist {
			label := interp(insn)
			fmt.Printf("X: %T %v\n", insn, insn)
			switch i := insn.(type) {
			default:
				panic(i)
			case ir_Var:
			case ir_StrLit:
			case ir_Call:
			case ir_Fail:
				return nil, nil
			}
		}
	}
}

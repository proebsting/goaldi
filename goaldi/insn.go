//  insn.go -- additional runtime instructions not part of the IR
//
//  These instructions replace others in a just-in-time fashion
//  when those other instructions are encountered during execution.

package main

import (
	g "goaldi/runtime"
)

//  Ins_Literal replaces Ir_NilLit, Ir_IntLit, Ir_RealLit, Ir_StrLit
type Ins_Literal struct {
	Lhs   string
	Value g.Value
}

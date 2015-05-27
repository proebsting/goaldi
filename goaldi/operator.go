//  operator.go -- interpret a unary or binary operator

package main

import (
	"fmt"
	"goaldi/ir"
	g "goaldi/runtime"
)

//  iOperator flag bits
const (
	rflag = 1 << iota // rvalue result wanted
	v0                // arg0 used as value
	v1                // arg1 used as value
	v2                // arg2 used as value
	VAR0              // arg0 used as variable
	VAR1              // arg1 used as variable
)

//  iOperator instruction
type iOperator struct {
	OpCode     uint16 // integer opcode for dispatching
	Flags      uint16 // flags
	Coord      string
	Lhs        string // may be nil
	Lhsclosure string // may be nil
	Fn         string // opstring for tracing purposes
	Arg0       interface{}
	Arg1       interface{}
	Arg2       interface{}
	FailLabel  string // may be nil
}

//  getOperator(i) returns the iOperator version of IR instruction i
func getOperator(i *ir.Ir_OpFunction) *iOperator {
	// get a new *copy* (by value) of this insn's table entry
	// with opcode and flags set
	s := string('0'+len(i.ArgList)) + i.Fn // e.g. "2:="
	f := opTable[s]                        // get new COPY of table entry
	if f.OpCode == oInvalid {              // if no entry found
		panic(g.Malfunction("No opcode found for " + s))
	}
	// fill in the fields for this particular instruction
	f.Fn = i.Fn
	f.Coord = i.Coord
	f.Lhs = i.Lhs
	f.Lhsclosure = i.Lhsclosure
	f.FailLabel = i.FailLabel
	// compute rval string to flag
	if i.Rval != "" {
		f.Flags |= rflag
	}
	// flatten argument array
	if len(i.ArgList) > 0 {
		f.Arg0 = i.ArgList[0]
	}
	if len(i.ArgList) > 1 {
		f.Arg1 = i.ArgList[1]
	}
	if len(i.ArgList) > 2 {
		f.Arg2 = i.ArgList[2]
	}
	// return faster version of instruction
	return &f
}

//  operate() executes an iOperator insn corresponding to a Goaldi operator
func operate(env *g.Env, f *pr_frame, i *iOperator) (g.Value, *g.Closure) {

	// load and possibly dereference arguments
	var lval, arg0, arg1, arg2 g.Value
	arg0 = argval(f, i.Arg0, i.Flags&VAR0)
	if (i.Flags & (v1 | VAR1)) != 0 {
		arg1 = argval(f, i.Arg1, i.Flags&VAR1)
	}
	if (i.Flags & v2) != 0 {
		arg2 = argval(f, i.Arg2, 0)
	}
	if (i.Flags & rflag) == 0 {
		lval = arg0
	}

	switch i.OpCode {
	default:
		panic(g.Malfunction(fmt.Sprintf("Unknown opcode %d", i.OpCode)))

	// fundamental operations
	case oIsNull:
		if g.Deref(arg0) == g.NilValue {
			return g.Return(arg0) // NOT dereferenced!
		} else {
			return g.Fail()
		}
	case oNotNull:
		if g.Deref(arg0) != g.NilValue {
			return g.Return(arg0) // NOT dereferenced!
		} else {
			return g.Fail()
		}
	case oIdentical:
		return g.Identical(arg0, arg1), nil
	case oNotIdentical:
		return g.NotIdentical(arg0, arg1), nil

	// assignment
	case oAssign:
		return arg0.(g.IVariable).Assign(arg1), nil
	case oRevAssign:
		return g.RevAssign(arg0, arg1)
	case oSwap:
		return g.Return(g.Swap(arg0, arg1))
	case oRevSwap:
		return g.RevSwap(arg0, arg1)

	// numeric operations
	case oLimit: // means e > 0, used with x \ e
		return g.ZERO.NumLT(arg0)
	case oNumerate:
		return arg0.(g.INumerate).Numerate(), nil
	case oNegate:
		return arg0.(g.INegate).Negate(), nil
	case oAdd:
		return arg0.(g.IAdd).Add(arg1), nil
	case oSub:
		return arg0.(g.ISub).Sub(arg1), nil
	case oMul:
		return arg0.(g.IMul).Mul(arg1), nil
	case oDiv:
		return arg0.(g.IDiv).Div(arg1), nil
	case oDivt:
		return arg0.(g.IDivt).Divt(arg1), nil
	case oMod:
		return arg0.(g.IMod).Mod(arg1), nil
	case oPower:
		return arg0.(g.IPower).Power(arg1), nil
	case oNumLT:
		return arg0.(g.INumLT).NumLT(arg1)
	case oNumLE:
		return arg0.(g.INumLE).NumLE(arg1)
	case oNumEQ:
		return arg0.(g.INumEQ).NumEQ(arg1)
	case oNumNE:
		return arg0.(g.INumNE).NumNE(arg1)
	case oNumGE:
		return arg0.(g.INumGE).NumGE(arg1)
	case oNumGT:
		return arg0.(g.INumGT).NumGT(arg1)
	case oToBy:
		return g.ToBy(arg0, arg1, arg2)

	// string operations
	case oConcat:
		return arg0.(g.IConcat).Concat(arg1), nil
	case oStrLT:
		return arg0.(g.IStrLT).StrLT(arg1), nil
	case oStrLE:
		return arg0.(g.IStrLE).StrLE(arg1), nil
	case oStrEQ:
		return arg0.(g.IStrEQ).StrEQ(arg1), nil
	case oStrNE:
		return arg0.(g.IStrNE).StrNE(arg1), nil
	case oStrGE:
		return arg0.(g.IStrGE).StrGE(arg1), nil
	case oStrGT:
		return arg0.(g.IStrGT).StrGT(arg1), nil

	// set operations
	case oUnion:
		return arg0.(g.IUnion).Union(arg1), nil
	case oSetDiff:
		return arg0.(g.ISetDiff).SetDiff(arg1), nil
	case oIntersect:
		return arg0.(g.IIntersect).Intersect(arg1), nil

	// indexing operations
	case oIndex:
		return g.Index(lval, g.Deref(arg0), arg1), nil
	case oSlice:
		return g.Deref(arg0).(g.ISlice).Slice(lval, arg1, arg2), nil
	case oSlicePlus:
		return deltaSlice(lval, arg0, arg1, arg2, +1)
	case oSliceMinus:
		return deltaSlice(lval, arg0, arg1, arg2, -1)

	// structure operations
	case oSize:
		return g.Size(arg0), nil
	case oTake:
		// always pass lval; ignored by all except @s (take from string)
		return g.Take(arg0, g.Deref(arg0)), nil
	case oChoose:
		return g.Choose(lval, g.Deref(arg0)), nil
	case oDispense:
		return g.Dispense(lval, g.Deref(arg0))

	// miscellaneous operations
	case oSend:
		return g.Send(arg0, g.Deref(arg0), arg1), nil // lval for s@:x
	case oCall:
		arglist := arg1.(*g.VList).Export().([]g.Value)
		return arg0.(g.ICall).Call(env, arglist, []string{})
	case oListPut:
		return arg0.(g.IListPut).ListPut(arg1), nil
	case oListCat:
		return arg0.(g.IListCat).ListCat(arg1), nil

	}
}

//  argval retrieves an argument for an operator and possibly deferences it
func argval(f *pr_frame, v interface{}, isvar uint16) g.Value {
	v = f.temps[v.(string)]
	if isvar == 0 {
		v = g.Deref(v)
	}
	if v == nil {
		panic("Go nil in Operator/ArgVal")
	}
	return g.Value(v)
}

//  deltaSlice handles x[i+:k] or x[i-:k] by calling x[i:j]
func deltaSlice(lval g.Value, arg0 g.Value, arg1 g.Value, arg2 g.Value,
	sign int) (g.Value, *g.Closure) {
	x := g.Deref(arg0).(g.ISlice)
	i := int(g.FloatVal(arg1))
	j := i + sign*int(g.FloatVal(arg2))
	if (i > 0 && j <= 0) || (i <= 0 && j > 0) { // if wraparound
		return nil, nil // fail
	}
	return x.Slice(lval, g.NewNumber(float64(i)), g.NewNumber(float64(j))), nil
}

//  opTable maps Ir_OpFunction strings to iOperator opcodes
var opTable = make(map[string]iOperator)

//  codeTable maps opcodes back to strings for checking against duplicates
var codeTable = make(map[int]string)

//  defOp creates a skeleton iOperator for replacing an Ir_OpFunction
func defOp(opstring string, opcode int, flags int) {
	if opTable[opstring].OpCode != 0 {
		panic(g.Malfunction("Duplicate opstring: " + opstring))
	}
	if e := codeTable[opcode]; e != "" {
		panic(g.Malfunction(fmt.Sprintf(
			`Duplicate opcode %d for "%s" and "%s"`, opcode, e, opstring)))
	}
	codeTable[opcode] = opstring
	opTable[opstring] = iOperator{OpCode: uint16(opcode), Flags: uint16(flags)}
}

//  This initializer defines the mapping of strings to opcodes
//  with flags for variable and value usage.
func init() {
	defOp("1/", oIsNull, VAR0)
	defOp("1\\", oNotNull, VAR0)
	defOp("2===", oIdentical, v0+v1)
	defOp("2~===", oNotIdentical, v0+v1)
	defOp("2:=", oAssign, VAR0+v1)
	defOp("2<-", oRevAssign, VAR0+v1)
	defOp("2:=:", oSwap, VAR0+VAR1)
	defOp("2<->", oRevSwap, VAR0+VAR1)
	defOp("1#", oLimit, v0)
	defOp("1+", oNumerate, v0)
	defOp("1-", oNegate, v0)
	defOp("2+", oAdd, v0+v1)
	defOp("2-", oSub, v0+v1)
	defOp("2*", oMul, v0+v1)
	defOp("2/", oDiv, v0+v1)
	defOp("2//", oDivt, v0+v1)
	defOp("2%", oMod, v0+v1)
	defOp("2^", oPower, v0+v1)
	defOp("2<", oNumLT, v0+v1)
	defOp("2<=", oNumLE, v0+v1)
	defOp("2=", oNumEQ, v0+v1)
	defOp("2~=", oNumNE, v0+v1)
	defOp("2>=", oNumGE, v0+v1)
	defOp("2>", oNumGT, v0+v1)
	defOp("3...", oToBy, v0+v1+v2)
	defOp("2||", oConcat, v0+v1)
	defOp("2<<", oStrLT, v0+v1)
	defOp("2<<=", oStrLE, v0+v1)
	defOp("2==", oStrEQ, v0+v1)
	defOp("2~==", oStrNE, v0+v1)
	defOp("2>>=", oStrGE, v0+v1)
	defOp("2>>", oStrGT, v0+v1)
	defOp("2++", oUnion, v0+v1)
	defOp("2--", oSetDiff, v0+v1)
	defOp("2**", oIntersect, v0+v1)
	defOp("2[]", oIndex, VAR0+v1)
	defOp("3[:]", oSlice, VAR0+v1+v2)
	defOp("3[+:]", oSlicePlus, VAR0+v1+v2)
	defOp("3[-:]", oSliceMinus, VAR0+v1+v2)
	defOp("1*", oSize, v0)
	defOp("1@", oTake, VAR0)
	defOp("1?", oChoose, VAR0)
	defOp("1!", oDispense, VAR0)
	defOp("2@:", oSend, VAR0+v1)
	defOp("2!", oCall, v0+v1)
	defOp("2put", oListPut, v0+v1)
	defOp("2|||", oListCat, v0+v1)
	codeTable = nil // no longer needed; was just for sanity check
}

//  without no preprocessor, we must list the opcode names again to define them
const (
	oInvalid = iota
	oIsNull
	oNotNull
	oIdentical
	oNotIdentical
	oAssign
	oRevAssign
	oSwap
	oRevSwap
	oLimit
	oNumerate
	oNegate
	oAdd
	oSub
	oMul
	oDiv
	oDivt
	oMod
	oPower
	oNumLT
	oNumLE
	oNumEQ
	oNumNE
	oNumGE
	oNumGT
	oToBy
	oConcat
	oStrLT
	oStrLE
	oStrEQ
	oStrNE
	oStrGE
	oStrGT
	oUnion
	oSetDiff
	oIntersect
	oIndex
	oSlice
	oSlicePlus
	oSliceMinus
	oSize
	oTake
	oChoose
	oDispense
	oSend
	oCall
	oListPut
	oListCat
)

//  operator.go -- interpret a unary or binary operator

package main

import (
	"goaldi/ir"
	g "goaldi/runtime"
)

//  operator -- implement IR unary or binary operator
func operator(env *g.Env, f *pr_frame, i *ir.Ir_OpFunction) (g.Value, *g.Closure) {
	op := string('0'+len(i.ArgList)) + i.Fn
	a := getArgs(f, nonDeref[op], i.ArgList)
	f.offv = a[0] // save potential offending value
	var lval g.Value
	if i.Rval == "" {
		lval = a[0] // pass non-nil lvalue if result is not an rvalue
	}

	switch op {
	default:
		panic(g.Malfunction("Unimplemented operator: " + op))

	// fundamental operations
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
		arglist := a[1].(*g.VList).Export().([]g.Value)
		return a[0].(g.ICall).Call(env, arglist, []string{})
	case "2put":
		return a[0].(g.IListPut).ListPut(a[1]), nil
	case "2|||":
		return a[0].(g.IListCat).ListCat(a[1]), nil

	// set operations
	case "2++":
		return a[0].(g.IUnion).Union(a[1]), nil
	case "2--":
		return a[0].(g.ISetDiff).SetDiff(a[1]), nil
	case "2**":
		return a[0].(g.IIntersect).Intersect(a[1]), nil

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

//  nonDeref gives the number of args that are NOT dereferenced for an operator
//  (with the default value of zero being correct for operators not listed here)
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
func deltaSlice(lval g.Value, a []g.Value, sign int) (g.Value, *g.Closure) {
	x := g.Deref(a[0]).(g.ISlice)
	i := int(a[1].(g.Numerable).ToNumber().Val())
	j := i + sign*int(a[2].(g.Numerable).ToNumber().Val())
	if (i > 0 && j <= 0) || (i <= 0 && j > 0) { // if wraparound
		return nil, nil // fail
	}
	return x.Slice(lval, g.NewNumber(float64(i)), g.NewNumber(float64(j))), nil
}

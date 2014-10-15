//  fnumber.go -- functions operating on numbers
//
//  In general, these do no error checking.

package goaldi

import (
	"math"
	"math/rand"
)

func init() {

	LibGoFunc("abs", math.Abs)
	LibGoFunc("ceil", math.Ceil)
	LibGoFunc("floor", math.Floor)

	LibProcedure("number", Number)
	LibProcedure("min", Min)
	LibProcedure("max", Max)

	LibGoFunc("log", Log)
	LibGoFunc("sqrt", math.Sqrt)

	LibGoFunc("sin", math.Sin)
	LibGoFunc("cos", math.Cos)
	LibGoFunc("tan", math.Tan)
	LibGoFunc("asin", math.Asin)
	LibGoFunc("acos", math.Acos)
	LibGoFunc("atan", Atan)

	LibGoFunc("seed", rand.Seed)
}

//------------------------------------  functions with Go interface

//  Atan(r1, r2) -- arctangent(r1/r2), default r2 = 1.0
func Atan(r1 float64, x2 interface{}) float64 {
	switch r2 := x2.(type) {
	case nil:
		return math.Atan(r1)
	case float64:
		return math.Atan2(r1, r2)
	case string:
		return math.Atan2(r1, MustParseNum(r2))
	default:
		return Atan(r1, x2.(Numerable).ToNumber().Val())
	}
}

//  Log(r1, r2) -- logarithm of r1 to base r2, default r2 = e
func Log(r1 float64, x2 interface{}) float64 {
	switch r2 := x2.(type) {
	case nil:
		return math.Log(r1)
	case float64:
		switch r2 {
		case 2.0:
			return math.Log2(r1)
		case 10.0:
			return math.Log10(r1)
		default:
			return math.Log(r1) / math.Log(r2)
		}
	case string:
		return Log(r1, MustParseNum(r2))
	default:
		return Log(r1, x2.(Numerable).ToNumber().Val())
	}
}

//------------------------------------  procedures with Goaldi interface

//  Number(x) -- return argument converted to number
func Number(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("number", a)
	return Return(ProcArg(a, 0, NilValue).(Numerable).ToNumber())
}

//  Min(n1, ...) -- return numeric minimum
func Min(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("min", a)
	v := ProcArg(a, 0, NilValue).(Numerable).ToNumber().Val()
	for i := 1; i < len(a); i++ {
		vi := a[i].(Numerable).ToNumber().Val()
		if vi < v {
			v = vi
		}
	}
	return Return(v)
}

//  Max(n1, ...) -- return numeric maximum
func Max(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("max", a)
	v := ProcArg(a, 0, NilValue).(Numerable).ToNumber().Val()
	for i := 1; i < len(a); i++ {
		vi := a[i].(Numerable).ToNumber().Val()
		if vi > v {
			v = vi
		}
	}
	return Return(v)
}

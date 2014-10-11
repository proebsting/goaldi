//  fmath.go -- numeric functions
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
	LibProcedure("min", min)
	LibProcedure("max", max)

	LibGoFunc("log", log)
	LibGoFunc("sqrt", math.Sqrt)

	LibGoFunc("sin", math.Sin)
	LibGoFunc("cos", math.Cos)
	LibGoFunc("tan", math.Tan)
	LibGoFunc("asin", math.Asin)
	LibGoFunc("acos", math.Acos)
	LibGoFunc("atan", atan)

	LibGoFunc("seed", rand.Seed)
}

//------------------------------------  functions with Go interface

//  atan(r1, r2) -- arctangent(r1/r2), default r2 = 1.0
func atan(r1 float64, x2 interface{}) float64 {
	switch r2 := x2.(type) {
	case nil:
		return math.Atan(r1)
	case float64:
		return math.Atan2(r1, r2)
	case string:
		return math.Atan2(r1, MustParseNum(r2))
	default:
		return atan(r1, x2.(Numerable).ToNumber().Val())
	}
}

//  log(r1, r2) -- logarithm of r1 to base r2, default r2 = e
func log(r1 float64, x2 interface{}) float64 {
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
		return log(r1, MustParseNum(r2))
	default:
		return log(r1, x2.(Numerable).ToNumber().Val())
	}
}

//------------------------------------  procedures with Goaldi interface

//  min(n1, ...) -- return numeric minimum
func min(env *Env, a ...Value) (Value, *Closure) {
	n := len(a)
	if n == 0 {
		panic(&RunErr{"min(): no arguments", nil})
	}
	v := a[0].(Numerable).ToNumber().Val()
	for i := 1; i < n; i++ {
		vi := a[i].(Numerable).ToNumber().Val()
		if vi < v {
			v = vi
		}
	}
	return Return(v)
}

//  max(n1, ...) -- return numeric maximum
func max(env *Env, a ...Value) (Value, *Closure) {
	n := len(a)
	if n == 0 {
		panic(&RunErr{"max(): no arguments", nil})
	}
	v := a[0].(Numerable).ToNumber().Val()
	for i := 1; i < n; i++ {
		vi := a[i].(Numerable).ToNumber().Val()
		if vi > v {
			v = vi
		}
	}
	return Return(v)
}

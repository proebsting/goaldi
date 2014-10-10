//  fmath.go -- numeric functions
//
//  In general, these do no error checking.

package goaldi

import (
	"math"
	"math/rand"
)

func init() {

	LibProc("abs", math.Abs)
	LibProc("ceil", math.Ceil)
	LibProc("floor", math.Floor)
	LibProc("min", math.Min) // not like Icon: only 2 args
	LibProc("max", math.Max) // not like Icon: only 2 args

	LibProc("log", log)
	LibProc("sqrt", math.Sqrt)

	LibProc("sin", math.Sin)
	LibProc("cos", math.Cos)
	LibProc("tan", math.Tan)
	LibProc("asin", math.Asin)
	LibProc("acos", math.Acos)
	LibProc("atan", atan)

	LibProc("seed", rand.Seed)
}

//  atan(r1, r2) -- arctangent(r1/r2), default r2 = 1.0
func atan(r1 float64, x2 interface{}) float64 {
	switch r2 := x2.(type) {
	case nil:
		return math.Atan(r1)
	case float64:
		return math.Atan2(r1, r2)
	case string:
		return math.Atan2(r1, NewString(r2).ToNumber().Val())
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
		return log(r1, NewString(r2).ToNumber().Val())
	default:
		return log(r1, x2.(Numerable).ToNumber().Val())
	}
}

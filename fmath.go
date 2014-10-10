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
	LibGoFunc("min", math.Min) // not like Icon: only 2 args
	LibGoFunc("max", math.Max) // not like Icon: only 2 args

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

//  fnumber.go -- functions operating on numbers
//
//  In general, these do no error checking.

package goaldi

import (
	"math"
	"math/rand"
)

//  Declare methods
var NumberMethods = map[string]interface{}{
	"type":   (*VNumber).Type,
	"copy":   (*VNumber).Copy,
	"string": (*VNumber).String,
	"image":  (*VNumber).GoString,
}

//  VNumber.Field implements methods
func (v *VNumber) Field(f string) Value {
	return GetMethod(NumberMethods, v, f)
}

func init() {
	// Goaldi procedures
	LibProcedure("number", Number)
	LibProcedure("min", Min)
	LibProcedure("max", Max)
	LibProcedure("log", Log)
	LibProcedure("atan", Atan)
	LibProcedure("gcd", GCD)
	// Go library functions
	LibGoFunc("abs", math.Abs)
	LibGoFunc("ceil", math.Ceil)
	LibGoFunc("floor", math.Floor)
	LibGoFunc("sqrt", math.Sqrt)
	LibGoFunc("seed", rand.Seed)
	LibGoFunc("sin", math.Sin)
	LibGoFunc("cos", math.Cos)
	LibGoFunc("tan", math.Tan)
	LibGoFunc("asin", math.Asin)
	LibGoFunc("acos", math.Acos)
}

//  Number(x) -- return argument converted to number, or fail
func Number(env *Env, a ...Value) (Value, *Closure) {
	// nonstandard entry; on panic, returns default nil values
	defer func() { recover() }()
	v := ProcArg(a, 0, NilValue)
	if n, ok := v.(Numerable); ok {
		return Return(n.ToNumber())
	} else {
		return Return(Import(v).(Numerable).ToNumber())
	}
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
	return Return(NewNumber(v))
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
	return Return(NewNumber(v))
}

//  Log(r1, r2) -- logarithm of r1 to base r2, default r2 = e
func Log(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("log", a)
	r1 := ProcArg(a, 0, NilValue).(Numerable).ToNumber().Val()
	r2 := ProcArg(a, 1, E).(Numerable).ToNumber().Val()
	if r2 == math.E {
		return Return(NewNumber(math.Log(r1)))
	} else {
		return Return(NewNumber(math.Log(r1) / math.Log(r2)))
	}
}

//  Atan(r1, r2) -- arctangent(r1/r2), default r2 = 1.0
func Atan(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("atan", a)
	r1 := ProcArg(a, 0, NilValue).(Numerable).ToNumber().Val()
	r2 := ProcArg(a, 1, ONE).(Numerable).ToNumber().Val()
	if r2 == 1.0 {
		return Return(NewNumber(math.Atan(r1)))
	} else {
		return Return(NewNumber(math.Atan2(r1, r2)))
	}
}

//  GCD(i, ...) -- greatest common divisor
//  Returns the GCD of one or more values, which are truncated to int.
//  Negative values are allowed.  Returns zero if all values are zero.
func GCD(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("gcd", args)
	a := int(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	if a < 0 {
		a = -a
	}
	for i := 1; i < len(args); i++ {
		b := int(args[i].(Numerable).ToNumber().Val())
		if b < 0 {
			b = -b
		}
		for b > 0 {
			a, b = b, a%b
		}
	}
	return Return(NewNumber(float64(a)))
}

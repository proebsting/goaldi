//  fnumlist.go -- functions operating on lists of numbers

package runtime

import (
	"math"
)

func init() {
	// Goaldi procedures
	DefLib(Min, "min", "n[]", "find minimum value")
	DefLib(Max, "max", "n[]", "find maximum value")
	DefLib(GCD, "gcd", "i[]", "find greatest common divisor")
	DefLib(Amean, "amean", "n[]", "compute arithmetic mean")
	DefLib(Gmean, "gmean", "n[]", "compute geometric mean")
	DefLib(Hmean, "hmean", "n[]", "compute harmonic mean")
	DefLib(Qmean, "qmean", "n[]", "compute quadratic mean")
}

//  min(n, ...) returns the smallest of its arguments.
func Min(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("min", args)
	v := FloatVal(ProcArg(args, 0, NilValue))
	for i := 1; i < len(args); i++ {
		vi := FloatVal(args[i])
		if vi < v {
			v = vi
		}
	}
	return Return(NewNumber(v))
}

//  max(n, ...) returns the largest of its arguments.
func Max(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("max", args)
	v := FloatVal(ProcArg(args, 0, NilValue))
	for i := 1; i < len(args); i++ {
		vi := FloatVal(args[i])
		if vi > v {
			v = vi
		}
	}
	return Return(NewNumber(v))
}

//  gcd(i,...) truncates its arguments to integer and
//  returns their greatest common divisor.
//  Negative values are allowed.
//  gcd() returns zero if all values are zero.
func GCD(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("gcd", args)
	a := int(FloatVal(ProcArg(args, 0, NilValue)))
	if a < 0 {
		a = -a
	}
	for i := 1; i < len(args); i++ {
		b := int(FloatVal(args[i]))
		if b < 0 {
			b = -b
		}
		for b > 0 {
			a, b = b, a%b
		}
	}
	return Return(NewNumber(float64(a)))
}

//  amean(n,...) returns the arithmetic mean, or simple average,
//  of its arguments.
func Amean(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("amean", args)
	t := FloatVal(ProcArg(args, 0, NilValue))
	for i := 1; i < len(args); i++ {
		t += FloatVal(ProcArg(args, i, NilValue))
	}
	return Return(NewNumber(float64(t) / float64(len(args))))
}

//  gmean(n,...) returns the geometric mean of its arguments,
//  which must all be strictly positive.
func Gmean(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("gmean", args)
	p := FloatVal(ProcArg(args, 0, NilValue))
	if p <= 0 {
		panic(NewExn("Nonpositive argument", p))
	}
	for i := 1; i < len(args); i++ {
		v := FloatVal(ProcArg(args, i, NilValue))
		if v <= 0 {
			panic(NewExn("Nonpositive argument", v))
		}
		p *= v
	}
	return Return(NewNumber(math.Exp(math.Log(p) / float64(len(args)))))
}

//  hmean(n,...) returns the harmonic mean of its arguments,
//  which must all be strictly positive.
func Hmean(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("hmean", args)
	v := FloatVal(ProcArg(args, 0, NilValue))
	if v <= 0 {
		panic(NewExn("Nonpositive argument", v))
	}
	t := 1 / v
	for i := 1; i < len(args); i++ {
		v = FloatVal(ProcArg(args, i, NilValue))
		if v <= 0 {
			panic(NewExn("Nonpositive argument", v))
		}
		t += 1 / v
	}
	return Return(NewNumber(float64(len(args)) / t))
}

//  qmean(n,...) returns the quadratic mean, or root mean square,
//  of its arguments.
func Qmean(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("qmean", args)
	v := FloatVal(ProcArg(args, 0, NilValue))
	t := v * v
	for i := 1; i < len(args); i++ {
		v = FloatVal(ProcArg(args, i, NilValue))
		t += v * v
	}
	return Return(NewNumber(math.Sqrt(float64(t) / float64(len(args)))))
}

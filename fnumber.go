//  fnumber.go -- functions operating on numbers
//
//  In general, these do no error checking.

package goaldi

import (
	"math"
	"math/rand"
)

//  Declare methods
var NumberMethods = MethodTable([]*VProcedure{
	DefMeth((*VNumber).Type, "type", "", "return number type"),
	DefMeth((*VNumber).Copy, "copy", "", "return number value"),
	DefMeth((*VNumber).String, "string", "", "return short string"),
	DefMeth((*VNumber).GoString, "image", "", "return full string"),
})

//  VNumber.Field implements methods
func (v *VNumber) Field(f string) Value {
	return GetMethod(NumberMethods, v, f)
}

func init() {
	// Goaldi procedures
	DefLib(Min, "min", "n[]", "find minimum value")
	DefLib(Max, "max", "n[]", "find maximum value")
	DefLib(Seq, "seq", "n,incr", "produce n to infinity")
	DefLib(Log, "log", "n,b", "compute logarithm to base b")
	DefLib(Atan, "atan", "x,y", "compute arctangent of x / y")
	DefLib(GCD, "gcd", "i[]", "find greatest common divisor")
	DefLib(RandGen, "randgen", "seed", "create independent random sequence")
	DefLib(RtoD, "rtod", "r", "convert radians to degrees")
	DefLib(DtoR, "dtor", "d", "convert degrees to radians")
	DefLib(IAnd, "iand", "i,j", "compute bitwise AND")
	DefLib(IOr, "ior", "i,j", "compute bitwise OR")
	DefLib(IXor, "ixor", "i,j", "compute bitwise exclusive OR")
	DefLib(IClear, "iclear", "i,j", "compute bitwise clear of i by j")
	DefLib(ICom, "icom", "i", "compute bitwise complement")
	DefLib(IShift, "ishift", "i,j", "compute bitwise shift of i by j")
	// Go library functions
	GoLib(math.Abs, "abs", "x", "compute absolute value")
	GoLib(math.Ceil, "ceil", "x", "round up to integer")
	GoLib(math.Floor, "floor", "x", "round down to integer")
	GoLib(math.Trunc, "integer", "x", "truncate to integer")
	GoLib(math.Cbrt, "cbrt", "x", "compute cube root")
	GoLib(math.Sqrt, "sqrt", "x", "compute square root")
	GoLib(math.Hypot, "hypot", "p,q", "return sqrt of p^2 + q^2")
	GoLib(math.Exp, "exp", "x", "return %e^x")
	GoLib(rand.Seed, "seed", "n", "set random number seed")
	GoLib(math.Sin, "sin", "x", "compute sine of x")
	GoLib(math.Cos, "cos", "x", "compute cosine of x")
	GoLib(math.Tan, "tan", "x", "compute tangent of x")
	GoLib(math.Asin, "asin", "x", "compute arcsine of x")
	GoLib(math.Acos, "acos", "x", "compute arccosine of x")
}

//  Number(x) -- return argument converted to number, or fail
func Number(env *Env, args ...Value) (Value, *Closure) {
	// nonstandard entry; on panic, returns default nil values to fail
	defer func() { recover() }()
	v := ProcArg(args, 0, NilValue)
	if n, ok := v.(Numerable); ok {
		return Return(n.ToNumber())
	} else {
		return Return(Import(v).(Numerable).ToNumber())
	}
}

//  Seq(n1,n2) -- produce n1 to infinity by n2
func Seq(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("seq", args)
	n1 := ProcArg(args, 0, ONE).(Numerable).ToNumber()
	n2 := ProcArg(args, 1, ONE).(Numerable).ToNumber()
	return ToBy(n1, INF, n2)
}

//  Min(n1, ...) -- return numeric minimum
func Min(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("min", args)
	v := ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val()
	for i := 1; i < len(args); i++ {
		vi := args[i].(Numerable).ToNumber().Val()
		if vi < v {
			v = vi
		}
	}
	return Return(NewNumber(v))
}

//  Max(n1, ...) -- return numeric maximum
func Max(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("max", args)
	v := ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val()
	for i := 1; i < len(args); i++ {
		vi := args[i].(Numerable).ToNumber().Val()
		if vi > v {
			v = vi
		}
	}
	return Return(NewNumber(v))
}

//  Log(r1, r2) -- logarithm of r1 to base r2, default r2 = e
func Log(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("log", args)
	r1 := ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val()
	r2 := ProcArg(args, 1, E).(Numerable).ToNumber().Val()
	if r2 == math.E {
		return Return(NewNumber(math.Log(r1)))
	} else {
		return Return(NewNumber(math.Log(r1) / math.Log(r2)))
	}
}

//  Atan(r1, r2) -- arctangent(r1/r2), default r2 = 1.0
func Atan(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("atan", args)
	r1 := ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val()
	r2 := ProcArg(args, 1, ONE).(Numerable).ToNumber().Val()
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

//  RandGen(i) -- new random generator seeded by i
//  Returns a new Go math.rand/Rand object whose methods may be called.
func RandGen(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("randgen", args)
	i := int64(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val())
	return Return(rand.New(rand.NewSource(i)))
}

//  DtoR(r1) -- convert degrees to radians
func DtoR(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("dtor", args)
	r := ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val()
	return Return(NewNumber(r * math.Pi / 180.0))
}

//  RtoD(r1) -- convert radians to degrees
func RtoD(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("rtod", args)
	r := ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val()
	return Return(NewNumber(r * 180.0 / math.Pi))
}

//  IAnd(i1, i2) -- bitwise AND of i1 and i2
func IAnd(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("iand", args)
	i1 := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	i2 := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(i1 & i2)))
}

//  IOr(i1, i2) -- bitwise OR of i1 and i2
func IOr(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("ior", args)
	i1 := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	i2 := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(i1 | i2)))
}

//  IXor(i1, i2) -- bitwise XOR of i1 and i2
func IXor(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("ixor", args)
	i1 := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	i2 := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(i1 ^ i2)))
}

//  IClear(i1, i2) -- bitwise clear of i1 by i2
func IClear(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("iclear", args)
	i1 := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	i2 := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(i1 &^ i2)))
}

//  IShift(i1, i2) -- bitwise shift of i1 by i2
func IShift(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("ishift", args)
	i1 := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	i2 := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	if i2 > 0 {
		return Return(NewNumber(float64(i1 << uint(i2))))
	} else {
		return Return(NewNumber(float64(i1 >> uint(-i2))))
	}
}

//  ICom(i) -- bitwise complement of i
func ICom(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("icom", args)
	i := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(^i)))
}

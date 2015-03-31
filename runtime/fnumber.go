//  fnumber.go -- functions operating on numbers
//
//  In general, these do no error checking.

package runtime

import (
	"encoding/binary"
	"math"
	"math/rand"
	"os"
)

func init() {
	// Goaldi procedures
	DefLib(Seq, "seq", "n,incr", "produce n to infinity")
	DefLib(Log, "log", "n,b", "compute logarithm to base b")
	DefLib(Atan, "atan", "y,x", "compute arctangent of y / x")
	DefLib(Randomize, "randomize", "", "irreproducibly seed random generation")
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
	GoLib(math.Abs, "abs", "n", "compute absolute value")
	GoLib(math.Ceil, "ceil", "n", "round up to integer")
	GoLib(math.Floor, "floor", "n", "round down to integer")
	GoLib(math.Trunc, "integer", "n", "truncate to integer")
	GoLib(math.Cbrt, "cbrt", "n", "compute cube root")
	GoLib(math.Sqrt, "sqrt", "n", "compute square root")
	GoLib(math.Hypot, "hypot", "x,y", "return sqrt of x^2 + y^2")
	GoLib(math.Exp, "exp", "n", "return e ^ x")
	GoLib(rand.Seed, "seed", "n", "set random number seed")
	GoLib(math.Sin, "sin", "n", "compute sine")
	GoLib(math.Cos, "cos", "n", "compute cosine")
	GoLib(math.Tan, "tan", "n", "compute tangent")
	GoLib(math.Asin, "asin", "n", "compute arcsine")
	GoLib(math.Acos, "acos", "n", "compute arccosine")
	GoLib(math.Sinh, "sinh", "n", "compute hyperbolic sine")
	GoLib(math.Cosh, "cosh", "n", "compute hyperbolic cosine")
	GoLib(math.Tanh, "tanh", "n", "compute hyperbolic tangent")
	GoLib(math.Asinh, "asinh", "n", "compute hyperbolic arcsine")
	GoLib(math.Acosh, "acosh", "n", "compute hyperbolic arccosine")
	GoLib(math.Atanh, "atanh", "n", "compute hyperbolic arccosine")
}

//  number(x) returns its argument converted to number,
//  or fails if it cannot be converted due to its form or dataype.
//  For string (or stringable) arguments,
//  number() trims leading and trailing spaces
//  and then accepts standard Go decimal forms (fixed and floating)
//  or Goaldi radix forms (101010b, 52o, 2Ax, 23r1J).
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

//  seq(n,incr) generates an endless sequence of values beginning at n
//  with increments of incr.
func Seq(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("seq", args)
	n1 := ProcArg(args, 0, ONE).(Numerable).ToNumber()
	n2 := ProcArg(args, 1, ONE).(Numerable).ToNumber()
	return ToBy(n1, INF, n2)
}

//  log(n, b) returns the logarithm of n to base b.
//  The default value of b is %e (2.7183...),
//  so log(n) returns the natural logarithm of n.
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

//  atan(y, x) returns the arctangent, in radians, of (y/x).
//  The default value of x is 1, so atan(y) returns the arctangent of y.
//  For the handling of special cases see http://golang.org/pkg/math/#Atan2.
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

//  randomize() seeds the random number generator
//  with an irreproducible value obtained from /dev/urandom.
func Randomize(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("randomize", args)
	var seed int64
	f, err := os.Open("/dev/urandom")
	if err != nil {
		panic(err)
	}
	err = binary.Read(f, binary.LittleEndian, &seed)
	if err != nil {
		panic(err)
	}
	f.Close()
	seed = seed & 0x0000FFFFFFFFFFFF // 48 bits
	rand.Seed(seed)
	return Return(NewNumber(float64(seed)))
}

//  randgen(i) returns a new random generator seeded by i.
//  The returned external value is a Go math.rand/Rand object
//  whose methods may be called from Goaldi to produce random values.
func RandGen(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("randgen", args)
	i := int64(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val())
	return Return(rand.New(rand.NewSource(i)))
}

//  dtor(d) returns the radian equivalent of the angle d given in degrees.
func DtoR(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("dtor", args)
	r := ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val()
	return Return(NewNumber(r * math.Pi / 180.0))
}

//  rtod(r) returns the degree equivalent of the angle r given in radians.
func RtoD(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("rtod", args)
	r := ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val()
	return Return(NewNumber(r * 180.0 / math.Pi))
}

//  iand(i, j) returns the bitwise AND of the values i and j truncated to integer.
func IAnd(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("iand", args)
	i := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	j := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(i & j)))
}

//  ior(i, j) returns the bitwise OR of the values i and j truncated to integer.
func IOr(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("ior", args)
	i := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	j := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(i | j)))
}

//  ixor(i, j) returns the bitwise exclusive OR
//  of the values i and j truncated to integer.
func IXor(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("ixor", args)
	i := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	j := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(i ^ j)))
}

//  iclear(i, j) returns the value of i cleared of those bits set in j,
//  after truncating both arguments to integer.
func IClear(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("iclear", args)
	i := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	j := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(i &^ j)))
}

//  ishift(i, j) shifts i by j bits and returns the result.
//  If j > 0, the shift is to the left with zero fill.
//  If j < 0, the shift is to the right with sign extension.
//  The arguments are both truncated to integer before operating.
func IShift(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("ishift", args)
	i := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	j := int64(ProcArg(args, 1, NilValue).(Numerable).ToNumber().Val())
	if j > 0 {
		return Return(NewNumber(float64(i << uint(j))))
	} else {
		return Return(NewNumber(float64(i >> uint(-j))))
	}
}

//  icom(i) truncates i to integer and returns its bitwise complement.
func ICom(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("icom", args)
	i := int64(ProcArg(args, 0, NilValue).(Numerable).ToNumber().Val())
	return Return(NewNumber(float64(^i)))
}

//  onumber.go -- operations on numbers

package runtime

import (
	"math"
	"math/rand"
)

// extract float64 value from arbitrary Value, or panic
func fval(v Value) float64 {
	if n, ok := v.(Numerable); ok {
		return float64(*(n.ToNumber()))
	} else {
		panic(NewExn("Not a number", v))
	}
}

//------------------------------------  Choose:  ?e

func (v1 *VNumber) Choose(unused Value) Value {
	n := v1.Val()
	if n < 0 {
		panic(NewExn("?n < 0", v1))
	} else if n == 0 {
		return NewNumber(rand.Float64())
	} else /* n > 0 */ {
		return NewNumber(float64(int(n * rand.Float64())))
	}
}

//------------------------------------  Dispense:  !e

func (v1 *VNumber) Dispense(unused Value) (Value, *Closure) {
	return ToBy(ONE, v1, ONE)
}

//------------------------------------  Numerate:  +e

type INumerate interface {
	Numerate() Value // +n
}

func (v1 *VString) Numerate() Value {
	return v1.ToNumber()
}

func (v1 *VNumber) Numerate() Value {
	return v1
}

//------------------------------------  Negate:  -e

type INegate interface {
	Negate() Value // -n
}

func (v1 *VString) Negate() Value {
	return v1.ToNumber().Negate()
}

func (v1 *VNumber) Negate() Value {
	return NewNumber(-float64(*v1))
}

//------------------------------------  Add:  e1 + e2

type IAdd interface {
	Add(Value) Value // n + n
}

func (v1 *VString) Add(v2 Value) Value {
	return v1.ToNumber().Add(v2)
}

func (v1 *VNumber) Add(v2 Value) Value {
	return NewNumber(float64(*v1) + fval(v2))
}

//------------------------------------  Sub:  e1 - e2

type ISub interface {
	Sub(Value) Value // n - n
}

func (v1 *VString) Sub(v2 Value) Value {
	return v1.ToNumber().Sub(v2)
}

func (v1 *VNumber) Sub(v2 Value) Value {
	return NewNumber(float64(*v1) - fval(v2))
}

//------------------------------------  Mul:  e1 * e2

type IMul interface {
	Mul(Value) Value // n * n
}

func (v1 *VString) Mul(v2 Value) Value {
	return v1.ToNumber().Mul(v2)
}

func (v1 *VNumber) Mul(v2 Value) Value {
	return NewNumber(float64(*v1) * fval(v2))
}

//------------------------------------  Div:  e1 / e2

type IDiv interface {
	Div(Value) Value // n / n
}

func (v1 *VString) Div(v2 Value) Value {
	return v1.ToNumber().Div(v2)
}

func (v1 *VNumber) Div(v2 Value) Value {
	return NewNumber(float64(*v1) / fval(v2))
}

//------------------------------------  Divt:  e1 // e2  (divide and truncate)

type IDivt interface {
	Divt(Value) Value // n // n
}

func (v1 *VString) Divt(v2 Value) Value {
	return v1.ToNumber().Divt(v2)
}

func (v1 *VNumber) Divt(v2 Value) Value {
	return NewNumber(float64(int64(float64(*v1) / fval(v2))))
}

//------------------------------------  Mod:  e1 % e2  (remainder)

type IMod interface {
	Mod(Value) Value // n % n
}

func (v1 *VString) Mod(v2 Value) Value {
	return v1.ToNumber().Mod(v2)
}

func (v1 *VNumber) Mod(v2 Value) Value {
	return NewNumber(math.Mod(float64(*v1), fval(v2)))
}

//------------------------------------  Power:  e1 ^ e2

type IPower interface {
	Power(Value) Value // n ^ n
}

func (v1 *VString) Power(v2 Value) Value {
	return v1.ToNumber().Power(v2)
}

func (v1 *VNumber) Power(v2 Value) Value {
	return NewNumber(math.Pow(float64(*v1), fval(v2)))
}

//------------------------------------  NumLT:  e1 < e2

type INumLT interface {
	NumLT(Value) (Value, *Closure) // n < n
}

func (v1 *VString) NumLT(v2 Value) (Value, *Closure) {
	return v1.ToNumber().NumLT(v2)
}

func (v1 *VNumber) NumLT(v2 Value) (Value, *Closure) {
	if float64(*v1) < fval(v2) {
		return Return(v2)
	} else {
		return Fail()
	}
}

//------------------------------------  NumLE:  e1 <= e2

type INumLE interface {
	NumLE(Value) (Value, *Closure) // n <= n
}

func (v1 *VString) NumLE(v2 Value) (Value, *Closure) {
	return v1.ToNumber().NumLE(v2)
}

func (v1 *VNumber) NumLE(v2 Value) (Value, *Closure) {
	if float64(*v1) <= fval(v2) {
		return Return(v2)
	} else {
		return Fail()
	}
}

//------------------------------------  NumEQ:  e1 = e2

type INumEQ interface {
	NumEQ(Value) (Value, *Closure) // n = n
}

func (v1 *VString) NumEQ(v2 Value) (Value, *Closure) {
	return v1.ToNumber().NumEQ(v2)
}

func (v1 *VNumber) NumEQ(v2 Value) (Value, *Closure) {
	if float64(*v1) == fval(v2) {
		return Return(v2)
	} else {
		return Fail()
	}
}

//------------------------------------  NumNE:  e1 ~= e2

type INumNE interface {
	NumNE(Value) (Value, *Closure) // n ~= n
}

func (v1 *VString) NumNE(v2 Value) (Value, *Closure) {
	return v1.ToNumber().NumNE(v2)
}

func (v1 *VNumber) NumNE(v2 Value) (Value, *Closure) {
	if float64(*v1) != fval(v2) {
		return Return(v2)
	} else {
		return Fail()
	}
}

//------------------------------------  NumGE:  e1 >= e2

type INumGE interface {
	NumGE(Value) (Value, *Closure) // n >= n
}

func (v1 *VString) NumGE(v2 Value) (Value, *Closure) {
	return v1.ToNumber().NumGE(v2)
}

func (v1 *VNumber) NumGE(v2 Value) (Value, *Closure) {
	if float64(*v1) >= fval(v2) {
		return Return(v2)
	} else {
		return Fail()
	}
}

//------------------------------------  NumGT:  e1 > e2

type INumGT interface {
	NumGT(Value) (Value, *Closure) // n > n
}

func (v1 *VString) NumGT(v2 Value) (Value, *Closure) {
	return v1.ToNumber().NumGT(v2)
}

func (v1 *VNumber) NumGT(v2 Value) (Value, *Closure) {
	if float64(*v1) > fval(v2) {
		return Return(v2)
	} else {
		return Fail()
	}
}

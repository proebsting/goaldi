//  omath.go -- math operations

package goaldi

import "math"

//  extract float64 value from arbitrary Value, or panic
func fval(v Value) float64 {
	if n, ok := v.(Numerable); ok {
		return float64(*(n.ToNumber()))
	} else {
		panic(&RunErr{"Not a number", v})
	}
}

//------------------------------------  Numerate:  +e

type INumerate interface {
	Numerate() Value
}

func (v1 *VString) Numerate() Value {
	return v1.ToNumber()
}

func (v1 *VNumber) Numerate() Value {
	return v1
}

//------------------------------------  Negate:  -e

type INegate interface {
	Negate() Value
}

func (v1 *VString) Negate() Value {
	return v1.ToNumber().Negate()
}

func (v1 *VNumber) Negate() Value {
	return NewNumber(-float64(*v1))
}

//------------------------------------  Add:  e1 + e2

type IAdd interface {
	Add(Value) Value
}

func (v1 *VString) Add(v2 Value) Value {
	return v1.ToNumber().Add(v2)
}

func (v1 *VNumber) Add(v2 Value) Value {
	return NewNumber(float64(*v1) + fval(v2))
}

//------------------------------------  Sub:  e1 - e2

type ISub interface {
	Sub(Value) Value
}

func (v1 *VString) Sub(v2 Value) Value {
	return v1.ToNumber().Sub(v2)
}

func (v1 *VNumber) Sub(v2 Value) Value {
	return NewNumber(float64(*v1) - fval(v2))
}

//------------------------------------  Mul:  e1 * e2

type IMul interface {
	Mul(Value) Value
}

func (v1 *VString) Mul(v2 Value) Value {
	return v1.ToNumber().Mul(v2)
}

func (v1 *VNumber) Mul(v2 Value) Value {
	return NewNumber(float64(*v1) * fval(v2))
}

//------------------------------------  Div:  e1 / e2

type IDiv interface {
	Div(Value) Value
}

func (v1 *VString) Div(v2 Value) Value {
	return v1.ToNumber().Div(v2)
}

func (v1 *VNumber) Div(v2 Value) Value {
	return NewNumber(float64(*v1) / fval(v2))
}

//------------------------------------  Divt:  e1 // e2  (divide and truncate)

type IDivt interface {
	Divt(Value) Value
}

func (v1 *VString) Divt(v2 Value) Value {
	return v1.ToNumber().Divt(v2)
}

func (v1 *VNumber) Divt(v2 Value) Value {
	return NewNumber(float64(int64(float64(*v1) / fval(v2))))
}

//------------------------------------  Mod:  e1 % e2  (remainder)

type IMod interface {
	Mod(Value) Value
}

func (v1 *VString) Mod(v2 Value) Value {
	return v1.ToNumber().Mod(v2)
}

func (v1 *VNumber) Mod(v2 Value) Value {
	return NewNumber(math.Mod(float64(*v1), fval(v2)))
}

//------------------------------------  Power:  e1 ^ e2

type IPower interface {
	Power(Value) Value
}

func (v1 *VString) Power(v2 Value) Value {
	return v1.ToNumber().Power(v2)
}

func (v1 *VNumber) Power(v2 Value) Value {
	return NewNumber(math.Pow(float64(*v1), fval(v2)))
}

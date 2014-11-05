//  calling.go -- procedure call / suspension / resumption interface
//
//  In the Go implementation, a Goaldi procedure or operation result
//  is the pair (Value, *Closure) with this meaning:
//
//  Value  *Closure   Interpretation
//  -----  -------   ---------------------------------
//   nil     nil     procedure failed
//  value    nil     procedure returned a value
//  value   resume   procedure suspended and can be resumed

package goaldi

import (
	"reflect"
)

type Closure struct {
	Go Resumer // start or resumption function, depending on context
}

//  Procedure function prototype
type Procedure func(env *Env, args ...Value) (Value, *Closure)

//  Procedure resumption prototype
type Resumer func() (Value, *Closure)

//  ProcArg(a,i,d) -- return procedure argument a[i], defaulting to d
func ProcArg(a []Value, i int, d Value) Value {
	if i < len(a) && a[i] != NilValue {
		return a[i]
	} else {
		return d
	}
}

//  Resume() executes the entry point in a Closure to produce the next result.
//  If the pointer is nil, failure is produced.
func (c *Closure) Resume() (Value, *Closure) {
	if c == nil {
		return Fail()
	}
	return c.Go()
}

//  Fail returns a failure indicator
func Fail() (Value, *Closure) {
	return nil, nil
}

//  Return returns a simple value as a duo
func Return(v Value) (Value, *Closure) {
	return v, nil
}

//  An MVFunc is like a Go "method value", a function bound to an object,
//  for example the "m.delete" part of the expression "m.delete(x)"
type MVFunc struct {
	v Value
	f interface{} // func(Value, ...Value)(Value, *Closure)
}

//  GetMethod(m,v,s) looks up method v.s in table m, panicking on failure.
func GetMethod(m map[string]interface{}, v Value, s string) *MVFunc {
	method := m[s]
	if method == nil {
		panic(&RunErr{"unrecognized method: " + s, v})
	}
	return &MVFunc{v, method}
}

//  MVFunc.Call(args) invokes the underlying method function.
func (mvf *MVFunc) Call(env *Env, args ...Value) (Value, *Closure) {
	arglist := make([]reflect.Value, 1+len(args))
	arglist[0] = reflect.ValueOf(mvf.v)
	for i, v := range args {
		arglist[i+1] = reflect.ValueOf(v)
	}
	method := reflect.ValueOf(mvf.f)
	result := method.Call(arglist)
	switch len(result) {
	case 0:
		return nil, nil
	case 1:
		return Value(result[0].Interface()), nil
	default:
		return Value(result[0].Interface()), (result[1].Interface().(*Closure))
	}
}

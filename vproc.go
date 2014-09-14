//  vproc.go -- VProcedure, the Goaldi type "procedure"

package goaldi

import (
	"reflect"
)

//  Procedure function prototype
type Procedure func(...Value) (Value, *Closure)

//  Procedure value
type VProcedure struct {
	name  string
	entry Procedure
}

//  NewProcedure(name, func) -- construct a procedure value
func NewProcedure(name string, entry Procedure) *VProcedure {
	return &VProcedure{name, entry}
}

//  VProcedure.String -- return "procname()"
func (v *VProcedure) String() string {
	return v.name + "()"
}

//  VProcedure.Type -- return "procedure"
func (v *VProcedure) Type() Value {
	return type_procedure
}

//  ICall interface
type ICall interface {
	Call(...Value) (Value, *Closure)
}

//  VProcedure.Call(args) -- invoke a procedure
func (v *VProcedure) Call(args ...Value) (Value, *Closure) {
	return v.entry(args...)
}

var type_procedure = NewString("procedure")

//  GoProcedure(name, func) -- construct a procedure from a Go function
func GoProcedure(name string, f interface{}) *VProcedure {
	ftype := reflect.TypeOf(f)
	fval := reflect.ValueOf(f)
	if fval.Kind() != reflect.Func {
		panic(&RunErr{"Not a func", f})
	}
	nargs := ftype.NumIn()
	nrtn := ftype.NumOut()
	if nrtn > 1 {
		panic(&RunErr{"Multiple returns", f}) //#%#% not yet
	}
	if ftype.IsVariadic() {
		panic(&RunErr{"Variadic function", f}) //#%#% not yet
	}

	passer := make([]func(Value) reflect.Value, nargs)
	for i := 0; i < nargs; i++ {
		passer[i] = passfunc(ftype.In(i))
	}

	pfun := func(args ...Value) (Value, *Closure) {
		in := make([]reflect.Value, nargs)
		for i := 0; i < nargs; i++ {
			if i < len(args) {
				in[i] = passer[i](args[i])
			} else {
				in[i] = passer[i](NIL)
			}
		}
		out := fval.Call(in)
		if nrtn == 1 {
			return out[0].Interface(), nil //#%#% Import(out[0]), nil
		} else {
			return NIL, nil
		}
	}
	return NewProcedure(name, pfun)
}

func passfunc(t reflect.Type) func(Value) reflect.Value {
	k := t.Kind()
	switch k {
	case reflect.Int:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(
				int(v.(Numerable).ToNumber().val()))
		}
	case reflect.Float64:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(
				float64(v.(Numerable).ToNumber().val()))
		}
	case reflect.String:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(
				string(v.(Stringable).ToString().val()))
		}
	default:
		panic(&RunErr{"Unimpl paramkind", t})
	}
}

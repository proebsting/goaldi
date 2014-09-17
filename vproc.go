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

	//  get information about the Go function
	ftype := reflect.TypeOf(f)
	fval := reflect.ValueOf(f)
	if fval.Kind() != reflect.Func {
		panic(&RunErr{"Not a func", f})
	}
	nargs := ftype.NumIn()
	nfixed := nargs
	if ftype.IsVariadic() {
		nfixed--
	}
	nrtn := ftype.NumOut()

	//  make an array of conversion functions, one per parameter
	passer := make([]func(Value) reflect.Value, nargs)
	for i := 0; i < nfixed; i++ {
		passer[i] = passfunc(ftype.In(i))
	}
	if nfixed < nargs { // if variadic
		passer[nfixed] = passfunc(ftype.In(nfixed).Elem())
	}

	//  define a func to convert arguments and call the underlying func
	pfun := func(args ...Value) (Value, *Closure) {
		//  convert fixed arguments from Goaldi values to needed Go type
		in := make([]reflect.Value, 0, len(args))
		for i := 0; i < nfixed; i++ {
			if i < len(args) {
				in = append(in, passer[i](args[i]))
			} else {
				in = append(in, passer[i](NIL))
			}
		}
		//  convert additional variadic arguments to final type
		for i := nfixed; i < len(args); i++ {
			in = append(in, passer[nfixed](args[i]))
		}
		//  call the Go function
		out := fval.Call(in)
		//  return the result   #%#% returns first result only!!!
		if nrtn >= 1 {
			return Import(out[0].Interface()), nil
		} else {
			return NIL, nil
		}
	}

	//  make this function a Goaldi procedure, and return it
	return NewProcedure(name, pfun)
}

//  passfunc returns a function that converts a Goaldi value
//  into a Go reflect.Value of the specified type
func passfunc(t reflect.Type) func(Value) reflect.Value {
	k := t.Kind()
	switch k {
	case reflect.Int:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(
				int(v.(Numerable).ToNumber().val()))
		}
	case reflect.Int64:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(
				int64(v.(Numerable).ToNumber().val()))
		}
	case reflect.Float64:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(
				float64(v.(Numerable).ToNumber().val()))
		}
	case reflect.String:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(
				v.(Stringable).ToString().ToUTF8())
		}
	case reflect.Interface: //#%#% assuming interface{}
		return func(v Value) reflect.Value {
			return reflect.ValueOf(Export(v)) // default conversion
		}
	default:
		panic(&RunErr{"Unimpl paramkind", t})
	}
}

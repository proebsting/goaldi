//  vproc.go -- VProcedure, the Goaldi type "procedure"

package goaldi

import (
	"reflect"
)

//  Procedure value
type VProcedure struct {
	Name  string
	Entry Procedure
}

//  NewProcedure(name, func) -- construct a procedure value
func NewProcedure(name string, entry Procedure) *VProcedure {
	return &VProcedure{name, entry}
}

//  VProcedure.String -- default conversion to Go string returns "P:procname"
func (v *VProcedure) String() string {
	return "P:" + v.Name
}

//  VProcedure.GoString -- convert to string for image() and printf("%#v")
func (v *VProcedure) GoString() string {
	return "procedure " + v.Name + "()"
}

//  VProcedure.Rank returns rProc
func (v *VProcedure) Rank() int {
	return rProc
}

//  VProcedure.Type -- return "procedure"
func (v *VProcedure) Type() Value {
	return type_procedure
}

var type_procedure = NewString("procedure")

//  VProcedure.Copy returns itself
func (v *VProcedure) Copy() Value {
	return v
}

//  VProcedure.Import returns itself
func (v *VProcedure) Import() Value {
	return v
}

//  VProcedure.Export returns the underlying function
//  (#%#% at least for now. should we wrap it somehow?)
func (v *VProcedure) Export() interface{} {
	return v.Entry
}

//  ICall interface
type ICall interface {
	Call(*Env, ...Value) (Value, *Closure)
}

//  VProcedure.Call(args) -- invoke a procedure
func (v *VProcedure) Call(env *Env, args ...Value) (Value, *Closure) {
	return v.Entry(env, args...)
}

//  Declare methods
var ProcedureMethods = map[string]interface{}{
	"type":  (*VProcedure).Type,
	"copy":  (*VProcedure).Copy,
	"image": Image,
}

//  VProcedure.Field implements methods
func (v *VProcedure) Field(f string) Value {
	return GetMethod(ProcedureMethods, v, f)
}

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
	pfun := func(env *Env, args ...Value) (Value, *Closure) {
		//  set up traceback recovery
		defer Traceback(name, args)
		//  convert fixed arguments from Goaldi values to needed Go type
		in := make([]reflect.Value, 0, len(args))
		var v reflect.Value
		for i := 0; i < nfixed; i++ {
			a := NilValue
			if i < len(args) {
				a = args[i]
			}
			v = passer[i](a)
			if !v.IsValid() {
				panic(&RunErr{"Cannot convert argument",
					args[i]})
			}
			in = append(in, v)
		}
		//  convert additional variadic arguments to final type
		if nfixed < nargs {
			for i := nfixed; i < len(args); i++ {
				v = passer[nfixed](args[i])
				if !v.IsValid() {
					panic(&RunErr{"Cannot convert argument",
						args[i]})
				}
				in = append(in, v)
			}
		}
		//  call the Go function
		out := fval.Call(in)
		//  return the result
		//  #%#% in the case of multiple return values, returns only the first!
		//  #%#% should at least panic if the second is type error and non-nil
		if nrtn >= 1 {
			return Import(out[0].Interface()), nil
		} else {
			return NilValue, nil
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
	// #%#% are other numeric types such as rune, int32, etc handled
	// #%#% acceptably by the default case below?
	case reflect.Int:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(int(v.(Numerable).ToNumber().Val()))
		}
	case reflect.Int64:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(int64(v.(Numerable).ToNumber().Val()))
		}
	case reflect.Float64:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(float64(v.(Numerable).ToNumber().Val()))
		}
	case reflect.String:
		return func(v Value) reflect.Value {
			return reflect.ValueOf(v.(Stringable).ToString().ToUTF8())
		}
	case reflect.Interface: // #%#% this assumes interface{}; should check
		// use default conversion
		break
	default:
		// check if convertible from numeric
		if reflect.TypeOf(1.0).ConvertibleTo(t) {
			return func(v Value) reflect.Value {
				return reflect.ValueOf(
					v.(Numerable).ToNumber().Val()).Convert(t)
			}
		}
		// otherwise, check if convertible from string
		if reflect.TypeOf("abc").ConvertibleTo(t) {
			return func(v Value) reflect.Value {
				return reflect.ValueOf(
					v.(Stringable).ToString().ToUTF8()).Convert(t)
			}
		}
		// otherwise, use default conversion
		break
	}
	// default conversion
	return func(v Value) reflect.Value {
		var inil interface{}
		x := Export(v) // default conversion
		if x == nil {
			return reflect.ValueOf(&inil).Elem() // nil is tricky
		} else {
			return reflect.ValueOf(x) // anything else
		}
	}
}

//  vproc.go -- VProcedure, the Goaldi type "procedure"
//
//  A VProcedure is created by the linker for each Go or Goaldi
//  procedure or method, and as a constructor for each VRecord.
//  Additional procedure can be created at runtime by Goaldi
//  "procedure" and "lambda" expressions.

package runtime

import (
	"fmt"
	"reflect"
	"strings"
)

var _ = fmt.Printf // enable debugging

//  Procedure value
type VProcedure struct {
	Name     string      // registered name
	Pnames   *[]string   // parameter names (nil if unknown)
	Variadic bool        // true if variadic
	RawCall  bool        // true if to use nonstandard raw argument lists
	GdProc   Procedure   // Goaldi-compatible function (possibly a shim)
	GoFunc   interface{} // underlying function
	Descr    string      // optional one-line description (used for stdlib)
}

const rProc = 50            // declare sort ranking
var _ ICore = &VProcedure{} // validate implementation

//  DefProc constructs a procedure from a Goaldi function and a description.
func DefProc(entry Procedure, name string, pspec string, descr string) *VProcedure {
	pnames, isvar := ParmsFromSpec(pspec)
	return NewProcedure(name, pnames, isvar, entry, entry, descr)
}

//  ParmsFromSpec turns a parameter spec into a pnames list and variadic flag
func ParmsFromSpec(pspec string) (*[]string, bool) {
	isvariadic := strings.HasSuffix(pspec, "[]")
	if isvariadic {
		pspec = strings.TrimSuffix(pspec, "[]")
	}
	pnames := strings.Split(pspec, ",")
	return &pnames, isvariadic
}

//  NewProcedure -- construct a procedure value
//  The result is variadic only if allowvar is true *and* entry is variadic.
func NewProcedure(name string, pnames *[]string, allowvar bool,
	entry Procedure, ufunc interface{}, descr string) *VProcedure {
	isvar := allowvar && reflect.TypeOf(entry).IsVariadic()
	return &VProcedure{name, pnames, isvar, false, entry, ufunc, descr}
}

//  ProcType is the procedure instance of type type.
var ProcType = NewType("procedure", "p", rProc, ProcCtor, nil,
	"proctype", "x", "succeed if procedure")

//  VProcedure.String -- default conversion to Go string returns "p:procname"
func (v *VProcedure) String() string {
	return "p:" + v.Name
}

//  VProcedure.GoString -- convert to string for image() and printf("%#v")
func (v *VProcedure) GoString() string {
	s := "procedure " + v.Name + "("
	if v.Pnames == nil {
		return s + "?)" // params unknown
	}
	d := ""
	for _, t := range *v.Pnames {
		s = s + d + t
		d = ","
	}
	if v.Variadic {
		s = s + "[]"
	}
	return s + ")"
}

//  VProcedure.Type -- return the procedure type
func (v *VProcedure) Type() IRank {
	return ProcType
}

//  VProcedure.Copy returns itself
func (v *VProcedure) Copy() Value {
	return v
}

//  VProcedure.Before compares two procs for sorting
func (a *VProcedure) Before(b Value, i int) bool {
	return a.Name < b.(*VProcedure).Name
}

//  VProcedure.Import returns itself
func (v *VProcedure) Import() Value {
	return v
}

//  VProcedure.Export returns the underlying function
//  (#%#% at least for now. should we wrap it somehow?)
func (v *VProcedure) Export() interface{} {
	return v.GdProc
}

//  VProcedure.Call invokes a procedure
func (v *VProcedure) Call(env *Env, args []Value, names []string) (Value, *Closure) {
	if v.RawCall {
		f := v.GoFunc.(func(*Env, []Value, []string) (Value, *Closure))
		return f(env, args, names)
	} else {
		args = ArgNames(v, args, names)
		return v.GdProc(env, args...)
	}
}

//  proctype(x) return x if x is a procedure, and fails otherwise.
//  proctype is the name of the result of main.type().
func ProcCtor(env *Env, args ...Value) (Value, *Closure) {
	x := ProcArg(args, 0, NilValue)
	if p, ok := x.(*VProcedure); ok {
		return Return(p)
	} else {
		return Fail()
	}
}

//  Go methods already converted to Goaldi procedures
var KnownMethods = make(map[uintptr]*VProcedure)

//  ImportMethod(val, name, meth) -- construct a Goaldi method from a Go method.
//  meth is a method struct such as returned by reflect.Type.MethodByName(),
//  not a bound method value such as returned by reflect.Value.MethodByName().
//  No GoShim flags are set, so errors and nil results get no special treatment.
func ImportMethod(val Value, name string, meth reflect.Method) Value {
	addr := meth.Func.Pointer()
	p := KnownMethods[addr]
	if p == nil {
		gofunc := meth.Func.Interface()
		proc := GoShim(name, gofunc, RNORM)
		p = NewProcedure(name, nil, true, proc, gofunc, "")
		KnownMethods[addr] = p
	}
	return MethodVal(p, Deref(val))
}

//  Flags for how GoShim should handle special return situations.
//  ETOSS and RNILF may both be set in which case ETOSS is applied first.
const (
	RNORM = 0 // normal return
	ETOSS = 1 // strip final error return and throw exception if not nil
	RNILF = 2 // turn a sole nil return value into failure
)

//  GoShim(name, func, retf) makes a shim for converting args to a Go function.
//  retf indicates the special handling, if any, of function returns.
func GoShim(name string, f interface{} /*func*/, retf int) Procedure {

	//  get information about the Go function
	ftype := reflect.TypeOf(f)
	fval := reflect.ValueOf(f)
	if fval.Kind() != reflect.Func {
		panic(NewExn("Not a func", f))
	}
	nargs := ftype.NumIn()
	nfixed := nargs
	if ftype.IsVariadic() {
		nfixed--
	}
	nrtn := ftype.NumOut()
	// clear inapplicable retf flags
	if nrtn == 0 {
		retf = 0 // no flags apply if there is no return
	} else {
		tn := ftype.Out(nrtn - 1)
		if tn.Name() != "error" || tn.PkgPath() != "" {
			retf = retf & ^ETOSS // last return is not of type "error"
		}
	}

	//  make an array of conversion functions, one per parameter
	passer := make([]func(Value) reflect.Value, nargs)
	for i := 0; i < nfixed; i++ {
		passer[i] = passfunc(ftype.In(i))
	}
	if nfixed < nargs { // if variadic
		passer[nfixed] = passfunc(ftype.In(nfixed).Elem())
	}

	// create a function that converts arguments and calls the underlying func
	return func(env *Env, args ...Value) (Value, *Closure) {
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
				panic(NewExn("Cannot convert argument", args[i]))
			}
			in = append(in, v)
		}
		//  convert additional variadic arguments to final type
		if nfixed < nargs {
			for i := nfixed; i < len(args); i++ {
				v = passer[nfixed](args[i])
				if !v.IsValid() {
					panic(NewExn("Cannot convert argument", args[i]))
				}
				in = append(in, v)
			}
		}
		//  call the Go function
		out := fval.Call(in)
		//  process the return values
		if nrtn == 0 {
			return Return(NilValue) // no return value: return nil
		}
		nrtn := nrtn // need private copy; may change later
		// if ETOSS is (still) set, check the final (or only) return value
		if (retf & ETOSS) != 0 {
			e := out[nrtn-1].Interface()
			if e != nil {
				panic(e) // throw error value as an exception
			} else if nrtn > 1 {
				nrtn-- // remove error return, keep the rest
			} else {
				return Return(NilValue) // nothing left, return nil
			}
		}
		// if there is (now) just one return value, return a simple value;
		// if RNILF is set, turn nil into failure
		if nrtn == 1 {
			r := Import(out[0].Interface()) // import the first return value
			if r == NilValue && (retf&RNILF) != 0 {
				return Fail()
			} else {
				return Return(r)
			}
		}
		// for multiple return values, make a list
		rlist := make([]Value, nrtn)
		for i := 0; i < nrtn; i++ {
			rlist[i] = Import(out[i].Interface())
		}
		return Return(InitList(rlist))
	}
}

//  passfunc returns a function that converts a Goaldi value
//  into a Go reflect.Value of the specified type
func passfunc(t reflect.Type) func(Value) reflect.Value {
	k := t.Kind()
	switch k {
	case reflect.Bool:
		return func(v Value) reflect.Value {
			var b bool
			switch x := v.(type) {
			case bool:
				b = x
			case vnil:
				b = false
			case *VNumber:
				b = (x.Val() != 0)
			default:
				b = true
			}
			return reflect.ValueOf(b)
		}
	case reflect.Interface: // #%#% this assumes interface{}; should check
		// use default conversion
		break
	default:
		// check if convertible from numeric
		if reflect.TypeOf(1.0).ConvertibleTo(t) {
			return func(v Value) reflect.Value {
				if reflect.TypeOf(v).ConvertibleTo(t) {
					return reflect.ValueOf(v).Convert(t)
				} else {
					return reflect.ValueOf(
						v.(Numerable).ToNumber().Val()).Convert(t)
				}
			}
		}
		// otherwise, check if convertible from string
		if reflect.TypeOf("abc").ConvertibleTo(t) {
			return func(v Value) reflect.Value {
				if reflect.TypeOf(v).ConvertibleTo(t) {
					return reflect.ValueOf(v).Convert(t)
				} else {
					return reflect.ValueOf(
						v.(Stringable).ToString().ToUTF8()).Convert(t)
				}
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

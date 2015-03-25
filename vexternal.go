//  vexternal.go -- interfacing with Go external values
//
//  External values are not a Goaldi type as such;
//  rather they are anything *not* a Goaldi type used as a Goladi value.

package goaldi

import (
	"fmt"
	"io"
	"reflect"
)

const rExternal = 99 // declare sort ranking

//  ExternalType defines the type "external", which is mostly just a stub
var ExternalType = NewType("external", "X", rExternal, External, nil,
	"external", "x", "export and re-import")

//  external(x) exports and then re-imports the value x.
func External(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("external", args)
	x := ProcArg(args, 0, NilValue)
	return Return(Import(Export(x)))
}

//  Import(x) builds a value of appropriate Goaldi type
func Import(x interface{}) Value {

	// must check first for a typed nil
	rv := reflect.ValueOf(x)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice:
		if rv.IsNil() {
			return NilValue
		}
	}

	switch v := x.(type) {

	case nil:
		return NilValue

	case IImport: // type declares its own import method incl Goaldi types
		return v.Import()
	case IExternal: // labels a user type that is to stay unconverted
		return v

	case bool:
		if v {
			return ONE
		} else {
			return ZERO
		}

	case string:
		return NewString(v)
	case []byte:
		return NewString(string(v))
	case []rune:
		return RuneString(v)

	case float32:
		return NewNumber(float64(v))
	case float64:
		return NewNumber(float64(v))
	case int:
		return NewNumber(float64(v)) //#%#% check vs MAX_EXACT?
	case int8:
		return NewNumber(float64(v))
	case int16:
		return NewNumber(float64(v))
	case int32:
		return NewNumber(float64(v))
	case int64:
		return NewNumber(float64(v)) //#%#% check vs MAX_EXACT?
	case uint:
		return NewNumber(float64(v)) //#%#% check vs MAX_EXACT?
	case uint8:
		return NewNumber(float64(v))
	case uint16:
		return NewNumber(float64(v))
	case uint32:
		return NewNumber(float64(v))
	case uint64:
		return NewNumber(float64(v)) //#%#% check vs MAX_EXACT?
	case uintptr:
		return NewNumber(float64(v)) //#%#% check vs MAX_EXACT?

	case io.Reader, io.Writer: // either reader or writer makes a file
		r, _ := x.(io.Reader)
		w, _ := x.(io.Writer)
		c, _ := x.(io.Closer)
		name := fmt.Sprintf("%T", x) // use type for name
		return NewFile(name, r, w, c)

	//#%#% add other cases?

	default:
		return x // external
	}
}

//  Export(v) returns the default Go representation of a Goaldi value
func Export(v Value) interface{} {
	if x, ok := v.(IExport); ok {
		return x.Export()
	} else {
		return v
	}
}

//  --------------------- trapped references (general) ----------------

type vGoTrap struct {
	target reflect.Value
}

func TrapValue(v reflect.Value) *vGoTrap {
	return &vGoTrap{v}
}

func (v *vGoTrap) Deref() Value {
	return Import(v.target.Interface())
}

func (v *vGoTrap) Assign(x Value) IVariable {
	v.target.Set(passfunc(v.target.Type())(x))
	return v
}

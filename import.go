//  import.go -- translation of Go values into Goaldi values
//  (#%#% experimental and incomplete #%#%)

package goaldi

import (
	"fmt"
	"io"
	"reflect"
)

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

//  --------------------- trapped references (for Go maps) ----------------

type vMapTrap struct {
	mapv reflect.Value // underlying Go map
	keyv reflect.Value // key converted to appropriate Go type
}

func TrapMap(mapv reflect.Value, key Value) *vMapTrap {
	return &vMapTrap{mapv, passfunc(mapv.Type().Key())(key)}
}

func (t *vMapTrap) Deref() Value {
	v := t.mapv.MapIndex(t.keyv)
	if v.IsValid() {
		return Import(v.Interface())
	} else {
		return nil // not found in map
	}
}

func (t *vMapTrap) Assign(x Value) IVariable {
	t.mapv.SetMapIndex(t.keyv, passfunc(t.mapv.Elem().Type())(x))
	return t
}

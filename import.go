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
		return NewString(string(v)) //#%#%???

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
		return NewForeign(x)
	}
}

//  Export(v) returns the default Go representation of a Goaldi value
func Export(v Value) interface{} {
	switch e := v.(type) {
	case IExport:
		return e.Export()
	default:
		return v
	}
}

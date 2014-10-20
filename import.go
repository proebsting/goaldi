//  import.go -- translation of Go values into Goaldi values
//  (#%#% experimental and incomplete #%#%)

package goaldi

import (
	"io"
	"reflect"
)

//  Import(x) builds a value of appropriate Goaldi type
func Import(x interface{}) Value {
	switch v := x.(type) {
	case IExternal: // labels a user type that is to stay unconverted
		return v
	case nil:
		return NewNil()
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
	//#%#% add other cases
	//#%#% see golang.org/src/pkg/fmt/print.go for reflection examples
	default:
		// if it implements Reader or Writer, make it a Goaldi file
		f := fileValue(x)
		if f != nil {
			return f
		}
		// unrecognized; use as is, but check for (typed) nil value
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice:
			if rv.IsNil() {
				return NilValue
			} else {
				return x
			}
		default:
			return x
		}
	}
}

//  fileValue(x) converts x to a Goaldi file if possible, else nil
func fileValue(x Value) *VFile {
	r, rok := x.(io.Reader)
	w, wok := x.(io.Writer)
	if !rok && !wok {
		return nil // neither Reader nor Writer
	}
	c, _ := x.(io.Closer)
	name := reflect.TypeOf(x).Name() // use underlying type as name
	return NewFile(name, r, w, c)
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

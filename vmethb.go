//  vmethb -- a bound method, or "method value"
//
//  This needs to be a first-class type because it is visible to a
//  Goaldi programmer.

package goaldi

import (
	"fmt"
	"reflect"
)

//  An VMethB is like a Go "method value", a function bound to an object,
//  for example the "m.delete" part of the expression "m.delete(x)"
type VMethB struct {
	Name string
	Val  Value
	Func interface{} // func(Value, ...Value)(Value, *Closure)
}

//  VMethB.String -- conversion to Go string returns "V:Name"
func (v *VMethB) String() string {
	return "V:" + v.Name
}

//  VMethB.GoString -- convert to Go string for image() and printf("%#v")
func (v *VMethB) GoString() string {
	return fmt.Sprintf("methodvalue (%v).%s", v.Val, v.Name)
}

//  VMethB.Rank returns rMethB
func (v *VMethB) Rank() int {
	return rMethB
}

//  VMethB.Type returns "methodvalue"
func (v *VMethB) Type() Value {
	return type_methodvalue
}

var type_methodvalue = NewString("methodvalue")

//  VMethB.Copy returns itself
func (v *VMethB) Copy() Value {
	return v
}

//  VMethB.Import returns itself
func (v *VMethB) Import() Value {
	return v
}

//  VMethB.Export returns itself
func (v *VMethB) Export() interface{} {
	return v
}

//  Declare required methods
var MethBMethods = map[string]interface{}{
	"type":   (*VMethB).Type,
	"copy":   (*VMethB).Copy,
	"string": (*VMethB).String,
	"image":  (*VMethB).GoString,
}

//  VMethB.Field implements methods
func (v *VMethB) Field(f string) Value {
	return GetMethod(MethBMethods, v, f)
}

//  GetMethod(m,v,s) looks up method v.s in table m, panicking on failure.
func GetMethod(m map[string]interface{}, v Value, s string) *VMethB {
	method := m[s]
	if method == nil {
		panic(&RunErr{"unrecognized method: " + s, v})
	}
	return &VMethB{s, v, method}
}

//  VMethB.Call(args) invokes the underlying method function.
func (mvf *VMethB) Call(env *Env, args ...Value) (Value, *Closure) {
	arglist := make([]reflect.Value, 1+len(args))
	arglist[0] = reflect.ValueOf(mvf.Val)
	for i, v := range args {
		arglist[i+1] = reflect.ValueOf(v)
	}
	method := reflect.ValueOf(mvf.Func)
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

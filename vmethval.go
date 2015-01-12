//  vmethval -- a "method value", a receiver bound with a method
//
//  This needs to be a first-class type because it is visible to a
//  Goaldi programmer, for example in "write(image([].pop))".

package goaldi

import (
	"fmt"
	"reflect"
)

//  An VMethVal is like a Go "method value", a function bound to an object,
//  for example the "m.delete" part of the expression "m.delete(x)"
type VMethVal struct {
	Name    string
	Val     Value
	Func    interface{} // func(Value, ...Value)(Value, *Closure)
	PassEnv bool        // pass environment when calling?
}

//  VMethVal.String -- conversion to Go string returns "V:Name"
func (v *VMethVal) String() string {
	return "V:" + v.Name
}

//  VMethVal.GoString -- convert to Go string for image() and printf("%#v")
func (v *VMethVal) GoString() string {
	return fmt.Sprintf("methodvalue (%v).%s", v.Val, v.Name)
}

//  VMethVal.Rank returns rMethVal
func (v *VMethVal) Rank() int {
	return rMethVal
}

//  VMethVal.Type returns "methodvalue"
func (v *VMethVal) Type() Value {
	return type_methodvalue
}

var type_methodvalue = NewString("methodvalue")

//  VMethVal.Copy returns itself
func (v *VMethVal) Copy() Value {
	return v
}

//  VMethVal.Identical -- check equality for === operator
func (s *VMethVal) Identical(x Value) Value {
	t, ok := x.(*VMethVal)
	if !ok {
		return nil // different types
	}
	if Identical(s.Val, t.Val) == nil {
		return nil // different values
	}
	if s.Name == t.Name {
		return s // same method of same value
	} else {
		return nil // different methods
	}
}

//  VMethVal.Import returns itself
func (v *VMethVal) Import() Value {
	return v
}

//  VMethVal.Export returns itself
func (v *VMethVal) Export() interface{} {
	return v
}

//  Declare required methods
var MethValMethods = MethodTable([]*GoProc{
	&GoProc{"type", (*VMethVal).Type, []string{}, "return methodvalue type"},
	&GoProc{"copy", (*VMethVal).Copy, []string{}, "return methodvalue"},
	&GoProc{"string", (*VMethVal).String, []string{}, "return short string"},
	&GoProc{"image", (*VMethVal).GoString, []string{}, "return string image"},
})

//  VMethVal.Field implements methods on methodvalues
func (v *VMethVal) Field(f string) Value {
	return GetMethod(MethValMethods, v, f)
}

//  GoProc describes a Go function to be used as a Goaldi procedure or method
type GoProc struct {
	Name   string      // name as seen from Goaldi
	Entry  interface{} // go func implmenting the procedure
	Pnames []string    // parameter names
	Descr  string      // one-line description
}

//  MethodTable makes a method table from a list of RecordMethods
func MethodTable(plist []*GoProc) map[string]*GoProc {
	t := make(map[string]*GoProc)
	for _, g := range plist {
		t[g.Name] = g
		//fmt.Printf("%s: %#v\n", g.Name, g.Entry)
	}
	return t
}

//  GetMethod(m,v,s) looks up method v.s in table m, panicking on failure.
func GetMethod(m map[string]*GoProc, v Value, s string) *VMethVal {
	method := m[s]
	if method == nil {
		panic(NewExn("unrecognized method: "+s, v))
	}
	return &VMethVal{s, v, method.Entry, false}
}

//  VMethVal.Call(args) invokes the underlying method function.
func (mvf *VMethVal) Call(env *Env, args ...Value) (Value, *Closure) {
	arglist := make([]reflect.Value, 2+len(args))
	arglist[0] = reflect.ValueOf(env)
	arglist[1] = reflect.ValueOf(mvf.Val)
	for i, v := range args {
		arglist[i+2] = reflect.ValueOf(v)
	}
	method := reflect.ValueOf(mvf.Func)
	if !mvf.PassEnv {
		arglist = arglist[1:]
	}
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

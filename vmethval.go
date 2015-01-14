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
	Name   string
	Val    Value
	Func   interface{} // func(Value, ...Value)(Value, *Closure)
	Pnames *[]string   // list of parameter names, if known
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
var MethValMethods = MethodTable([]*VProcedure{
	DefMeth("type", (*VMethVal).Type, []string{}, "return methodvalue type"),
	DefMeth("copy", (*VMethVal).Copy, []string{}, "return methodvalue"),
	DefMeth("string", (*VMethVal).String, []string{}, "return short string"),
	DefMeth("image", (*VMethVal).GoString, []string{}, "return string image"),
})

//  VMethVal.Field implements methods on methodvalues
func (v *VMethVal) Field(f string) Value {
	return GetMethod(MethValMethods, v, f)
}

//  DefMeth defines a method implemented in Go as a VProcedure
func DefMeth(name string, entry interface{}, pnames []string, descr string) *VProcedure {
	p := NewProcedure(name, &pnames, true, nil, entry, descr)
	p.IsMethod = true
	return p
}

//  MethodTable makes a method table from a list of VProcedures
func MethodTable(plist []*VProcedure) map[string]*VProcedure {
	t := make(map[string]*VProcedure)
	for _, g := range plist {
		t[g.Name] = g
	}
	return t
}

//  GetMethod(m,v,s) looks up method v.s in table m, panicking on failure.
func GetMethod(m map[string]*VProcedure, v Value, s string) *VMethVal {
	method := m[s]
	if method == nil {
		panic(NewExn("unrecognized method: "+s, v))
	}
	return &VMethVal{s, v, method.GoFunc, method.Pnames}
}

//  VMethVal.Call invokes the underlying method function.
func (mvf *VMethVal) Call(env *Env, args []Value, names []string) (Value, *Closure) {
	args = ArgNames(args, names, mvf, mvf.Pnames)
	arglist := make([]reflect.Value, 2+len(args))
	arglist[0] = reflect.ValueOf(env)
	arglist[1] = reflect.ValueOf(mvf.Val)
	for i, v := range args {
		arglist[i+2] = reflect.ValueOf(v)
	}
	method := reflect.ValueOf(mvf.Func)
	mtype := reflect.TypeOf(mvf.Func)
	if mtype.NumIn() == 0 || !reflect.TypeOf(env).AssignableTo(mtype.In(0)) {
		arglist = arglist[1:] // remove env argument if not wanted
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

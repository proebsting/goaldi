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
	Proc *VProcedure
	Val  Value
}

//  MethodVal(p,v) builds a VMethVal struct representing the expression p.v
func MethodVal(p *VProcedure, v Value) *VMethVal {
	return &VMethVal{p, v}
}

//  MethValType is the methodvalue instance of type type.
var MethValType = NewType("m", rMethVal, MethodValue, nil,
	"methodvalue", "m", "succeed if methodvalue")

//  VMethVal.String -- conversion to Go string returns "m:Name"
func (v *VMethVal) String() string {
	return "m:" + v.Proc.Name
}

//  VMethVal.GoString -- convert to Go string for image() and printf("%#v")
func (v *VMethVal) GoString() string {
	return fmt.Sprintf("methodvalue (%v).%s", v.Val, v.Proc.Name)
}

//  VMethVal.Type returns the meethodvalue type
func (v *VMethVal) Type() IRank {
	return MethValType
}

//  VMethVal.Copy returns itself
func (v *VMethVal) Copy() Value {
	return v
}

//  VMethVal.Before compares two methodvalues for sorting
func (a *VMethVal) Before(x Value, i int) bool {
	b := x.(*VMethVal)
	if a.Proc != b.Proc {
		return a.Proc.Before(b.Proc, i)
	} else {
		return LT(a.Val, b.Val, -1)
	}
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
	if s.Proc == t.Proc {
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

//  The "constructor" returns its argument if methodvalue and otherwise fails.
func MethodValue(env *Env, args ...Value) (Value, *Closure) {
	x := ProcArg(args, 0, NilValue)
	if v, ok := x.(*VMethVal); ok {
		return Return(v)
	} else {
		return Fail()
	}
}

//  DefMeth defines a method implemented in Go as a VProcedure
func DefMeth(entry interface{}, name string, pspec string, descr string) *VProcedure {
	pnames, isvar := ParmsFromSpec(pspec)
	return NewProcedure(name, pnames, isvar, nil, entry, descr)
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
//  Defaults methods are additionally provided for all types.
func GetMethod(m map[string]*VProcedure, v Value, s string) *VMethVal {
	if m != nil {
		method := m[s]
		if method != nil {
			return MethodVal(method, v)
		}
	}
	if mv := UniMethod(v, s); mv != nil {
		return mv
	}
	panic(NewExn("unrecognized method: "+s, v))
}

//  UniMethod(v,s) finds one of the universal methods defined on all types
func UniMethod(v Value, s string) *VMethVal {
	if m := UniMethods[s]; m != nil {
		return MethodVal(m, v)
	} else {
		return nil
	}
}

//  Declare universal methods
var UniMethods = map[string]*VProcedure{
	"type":     DefProc(Type, "type", "", "return type of value"),
	"string":   DefProc(String, "string", "", "render value as string"),
	"image":    DefProc(Image, "image", "", "return detailed string image"),
	"copy":     DefProc(Copy, "copy", "", "copy value"),
	"external": DefProc(External, "external", "", "export and re-import"),
}

//  VMethVal.Call invokes the underlying method function.
func (mvf *VMethVal) Call(env *Env, args []Value, names []string) (Value, *Closure) {
	p := mvf.Proc // procedure description
	args = ArgNames(p, args, names)
	arglist := make([]reflect.Value, 2+len(args))
	arglist[0] = reflect.ValueOf(env)
	arglist[1] = reflect.ValueOf(mvf.Val)
	for i, v := range args {
		arglist[i+2] = reflect.ValueOf(v)
	}
	f := p.GoFunc // there's always a Go function
	if p.GdProc != nil {
		f = p.GdProc // use the Goaldi version if provided
	}
	method := reflect.ValueOf(f)
	mtype := reflect.TypeOf(f)
	if mtype.NumIn() == 0 || !mtype.In(0).ConvertibleTo(reflect.TypeOf(env)) {
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

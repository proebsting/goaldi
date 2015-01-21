//  vctor.go -- record constructor information
//
//  Defines the interpretation of a VRecord object that points to it,
//  and constructs instances of it.

package goaldi

import (
	"fmt"
	"regexp"
)

var _ = fmt.Printf // enable debugging

//  VCtor is the constructor structure
type VCtor struct {
	Name    string                 // type name
	Flist   []string               // ordered list of field names
	Ctor    *VProcedure            // pseudo-constructor for argname handling
	Members map[string]interface{} // field and method table
}

//  ConstructorType is the constructor instance of type type
var ConstructorType = NewType("R", rCtor, Constructor, nil,
	"constructor", "name,fields[]", "build a record constructor")

//  A Constructor is also a type, which means it must implement Rank()
func (v *VCtor) Rank() int {
	return rRecord // if this is a type, its value is a record
}

//  NewCtor(name, fields) -- construct new definition
//  Panics if a field name is duplicated.
func NewCtor(name string, fields []string) *VCtor {
	cproc := NewProcedure(name, &fields, false, nil, (*VCtor).New, "")
	ctor := &VCtor{name, fields, cproc, make(map[string]interface{})}
	for i, s := range fields {
		if ctor.Members[s] != nil {
			panic(NewExn("duplicate field name", s))
		}
		ctor.Members[s] = i // enter field-to-index mapping
	}
	return ctor
}

//  AddMethod(name, procedure) -- add a method for this record type
//  Returns false if rejected as a duplicate.
func (v *VCtor) AddMethod(name string, vproc *VProcedure) bool {
	if v.Members[name] != nil {
		return false // this is a duplicate
	}
	p := *vproc               // copy original VProcedure struct
	p.Name = name             // set unqualified name
	pnames := (*p.Pnames)[1:] // trim explicit "self" parameter
	p.Pnames = &pnames        // and store updated list
	if v.Members[name] != nil {
		return false
	}
	v.Members[name] = &p
	return true
}

//  VCtor.New(values) -- create a new underlying record object
func (v *VCtor) New(a []Value) *VRecord {
	r := &VRecord{v, make([]Value, len(v.Flist))}
	for i := range r.Data {
		if i < len(a) {
			r.Data[i] = a[i]
		} else {
			r.Data[i] = NilValue
		}
	}
	return r
}

//  VCtor.String -- conversion to Go string returns "R:name"
func (v *VCtor) String() string {
	return "R:" + v.Name
}

//  VCtor.GoString -- convert to Go string for image() and printf("%#v")
func (v *VCtor) GoString() string {
	s := "constructor " + v.Name + "("
	d := ""
	for _, t := range v.Flist {
		s = s + d + t
		d = ","
	}
	return s + ")"
}

//  VCtor.Type returns the constructor type
func (v *VCtor) Type() IRank {
	return ConstructorType
}

//  VCtor.Copy returns itself
func (v *VCtor) Copy() Value {
	return v
}

//  VCtor.Import returns itself
func (v *VCtor) Import() Value {
	return v
}

//  VCtor.Export returns itself
func (v *VCtor) Export() interface{} {
	return v
}

//  VCtor.Dispense() implements !D to generate the field names
func (v *VCtor) Dispense(unused Value) (Value, *Closure) {
	var c *Closure
	i := -1
	c = &Closure{func() (Value, *Closure) {
		i++
		if i < len(v.Flist) {
			return NewString(v.Flist[i]), c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

//  VCtor.Call() implements a record constructor
func (v *VCtor) Call(env *Env, args []Value, names []string) (Value, *Closure) {
	args = ArgNames(v.Ctor, args, names)
	return Return(v.New(args))
}

//  Constructor(name, fields[]) builds a record constructor dynamically
func Constructor(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("constructor", args)
	name := Identifier(ProcArg(args, 0, NilValue))
	fields := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		fields[i-1] = Identifier(args[i])
	}
	return Return(NewCtor(name, fields))
}

//  Identifier converts its argument to a Go string and validates its form
func Identifier(x Value) string {
	s := x.(Stringable).ToString().ToUTF8()
	if !idPattern.MatchString(s) {
		panic(NewExn("Not an identifier", s))
	}
	return s
}

var idPattern = regexp.MustCompile("^[A-Za-z_][[0-9A-Za-z_]*$")

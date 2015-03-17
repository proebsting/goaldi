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
	VType                 // embedded type struct (VCtor extends VType, sort of)
	Parent *VCtor         // parent type
	Flist  []string       // ordered list of field names
	Fmap   map[string]int // map of names to indexes (1-based)
}

//  NewCtor(name, parent, fields) -- make constructor: name extends parent
//  Panics if a field name is duplicated.
func NewCtor(name string, parent *VCtor, newfields []string) *VCtor {

	// combine the parent's fields with the new fields
	fields := []string{}
	if parent != nil {
		fields = append(fields, parent.Flist...)
	}
	fields = append(fields, newfields...)

	// now build the structures
	cproc := NewProcedure(name, &fields, false, nil, (*VCtor).New, "")
	meth := make(map[string]*VProcedure)
	fmap := make(map[string]int)
	ctype := VType{name, "R", rRecord, cproc, meth}
	ctor := &VCtor{ctype, parent, fields, fmap}
	for i, s := range fields {
		if ctor.Methods[s] != nil || ctor.Fmap[s] != 0 {
			panic(NewExn("duplicate field name", s))
		}
		ctor.Fmap[s] = i + 1 // enter field-to-index mapping
	}
	return ctor
}

//  AddMethod(name, procedure) -- add a method for this record type
//  Returns false if rejected as a duplicate.
func (v *VCtor) AddMethod(name string, vproc *VProcedure) bool {
	if v.Methods[name] != nil {
		return false // this is a duplicate
	}
	p := *vproc               // copy original VProcedure struct
	p.Name = name             // set unqualified name
	pnames := (*p.Pnames)[1:] // trim explicit "self" parameter
	p.Pnames = &pnames        // and store updated list
	if v.Methods[name] != nil {
		return false
	}
	v.Methods[name] = &p
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

//  Declare static constructor
func init() {
	DefLib(Constructor,
		"constructor", "name,fields[]", "build a record constructor")
}

//  VCtor.Field -- implement C.id to override methods in VType
func (c *VCtor) Field(f string) Value {
	// check first for field index
	i := c.Fmap[f]
	if i > 0 {
		return NewNumber(float64(i)) // return Goaldi index of field f
	}
	// next check for universal method
	// (must pass VCtor, not Vtype, for c.copy() or c.image())
	m := UniMethod(c, f)
	if m != nil {
		return m
	}
	// finally fall back to methods defined by VType
	return GetMethod(TypeMethods, &c.VType, f)
}

//  VCtor.Size() returns the number of fields.
func (c *VCtor) Size() Value {
	return NewNumber(float64(len(c.Flist)))
}

//  VCtor.Index -- implement C[x] to return name of field i.
func (c *VCtor) Index(lval Value, x Value) Value {
	i, isNumber := c.Lookup(x)
	if i < 0 {
		return nil // fail: subscript out of range
	}
	if isNumber {
		// for numeric argument, return field name
		return NewString(c.Flist[i])
	} else {
		// for string argument, return corresponding Goaldi index
		return NewNumber(float64(i + 1))
	}
}

//  VCtor.Lookup(x) converts x (s or n) to zero-based Go index, or -1 to fail.
func (c *VCtor) Lookup(x Value) (index int, isNumber bool) {
	n := len(c.Flist)
	// if this is a string, check first for matching field name
	if s, ok := x.(*VString); ok {
		key := s.ToUTF8()
		i := c.Fmap[key]
		if i > 0 {
			return i - 1, false
		} else {
			return -1, false // fail: name not found
		}
		k := s.TryNumber()
		if k == nil {
			return -1, false // fail: unmatched string, not a number
		}
		x = k
	}
	// not a string; must be a number, else throw error
	i := int(x.(Numerable).ToNumber().Val())
	i = GoIndex(i, n)
	if i < n {
		return i, true // in range
	} else {
		return -1, true // fail: not in range
	}
}

//  VCtor.GoString -- convert to Go string for image() and printf("%#v")
func (v *VCtor) GoString() string {
	s := "constructor " + v.TypeName + "("
	d := ""
	for _, t := range v.Flist {
		s = s + d + t
		d = ","
	}
	return s + ")"
}

//  VCtor.Copy returns itself
func (v *VCtor) Copy() Value {
	return v
}

//  VCtor.Before compares itself with a constructor or type value
func (a *VCtor) Before(b Value, i int) bool {
	switch t := b.(type) {
	case *VCtor:
		return a.TypeName < t.TypeName
	case *VType:
		return rRecord < t.SortRank
	default:
		panic(Malfunction("unexpected type in VCtor.Before"))
	}
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

//  constructor(name, field...) builds a record constructor for creating
//  records with the given type name and field list.
//  There is no requirement or guarantee that record names be distinct.
func Constructor(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("constructor", args)
	name := Identifier(ProcArg(args, 0, NilValue))
	fields := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		fields[i-1] = Identifier(args[i])
	}
	return Return(NewCtor(name, nil, fields))
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

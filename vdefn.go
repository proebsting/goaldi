//  vdefn.go -- struct definition information
//
//  Defines the interpretation of a vstruct object that points to it,
//  and constructs objects for it.

package goaldi

import (
	"fmt"
)

type VDefn struct {
	Name  string   // type name
	Flist []string // ordered list of fields names
	//#%#% could add a hash map... but is it worth it?
}

//  NewDefn(name, fields) -- construct new definition
func NewDefn(name string, fields []string) *VDefn {
	return &VDefn{name, fields}
}

//  VDefn.New() -- create a new underlying struct object
func (v *VDefn) New() *VStruct {
	r := &VStruct{v, make([]Value, len(v.Flist))}
	for i := range r.Data {
		r.Data[i] = NilValue
	}
	return r
}

//  VDefn.String -- conversion to Go string returns "name:"
func (v *VDefn) String() string {
	return v.Name + ":"
}

//  VDefn.GoString -- image returns "name:n"
func (v *VDefn) GoString() string {
	return fmt.Sprintf("%s:%d", v.Name, len(v.Flist))
}

//  VDefn.Type returns "defn"
func (v *VDefn) Type() Value {
	return type_defn
}

var type_defn = NewString("defn")

//  VDefn.Copy returns itself
func (v *VDefn) Copy() Value {
	return v
}

//  VDefn.Import returns itself
func (v *VDefn) Import() Value {
	return v
}

//  VDefn.Export returns itself
func (v *VDefn) Export() interface{} {
	return v
}

//  VDefn.Dispense() implements !D to generate the field names
func (v *VDefn) Dispense(unused IVariable) (Value, *Closure) {
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

//  VDefn.Call() implements a struct constructor  //#%#%#% TO BE WRITTEN

//  Declare required methods
var Defnethods = map[string]interface{}{
	"type":  (*VDefn).Type,
	"copy":  (*VDefn).Copy,
	"image": Image,
}

//  VDefn.Field implements methods
func (v *VDefn) Field(f string) Value {
	return GetMethod(Defnethods, v, f)
}

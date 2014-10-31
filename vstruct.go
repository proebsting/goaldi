//  vstruct.go -- a user-defined (usually) structure

package goaldi

import (
	"bytes"
	"fmt"
)

type VStruct struct {
	Defn *VDefn  // underlying struct definition
	Data []Value // current data values
}

//  VStruct.String -- conversion to Go string returns "name{}"
func (v *VStruct) String() string {
	return v.Defn.Name + "{}"
}

//  VStruct.GoString -- returns string for image() and printf("%#v")
func (v *VStruct) GoString() string {
	if len(v.Data) == 0 {
		return v.Defn.Name + "{}"
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s{", v.Defn.Name)
	for _, x := range v.Data {
		fmt.Fprintf(&b, "%v,", x)
	}
	s := b.Bytes()
	s[len(s)-1] = '}'
	return string(s)
}

//  VStruct.Type returns the defined struct name
func (v *VStruct) Type() Value {
	return NewString(v.Defn.Name)
}

//  VStruct.Copy returns a distinct copy of itself
func (v *VStruct) Copy() Value {
	r := &VStruct{v.Defn, make([]Value, len(v.Data))}
	copy(r.Data, v.Data)
	return r
}

//  VStruct.Import returns itself
func (v *VStruct) Import() Value {
	return v
}

//  VStruct.Export returns itself
func (v *VStruct) Export() interface{} {
	return v
}

//  VStruct.Dispense() implements !D to generate the field values
func (v *VStruct) Dispense(unused IVariable) (Value, *Closure) {
	var c *Closure
	i := -1
	c = &Closure{func() (Value, *Closure) {
		i++
		if i < len(v.Data) {
			return v.Data[i], c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

//  VStruct.Call() implements a struct constructor  //#%#%#% TO BE WRITTEN

//  VStruct.Field() implements a field reference
func (v *VStruct) Field(f string) Value {
	d := v.Defn
	for i, s := range d.Flist {
		if s == f {
			return Trapped(&v.Data[i])
		}
	}
	panic(&RunErr{"Field not found: " + f, v})
}

//  ostruct.go -- operations on user-defined structures

package goaldi

import (
	"math/rand"
)

//  Declare standard methods
var RecordMethods = map[string]interface{}{
	"type":   (*VRecord).Type,
	"copy":   (*VRecord).Copy,
	"string": (*VRecord).String,
	"image":  (*VRecord).GoString,
}

//  VRecord.Field() implements a field reference R.k
func (v *VRecord) Field(f string) Value {
	//  check first for record field
	d := v.Defn
	for i, s := range d.Flist {
		if s == f {
			return Trapped(&v.Data[i])
		}
	}
	//  check for explicit method
	if m := d.Methods[f]; m != nil {
		return &VMethVal{f, v, m.Entry, true}
	}
	//  check for standard method
	if m := RecordMethods[f]; m != nil {
		return &VMethVal{f, v, m, false}
	}
	//  neither one found
	panic(NewExn("Field not found: "+f, v))
}

//  VRecord.Index(u, k) implements an indexed reference R[k]
func (v *VRecord) Index(lval Value, x Value) Value {
	n := len(v.Data)
	// if this is a string, check first for matching field name
	if s, ok := x.(*VString); ok {
		key := s.ToUTF8()
		for i, f := range v.Defn.Flist {
			if f == key {
				if lval == nil {
					return v.Data[i]
				} else {
					return Trapped(&v.Data[i])
				}
			}
		}
		k := s.TryNumber()
		if k == nil {
			return nil // fail: unmatched string, not a number
		}
		x = k
	}
	// otherwise index by number
	i := int(x.(Numerable).ToNumber().Val())
	i = GoIndex(i, n)
	if i >= n {
		return nil // fail: subscript out of range
	} else if lval == nil {
		return v.Data[i]
	} else {
		return Trapped(&v.Data[i])
	}
}

//  VRecord.Size() implements *R, returning the number of fields
func (v *VRecord) Size() Value {
	return NewNumber(float64(len(v.Data)))
}

//  VRecord.Choose() implements ?R
func (v *VRecord) Choose(lval Value) Value {
	n := len(v.Data)
	if n == 0 {
		return nil
	} else if lval == nil {
		return v.Data[rand.Intn(n)]
	} else {
		return Trapped(&v.Data[rand.Intn(n)])
	}
}

//  VRecord.Dispense() implements !R to generate the field values
func (v *VRecord) Dispense(lval Value) (Value, *Closure) {
	var c *Closure
	i := -1
	c = &Closure{func() (Value, *Closure) {
		i++
		if i >= len(v.Data) {
			return Fail()
		} else if lval == nil {
			return v.Data[i], c
		} else {
			return Trapped(&v.Data[i]), c
		}
	}}
	return c.Resume()
}
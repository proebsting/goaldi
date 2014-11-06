//  ostruct.go -- operations on user-defined structures

package goaldi

import (
	"math/rand"
)

//  Declare standard methods
var StructMethods = map[string]interface{}{
	"type":  (*VStruct).Type,
	"copy":  (*VStruct).Copy,
	"image": Image,
}

//  VStruct.Field() implements a field reference S.k
func (v *VStruct) Field(f string) Value {
	//  check first for record field
	d := v.Defn
	for i, s := range d.Flist {
		if s == f {
			return Trapped(&v.Data[i])
		}
	}
	//  check for standard method
	if m := StructMethods[f]; m != nil {
		return &VMethB{f, v, m}
	}
	//  neither one found
	panic(&RunErr{"Field not found: " + f, v})
}

//  VStruct.Index(u, k) implements an indexed reference S[k]
func (v *VStruct) Index(lval IVariable, x Value) Value {
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

//  VStruct.Size() implements *S, returning the number of fields
func (v *VStruct) Size() Value {
	return NewNumber(float64(len(v.Data)))
}

//  VStruct.Choose() implements ?S
func (v *VStruct) Choose(lval IVariable) Value {
	n := len(v.Data)
	if n == 0 {
		return nil
	} else if lval == nil {
		return v.Data[rand.Intn(n)]
	} else {
		return Trapped(&v.Data[rand.Intn(n)])
	}
}

//  VStruct.Dispense() implements !S to generate the field values
func (v *VStruct) Dispense(lval IVariable) (Value, *Closure) {
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

//  ostruct.go -- operations on user-defined structures

package goaldi

import (
	"math/rand"
)

//  VStruct.Field() implements a field reference S.k
func (v *VStruct) Field(f string) Value {
	d := v.Defn
	for i, s := range d.Flist {
		if s == f {
			return Trapped(&v.Data[i])
		}
	}
	panic(&RunErr{"Field not found: " + f, v})
}

//  VStruct.Index(u, k) implements an indexed reference S[k]
func (v *VStruct) Index(unused IVariable, x Value) Value {
	n := len(v.Data)
	// if this is a string, check first for matching field name
	if s, ok := x.(*VString); ok {
		key := s.ToUTF8()
		for i, f := range v.Defn.Flist {
			if f == key {
				return Trapped(&v.Data[i])
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
	if i < n {
		return Trapped(&v.Data[i])
	} else {
		return nil // fail: subscript out of range
	}
}

//  VStruct.Size() implements *S, returning the number of fields
func (v *VStruct) Size() Value {
	return NewNumber(float64(len(v.Data)))
}

//  VStruct.Choose() implements ?S
func (v *VStruct) Choose(unused IVariable) Value {
	n := len(v.Data)
	if n == 0 {
		return nil
	} else {
		return Trapped(&v.Data[rand.Intn(n)])
	}
}

//  VStruct.Dispense() implements !S to generate the field values
func (v *VStruct) Dispense(unused IVariable) (Value, *Closure) {
	var c *Closure
	i := -1
	c = &Closure{func() (Value, *Closure) {
		i++
		if i < len(v.Data) {
			return Trapped(&v.Data[i]), c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

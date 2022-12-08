//  orecord.go -- operations on user-defined Goaldi record structures

package runtime

import (
	"fmt"
	"math/rand"
)

var _ = fmt.Printf // enable debugging

// VRecord.Field() implements a field reference R.k
func (v *VRecord) Field(f string) Value {
	d := v.Ctor
	i := d.Fmap[f]
	if i > 0 {
		return Trapped(&v.Data[i-1])
	}
	for d != nil {
		m := d.Methods[f]
		if m != nil {
			return MethodVal(m, v)
		}
		d = d.Parent
	}
	//  check for standard method
	if mv := UniMethod(v, f); mv != nil {
		return mv
	}
	//  nothing found
	panic(NewExn("Field not found: "+f, v))
}

// VRecord.Index(lval, x) implements an indexed reference R[x]
func (v *VRecord) Index(lval Value, x Value) Value {
	i, _ := v.Ctor.Lookup(x)
	if i < 0 {
		return nil // fail: not found
	} else if lval == nil {
		return v.Data[i] // return value
	} else {
		return Trapped(&v.Data[i]) // return trapped lvalue
	}
}

// VRecord.Size() implements *R, returning the number of fields
func (v *VRecord) Size() Value {
	return NewNumber(float64(len(v.Data)))
}

// VRecord.Choose() implements ?R
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

// VRecord.Dispense() implements !R to generate the field values
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

//  olist.go -- list operations

package goaldi

import (
	"math/rand"
)

//------------------------------------  Size:  *e

func (v *VList) Size() Value {
	return NewNumber(float64(len(v.data)))
}

//------------------------------------  Choose:  ?e

func (v *VList) Choose(lval IVariable) Value {
	n := len(v.data)
	if n == 0 {
		return nil // fail
	} else {
		return &vListRef{v, rand.Intn(n)}
	}
}

//------------------------------------  Dispense:  !e

func (v *VList) Dispense(lval IVariable) (Value, *Closure) {
	i := -1
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		i++
		if i < len(v.data) {
			return &vListRef{v, i}, c
		} else {
			return nil, nil
		}
	}}
	return c.Resume()
}

//------------------------------------  Index:  e1[e2]

func (v *VList) Index(lval IVariable, x Value) Value {
	n := len(v.data)
	i := int(x.(Numerable).ToNumber().Val())
	i = GoIndex(i, n)
	if i < n {
		return &vListRef{v, i}
	} else {
		return nil // fail: subscript out of range
	}
}

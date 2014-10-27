//  omap.go -- map operations

package goaldi

import (
	"math/rand"
)

//------------------------------------  Size:  *e

func (v *VMap) Size() Value {
	return NewNumber(float64(len(v.data)))
}

//------------------------------------  Choose:  ?e

func (v *VMap) Choose(lval IVariable) Value {
	n := len(v.klist)
	if n == 0 {
		return nil
	} else {
		return v.klist[rand.Intn(n)]
	}
}

//------------------------------------  Dispense:  !e

func (v *VMap) Dispense(lval IVariable) (Value, *Closure) {
	i := 0
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		if i >= len(v.klist) {
			return Fail()
		}
		v := v.klist[i]
		i++
		return v, c
	}}
	return c.Resume()
}

//------------------------------------  Index:  e1[e2]

func (v *VMap) Index(lval IVariable, key Value) Value {
	return &vMapSlot{v, key}
}

//------------------------------------  Member:  e1.member(32)

func (v *VMap) Member(key Value) Value {
	if v.data[MapIndex(key)] != nil {
		return key
	} else {
		return nil
	}
}

//------------------------------------  Delete:  e1.delete(32)

func (v *VMap) Delete(key Value) Value {
	x := MapIndex(key)
	delete(v.data, x)
	if len(v.data) != len(v.klist) {
		// delete was successful; need to remove from klist
		for i := len(v.klist) - 1; i >= 0; i-- {
			if Identical(key, v.klist[i]) != nil {
				v.klist[i] = v.klist[len(v.klist)-1]
				v.klist = v.klist[:len(v.klist)-1]
				break
			}
		}
	}
	if len(v.data) != len(v.klist) {
		panic(&RunErr{"inconsistent map", v})
	}
	return v
}

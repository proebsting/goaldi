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

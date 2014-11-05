//  omap.go -- map operations

package goaldi

import (
	"math/rand"
)

var kvstruct = NewDefn("mapElem", []string{"key", "value"})

//  VMap.Entry(i) -- return an initialized {key,value} struct copying entry i
func (m *VMap) Entry(i int) *VStruct {
	k := m.klist[i]
	return kvstruct.New([]Value{k, m.data[MapIndex(k)]})
}

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
		return v.Entry(rand.Intn(n))
	}
}

//------------------------------------  Dispense:  !e

func (v *VMap) Dispense(lval IVariable) (Value, *Closure) {
	i := -1
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		i++
		if i < len(v.klist) {
			return v.Entry(i), c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

//------------------------------------  Index:  e1[e2]

func (v *VMap) Index(lval IVariable, key Value) Value {
	return &vMapSlot{v, key}
}

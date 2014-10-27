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

//------------------------------------  Field:  e1.s

func (v *VMap) Field(f string) Value {
	//#%#% checking first for "member" and "delete",
	//#%#% allowing any other string as a index
	switch f {
	case "member":
		return MVFunc(v.Member)
	case "delete":
		return MVFunc(v.Delete)
	default:
		return &vMapSlot{v, NewString(f)}
	}
}

//------------------------------------  Member:  e1.member(e2)

func (v *VMap) Member(args ...Value) (Value, *Closure) {
	key := args[0]
	if v.data[MapIndex(key)] != nil {
		return Return(key)
	} else {
		return Fail()
	}
}

//------------------------------------  Delete:  e1.delete(e2)

func (v *VMap) Delete(args ...Value) (Value, *Closure) {
	key := args[0]
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
	return Return(v)
}

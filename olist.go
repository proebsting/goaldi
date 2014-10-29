//  olist.go -- list operations

package goaldi

import (
	"math/rand"
)

//------------------------------------  Size:  *L

func (v *VList) Size() Value {
	return NewNumber(float64(len(v.data)))
}

//------------------------------------  Choose:  ?L

func (v *VList) Choose(lval IVariable) Value {
	n := len(v.data)
	if n == 0 {
		return nil // fail
	} else {
		return &vListRef{v, rand.Intn(n)}
	}
}

//------------------------------------  Dispense:  !L

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

//------------------------------------  Index:  L[i]

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

//------------------------------------  Slice:  L[i:j]

func (v *VList) Slice(lval IVariable, x Value, y Value) Value {
	i := int(x.(Numerable).ToNumber().Val())
	j := int(y.(Numerable).ToNumber().Val())
	n := len(v.data)
	i = GoIndex(i, n)
	j = GoIndex(j, n)
	if i > n || j > n {
		return nil // subscript out of range
	}
	if i > j {
		i, j = j, i // indexing was backwards
	}
	m := j - i
	a := make([]Value, m, m)
	if v.rev {
		copy(a, v.data[n-j:n-i])
		ReverseValues(a)
	} else {
		copy(a, v.data[i:j])
	}
	return InitList(a)
}

//------------------------------------  ListCat:  L1 ||| L2

type IListCat interface {
	ListCat(Value) Value
}

func (v *VList) ListCat(x Value) Value {
	return InitList(append(v.Export().([]Value), x.(*VList).Export().([]Value)...))
}

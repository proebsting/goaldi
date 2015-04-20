//  olist.go -- list operations

package runtime

import (
	"math/rand"
)

//------------------------------------  Size:  *L

func (v *VList) Size() Value {
	return NewNumber(float64(len(v.data)))
}

//------------------------------------  Take:  @L

func (v *VList) Take(lval Value) Value {
	r, _ := v.Snip(true, "@L", nil)
	return r
}

//------------------------------------  Choose:  ?L

func (v *VList) Choose(lval Value) Value {
	n := len(v.data)
	if n == 0 {
		return nil // fail
	} else if lval == nil {
		return v.data[rand.Intn(n)]
	} else {
		return &vListRef{v, rand.Intn(n)}
	}
}

//------------------------------------  Dispense:  !L

func (v *VList) Dispense(lval Value) (Value, *Closure) {
	i := -1
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		i++
		if i >= len(v.data) {
			return nil, nil
		} else if lval != nil {
			return &vListRef{v, i}, c
		} else {
			return v.Elem(i), c
		}
	}}
	return c.Resume()
}

//------------------------------------  Index:  L[i]

func (v *VList) Index(lval Value, x Value) Value {
	n := len(v.data)
	i := int(x.(Numerable).ToNumber().Val())
	i = GoIndex(i, n)
	if i >= n {
		return nil // fail: subscript out of range
	} else if lval != nil {
		return &vListRef{v, i}
	} else {
		return v.Elem(i)
	}
}

//------------------------------------  Slice:  L[i:j]

func (v *VList) Slice(lval Value, x Value, y Value) Value {
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
	a := make([]Value, m)
	if v.rev {
		copy(a, v.data[n-j:n-i])
		ReverseValues(a)
	} else {
		copy(a, v.data[i:j])
	}
	return InitList(a)
}

//------------------------------------  ListPut: used in [: expr :]

type IListPut interface {
	ListPut(Value) Value // [: ...v... :]
}

func (v *VList) ListPut(x Value) Value {
	v.Grow(false, "[:put:]", x)
	return v
}

//------------------------------------  ListCat:  L1 ||| L2

type IListCat interface {
	ListCat(Value) Value // L ||| L
}

func (v *VList) ListCat(x Value) Value {
	return InitList(append(v.Export().([]Value), x.(*VList).Export().([]Value)...))
}

//  oset.go -- set operations

package runtime

import (
	"math/rand"
	"reflect"
)

//------------------------------------  Size:  *S

func (S *VSet) Size() Value {
	return NewNumber(float64(len(*S)))
}

//------------------------------------  Choose:  ?S

func (S *VSet) Choose(lval Value) Value {
	n := len(*S)
	if n == 0 {
		return nil // fail
	}
	vlist := reflect.ValueOf(*S).MapKeys()
	x := vlist[rand.Intn(n)].Interface()
	return Import(x) // convert back from GoKey
}

//------------------------------------  Take:  @S

func (S *VSet) Take(lval Value) Value {
	for v := range *S { // for just one
		delete(*S, v)
		return Import(v) // convert back from GoKey
	}
	return nil // must have been empty: fail
}

//------------------------------------  Dispense:  !S

func (S *VSet) Dispense(lval Value) (Value, *Closure) {
	vlist := reflect.ValueOf(*S).MapKeys()
	i := -1
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		i++
		if i >= len(vlist) {
			return nil, nil
		} else {
			return Import(vlist[i].Interface()), c
		}
	}}
	return c.Resume()
}

//------------------------------------  Send:  S @: x

func (S *VSet) Send(lval Value, x Value) Value {
	(*S)[GoKey(x)] = true
	return x
}

//------------------------------------  Index:  S[x]

func (S *VSet) Index(lval Value, x Value) Value {
	if (*S)[GoKey(x)] {
		return x // found x in set
	} else {
		return nil // fail
	}
}

//------------------------------------  Union: S1 ++ S2

type IUnion interface {
	Union(Value) Value // S ++ S
}

func (S1 *VSet) Union(x Value) Value {
	S2 := x.(*VSet)
	S3 := NewSet(EMPTYLIST)
	for k := range *S1 {
		(*S3)[k] = true
	}
	for k := range *S2 {
		(*S3)[k] = true
	}
	return S3
}

//------------------------------------  SetDiff: S1 -- S2

type ISetDiff interface {
	SetDiff(Value) Value // S -- S
}

func (S1 *VSet) SetDiff(x Value) Value {
	S2 := x.(*VSet)
	S3 := NewSet(EMPTYLIST)
	for k := range *S1 {
		if !(*S2)[k] {
			(*S3)[k] = true
		}
	}
	return S3
}

//------------------------------------  Intersect: S1 ** S2

type IIntersect interface {
	Intersect(Value) Value // S ** S
}

func (S1 *VSet) Intersect(x Value) Value {
	S2 := x.(*VSet)
	S3 := NewSet(EMPTYLIST)
	for k := range *S1 {
		if (*S2)[k] {
			(*S3)[k] = true
		}
	}
	return S3
}

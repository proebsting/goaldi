//  oset.go -- set operations

package goaldi

import (
	"math/rand"
	"reflect"
)

//------------------------------------  Size:  *S

func (S VSet) Size() Value {
	return NewNumber(float64(len(S)))
}

//------------------------------------  Choose:  ?S

func (S VSet) Choose(lval Value) Value {
	n := len(S)
	if n == 0 {
		return nil // fail
	}
	vlist := reflect.ValueOf(S).MapKeys()
	x := vlist[rand.Intn(n)].Interface()
	return Import(x) // convert back from GoKey
}

//------------------------------------  Take:  @S

func (S VSet) Take() Value {
	for v := range S { // for just one
		delete(S, v)
		return Import(v) // convert back from GoKey
	}
	return nil // must have been empty: fail
}

//------------------------------------  Dispense:  !S

func (S VSet) Dispense(lval Value) (Value, *Closure) {
	vlist := reflect.ValueOf(S).MapKeys()
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

//------------------------------------  Index:  S[x]

func (S VSet) Index(lval Value, x Value) Value {
	if S[GoKey(x)] {
		return x // found x in set
	} else {
		return nil // fail
	}
}

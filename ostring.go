//  ostring.go -- string operations

package goaldi

import (
	"math/rand"
)

//  sval -- extract VString value from arbitrary Value, or panic
func sval(v Value) *VString {
	if n, ok := v.(Stringable); ok {
		return n.ToString()
	} else {
		panic(&RunErr{"Not a string", v})
	}
}

//------------------------------------  Size:  *e

func (s *VString) Size() Value {
	return NewNumber(float64(s.length()))
}

//------------------------------------  Choose:  ?e

func (s *VString) Choose(lval IVariable) Value {
	n := s.length()
	if n == 0 {
		return nil // fail
	}
	i := rand.Intn(n)
	return s.slice(lval, i, i+1)
}

//------------------------------------  Dispense:  !e

func (s *VString) Dispense(lval IVariable) (Value, *Closure) {
	i := -1
	n := s.length()
	var f *Closure
	f = &Closure{func() (Value, *Closure) {
		i++
		if i < n {
			return s.slice(lval, i, i+1), f
		} else {
			return nil, nil
		}
	}}
	return f.Resume()
}

//------------------------------------  Concat:  e1 || e2

type IConcat interface {
	Concat(Value) Value
}

func (s *VNumber) Concat(t Value) Value {
	return s.ToString().Concat(t)
}

func (s *VString) Concat(x Value) Value {
	t := sval(x)
	return scat(s, 0, s.length(), t, 0, t.length(), EMPTY, 0, 0)
}

//------------------------------------  Index:  e1[e2]

func (s *VNumber) Index(lval IVariable, x Value) Value {
	return s.ToString().Index(lval, x)
}

func (s *VString) Index(lval IVariable, x Value) Value {
	i := int(x.(Numerable).ToNumber().Val())
	n := s.length()
	i = GoIndex(i, n)
	if i < n {
		return s.slice(lval, i, i+1) // return 1-char slice
	} else {
		return nil // subscript out of range
	}
}

//------------------------------------  Slice:  e1[e2:e3]

func (s *VNumber) Slice(lval IVariable, x Value, y Value) Value {
	return s.ToString().Slice(lval, x, y)
}

func (s *VString) Slice(lval IVariable, x Value, y Value) Value {
	i := int(x.(Numerable).ToNumber().Val())
	j := int(y.(Numerable).ToNumber().Val())
	n := s.length()
	i = GoIndex(i, n)
	j = GoIndex(j, n)
	if i > n || j > n {
		return nil // subscript out of range
	}
	if i > j {
		i, j = j, i // indexing was backwards
	}
	return s.slice(lval, i, j) // return slice
}

//------------------------------------  StrLT:  e1 << e2

type IStrLT interface {
	StrLT(Value) Value
}

func (s *VNumber) StrLT(x Value) Value {
	return s.ToString().StrLT(x)
}

func (s *VString) StrLT(x Value) Value {
	if s.compare(sval(x)) < 0 {
		return x
	} else {
		return nil
	}
}

//------------------------------------  StrLE:  e1 <<= e2

type IStrLE interface {
	StrLE(Value) Value
}

func (s *VNumber) StrLE(x Value) Value {
	return s.ToString().StrLE(x)
}

func (s *VString) StrLE(x Value) Value {
	if s.compare(sval(x)) <= 0 {
		return x
	} else {
		return nil
	}
}

//------------------------------------  StrEQ:  e1 == e2

type IStrEQ interface {
	StrEQ(Value) Value
}

func (s *VNumber) StrEQ(x Value) Value {
	return s.ToString().StrEQ(x)
}

func (s *VString) StrEQ(x Value) Value {
	t := sval(x)
	if s.length() != t.length() {
		return nil // can't match if lengths differ
	}
	if s.compare(t) == 0 {
		return x
	} else {
		return nil
	}
}

//------------------------------------  StrNE:  e1 ~== e2

type IStrNE interface {
	StrNE(Value) Value
}

func (s *VNumber) StrNE(x Value) Value {
	return s.ToString().StrNE(x)
}

func (s *VString) StrNE(x Value) Value {
	t := sval(x)
	if s.length() != t.length() {
		return x // can't be equal if lengths differ
	}
	if s.compare(t) != 0 {
		return x
	} else {
		return nil
	}
}

//------------------------------------  StrGE:  e1 >>= e2

type IStrGE interface {
	StrGE(Value) Value
}

func (s *VNumber) StrGE(x Value) Value {
	return s.ToString().StrGE(x)
}

func (s *VString) StrGE(x Value) Value {
	if s.compare(sval(x)) >= 0 {
		return x
	} else {
		return nil
	}
}

//------------------------------------  StrGT:  e1 >> e2

type IStrGT interface {
	StrGT(Value) Value
}

func (s *VNumber) StrGT(x Value) Value {
	return s.ToString().StrGT(x)
}

func (s *VString) StrGT(x Value) Value {
	if s.compare(sval(x)) > 0 {
		return x
	} else {
		return nil
	}
}

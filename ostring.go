//  ostring.go -- string operations

package goaldi

//  sval -- extract VString value from arbitrary Value, or panic
func sval(v Value) *VString {
	if n, ok := v.(Stringable); ok {
		return n.ToString()
	} else {
		panic(&RunErr{"Not a string", v})
	}
}

//------------------------------------  Size:  *e1

func (s *VNumber) Size() Value {
	return s.ToString().Size()
}

func (s *VString) Size() Value {
	return NewNumber(float64(s.length()))
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

func (s *VNumber) Index(x Value) (Value, *Closure) {
	return s.ToString().Index(x)
}

func (s *VString) Index(x Value) (Value, *Closure) {
	i := int(x.(Numerable).ToNumber().val())
	n := s.length()
	if i > 0 {
		i-- // convert to zero-based
	} else {
		i = n + i // count backwards from end
	}
	if i >= 0 || i > n {
		return Return(s.slice(i, i+1)) // return 1-char slice
	} else {
		return Fail() // subscript out of range
	}
}

//------------------------------------  LEqual:  e1 == e2

//#%#%#% incomplete. need num version, interface, 5more oprs, etc

func (s *VString) LEqual(x Value) Value {
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

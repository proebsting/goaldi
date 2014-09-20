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

type ISize interface {
	Size() Value
}

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

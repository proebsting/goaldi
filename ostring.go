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

//  VString.compare -- compare two strings, return <0, 0, or >0
func (s *VString) compare(t *VString) int {
	// check for easy case
	if s == t {
		return 0
	}
	// extract fields
	sl := s.low
	tl := t.low
	sh := s.high
	th := t.high
	sn := len(sl)
	tn := len(tl)
	// compare runes until one differs
	for i := 0; i < sn && i < tn; i++ {
		sr := rune(sl[i])
		tr := rune(tl[i])
		if sh != nil {
			sr |= rune(sh[i] << 8)
		}
		if th != nil {
			tr |= rune(th[i] << 8)
		}
		if sr != tr {
			return int(sr) - int(tr)
		}
	}
	// reached the end of one or both strings
	return sn - tn
}

//------------------------------------  Size:  *e1

type ISize interface {
	Size() Value
}

func (s *VNumber) Size() Value {
	return s.ToString().Size()
}

func (s *VString) Size() Value {
	return NewNumber(float64(len(s.low)))
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
	ns := len(s.low)
	nt := len(s.low) + len(t.low)
	low := make([]uint8, nt, nt)
	copy(low[:ns], s.low)
	copy(low[ns:], t.low)
	if s.high == nil && t.high == nil {
		return &VString{low, nil}
	}
	high := make([]uint16, nt, nt)
	if s.high != nil {
		copy(high[:ns], s.high)
	}
	if t.high != nil {
		copy(high[ns:], t.high)
	}
	return &VString{low, high}
}

//------------------------------------  LEqual:  e1 == e2

//#%#%#% incomplete. need num version, interface, 5more oprs, etc

func (s *VString) LEqual(x Value) Value {
	t := sval(x)
	if len(s.low) != len(t.low) {
		return nil // can't match if lengths differ
	}
	if s.compare(t) == 0 {
		return x
	} else {
		return nil
	}
}

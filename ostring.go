//  ostring.go -- string operations

//#%#%  not quite right yet.  not guaranteed that e2 is a VString.

package goaldi

//  e1 || e2

func (v1 *VNumber) Concat(v2 Value) Value {
	return v1.ToString().Concat(v2)
}

func (v1 *VString) Concat(v2 Value) Value {
	s1 := v1
	s2 := v2.(Stringable).ToString()
	return NewString(string(*s1) + string(*s2))
}

//  ostring.go -- string operations

//#%#%  not quite right yet.  not guaranteed that e2 is a VString.

package goaldi

//  e1 || e2

func (v1 *VNumber) Concat(v2 IString) (IString, *Closure) {
	return v1.ToString().Concat(v2)
}

func (v1 *VString) Concat(v2 IString) (IString, *Closure) {
	return NewString(string(*v1) + string(*(v2.(*VString)))), nil
}

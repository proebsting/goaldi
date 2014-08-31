package goaldi

//#%#%#% NEEDS REVIEW. PROBABLY BUGGY.
//#%#%#% NEEDS CONVERSION TO O-O FORM.

//  OToBy implements e1 to e2 by e3
func OToBy(e1, e2, e3 Value) (Value, *Closure) {
	n1 := e1.AsNumber()
	if n1 == nil {
		return Throwf("ToBy: e1 bad: %s", e1)
	}
	n2 := e2.AsNumber()
	if n2 == nil {
		return Throwf("ToBy: e2 bad: %s", e2)
	}
	n3 := NewNumber(1)
	if e3 != nil {
		n3 = e3.AsNumber()
		if n3 == nil {
			return Throwf("ToBy: e3 bad: %s", e3)
		}
		if *n3 == 0 {
			return Throwf("ToBy: by 0")
		}
	}
	v1 := *n1
	v2 := *n2
	v3 := *n3
	v1 -= v3
	var f *Closure
	f = &Closure{func() (Value, *Closure) {
		v1 += v3
		if (v3 > 0 && v1 <= v2) || (v3 < 0 && v1 >= v2) {
			return NewNumber(float64(v1)), f
		} else {
			return Fail()
		}
	}}
	return f.Resume()
}

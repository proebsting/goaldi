//  omisc.go -- miscellaneous runtime operations

package goaldi

//  ToBy -- implement "e1 to e2 by e3"
func ToBy(e1 Value, e2 Value, e3 Value) (Value, *Closure) {
	n1 := e1.(Numerable).ToNumber()
	if n1 == nil {
		panic(&RunErr{"ToBy: e1 bad", e1})
	}
	n2 := e2.(Numerable).ToNumber()
	if n2 == nil {
		panic(&RunErr{"ToBy: e2 bad", e2})
	}
	n3 := e3.(Numerable).ToNumber()
	if n3 == nil {
		panic(&RunErr{"ToBy: e3 bad", e3})
	}
	if *n3 == 0 {
		panic(&RunErr{"ToBy: by 0", nil})
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

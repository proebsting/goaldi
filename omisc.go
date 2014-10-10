//  omisc.go -- miscellaneous runtime operations

package goaldi

//  GoIndex(i, n) -- return Go index for i into length n, or n+1 if out of range
//  i is a 1-based index and may be nonpositive.  i==n or i==0 is in range.
//  The caller may want to check for result<n or result<=n depending on context.
func GoIndex(i int, n int) int {
	if i > 0 {
		i-- // convert to zero-based
	} else {
		i = n + i // count backwards from end
	}
	if i >= 0 && i <= n {
		return i // index is valid
	} else {
		return n + 1 // index is out of range
	}
}

//  VNumber.ICall -- implement i(e1, e2, e3...)
func (v *VNumber) Call(env *Env, args ...Value) (Value, *Closure) {
	i := GoIndex(int(v.Val()), len(args))
	if i < len(args) {
		return Return(args[i])
	} else {
		return Fail()
	}
}

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

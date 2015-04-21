//  ofile.go -- operators applied to files

package runtime

//  VFile.Dispense() implements the !f operator
func (f *VFile) Dispense(unused Value) (Value, *Closure) {
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		s := f.ReadLine()
		if s != nil {
			return s, c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

//  VFile.Take(lval) implements the @f operator
func (f *VFile) Take(lval Value) Value {
	s := f.ReadLine()
	if s != nil {
		return s
	} else {
		return nil
	}
}

//  VFile.Send(x) implements f @: x
func (f *VFile) Send(x Value) Value {
	Wrt(f, nil, nlByte, []Value{x})
	return x
}

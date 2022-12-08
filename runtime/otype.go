//  otype.go -- operations on system-defined Goaldi types

package runtime

// VType.Size() implements *t, returning 0 for a non-record type.
func (v *VType) Size() Value {
	return ZERO
}

// VType.Index(lval, x) fails immediately for a non-record type.
func (v *VType) Index(lval Value, x Value) Value {
	return nil
}

// VType.Dispense() fails immediately for a non-record type.
func (v *VType) Dispense(lval Value) (Value, *Closure) {
	return Fail()
}

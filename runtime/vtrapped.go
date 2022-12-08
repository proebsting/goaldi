//  vtrapped.go -- trapped variables and assignment functions and methods

package runtime

import (
	"fmt"
)

type VTrapped struct {
	Target *Value // pointer to target
}

const rTrapped = 3            // declare sort ranking
var _ IVariable = &VTrapped{} // validate implementation

// Trapped(v) -- create a simple (unindexed) trapped variable
func Trapped(target *Value) *VTrapped {
	return &VTrapped{target}
}

// NewVariable(x) returns a new trapped variable initialized to x.
func NewVariable(x Value) *VTrapped {
	return Trapped(&x)
}

// TrappedType is the instance of type type.
var TrappedType = NewType("trapped", "v", rTrapped, nil, nil,
	"trapped", "", "")

// VTrapped.String() -- show string representation: "V(<value>)"
func (t *VTrapped) String() string {
	return fmt.Sprintf("V(%v)", t.Deref())
}

// VTrapped.GoString() -- show string representation for traceback
// (This may be visible, for example, as an "offending value".)
func (t *VTrapped) GoString() string {
	return fmt.Sprintf("Variable(%#v)", t.Deref())
}

// VTrapped.Type() -- return the trapped type (shouldn't be used)
func (t *VTrapped) Type() IRank {
	return TrappedType
}

// VTrapped.Deref() -- extract value of for use as an rvalue
func (t *VTrapped) Deref() Value {
	// later becomes more complicated with tvsubs, tvstr
	return *t.Target
}

// VTrapped.Assign -- store value in target variable
func (t *VTrapped) Assign(v Value) IVariable {
	*t.Target = v
	return t
}

// Deref(v) -- dereference a value only if it implements IVariable
func Deref(v Value) Value {
	if d, ok := v.(IVariable); ok {
		return d.Deref()
	} else {
		return v
	}
}

// RevAssign -- implement reversible assignment (<-)
func RevAssign(e1 Value, e2 Value) (IVariable, *Closure) {
	v1 := e1.(IVariable)
	x1 := v1.Deref()
	v1.Assign(e2)
	return v1, &Closure{func() (Value, *Closure) {
		v1.Assign(x1)
		return Fail()
	}}
}

// Swap -- implement the exchange operation (:=:)
// Works for any two arguments implementing the IVariable interface
func Swap(e1 Value, e2 Value) IVariable {
	v1 := e1.(IVariable)
	v2 := e2.(IVariable)
	x1 := v1.Deref()
	x2 := v2.Deref()
	v1.Assign(x2)
	v2.Assign(x1)
	return v1
}

// RevSwap -- implement reversible exchange (<->)
func RevSwap(e1 Value, e2 Value) (IVariable, *Closure) {
	v1 := e1.(IVariable)
	v2 := e2.(IVariable)
	x1 := v1.Deref()
	x2 := v2.Deref()
	v1.Assign(x2)
	v2.Assign(x1)
	return v1, &Closure{func() (Value, *Closure) {
		v1.Assign(x1)
		v2.Assign(x2)
		return Fail()
	}}
}

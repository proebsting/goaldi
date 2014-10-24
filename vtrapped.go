//  vtrapped.go -- trapped variables and assignment functions and methods

package goaldi

import (
	"fmt"
)

type VTrapped struct {
	Target *Value // pointer to target
	// other stuff later for tvsubs, tvstr
}

//  Trapped(v) -- create a simple trapped (unindexed) variable
func Trapped(target *Value) *VTrapped {
	return &VTrapped{target}
}

//  VTrapped.Deref() -- extract value of for use as an rvalue
func (t *VTrapped) Deref() Value {
	// later becomes more complicated with tvsubs, tvstr
	return *t.Target
}

//  VTrapped.GoaldiValue -- Declare this to be a Goaldi value
func (*VTrapped) GoaldiValue() {}

//  VTrapped.String() -- show string representation: produces (&value)
//  #%#% should make this smarter
func (t *VTrapped) String() string {
	return fmt.Sprintf("(&%v)", t.Deref())
}

//  VTrapped.GoString() -- show string representation for traceback
func (t *VTrapped) GoString() string {
	return t.String()
}

//  VTrapped.Assign -- store value in target variable
func (t *VTrapped) Assign(v Value) IVariable {
	*t.Target = v
	return t
}

//  Deref(v) -- dereference a value only if it implements IVariable
func Deref(v Value) Value {
	if d, ok := v.(IVariable); ok {
		return d.Deref()
	} else {
		return v
	}
}

//  RevAssign -- implment reversible assignment (<-)
func RevAssign(e1 Value, e2 Value) (IVariable, *Closure) {
	v1 := e1.(IVariable)
	x1 := v1.Deref()
	v1.Assign(e2)
	return v1, &Closure{func() (Value, *Closure) {
		v1.Assign(x1)
		return Fail()
	}}
}

//  Swap -- implement the exchange operation (:=:)
//  Works for any two arguments implementing the IVariable interface
func Swap(e1 Value, e2 Value) IVariable {
	v1 := e1.(IVariable)
	v2 := e2.(IVariable)
	x1 := v1.Deref()
	x2 := v2.Deref()
	v1.Assign(x2)
	v2.Assign(x1)
	return v1
}

//  RevSwap -- implment reversible exchange (<->)
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

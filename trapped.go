//  trapped.go -- trapped variables

package goaldi

import (
	"fmt"
)

type Vtrapped struct {
	Target	*Value	// pointer to target
			// other stuff later for tvsubs, tvstr
}

var _ Value = &Vtrapped{}	// assert that Vtrapped implements Value

//  create a simple trapped (unindexed) variable
func Trapped(target *Value) *Vtrapped {
	return &Vtrapped{target}
}

//  extract value of a trapped variable for use as an rvalue
func (t *Vtrapped) Deref() Value {
	// later becomes more complicated with tvsubs, tvstr
	return *t.Target
}

//  show trapped variable as a string for debugging: produces [[value]]
func (t *Vtrapped) String() string {
	return fmt.Sprintf("[[%s]]", (*(t.Target)).String())
}

//  assign value to trapped variable
func (t *Vtrapped) Assign(v Value) *Vtrapped {
	*t.Target = v.Deref()
	return t
}

//  implement Value interface

func (t *Vtrapped) AsString() *VString { return t.Deref().AsString() }
func (t *Vtrapped) AsNumber() *VNumber { return t.Deref().AsNumber() }

func (t *Vtrapped) Add(v2 Value) (Value, *Closure)  { return t.Deref().Add(v2) }
func (t *Vtrapped) Mult(v2 Value) (Value, *Closure)  { return t.Deref().Mult(v2) }

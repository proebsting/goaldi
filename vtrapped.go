//  vtrapped.go -- VTrapped, for Goaldi trapped variables

package goaldi

import (
	"fmt"
)

type VTrapped struct {
	Target *Value // pointer to target
	// other stuff later for tvsubs, tvstr
}

var _ Value = &VTrapped{} // assert that VTrapped implements Value

//  create a simple trapped (unindexed) variable
func Trapped(target *Value) *VTrapped {
	return &VTrapped{target}
}

//  extract value of a trapped variable for use as an rvalue
func (t *VTrapped) Deref() Value {
	// later becomes more complicated with tvsubs, tvstr
	return *t.Target
}

//  show trapped variable as a string for debugging: produces [[value]]
func (t *VTrapped) String() string {
	return fmt.Sprintf("[[%s]]", (*(t.Target)).String())
}

//  assign value to trapped variable
func (t *VTrapped) Assign(v Value) *VTrapped {
	*t.Target = v.Deref()
	return t
}

//  implement Value interface

func (t *VTrapped) AsString() *VString { return t.Deref().AsString() }
func (t *VTrapped) AsNumber() *VNumber { return t.Deref().AsNumber() }

func (t *VTrapped) Add(v2 Value) (Value, *Closure)  { return t.Deref().Add(v2) }
func (t *VTrapped) Mult(v2 Value) (Value, *Closure) { return t.Deref().Mult(v2) }

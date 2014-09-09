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
func (t *VTrapped) Deref() (Value, *Closure) {
	// later becomes more complicated with tvsubs, tvstr
	return *t.Target, nil
}

//  show trapped variable as a string for debugging: produces [[value]]
//  #%#% should make this smarter
func (t *VTrapped) String() string {
	return fmt.Sprintf("[[%v]]", (*(t.Target)))
}

//  assign value to trapped variable
func (t *VTrapped) Assign(v Value) (IVariable, *Closure) {
	*t.Target = v
	return t, nil
}

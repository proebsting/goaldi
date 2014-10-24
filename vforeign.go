//  vforeign.go -- a wrapper for Go values in Goaldi
//
//  A foreign value references a Go value that could not be converted.

package goaldi

import (
	"fmt"
)

type VForeign struct {
	GoVal interface{}
}

//  NewForeign(value) -- construct new Foreign value
func NewForeign(v interface{}) *VForeign {
	return &VForeign{v}
}

//  VForeign.GoaldiValue -- Declare this to be a Goaldi value
func (*VForeign) GoaldiValue() {}

//  VForeign.String -- conversion to Go string returns "foreign(%v)"
func (v *VForeign) String() string {
	return fmt.Sprint("foreign(%v)", v.GoVal)
}

//  VForeign.GoString -- conversion to Go string returns "foreign(%#v)"
func (v *VForeign) GoString() string {
	return fmt.Sprint("foreign(%#v)", v.GoVal)
}

//  VForeign.Type returns "foreign"
func (v *VForeign) Type() Value {
	return type_foreign
}

var type_foreign = NewString("foreign")

//  VForeign.Export returns the underlying Go value
func (v *VForeign) Export() interface{} {
	return v.GoVal
}

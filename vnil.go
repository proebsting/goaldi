//  vnil.go -- VNil, the Goaldi type "nil"

package goaldi

import ()

//  The VNil strict contains no data
type VNil struct {
}

//  NIL is the one Goaldi nil value
var NIL = &VNil{}

//  NilVal is the same thing but of type Value
var NilVal Value = NIL

//  NewNil returns the singleton Goaldi nil value
func NewNil() *VNil {
	return NIL
}

//  VNil.String returns "nil" as a Go string
func (v *VNil) String() string {
	return "nil"
}

//  VNil.Type returns "nil"
func (v *VNil) Type() Value {
	return type_nil
}

//  VNil.Export returns a nil
func (v *VNil) Export() interface{} {
	return nil
}

var type_nil = NewString("nil")

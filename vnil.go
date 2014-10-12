//  vnil.go -- VNil, the Goaldi type "nil"

package goaldi

import ()

//  The VNil strict contains no data
type VNil struct {
}

//  NIL is the one Goaldi nil value
var NIL = &VNil{}

//  NilValue is the same thing but of type Value
var NilValue Value = NIL

//  NewNil returns the singleton Goaldi nil value
func NewNil() *VNil {
	return NIL
}

//  VNil.String -- default conversion to Go string returns "~"
func (v *VNil) String() string {
	return "~" //#%#% ??
}

//  VNil.GoString -- convert to string "%nil" for image() and printf("%#v")
func (v *VNil) GoString() string {
	return "%nil" //#%#% ??
}

//  VNil.Type returns "nil"
func (v *VNil) Type() Value {
	return type_nil
}

var type_nil = NewString("nil")

//  VNil.Export returns a Go nil
func (v *VNil) Export() interface{} {
	return nil
}

//  nil.go -- the Goaldi type "nil"

package goaldi

import ()

//  The VNil strict contains no data
type VNil struct {
	Stubs
}

//  NIL is the one Goaldi nil value
var NIL = &VNil{}

//  NilVal is the same thing but of type Value
var NilVal Value = NIL

//  NewNil returns the singleton Goaldi nil value
func NewNil() *VNil {
	return NIL
}

// //  VNil.AsString returns nil, signifying no implicit conversion to VString
// func (v *VNil) AsString() *VString {
// 	return nil
// }
//
// //  VNil.AsString returns nil, signifying no implicit conversion to VNumber
// func (v *VNil) AsNumber() *VNumber {
// 	return nil
// }

//  VNil.String returns "%nil%" as a Go string for printing
func (v *VNil) String() string {
	return "%nil%"
}

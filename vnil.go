//  vnil.go -- vnil, the Goaldi type "nil"

package goaldi

import ()

//  The vnil strict contains no data and is not exported.
type vnil struct {
}

//  NilValue is the one and only nil value.
//  For convenience, its type is Value, not vnil.
var NilValue Value = &vnil{}

//  vnil.String -- default conversion to Go string returns "~"
func (v *vnil) String() string {
	return "~" //#%#% ??
}

//  vnil.GoString -- convert to string "%nil" for image() and printf("%#v")
func (v *vnil) GoString() string {
	return "%nil" //#%#% ??
}

//  vnil.Type returns "nil"
func (v *vnil) Type() Value {
	return type_nil
}

var type_nil = NewString("nil")

//  vnil.Export returns a Go nil
func (v *vnil) Export() interface{} {
	return nil
}

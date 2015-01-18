//  vnil.go -- vnil, the Goaldi type "nil"

package goaldi

import ()

//  The constructor named "nil" is not a global because "nil" is reserved.
var NilType = NewType(rNil, Nil, "nil", "", "return nil value")

//  The vnil struct contains no data and is not exported.
type vnil struct {
}

//  NilValue is the one and only nil value.
//  For convenience, its type is Value, not vnil.
var NilValue Value = &vnil{}

//  The constructor just returns the only nil value.
func Nil(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("nil", args)
	return Return(NilValue)
}

//  vnil.String -- default conversion to Go string returns "~"
func (v *vnil) String() string {
	return "~"
}

//  vnil.GoString -- convert to string "nil" for image() and printf("%#v")
func (v *vnil) GoString() string {
	return "nil"
}

//  vnil.Type returns the nil type
func (v *vnil) Type() IRank {
	return NilType
}

//  vnil.Copy returns itself
func (v *vnil) Copy() Value {
	return v
}

//  vnil.Import returns itself
func (v *vnil) Import() Value {
	return v
}

//  vnil.Export returns a Go nil
func (v *vnil) Export() interface{} {
	return nil
}

//  Declare methods
var NilMethods = MethodTable([]*VProcedure{
	DefMeth((*vnil).Type, "type", "", "return nil type"),
	DefMeth((*vnil).Copy, "copy", "", "return nil value"),
	DefMeth((*vnil).String, "string", "", "return \"~\""),
	DefMeth((*vnil).GoString, "image", "", "return \"nil\""),
})

//  vnil.Field implements methods
func (v *vnil) Field(f string) Value {
	return GetMethod(NilMethods, v, f)
}

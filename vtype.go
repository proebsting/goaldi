//  vtype.go -- VType, the Goaldi type "type"

package goaldi

import (
	"fmt"
)

var _ = fmt.Printf // enable debugging

//  ranking of types for sorting
const (
	rNil = iota
	rType
	rNumber
	rString
	rFile
	rChannel
	rDefn
	rMethVal
	rProc
	rList
	rTable
	rRecord
	rExternal
)

//  The global named "type"
var TypeType = NewType(Type, "type", "x", "return type of value")

//  A type value structure
type VType struct {
	Name string      // type name
	Ctor *VProcedure // standard constructor procedure
}

//  NewType defines and registers a Goaldi standard (not a record) type.
//  The constructor procedure is installed in the standard library
//  (but remains inaccessible for reserved names "nil" and "procedure").
func NewType(ctor Procedure,
	name string, pspec string, descr string) *VType {
	proc := DefProc(ctor, name, pspec, descr)
	t := &VType{name, proc}
	StdLib[name] = t
	return t
}

//  Declare methods on a type value
var TypeMethods = MethodTable([]*VProcedure{
	DefMeth((*VType).Type, "type", "", "return type type"),
	DefMeth((*VType).Copy, "copy", "", "return type value"),
	DefMeth((*VType).String, "string", "", "return type name"),
	DefMeth((*VType).GoString, "image", "", "return type image"),
})

//  VType.Field implements methods
func (v *VType) Field(f string) Value {
	return GetMethod(TypeMethods, v, f)
}

//  VType.String -- default conversion to Go string returns type name
func (t VType) String() string {
	return "t:" + t.Name
}

//  VType.GoString -- convert to Go string for image() and printf("%#v")
func (t VType) GoString() string {
	return "type " + t.Name
}

//  VType.Rank -- return rType
func (t VType) Rank() int {
	return rType
}

//  VType.Type -- return the type "type"
func (t VType) Type() Value {
	return TypeType
}

//  VType.Copy returns itself
func (t VType) Copy() Value {
	return t
}

//  VType.Import returns itself
func (t VType) Import() Value {
	return t
}

//  VType.Export returns itself.
func (t VType) Export() interface{} {
	return t
}

//  VType.Call invokes the constructor procedure for a type.
func (t *VType) Call(env *Env, args []Value, names []string) (Value, *Closure) {
	return t.Ctor.Call(env, args, names)
}

//  Type(v) -- construct (or sometimes just find) an instance of type v
func Type(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("type", args)
	v := ProcArg(args, 0, NilValue)
	if t, ok := v.(IType); ok {
		return Return(t.Type())
	} else {
		return Return(ExternalType)
	}
}

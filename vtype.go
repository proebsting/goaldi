//  vtype.go -- VType, the Goaldi type "type"

package goaldi

import (
	"fmt"
)

var _ = fmt.Printf // enable debugging

//  ranking of types for sorting
const (
	rNil = iota
	rTrapped
	rType
	rNumber
	rString
	rFile
	rChannel
	rMethVal
	rProc
	rList
	rTable
	rRecord
	rExternal
)

//  The global named "type"
var TypeType = NewType("t", rType, Type, TypeMethods,
	"type", "x", "return type of value")

//  A type value structure
type VType struct {
	TypeName string                 // type name
	Abbr     string                 // one-character abbreviation
	SortRank int                    // rank for sorting
	Ctor     *VProcedure            // standard constructor procedure
	Methods  map[string]*VProcedure // method table
}

//  NewType defines and registers a Goaldi standard (not a record) type.
//  The constructor procedure is installed in the standard library
//  (but remains inaccessible for reserved names "nil" and "procedure").
//  A nil constructor indicates an internal type (i.e. Trapped),
//  and such a type is not installed in the library.
func NewType(abbr string, rank int, ctor Procedure,
	mtable map[string]*VProcedure,
	name string, pspec string, descr string) *VType {
	proc := DefProc(ctor, name, pspec, descr)
	t := &VType{name, abbr, rank, proc, mtable}
	if ctor != nil {
		StdLib[name] = t
	}
	return t
}

//  Declare methods
var TypeMethods = MethodTable([]*VProcedure{
	DefMeth((*VType).Name, "name", "", "get type name"),
	DefMeth((*VType).Char, "char", "", "get abbreviation character"),
})

//  VType.String -- default conversion to Go string returns type name
func (t *VType) String() string {
	return "t:" + t.TypeName
}

//  VType.GoString -- convert to Go string for image() and printf("%#v")
func (t *VType) GoString() string {
	return "type " + t.TypeName
}

//  VType.Type -- return the type "type"
func (t *VType) Type() IRank {
	return TypeType
}

//  VType.Copy returns itself
func (t *VType) Copy() Value {
	return t
}

//  VType.Before compares itself with a constructor or type value
func (a *VType) Before(b Value, i int) bool {
	switch t := b.(type) {
	case *VType:
		return a.SortRank < t.SortRank
	case *VCtor:
		return a.SortRank < rRecord
	default:
		panic(Malfunction("unexpected type in VType.Before"))
	}
}

//  VType.Import returns itself
func (t *VType) Import() Value {
	return t
}

//  VType.Export returns itself.
func (t *VType) Export() interface{} {
	return t
}

//  VType.Rank returns the sorting rank.
func (t *VType) Rank() int {
	return t.SortRank
}

//  VType.Call invokes the constructor procedure for a type.
func (t *VType) Call(env *Env, args []Value, names []string) (Value, *Closure) {
	return t.Ctor.Call(env, args, names)
}

//  VType.Name returns the type name
func (t *VType) Name(args ...Value) (Value, *Closure) {
	return Return(NewString(t.TypeName))
}

//  VType.Char returns the abbreviation character.
func (t *VType) Char(args ...Value) (Value, *Closure) {
	return Return(NewString(t.Abbr))
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

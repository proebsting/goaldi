//  interfaces.go -- interfaces for implementing types and operations

package goaldi

import (
	"fmt"
)

//  Any Go value can be a Goaldi value.
//  This identifier is intended to designate a Goaldi value context.
type Value interface{}

//  IExternal -- declares an external type to be a Goaldi external
//  (to prevent conversion of something that otherwise might be converted.)
type IExternal interface {
	GoaldiExternal()
}

//  ICore -- interfaces required of all Goaldi types
type ICore interface {
	fmt.Stringer   // for string() and printf("%v") (v.String())
	fmt.GoStringer // for image() and printf("%#v") (v.GoString())
	IType          // for "type()", and for ranking when sorting
	ICopy          // for "copy()"
	IImport        // for returning self to Import()
	IExport        // for passing to a Go function as interface{} value
	// optional:  Numerable and Stringable, if implicitly convertible
	// optional:  IField, for implementing methods
	// optional:  IIdentical, if === requires more than pointer comparison
	Before(Value, int) bool // compare value of same type on field i
}

//  Confirm implementation of core interfaces by all types
var _ ICore = NilValue.(*vnil)
var _ ICore = &VType{}
var _ ICore = NewNumber(1)
var _ ICore = NewString("a")
var _ ICore = &VFile{}
var _ ICore = NewChannel(0)
var _ ICore = &VCtor{}
var _ ICore = &VMethVal{}
var _ ICore = &VProcedure{}
var _ ICore = &VList{}
var _ ICore = &VTable{}
var _ ICore = &VRecord{}

//  IRank designates anything usable as a type: VType or VCtor
type IRank interface {
	Rank() int                            // return rank for sorting
	Name(args ...Value) (Value, *Closure) // return type name to Goaldi
	Char(args ...Value) (Value, *Closure) // return type char to Goaldi
}

type IType interface {
	Type() IRank // return type for type()
}

type ICopy interface {
	Copy() Value // return copy of value
}

//  IImport -- convert to Goaldi value
type IImport interface {
	Import() Value // return value to be imported as Goaldi value
}

//  IExport -- convert Goaldi value to Go value
type IExport interface {
	Export() interface{} // return value for export to Go function
}

//  Interfaces for implicit conversion (also requires operator methods)
type Stringable interface {
	ToString() *VString // if implicitly convertible to string
}
type Numerable interface {
	ToNumber() *VNumber // if implicitly convertible to number
}

//  IVariable -- an assignable trapped variable (simple or subscripted)
type IVariable interface {
	Deref() Value           // return dereferenced value
	Assign(Value) IVariable // assign value
}

var _ IVariable = &VTrapped{} // confirm implementation by VTrapped

//  Interfaces for operations that can produce substring variables
//  when applied to strings.
//  If the lval argument is nil, an rvalue is wanted.
//  If not, it is just an lvalue flag for most operations, but for
//  substring assignment it is the actual underlying string to replace.

type IChoose interface { // ?x
	Choose(lval Value) Value
}

type IDispense interface { // !x
	Dispense(lval Value) (Value, *Closure)
}

type IIndex interface { // x[y]
	Index(lval Value, y Value) Value
}

type ISlice interface { // x[i:j]
	Slice(lval Value, i Value, j Value) Value
}

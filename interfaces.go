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
	fmt.Stringer   // for printing (v.String())
	fmt.GoStringer // for image() and printf("%#v") (v.GoString())
	IType          // for "type()"
	ICopy          // for "copy()"
	IImport        // for returning self to Import()
	IExport        // for passing to a Go function as interface{} value
	IField         // for implementing methods
	// optional:  Numerable and Stringable, if implicitly convertible
	// optional:  IIdentical, if === requires more than pointer comparison
}

//  Confirm implementation of core interfaces by all types
var _ ICore = NilValue.(*vnil)
var _ ICore = NewNumber(1)
var _ ICore = NewString("a")
var _ ICore = &VFile{}
var _ ICore = &VProcedure{}
var _ ICore = &VDefn{}
var _ ICore = &VStruct{}
var _ ICore = &VList{}
var _ ICore = &VMap{}

type IType interface {
	Type() Value // return name of type for type()
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
//  If the IVariable argument is nil, a value is wanted.
//  If not, it is just a flag for most datatypes but is the actual
//  underlying value to be replaced by substring assignment.

type IChoose interface {
	Choose(IVariable) Value
}

type IDispense interface {
	Dispense(IVariable) (Value, *Closure)
}

type IIndex interface {
	Index(IVariable, Value) Value
}

type ISlice interface {
	Slice(IVariable, Value, Value) Value
}

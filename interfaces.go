//  interfaces.go -- interfaces for implementing types and operations

package goaldi

import (
	"fmt"
)

//  Any Go value can be a Goaldi value.
//  Use of this interface is intended to designate a Goaldi value context.
type Value interface{}

//  IExternal -- declares an external type to be a Goaldi external
//  (i.e. tells Goaldi to keeps hands off even it it looks convertible)
type IExternal interface {
	ExternalType() string // return type name for external value
}

//  IImport -- for an external type that declares its own converter method
type IImport interface {
	Import() Value
}

//  ICore -- interfaces required of all Goaldi types
type ICore interface {
	fmt.Stringer   // for printing (v.String())
	fmt.GoStringer // for image() and printf("%#v") (v.GoString())
	IType          // for "Type()"
	IExport        // for passing to a Go function as interface{} value
	// optional:  Numerable and Stringable, if implicitly convertible
	// optional:  IIdentical, if === requires more than pointer comparison
}

var _ ICore = NewNil()       // confirm implementation by VNil
var _ ICore = NewNumber(1)   // confirm implementation by VNumber
var _ ICore = NewString("a") // confirm implementation by VString
var _ ICore = &VProcedure{}  // confirm implementation by VProcedure

type IType interface {
	Type() Value // return name of type for type()
}

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

//  IIdentical -- needed for types where === is not just a pointer match
type IIdentical interface {
	Identical(Value) Value
}

var _ IIdentical = NewNumber(1)   // confirm implementation by VNumber
var _ IIdentical = NewString("a") // confirm implementation by VString

//  IVariable -- an assignable trapped variable (simple or subscripted)
type IVariable interface {
	Deref() Value           // return dereferenced value
	Assign(Value) IVariable // assign value
}

var _ IVariable = &VTrapped{} // confirm implementation by VTrapped

//  Interfaces for indexing operations that can produce variables
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

//  Other interfaces implemented by multiple types

type ISize interface {
	Size() Value
}

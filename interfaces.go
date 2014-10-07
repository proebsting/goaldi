//  interfaces.go -- interfaces for implementing types and operations

package goaldi

import (
	"fmt"
)

//  A Value can now be anything
//  Use of this interface is intended to designate a Goaldi value
type Value interface{}

//  Interfaces for implicitly convertible values
type Stringable interface {
	ToString() *VString
}
type Numerable interface {
	ToNumber() *VNumber
}

//  ICore -- should be implemented by all Goaldi types
type ICore interface {
	fmt.Stringer // i.e. v.String() -> string
	IType
	IExport
}

var _ ICore = NewNil()       // confirm implementation by VNil
var _ ICore = NewNumber(1)   // confirm implementation by VNumber
var _ ICore = NewString("a") // confirm implementation by VString
var _ ICore = &VProcedure{}  // confirm implementation by VProcedure

//  IVariable -- assignable trapped variable
type IVariable interface {
	Deref() Value           // return dereferenced value
	Assign(Value) IVariable // assign value
}

var _ IVariable = &VTrapped{} // confirm implementation by VTrapped

//  IExternal -- declares an external type to be a Goaldi external
//  (i.e. tells Goaldi to keeps hands off even it it looks convertible)
type IExternal interface {
	ExternalType() string // return type name for external value
}

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

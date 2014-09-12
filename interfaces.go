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

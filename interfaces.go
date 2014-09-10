//  interfaces.go -- interfaces for implementing types and operations

package goaldi

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
	String() string // return (Go) string image for printing
	// Type()?
	// Copy()?
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

//  interfaces for math operations

type INegate interface { // -e
	Negate() Value
}

var _ INegate = NewNumber(1)
var _ INegate = NewString("1")

type IAdd interface { // e1 + e2
	Add(Value) Value
}

var _ IAdd = NewNumber(1)
var _ IAdd = NewString("1")

type IMult interface { // e1 * e2
	Mult(Value) Value
}

var _ IMult = NewNumber(1)
var _ IMult = NewString("1")

//  interfaces for string operations
type IConcat interface { // e1 || e2
	Concat(Value) Value
}

var _ IConcat = NewNumber(1)
var _ IConcat = NewString("a")

//  IExternal -- declares an external type to be a Goaldi external
//  (i.e. tells Goaldi to keeps hands off even it it looks convertible)
type IExternal interface {
	ExternalType() string // return type name for external value
}

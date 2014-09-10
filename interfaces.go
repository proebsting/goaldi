//  interfaces.go -- interfaces for implementing types and operations

//#%#% Add INumerbable, IStringable for 2nd operand of math/string oprns?

package goaldi

//  A Value can now be anything
type Value interface{}

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

//  IMath -- arithmetic operations
type IMath interface {
	Negate() IMath    // -e
	Add(IMath) IMath  // e1 + e2
	Mult(IMath) IMath // e1 * e2
}

var _ IMath = NewNumber(1)   // confirm implementation by VNumber
var _ IMath = NewString("a") // confirm implementation by VString

//  IString -- string operations
type IString interface {
	Concat(IString) IString // e1 || e2
}

var _ IString = NewNumber(1)   // confirm implementation by VNumber
var _ IString = NewString("a") // confirm implementation by VString

//  IExternal -- declares an external type to be a Goaldi external
//  (i.e. tells Goaldi to keeps hands off even it it looks convertible)
type IExternal interface {
	ExternalType() string // return type name for external value
}

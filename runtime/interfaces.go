//  interfaces.go -- interfaces for implementing types and operations

package runtime

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

//	------------------------- Core Interfaces -------------------------
//  These interfaces must be implemented by all Goaldi types.

type ICore interface {
	fmt.Stringer   // for string() and printf("%v") (v.String())
	fmt.GoStringer // for image() and printf("%#v") (v.GoString())
	IType          // for "type()", and for ranking when sorting
	ICopy          // for "copy()"
	IImport        // for returning self to Import()
	IExport        // for passing to a Go function as interface{} value
	// add: IIdentical, if === requires other than pointer comparison
	// add: Numerable and Stringable, if implicitly convertible
	// add: IField, if any methods are defined on type
	Before(Value, int) bool // compare value of same type on field i
}

type IType interface {
	Type() IRank // return type value (that implements IRank) for type()
}

type ICopy interface {
	Copy() Value // return copy of value
}

type IImport interface {
	Import() Value // convert to Goaldi value (a no-op for built-in types)
}

type IExport interface {
	Export() interface{} // convert for export to Go function
}

//	------------------------- Operator Interfaces -------------------------
//  These interfaces are associated with particular operators.
//  Interfaces limited to a single type are declared within that type.

type ISize interface {
	Size() Value // *x
}

type ITake interface {
	Take() Value // @x
}

type IField interface {
	Field(string) Value // x.id
}

//  IIdentical -- for types where === is not just a pointer comparison
type IIdentical interface {
	Identical(Value) Value // === and ~===
}

//  Interfaces for operations that can produce lvalues, and can
//  produce substring variables when applied to strings.
//  If the lval argument is nil, an rvalue is wanted.
//  If not, it is just an lvalue flag for most operations, but for
//  substring assignment it is the actual underlying string to replace.

type IChoose interface {
	Choose(lval Value) Value // ?x
}

type IDispense interface {
	Dispense(lval Value) (Value, *Closure) // !x
}

type IIndex interface {
	Index(lval Value, y Value) Value // x[y]
}

type ISlice interface {
	Slice(lval Value, i Value, j Value) Value // x[i:j]
}

//  Interfaces for implicit conversion (these also require operator methods)

type Stringable interface {
	ToString() *VString // if implicitly convertible to string
}

type Numerable interface {
	ToNumber() *VNumber // if implicitly convertible to number
}

//	------------------------- Miscellaneous Interfaces -------------------------

//  IRank designates anything usable as a type: VType or VCtor
type IRank interface {
	Rank() int                            // return rank for sorting
	Name(args ...Value) (Value, *Closure) // return type name to Goaldi
	Char(args ...Value) (Value, *Closure) // return type char to Goaldi
}

//  IVariable -- an assignable trapped variable (simple or subscripted)
type IVariable interface {
	Deref() Value           // return dereferenced value
	Assign(Value) IVariable // assign value
}

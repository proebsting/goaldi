//  frecord.go -- library routines for record types

package runtime

import (
	"bytes"
	"fmt"
)

var _ = fmt.Printf // enable debugging

//  Declare library procedures
func init() {
	GoLib(Tuple, "tuple", "id:e...", "create anonymous record")
	StdLib["tuple"].(*VProcedure).RawCall = true // add magic bit
}

//  tuple(id:e, ...) creates an anonymous record value.
//  Each argument must be named.
//  Each distinct identifier list defines a new type,
//  all of which have the name "tuple".
func Tuple(env *Env, args []Value, names []string) (Value, *Closure) {
	//  Note the special RawCall argument list (and special registration above).
	defer Traceback("tuple", args)
	if len(names) < len(args) {
		panic(NewExn("Unnamed tuple arguments not allowed"))
	}
	t := TupleType(names)
	return Return(t.New(args))
}

//  Table of known tuples, indexed by stringified list of fields
var KnownTuples = make(map[string]*VCtor)

//  TupleType(names) finds or makes a type for constructing a tuple
func TupleType(names []string) *VCtor {
	// make a string of the field names e.g. "a,b,c,"
	var b bytes.Buffer
	for _, s := range names {
		b.WriteString(s)
		b.WriteByte(',')
	}
	s := b.String()
	// check for already known type
	t := KnownTuples[s]
	if t == nil {
		t = NewCtor("tuple", nil, names)
		KnownTuples[s] = t
	}
	return t
}

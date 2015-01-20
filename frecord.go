//  frecord.go -- library routines for record types

package goaldi

import (
	"bytes"
	"fmt"
)

var _ = fmt.Printf // enable debugging

//  Declare standard methods
var RecordMethods = MethodTable([]*VProcedure{
	DefMeth((*VRecord).Type, "type", "", "return type of record"),
	DefMeth((*VRecord).Copy, "copy", "", "duplicate record"),
	DefMeth((*VRecord).String, "string", "", "return short string"),
	DefMeth((*VRecord).GoString, "image", "", "return string image"),
})

//  Declare library procedures
func init() {
	GoLib(Tuple, "tuple", "id:e...", "create anonymous record")
	StdLib["tuple"].(*VProcedure).RawCall = true // add magic bit
}

//  Tuple(id:v, ...) creates an anonymous record dynamically
//  Note the special argument list here (and special registration above).
func Tuple(env *Env, args []Value, names []string) (Value, *Closure) {
	defer Traceback("tuple", args)
	if len(names) < len(args) {
		panic(NewExn("unnamed tuple arguments not allowed"))
	}
	t := TupleType(names)
	return Return(t.New(args))
}

//  Table of known tuples, indexed by stringified list of fields
var KnownTuples = make(map[string]*VDefn)

//  TupleType(names) finds or makes a type for constructing a tuple
func TupleType(names []string) *VDefn {
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
		t = NewDefn("tuple", names)
		KnownTuples[s] = t
	}
	return t
}

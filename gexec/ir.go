//  ir.go -- intermediate representation structures

package main

import (
	"reflect"
)

//  list of IR struct examples for use by JSON input converter
var irlist = [...]interface{}{
	&ir_Record{},
	&ir_Global{},
	&ir_Initial{},
	&ir_Function{},
	&ir_chunk{},
	&ir_NoOp{}, // not normally seen, but allowed as a comment
	&ir_Catch{},
	&ir_EnterScope{},
	&ir_ExitScope{},
	&ir_Var{},
	&ir_Key{},
	&ir_NilLit{},
	&ir_IntLit{},
	&ir_RealLit{},
	&ir_StrLit{},
	&ir_MakeClosure{},
	&ir_Move{},
	&ir_MoveLabel{},
	&ir_MakeList{},
	&ir_Field{},
	&ir_OpFunction{},
	&ir_Call{},
	&ir_ResumeValue{},
	&ir_Goto{},
	&ir_IndirectGoto{},
	&ir_Succeed{},
	&ir_Fail{},
	&ir_Create{},
	&ir_CoRet{},
	&ir_CoFail{},
	&ir_Select{},
	&ir_SelectCase{},
	&ir_NoValue{},     // seen only if unoptimized; not implemented
	&ir_Unreachable{}, // seen only if unoptimized; not implemented
}

//  struct table indexed by type names
var irtable = make(map[string]reflect.Type)

func init() {
	for _, ir := range irlist {
		t := reflect.TypeOf(ir).Elem()
		irtable[t.Name()] = t
	}
}

//  intermediate representation struct definitions
//  all fields must be capitalized for access by the reflection package

type ir_Record struct {
	Coord      string
	Name       string
	Extends    string
	Extendspkg string
	FieldList  []string
}

type ir_Initial struct {
	Coord string
	Fn    string
}

type ir_Global struct {
	Coord    string
	NameList []string
	Fn       string
}

type ir_Function struct {
	Coord       string
	Name        string
	ParamList   []string
	Accumulate  string // may be nil
	LocalList   []string
	StaticList  []string
	UnboundList []string
	CodeList    []ir_chunk
	CodeStart   string
	Parent      string
}

type ir_chunk struct {
	Label    string
	InsnList []interface{} // heterogeneous
}

type ir_NoOp struct {
	Coord   string
	Comment string
}

type ir_Catch struct {
	Coord string
	Fn    string
}

type ir_EnterScope struct {
	Coord       string
	NameList    []string
	DynamicList []string
	Scope       string
}

type ir_ExitScope struct {
	Coord    string
	NameList []string
	Scope    string
}

type ir_Var struct {
	Coord string
	Lhs   string
	Name  string
	Scope string
}

type ir_Key struct {
	Coord string
	Lhs   string // may be nil
	Name  string
	Scope string
}

type ir_NilLit struct {
	Coord string
	Lhs   string
}

type ir_IntLit struct {
	Coord string
	Lhs   string
	Val   string
}

type ir_RealLit struct {
	Coord string
	Lhs   string
	Val   string
}

type ir_StrLit struct {
	Coord string
	Lhs   string
	Len   string // length of the UTF-8 encoding
	Val   string // individual bytes of the UTF-8 encoding
}

type ir_MakeClosure struct {
	Coord string
	Lhs   string
	Name  string
}

type ir_Move struct {
	Coord string
	Lhs   string
	Rhs   string
}

type ir_MoveLabel struct {
	Coord string
	Lhs   string
	Label string
}

type ir_MakeList struct {
	Coord     string
	Lhs       string
	ValueList []interface{} // heterogeneous
}

type ir_Field struct {
	Coord string
	Lhs   string // may be nil
	Expr  string
	Field string
	Rval  string // may be nil
}

type ir_OpFunction struct {
	Coord      string
	Lhs        string // may be nil
	Lhsclosure string // may be nil
	Fn         string
	ArgList    []interface{} // heterogeneous
	Rval       string        // may be nil
	FailLabel  string        // may be nil
}

type ir_Call struct {
	Coord      string
	Lhs        string
	Lhsclosure string
	Fn         string
	ArgList    []interface{} // heterogeneous
	NameList   []string
	FailLabel  string // may be nil
	Scope      string
}

type ir_ResumeValue struct {
	Coord      string
	Lhs        string // may be nil
	Lhsclosure string
	Closure    string
	FailLabel  string // may be nil
}

type ir_Goto struct {
	Coord       string
	TargetLabel string
}

type ir_IndirectGoto struct {
	Coord          string
	TargetTmpLabel string
	LabelList      []string
}

type ir_Succeed struct {
	Coord       string
	Expr        string
	ResumeLabel string // may be nil
}

type ir_Fail struct {
	Coord string
}

type ir_Create struct {
	Coord      string
	Lhs        string
	CoexpLabel string
}

type ir_CoRet struct {
	Coord       string
	Value       string
	ResumeLabel string
}

type ir_CoFail struct {
	Coord string
}

type ir_Select struct {
	Coord     string
	CaseList  []ir_SelectCase
	FailLabel string
}

type ir_SelectCase struct {
	Coord     string
	Kind      string // "send" | "receive" | "default"
	Lhs       string
	Rhs       string
	BodyLabel string
}

type ir_NoValue struct {
	Coord string
	Lhs   string
}

type ir_Unreachable struct {
	Coord string
}

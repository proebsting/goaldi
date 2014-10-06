//  ir.go -- intermediate representation structures

package main

import (
	"fmt"
	"reflect"
)

//  list of IR struct examples for use by JSON input converter
var irlist = [...]interface{}{
	&ir_Invocable{},
	&ir_Link{},
	&ir_Record{},
	&ir_Global{},
	&ir_Function{},
	&ir_chunk{},
	&ir_Var{},
	&ir_Key{},
	&ir_IntLit{},
	&ir_RealLit{},
	&ir_StrLit{},
	&ir_CsetLit{},
	&ir_operator{},
	&ir_MakeClosure{},
	&ir_Move{},
	&ir_MoveLabel{},
	&ir_Deref{},
	&ir_Assign{},
	&ir_MakeList{},
	&ir_Field{},
	&ir_OpFunction{},
	&ir_Call{},
	&ir_ResumeValue{},
	&ir_EnterInit{},
	&ir_Goto{},
	&ir_IndirectGoto{},
	&ir_Succeed{},
	&ir_Fail{},
	&ir_Create{},
	&ir_CoRet{},
	&ir_CoFail{},
	&ir_ScanSwap{},
	&ir_Unreachable{},
	&ir_coordinate{},
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

type ir_Invocable struct {
	Coord    *ir_coordinate
	NameList []string
	All      string
}

type ir_Link struct {
	Coord    *ir_coordinate
	NameList []string
}

type ir_Record struct {
	Coord     *ir_coordinate
	Name      string
	FieldList []string
}

type ir_Global struct {
	Coord    *ir_coordinate
	NameList []string
}

type ir_Function struct {
	Coord      *ir_coordinate
	Name       string
	ParamList  []string
	Accumulate string // may be nil
	LocalList  []string
	StaticList []string
	CodeList   []ir_chunk
	CodeStart  string
	Lvalset    []string
}

type ir_chunk struct {
	Label    string
	InsnList []interface{} // heterogeneous
}

type ir_Var struct {
	Coord *ir_coordinate
	Lhs   string
	Name  string
}

type ir_Key struct {
	Coord     *ir_coordinate
	Lhs       string // may be nil
	Name      string
	FailLabel string // may be nil
}

type ir_IntLit struct {
	Coord *ir_coordinate
	Lhs   string
	Val   string
}

type ir_RealLit struct {
	Coord *ir_coordinate
	Lhs   string
	Val   string
}

type ir_StrLit struct {
	Coord *ir_coordinate
	Lhs   string
	Len   string
	Val   string
}

type ir_CsetLit struct {
	Coord *ir_coordinate
	Lhs   string
	Len   string
	Val   string
}

type ir_operator struct {
	Name  string
	Arity string
	Rval  string // may be nil
}

func (i ir_operator) String() string {
	return fmt.Sprintf("%s%s %v", i.Arity, i.Name, i.Rval)
}

type ir_MakeClosure struct {
	Coord *ir_coordinate
	Lhs   string
	Name  string
}

type ir_Move struct {
	Coord *ir_coordinate
	Lhs   string
	Rhs   string
}

type ir_MoveLabel struct {
	Coord *ir_coordinate
	Lhs   string
	Label string
}

type ir_Deref struct {
	Coord *ir_coordinate
	Lhs   string
	Value string
}

type ir_Assign struct {
	Coord  *ir_coordinate
	Target string
	Value  string
}

type ir_MakeList struct {
	Coord     *ir_coordinate
	Lhs       string
	ValueList []interface{} // heterogeneous
}

type ir_Field struct {
	Coord     *ir_coordinate
	Lhs       string // may be nil
	Expr      string
	Field     string
	FailLabel string
}

type ir_OpFunction struct {
	Coord      *ir_coordinate
	Lhs        string // may be nil
	Lhsclosure string // may be nil
	Fn         *ir_operator
	ArgList    []interface{} // heterogeneous
	FailLabel  string        // may be nil
}

type ir_Call struct {
	Coord      *ir_coordinate
	Lhs        string
	Lhsclosure string
	Fn         string
	ArgList    []interface{} // heterogeneous
	FailLabel  string        // may be nil
}

type ir_ResumeValue struct {
	Coord      *ir_coordinate
	Lhs        string // may be nil
	Lhsclosure string
	Closure    string
	FailLabel  string // may be nil
}

type ir_EnterInit struct {
	Coord      *ir_coordinate
	StartLabel string
}

type ir_Goto struct {
	Coord       *ir_coordinate
	TargetLabel string
}

type ir_IndirectGoto struct {
	Coord          *ir_coordinate
	TargetTmpLabel string
}

type ir_Succeed struct {
	Coord       *ir_coordinate
	Expr        string
	ResumeLabel string // may be nil
}

type ir_Fail struct {
	Coord *ir_coordinate
}

type ir_Create struct {
	Coord      *ir_coordinate
	Lhs        string
	CoexpLabel string
}

type ir_CoRet struct {
	Coord       *ir_coordinate
	Value       string
	ResumeLabel string
}

type ir_CoFail struct {
	Coord *ir_coordinate
}

type ir_ScanSwap struct {
	Coord   *ir_coordinate
	Subject string
	Pos     string
}

type ir_Unreachable struct {
	Coord *ir_coordinate
}

type ir_coordinate struct {
	File   string
	Line   string
	Column string
}

func (i ir_coordinate) String() string {
	return fmt.Sprintf("%s:%s:%s", i.File, i.Line, i.Column)
}

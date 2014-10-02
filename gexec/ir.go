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
	&ir_Tmp{},
	&ir_TmpLabel{},
	&ir_TmpClosure{},
	&ir_Label{},
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
	Accumulate string // may be nil -> ""
	LocalList  []string
	StaticList []string
	CodeList   []ir_chunk
	CodeStart  *ir_Label
	Lvalset    []string
}

type ir_chunk struct {
	Label    *ir_Label
	InsnList []interface{} // heterogeneous
}

type ir_Tmp struct {
	Name string
}

func (i ir_Tmp) String() string {
	return fmt.Sprintf("%s", i.Name)
}

type ir_TmpLabel struct {
	Name string
}

func (i ir_TmpLabel) String() string {
	return fmt.Sprintf("%s", i.Name)
}

type ir_TmpClosure struct {
	Name string
}

func (i ir_TmpClosure) String() string {
	return fmt.Sprintf("%s", i.Name)
}

type ir_Label struct {
	Value string
}

func (i ir_Label) String() string {
	return fmt.Sprintf("%s", i.Value)
}

type ir_Var struct {
	Coord *ir_coordinate
	Lhs   *ir_Tmp
	Name  string
}

type ir_Key struct {
	Coord     *ir_coordinate
	Lhs       *ir_Tmp // may be nil
	Name      string
	FailLabel *ir_Label
}

type ir_IntLit struct {
	Coord *ir_coordinate
	Lhs   *ir_Tmp
	Val   string
}

type ir_RealLit struct {
	Coord *ir_coordinate
	Lhs   *ir_Tmp
	Val   string
}

type ir_StrLit struct {
	Coord *ir_coordinate
	Lhs   *ir_Tmp
	Len   string
	Val   string
}

type ir_CsetLit struct {
	Coord *ir_coordinate
	Lhs   *ir_Tmp
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
	Lhs   *ir_Tmp
	Name  string
}

type ir_Move struct {
	Coord *ir_coordinate
	Lhs   *ir_Tmp
	Rhs   *ir_Tmp
}

type ir_MoveLabel struct {
	Coord *ir_coordinate
	Lhs   *ir_TmpLabel
	Label *ir_Label
}

type ir_Deref struct {
	Coord *ir_coordinate
	Lhs   *ir_Tmp
	Value *ir_Tmp
}

type ir_Assign struct {
	Coord  *ir_coordinate
	Target *ir_Tmp
	Value  *ir_Tmp
}

type ir_MakeList struct {
	Coord     *ir_coordinate
	Lhs       *ir_Tmp
	ValueList []interface{} // heterogeneous
}

type ir_Field struct {
	Coord     *ir_coordinate
	Lhs       *ir_Tmp // may be nil
	Expr      *ir_Tmp
	Field     string
	FailLabel *ir_Label
}

type ir_OpFunction struct {
	Coord      *ir_coordinate
	Lhs        *ir_Tmp        // may be nil
	Lhsclosure *ir_TmpClosure // may be nil
	Fn         *ir_operator
	ArgList    []interface{} // heterogeneous
	FailLabel  *ir_Label     // may be nil
}

type ir_Call struct {
	Coord      *ir_coordinate
	Lhs        *ir_Tmp
	Lhsclosure *ir_TmpClosure
	Fn         *ir_Tmp
	ArgList    []interface{} // heterogeneous
	FailLabel  *ir_Label     // may be nil
}

type ir_ResumeValue struct {
	Coord      *ir_coordinate
	Lhs        *ir_Tmp // may be nil
	Lhsclosure *ir_TmpClosure
	Closure    *ir_TmpClosure
	FailLabel  *ir_Label // may be nil
}

type ir_EnterInit struct {
	Coord      *ir_coordinate
	StartLabel *ir_Label
}

type ir_Goto struct {
	Coord       *ir_coordinate
	TargetLabel *ir_Label
}

type ir_IndirectGoto struct {
	Coord          *ir_coordinate
	TargetTmpLabel *ir_TmpLabel
}

type ir_Succeed struct {
	Coord       *ir_coordinate
	Expr        *ir_Tmp
	ResumeLabel *ir_Label // may be nil
}

type ir_Fail struct {
	Coord *ir_coordinate
}

type ir_Create struct {
	Coord      *ir_coordinate
	Lhs        *ir_Tmp
	CoexpLabel *ir_Label
}

type ir_CoRet struct {
	Coord       *ir_coordinate
	Value       *ir_Tmp
	ResumeLabel *ir_Label
}

type ir_CoFail struct {
	Coord *ir_coordinate
}

type ir_ScanSwap struct {
	Coord   *ir_coordinate
	Subject *ir_Tmp
	Pos     *ir_Tmp
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

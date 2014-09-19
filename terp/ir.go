//  ir.go -- intermediate representation structures

package main

//  list of IR structures for use by JSON input converter
var irlist = [...]struct {
	name    string
	example interface{}
}{
	{"ir_Invocable", ir_Invocable{}},
	{"ir_Link", ir_Link{}},
	{"ir_Record", ir_Record{}},
	{"ir_Global", ir_Global{}},
	{"ir_Function", ir_Function{}},
	{"ir_chunk", ir_chunk{}},
	{"ir_Tmp", ir_Tmp{}},
	{"ir_TmpLabel", ir_TmpLabel{}},
	{"ir_TmpClosure", ir_TmpClosure{}},
	{"ir_Label", ir_Label{}},
	{"ir_Var", ir_Var{}},
	{"ir_Key", ir_Key{}},
	{"ir_IntLit", ir_IntLit{}},
	{"ir_RealLit", ir_RealLit{}},
	{"ir_StrLit", ir_StrLit{}},
	{"ir_CsetLit", ir_CsetLit{}},
	{"ir_operator", ir_operator{}},
	{"ir_Move", ir_Move{}},
	{"ir_MoveLabel", ir_MoveLabel{}},
	{"ir_Deref", ir_Deref{}},
	{"ir_Assign", ir_Assign{}},
	{"ir_MakeList", ir_MakeList{}},
	{"ir_Field", ir_Field{}},
	{"ir_OpFunction", ir_OpFunction{}},
	{"ir_Call", ir_Call{}},
	{"ir_ResumeValue", ir_ResumeValue{}},
	{"ir_EnterInit", ir_EnterInit{}},
	{"ir_Goto", ir_Goto{}},
	{"ir_IndirectGoto", ir_IndirectGoto{}},
	{"ir_Succeed", ir_Succeed{}},
	{"ir_Fail", ir_Fail{}},
	{"ir_Create", ir_Create{}},
	{"ir_CoRet", ir_CoRet{}},
	{"ir_CoFail", ir_CoFail{}},
	{"ir_ScanSwap", ir_ScanSwap{}},
	{"ir_Unreachable", ir_Unreachable{}},
	{"ir_coordinate", ir_coordinate{}},
}

//  intermediate representation struct definitions
//  all fields must be capitalized for access by the reflection package

//  #%#% field types are only partially defined
//  #%#% slice fields can't be typed without first enhancing jfix()

type ir_Invocable struct {
	Coord         ir_coordinate
	NameList, All interface{}
}
type ir_Link struct {
	Coord    ir_coordinate
	NameList interface{}
}
type ir_Record struct {
	Coord           ir_coordinate
	Name, FieldList interface{}
}
type ir_Global struct {
	Coord    ir_coordinate
	NameList interface{}
}
type ir_Function struct {
	Coord ir_coordinate
	Name, ParamList, Accumulate, LocalList,
	StaticList, CodeList, CodeStart, Lvalset interface{}
}
type ir_chunk struct {
	Label    ir_Label
	InsnList interface{}
}
type ir_Tmp struct{ Name interface{} }
type ir_TmpLabel struct{ Name interface{} }
type ir_TmpClosure struct{ Name interface{} }
type ir_Label struct{ Value interface{} }
type ir_Var struct {
	Coord     ir_coordinate
	Lhs, Name interface{}
}
type ir_Key struct {
	Coord                ir_coordinate
	Lhs, Name, FailLabel interface{}
}
type ir_IntLit struct {
	Coord    ir_coordinate
	Lhs, Val interface{}
}
type ir_RealLit struct {
	Coord    ir_coordinate
	Lhs, Val interface{}
}
type ir_StrLit struct {
	Coord         ir_coordinate
	Lhs, Len, Val interface{}
}
type ir_CsetLit struct {
	Coord         ir_coordinate
	Lhs, Len, Val interface{}
}
type ir_operator struct{ Name, Arity, Rval interface{} }
type ir_Move struct {
	Coord    ir_coordinate
	Lhs, Rhs interface{}
}
type ir_MoveLabel struct {
	Coord      ir_coordinate
	Lhs, Label interface{}
}
type ir_Deref struct {
	Coord      ir_coordinate
	Lhs, Value interface{}
}
type ir_Assign struct {
	Coord         ir_coordinate
	Target, Value interface{}
}
type ir_MakeList struct {
	Coord          ir_coordinate
	Lhs, ValueList interface{}
}
type ir_Field struct {
	Coord                       ir_coordinate
	Lhs, Expr, Field, FailLabel interface{}
}
type ir_OpFunction struct {
	Coord                                   ir_coordinate
	Lhs, Lhsclosure, Fn, ArgList, FailLabel interface{}
}
type ir_Call struct {
	Coord                                   ir_coordinate
	Lhs, Lhsclosure, Fn, ArgList, FailLabel interface{}
}
type ir_ResumeValue struct {
	Coord                               ir_coordinate
	Lhs, Lhsclosure, Closure, FailLabel interface{}
}
type ir_EnterInit struct {
	Coord      ir_coordinate
	StartLabel interface{}
}
type ir_Goto struct {
	Coord       ir_coordinate
	TargetLabel interface{}
}
type ir_IndirectGoto struct {
	Coord          ir_coordinate
	TargetTmpLabel interface{}
}
type ir_Succeed struct {
	Coord             ir_coordinate
	Expr, ResumeLabel interface{}
}
type ir_Fail struct{ Coord ir_coordinate }
type ir_Create struct {
	Coord           ir_coordinate
	Lhs, CoexpLabel interface{}
}
type ir_CoRet struct {
	Coord              ir_coordinate
	Value, ResumeLabel interface{}
}
type ir_CoFail struct{ Coord ir_coordinate }
type ir_ScanSwap struct {
	Coord        ir_coordinate
	Subject, Pos interface{}
}
type ir_Unreachable struct{ Coord ir_coordinate }

type ir_coordinate struct {
	File string
	Line, Column/*also*/ string
}

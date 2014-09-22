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
	{"ir_MakeClosure", ir_MakeClosure{}},
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
//  #%#% n.b. jfix can't handle slices of other than string or interface{}

type ir_Invocable struct {
	Coord    ir_coordinate
	NameList []string
	All      interface{}
}
type ir_Link struct {
	Coord    ir_coordinate
	NameList []string
}
type ir_Record struct {
	Coord     ir_coordinate
	Name      string
	FieldList []string
}
type ir_Global struct {
	Coord    ir_coordinate
	NameList []string
}
type ir_Function struct {
	Coord      ir_coordinate
	Name       string
	ParamList  []string
	Accumulate interface{}
	LocalList  []string
	StaticList []string
	CodeList   []interface{} //#%#% []ir_chunk
	CodeStart  ir_Label
	Lvalset    []string
}
type ir_chunk struct {
	Label    ir_Label
	InsnList []interface{} // heterogeneous
}
type ir_Tmp struct {
	Name string
}
type ir_TmpLabel struct {
	Name string
}
type ir_TmpClosure struct {
	Name string
}
type ir_Label struct {
	Value string
}
type ir_Var struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Name  string
}
type ir_Key struct {
	Coord     ir_coordinate
	Lhs       ir_Tmp
	Name      string
	FailLabel string
}
type ir_IntLit struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Val   string
}
type ir_RealLit struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Val   string
}
type ir_StrLit struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Len   string
	Val   string
}
type ir_CsetLit struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Len   string
	Val   string
}
type ir_operator struct {
	Name  string
	Arity string
	Rval  interface{}
}
type ir_MakeClosure struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Name  string
}
type ir_Move struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Rhs   interface{}
}
type ir_MoveLabel struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Label ir_Label
}
type ir_Deref struct {
	Coord ir_coordinate
	Lhs   ir_Tmp
	Value string
}
type ir_Assign struct {
	Coord  ir_coordinate
	Target interface{}
	Value  string
}
type ir_MakeList struct {
	Coord     ir_coordinate
	Lhs       ir_Tmp
	ValueList []interface{} // heterogeneous
}
type ir_Field struct {
	Coord     ir_coordinate
	Lhs       ir_Tmp
	Expr      interface{}
	Field     string
	FailLabel string
}
type ir_OpFunction struct {
	Coord      ir_coordinate
	Lhs        ir_Tmp
	Lhsclosure interface{}
	Fn         interface{}
	ArgList    []interface{} // heterogeneous
	FailLabel  ir_Label
}
type ir_Call struct {
	Coord      ir_coordinate
	Lhs        ir_Tmp
	Lhsclosure interface{}
	Fn         interface{}
	ArgList    []interface{} // heterogeneous
	FailLabel  ir_Label
}
type ir_ResumeValue struct {
	Coord      ir_coordinate
	Lhs        ir_Tmp
	Lhsclosure interface{}
	Closure    interface{}
	FailLabel  ir_Label
}
type ir_EnterInit struct {
	Coord      ir_coordinate
	StartLabel ir_Label
}
type ir_Goto struct {
	Coord       ir_coordinate
	TargetLabel ir_Label
}
type ir_IndirectGoto struct {
	Coord          ir_coordinate
	TargetTmpLabel ir_Label
}
type ir_Succeed struct {
	Coord       ir_coordinate
	Expr        interface{}
	ResumeLabel ir_Label
}
type ir_Fail struct {
	Coord ir_coordinate
}
type ir_Create struct {
	Coord      ir_coordinate
	Lhs        ir_Tmp
	CoexpLabel ir_Label
}
type ir_CoRet struct {
	Coord       ir_coordinate
	Value       string
	ResumeLabel ir_Label
}
type ir_CoFail struct {
	Coord ir_coordinate
}
type ir_ScanSwap struct {
	Coord   ir_coordinate
	Subject interface{}
	Pos     interface{}
}
type ir_Unreachable struct {
	Coord ir_coordinate
}
type ir_coordinate struct {
	File   string
	Line   string
	Column string
}

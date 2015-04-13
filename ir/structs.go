//  structs.go -- intermediate representation structures

package ir

import (
	"reflect"
)

//  list of IR struct examples for use by JSON input converter
var irlist = [...]interface{}{
	&Ir_Record{},
	&Ir_Global{},
	&Ir_Initial{},
	&Ir_Function{},
	&Ir_chunk{},
	&Ir_NoOp{}, // not normally seen, but allowed as a comment
	&Ir_Catch{},
	&Ir_EnterScope{},
	&Ir_ExitScope{},
	&Ir_Var{},
	&Ir_Key{},
	&Ir_NilLit{},
	&Ir_IntLit{},
	&Ir_RealLit{},
	&Ir_StrLit{},
	&Ir_MakeClosure{},
	&Ir_Move{},
	&Ir_MoveLabel{},
	&Ir_MakeList{},
	&Ir_Field{},
	&Ir_OpFunction{},
	&Ir_Call{},
	&Ir_ResumeValue{},
	&Ir_Goto{},
	&Ir_IndirectGoto{},
	&Ir_Succeed{},
	&Ir_Fail{},
	&Ir_Create{},
	&Ir_CoRet{},
	&Ir_CoFail{},
	&Ir_Select{},
	&Ir_SelectCase{},
	&Ir_NoValue{},     // seen only if unoptimized; not implemented
	&Ir_Unreachable{}, // seen only if unoptimized; not implemented
}

//  struct table indexed by type names
var irtable = make(map[string]reflect.Type)

func init() {
	for _, ir := range irlist {
		t := reflect.TypeOf(ir).Elem()
		name := t.Name()
		irtable[Capitalize(name)] = t
		irtable[DeCapit(name)] = t
	}
}

//  intermediate representation struct definitions
//  all fields must be capitalized for access by the reflection package

type Ir_Record struct {
	Coord      string
	Name       string
	ExtendsRec string
	ExtendsPkg string
	FieldList  []string
	Namespace  string
}

type Ir_Initial struct {
	Coord     string
	Fn        string
	Namespace string
}

type Ir_Global struct {
	Coord     string
	Name      string
	Namespace string
	Fn        string
}

type Ir_Function struct {
	Coord       string
	Name        string
	ParamList   []string
	Accumulate  string // may be nil
	LocalList   []string
	StaticList  []string
	UnboundList []string
	CodeList    []Ir_chunk
	CodeStart   string
	Parent      string
	Namespace   string
}

type Ir_chunk struct {
	Label    string
	InsnList []interface{} // heterogeneous
}

type Ir_NoOp struct {
	Coord   string
	Comment string
}

type Ir_Catch struct {
	Coord string
	Lhs   string
	Fn    string
}

type Ir_EnterScope struct {
	Coord       string
	NameList    []string
	DynamicList []string
	Scope       string
	ParentScope string
}

type Ir_ExitScope struct {
	Coord       string
	NameList    []string
	DynamicList []string
	Scope       string
}

type Ir_Var struct {
	Coord     string
	Lhs       string
	Name      string
	Namespace string
	Scope     string
	Rval      string // may be nil
}

type Ir_Key struct {
	Coord string
	Lhs   string // may be nil
	Name  string
	Scope string
}

type Ir_NilLit struct {
	Coord string
	Lhs   string
}

type Ir_IntLit struct {
	Coord string
	Lhs   string
	Val   string
}

type Ir_RealLit struct {
	Coord string
	Lhs   string
	Val   string
}

type Ir_StrLit struct {
	Coord string
	Lhs   string
	Len   string // length of the UTF-8 encoding
	Val   string // individual bytes of the UTF-8 encoding
}

type Ir_MakeClosure struct {
	Coord string
	Lhs   string
	Name  string
}

type Ir_Move struct {
	Coord string
	Lhs   string
	Rhs   string
}

type Ir_MoveLabel struct {
	Coord string
	Lhs   string
	Label string
}

type Ir_MakeList struct {
	Coord     string
	Lhs       string
	ValueList []interface{} // heterogeneous
}

type Ir_Field struct {
	Coord string
	Lhs   string // may be nil
	Expr  string
	Field string
	Rval  string // may be nil
}

type Ir_OpFunction struct {
	Coord      string
	Lhs        string // may be nil
	Lhsclosure string // may be nil
	Fn         string
	ArgList    []interface{} // heterogeneous
	Rval       string        // may be nil
	FailLabel  string        // may be nil
}

type Ir_Call struct {
	Coord      string
	Lhs        string
	Lhsclosure string
	Fn         string
	ArgList    []interface{} // heterogeneous
	NameList   []string
	FailLabel  string // may be nil
	Scope      string
}

type Ir_ResumeValue struct {
	Coord      string
	Lhs        string // may be nil
	Lhsclosure string
	Closure    string
	FailLabel  string // may be nil
}

type Ir_Goto struct {
	Coord       string
	TargetLabel string
}

type Ir_IndirectGoto struct {
	Coord          string
	TargetTmpLabel string
	LabelList      []string
}

type Ir_Succeed struct {
	Coord       string
	Expr        string
	ResumeLabel string // may be nil
}

type Ir_Fail struct {
	Coord string
}

type Ir_Create struct {
	Coord      string
	Lhs        string
	CoexpLabel string
}

type Ir_CoRet struct {
	Coord       string
	Value       string
	ResumeLabel string
}

type Ir_CoFail struct {
	Coord string
}

type Ir_Select struct {
	Coord     string
	CaseList  []Ir_SelectCase
	FailLabel string
}

type Ir_SelectCase struct {
	Coord     string
	Kind      string // "send" | "receive" | "default"
	Lhs       string
	Rhs       string
	BodyLabel string
}

type Ir_NoValue struct {
	Coord string
	Lhs   string
}

type Ir_Unreachable struct {
	Coord string
}

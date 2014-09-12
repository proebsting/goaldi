//  vproc.go -- VProcedure, the Goaldi type "procedure"

package goaldi

//  Procedure function prototype
type Procedure func(...Value) (Value, *Closure)

//  Procedure value
type VProcedure struct {
	name  string
	entry Procedure
}

//  NewProcedure(name, func) -- construct a procedure value
func NewProcedure(name string, entry Procedure) *VProcedure {
	return &VProcedure{name, entry}
}

//  VProcedure.String -- return "procname()"
func (v *VProcedure) String() string {
	return v.name + "()"
}

//  VProcedure.Type -- return "procedure"
func (v *VProcedure) Type() Value {
	return type_procedure
}

//  ICall interface
type ICall interface {
	Call(...Value) (Value, *Closure)
}

//  VProcedure.Call(args) -- invoke a procedure
func (v *VProcedure) Call(args ...Value) (Value, *Closure) {
	return v.entry(args)
}

var type_procedure = NewString("procedure")

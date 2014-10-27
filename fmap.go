//  fmap.go -- map functions

package goaldi

//  This init function adds a set of Go functions to the standard library.
//  Note that delete() and member() are here as static (not method) funcs.
func init() {
	// Goaldi procedures
	LibProcedure("map", Map)
	LibProcedure("member", Member)
	LibProcedure("delete", Delete)
}

//  Map() -- return a new map
func Map(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("map", a)
	return Return(NewMap())
}

type IMember interface {
	Member(Value) Value
}

//  Member(x, k) -- check k for membership in x
func Member(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("member", a)
	x := ProcArg(a, 0, NilValue)
	k := ProcArg(a, 1, NilValue)
	return Return(x.(IMember).Member(k))
}

type IDelete interface {
	Delete(Value) Value
}

//  Delete(x, k) -- delete k from x
func Delete(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("delete", a)
	x := ProcArg(a, 0, NilValue)
	k := ProcArg(a, 1, NilValue)
	return Return(x.(IDelete).Delete(k))
}

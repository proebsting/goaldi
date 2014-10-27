//  fmap.go -- map functions

package goaldi

//  This init function adds a set of Go functions to the standard library.
func init() {
	// Goaldi procedures
	LibProcedure("map", Map)
}

//  Map() -- return a new map
func Map(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("map", a)
	return Return(NewMap())
}

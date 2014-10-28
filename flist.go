//  flist.go -- list functions

package goaldi

//  This init function adds a set of Go functions to the standard library.
func init() {
	// Goaldi procedures
	LibProcedure("list", List)
}

//  List() -- return a new list
func List(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("list", a)
	n := int(ProcArg(a, 0, ZERO).(Numerable).ToNumber().Val())
	return Return(NewList(n))
}

//------------------------------------  Field:  L.s  implements methods

func (v *VList) Field(f string) Value {
	switch f {
	case "push":
		return MVFunc(v.Push)
	case "pop":
		return MVFunc(v.Pop)
	case "get":
		return MVFunc(v.Get)
	case "put":
		return MVFunc(v.Put)
	case "pull":
		return MVFunc(v.Pull)
	default:
		panic(&RunErr{"Undefined method: " + f, v})
	}
}

//------------------------------------  Member:  L.push(x...)

func (v *VList) Push(args ...Value) (Value, *Closure) {
	return v.Grow(true, "L.push", args...)
}

//------------------------------------  Member:  L.pop(x...)

func (v *VList) Pop(args ...Value) (Value, *Closure) {
	return v.Snip(true, "L.pop", args...)
}

//------------------------------------  Member:  L.get(x...)

func (v *VList) Get(args ...Value) (Value, *Closure) {
	return v.Snip(true, "L.get", args...)
}

//------------------------------------  Member:  L.put(x...)

func (v *VList) Put(args ...Value) (Value, *Closure) {
	return v.Grow(false, "L.put", args...)
}

//------------------------------------  Member:  L.pull(x...)

func (v *VList) Pull(args ...Value) (Value, *Closure) {
	return v.Snip(false, "L.pull", args...)
}

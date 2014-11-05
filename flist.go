//  flist.go -- list functions

package goaldi

//  Declare methods
var ListMethods = map[string]interface{}{
	"type":  (*VList).Type,
	"copy":  (*VList).Copy,
	"image": Image,
	"push":  (*VList).Push,
	"pop":   (*VList).Pop,
	"get":   (*VList).Get,
	"put":   (*VList).Put,
	"pull":  (*VList).Pull,
}

//  VList.Field implements methods
func (v *VList) Field(f string) Value {
	return GetMethod(ListMethods, v, f)
}

//  Declare constructor function
func init() {
	LibProcedure("list", List)
}

//  List(n, x) -- return a new list of n elements initialize to x
func List(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("list", a)
	n := int(ProcArg(a, 0, ZERO).(Numerable).ToNumber().Val())
	x := ProcArg(a, 1, NilValue)
	return Return(NewList(n, x))
}

//------------------------------------  Push:  L.push(x...)

func (v *VList) Push(args ...Value) (Value, *Closure) {
	return v.Grow(true, "L.push", args...)
}

//------------------------------------  Pop:  L.pop(x...)

func (v *VList) Pop(args ...Value) (Value, *Closure) {
	return v.Snip(true, "L.pop", args...)
}

//------------------------------------  Get:  L.get(x...)

func (v *VList) Get(args ...Value) (Value, *Closure) {
	return v.Snip(true, "L.get", args...)
}

//------------------------------------  Put:  L.put(x...)

func (v *VList) Put(args ...Value) (Value, *Closure) {
	return v.Grow(false, "L.put", args...)
}

//------------------------------------  Pull:  L.pull(x...)

func (v *VList) Pull(args ...Value) (Value, *Closure) {
	return v.Snip(false, "L.pull", args...)
}

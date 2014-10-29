//  fmap.go -- map functions and methods

package goaldi

//  Declare methods
var MapMethods = map[string]interface{}{
	"type":   (*VMap).Type,
	"copy":   (*VMap).Copy,
	"image":  Image,
	"member": (*VMap).Member,
	"delete": (*VMap).Delete,
}

//  VMap.Field implements methods
func (v *VMap) Field(f string) Value {
	return GetMethod(MapMethods, v, f)
}

//  Declare constructor function
func init() {
	// Goaldi procedures
	LibProcedure("map", Map)
}

//  Map() -- return a new map
func Map(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("map", a)
	return Return(NewMap())
}

//------------------------------------  Member:  M.member(x)

func (v *VMap) Member(args ...Value) (Value, *Closure) {
	defer Traceback("M.member", args)
	key := args[0]
	if v.data[MapIndex(key)] != nil {
		return Return(key)
	} else {
		return Fail()
	}
}

//------------------------------------  Delete:  M.delete(x)

func (v *VMap) Delete(args ...Value) (Value, *Closure) {
	defer Traceback("M.delete", args)
	key := args[0]
	x := MapIndex(key)
	delete(v.data, x)
	if len(v.data) != len(v.klist) {
		// delete was successful; need to remove from klist
		for i := len(v.klist) - 1; i >= 0; i-- {
			if Identical(key, v.klist[i]) != nil {
				v.klist[i] = v.klist[len(v.klist)-1]
				v.klist = v.klist[:len(v.klist)-1]
				break
			}
		}
	}
	if len(v.data) != len(v.klist) {
		panic(&RunErr{"inconsistent map", v})
	}
	return Return(v)
}

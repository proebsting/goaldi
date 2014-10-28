//  fmap.go -- map functions and methods

package goaldi

//  Declare constructor function for standard library
func init() {
	// Goaldi procedures
	LibProcedure("map", Map)
}

//  Map() -- return a new map
func Map(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("map", a)
	return Return(NewMap())
}

//------------------------------------  Field:  e1.s  implements methods

func (v *VMap) Field(f string) Value {
	//#%#% check first for "member" and "delete" method references,
	//#%#% but allow any other string as a index (?!)
	switch f {
	case "member":
		return MVFunc(v.Member)
	case "delete":
		return MVFunc(v.Delete)
	default:
		return &vMapSlot{v, NewString(f)}
	}
}

//------------------------------------  Member:  e1.member(e2)

func (v *VMap) Member(args ...Value) (Value, *Closure) {
	defer Traceback("M.member", args)
	key := args[0]
	if v.data[MapIndex(key)] != nil {
		return Return(key)
	} else {
		return Fail()
	}
}

//------------------------------------  Delete:  e1.delete(e2)

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

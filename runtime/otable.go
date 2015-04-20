//  otable.go -- table operations

package runtime

import (
	"math/rand"
	"reflect"
)

//  VTable.Size -- return the table size
func (T *VTable) Size() Value {
	return NewNumber(float64(len(T.data)))
}

//  VTable.Choose -- return random element as (key,value) pair
func (T *VTable) Choose(lval Value) Value {
	return ChooseMap(T.data)
}

//  VTable.Take -- return random element as (key,value) pair
func (T *VTable) Take(lval Value) Value {
	return TakeMap(T.data)
}

//  VTable.Dispense -- generate table contents as (key,value) pairs
func (T *VTable) Dispense(lval Value) (Value, *Closure) {
	return DispenseMap(T.data)
}

//  T.Index(lval, x) implements the [] operator.
func (T *VTable) Index(lval Value, x Value) Value {
	return TrapMap(T, x)
}

//  ChooseMap returns a key/value pair from any Goaldi table or Go map
func ChooseMap(T interface{} /*anymap*/) Value {
	mv := reflect.ValueOf(T)
	klist := mv.MapKeys()
	n := len(klist)
	if n == 0 { // if map empty
		return nil // fail
	}
	i := rand.Intn(n)
	k := Import(klist[i].Interface())
	v := Import(mv.MapIndex(klist[i]).Interface())
	return ElemType.New([]Value{k, v})
}

//  TakeMap returns a key/value pair after removal from the underlying map.
func TakeMap(T interface{} /*anymap*/) Value {
	kv := ChooseMap(T)
	if kv != nil {
		key := kv.(*VRecord).Data[0]
		GoMapDelete(T, key)
	}
	return kv
}

//  DispenseMap generates key/value pairs for any Goaldi table or Go map
func DispenseMap(T interface{} /*anymap*/) (Value, *Closure) {
	mv := reflect.ValueOf(T)
	klist := mv.MapKeys()
	i := -1
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		for {
			i++
			if i >= len(klist) {
				return Fail()
			}
			x := mv.MapIndex(klist[i])
			if x.IsValid() { // if didn't disappear while suspended
				k := Import(klist[i].Interface())
				v := Import(x.Interface())
				return ElemType.New([]Value{k, v}), c
			}
		}
		return Fail()
	}}
	return c.Resume()
}

//  -------------------------- trapped references ---------------------

//  vMapTrap is a trapped reference T[k] into a Goaldi table or Go map
type vMapTrap struct {
	gmap  bool          // true if a Goaldi (not Go) map
	dfval Value         // default value if a Goaldi map
	mapv  reflect.Value // underlying Go map
	keyv  reflect.Value // key converted to appropriate Go type
}

//  TrapMap(T,k) creates a trapped variable for T[k]
func TrapMap(x Value, key Value) *vMapTrap {
	if T, ok := x.(*VTable); ok {
		// this is a Goaldi table; must convert string or number key
		tv := reflect.ValueOf(T.data)
		return &vMapTrap{true, T.dfval, tv, reflect.ValueOf(GoKey(key))}
	} else {
		tv := reflect.ValueOf(x)
		// otherwise, key will be converted by passfunc
		return &vMapTrap{false, nil, tv, passfunc(tv.Type().Key())(key)}
	}
}

//  vMapTrap.Exists() returns true if the reference matches an existing key
func (t *vMapTrap) Exists() bool {
	return t.mapv.MapIndex(t.keyv).IsValid()
}

//  vMapTrap.Deref() returns the indexed value, or the default if not found
func (t *vMapTrap) Deref() Value {
	v := t.mapv.MapIndex(t.keyv)
	if v.IsValid() {
		return Import(v.Interface()) // identity function for VTable values
	} else {
		return t.dfval // not found in map
	}
}

//  vMapTrap.Assign(x) stores x as a map entry using the trapped key
func (t *vMapTrap) Assign(x Value) IVariable {
	if t.gmap { // if Goaldi table
		t.mapv.SetMapIndex(t.keyv, reflect.ValueOf(x))
	} else {
		t.mapv.SetMapIndex(t.keyv, passfunc(t.mapv.Type().Elem())(x))
	}
	return t
}

//  vMapTrap.Delete() removes the entry, if any, associated with the trapped key
func (t *vMapTrap) Delete() {
	t.mapv.SetMapIndex(t.keyv, reflect.Value{})
}

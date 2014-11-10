//  vmap.go -- VMap, the Goaldi type "map"

package goaldi

import (
	"bytes"
	"fmt"
	"reflect"
)

//  VMap implements a native Goaldi map.
//  It behaves similarly to an external map except tha
//  only strings and numbers are converted before use as keys.
//  (Unconverted "identical" values would be seen as distinct.)
type VMap map[Value]Value

//  NewMap -- construct a new Goaldi map
func NewMap() VMap {
	return make(map[Value]Value)
}

//  VMap.String -- default conversion to Go string returns "M:size"
func (m VMap) String() string {
	return fmt.Sprintf("M:%d", len(m))
}

//  VMap.GoString -- convert to Go string for image() and printf("%#v")
//
//  For utility and reproducibility, we assume it's worth the cost
//  to sort the map in key order.
func (m VMap) GoString() string {
	if len(m) == 0 {
		return "map{}"
	}
	l, _ := m.Sort(ONE) // sort on key values
	var b bytes.Buffer
	fmt.Fprintf(&b, "map{")
	for _, e := range l.(*VList).data {
		r := e.(*VStruct)
		fmt.Fprintf(&b, "%v:%v,", r.Data[0], r.Data[1])
	}
	s := b.Bytes()
	s[len(s)-1] = '}'
	return string(s)
}

//  VMap.Rank returns rMap
func (v VMap) Rank() int {
	return rMap
}

//  VMap.Type -- return "map"
func (m VMap) Type() Value {
	return type_map
}

var type_map = NewString("map")

//  VMap.Copy returns a duplicate of itself
func (m VMap) Copy() Value {
	r := NewMap()
	for k, v := range m {
		r[k] = v
	}
	return r
}

//  VMap.Identical checks equality for the === operator
func (m VMap) Identical(x Value) Value {
	m2, ok := x.(VMap)
	if ok && reflect.ValueOf(m).Pointer() == reflect.ValueOf(m2).Pointer() {
		return x
	} else {
		return nil
	}
}

//  VMap.Import returns itself
func (v VMap) Import() Value {
	return v
}

//  VMap.Export returns itself.
//  Go extensions may wish to use v.Index().Deref(), v.Delete(), etc.
//  to ensure proper conversion of keys.
func (v VMap) Export() interface{} {
	return v
}

//  -------------------------- trapped references ---------------------

//  vMapTrap is a trapped map reference m[k] to a Goaldi or Go map
type vMapTrap struct {
	mapv reflect.Value // underlying Go map
	keyv reflect.Value // key converted to appropriate Go type
}

//  TrapMap(m,k) creates a trapped variable for m[k]
func TrapMap(m Value, key Value) *vMapTrap {
	mv := reflect.ValueOf(m)
	if _, ok := mv.Interface().(VMap); ok {
		// this is a native map; must convert string or number key
		switch t := key.(type) {
		case *VString:
			key = t.ToUTF8()
		case *VNumber:
			key = t.Val()
		default:
			// nothing: use key as is
		}
	} else {
		// convert key to export form for external map
		key = Export(key)
	}
	return &vMapTrap{mv, passfunc(mv.Type().Key())(key)}
}

//  vMapTrap.Exists() returns true if the reference matches an existing key
func (t *vMapTrap) Exists() bool {
	return t.mapv.MapIndex(t.keyv).IsValid()
}

//  vMapTrap.Deref() returns the indexed value, or the nil value if not found
func (t *vMapTrap) Deref() Value {
	v := t.mapv.MapIndex(t.keyv)
	if v.IsValid() {
		return Import(v.Interface()) // identity function for VMap values
	} else {
		return NilValue // not found in map
	}
}

//  vMapTrap.Assign(x) stores x as a map entry using the trapped key
func (t *vMapTrap) Assign(x Value) IVariable {
	t.mapv.SetMapIndex(t.keyv, passfunc(t.mapv.Type().Elem())(x))
	return t
}

//  vMapTrap.Delete() removes the entry, if any, associated with the trapped key
func (t *vMapTrap) Delete() {
	t.mapv.SetMapIndex(t.keyv, reflect.Value{})
}

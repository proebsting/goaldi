//  vmap.go -- VMap, the Goaldi type "map"

package goaldi

import (
	"fmt"
)

//  A Goaldi map wraps a Go map.
//  Goaldi numbers and strings must be converted for use as keys,
//  because otherwise two equal values might be treated distinctly.
//
//  To keep ?m and !m fast (at the cost of slow deletion)
//  we keep an additional separate list of Goaldi keys.
type VMap struct {
	data  map[Value]Value // map indexed by canonical values
	klist []Value         // list of Goaldi keys; len(klist)==len(data)
}

//  NewMap -- construct a new Goaldi map
func NewMap() *VMap {
	return &VMap{make(map[Value]Value), nil}
}

//  VMap.String -- default conversion to Go string
func (v *VMap) String() string {
	return fmt.Sprint("map()")
}

//  VMap.GoString -- convert to Go string for image() and printf("%#v")
func (v *VMap) GoString() string {
	return fmt.Sprintf("map(%d)", len(v.data))
}

//  VMap.Type -- return "map"
func (v *VMap) Type() Value {
	return type_map
}

var type_map = NewString("map")

//  VMap.Copy returns a duplicate of itself
func (v *VMap) Copy() Value {
	w := NewMap()
	for _, k := range v.klist {
		vx := v.Index(nil, k).(IVariable)
		wx := w.Index(nil, k).(IVariable)
		wx.Assign(vx.Deref())
	}
	return w
}

//  VMap.Import returns itself
func (v *VMap) Import() Value {
	return v
}

//  VMap.Export returns itself.
//  Go extensions should use v.Index().Deref(), v.Delete(), etc. for access.
func (v *VMap) Export() interface{} {
	return v
}

//  -------------------------- trapped references ---------------------

type vMapSlot struct {
	gmap *VMap // Goaldi map
	gkey Value // Goaldi key
}

//  vMapSlot.Deref() -- extract value of reference for use as an rvalue
func (ms *vMapSlot) Deref() Value {
	r := ms.gmap.data[MapIndex(ms.gkey)]
	if r == nil {
		return NilValue
	} else {
		return r
	}
}

//  vMapSlot.String() -- show string representation: produces (map[k])
func (ms *vMapSlot) String() string {
	return fmt.Sprintf("(map[%v])", ms.gkey)
}

//  vMapSlot.Assign -- store value in map
func (ms *vMapSlot) Assign(v Value) IVariable {
	m := ms.gmap
	m.data[MapIndex(ms.gkey)] = v
	if len(m.data) > len(m.klist) { // if the map grew
		m.klist = append(m.klist, ms.gkey) // note the new key
	}
	if len(m.data) != len(m.klist) {
		panic(&RunErr{"inconsistent map", m})
	}
	return ms
}

//  -------------------------- internal function ---------------------

//  MapIndex(v) -- return canonical value for indexing a map
func MapIndex(v Value) interface{} {
	switch t := v.(type) {
	case *VString:
		return t.ToUTF8()
	case *VNumber:
		return t.Val()
	default:
		return v
	}
}

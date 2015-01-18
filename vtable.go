//	vtable.go -- VTable, the Goaldi type "table"
//
//	Implementation:
//	A Goaldi table is just a type name VTable attached to a Go map[Value]Value.
//	This distinguishes it from an external Go map and allows attaching
//	(internal) methods.  Goaldi string and number indexes are converted
//	to Go string and float64 values.

package goaldi

import (
	"bytes"
	"fmt"
	"reflect"
)

//  VTable implements a native Goaldi table.
//  It behaves similarly to an external Go map except that
//  only strings and numbers are converted before use as keys.
//  (Unconverted "identical" values would be seen as distinct.)
type VTable map[Value]Value

//  NewTable -- construct a new Goaldi table
func NewTable() VTable {
	return make(map[Value]Value)
}

//  TableType is the table instance of type type.
var TableType = NewType("table",
	DefProc(Table, "table", "", "create an empty table"))

//  VTable.String -- default conversion to Go string returns "T:size"
func (m VTable) String() string {
	return fmt.Sprintf("T:%d", len(m))
}

//  VTable.GoString -- convert to Go string for image() and printf("%#v")
//
//  For utility and reproducibility, we assume it's worth the cost
//  to sort the table in key order.
func (m VTable) GoString() string {
	if len(m) == 0 {
		return "table{}"
	}
	l, _ := m.Sort(ONE) // sort on key values
	var b bytes.Buffer
	fmt.Fprintf(&b, "table{")
	for _, e := range l.(*VList).data {
		r := e.(*VRecord)
		fmt.Fprintf(&b, "%v:%v,", r.Data[0], r.Data[1])
	}
	s := b.Bytes()
	s[len(s)-1] = '}'
	return string(s)
}

//  VTable.Rank returns rTable
func (v VTable) Rank() int {
	return rTable
}

//  VTable.Type -- return the table type
func (m VTable) Type() Value {
	return TableType
}

//  VTable.Copy returns a duplicate of itself
func (m VTable) Copy() Value {
	r := NewTable()
	for k, v := range m {
		r[k] = v
	}
	return r
}

//  VTable.Import returns itself
func (v VTable) Import() Value {
	return v
}

//  VTable.Export returns itself.
//  Go extensions may wish to use v.Index().Deref(), v.Delete(), etc.
//  to ensure proper conversion of keys.
func (v VTable) Export() interface{} {
	return v
}

//  -------------------------- trapped references ---------------------

//  vMapTrap is a trapped reference m[k] into a Goaldi table or Go map
type vMapTrap struct {
	mapv reflect.Value // underlying Go map
	keyv reflect.Value // key converted to appropriate Go type
}

//  TrapMap(m,k) creates a trapped variable for m[k]
func TrapMap(m Value, key Value) *vMapTrap {
	mv := reflect.ValueOf(m)
	if _, ok := mv.Interface().(VTable); ok {
		// this is a Goaldi table; must convert string or number key
		switch t := key.(type) {
		case *VString:
			key = t.ToUTF8()
		case *VNumber:
			key = t.Val()
		default:
			// nothing: use key as is
		}
	} // else key will be converted by passfunc
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
		return Import(v.Interface()) // identity function for VTable values
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

//	vtable.go -- VTable, the Goaldi type "table"

package goaldi

import (
	"bytes"
	"fmt"
)

//	A Goaldi table combines a default value with a Go map[Value]Value under
//  the name of VTable.  Goaldi string and number indexes are converted to Go
//	Go string and float64 values by the GoKey() function (also used for sets).
type VTable struct {
	data  map[Value]Value // underlying Go map
	dfval Value           // default value
}

//  NewTable -- construct a new Goaldi table
func NewTable(dfval Value) *VTable {
	return &VTable{make(map[Value]Value), dfval}
}

//  TableType is the table instance of type type.
var TableType = NewType("table", "T", rTable, Table, TableMethods,
	"table", "x", "create an empty table")

//  VTable.String -- default conversion to Go string returns "T:size"
func (T *VTable) String() string {
	return fmt.Sprintf("T:%d", len(T.data))
}

//  VTable.GoString -- convert to Go string for image() and printf("%#v")
//
//  For utility and reproducibility, we pay the cost to sort into key order.
func (T *VTable) GoString() string {
	if len(T.data) == 0 {
		return "table{}"
	}
	l, _ := T.Sort(ONE) // sort on key values
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

//  VTable.Type -- return the table type
func (T *VTable) Type() IRank {
	return TableType
}

//  VTable.Copy returns a duplicate of itself
func (T *VTable) Copy() Value {
	r := NewTable(T.dfval)
	for k, v := range T.data {
		r.data[k] = v
	}
	return r
}

//  VTable.Before compares two tables for sorting
func (T *VTable) Before(b Value, i int) bool {
	return false // no ordering defined
}

//  VTable.Import returns itself
func (T *VTable) Import() Value {
	return T
}

//  VTable.Export returns its underlying Go map.  The table default is lost.
//  Go extensions may wish to use GoKey() for proper conversion of keys.
func (T *VTable) Export() interface{} {
	return T.data
}

//  vrecord.go -- a user-defined (usually) Goaldi record structure

package goaldi

import (
	"bytes"
	"fmt"
)

type VRecord struct {
	Defn *VDefn  // underlying type definition
	Data []Value // current data values
}

//  VRecord.String -- conversion to Go string returns "name{}"
func (v *VRecord) String() string {
	return v.Defn.Name + "{}"
}

//  VRecord.GoString -- returns string for image() and printf("%#v")
func (v *VRecord) GoString() string {
	if len(v.Data) == 0 {
		return v.Defn.Name + "{}"
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s{", v.Defn.Name)
	for i, x := range v.Data {
		fmt.Fprintf(&b, "%v:%v,", v.Defn.Flist[i], x)
	}
	s := b.Bytes()
	s[len(s)-1] = '}'
	return string(s)
}

//  VRecord.Type returns the underlying constructor
func (v *VRecord) Type() IRanking {
	return v.Defn
}

//  VRecord.Copy returns a distinct copy of itself
func (v *VRecord) Copy() Value {
	r := &VRecord{v.Defn, make([]Value, len(v.Data))}
	copy(r.Data, v.Data)
	return r
}

//  VRecord.Import returns itself
func (v *VRecord) Import() Value {
	return v
}

//  VRecord.Export returns itself
func (v *VRecord) Export() interface{} {
	return v
}

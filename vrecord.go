//  vrecord.go -- a user-defined (usually) Goaldi record structure

package goaldi

import (
	"bytes"
	"fmt"
)

type VRecord struct {
	Ctor *VCtor  // underlying type definition
	Data []Value // current data values
}

//  VRecord.String -- conversion to Go string returns "name{}"
func (v *VRecord) String() string {
	return v.Ctor.RecName + "{}"
}

//  VRecord.GoString -- returns string for image() and printf("%#v")
func (v *VRecord) GoString() string {
	if len(v.Data) == 0 {
		return v.Ctor.RecName + "{}"
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s{", v.Ctor.RecName)
	for i, x := range v.Data {
		fmt.Fprintf(&b, "%v:%v,", v.Ctor.Flist[i], x)
	}
	s := b.Bytes()
	s[len(s)-1] = '}'
	return string(s)
}

//  VRecord.Type returns the underlying constructor
func (v *VRecord) Type() IRank {
	return v.Ctor
}

//  VRecord.Copy returns a distinct copy of itself
func (v *VRecord) Copy() Value {
	r := &VRecord{v.Ctor, make([]Value, len(v.Data))}
	copy(r.Data, v.Data)
	return r
}

//  VRecord.Before compares two records for sorting on field i
func (a *VRecord) Before(x Value, i int) bool {
	b := x.(*VRecord)
	if a.Ctor != b.Ctor {
		// different record types; order by type name
		return a.Ctor.RecName < b.Ctor.RecName
	}
	if i >= 0 && len(a.Data) > i && len(b.Data) > i {
		// both sides have an item i
		return LT(a.Data[i], b.Data[i], -1)
	} else {
		// put missing one first; otherwise #%#% we don't care
		return len(a.Data) < len(b.Data)
	}
}

//  VRecord.Import returns itself
func (v *VRecord) Import() Value {
	return v
}

//  VRecord.Export returns itself
func (v *VRecord) Export() interface{} {
	return v
}

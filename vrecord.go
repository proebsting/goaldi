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
	return v.Ctor.Name + "{}"
}

//  VRecord.GoString -- returns string for image() and printf("%#v")
func (v *VRecord) GoString() string {
	if len(v.Data) == 0 {
		return v.Ctor.Name + "{}"
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s{", v.Ctor.Name)
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

//  VRecord.Import returns itself
func (v *VRecord) Import() Value {
	return v
}

//  VRecord.Export returns itself
func (v *VRecord) Export() interface{} {
	return v
}

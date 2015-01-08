//  vstruct.go -- a user-defined (usually) structure

package goaldi

import (
	"bytes"
	"fmt"
)

type VStruct struct {
	Defn *VDefn  // underlying struct definition
	Data []Value // current data values
}

//  VStruct.String -- conversion to Go string returns "name{}"
func (v *VStruct) String() string {
	return v.Defn.Name + "{}"
}

//  VStruct.GoString -- returns string for image() and printf("%#v")
func (v *VStruct) GoString() string {
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

//  VStruct.Rank returns rStruct
func (v *VStruct) Rank() int {
	return rStruct
}

//  VStruct.Type returns the defined struct name
func (v *VStruct) Type() Value {
	return NewString(v.Defn.Name)
}

//  VStruct.Copy returns a distinct copy of itself
func (v *VStruct) Copy() Value {
	r := &VStruct{v.Defn, make([]Value, len(v.Data))}
	copy(r.Data, v.Data)
	return r
}

//  VStruct.Import returns itself
func (v *VStruct) Import() Value {
	return v
}

//  VStruct.Export returns itself
func (v *VStruct) Export() interface{} {
	return v
}

//  json.go -- JSON manipulation

package main

import (
	"fmt"
	"reflect"
)

//  jdump -- print contents of generic JSON data tree
//  (does not recurse into arrays inside typed structs)
func jdump(jtree interface{}) {
	fmt.Printf("JSON data:")
	jdu("", jtree)
	fmt.Printf("\n")
}

func jdu(indent string, jtree interface{}) {
	switch x := jtree.(type) {
	case []interface{}:
		for _, v := range x {
			jdu("   "+indent, v)
			fmt.Printf("\n%s-----------------------------",
				indent)
		}
	case map[string]interface{}:
		for k, v := range x {
			fmt.Printf("\n%s%v: ", indent, k)
			jdu("   "+indent, v)
		}
	default:
		fmt.Printf("%v", x)
	}
}

//  jfix -- replace maps by IR structs in Json tree
func jfix(jtree interface{}) interface{} {
	switch x := jtree.(type) {
	case []interface{}:
		for i, v := range x {
			x[i] = jfix(v)
		}
		return x
	case map[string]interface{}:
		for k, v := range x {
			x[k] = jfix(v)
		}
		return structFor(x)
	default:
		return jtree
	}
}

//  initialize IR mapping table
func init() {
	for _, ir := range irlist {
		irtable[ir.name] = reflect.TypeOf(ir.example)
	}
}

var irtable = make(map[string]reflect.Type)

//  structFor -- return IR struct equivalent to map
func structFor(m map[string]interface{}) interface{} {
	tag := m["tag"].(string)
	if tag == "" {
		panic(fmt.Sprintf("no tag in %v", m))
	}
	rtype := irtable[tag]
	if rtype == nil {
		panic(fmt.Sprintf("unrecognized IR tag %s", tag))
	}
	resultp := reflect.New(rtype)
	result := resultp.Elem()
	for key, val := range m {
		key = Capitalize(key)
		f := result.FieldByName(key)
		if f.CanSet() {
			if val != nil {
				f.Set(reflect.ValueOf(val))
			}
		} else if key != "Tag" {
			panic("cannot set key " + key)
		}
	}
	return result.Interface()
}

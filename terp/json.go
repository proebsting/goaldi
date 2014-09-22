//  json.go -- JSON manipulation

package main

import (
	"fmt"
	"reflect"
)

//  jdump -- print contents of generic JSON data tree
//  (does not recurse into arrays inside typed structs)
func jdump(jtree interface{}) {
	tally := make(map[string]int)
	fmt.Printf("JSON data:")
	jdu("", jtree, tally)
	fmt.Printf("\n\nStruct field types:\n")
	for k, v := range tally {
		fmt.Printf("field %-45s %3d\n", k, v)
	}
}

func jdu(indent string, jtree interface{}, tally map[string]int) {
	switch x := jtree.(type) {
	case []interface{}:
		for _, v := range x {
			fmt.Printf("\n%s----------------------------- ",
				indent)
			jdu("   "+indent, v, tally)
		}
	case map[string]interface{}:
		for k, v := range x {
			fmt.Printf("\n%s%v: ", indent, k)
			jdu("   "+indent, v, tally)
			if submap, ok := v.(map[string]interface{}); ok {
				tx := fmt.Sprintf("%15s.%-12s %s",
					x["tag"], k, submap["tag"])
				tally[tx]++
			} else if k != "tag" {
				tx := fmt.Sprintf("%15s.%-12s %T",
					x["tag"], k, v)
				tally[tx]++
			}
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
		t := reflect.TypeOf(ir)
		irtable[t.Name()] = t
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
		setField(f, key, val)
	}
	return result.Interface()
}

//  setField -- set field in struct
//#%#%  can't handle destination slices other than []string and []interface{}
func setField(f reflect.Value, key string, val interface{}) {
	if key == "Tag" || val == nil {
		return // nothing to do
	}
	if !f.CanSet() {
		panic("cannot set key " + key)
	}
	t := f.Type()
	if t.Kind() != reflect.Slice || t.Elem().Kind() == reflect.Interface {
		f.Set(reflect.ValueOf(val))
		return
	}
	// we have to make a typed slice and copy in the elements
	resultp := reflect.New(t)
	//#%#% fmt.Printf("%T %v\n", resultp, resultp)
	result := *(resultp.Interface().(*[]string))
	//#%#% fmt.Printf("%T %v\n", result, result)
	for _, v := range val.([]interface{}) {
		//#%#% fmt.Printf("%T %v\n", v, v)
		result = append(result, v.(string))
	}
	f.Set(reflect.ValueOf(result))
}

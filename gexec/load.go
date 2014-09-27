//  load.go -- read intermediate representation from JSON file

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

//  load -- read a single JSON-encoded IR file as a tree of objects
func load(fname string) []interface{} {

	babble("loading %s", fname)

	//  open the file
	var gfile *os.File
	var err error
	if fname == "-" {
		gfile = os.Stdin
	} else {
		gfile, err = os.Open(fname)
		checkError(err)
	}
	gcode := bufio.NewReader(gfile)

	//  skip initial comment lines (e.g. #!/usr/bin/env gdx...)
	for {
		b, e := gcode.Peek(1)
		if e != nil || b[0] != '#' {
			break
		}
		gcode.ReadBytes('\n')
	}

	//  load the JSON-encoded program
	jd := json.NewDecoder(gcode)
	var jtree []interface{}
	checkError(jd.Decode(&jtree))
	if opt_jdump {
		jdump(jtree)
	}
	jtree = jstructs(jtree).([]interface{})
	if opt_adump {
		dumptree("", jtree)
		fmt.Println()
	}
	return jtree
}

//  dumptree -- print a human-readable listing of the IR
func dumptree(indent string, x interface{}) {
	switch t := x.(type) {
	case nil:
		return
	case []interface{}:
		for _, v := range t {
			dumptree(indent, v)
		}
	case []ir_chunk:
		for _, v := range t {
			dumptree(indent, v)
		}
	case ir_Function:
		fmt.Printf("\n%sproc %s  %v  start %v\n",
			indent, t.Name, t.Coord,
			t.CodeStart.Value)
		fmt.Printf("%s   param %v\n", indent, t.ParamList)
		fmt.Printf("%s   local %v\n", indent, t.LocalList)
		fmt.Printf("%s   static %v\n", indent, t.StaticList)
		dumptree(indent, t.CodeList)
	case ir_chunk:
		fmt.Printf("%s%s:\n", indent, t.Label.Value)
		dumptree(indent+"   ", t.InsnList)
	default:
		fmt.Printf("%s%T %v\n", indent, x, x)
	}
}

//  jdump -- print the contents of a generic JSON data tree
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

//  jstructs -- replace maps by IR structs in Json tree
func jstructs(jtree interface{}) interface{} {
	switch x := jtree.(type) {
	case []interface{}:
		for i, v := range x {
			x[i] = jstructs(v)
		}
		return x
	case map[string]interface{}:
		for k, v := range x {
			x[k] = jstructs(v)
		}
		return structFor(x)
	default:
		return jtree
	}
}

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
	result := resultp.Elem()
	for _, v := range val.([]interface{}) {
		result = reflect.Append(result, reflect.ValueOf(v))
	}
	f.Set(result)
}

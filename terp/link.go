//  link.go -- loading and linking

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

//  load reads a single JSON-encoded IR file.
func load(fname string) UNKNOWN {

	babble("loading %s", fname)

	//  open the file
	gfile, err := os.Open(fname)
	checkError(err)
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
	var jtree interface{}
	jd.Decode(&jtree)
	// jdump(jtree) //#%#%#%
	jtree = jfix(jtree)
	dumptree("", jtree) //#%#%#%
	return jtree
}

//  link combines IR files to make a complete program.
func link(parts []UNKNOWN) UNKNOWN {
	babble("linking")
	return nil
}

//  dumptree prints a human-readable listing of the IR
func dumptree(indent string, x interface{}) {
	switch t := x.(type) {
	case nil:
		return
	case []interface{}:
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

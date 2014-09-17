//  link.go -- loading and linking

package main

import (
	"bufio"
	"encoding/json"
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

	_ = jd
	return nil
}

//  link combines IR files to make a complete program.
func link(parts []UNKNOWN) UNKNOWN {
	babble("linking")
	return nil
}

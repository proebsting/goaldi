#!/usr/bin/env goaldi
#
#  libdoc.gd efile godoc -- extract Goaldi library documentation
#
#  efile is a file produced by "goaldi -E"
#  godoc is a file produced by "go doc -all pkg" for all needed packages
#
#  Parsing assume relatively simple types like those we are actually using.
#  This probably can't handle double indirection or multiple parenthesis levels.

global exclude := [				# funcname patterns to exclude from list
	"goaldi/extensions",
	"hash/",
	"archive/zip",
]

global funcs := table()			# function documentation indexed by name

global curpkg := "unknown."		# current package name
global curlist := []			# current entry being accumulated


#  main logic

procedure main(ename, gname) {
	local e := file(ename)
	local g := file(gname)
	loaddoc(g)
	gendoc(e)
}


#  loaddoc(g) -- load function documentation from the "godoc" output in file g

procedure loaddoc(g) {
	while local line := read(g) do {
		if *line = 0 | line[1] == !" \t" then {
			# this is a continuation of the existing entry
			curlist.put(line)
		} else {
			# this is a new entry; check it out
			local f := fields(line)
			case f[1] of {
				"package": {	# package: just remember the name
					curlist := []
					if !!contains(line, " // import ") then {
						# not a false hit; remember the name
						curpkg := f[2] || "."
					}
				}
				"func": {		# func: register and start accumulating lines
					newfunc(line)
				}
				default: {		# anything else: ignore it
					curlist := []	# anonymous bitbucket for subsequent lines
				}
			}
		}
	}
}


#  newfunc(line) -- parse a line of the form "func [(type rcvr)] name(..."

procedure newfunc(line) {

	# break apart the line using a regular expression
	/static rx := regex(`^func (\([^)]*\) )?([^(]+)`)
	local frags := rx.FindStringSubmatch(line)
	local rcvr := frags[2]	# receiver, if a method
	local name := frags[3]	# function name

	if rcvr ~== "" then {	# if a method
		# strip out just type name and then add a dot
		rcvr := split(rcvr, " ")[2][1:-1]
		if rcvr[1] == "*" then rcvr := rcvr[2:0]
		rcvr := rcvr || "."
	}

	# register the function and start accumulating lines
	local key := curpkg || rcvr || split(name, "(")[1]
	funcs[key] := curlist := [line]
	return
}


#  gendoc(e) -- extract functions listed in "goaldi -E" output file e

procedure gendoc(e) {
	/static prx := regex(`(.*)\(\*(.*)\)(.*)`)	# pointer receiver pattern

	while local line := read(e) do {
		if !!contains(line, " -- ") then {

			# this is a line specifying a stdlib procedure or method
			local words := split(trim(line, " "), " ")
			local fspec := words[-1]
			local descr := trim(line[1:-*fspec], " ")
			local func := split(fspec, "/")[-1]

			if !!contains(fspec, !exclude) then {
				continue
			}

			if descr[2] == " " then {
				descr := descr[3:0]			# remove constructor result char
			}

			if fspec[1+:14] == "goaldi/runtime" then {
				fspec[1+:14] := "goaldi"	# simplify
			}
			if \(local x := prx.FindStringSubmatch(func)) then {
				# change "goaldi.(*VFile).FRead" to "goaldi.VFile.FRead"
				func := x[2] || x[3] || x[4]
			}

			local doc := funcs[func]
			if \doc then {
				show(descr, fspec, doc)
			} else {
				%stderr.write("missing ", func, "  (needed for ", words[1], ")")
			}
		}
	}
}


#  show(descr, fspec, doc) -- generate documentation for one function
#
#  Each call produces one entry in an ASCIIdoc "labeled list"

procedure show(descr, fspec, doc) {

	# trim trailing blank lines from godoc output
	while doc[-1] == "" do {
		doc.pull()
	}

	# build a boilerplate header from available parts
	write()
	writes(descr)
	if fspec[1+:6] ~== "goaldi" then {
		^w := split(fspec, ".")
		^u := "http://golang.org/pkg/" || w[1] || "#" || w[2]
		writes(" [silver]_(", u, "[", fspec, "])_")
	}
	write("::")

	# skip godoc first line (in favor of descr just written) and copy the rest
	every ^s := doc[2 to *doc] do {
		if *s = 0 then {
			write("+")
		} else if s[1:5] == "    " then {
			write(s[5:0])
		} else {
			write(s)
		}
	}
}

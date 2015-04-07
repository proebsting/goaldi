#	main.gd -- main program for Goaldi front-end translator
#
#   This program is run (interpreted) by the Go main program
#   if its first argument is not "-x".

global USAGE := "goaldi [options] file.gd... [--] [arg...]"

record optf (flag, meaning)
global optlist := [
	optf("-c", "compile only, IR code to file.gir"),
	optf("-a", "compile only, IR code to file.gir, assembly to file.gia"),
	optf("-l", "load and link but do not execute"),
	optf("-t", "show CPU timings"),
	optf("-v", "issue verbose commentary"),
	optf("-D", "dump Go stack on panic"),
	optf("-E", "show initial environment"),
	optf("-G", "compile to file.go (EXPERIMENTAL)"),
	optf("-N", "inhibit optimization"),
	optf("-P", "produce ./PROFILE file (Linux)"),
	optf("-T", "trace IR instruction execution"),
]
global gxopts := "ltvDEPT"	# options passed to goaldi interpreter


#  main program -- see code above for usage 

procedure main(args[]) {

	#  process options
	^opts := getopts(args, optlist)
	^gxargs := []
	every ^c := !gxopts do {
		if \opts[c] then gxargs.put("-" || c)
	}

	#  source files are the first file, always,
	#  plus any following that end in ".gd"
	^srclist := [@args] | usage()
	while args[1][-3:0] == ".gd" do {
		srclist.put(@args)
	}
	if args[1] == "--" then {
		@args					# discard "--" separator argument
	}

	#  translate source files to IR code
	every ^iname := !srclist do {
		^ibase := if iname[-3:0]==".gd" then iname[1:-3] else iname
		^oname := ibase || ".gir"
		if \opts["G"] then {
			oname := ibase || ".go"
		} else if /opts["a"] & /opts["c"] then {
			# need a temporary name
			#%#% need a better way to do this, and to ensure deletion
			/static ntemps := 0
			oname := "/tmp/gd-" || getpid() || "-" || ?10000 || "-" ||
				(ntemps +:= 1) || ".gir"
		}
		translate(iname, oname, opts)
		if \opts["a"] then {
			with %stdout := file(ibase || ".gia", "w") do {
				gexec(["-l", "-A", oname])
				%stdout.close()
			}
		}
		gxargs.put(oname)
	}
	if \opts["a"] | \opts["c"] | \opts["G"] then {
		return
	}

	#  execute the translated files (already put on gxargs list)
	gxargs.put("--")
	every gxargs.put(!args)
	gexec(gxargs)

}


#  translate(iname, oname, opts) -- translate one file

procedure translate(iname, oname, opts) {
	^ofile := file(oname, "w")
	^pipeline := create !file(iname)
	pipeline := create lex(pipeline, iname)
	pipeline := create parse(pipeline)
	pipeline := create ast2ir(pipeline)
	if /opts["N"] then {
		pipeline := create optim(pipeline, ["-O"])
	}
	if \opts["G"] then {
		go_File(ofile, pipeline)
	} else {
		json_File(ofile, pipeline)
	}
}


#  gexec(arglist) -- run goaldi process with given arglist (preceded by -x)

procedure gexec(arglist) {
	arglist.push("-x")
	arglist.push("goaldi")
	^c := command ! arglist
	c.Stdin := osfile(0)
	c.Stdout := %stdout
	c.Stderr := %stderr
	^r := c.Run()
	if \r then throw(r)		#%#% later make this nicer
}


#  getopts(args,optlist) -- simplified command option processing
#
#  Processes only one-character flag options; does not handle values.
#  Removes option arguments from args and returns a table of flags.
#  Aborts on error.

procedure getopts(args, optlist) {
	^allowed := set()
	every ^o := !optlist do {
		allowed.insert(o.flag[2])
	}
	^seen := table()
	while args[1][1] == "-" do {
		^flags := args.get()[2:0]
		if flags == "-" then {
			break	# exit on "--"
		}
		every ^c := !flags do {
			if allowed[c] then {
				seen[c] := c
			} else {
				%stderr.write("unrecognized option: -", c)
				usage()
			}
		}
	}
	return seen
}


#  usage() -- write usage message, list legal options, abort
procedure usage() {
	%stderr.write("usage: ", USAGE)
	every ^o := !optlist do {
		%stderr.write("  ", o.flag, "  ", o.meaning)
	}
	exit(1)
}

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
	optf("-A", "dump assembly listing to stdout before execution"),
	optf("-D", "dump Go stack on panic"),
	optf("-E", "show initial environment"),
	optf("-G", "compile to file.go (SECRET)"),
	optf("-I", "trace initialization ordering"),
	optf("-N", "inhibit optimization"),
	optf("-P", "produce ./PROFILE file (Linux)"),
	optf("-T", "trace IR instruction execution"),
]
global gxopts := "ltADEIPT"	# options passed to goaldi interpreter


#  main program -- see code above for usage 

procedure main(args[]) {
	^t0 := cputime()
	randomize()					# for irreproducible temp file names

	#  process options
	^opts := getopts(args, optlist)
	^gxargs := []
	every ^c := !gxopts do {
		if \opts[c] then gxargs.put("-" || c)
	}
	if /opts["c"] & /opts["a"] & /opts["G"] then {
		gxargs.put("-#")		# delete temp files after loading
	}
	if \opts["t"] then {
		fprintf(%stderr, "%7.3f startup\n", t0)
	}

	#  source files are the first file, always,
	#  plus any following arguments that end in ".gd"
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
			/static ntemps := 0
			oname := "/tmp/gd-" || getpid() || "-" || ?10000 || "-" ||
				(ntemps +:= 1) || ".gir"
		}
		translate(iname, oname, opts)
		if \opts["a"] then {
			^out := file(ibase || ".gia", "w")
			gexec(["-l", "-A", oname], out)
			out.close()
		}
		gxargs.put(oname)
		if \opts["t"] then {
			^t := cputime()
			fprintf(%stderr, "%7.3f translation (%s)\n", t - t0, iname)
			t0 := t
		}
	}
	if \opts["a"] | \opts["c"] | \opts["G"] then {
		return
	}

	#  execute the translated files (already put on gxargs list)
	gxargs.put("--")		# end of arguments to interpreter
	every gxargs.put(!args)	# program arguments
	exit(gexec(gxargs))		# run program and exit with its exit code
}


#  translate(iname, oname, opts) -- translate one file

procedure translate(iname, oname, opts) {
	^ofile := file(oname, "w")
	^ifile := file(iname, "f") | stop(osargs()[1], ": Cannot open: ", iname)
	^pipeline := create !ifile
	pipeline := (create lex(pipeline, iname)).buffer(1000)
		# buffer size 1000 gives perhaps 1% speedup. bufsize 100 doesn't help.
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


#  gexec(arglist,stdout) -- run goaldi process (with args preceded by -x)
#
#  returns the exit status.

procedure gexec(arglist, stdout) {
	%stdout.flush()
	%stderr.flush()
	arglist.push("-x")
	arglist.push(osargs()[1])
	^c := command ! arglist
	c.Stdin := osfile(0)
	c.Stdout := \stdout | osfile(1)
	c.Stderr := osfile(2)
	c.Run()
	return c.ProcessState.Sys().ExitStatus()
}


#  getopts(args,optlist) -- simplified command option processing
#
#  Processes only one-character flag options; does not handle values.
#  Removes option arguments from args and returns a table of flags.
#  Aborts on error.

procedure getopts(args, optlist) {
	^allowed := set()
	every ^o := !optlist do {
		allowed.put(o.flag[2])
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
				%stderr.write("Unrecognized option: -", c)
				usage()
			}
		}
	}
	return seen
}


#  usage() -- write usage message, list legal options, abort
procedure usage() {
	%stderr.write("Usage: ", USAGE)
	every ^o := !optlist do {
		if not !!contains(o.meaning, "SECRET") then {
			%stderr.write("  ", o.flag, "  ", o.meaning)
		}
	}
	exit(1)
}

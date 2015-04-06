#	gtran.gd -- main program for Goaldi front-end translator
#
#	usage:  gtran [-N] [-G] file...
#
#	-N		inhibit optimization
#	-G		generate Go code in file.go instead of JSON to stdout

procedure main(args[]) {
	local opts := options(args, "NG")
	every translate(!args, opts)
}

#  options(args,optstring) -- simplified command option processing
#
#  Interprets any char of optstring as an allowable flag argument.
#  Does not allow value arguments.
#  Removes option arguments from args and returns a table of flags.
#  Aborts on error.

procedure options(args, optstring) {
	every (^allowed := set()).insert(!optstring)
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
				stop("unrecognized option: -", c)
			}
		}
	}
	return seen
}

#  translate(fname, opts) -- translate one file

procedure translate(fname, opts) {
	^fbase := if fname[-3:0]==".gd" then fname[1:-3] else fname
	^pipeline := create !file(fname)
	pipeline := create lex(pipeline, fname)
	pipeline := create parse(pipeline)
	pipeline := create ast2ir(pipeline)
	if /opts["N"] then {
		pipeline := create optim(pipeline, ["-O"])
	}
	if \opts["G"] then {
		go_File(file(fbase || ".go", "w"), pipeline)
	} else {
		json_File(%stdout, pipeline)
	}
}

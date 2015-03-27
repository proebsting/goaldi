#  gtran.gd -- main program for Goaldi front-end translator
#  (quick-and-dirty experimental version)

procedure main(args[]) {
	local opts := options(args, "N")
	every ^fname := !args do {
		^fbase := if fname[-3:0]==".gd" then fname[1:-3] else fname
		^pipeline := create !open(fbase || ".gd")
		pipeline := create lex(pipeline, fbase || ".gd")
		pipeline := create parse(pipeline)
		pipeline := create ast2ir(pipeline)
		if /opts["N"] then {
			pipeline := create optim(pipeline, ["-O"])
		}
		pipeline := create json_File(pipeline)
		pipeline := create stdout(pipeline)
		@pipeline	# wait for processes to finish and close
	}
}

#  pipeline component to copy its contents to a file
procedure tee(src, fname) {
	^f := open(fname, "w")
	while ^v := @src do {
		if type(v) === string then {
			f.write(v)
		} else {
			f.write(v.image())
		}
		suspend v
	}
	f.close()
}

#  (terminal) pipeline component to write stream to stdout
procedure stdout(src) {
	while write(@src)
}

#  (terminal) pipeline component to toss everything into a black hole
procedure sink(src) {
	while @src
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

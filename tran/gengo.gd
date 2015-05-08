#  gengo.gd -- generate Go code from IR.  (INCOMPLETE, EXPERIMENTAL)

#	Conventions:
#	user symbols map to the form _symbol
#	generated global/initial procedure names end up prefixed by __
#	non-prefixed names are reserved for use by generated code
#	(e.g. env, args, frame, p_xxx etc)

procedure go_File(f, irgen) {
	local init
	while local p := @irgen do {
		/init := go_start(f, p)
		p.go(f)
	}
}

procedure go_start(f, p) {
	^namespace := ("" ~== \p.namespace) | "main"
	f.write(`package `, namespace)
	f.write(`import g "goaldi/runtime"`)
	return namespace
}

procedure ir_Global.go(f) {
	f.write()
	f.writes(`var _`, self.name, ` g.Value = g.NewVariable(`)
	if \self.fn then {
		#%#% not quite right: assumes that initialization function
		#%#% returns a value rather than implementing an assignment
		f.write(`g.ResultOf(p`, goname(self.fn), `))`)
	} else {
		f.write(`g.NilValue)`)
	}
}

procedure ir_Record.go(f) {
	^flist := "[]string{"
	^parent := "nil"
	every flist ||:= image(!self.fieldList) || ","
	f.write()
	f.write(`var _`, self.name, ` g.Value = g.NewCtor("`,
		self.name, `", `, parent, `, `, flist, `})`)
}

procedure ir_Initial.go(f) {
	f.write(`func init() { g.Run(p`, goname(self.fn), `, []g.Value{}) }`)
}

procedure ir_Function.go(f) {
	^gname := goname(self.name)

	# generate global symbol
	f.write()
	^plist := "&[]string{"
	every plist ||:= image(!self.paramList) || ","
	f.write(`var `, gname, ` g.Value = g.NewProcedure("`,
		self.name, `", `, "\n\t", plist, "},\n\t",
		if self.accumulate ~=== "" then "true" else "false",
		`, p`, gname, `, p`, gname, `, "")`)

	# generate prologue for implementing function
	f.write()
	f.write(`func p`, gname,
		`(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {`)

	# generate code
	every ^c := !self.codeList do {
		c.go(f)
	}

	f.write(`	return nil, nil`)
	f.write(`}`)
}

procedure ir_chunk.go(f) {
	every ^c := !self.insnList do {
		c.go(f)
	}
}

procedure ir_Fail.go(f) {
}

procedure ir_Call.go(f) {
}

procedure ir_OpFunction.go(f) {
}

procedure ir_Var.go(f) {
}

procedure ir_RealLit.go(f) {
}

procedure ir_EnterScope.go(f) {
}

#	turn a user or generated symbol into the version used in generated Go code:
#   prefix by "_" and replace all "$" by "_"
procedure goname(s) {
    ^t := "_"
	every ^c := !s do
		t ||:= if c == "$" then "_" else c
	return t
}

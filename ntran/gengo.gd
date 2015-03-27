#  gengo.gd -- generate Go code from IR.  (INCOMPLETE, EXPERIMENTAL)

procedure go_File(f, irgen) {
	local init
	while local p := @irgen do {
		/init := go_start(f, p)
		case type(p) of {
			default:		throw("unrecognized type", p)
			ir_Global:		suspend go_global(f, p)
			ir_Record:		suspend go_record(f, p)
			ir_Initial:		suspend go_initial(f, p)
			ir_Function:	suspend go_function(f, p)
		}
	}
}

procedure go_start(f, p) {
	^namespace := ("" ~== \p.namespace) | "main" 
	f.write(`package `, namespace)
	f.write(`import g "goaldi"`)
	return namespace
}

procedure go_global(f, p) {
	f.write(`var `, p.name, ` = g.NewVariable(g.NilValue)`)
}

procedure go_record(f, p) {
	^flist := "[]string{"
	^parent := "nil"
	every flist ||:= image(!p.fieldList) || ","
	f.write(`var `, p.name, ` = g.NewCtor("`,
		p.name, `", `, parent, `, `, flist, `})`)
}

procedure go_initial(f, p) {
}

procedure go_function(f, p) {
}

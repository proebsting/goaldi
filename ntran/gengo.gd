#  gengo.gd -- generate Go code from IR.  (INCOMPLETE, EXPERIMENTAL)

procedure go_File(irgen) {
	local init
	while local p := @irgen do {
		suspend /init := go_start(p)
		case type(p) of {
			default:		throw("unrecognized type", p)
			ir_Global:		suspend go_global(p)
			ir_Record:		suspend go_record(p)
			ir_Initial:		suspend go_initial(p)
			ir_Function:	suspend go_function(p)
		}
	}
}

procedure go_start(p) {
	^namespace := ("" ~== \p.namespace) | "main" 
	suspend "package " || namespace
	suspend `import g "goaldi"`
}

procedure go_global(p) {
	suspend "var " || p.name || " = g.NewVariable(g.NilValue)"
}

procedure go_record(p) {
	^flist := "[]string{"
	^parent := "nil"
	every flist ||:= image(!p.fieldList) || ","
	suspend "var " || p.name || ` = g.NewCtor("` || p.name || `", ` ||
		parent || ", " || flist || "})"
}

procedure go_initial(p) {
}

procedure go_function(p) {
}

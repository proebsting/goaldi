#  gengo.gd -- generate Go code from IR.  (INCOMPLETE, EXPERIMENTAL)

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
	f.write(`import g "goaldi"`)
	return namespace
}

procedure ir_Global.go(f) {
	f.write(`var `, self.name, ` = g.NewVariable(g.NilValue)`)
}

procedure ir_Record.go(f) {
	^flist := "[]string{"
	^parent := "nil"
	every flist ||:= image(!self.fieldList) || ","
	f.write(`var `, self.name, ` = g.NewCtor("`,
		self.name, `", `, parent, `, `, flist, `})`)
}

procedure ir_Initial.go(f) {
	throw(self)
}

procedure ir_Function.go(f) {
	every ^c := !self.codeList do {
		c.go(f)
	}
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

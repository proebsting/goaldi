#  irgen.gd -- translate abstract syntax trees to intermediate representation.

global ir_deferred
global ir_declare_set
global ir_namespace

record ir_info(start, resume, failure, success, x)
record ir_loopinfo(nextlabel, continueTmp, value, bounded, rval, indirects)

record ir_stacks(tmp, current_proc, createflag, declare_set, loop_stack, localSet, staticSet, globalSet, syms, parent)

record ir_scope(parent, Static, Dynamic)

procedure ir_a_Paired(p, st, target, bounded, rval) {
	local tmp
	local lhs
	local rhs
	local L
	local R
	local Ltmp
	local i
	local c

	ir_init(p)

	c := p.coord

	lhs := ir_tmp(st)
	rhs := ir_tmp(st)
	Ltmp := ir_tmp(st)
	tmp := ir_tmp(st)

	suspend ir(p.fn, st, tmp, bounded, nil)

	suspend ir_chunk(p.ir.start, [ ir_Goto(c, p.fn.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(c, p.fn.ir.resume) ])
	suspend ir_chunk(p.fn.ir.failure, [ ir_Goto(c, p.ir.failure) ])


	every i := *p.leftList to 1 by -1 do {
		L := p.leftList[i]
		R := p.rightList[i]
		suspend ir(L, st, lhs, nil, "rval")
		suspend ir(R, st, rhs, nil, "rval")
		suspend ir_chunk(L.ir.success, [
			ir_opfn(p.coord, Ltmp, nil, ir_operator("[]", 2, nil), [ tmp, lhs ], L.ir.resume),
			ir_Goto(c, R.ir.start),
			])
		suspend ir_chunk(L.ir.failure, [ ir_Goto(c, p.leftList[i+1].ir.start) ])
		suspend ir_chunk(R.ir.failure, [ ir_Goto(c, L.ir.resume) ])
		suspend ir_chunk(R.ir.success, [
			ir_opfn(p.coord, nil, nil, ir_operator(":=", 2, "rval"), [ Ltmp, rhs ], R.ir.resume),
			ir_Goto(c, R.ir.resume),
			])
	}
	if *p.leftList > 0 then {
		suspend ir_chunk(p.leftList[-1].ir.failure, [
			ir_Move(c, target, tmp),
			ir_Goto(c, p.ir.success),
			])
		suspend ir_chunk(p.fn.ir.success, [ ir_Goto(c, p.leftList[1].ir.start) ])
	} else {
		suspend ir_chunk(p.fn.ir.success, [
			ir_Move(c, target, tmp),
			ir_Goto(c, p.ir.success),
			])
	}
}

# record a_NoOp( )
procedure ir_a_NoOp(p, st, target, bounded, rval) {
	local c

	ir_init(p)

	c := (\p.coord | "0:0")
	suspend ir_chunk(p.ir.start, [ ir_Goto(c, p.ir.success) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(c, p.ir.failure) ])
}

# record a_Field( expr field )
procedure ir_a_Field(p, st, target, bounded, rval) {
	local t
	ir_init(p)
	t := ir_value(p, st, target)
	suspend ir(p.expr, st, t, nil, rval)

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])
	suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.expr.ir.resume) ])
	suspend ir_chunk(p.expr.ir.success, [
		ir_Field(p.coord, target, t, p.field.id, rval),
		ir_Goto(p.coord, p.ir.success),
		])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Limitation( expr limit )
procedure ir_a_Limitation(p, st, target, bounded, rval) {
	local c
	local t
	local one

	ir_init(p)
	c := ir_tmp(st)
	t := ir_tmp(st)
	one := ir_tmp(st)

	suspend ir(p.limit, st, t, nil, "rval")
	suspend ir(p.expr, st, target, bounded, rval)

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.limit.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [
		ir_opfn(p.coord, c, nil, ir_operator(">", 2, "rval"), [ t, c ], p.limit.ir.resume),
		ir_IntLit(p.coord, one, 1),
		ir_opfn(p.coord, c, nil, ir_operator("+", 2, "rval"), [ c, one ], p.expr.ir.resume),
		ir_Goto(p.coord, p.expr.ir.resume),
		])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.limit.ir.resume) ])
	suspend ir_chunk(p.limit.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.expr.ir.success, [ ir_Goto(p.coord, p.ir.success) ])
	suspend ir_chunk(p.limit.ir.success, [
		ir_opfn(p.coord, t, nil, ir_operator("#", 1, "rval"), [ t ], p.limit.ir.resume),
		ir_IntLit(p.coord, c, 1),
		ir_Goto(p.coord, p.expr.ir.start),
		])
}

# record a_Not( expr )
procedure ir_a_Not(p, st, target, bounded, rval) {
	ir_init(p)

	suspend ir(p.expr, st, nil, "always bounded", "rval")

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.expr.ir.success, [
		ir_Goto(p.coord, p.ir.failure),
		])
	suspend ir_chunk(p.expr.ir.failure, [
		ir_NilLit(p.coord, target),
		ir_Goto(p.coord, p.ir.success),
		])
}

# record a_Alt( eList )
procedure ir_a_Alt(p, st, target, bounded, rval) {
	local t
	local tmpst
	local i
	local indirects

	indirects := []

	ir_init(p)
	/bounded & (t := ir_tmploc(st))

	every i := 1 to *p.eList do {
		suspend ir(p.eList[i], st, target, bounded, rval)
	}

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.eList[1].ir.start) ])

	every i := 1 to *p.eList do {
		if /bounded then {
			indirects.put(p.eList[i].ir.resume)
			suspend ir_chunk(p.eList[i].ir.success, [
				ir_MoveLabel(p.coord, t, p.eList[i].ir.resume),
				ir_Goto(p.coord, p.ir.success),
				])
		} else {
			suspend ir_chunk(p.eList[i].ir.success, [
				ir_Goto(p.coord, p.ir.success),
				])
		}
		suspend ir_chunk(p.eList[i].ir.failure,[ir_Goto(p.coord, p.eList[i+1].ir.start)])
	}
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, t, indirects) ])
	suspend ir_chunk(p.eList[-1].ir.failure, [ ir_Goto(p.coord, p.ir.failure)])
}

# record a_ExcAlt( eList )
procedure ir_a_ExcAlt(p, st, target, bounded, rval) {
	local i
	local t
	local R
	local F
	local ch
	local oldt
	local indirectsR
	local indirectsF

	indirectsR := []

	ir_init(p)
	/bounded & (R := ir_tmploc(st))
	F := ir_tmploc(st)

	every i := 1 to *p.eList do {
		suspend ir(p.eList[i], st, target, bounded, rval)
	}


	t := nil
	every i := *p.eList to 1 by -1 do {
		if /t then {
			suspend ir_chunk(p.eList[i].ir.failure,[ir_Goto(p.coord, p.ir.failure)])
		} else {
			suspend ir_chunk(p.eList[i].ir.failure,[ir_IndirectGoto(p.coord, F, [t, p.ir.failure])])
		}

		indirectsR.put(p.eList[i].ir.resume)
		oldt := t
		t := ir_label(p.eList[i], "prefix")
		ch := ir_chunk(t, [
			ir_MoveLabel(p.coord, F, \oldt | p.ir.failure),
			ir_MoveLabel(p.coord, R, p.eList[i].ir.resume),
			ir_Goto(p.coord, p.eList[i].ir.start),
			])
		suspend ch
		if /bounded then {
			suspend ir_chunk(p.eList[i].ir.success, [
				ir_MoveLabel(p.coord, F, p.ir.failure),
				ir_Goto(p.coord, p.ir.success),
				])
		} else {
			suspend ir_chunk(p.eList[i].ir.success, [
				ir_Goto(p.coord, p.ir.success),
				])
		}
	}

	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, R, indirectsR) ])
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, t) ])	 # subtle t from previous loop
}

# record a_RepAlt( expr )
procedure ir_a_RepAlt(p, st, target, bounded, rval) {
	local t

	ir_init(p)
	/bounded & (t := ir_tmploc(st))
	suspend ir(p.expr, st, target, bounded, rval)

	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.expr.ir.resume) ])
	if /bounded then {
		suspend ir_chunk(p.ir.start, [
			ir_MoveLabel(p.coord, t, p.ir.failure),
			ir_Goto(p.coord, p.expr.ir.start),
			])
		suspend ir_chunk(p.expr.ir.success, [
			ir_MoveLabel(p.coord, t, p.ir.start),
			ir_Goto(p.coord, p.ir.success),
			])
		suspend ir_chunk(p.expr.ir.failure, [
			ir_IndirectGoto(p.coord, t, [p.ir.failure, p.ir.start]),
			])
	} else {
		suspend ir_chunk(p.ir.start, [
			ir_Goto(p.coord, p.expr.ir.start),
			])
		suspend ir_chunk(p.expr.ir.success, [
			ir_Goto(p.coord, p.ir.success),
			])
		suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	}
}

# record a_Case( expr clauseList dflt )
procedure ir_a_Case(p, st, target, bounded, rval) {
	local t
	local L
	local i
	local e
	local v
	local x
	local oiu
	local indirects

	indirects := []

	/p.dflt := a_Fail(p.coord)

	ir_init(p)
	/bounded & (t := ir_tmploc(st))
	e := ir_tmp(st)
	v := (\target | ir_tmp(st))

	suspend ir(p.expr, st, e, "always bounded", "rval")

	every i := 1 to *p.clauseList do {
		suspend ir(p.clauseList[i].expr, st, v, nil, "rval")
		suspend ir(p.clauseList[i].body, st, target, bounded, rval)
	}
	suspend ir(p.dflt, st, target, bounded, rval)

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])

	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])

	L := p.clauseList
	if *L = 0 then {
		suspend ir_chunk(p.expr.ir.success, [ ir_Goto(p.coord, p.dflt.ir.start) ])
	} else {
		suspend ir_chunk(p.expr.ir.success, [ ir_Goto(p.coord, L[1].expr.ir.start) ])
		every i := 1 to *L do {
			suspend ir_chunk(L[i].expr.ir.success, [
				ir_opfn(p.coord, nil, nil, ir_operator("===", 2, "rval"), [ e, v ],
						L[i].expr.ir.resume),
				ir_Goto(p.coord, L[i].body.ir.start),
				])
			suspend ir_chunk(L[i].expr.ir.failure,
							 [ ir_Goto(p.coord, L[i+1].expr.ir.start) ])
			if /bounded then {
				indirects.put(L[i].body.ir.resume)
				suspend ir_chunk(L[i].body.ir.success, [
					ir_MoveLabel(p.coord, t, L[i].body.ir.resume),
					ir_Goto(p.coord, p.ir.success),
					])
			} else {
				suspend ir_chunk(L[i].body.ir.success, [
					ir_Goto(p.coord, p.ir.success),
					])
			}
			suspend ir_chunk(L[i].body.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
		}
		suspend ir_chunk(L[-1].expr.ir.failure, [ ir_Goto(p.coord, p.dflt.ir.start) ])
	}

	if /bounded then {
		indirects.put(p.dflt.ir.resume)
		suspend ir_chunk(p.dflt.ir.success, [
			ir_MoveLabel(p.coord, t, p.dflt.ir.resume),
			ir_Goto(p.coord, p.ir.success),
			])
	} else {
		suspend ir_chunk(p.dflt.ir.success, [
			ir_Goto(p.coord, p.ir.success),
			])
	}
	suspend ir_chunk(p.dflt.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, t, indirects) ])
}

# record a_Every( expr body )
procedure ir_a_Every(p, st, target, bounded, rval) {
	/p.body := a_Fail(p.coord)

	ir_init_loop(p, st, target, bounded, rval)
	st.loop_stack.put(p)
	suspend ir(p.expr, st, nil, nil, "rval")
	suspend ir(p.body, st, nil, "always bounded", "rval")
	st.loop_stack.pull()

	suspend ir_chunk(p.ir.x.nextlabel, [ ir_Goto(p.coord, p.expr.ir.resume) ])
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, p.ir.x.continueTmp, p.ir.x.indirects) ])
	suspend ir_chunk(p.expr.ir.success, [ ir_Goto(p.coord, p.body.ir.start) ])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.body.ir.success, [
		ir_Goto(p.coord, p.expr.ir.resume),
		])
	suspend ir_chunk(p.body.ir.failure, [ ir_Goto(p.coord, p.expr.ir.resume) ])
}

# record a_Sectionop( op val left right )
procedure ir_a_Sectionop(p, st, target, bounded, rval) {
	local vv
	local lv
	local rv

	ir_init(p)
	vv := ir_value(p.val, st, nil)
	lv := ir_value(p.left, st, nil)
	rv := ir_value(p.right, st, target)

	suspend ir(p.val, st, vv, nil, rval)
	suspend ir(p.left, st, lv, nil, "rval")
	suspend ir(p.right, st, rv, nil, "rval")

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.val.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.right.ir.resume) ])
	suspend ir_chunk(p.val.ir.success, [ ir_Goto(p.coord, p.left.ir.start) ])
	suspend ir_chunk(p.val.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.left.ir.success, [ ir_Goto(p.coord, p.right.ir.start) ])
	suspend ir_chunk(p.left.ir.failure, [ ir_Goto(p.coord, p.val.ir.resume) ])
	suspend ir_chunk(p.right.ir.success, [
		ir_opfn(p.coord, target, nil, ir_operator(p.op, 3, rval), [ vv, lv, rv], p.right.ir.resume),
		ir_Goto(p.coord, p.ir.success),
		])
	suspend ir_chunk(p.right.ir.failure, [ ir_Goto(p.coord, p.left.ir.resume) ])
}

procedure ir_a_Parallel(p, st, target, bounded, rval) {
	local value
	local L
	local i
	local j
	local resumeTmps
	local t
	local success

	\p.coord | throw("/p.coord", p)
	ir_init(p)

	L := p.exprList

	suspend ir(L[1 to *L-1], st, nil,  bounded, "rval")
	suspend ir(L[-1],        st, target, bounded, rval)

	if \bounded then {
		# this is so simple that it is special-cased
		suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, L[1].ir.start) ])
		every i := 1 to *L do {
			suspend ir_chunk(L[i].ir.success,[ir_Goto(p.coord,L[i+1].ir.start)])
			suspend ir_chunk(L[i].ir.failure,[ir_Goto(p.coord, p.ir.failure)])
		}
		suspend ir_chunk(L[-1].ir.success, [ ir_Goto(p.coord, p.ir.success) ])
		return fail
	}

	resumeTmps := table()
	every resumeTmps[L[2 to *L]] := ir_tmploc(st)

	success := ir_label(p, "p_op")

	t := []
	every i := (!resumeTmps).key do {
		t.put(ir_MoveLabel(p.coord, resumeTmps[i], i.ir.start))
	}
	suspend ir_chunk(p.ir.start, t ||| [
		ir_Goto(p.coord, L[1].ir.start),
		])

	suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, L[1].ir.resume) ])

	every i := 1 to *L & j := L[i+1] do {
		suspend ir_chunk(L[i].ir.success, [
			ir_IndirectGoto(p.coord, resumeTmps[j], [j.ir.start, j.ir.resume] ),
			])
	}
	suspend ir_chunk((!L).ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(L[-1].ir.success, [ ir_Goto(p.coord, success) ])

	t := []
	every i := (!resumeTmps).key do {
		t.put(ir_MoveLabel(p.coord, resumeTmps[i], i.ir.resume))
	}
	suspend ir_chunk(success, t ||| [ ir_Goto(p.coord, p.ir.success) ])
}

# record a_Call( fn args )
procedure ir_a_Call(p, st, target, bounded, rval) {
	local value
	local L
	local i
	local fn
	local args
	local clsr

	\p.coord | throw("/p.coord", p)
	every /(!p.args.exprList) := a_Nil(p.coord)
	type(p.args) === a_Arglist | throw("not type a_Arglist", p.args)

	ir_init(p)
	value := ir_tmp(st)
	clsr := ir_tmpclosure(st)
	fn := ir_tmp(st)
	args := []
	every i := !p.args.exprList do args.put(ir_value(i, st, nil))

	suspend ir(p.fn, st, fn, nil, "rval")
	every i := 1 to *p.args.exprList do {
		suspend ir(p.args.exprList[i], st, args[i], nil, "rval")
	}

	L := [p.fn] ||| p.args.exprList
	\p.coord | throw("/p.coord", p)

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.fn.ir.start) ])
	suspend ir_chunk(p.ir.resume, [
		ir_ResumeValue(p.coord, target, clsr, clsr, L[-1].ir.resume),
		ir_Goto(p.coord, p.ir.success),
		])
	every i := 1 to *L do {
		suspend ir_chunk(L[i].ir.success, [ ir_Goto(p.coord, L[i+1].ir.start) ])
		suspend ir_chunk(L[i].ir.failure, [ ir_Goto(p.coord, L[i-1].ir.resume) ])
	}
	suspend ir_chunk(L[ 1].ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(L[-1].ir.success, [
		ir_Call(p.coord, value, clsr, fn, args, p.args.nameList, L[-1].ir.resume, ir_stname(st.syms)),
		ir_Move(p.coord, target, value),
		ir_Goto(p.coord, p.ir.success),
		])
}

procedure ir_conjunction(p, st, target, bounded, rval) {
	ir_init(p)
	suspend ir(p.left, st, nil, nil, "rval")
	suspend ir(p.right, st, target, bounded, rval)
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.left.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.right.ir.resume) ])
	suspend ir_chunk(p.left.ir.success, [ ir_Goto(p.coord, p.right.ir.start) ])
	suspend ir_chunk(p.left.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.right.ir.success, [ ir_Goto(p.coord, p.ir.success) ])
	suspend ir_chunk(p.right.ir.failure, [ ir_Goto(p.coord, p.left.ir.resume) ])
}

procedure ir_augmented_assignment(p, target, bounded, rval, lv, rv, tmp) {
	local op

		op := p.op[1:-2]
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.right.ir.resume) ])
	suspend ir_chunk(p.right.ir.success, [
		ir_opfn(p.coord, tmp, nil, ir_operator(op, 2, "rval"), [ lv, rv ], p.right.ir.resume),
		ir_opfn(p.coord, target, nil, ir_operator(":=", 2, rval), [ lv, tmp ],
				p.right.ir.resume),
		ir_Goto(p.coord, p.ir.success),
		])
}

procedure ir_binary(p, target, bounded, rval, lv, rv, clsr, funcs) {
	local args

	args := [ lv, rv ]
	if funcs.member(p.op) then {
		/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.right.ir.resume) ])
		suspend ir_chunk(p.right.ir.success, [
			ir_opfn(p.coord, target, nil, ir_operator(p.op, 2, rval), args,
					p.right.ir.resume),
			ir_Goto(p.coord, p.ir.success),
			])
	} else {
		suspend ir_chunk(p.ir.resume, [
			ir_ResumeValue(p.coord, target, clsr, clsr, p.right.ir.resume),
			ir_Goto(p.coord, p.ir.success),
			])
		suspend ir_chunk(p.right.ir.success, [
			ir_opfn(p.coord, target, clsr, ir_operator(p.op, 2, rval), args,
					p.right.ir.resume),
			# ir_Move(p.coord, target, clsr),
			ir_Goto(p.coord, p.ir.success),
			])
	}
}

procedure ir_rval(op, arity, arg, parent) {
	if op == (":=:"|"<->") then {
		return nil
	} else if !!contains(op, ":=" | "<-") & arg = 1 then {
		return nil
	} else if op[1] == "@" then {
		return nil	# need lvalue for @s or s@:x
	} else if op == "[]" & arg = 1 then {
		return parent
	} else if op == "!" & arity = 1 then {
		return parent
	} else if op == "?" & arity = 1 then {
		return parent
	} else if op == "/" & arity = 1 then {
		return parent
	} else if op == "\\" & arity = 1 then {
		return parent
	} else {
		return "rval"
	}
}

# record a_Binop( op left right )
procedure ir_a_Binop(p, st, target, bounded, rval) {
	local clsr
	local tmp
	local op
	local lv
	local rv
	static funcs	# functions for which resumption fails immediately.
	/funcs := set([ ":=", ":=:", "&", ".", "[]", "+", "-", "/",
			"*", "%", "^", "**", "++", "--", "<", "<=", "=", "~=",
			">=", ">", "<<", "<<=", "==", "~==", ">>=", ">=", ">>",
			"===", "~===", "|||", "||" ])

	/p.right := a_Nil(p.coord)

	if p.op == "&" then {
		suspend ir_conjunction(p, st, target, bounded, rval)
		return fail
	}

	ir_init(p)
	if not funcs.member(p.op) &
			not (funcs.member(p.op[1:-2]) & p.op[-2:0] == ":=") then {
		clsr := ir_tmpclosure(st)
	}
	lv := ir_value(p.left, st, nil)
	rv := ir_value(p.right, st, target)
	tmp := (\target | ir_tmp(st))

	suspend ir(p.left,  st, lv, nil, ir_rval(p.op, 2, 1, rval))
	suspend ir(p.right, st, rv, nil, ir_rval(p.op, 2, 2, rval))

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.left.ir.start) ])
	suspend ir_chunk(p.left.ir.success, [ ir_Goto(p.coord, p.right.ir.start) ])
	suspend ir_chunk(p.left.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.right.ir.failure, [ ir_Goto(p.coord, p.left.ir.resume) ])

	if *p.op > 2 & p.op[-2:0] == ":=" then {
		suspend ir_augmented_assignment(p, target, bounded, rval, lv, rv, tmp)
	} else {
		suspend ir_binary(p, target, bounded, rval, lv, rv, clsr, funcs)
	}
}

procedure ir_unary(p, target, bounded, rval, v, clsr, funcs) {
	if funcs.member(p.op) then {
		/bounded & suspend ir_chunk(p.ir.resume, [ir_Goto(p.coord, p.operand.ir.resume)])
		suspend ir_chunk(p.operand.ir.success, [
			ir_opfn(p.coord, target, nil, ir_operator(p.op, 1, rval), [ v ], p.operand.ir.resume),
			ir_Goto(p.coord, p.ir.success),
			])
	} else {
		suspend ir_chunk(p.ir.resume, [
			ir_ResumeValue(p.coord, target, clsr, clsr, p.operand.ir.resume),
			ir_Goto(p.coord, p.ir.success),
			])
		suspend ir_chunk(p.operand.ir.success, [
			ir_opfn(p.coord, target, clsr, ir_operator(p.op, 1, rval), [ v ], p.operand.ir.resume),
			# ir_Move(p.coord, target, closure),
			ir_Goto(p.coord, p.ir.success),
			])
	}
}

# record a_Unop( op operand )
procedure ir_a_Unop(p, st, target, bounded, rval) {
	local closure
	local v
	local t
	static funcs	# functions for which resumption fails immediately.
	/funcs := set([ "/", "\\", "*", "?", "+", "-", "~", "^", "@" ])

	ir_init(p)
	if not funcs.member(p.op) then {
		closure := ir_tmpclosure(st)
	}
	v := ir_value(p.operand, st, target)

	suspend ir(p.operand, st, v, nil, ir_rval(p.op, 1, 1, rval))

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.operand.ir.start) ])
	suspend ir_unary(p, target, bounded, rval, v, closure, funcs)
	suspend ir_chunk(p.operand.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Package( name )
procedure ir_a_Package(p, st, target, bounded, rval) {
	ir_namespace := p.name.id
	return fail						# generate no code
}

# record a_If( expr thenexpr elseexpr )
procedure ir_a_If(p, st, target, bounded, rval) {
	local t

	/p.elseexpr := a_Fail(p.coord)

	ir_init(p)
	/bounded & (t := ir_tmploc(st))

	suspend ir(p.expr, st, nil, "always bounded", "rval")
	suspend ir(p.thenexpr, st, target, bounded, rval)
	suspend ir(p.elseexpr, st, target, bounded, rval)

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, t,
		[p.thenexpr.ir.resume, p.elseexpr.ir.resume]) ])
	if /bounded then {
		suspend ir_chunk(p.expr.ir.success, [
			ir_MoveLabel(p.coord, t, p.thenexpr.ir.resume),
			ir_Goto(p.coord, p.thenexpr.ir.start),
			])
		suspend ir_chunk(p.expr.ir.failure, [
			ir_MoveLabel(p.coord, t, p.elseexpr.ir.resume),
			ir_Goto(p.coord, p.elseexpr.ir.start),
			])
	} else {
		suspend ir_chunk(p.expr.ir.success, [
			ir_Goto(p.coord, p.thenexpr.ir.start),
			])
		suspend ir_chunk(p.expr.ir.failure, [
			ir_Goto(p.coord, p.elseexpr.ir.start),
			])
	}
	suspend ir_chunk(p.thenexpr.ir.success, [ ir_Goto(p.coord, p.ir.success) ])
	suspend ir_chunk(p.thenexpr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.elseexpr.ir.success, [ ir_Goto(p.coord, p.ir.success) ])
	suspend ir_chunk(p.elseexpr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Global( id)
procedure ir_a_Global(p, st, target, bounded, rval) {
	local fn
	local i
	local susp
	static counter
	/counter := 0

	/st | throw("\st", st)
	if \p.expr then {
		st := ir_stacks(0)
		st.localSet := set([])
		st.staticSet := set([])
		st.globalSet := set([])
		st.loop_stack := []

		st.current_proc := "$global$" || counter
		counter +:= 1

		fn := ir_a_ProcDecl1(st, p.expr, [], nil, p.coord)
		suspend fn
	}
	return ir_Global(p.coord, p.id, \(\st).current_proc | nil, ir_namespace)
}

# record a_Initial( expr )
procedure ir_a_Initial(p, st, target, bounded, rval) {
	local fn
	local i
	static counter
	/counter := 0

	/st | throw("\st", st)
	st := ir_stacks(0)
	st.localSet := set([])
	st.staticSet := set([])
	st.globalSet := set([])

	st.current_proc := "$initial$" || counter
	counter +:= 1

	fn := ir_a_ProcDecl1(st, p.expr, [], nil, p.coord)
	suspend fn
	i := ir_Initial(p.coord, st.current_proc, ir_namespace)
	return i
}

procedure ir_value(p, st, target) {
	return ( \target | ir_tmp(st))
}

# record a_Intlit( int )
procedure ir_a_Intlit(p, st, target, bounded, rval) {

	if type(target) === ir_IntLit then target := nil

	ir_init(p)

	if not (p.int := integer(p.int)) then {
		semantic_error(p.int || ": Illegal integer literal", p.coord)
	}
	suspend ir_chunk(p.ir.start, [
		ir_IntLit(p.coord, target, p.int),
		ir_Goto(p.coord, p.ir.success),
		])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Reallit( real )
procedure ir_a_Reallit(p, st, target, bounded, rval) {

	if type(target) === ir_RealLit then target := nil

	ir_init(p)

	if not (p.real := number(p.real)) then {
		semantic_error(p.real || ": Illegal real literal", p.coord)
	}
	suspend ir_chunk(p.ir.start, [
		ir_RealLit(p.coord, target, p.real),
		ir_Goto(p.coord, p.ir.success),
		])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Stringlit( str )
procedure ir_a_Stringlit(p, st, target, bounded, rval) {

	if type(target) === ir_StrLit then target := nil

	ir_init(p)

	suspend ir_chunk(p.ir.start, [
		ir_StrLit(p.coord, target, *p.str, p.str),
		ir_Goto(p.coord, p.ir.success),
		])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_ProcDecl( ident paramList accumulate localsList code )
procedure ir_a_ProcDeclNested(p, st, target, bounded, rval) {
	local f
	local scope
	local st1
	static counter
	/counter := 0

	counter +:= 1

	ir_init(p)

	p.ident.id := "$" || st.current_proc || "$nested$" || counter

	suspend ir_chunk(p.ir.start, [
		ir_MakeClosure(p.coord, target, p.ident.id),
		ir_Goto(p.coord, p.ir.success),
		])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])

	scope := st.syms
	st1 := ir_stacks(0)
	st1.parent := st.current_proc
	st1.syms := scope
	f := ir_a_ProcDecl0(p, st1)
	ir_deferred.push(f)
}

# record a_ProcDecl( ident paramList accumulate localsList code )
procedure ir_a_ProcDecl(p, st, target, bounded, rval) {
	st := ir_stacks(0)
	suspend ir_a_ProcDecl0(p, st)
}

procedure ir_a_ProcDecl0(p, st) {
	local locals
	local statics
	local globals
	local params
	local i
	local v
	local code
	local s
	local f

	st.localSet := set([])
	st.staticSet := set([])
	st.globalSet := set([])
	st.current_proc := p.ident.id

	v := set([])
	st.syms := ir_scope(st.syms)
	st.syms.Static := table()
	st.syms.Dynamic := table()
	params := []
	every i := !p.paramList do {
		if v.member(i.id) then {
			semantic_error(image(i.id) || ": Redeclared identifier", i.coord)
		}
		v.put(i.id)
		s := i.id || mkSuffix(st.syms)
		st.syms.Static[i.id] := s
		params.put(s)
	}
	if ir_declare_set.member(p.ident.id) then {
		semantic_error(image(p.ident.id) || ": Inconsistent redeclaration",
					   p.ident.coord)
	}
	ir_declare_set.put(p.ident.id)

	f := ir_a_ProcDecl1(st, p.code, params, p.accumulate, p.ident.coord)
	st.syms := st.syms.parent	# this is very odd, why is it here?
								# Should it be in the caller?
	return f
}

procedure ir_a_ProcDecl1(st, body, params, accumulate, coord) {
	local code
	local locals
	local statics
	local globals

	code := []
	every code.put(ir(body, st))

	locals := []
	every locals.put(!st.localSet)
	statics := []
	every statics.put(!st.staticSet)
	globals := []
	every globals.put(!st.globalSet)

	return ir_Function(coord, st.current_proc, params,
					  accumulate, locals, statics, globals, code,
					  body.ir.start, st.parent, ir_namespace)
}

# record a_ProcCode( body )
procedure ir_a_ProcCode(p, st, target, bounded, rval) {
	ir_init(p)

	st.loop_stack := []

	suspend ir(p.body, st, nil, "always bounded", "rval")

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.body.ir.start) ])
	suspend ir_chunk(p.ir.resume, [ ir_Unreachable(p.coord) ])
	suspend ir_chunk(p.body.ir.success, [ ir_Fail(p.coord) ])
	suspend ir_chunk(p.body.ir.failure, [ ir_Fail(p.coord) ])
}

# record a_Record( ident idlist )
procedure ir_a_Record(p, st, target, bounded, rval) {
	local fields
	local v
	local i

	if ir_declare_set.member(p.ident.id) then {
		semantic_error(image(p.ident.id) || ": Inconsistent redeclaration",
					   p.ident.coord)
	}
	ir_declare_set.put(p.ident.id)
	v := set([])
	every i := !p.idlist do {
		if v.member(i.id) then {
			semantic_error(image(i.id) || ": Redeclared identifier", i.coord)
		}
		v.put(i.id)
	}
	fields := []
	every fields.put((!p.idlist).id)
	return ir_Record(p.ident.coord, p.ident.id, (\p.extendsRec).id |nil, (\p.extendsPkg).id | nil, fields, ir_namespace)
}

# record a_Repeat( expr )
# old
procedure ir_a_RepeatX(p, st, target, bounded, rval) {
	ir_init_loop(p, st, target, bounded, rval)
	st.loop_stack.put(p)
	suspend ir(p.body, st, nil, "always bounded", "rval")
	st.loop_stack.pull()

	suspend ir_chunk(p.ir.x.nextlabel, [ ir_Goto(p.coord, p.body.ir.start) ])
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.body.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, p.ir.x.continueTmp, p.ir.x.indirects) ])
	suspend ir_chunk(p.body.ir.success, [
		ir_Goto(p.coord, p.ir.start),
		])
	suspend ir_chunk(p.body.ir.failure, [ ir_Goto(p.coord, p.ir.start) ])
}

# record a_Catch( expr )
procedure ir_a_Catch(p, st, target, bounded, rval) {
	local t
	local mk

	ir_init(p)
	t := ir_tmp(st)

	suspend ir(p.expr, st, t, "always bounded", "rval")

	suspend ir_chunk(p.ir.start,        [ ir_Goto(p.coord, p.expr.ir.start) ])
	# should this really fail on resumption?
	suspend ir_chunk(p.ir.resume,       [ ir_Goto(p.coord, p.ir.failure) ])

	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.expr.ir.success, [
		ir_Catch(p.coord, target, t),
		ir_Goto(p.coord, p.ir.success),
		])
}

# record a_Return( expr )
procedure ir_a_Return(p, st, target, bounded, rval) {
	local t
	local mk

	/st.createflag | semantic_error("Invalid context for return", p.coord)

	/p.expr := a_Nil(p.coord)

	ir_init(p)
	t := ir_tmp(st)
	mk := ir_tmp(st)

	suspend ir(p.expr, st, t, "always bounded", nil)

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Fail(p.coord) ])
		suspend ir_chunk(p.expr.ir.success, [
			ir_Succeed(p.coord, t, nil),
			])
		suspend ir_chunk(p.expr.ir.failure, [ ir_Fail(p.coord) ])
}

# record a_Fail( )
procedure ir_a_Fail(p, st, target, bounded, rval) {
	ir_init(p)
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.ir.failure) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Unreachable(p.coord) ])
}

# record a_Suspend( expr body )
procedure ir_a_Suspend(p, st, target, bounded, rval) {
	local t
	local susp

	/st.createflag | semantic_error("Invalid context for suspend", p.coord)

	/p.body := a_Fail(p.coord) & /p.expr := a_Nil(p.coord)

	ir_init_loop(p, st, target, bounded, rval)
	t := ir_label(p, "suspend")
	susp := ir_tmp(st)

	st.loop_stack.put(p)
	suspend ir(p.expr, st, susp, nil, "rval")
	suspend ir(p.body, st, nil, "always bounded", nil)
	st.loop_stack.pull()

	suspend ir_chunk(p.ir.x.nextlabel, [ ir_Goto(p.coord, p.expr.ir.resume) ])
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, p.ir.x.continueTmp, p.ir.x.indirects) ])
		suspend ir_chunk(p.expr.ir.success, [ ir_Succeed(p.coord, susp, t) ])
		suspend ir_chunk(t, [ ir_Goto(p.coord, p.body.ir.start) ])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.body.ir.success, [
		ir_Goto(p.coord, p.expr.ir.resume),
		])
	suspend ir_chunk(p.body.ir.failure, [ ir_Goto(p.coord, p.expr.ir.resume) ])
}

# record a_Until( expr body )
procedure ir_a_Until(p, st, target, bounded, rval) {
	/p.body := a_Fail(p.coord)

	ir_init_loop(p, st, target, bounded, rval)
	st.loop_stack.put(p)
	suspend ir(p.expr, st, nil, "always bounded", "rval")
	suspend ir(p.body, st, nil, "always bounded", "rval")
	st.loop_stack.pull()

	suspend ir_chunk(p.ir.x.nextlabel, [ ir_Goto(p.coord, p.expr.ir.start) ])
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, p.ir.x.continueTmp, p.ir.x.indirects) ])
	suspend ir_chunk(p.expr.ir.success, [
		ir_Goto(p.coord, p.ir.failure),
		])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.body.ir.start) ])
	suspend ir_chunk(p.body.ir.success, [
		ir_Goto(p.coord, p.expr.ir.start),
		])
	suspend ir_chunk(p.body.ir.failure, [ ir_Goto(p.coord, p.expr.ir.start) ])
}

# record a_While( expr body )
procedure ir_a_While(p, st, target, bounded, rval) {
	/p.body := a_Fail(p.coord)

	ir_init_loop(p, st, target, bounded, rval)
	st.loop_stack.put(p)
	suspend ir(p.expr, st, nil, "always bounded", "rval")
	suspend ir(p.body, st, nil, "always bounded", "rval")
	st.loop_stack.pull()

	suspend ir_chunk(p.ir.x.nextlabel, [ ir_Goto(p.coord, p.expr.ir.start) ])
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.expr.ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, p.ir.x.continueTmp, p.ir.x.indirects) ])
	suspend ir_chunk(p.expr.ir.success, [
		ir_Goto(p.coord, p.body.ir.start),
		])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.body.ir.success, [
		ir_Goto(p.coord, p.expr.ir.start),
		])
	suspend ir_chunk(p.body.ir.failure, [ ir_Goto(p.coord, p.expr.ir.start) ])
}

# record a_Repeat( body expr )
# correct
procedure ir_a_Repeat(p, st, target, bounded, rval) {
	/p.expr := a_Fail(p.coord)

	ir_init_loop(p, st, target, bounded, rval)
	st.loop_stack.put(p)
	suspend ir(p.expr, st, nil, "always bounded", "rval")
	suspend ir(p.body, st, nil, "always bounded", "rval")
	st.loop_stack.pull()

	suspend ir_chunk(p.ir.x.nextlabel, [ ir_Goto(p.coord, p.expr.ir.start) ])

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.body.ir.start) ])

	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, p.ir.x.continueTmp, p.ir.x.indirects) ])

	suspend ir_chunk(p.body.ir.success, [ ir_Goto(p.coord, p.expr.ir.start) ])
	suspend ir_chunk(p.body.ir.failure, [ ir_Goto(p.coord, p.expr.ir.start) ])

	suspend ir_chunk(p.expr.ir.success, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.body.ir.start) ])
}

# record a_Create( expr )
procedure ir_a_Create(p, st, target, bounded, rval) {
	local t

	\p.coord | throw("/p.coord", p)
	ir_init(p)
	t := (\target | ir_tmp(st))

	st.createflag := 1
	suspend ir(p.expr, st, t, nil, nil)
	st.createflag := nil

	suspend ir_chunk(p.ir.start, [
		ir_Create(p.coord, target, p.expr.ir.start),
		ir_Goto(p.coord, p.ir.success),
		])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.expr.ir.success, [ ir_CoRet(p.coord, t, p.expr.ir.resume) ])
	suspend ir_chunk(p.expr.ir.failure, [ ir_CoFail(p.coord) ])
}

procedure mkSuffix(t) {
	/static T := table()
	/T[t] := ":" || (*T + 1)
	return T[t]
}

procedure ir_stname(s) {
	return (mkSuffix(\s) | nil)
}

# record a_Local( id )
procedure ir_a_Local(p, st, target, bounded, rval) {
	local s

	/st.syms.Static := table()

	/st.syms.Static[p.id] | semantic_error(image(p.id) || ": Redeclared identifier", p.coord)

	s := p.id || ":" || image(st.syms)
	s := p.id || mkSuffix(st.syms)

	st.syms.Static[p.id] := s
	st.localSet.put(s)
	suspend ir_a_Ident(p, st, target, bounded, rval)
}

# record a_Static( id )
procedure ir_a_Static(p, st, target, bounded, rval) {
	local s

	/st.syms.Static := table()

	/st.syms.Static[p.id] | semantic_error(image(p.id) || ": Redeclared identifier", p.coord)

	s := p.id || mkSuffix(st.syms)
	st.syms.Static[p.id] := s

	st.staticSet.put(s)
	suspend ir_a_Ident(p, st, target, bounded, rval)
}

procedure ir_a_With(p, st, target, bounded, rval) {
	local L
	local i
	local dynamics
	local parent
	local tmp
	local tmpkey
	local newscope

	parent := ir_stname(st.syms)
	newscope := ir_scope(st.syms)

	ir_init(p)

	st.syms := newscope
	st.syms.Static := table()
	st.syms.Dynamic := table()
	st.syms.Dynamic[p.id] := p.id
	suspend ir(p.expr, st, target, bounded, rval)
	st.syms := st.syms.parent

	if \p.init then {
		tmp := ir_tmp(st)
		tmpkey := ir_tmp(st)
		suspend ir(p.init, st, tmp, "bounded", "rval")
		suspend ir_chunk(p.init.ir.success, [
			ir_EnterScope(p.coord, [], [p.id], ir_stname(newscope), parent),
			ir_Key(p.coord, tmpkey, p.id, ir_stname(newscope), nil),
			ir_opfn(p.coord, nil, nil, ir_operator(":=", 2, "rval"), [ tmpkey, tmp ], nil),
			ir_Goto(p.coord, p.expr.ir.start),
			])
		suspend ir_chunk(p.init.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
		suspend ir_chunk(p.ir.start, [
			ir_Goto(p.coord, p.init.ir.start),
		])
	} else {
		suspend ir_chunk(p.ir.start, [
			ir_EnterScope(p.coord, [], [p.id], ir_stname(newscope), parent),
			ir_Goto(p.coord, p.expr.ir.start),
		])
	}

	if /bounded then {
		suspend ir_chunk(p.expr.ir.success, [
			ir_ExitScope(p.coord, [], [p.id], ir_stname(newscope)),
			ir_Goto(p.coord, p.ir.success),
			])
		suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.expr.ir.resume) ])
	} else {
		suspend ir_chunk(p.expr.ir.success, [ ir_Goto(p.coord, p.ir.success) ])
	}

	suspend ir_chunk(p.expr.ir.failure, [
		ir_ExitScope(p.coord, [], [p.id], ir_stname(newscope)),
		ir_Goto(p.coord, p.ir.failure),
		])


}

procedure ir_lookupD(st, name) {
	local s
	local syms
	syms := st.syms
	while \syms do {
		if s := \(\syms.Dynamic)[name] then return syms
		syms := syms.parent
	}
	return nil
}

procedure ir_lookupS(st, name) {
	local s
	local syms
	syms := st.syms
	while \syms do {
		if s := \(\syms.Static)[name] then return syms
		syms := syms.parent
	}
	st.globalSet.put(name)
	return nil
}

# record a_Ident( id )
procedure ir_a_Ident(p, st, target, bounded, rval) {
	local s
	local T

	if type(\target) ~=== ir_Tmp then {
		# #%#%# prevents nasty interaction with targeting.
		# %#%#% probably a symptom of bad design....
		target := nil
	}

	ir_init(p)

	if /p.namespace then {
		T := ir_lookupS(st, p.id)
		s := ((\T).Static[p.id] | p.id)
	} else {
		type(p.namespace) === string | throw("namespace not string", p)
		s := p.id
		st.globalSet.put(p.namespace || "::" || p.id)
	}

	suspend ir_chunk(p.ir.start, [
		ir_Var(p.coord, target, s, p.namespace, ir_stname(T), rval),
		ir_Goto(p.coord, p.ir.success),
		])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Continue( )
procedure ir_a_Continue(p, st, target, bounded, rval) {
	local curloop
	local sl
	local mk

	ir_init(p)
	if /p.name then {
			st.loop_stack[1] | semantic_error("Invalid context for continue", p.coord)

			curloop := st.loop_stack[-1]
	} else {
		if \(curloop <- st.loop_stack[*st.loop_stack to 1 by -1]).name == p.name then {
			# set curloop above
		} else {
			semantic_error("Undeclared continue identifier", p.coord)
		}
	}
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, curloop.ir.x.nextlabel) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Unreachable(p.coord) ])
}

# record a_Break( expr )
procedure ir_a_Break(p, st, target, bounded, rval) {
	local curloop
	local sl
	local mk

	ir_init(p)
	if /p.name then {
			st.loop_stack[1] | semantic_error("Invalid context for break", p.coord)

			curloop := st.loop_stack[-1]
	} else {
		if \(curloop <- st.loop_stack[*st.loop_stack to 1 by -1]).name == p.name then {
			# set curloop above
		} else {
			semantic_error("Undeclared break identifier", p.coord)
		}
	}
	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, curloop.ir.failure) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Unreachable(p.coord) ])
}

# record a_Yield( expr name )
procedure ir_a_Yield(p, st, target, bounded, rval) {
	local curloop
	local sl
	local mk
	local clx

	/p.expr := a_Nil(p.coord)
	ir_init(p)

	if /p.name then {
			st.loop_stack[1] | semantic_error("Invalid context for yield", p.coord)

			curloop := st.loop_stack[-1]
	} else {
		if \(curloop <- st.loop_stack[*st.loop_stack to 1 by -1]).name == p.name then {
			# set curloop above
		} else {
			semantic_error("Undeclared yield identifier", p.coord)
		}
	}
	clx := curloop.ir.x
	suspend ir(p.expr, st, clx.value, clx.bounded, clx.rval)

	if /clx.bounded then {
		clx.indirects.put(p.ir.resume)
		suspend ir_chunk(p.ir.start, [
		   ir_MoveLabel(p.coord, clx.continueTmp, p.ir.resume),
		   ir_Goto(p.coord, p.expr.ir.start),
		])
	} else {
		suspend ir_chunk(p.ir.start, [
			ir_Goto(p.coord, p.expr.ir.start),
		])
	}
	/clx.bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.expr.ir.resume) ])
	suspend ir_chunk(p.expr.ir.success, [ ir_Goto(p.coord, curloop.ir.success) ])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Break( expr )
procedure ir_a_BreakX(p, st, target, bounded, rval) {
	local curloop
	local oldloops
	local clx
	local mk

	st.loop_stack[1] | semantic_error("Invalid context for break", p.coord)

	/p.expr := a_Nil(p.coord)

	ir_init(p)

	curloop := st.loop_stack[-1]
	oldloops := st.loop_stack
	st.loop_stack := st.loop_stack[1:-1]
	clx := curloop.ir.x
	suspend ir(p.expr, st, clx.value, clx.bounded, clx.rval)
	st.loop_stack := oldloops

		if /clx.bounded then {
			curloop.indirects.put(p.ir.resume)
			suspend ir_chunk(p.ir.start, [
				ir_MoveLabel(p.coord, clx.continueTmp, p.ir.resume),
				ir_Goto(p.coord, p.expr.ir.start),
				])
		} else {
			suspend ir_chunk(p.ir.start, [
				ir_Goto(p.coord, p.expr.ir.start),
				])
		}
	/clx.bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.expr.ir.resume) ])
	suspend ir_chunk(p.expr.ir.success, [ ir_Goto(p.coord, curloop.ir.success) ])
	suspend ir_chunk(p.expr.ir.failure, [ ir_Goto(p.coord, curloop.ir.failure) ])
}

# record a_ToBy( fromexpr toexpr byexpr )
procedure ir_a_ToBy(p, st, target, bounded, rval) {
	local clsr
	local fv
	local tv
	local bv

	/p.byexpr := a_Intlit(1, p.coord)

	ir_init(p)
	clsr := ir_tmpclosure(st)
	fv := ir_value(p.fromexpr, st, nil)
	tv := ir_value(p.toexpr, st, nil)
	bv := ir_value(p.byexpr, st, target)

	suspend ir(p.fromexpr, st, fv, nil, "rval")
	suspend ir(p.toexpr, st, tv, nil, "rval")
	suspend ir(p.byexpr, st, bv, nil, "rval")

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.fromexpr.ir.start) ])
	suspend ir_chunk(p.ir.resume, [
		ir_ResumeValue(p.coord, target, clsr, clsr, p.byexpr.ir.resume),
		ir_Goto(p.coord, p.ir.success),
		])
	suspend ir_chunk(p.fromexpr.ir.success, [ ir_Goto(p.coord, p.toexpr.ir.start) ])
	suspend ir_chunk(p.fromexpr.ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.toexpr.ir.success, [ ir_Goto(p.coord, p.byexpr.ir.start) ])
	suspend ir_chunk(p.toexpr.ir.failure, [ ir_Goto(p.coord, p.fromexpr.ir.resume) ])
	suspend ir_chunk(p.byexpr.ir.success, [
		ir_opfn(p.coord, target, clsr, ir_operator("...", 3, "rval"), [ fv, tv, bv ], p.byexpr.ir.resume),
		# ir_Move(p.coord, target, closure),
		ir_Goto(p.coord, p.ir.success),
		])
	suspend ir_chunk(p.byexpr.ir.failure, [ ir_Goto(p.coord, p.toexpr.ir.resume) ])

}

# record a_Select( caseList, dflt )
procedure ir_a_Select(p, st, target, bounded, rval) {
	local i
	local left
	local right
	local caseList
	local sc
	local c
	local t
	local selectlab
	local lab
	local dest
	local indirects

	indirects := []

	ir_init(p)

	/bounded & (t := ir_tmploc(st))

	selectlab := ir_label(p, "select")

	caseList := []
	every i := *p.caseList to 1 by -1 do {
		c := p.caseList[i]
		suspend ir(c.body, st, target, bounded, rval)
		suspend ir_chunk(c.body.ir.failure, [ ir_Goto(c.coord, p.ir.failure) ])
		suspend ir_chunk(c.body.ir.success, [ ir_Goto(c.coord, p.ir.success) ])
		if /bounded then {
			lab := ir_label(c, "setup")
			indirects.put(c.body.ir.resume)
			suspend ir_chunk(lab, [
				ir_MoveLabel(c.coord, t, c.body.ir.resume),
				ir_Goto(c.coord, c.body.ir.start),
				])
		} else {
			lab := c.body.ir.start
		}

		left := ir_tmp(st)
		right := ir_tmp(st)
		sc := ir_SelectCase(c.coord, c.kind, left, right, lab)
		caseList.put(sc)
		suspend ir(c.left, st, left, nil, nil)
		suspend ir(c.right, st, right, "bounded", "rval")
		suspend ir_chunk(c.left.ir.success, [ ir_Goto(c.coord, c.right.ir.start) ])
		dest := p.caseList[i+1].left.ir.start | selectlab
		suspend ir_chunk(c.left.ir.failure, [
			ir_NoValue(c.coord, left),
			ir_Goto(c.coord, dest),
			])
		dest := p.caseList[i+1].left.ir.start | selectlab
		suspend ir_chunk(c.right.ir.success, [ ir_Goto(c.coord, dest) ])
		suspend ir_chunk(c.right.ir.failure, [ ir_Goto(c.coord, c.left.ir.resume) ])
	}
	if c := \p.dflt then {
		suspend ir(c.body, st, target, bounded, rval)
		suspend ir_chunk(c.body.ir.failure, [ ir_Goto(c.coord, p.ir.failure) ])
		suspend ir_chunk(c.body.ir.success, [ ir_Goto(c.coord, p.ir.success) ])
		if /bounded then {
			lab := ir_label(c, "setup")
			indirects.put(c.body.ir.resume)
			suspend ir_chunk(lab, [
				ir_MoveLabel(c.coord, t, c.body.ir.resume),
				ir_Goto(c.coord, c.body.ir.start),
				])
		} else {
			lab := c.body.ir.start
		}
		sc := ir_SelectCase(c.coord, "default", nil, nil, c.body.ir.start)
		caseList.put(sc)
	}
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_IndirectGoto(p.coord, t, indirects) ])

	suspend ir_chunk(selectlab, [ ir_Select(p.coord, caseList, p.ir.failure) ])

	if p.caseList[1] then {
		suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, p.caseList[1].left.ir.start) ])
	} else {
		suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, selectlab) ])
	}
}

# record a_Mutual( exprList )
procedure ir_a_Mutual(p, st, target, bounded, rval) {
	local L
	local i

	if *p.exprList = 0 then {
		p.exprList := [ a_Nil(p.coord) ]
	} else {
		every /(!p.exprList) := a_Nil(p.coord)
	}

	ir_init(p)

	every i := 1 to *p.exprList-1 do {
		suspend ir(p.exprList[i], st, nil, nil, "rval")
	}
	suspend ir(p.exprList[-1], st, target, bounded, rval)
	L := p.exprList

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, L[1].ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, L[-1].ir.resume) ])

	every i := 1 to *L do {
		suspend ir_chunk(L[i].ir.success, [ ir_Goto(p.coord, L[i+1].ir.start) ])
		suspend ir_chunk(L[i].ir.failure, [ ir_Goto(p.coord, L[i-1].ir.resume) ])
	}
	suspend ir_chunk(L[-1].ir.success, [ ir_Goto(p.coord, p.ir.success) ])
	suspend ir_chunk(L[ 1].ir.failure, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Compound( exprList )
procedure ir_a_Compound(p, st, target, bounded, rval) {
	local L
	local i
	local locals
	local failure
	local dynamics
	local parent

	parent := ir_stname(st.syms)

	st.syms := ir_scope(st.syms)

	failure := ir_label(p, "exit")

	every /(!p.exprList) := a_Nil(p.coord)

	ir_init(p)

	every i := 1 to *p.exprList-1 do {
		suspend ir(p.exprList[i], st, nil, "always bounded", "rval")
	}
	suspend ir(p.exprList[-1], st, target, bounded, rval)

	L := p.exprList
	*L > 0 | throw("*L=0", L)

	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, L[-1].ir.resume) ])

	every i := 1 to *p.exprList-1 do {
		suspend ir_chunk(L[i].ir.success, [
			ir_Goto(p.coord, L[i+1].ir.start),
			])
		suspend ir_chunk(L[i].ir.failure, [ ir_Goto(p.coord, L[i+1].ir.start) ])
	}

	locals := []
	every i := !\st.syms.Static do {
		if st.localSet.member(i.value) then {
			locals.put(i.value)
		}
	}
	dynamics := []
	every i := !\st.syms.Dynamic do {
		dynamics.put(i)
	}

	# should the following have ExitScope if this is bounded?
	if /bounded then {
		suspend ir_chunk(L[-1].ir.success, [
			ir_ExitScope(p.coord, locals, dynamics, ir_stname(st.syms)),
			ir_Goto(p.coord, p.ir.success),
			])
	} else {
		suspend ir_chunk(L[-1].ir.success, [ ir_Goto(p.coord, p.ir.success) ])
	}

	suspend ir_chunk(L[-1].ir.failure, [ ir_Goto(p.coord, failure) ])

	suspend ir_chunk(p.ir.start, [
		ir_EnterScope(p.coord, locals, dynamics, ir_stname(st.syms), parent),
		ir_Goto(p.coord, L[1].ir.start),
	])
	suspend ir_chunk(failure, [
		ir_ExitScope(p.coord, locals, dynamics, ir_stname(st.syms)),
		ir_Goto(p.coord, p.ir.failure),
	])

	st.syms := st.syms.parent
}

procedure ir_a_Nil(p, st, target, bounded, rval) {
	ir_init(p)
	suspend ir_chunk(p.ir.start, [
		ir_NilLit(p.coord, target),
		ir_Goto(p.coord, p.ir.success),
		])
	suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_Key( id )
procedure ir_a_Key(p, st, target, bounded, rval) {
	local s

	ir_init(p)
	\p.coord | throw("/p.coord", p)

	suspend ir_chunk(p.ir.start, [
		ir_Key(p.coord, target, p.id, ir_stname(ir_lookupD(st,p.id)), rval),
		ir_Goto(p.coord, p.ir.success),
		])
	suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
}

# record a_ListConstructor( exprList )
procedure ir_a_ListConstructor(p, st, target, bounded, rval) {
	local L
	local i
	local args

	\p.coord | throw("/p.coord", p)
	every /(!p.exprList) := a_Nil(p.coord)

	ir_init(p)

	args := []
	if \target then {
		every i := !p.exprList do args.put(ir_value(i, st, nil))
	} else {
		every !p.exprList do args.put(nil)
	}

	every i := 1 to *p.exprList do {
		suspend ir(p.exprList[i], st, args[i], nil, "rval")
	}

	L := ir_make_sentinel(p.exprList)

	suspend ir_chunk(p.ir.start, [ ir_Goto(p.coord, L[1].ir.start) ])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, L[-1].ir.resume) ])

	suspend ir_chunk(L[1].ir.start, [ ir_Goto(p.coord, L[2].ir.start) ])
	suspend ir_chunk(L[1].ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
	every i := 2 to *p.exprList-1 do {
		suspend ir_chunk(L[i].ir.success, [ ir_Goto(p.coord, L[i+1].ir.start) ])
		suspend ir_chunk(L[i].ir.failure, [ ir_Goto(p.coord, L[i-1].ir.resume) ])
	}
	\p.coord | throw("/p.coord", p)
	if \target then {
		suspend ir_chunk(L[-1].ir.start, [
			ir_MakeList(p.coord, target, args),
			ir_Goto(p.coord, p.ir.success),
			])
	} else {
		suspend ir_chunk(L[-1].ir.start, [ ir_Goto(p.coord, p.ir.success) ])
	}
	suspend ir_chunk(L[-1].ir.resume, [ ir_Goto(p.coord, L[-2].ir.resume) ])

}

# record a_ListComprehension( expr )
procedure ir_a_ListComprehension(p, st, target, bounded, rval) {
	local tmp
	local v

	\p.coord | throw("/p.coord", p)
	ir_init(p)

	tmp := ir_tmp(st)
	v := ir_tmp(st)

	suspend ir(p.expr, st, tmp, nil, "rval")

	suspend ir_chunk(p.ir.start, [
		ir_MakeList(p.coord, v, []),
		ir_Goto(p.coord, p.expr.ir.start),
		])
	/bounded & suspend ir_chunk(p.ir.resume, [ ir_Goto(p.coord, p.ir.failure) ])
	suspend ir_chunk(p.expr.ir.success, [
		ir_opfn(p.coord, nil, nil, ir_operator("put", 2, rval), [ v, tmp ], p.ir.failure),
		ir_Goto(p.coord, p.expr.ir.resume),
		])
	suspend ir_chunk(p.expr.ir.failure, [
		ir_Move(p.coord, target, v),
		ir_Goto(p.coord, p.ir.success),
		])

}

procedure ir_outer(p) {
	case type(p) of {
		a_ProcDecl : suspend ir_a_ProcDecl(p)
		a_Global : suspend ir_a_Global(p)
		a_Record : suspend ir_a_Record(p)
		a_Package : suspend ir_a_Package(p)
		a_Initial : suspend ir_a_Initial(p)
		default : throw("unrecognized type", p)
	}
}

procedure ir(p, st, target, bounded, rval) {
	case type(p) of {
		a_NoOp : suspend ir_a_NoOp(p, st, target, bounded, rval)
		a_Field : suspend ir_a_Field(p, st, target, bounded, rval)
		a_Call : suspend ir_a_Call(p, st, target, bounded, rval)
		a_Paired : suspend ir_a_Paired(p, st, target, bounded, rval)
		a_Limitation : suspend ir_a_Limitation(p, st, target, bounded, rval)
		a_Not : suspend ir_a_Not(p, st, target, bounded, rval)
		a_Alt : suspend ir_a_Alt(p, st, target, bounded, rval)
		a_ExcAlt : suspend ir_a_ExcAlt(p, st, target, bounded, rval)
		a_RepAlt : suspend ir_a_RepAlt(p, st, target, bounded, rval)
		a_Case : suspend ir_a_Case(p, st, target, bounded, rval)
		a_Select : suspend ir_a_Select(p, st, target, bounded, rval)
		a_Every : suspend ir_a_Every(p, st, target, bounded, rval)
		a_Sectionop : suspend ir_a_Sectionop(p, st, target, bounded, rval)
		a_Binop : suspend ir_a_Binop(p, st, target, bounded, rval)
		a_Unop : suspend ir_a_Unop(p, st, target, bounded, rval)
		a_If : suspend ir_a_If(p, st, target, bounded, rval)
		a_Initial : suspend ir_a_Initial(p, st, target, bounded, rval)
		a_Intlit : suspend ir_a_Intlit(p, st, target, bounded, rval)
		a_Reallit : suspend ir_a_Reallit(p, st, target, bounded, rval)
		a_Stringlit : suspend ir_a_Stringlit(p, st, target, bounded, rval)
		a_ProcDecl : suspend ir_a_ProcDeclNested(p, st, target, bounded, rval)
		a_ProcCode : suspend ir_a_ProcCode(p, st, target, bounded, rval)
		a_Repeat : suspend ir_a_Repeat(p, st, target, bounded, rval)
		a_Return : suspend ir_a_Return(p, st, target, bounded, rval)
		a_Nil : suspend ir_a_Nil(p, st, target, bounded, rval)
		a_Catch : suspend ir_a_Catch(p, st, target, bounded, rval)
		a_Fail : suspend ir_a_Fail(p, st, target, bounded, rval)
		a_Suspend : suspend ir_a_Suspend(p, st, target, bounded, rval)
		a_While : suspend ir_a_While(p, st, target, bounded, rval)
		a_With : suspend ir_a_With(p, st, target, bounded, rval)
		a_Create : suspend ir_a_Create(p, st, target, bounded, rval)
		a_Ident : suspend ir_a_Ident(p, st, target, bounded, rval)
		a_Continue : suspend ir_a_Continue(p, st, target, bounded, rval)
		a_Break : suspend ir_a_Break(p, st, target, bounded, rval)
		a_Yield : suspend ir_a_Yield(p, st, target, bounded, rval)
		a_ToBy : suspend ir_a_ToBy(p, st, target, bounded, rval)
		a_Mutual : suspend ir_a_Mutual(p, st, target, bounded, rval)
		a_Parallel : suspend ir_a_Parallel(p, st, target, bounded, rval)
		a_Compound : suspend ir_a_Compound(p, st, target, bounded, rval)
		a_ListConstructor : suspend ir_a_ListConstructor(p, st, target,
													   bounded, rval)
		a_ListComprehension : suspend ir_a_ListComprehension(p, st, target,
													   bounded, rval)
		a_Key : suspend ir_a_Key(p, st, target, bounded, rval)
		a_Local : suspend ir_a_Local(p, st, target, bounded, rval)
		a_Static : suspend ir_a_Static(p, st, target, bounded, rval)
		default : throw("unrecognized type", p)
	}
}

procedure ir_opfn(coord, lhs, lhsclsr, op, args, failLabel) {
	static neverfail
	\neverfail | {
		neverfail := list(3)
		neverfail[1] := set([ "#", "+", "-", "~", "^", "*", "." ])
		neverfail[2] := set([
			"+", "-", "*", "/", "%", "^",
			"++", "--", "**",
			"||", "|||",
			".", "&",
			":=", ":=:",
		])
		neverfail[3] := set([ ])
	}

	op.arity = *args | throw("no arity", op)
	if neverfail[op.arity].member(op.name) then {
		failLabel := nil
	}
	return ir_OpFunction(coord, lhs, lhsclsr, op.name, args, op.rval, failLabel)
}

procedure ir_init(p) {
	p.ir := ir_info()
	p.ir.start := ir_label(p, "start")
	p.ir.resume := ir_label(p, "resume")
	p.ir.success := ir_label(p, "success")
	p.ir.failure := ir_label(p, "failure")
	return p
}

procedure ir_init_loop(p, st, target, bounded, rval) {
	ir_init(p)
	p.ir.x := ir_loopinfo()
	/bounded & (p.ir.x.continueTmp := ir_tmploc(st))
	p.ir.x.nextlabel := ir_label(p, "next")
	p.ir.x.value := target
	p.ir.x.bounded := bounded
	p.ir.x.rval := rval
	p.ir.x.indirects := []
	return p
}

procedure ir_label(p, suffix) {
	return ir_Label(ir_naming(p, suffix))
}

procedure ir_naming(p, suffix) {
	# kludge
	/static T := table()
	/T[p] := *T
	return type(p).name() || "_" || T[p] || "_" || suffix
}

procedure ir_key(str) {
	local k

	static keytable
	/keytable := table()
	/keytable[str] := str
	return str
}

procedure ir_tmp(st) {
	st.tmp +:= 1
	return ir_Tmp("tmp" || st.tmp)
}

procedure ir_tmploc(st) {
	st.tmp +:= 1
	return ir_TmpLabel("loc" || st.tmp)
}

procedure ir_tmpclosure(st) {
	st.tmp +:= 1
	return ir_TmpClosure("closure" || st.tmp)
}

procedure ir_make_sentinel(L) {
	L.put(ir_init(a_NoOp()))
	L.push(ir_init(a_NoOp()))
	return L
}

procedure semantic_error(msg, coord) {
	%stderr.writes("At ", \coord, ": ")
	stop(msg)
}

procedure ast2ir(parse, flagList) {
	local p
	local k

	ir_declare_set := set([])
	ir_deferred := []

	while p := @parse do {
		suspend ir_outer(p)
		while k := ir_deferred.pop() do {
				suspend k
		}
	}
}

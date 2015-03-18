#  optimize.gd -- intermediate representation optimizer routines.

procedure optim_optimize(t, start) {
	local new
	local i

	new := t
	#%#% for unknown reasons, but definitely a bug,
	#%#% the following two lines cannot be reordered.
	new := optim_dead_assignment(new, start)
	new := optim_goto_elimination(new, start)
	new := optim_copy_propagation(new, start)
	return new
}

procedure optim_dead_assignment(t, start) {
	local k
	local chunk
	local new
	local c
	local i
	local pair

	new := table()
	every pair := !t & k := pair.key & chunk := pair.value do {
		c := []
		every i := !chunk do {
			case type(i) of {
				ir_MoveLabel |
				ir_IntLit |
				ir_NilLit |
				ir_RealLit |
				ir_StrLit |
				ir_NilLit |
				ir_Var |
				ir_Move : if /i.lhs then continue
			}
			c.put(i)
		}
		new[k] := c
	}
	return new
}

procedure optim_copy_propagation(t, start) {
	local rename
	local uses
	local defs
	local chunk
	local insn
	local newchunk
	local changes
	local k
	local i
	local pair

	repeat {
		changes := nil
		uses := table(0)
		defs := table(0)
		every chunk := !t do {
			every insn := !chunk do {
				optim_def_use(insn, uses, defs)
			}
		}
		rename := table()
		# every k := key(t) & chunk := t[k] do {
		every pair := !t & k := pair.key & chunk := pair.value do {
			newchunk := []
			every insn := !chunk do {
				if type(insn) === ir_Move &
						defs[insn.lhs] = 1 &
						uses[insn.rhs] = 1 then {
					rename[insn.lhs] := insn.rhs
					changes := 1
				} else {
					newchunk.put(insn)
				}
			}
			t[k] := newchunk
		}
		if /changes then break;
		every chunk := !t do {
			every i := 1 to *chunk do {
				chunk[i] := optim_rewrite(chunk[i], rename)
			}
		}
	}
	return t
}

procedure optim_def_use(p, uses, defs) {
	local sc
	case type(p) of {
		ir_Field : {
			defs[\p.lhs] +:= 1
			uses[p.expr] +:= 1
		}
		ir_Move :  {
			defs[\p.lhs] +:= 1
			uses[p.rhs] +:= 1
		}
		ir_MoveLabel :  {
			defs[\p.lhs] +:= 1
			uses[p.label] +:= 1
		}
		ir_Deref :  {
			defs[\p.lhs] +:= 1
			uses[p.value] +:= 1
		}
		ir_Goto :		{ uses[p.targetLabel] +:= 1 }
		ir_IndirectGoto :  { uses[p.targetTmpLabel] +:= 1 }
		ir_Var :		{ defs[\p.lhs] +:= 1 }
		ir_Key :		{ defs[\p.lhs] +:= 1 }
		ir_NilLit :		{ defs[\p.lhs] +:= 1 }
		ir_IntLit :		{ defs[\p.lhs] +:= 1 }
		ir_RealLit :	{ defs[\p.lhs] +:= 1 }
		ir_StrLit :		{ defs[\p.lhs] +:= 1 }
		ir_NoValue :	{ defs[\p.lhs] +:= 1 }
		ir_MakeClosure :{ defs[\p.lhs] +:= 1 }
		ir_Succeed :   {
			uses[p.expr] +:= 1
			uses[\p.resumeLabel] +:= 1
			}
		ir_Fail :  { }
		ir_ResumeValue :    {
			# uses[p.value] +:= 1
			uses[\p.failLabel] +:= 1
			defs[\p.lhs] +:= 1
			}
		ir_MakeList :  {
			defs[\p.lhs] +:= 1
			every uses[!p.valueList] +:= 1
			}
		ir_Call |
		ir_OpFunction: {
			defs[\p.lhs] +:= 1
			uses[p.fn] +:= 1
			uses[\p.failLabel] +:= 1
			every uses[!p.argList] +:= 1
			}
		ir_Select : {
			uses[\p.failLabel] +:= 1
			every sc := !p.caseList do {
				uses[\sc.lhs] +:= 1
				uses[\sc.rhs] +:= 1
				uses[sc.bodyLabel] +:= 1
			}
			}
		ir_Create :    {
			defs[\p.lhs] +:= 1
			uses[p.coexpLabel] +:= 1
			}
		ir_CoRet : {
			uses[p.value] +:= 1
			uses[p.resumeLabel] +:= 1
			}
		ir_CoFail :    { }
		ir_NoOp :    { }
		ir_EnterScope :    { }
		ir_ExitScope :    { }
		ir_Catch :   { uses[p.fn] +:= 1 }

		ir_Unreachable:{ }

		default :   { throw("unrecognized type", p) }
	}
}

procedure optim_rename(p, rename) {
	while rename.member(p) do {
		p := rename[p]
	}
	return p
}

procedure optim_rewrite(p, rename) {
	local i
	local sc

	case type(p) of {
		ir_Move :  {
			p.lhs := optim_rename(\p.lhs, rename);
			p.rhs := optim_rename(p.rhs, rename) }
		ir_MoveLabel :  {
			p.lhs := optim_rename(\p.lhs, rename);
			p.label := optim_rename(p.label, rename) }
		ir_Deref :  {
			p.lhs := optim_rename(\p.lhs, rename);
			p.value := optim_rename(p.value, rename) }
		ir_Goto :  { }
		ir_IndirectGoto :  { }
		ir_MakeClosure :   { p.lhs := optim_rename(\p.lhs, rename) }
		ir_Var :   { p.lhs := optim_rename(\p.lhs, rename) }
		ir_Key :   { p.lhs := optim_rename(\p.lhs, rename) }
		ir_IntLit :    { p.lhs := optim_rename(\p.lhs, rename) }
		ir_NilLit :    { p.lhs := optim_rename(\p.lhs, rename) }
		ir_RealLit :   { p.lhs := optim_rename(\p.lhs, rename) }
		ir_StrLit :    { p.lhs := optim_rename(\p.lhs, rename) }
		ir_NoValue :   { p.lhs := optim_rename(\p.lhs, rename) }
		ir_Succeed :   {
			p.expr := optim_rename(p.expr, rename);
			p.resumeLabel := optim_rename(\p.resumeLabel, rename); }
		ir_Fail :  { }
		ir_ResumeValue :    {
			# p.closure := optim_rename(p.closure, rename);
			p.failLabel := optim_rename(\p.failLabel, rename);
			p.lhs := optim_rename(\p.lhs, rename); }
		ir_MakeList :  {
			p.lhs := optim_rename(\p.lhs, rename);
			every i := 1 to *p.valueList do {
				p.valueList[i] := optim_rename(p.valueList[i], rename);
			}
			}
		ir_Field : {
			p.lhs := optim_rename(\p.lhs, rename)
			p.expr := optim_rename(p.expr, rename)
			}
		ir_Call |
		ir_OpFunction:{
			p.lhs := optim_rename(\p.lhs, rename);
			every i := 1 to *p.argList do {
				p.argList[i] := optim_rename(p.argList[i], rename);
			}
			p.failLabel := optim_rename(\p.failLabel, rename);
			}
		ir_Select :  {	#%#% untested
			p.failLabel := optim_rename(\p.failLabel, rename);
			every sc := !p.caseList do {
				sc.lhs := optim_rename(\sc.lhs, rename);
				sc.rhs := optim_rename(\sc.rhs, rename);
				sc.bodyLabel := optim_rename(\sc.bodyLabel, rename);
			}
		}
		ir_Create :    {
			p.lhs := optim_rename(\p.lhs, rename);
			p.coexpLabel := optim_rename(p.coexpLabel, rename); }
		ir_CoRet : {
			p.value := optim_rename(p.value, rename);
			p.resumeLabel := optim_rename(p.resumeLabel, rename);}
		ir_CoFail :    { }
		ir_Unreachable:{ }
		ir_NoOp :    { }
		ir_EnterScope :    { }
		ir_ExitScope :    { }
		ir_Catch :   { p.fn := optim_rename(p.fn, rename) }

		default :   { throw("unrecognized type", p) }
	}
	return p;
}

procedure optim_indirect_elimination(t) {
	local chunk
	local insn

	every chunk := !t & insn := chunk[-1] &
			type(insn) === ir_IndirectGoto &
			*insn.labelList = 1 do {
		chunk[-1] := ir_Goto(insn.coord, insn.labelList[1])
	}
}

procedure optim_goto_elimination(t, start) {
	local new
	local i

	optim_indirect_elimination(t)
	new := table()
	optim_goto_transitive(t, new, start)
	optim_test_elimination(new, start)
	optim_fallthrough(new)
	return new
}

procedure optim_fallthrough(new) {
	local lab
	local insn
	local chunk
	local refcount

	refcount := optim_refcount(new)

	every lab := (!new).key do {
		while chunk := \new[lab] &
				insn := chunk[-1] &
				type(insn) === ir_Goto &
				type(insn.targetLabel) === ir_Label &
				refcount[insn.targetLabel] = 1 do {
			\new[insn.targetLabel] | throw("/new", insn.targetLabel)
			new[lab] := chunk[1:-1] ||| new[insn.targetLabel]
			new.delete(insn.targetLabel)
		}
	}
}

procedure optim_refcount(t) {
	local refcount

	refcount := table(0)
	every optim_refcountX(refcount, !!t)
	return refcount
}


procedure optim_goto_chain(lab, t) {
	local chunk
	local seen

	seen := set([])
	while chunk := \t[lab] &
			type(chunk[1]) === ir_Goto &
			not seen.member(chunk[1].targetLabel) do {
		lab := chunk[1].targetLabel
		seen.insert(lab)
	}
	return lab
}

procedure optim_goto_transitive(t, new, lab) {
	local p
	local sc
	local i

	\t[lab] | throw("/t[lab]", t, lab)

	if \new[lab] then {
		return
	}
	new[lab] := t[lab]

	every p := !t[lab] do {
		case type(p) of {
			ir_CoFail |
			ir_Unreachable |
			ir_MakeList |
			ir_operator |
			ir_Fail |
			ir_Var |
			ir_IntLit |
			ir_NilLit |
			ir_RealLit |
			ir_StrLit |
			ir_NoValue |
			ir_Tmp |
			ir_TmpLabel |
			ir_Move |
			ir_MakeClosure |
			ir_NoOp |
			ir_EnterScope |
			ir_ExitScope |
			ir_Catch |
			ir_Deref : {
				# nothing
			}
			ir_Label : {  # ir_Label : ( value )
				throw("case ir_Label", p)
			}
			ir_IndirectGoto : {
				every i := 1 to *p.labelList do {
					p.labelList[i] := optim_goto_chain(p.labelList[i], t)
				}
			}
			ir_MoveLabel : {  # ir_MoveLabel : ( lhs label )
				p.label := optim_goto_chain(p.label, t)
				optim_goto_transitive(t, new, p.label)
			}
			ir_Goto : {  # ir_Goto : ( targetLabel )
				p.targetLabel := optim_goto_chain(p.targetLabel, t)
				optim_goto_transitive(t, new, p.targetLabel)
			}
			ir_Succeed : {  # ir_Succeed : ( expr resumeLabel )
				p.resumeLabel := optim_goto_chain(\p.resumeLabel, t)
				optim_goto_transitive(t, new, \p.resumeLabel)
			}
			ir_Field : {  # ir_Field : ( lhs expr field )
			}
			ir_Key : {}
			ir_Call |
			ir_OpFunction : {  # ir_OpFunction : ( lhs fn argList failLabel )
				p.failLabel := optim_goto_chain(\p.failLabel, t)
				optim_goto_transitive(t, new, \p.failLabel)
			}
			ir_ResumeValue : {  # ir_ResumeValue : ( lhs value failLabel )
				p.failLabel := optim_goto_chain(\p.failLabel, t)
				optim_goto_transitive(t, new, \p.failLabel)
			}
			ir_Select : {  # ir_Select : ( caselist )
				p.failLabel := optim_goto_chain(\p.failLabel, t)
				optim_goto_transitive(t, new, \p.failLabel)
				every sc := !p.caseList do {
					# ir_SelectCase : ( kind lhs rhs bodyLabel )
					sc.bodyLabel := optim_goto_chain(sc.bodyLabel, t)
					optim_goto_transitive(t, new, sc.bodyLabel)
				}
			}
			ir_Create : {  # ir_Create : ( lhs location )
				p.coexpLabel := optim_goto_chain(p.coexpLabel, t)
				optim_goto_transitive(t, new, p.coexpLabel)
			}
			ir_CoRet : {  # ir_CoRet : ( value resumeLabel )
				p.resumeLabel := optim_goto_chain(p.resumeLabel, t)
				optim_goto_transitive(t, new, p.resumeLabel)
			}
			default : throw("unrecognized type", p)
		}
	}
}

procedure optim_refcountX(refcount, p) {
	case type(p) of {
		ir_CoFail |
		ir_Unreachable |
		ir_MakeList |
		ir_operator |
		ir_Fail |
		ir_Var |
		ir_IntLit |
		ir_NilLit |
		ir_RealLit |
		ir_StrLit |
		ir_NoValue |
		ir_Tmp |
		ir_MakeClosure |
		ir_TmpLabel : {
			# nothing
			}
		ir_Deref : {
			refcount[p.value] +:= 1
			}
		ir_Label : {  # ir_Label : ( value )
			}
		ir_Move : {  # ir_Move : ( lhs rhs )
			refcount[p.rhs] +:= 1
			}
		ir_MoveLabel : {  # ir_MoveLabel : ( lhs label )
			refcount[p.label] +:= 1
			}
		ir_Goto : {  # ir_Goto : ( targetLabel )
			refcount[p.targetLabel] +:= 1
			}
		ir_IndirectGoto : {  # ir_IndirectGoto : ( targetTmpLabel )
			refcount[p.targetTmpLabel] +:= 1
			}
		ir_Succeed : {  # ir_Succeed : ( expr resumeLabel )
			refcount[p.resumeLabel] +:= 1
			}
		ir_Field : {  # ir_Field : ( lhs expr field )
			}
		ir_Key : {}
		ir_Call |
		ir_OpFunction : {  # ir_OpFunction : ( lhs fn argList failLabel )
			refcount[p.failLabel] +:= 1
			}
		ir_ResumeValue : {  # ir_ResumeValue : ( lhs value failLabel )
			refcount[\p.failLabel] +:= 1
			}
		ir_Select: {
			refcount[\p.failLabel] +:= 1
			every refcount[(!p.caseList).bodyLabel] +:= 1
			}
		ir_Create : {  # ir_Create : ( lhs coexpLabel )
			refcount[p.coexpLabel] +:= 1
			}
		ir_CoRet : {  # ir_CoRet : ( value resumeLabel )
			refcount[p.resumeLabel] +:= 1
			}
		ir_NoOp :    { }
		ir_EnterScope :    { }
		ir_ExitScope :    { }
		ir_Catch :   { }
		default : throw("unrecognized type", p)
	}
}

procedure optim_test_elimination(t, start) {
	local chunk
	local insn

	every chunk := !t &
		    type(chunk[-1]) === ir_Goto &
		    type(chunk[-2]) === (ir_Call | ir_OpFunction | ir_Select) &
			chunk[-1].targetLabel === \chunk[-2].failLabel do {
		chunk[-2].failLabel := nil
	}
}

procedure optim(irgen, flagList) {
	local p
	local T
	local L
	local i

	while p := @irgen do {
		case type(p) of {
			ir_Record |
			ir_Global |
			ir_Initial : {
				suspend p
			}
			ir_Function : {
				if match("-O", !flagList) then {
					T := table()
					every i := !p.codeList do {
						T[i.label] := i.insnList
					}
					T := optim_optimize(T, p.codeStart)
					L := []
					every i := (!T).key do {
						L.put(ir_chunk(i, T[i]))
					}
					p.codeList := L
				}
				suspend p
			}
		default: throw("unrecognized type", p)
		}
	}
}

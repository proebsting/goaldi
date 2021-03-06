#  ir.gd -- data structures for the Goaldi intermediate representation.

record ir_Record(coord, name, extendsRec, extendsPkg, fieldList, namespace)
record ir_Global(coord, name, fn, namespace)
record ir_Initial(coord, fn, namespace)
record ir_Function(coord, name, paramList, accumulate,
	localList, staticList, unboundList, codeList, codeStart,
	parent, namespace, tempCount)
record ir_chunk(label, insnList)

record ir_NoOp(coord, comment)
record ir_Catch(coord, lhs, fn)
record ir_EnterScope(coord, nameList, dynamicList, scope, parentScope)
record ir_ExitScope(coord, nameList, dynamicList, scope)

record ir_Tmp(name)
record ir_TmpLabel(name)
record ir_TmpClosure(name)
record ir_Label(value)

record ir_Var(coord, lhs, name, namespace, scope, rval)
record ir_Key(coord, lhs, name, scope, rval)
record ir_IntLit(coord, lhs, val)
record ir_NilLit(coord, lhs)
record ir_RealLit(coord, lhs, val)
record ir_StrLit(coord, lhs, len, val)	# UTF-8 encoded string

record ir_operator(name, arity, rval)

record ir_MakeClosure(coord, lhs, name)
record ir_Move(coord, lhs, rhs)
record ir_MoveLabel(coord, lhs, label)
record ir_Deref(coord, lhs, value)
record ir_MakeList(coord, lhs, valueList)
record ir_Field(coord, lhs, expr, field, rval)
record ir_OpFunction(coord, lhs, lhsclosure, fn, argList, rval, failLabel)
record ir_Call(coord, lhs, lhsclosure, fn, argList, nameList, failLabel, scope)
record ir_ResumeValue(coord, lhs, lhsclosure, closure, failLabel)

record ir_Goto(coord, targetLabel)
record ir_IndirectGoto(coord, targetTmpLabel, labelList)
record ir_Succeed(coord, expr, resumeLabel)
record ir_Fail(coord)

record ir_Create(coord, lhs, coexpLabel, scope)
record ir_CoRet(coord, value, resumeLabel)
record ir_CoFail(coord)

record ir_Select(coord, caseList, failLabel)
record ir_SelectCase(coord, kind, lhs, rhs, bodyLabel)
record ir_NoValue(coord, lhs)

record ir_Unreachable(coord)

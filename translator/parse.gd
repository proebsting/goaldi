#  parse.gd -- LL(1) parser that maps tokens to abstract syntax trees.

global parse_tok        # This is the co_expression that represents the scanner
global parse_tok_rec    # This is the current token record

record parse_named(nameList, exprList)

procedure parse(lex) {
	local d

	parse_tok := lex
	parse_tok_rec := @parse_tok
	suspend parse_program()
}

#  CASE  expr  OF  LBRACE  cclause {  SEMICOL  cclause }  RBRACE
procedure parse_do_case() {
	local e
	local body
	local element
	local dflt
	local coord

	coord := parse_tok_rec.coord
	parse_eat_token()
	e := parse_expr()
	parse_match_token(lex_OF)
	parse_match_token(lex_LBRACE)
	body := []
	element := parse_cclause()
	if element.expr === lex_DEFAULT then {
		dflt := element.body
	} else {
		body.put(element)
	}
	while parse_tok_rec === lex_SEMICOL do {
		parse_match_token(lex_SEMICOL)
		if parse_tok_rec === lex_RBRACE then {
			break
		}
		element := parse_cclause()
		if element.expr === lex_DEFAULT then {
			(/dflt := element.body) | parse_error("multiple default clauses")
		} else {
			body.put(element)
		}
	}
	parse_match_token(lex_RBRACE)
	return a_Case(e, body, dflt, coord)
}

procedure parse_cclause() {
# case-clause
	local e
	local body
	local coord
	static cclause_set
	/cclause_set := set([lex_AT, lex_BACKSLASH, lex_BANG,
		lex_BAR, lex_BREAK, lex_CARET, lex_CASE,
		lex_CONCAT, lex_CREATE, lex_DIFF,
		lex_DOT, lex_EQUIV, lex_EVERY, lex_FAIL,
		lex_IDENT, lex_IF, lex_INTER, lex_INTLIT,
		lex_LBRACE, lex_LBRACK, lex_LCONCAT, lex_LPAREN,
		lex_MINUS, lex_NEQUIV, lex_CONTINUE, lex_NMEQ,
		lex_NMNE, lex_NOT, lex_PLUS, lex_QMARK,
		lex_REALLIT, lex_REPEAT, lex_RETURN, lex_SEQ,
		lex_SLASH, lex_SNE, lex_STAR, lex_STRINGLIT,
		lex_SUSPEND, lex_UNION,
		lex_SLASHSLASH,
		lex_MOD,
		lex_NIL,
		lex_YIELD,
		lex_LAMBDA,
		lex_WHILE,
		lex_WITH,
		lex_SELECT,
		lex_CATCH,
		lex_LOCAL,
		lex_STATIC,
		lex_LCOMP,
		lex_PROCEDURE])
	if parse_tok_rec === lex_DEFAULT then {
		e := lex_DEFAULT
		parse_eat_token()
		coord := parse_tok_rec.coord
		parse_match_token(lex_COLON)
		body := parse_expr()
	} else if cclause_set.member(parse_tok_rec) then {
		e := parse_expr()
		coord := parse_tok_rec.coord
		parse_match_token(lex_COLON)
		body := parse_expr()
	} else {
		parse_error("\""||parse_tok_rec.str||"\": invalid case clause")
	}
	return a_Cclause(e, body, coord)
}

procedure parse_do_select() {
#  SELECT  LBRACE  selcase {  SEMICOL  selcase }  RBRACE
	local caseList
	local element
	local dflt
	local coord

	coord := parse_tok_rec.coord
	parse_eat_token()
	parse_match_token(lex_LBRACE)
	caseList := []
	repeat {
		if parse_tok_rec === lex_RBRACE then {
			break
		}
		element := parse_selcase()
		if element.kind === "default" then {
			(/dflt := element) | parse_error("more than one default clause")
		} else {
			caseList.put(element)
		}
		parse_match_token(parse_tok_rec === lex_SEMICOL) | break
	}
	parse_match_token(lex_RBRACE)
	return a_Select(caseList, dflt, coord)
}

procedure parse_selcase() {
# select-case:  select-condition : body
	local sc

	sc := a_SelectCase()
	parse_selectby(sc) | parse_error("Malformed select condition")
	sc.coord := parse_tok_rec.coord
	parse_match_token(lex_COLON)
	sc.body := parse_expr()
	return sc
}

procedure parse_selectby(sc) {
# select-condition:  default  |  expr := @expr  |  expr @: expr
# sets sc.kind,left,right; fails on parse error
	local e
	local r

	if parse_tok_rec === lex_DEFAULT then {
		parse_eat_token()
		sc.kind := "default"
		return
	}
	e := parse_expr()
	type(e) === a_Binop | return fail
	sc.left := e.left
	case e.op of {
		default: {
			return fail
		}
		"@:": {
			sc.kind := "send"
			sc.right := e.right
			return
		}
		":=": {
			r := e.right
			type(r) === a_Unop | return fail
			r.op === "@" | return fail
			sc.kind := "receive"
			sc.right := r.operand
			return
		}
	}
}

procedure parse_compound() {
#  nexpr {  SEMICOL  nexpr }
	local l
	local e

	l := [parse_nexpr()]
	while parse_tok_rec === lex_SEMICOL do {
		parse_eat_token()
		e := parse_nexpr()
		l.put(e)
	}
	return l
}

procedure parse_decl() {
	case parse_tok_rec of {
		lex_RECORD    : return parse_do_record()
		lex_PROCEDURE : return parse_do_proc()
		lex_GLOBAL    : return parse_do_global()
		lex_INITIAL   : return parse_do_initial()
		default   : parse_error("Expecting parse_declaration")
	}
}

procedure parse_do_every() {
#  EVERY  expr [  DO  expr ]
	local e
	local body
	local coord
	local id
	coord := parse_tok_rec.coord
	parse_match_token(lex_EVERY)
	if parse_tok_rec === lex_COLON then {
		parse_eat_token()
		id := parse_match_token(lex_IDENT)
	}
	e := parse_expr()
	if parse_tok_rec === lex_DO then {
		parse_match_token(lex_DO)
		body := parse_expr()
	} else {
		body := nil
	}
	return a_Every(e, body, id, coord)
}

procedure parse_expr() {
#  expr1x {  ANDAND  expr1x }
	local ret
	local L

	ret := parse_expr1x()
	if parse_tok_rec === lex_ANDAND then {
		L := [ret]
		while parse_tok_rec === lex_ANDAND do {
			parse_eat_token()
			ret := parse_expr1x()
			L.put(ret)
		}
		return a_Parallel(L, ret.coord)
	} else {
		return ret
	}
}

procedure parse_expr1x() {
#  expr1 {  AND  expr1 }
	local ret
	local op
	local right

	ret := parse_expr1()
	while parse_tok_rec === lex_AND do {
		op := parse_eat_token()
		right := parse_expr1()
		ret := a_Binop(op, ret, right, ret.coord)
	}
	return ret
}

procedure parse_expr1() {
	local ret
	local op
	local right
	local coord
	static expr1_set

	/expr1_set := set([lex_ASSIGN, lex_AUGAND, lex_AUGCARET,
		lex_AUGCONCAT, lex_AUGDIFF, lex_AUGEQUIV,
		lex_AUGINTER, lex_AUGLCONCAT, lex_AUGMINUS,
		lex_AUGMOD, lex_AUGNEQUIV, lex_AUGNMEQ, lex_AUGNMGE,
		lex_AUGNMGT, lex_AUGNMLE, lex_AUGNMLT, lex_AUGNMNE,
		lex_AUGPLUS, lex_AUGSEQ, lex_AUGSGE,
		lex_AUGSGT, lex_AUGSLASH, lex_AUGSLE, lex_AUGSLT,
		lex_AUGSNE, lex_AUGSTAR, lex_AUGUNION,
		lex_AUGSLASHSLASH,
		lex_REVASSIGN, lex_REVSWAP, lex_SWAP, lex_ATCOLON])
	#  expr2 {  expr1op  expr1 } (Right Associative)
	ret := parse_expr2()
	while expr1_set.member(parse_tok_rec) do {
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		right := parse_expr1()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parse_expr10() {
	local op
	local operand
	local tmp_tok
	local coord
	static expr10_set1
	static expr10_set2
	static expr10_set3

	/expr10_set1 := set([lex_BREAK, lex_CASE, lex_CREATE,
		lex_EVERY, lex_FAIL, lex_IDENT,
		lex_IF, lex_INTLIT, lex_LBRACE, lex_LBRACK,
		lex_LPAREN, lex_CONTINUE, lex_REALLIT, lex_REPEAT,
		lex_RETURN, lex_STRINGLIT, lex_SUSPEND,
		lex_YIELD,
		lex_MOD,
		lex_NIL,
		lex_LAMBDA,
		lex_WHILE,
		lex_WITH,
		lex_SELECT,
		lex_CATCH,
		lex_LOCAL,
		lex_CARET,
		lex_STATIC,
		lex_LCOMP,
		lex_PROCEDURE ])
	/expr10_set2 := set([lex_AT, lex_NOT, lex_BAR, lex_DOT, lex_BANG,
		lex_PLUS, lex_STAR, lex_SLASH, 
		lex_MINUS, 
		lex_QMARK,
		lex_BACKSLASH])
	/expr10_set3 := set([lex_CONCAT, lex_LCONCAT, lex_UNION, lex_INTER,
		lex_SLASHSLASH,
		lex_DIFF])

	if expr10_set1.member(parse_tok_rec) then {
		return parse_expr11a()
	} else if expr10_set2.member(parse_tok_rec) then {
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		operand := parse_expr10()
		case (op) of {
			"|":        return a_RepAlt(operand, coord)
			"not":      return a_Not(operand, coord)
			default:    return a_Unop(op, operand, coord)
		}
	} else if expr10_set3.member(parse_tok_rec) then {
		tmp_tok := parse_tok_rec
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		operand := parse_expr10()
		case tmp_tok of {
			lex_SLASHSLASH:	return a_Unop("/", operand, coord)
			lex_CONCAT:		return a_RepAlt(operand, coord)
			lex_LCONCAT:	return a_RepAlt(operand, coord)
			lex_UNION:		return a_Unop("+", operand, coord)
			lex_INTER:		return a_Unop("*", a_Unop("*",operand,coord), coord)
			lex_DIFF:		return a_Unop("-", a_Unop("-",operand,coord), coord)
			default:		throw("unrecognized token", tmp_tok)
		}
	} else {
		parse_error("\""||parse_tok_rec.str||"\": expression expected")
	}
}

procedure parse_braced() {
	local e
	local coord

	coord := parse_tok_rec.coord
	parse_match_token(lex_LBRACE)
	e := parse_compound()
	parse_match_token(lex_RBRACE)
	return a_Compound(e, coord)
}

procedure parse_expr11() {
	local e
	local id
	local coord
	local ns
	case parse_tok_rec of {
		lex_INTLIT  |
		lex_REALLIT |
		lex_STRINGLIT   :   # literal
			return parse_literal()
		lex_CATCH   :
			return parse_do_catch()
		lex_NIL    : {
			coord := parse_tok_rec.coord
			parse_eat_token()
			return a_Nil(coord)
			}
		lex_FAIL    : {
			coord := parse_tok_rec.coord
			parse_eat_token()
			return a_Fail(coord)
			}
		lex_RETURN  |
		lex_SUSPEND :   # return
			return parse_do_return()
		lex_IF  :   # if
			return parse_do_if()
		lex_CASE    :   # case
			return parse_do_case()
		lex_SELECT  :   # select
			return parse_do_select()
		lex_WHILE   :   # while
			return parse_do_while()
		lex_EVERY   :   # every
			return parse_do_every()
		lex_REPEAT  :   # repeat
			return parse_do_repeat()
		lex_CREATE  :   # CREATE  expr
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			e := parse_expr()
			return a_Create(e, coord)
			}
		lex_IDENT   :   # IDENT
			{
			coord := parse_tok_rec.coord
			id := parse_eat_token()
			ns := nil
			if parse_tok_rec === lex_COLONCOLON then {
				parse_eat_token()
				ns := id
				id := parse_match_token(lex_IDENT)
			}
			return a_Ident(id, ns, coord)
			}
		lex_CARET   |
		lex_LOCAL   :
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			id := parse_match_token(lex_IDENT)
			return a_Local(id, coord)
			}
		lex_STATIC   :   # STATIC
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			id := parse_match_token(lex_IDENT)
			return a_Static(id, coord)
			}
		lex_CONTINUE    :   # CONTINUE
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			if parse_tok_rec === lex_COLON then {
				parse_eat_token()
				id := parse_match_token(lex_IDENT)
			}
			return a_Continue(id, coord)
			}
		lex_YIELD   :   # YIELD  [ : IDENT ] nexpr
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			if parse_tok_rec === lex_COLON then {
				parse_eat_token()
				id := parse_match_token(lex_IDENT)
			}
			e := parse_nexpr()
			return a_Yield(e, id, coord)
			}
		lex_BREAK   :   # BREAK  nexpr
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			if parse_tok_rec === lex_COLON then {
				parse_eat_token()
				id := parse_match_token(lex_IDENT)
			}
			return a_Break(id, coord)
			}
		lex_LPAREN  :   # LPAREN  exprlist  RPAREN
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			e := parse_exprlist()
			parse_match_token(lex_RPAREN)
			return a_Mutual(e, coord)
			}
		lex_LBRACE  :   # LBRACE  compound  RBRACE
			{
			return parse_braced()
			}
		lex_LBRACK  :   # LBRACK  exprlist  RBRACK
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			e := parse_exprlist()
			parse_match_token(lex_RBRACK)
			return a_ListConstructor(e, coord)
			}
		lex_LCOMP  :   # LCOMP  expr  RCOMP
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			e := parse_expr()
			parse_match_token(lex_RCOMP)
			return a_ListComprehension(e, coord)
			}
		lex_MOD :   # MOD IDENT
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			return a_Key(parse_match_token(lex_IDENT), coord)
			}
		lex_PROCEDURE : {
				return parse_do_proc("noident")
			}
		lex_LAMBDA : {
				return parse_do_lambda()
			}
		lex_WITH : {
				return parse_do_with()
			}
		default : {
			parse_error("Expecting parse_expression")
		}
	}
}

procedure parse_expr11a() {
	local left
	local right
	local op
	local oprec
	static expr11_set

	/expr11_set := set([lex_DOT, lex_LBRACE, lex_LBRACK, lex_LPAREN, lex_LCOMP])
	#  expr11 {  expr11suffix }
	left := parse_expr11()
	while expr11_set.member(parse_tok_rec) do {
		left := parse_expr11suffix(left)
	}
	return left
}

procedure parse_expr11suffix(lhs) {
	local left
	local right
	local l
	local op
	local id
	local coord
	local x
	local idcoord
	local lefts
	local rights
	static expr11suffix_set1
	static expr11suffix_set2

	/expr11suffix_set1 := set([lex_AT, lex_BACKSLASH, lex_BANG,
		lex_BAR, lex_BREAK, lex_CARET, lex_CASE,
		lex_COMMA, lex_CONCAT, lex_CREATE,
		lex_DIFF, lex_DOT, lex_EQUIV,
		lex_EVERY, lex_FAIL, lex_IDENT, lex_IF,
		lex_INTER, lex_INTLIT, lex_LBRACE, lex_LBRACK,
		lex_LCONCAT, lex_LPAREN, lex_MINUS,
		lex_NEQUIV, lex_CONTINUE, lex_NMEQ, lex_NMNE,
		lex_NOT, lex_PLUS, lex_QMARK, lex_REALLIT,
		lex_REPEAT, lex_RETURN, lex_SEQ, lex_SLASH,
		lex_SNE, lex_STAR, lex_STRINGLIT,
		lex_SUSPEND, lex_UNION,
		lex_SLASHSLASH,
		lex_MOD,
		lex_NIL,
		lex_ANDAND,
		lex_YIELD,
		lex_LAMBDA,
		lex_WHILE,
		lex_WITH,
		lex_SELECT,
		lex_CATCH,
		lex_LOCAL,
		lex_STATIC,
		lex_LCOMP,
		lex_PROCEDURE ])
	/expr11suffix_set2 := set([lex_COLON, lex_MCOLON, lex_PCOLON])

	case parse_tok_rec of {
		lex_LBRACE  :   # LBRACE [  { e : e [ COMMA ] ]  RBRACE
			{
			lefts := []
			rights := []
			coord := parse_tok_rec.coord
			parse_eat_token()
			if parse_tok_rec ~=== lex_RBRACE then {
				while parse_tok_rec ~=== lex_RBRACE do {
					left := parse_expr()
					parse_match_token(lex_COLON)
					right := parse_expr()
					lefts.put(left)
					rights.put(right)
					if parse_tok_rec === lex_COMMA then {
						parse_eat_token()
					} else {
						break
					}
				}
			}
			parse_match_token(lex_RBRACE)
			return a_Paired(lhs, lefts, rights, coord)
			}
		lex_LPAREN  :   # LPAREN  exprlist  RPAREN
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			l := parse_named_exprlist()
			parse_match_token(lex_RPAREN)
				if type(l) === parse_named then {
					return a_Call(lhs, a_Arglist(l.exprList, l.nameList), coord)
				} else {
					return a_Call(lhs, a_Arglist(l, nil), coord)
				}
			}
		lex_DOT :   # DOT  IDENT
			{
			coord := parse_tok_rec.coord
			x := parse_eat_token()
			idcoord := parse_tok_rec.coord
			id := parse_match_token(lex_IDENT)
			return a_Field(lhs, a_Ident(id, nil, idcoord), coord)
			}
		lex_LBRACK  :   # LBRACK  expr [  sectop  expr ]  RBRACK
			{
			coord := parse_tok_rec.coord
			parse_eat_token()
			left := parse_nexpr()
			if expr11suffix_set2.member(parse_tok_rec) then {
				coord := parse_tok_rec.coord
				op := "[" || parse_eat_token() || "]"
				right := parse_expr()
				lhs := a_Sectionop(op, lhs, left, right, coord)
			} else {
				lhs := a_Binop("[]", lhs, left, coord)
				while parse_tok_rec === lex_COMMA do {
					parse_match_token(lex_COMMA)
					left := parse_nexpr()
					lhs := a_Binop("[]", lhs, left, coord)
				}
			}
			parse_match_token(lex_RBRACK)
			return lhs
			}
		default : {
			parse_error("Malformed argument list")
		}
	}
}

procedure parse_expr2() {
#  expr3 {  TO  expr2 [  BY  expr3 ] }
	local e1
	local e2
	local e3
	local ret
	local coord

	e1 := parse_expr3()
	e2 := nil
	e3 := nil
	ret := e1
	while parse_tok_rec === lex_TO do {
	coord := parse_tok_rec.coord
		parse_eat_token()
		e2 := parse_expr3()
		if parse_tok_rec === lex_BY then {
			parse_match_token(lex_BY)
			e3 := parse_expr3()
		}
		ret := a_ToBy(e1, e2, e3, coord)
		e1 := ret
	}
	return ret
}

procedure parse_expr3() {
#  expr3a {  TILDEBAR  expr3 }
	local ret
	local a

	ret := parse_expr3a()
	while parse_tok_rec === lex_TILDEBAR do {
		/a := a_ExcAlt([ret], parse_tok_rec.coord)
		parse_eat_token()
		a.eList.put(parse_expr3a())
	}
	ret := \a
	return ret
}

procedure parse_expr3a() {
#  expr4 {  BAR  expr3 }
	local ret
	local a

	ret := parse_expr4()
	while parse_tok_rec === lex_BAR do {
		/a := a_Alt([ret], parse_tok_rec.coord)
		parse_eat_token()
		a.eList.put(parse_expr4())
	}
	ret := \a
	return ret
}

procedure parse_expr4() {
	local ret
	local op
	local right
	local coord
	static expr4_set
	/expr4_set := set([lex_EQUIV, lex_NEQUIV, lex_NMEQ, lex_NMGE,
		lex_NMGT, lex_NMLE, lex_NMLT, lex_NMNE, lex_SEQ,
		lex_SGE, lex_SGT, lex_SLE, lex_SLT, lex_SNE])
	#  expr5 {  expr4op  expr4 }
	ret := parse_expr5()

	while expr4_set.member(parse_tok_rec) do {
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		right := parse_expr5()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parse_expr5() {
#  expr6 {  expr5op  expr5 }
	local ret
	local right
	local op
	local coord

	ret := parse_expr6()
	while parse_tok_rec === lex_CONCAT | parse_tok_rec === lex_LCONCAT do {
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		right := parse_expr6()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}


procedure parse_expr6() {
#  expr7 {  expr6op  expr6 }
	local ret
	local op
	local right
	local coord
	static expr6_set

	/expr6_set := set([lex_DIFF, lex_MINUS, lex_PLUS, lex_UNION])

	ret := parse_expr7()
	while expr6_set.member(parse_tok_rec) do {
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		right := parse_expr7()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parse_expr7() {
#  expr8 {  expr7op  expr7 }
	local ret
	local op
	local right
	local coord
	static expr7_set

	/expr7_set := set([lex_INTER, lex_MOD, lex_SLASH, lex_STAR, lex_SLASHSLASH])

	ret := parse_expr8()

	while expr7_set.member(parse_tok_rec) do {
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		right := parse_expr8()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parse_expr8() {
#  expr9 {  CARET  expr8 }  (Right Associative)
	local ret
	local op
	local right
	local coord

	ret := parse_expr9()
	while parse_tok_rec === lex_CARET do {
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		right := parse_expr8()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parse_expr9() {
#  expr10 {  expr9op  expr9 }
	local ret
	local op
	local right
	local coord

	ret := parse_expr10()
	while parse_tok_rec === ( lex_AT | lex_BACKSLASH | lex_BANG ) do {
		coord := parse_tok_rec.coord
		op := parse_eat_token()
		right := parse_expr10()
		if op == "\\" then {
			ret := a_Limitation(ret, right, coord)
		} else {
			ret := a_Binop(op, ret, right, coord)
		}
	}
	return ret
}

procedure new_parse_exprlist() {
	local L
	local e
	#  nexpr {  COMMA  nexpr }
	#  [ expr { COMMA expr } [ COMMA ] ]
	static expr_set

	/expr_set := set([lex_AT, lex_BACKSLASH, lex_BANG, lex_BAR,
		lex_BREAK, lex_CARET, lex_CASE, lex_CONCAT, lex_CREATE,
		lex_DIFF, lex_DOT, lex_EQUIV, lex_EVERY,
		lex_FAIL, lex_IDENT, lex_IF, lex_INTER, lex_INTLIT,
		lex_LBRACE, lex_LBRACK, lex_LCONCAT, lex_LPAREN,
		lex_MINUS, lex_NEQUIV, lex_CONTINUE, lex_NMEQ, lex_NMNE,
		lex_NOT, lex_PLUS, lex_QMARK, lex_REALLIT, lex_REPEAT,
		lex_RETURN, lex_SEQ, lex_SLASH, lex_SNE, lex_STAR,
		lex_STRINGLIT, lex_SUSPEND, lex_UNION,
		lex_WHILE,
		lex_WITH,
		lex_MOD,
		lex_NIL,
		lex_SLASHSLASH,
		lex_ANDAND,
		lex_YIELD,
		lex_LAMBDA,
		lex_SELECT,
		lex_CATCH,
		lex_LOCAL,
		lex_STATIC,
		lex_LCOMP,
		lex_PROCEDURE ])

	L := []
	while expr_set.member(parse_tok_rec) do {
		e := parse_expr()
		L.put(e)
		if (parse_tok_rec === lex_COMMA) then {
			parse_eat_token()
		} else {
			break
		}
	}
	return L
}

procedure old_parse_exprlist() {
	local l
	local e

	e := parse_nexpr()
	if \e | (parse_tok_rec === lex_COMMA) then {
		l := [ e ]
	} else {
		l := []
	}
	while parse_tok_rec === lex_COMMA do {
		parse_eat_token()
		e := parse_nexpr()
		l.put(e)
	}
	return l
}

procedure parse_named_exprlist0(id) {
	local e
	local N
	local E
	parse_match_token(lex_COLON)
	N := []
	E := []
	e := parse_expr()
	N.put(id.id)
	E.put(e)

	while parse_tok_rec === lex_COMMA do {
		parse_eat_token()
		if parse_tok_rec ~=== lex_IDENT then {
			break
		}
		id := parse_match_token(lex_IDENT)
		parse_match_token(lex_COLON)
		e := parse_nexpr()
		N.put(id)
		E.put(e)
	}
	return parse_named(N, E)
}

procedure parse_named_exprlistX() {
	local l
	local e

	e := parse_nexpr()
	if \e & type(e) === a_Ident & parse_tok_rec === lex_COLON then {
		return parse_named_exprlist0(e)
	} else if \e | (parse_tok_rec === lex_COMMA) then {
		l := [ e ]
	} else {
		l := []
	}
	while parse_tok_rec === lex_COMMA do {
		parse_eat_token()
		e := parse_nexpr()
		l.put(e)
	}
	if /l[-1] then {
		l := l[1:-1]
	}
	return l
}

procedure parse_named_exprlist() {
	local L
	local N
	local e

	L := []
	N := []
	repeat {
		e := parse_nexpr()
		if parse_tok_rec === lex_COMMA then {
			L.put(e)
		} else if \e then {
			if type(e) === a_Ident & parse_tok_rec === lex_COLON then {
				N.put(e.id)
				parse_match_token(lex_COLON)
				e := parse_expr()
				L.put(e)
			} else {
				L.put(e)
			}
		} else {
		}
		if parse_tok_rec ~=== lex_COMMA then {
			break
		}
		parse_match_token(lex_COMMA)
	}

	if *N > 0 then {
		return parse_named(N, L)
	} else {
		return L
	}
}

procedure parse_exprlist() {
	local l
	local e

	e := parse_nexpr()
	if \e | (parse_tok_rec === lex_COMMA) then {
		l := [ e ]
	} else {
		l := []
	}
	while parse_tok_rec === lex_COMMA do {
		parse_eat_token()
		e := parse_nexpr()
		l.put(e)
	}
	if /l[-1] then {
		l := l[1:-1]
	}
	return l
}

#  PACKAGE  name, coord
procedure parse_do_package() {
	local coord
	local coord2
	local id

	coord := parse_tok_rec.coord
	parse_match_token(lex_PACKAGE)

	coord2 := parse_tok_rec.coord
	id := a_Ident(parse_match_token(lex_IDENT), nil, coord2)

	return a_Package(id, coord)
}


procedure parse_idlist() {
#  IDENT {  COMMA  IDENT }
	local l
	local id
	local coord

	coord := parse_tok_rec.coord

	l := [a_Ident(parse_match_token(lex_IDENT), nil, coord)]
	while parse_tok_rec === lex_COMMA do {
		parse_eat_token()
		coord := parse_tok_rec.coord
		if parse_tok_rec ~=== lex_IDENT then {
			break
		}
		id := a_Ident(parse_match_token(lex_IDENT), nil, coord)
		l.put(id)
	}
	return l
}

procedure parse_do_if() {
#  IF  expr  THEN  expr [  ELSE  expr ]
	local ex
	local theparse_nexpr
	local elseparse_expr
	local coord

	coord := parse_tok_rec.coord
	parse_match_token(lex_IF)
	ex := parse_expr()
	parse_match_token(lex_THEN)
	theparse_nexpr := parse_expr()
	elseparse_expr := nil
	if parse_tok_rec === lex_ELSE then {
		parse_eat_token()
		elseparse_expr := parse_expr()
	}
	return a_If(ex, theparse_nexpr, elseparse_expr, coord)
}

#  GLOBAL  idlist, coord
procedure parse_do_global() {
	local coord
	local e
	local id

	coord := parse_tok_rec.coord
	parse_match_token(lex_GLOBAL)
	id := parse_match_token(lex_IDENT)
	if parse_tok_rec === lex_ASSIGN then {
		parse_match_token(lex_ASSIGN)
		e := parse_expr()
		e := a_Binop(":=", a_Ident(id, nil, coord), e)
		e := a_ProcCode(e)
	}
	return a_Global(id, e, coord)
}

procedure parse_do_initial() {
# INITIAL expr
	local coord
	local e

	coord := parse_tok_rec.coord
	parse_match_token(lex_INITIAL)
	e := parse_braced()
	return a_Initial(a_ProcCode(e), coord)
}

procedure parse_literal() {
	local coord
	coord := parse_tok_rec.coord
	case parse_tok_rec of {
		lex_INTLIT  :   # INTLIT
			return a_Intlit(integer(parse_eat_token()), coord)
		lex_REALLIT :   # REALLIT
			return a_Reallit(number(parse_eat_token()), coord)
		lex_STRINGLIT   :   # STRINGLIT
			return a_Stringlit(parse_eat_token(), coord)
		default :
			parse_error("Expecting parse_literal")
	}
}

procedure parse_nexpr() {
# [  expr ]
	static nexpr_set

	/nexpr_set := set([lex_AT, lex_BACKSLASH, lex_BANG, lex_BAR,
		lex_BREAK, lex_CARET, lex_CASE, lex_CONCAT, lex_CREATE,
		lex_DIFF, lex_DOT, lex_EQUIV, lex_EVERY,
		lex_FAIL, lex_IDENT, lex_IF, lex_INTER, lex_INTLIT,
		lex_LBRACE, lex_LBRACK, lex_LCONCAT, lex_LPAREN,
		lex_MINUS, lex_NEQUIV, lex_CONTINUE, lex_NMEQ, lex_NMNE,
		lex_NOT, lex_PLUS, lex_QMARK, lex_REALLIT, lex_REPEAT,
		lex_RETURN, lex_SEQ, lex_SLASH, lex_SNE, lex_STAR,
		lex_STRINGLIT, lex_SUSPEND, lex_UNION,
		lex_WHILE,
		lex_WITH,
		lex_MOD,
		lex_NIL,
		lex_SLASHSLASH,
		lex_ANDAND,
		lex_YIELD,
		lex_LAMBDA,
		lex_SELECT,
		lex_CATCH,
		lex_LOCAL,
		lex_STATIC,
		lex_LCOMP,
		lex_PROCEDURE ])
	if nexpr_set.member(parse_tok_rec) then {
		return parse_expr()
	}
	return nil
}

procedure parse_do_lambda() {
	local paramList
	local accumulate
	local loc
	local init
	local nexprList
	local e
	local coord
	local idcoord
	local endcoord
	local body
	static do_proc_set

	/do_proc_set := set([lex_AT, lex_BACKSLASH, lex_BANG, lex_BAR,
		lex_BREAK, lex_CARET, lex_CASE, lex_CONCAT, lex_CREATE,
		lex_DIFF, lex_DOT, lex_EQUIV, lex_EVERY,
		lex_FAIL, lex_IDENT, lex_IF, lex_INTER, lex_INTLIT,
		lex_LBRACE, lex_LBRACK, lex_LCONCAT, lex_LPAREN,
		lex_MINUS, lex_NEQUIV, lex_CONTINUE, lex_NMEQ, lex_NMNE,
		lex_NOT, lex_PLUS, lex_QMARK, lex_REALLIT, lex_REPEAT,
		lex_RETURN, lex_SEMICOL, lex_SEQ, lex_SLASH,
		lex_SNE, lex_STAR, lex_STRINGLIT, lex_SUSPEND,
		lex_UNION, lex_WHILE,
		lex_WITH,
		lex_SLASHSLASH,
		lex_MOD,
		lex_NIL,
		lex_ANDAND,
		lex_YIELD,
		lex_LAMBDA,
		lex_SELECT,
		lex_CATCH,
		lex_LOCAL,
		lex_STATIC,
		lex_LCOMP,
		lex_PROCEDURE ])
	#  lambdahead  expr
	coord := parse_tok_rec.coord
	parse_match_token(lex_LAMBDA)

	parse_match_token(lex_LPAREN)
	paramList := []
	if parse_tok_rec === lex_IDENT then {
		paramList := parse_idlist()
		if parse_tok_rec === lex_LBRACK then {
			parse_eat_token()
			parse_match_token(lex_RBRACK)
			accumulate := 1
		}
	}
	parse_match_token(lex_RPAREN)

#     parse_match_token(lex_SEMICOL)
	if parse_tok_rec === lex_SEMICOL then {
		# horrible hack to get rid of inserted semicolon when
		# opening brace is on the next line
		parse_match_token(lex_SEMICOL)
	}

	body := parse_expr()
	endcoord := parse_tok_rec.coord
	return a_ProcDecl(a_Ident(nil, nil, idcoord), paramList, accumulate,
		a_ProcCode(
			a_Compound([a_Suspend(body, nil, nil, coord)], coord), endcoord),
		coord, endcoord)
}

procedure parse_do_proc(noident) {
	local ident
	local paramList
	local accumulate
	local loc
	local init
	local nexprList
	local e
	local coord
	local idcoord
	local endcoord
	local body
	static do_proc_set

	/do_proc_set := set([lex_AT, lex_BACKSLASH, lex_BANG, lex_BAR,
		lex_BREAK, lex_CARET, lex_CASE, lex_CONCAT, lex_CREATE,
		lex_DIFF, lex_DOT, lex_EQUIV, lex_EVERY,
		lex_FAIL, lex_IDENT, lex_IF, lex_INTER, lex_INTLIT,
		lex_LBRACE, lex_LBRACK, lex_LCONCAT, lex_LPAREN,
		lex_MINUS, lex_NEQUIV, lex_CONTINUE, lex_NMEQ, lex_NMNE,
		lex_NOT, lex_PLUS, lex_QMARK, lex_REALLIT, lex_REPEAT,
		lex_RETURN, lex_SEMICOL, lex_SEQ, lex_SLASH,
		lex_SNE, lex_STAR, lex_STRINGLIT, lex_SUSPEND,
		lex_UNION, lex_WHILE,
		lex_WITH,
		lex_SLASHSLASH,
		lex_MOD,
		lex_NIL,
		lex_ANDAND,
		lex_YIELD,
		lex_LAMBDA,
		lex_SELECT,
		lex_CATCH,
		lex_LOCAL,
		lex_STATIC,
		lex_LCOMP,
		lex_PROCEDURE ])
	#  prochead  SEMICOL  locals  initial  procbody  END
	coord := parse_tok_rec.coord
	parse_match_token(lex_PROCEDURE)
	idcoord := parse_tok_rec.coord
	if /noident then {
		ident := a_Ident(parse_match_token(lex_IDENT), nil, idcoord)
		if parse_tok_rec === lex_DOT then {	# it's a method!
			parse_eat_token()
			ident.id ||:= "." || parse_match_token(lex_IDENT)
		}
	} else {
		ident := a_Ident(nil, nil, idcoord)
	}
	parse_match_token(lex_LPAREN)
	paramList := []
	if parse_tok_rec === lex_IDENT then {
		paramList := parse_idlist()
		if parse_tok_rec === lex_LBRACK then {
			parse_eat_token()
			parse_match_token(lex_RBRACK)
			accumulate := 1
		}
	}
	if !!contains(\ident.id, ".") then {	# if a method
		# add an initial "self" parameter; note NOT a reserved word
		paramList.push(a_Ident("self", nil, idcoord))
	}

	parse_match_token(lex_RPAREN)
#     parse_match_token(lex_SEMICOL)
	if parse_tok_rec === lex_SEMICOL then {
		# horrible hack to get rid of inserted semicolon when
		# opening brace is on the next line
		parse_match_token(lex_SEMICOL)
	}

	body := parse_braced()
#    nexprList := []
#    while do_proc_set.member(parse_tok_rec) do {
#        e := parse_nexpr()
#        put(nexprList, e)
#        parse_match_token(lex_SEMICOL)
#    }
	endcoord := parse_tok_rec.coord
#    parse_match_token(lex_END)
	return a_ProcDecl(ident, paramList, accumulate,
		# a_ProcCode(a_ProcBody(nexprList, endcoord), endcoord),
		a_ProcCode(body, endcoord), coord, endcoord)
}

procedure parse_program() {
# {  decl }
	static program_set
	local d

	/program_set := set([lex_GLOBAL, lex_PROCEDURE, lex_INITIAL, lex_RECORD])

	if  parse_tok_rec === lex_PACKAGE then {
		d := parse_do_package()
		suspend d
	}
	while parse_tok_rec ~=== lex_EOFX do {
		if \d then {
			parse_match_token(lex_SEMICOL)
		}
		if parse_tok_rec === lex_EOFX then {
			break
		}
		if program_set.member(parse_tok_rec) then {
			d := parse_decl()
			suspend d
		} else {
			d := parse_expr()
			# suspend d
		}
	}
}

procedure parse_do_record() {
#  RECORD  IDENT  LPAREN [  idlist ]  RPAREN
	local id
	local l
	local coord
	local idcoord
	local ex
	local expkg
	local excoord
	coord := parse_tok_rec.coord
	parse_match_token(lex_RECORD)
	idcoord := parse_tok_rec.coord
	id := a_Ident(parse_match_token(lex_IDENT), nil, idcoord)

	if parse_tok_rec === lex_EXTENDS then {
		parse_eat_token()
		excoord := parse_tok_rec.coord
		ex := a_Ident(parse_match_token(lex_IDENT), nil, excoord)
		if parse_tok_rec === lex_DOT then {
		parse_eat_token()
		expkg := ex
		excoord := parse_tok_rec.coord
		ex := a_Ident(parse_match_token(lex_IDENT), nil, excoord)
		}
	}

	parse_match_token(lex_LPAREN)
	l := []
	if parse_tok_rec === lex_IDENT then {
		l := parse_idlist()
	}
	parse_match_token(lex_RPAREN)
	return a_Record(id, ex, expkg, l, coord)
}

procedure parse_do_repeat() {
	local b
	local e
	local coord
	local id
	#  REPEAT  expr
	coord := parse_tok_rec.coord
	parse_match_token(lex_REPEAT)
	if parse_tok_rec === lex_COLON then {
		parse_eat_token()
		id := parse_match_token(lex_IDENT)
	}
	b := parse_expr()
	if parse_tok_rec === lex_UNTIL then {
		parse_match_token(lex_UNTIL)
		e := parse_expr()
	}
	return a_Repeat(b, e, id, coord)
}

procedure parse_do_catch() {
	local e
	local coord
	#  CATCH  expr
	coord := parse_tok_rec.coord
	parse_match_token(lex_CATCH)
	e := parse_expr()
	return a_Catch(e, coord)
}

procedure parse_do_return() {
	local e
	local coord
	local doparse_expr
	local id

	coord := parse_tok_rec.coord
	case parse_tok_rec of {
		lex_RETURN  :   # RETURN  nexpr
			{
			parse_eat_token()
			e := parse_nexpr()
			return a_Return(e, coord)
			}
		lex_SUSPEND :   # SUSPEND  expr [  DO  expr ]
			{
			parse_eat_token()
			if parse_tok_rec === lex_COLON then {
				parse_eat_token()
				id := parse_match_token(lex_IDENT)
			}
			e := parse_nexpr()
			doparse_expr := nil
			if parse_tok_rec === lex_DO then {
				parse_eat_token()
				doparse_expr := parse_expr()
			}
			return a_Suspend(e, doparse_expr, id, coord)
			}
		default :
			parse_error("Expecting lex_FAIL, lex_RETURN, or SUSPEND")
	}
}

#  WITH  % id [ := expr ] DO  expr
procedure parse_do_with() {
	local e
	local init
	local coord
	local id
	local current
	local root
	local tmp

	coord := parse_tok_rec.coord
	parse_match_token(lex_WITH)

	while parse_tok_rec === lex_MOD do {
		parse_match_token(lex_MOD)
		id := parse_match_token(lex_IDENT)
		init := nil
		if parse_tok_rec === lex_ASSIGN then {
			parse_eat_token()
			init := parse_expr()
		}
		tmp := a_With(id, init, nil, coord)
		/root:= tmp
		(\current).expr := tmp
		current := tmp
		if parse_tok_rec ~=== lex_COMMA then {
			break
		}
	parse_match_token(lex_COMMA)
	}

	parse_match_token(lex_DO)
	current.expr := parse_braced()
	return root
}

#  WHILE  expr [  DO  expr ]
procedure parse_do_while() {
	local e
	local dparse_expr
	local coord
	local id

	coord := parse_tok_rec.coord
	parse_match_token(lex_WHILE)
	if parse_tok_rec === lex_COLON then {
		parse_eat_token()
		id := parse_match_token(lex_IDENT)
	}
	e := parse_expr()
	dparse_expr := nil
	if parse_tok_rec === lex_DO then {
		parse_eat_token()
		dparse_expr := parse_expr()
	}
	return a_While(e, dparse_expr, id, coord)
}

procedure parse_match_token(which_token) {
	local saved
	saved := parse_tok_rec.str
	if parse_tok_rec === which_token then {
		parse_tok_rec := @parse_tok
		return saved
	} else {
		if which_token ~=== lex_IDENT then {
			parse_error("Expecting "|| which_token.str ||
				", but found " || parse_tok_rec.str)
		} else {
			parse_error("Expecting identifier")
		}
	}
}

procedure parse_eat_token() {
	local saved
	saved := parse_tok_rec.str
	parse_tok_rec := @parse_tok
	return saved
}


procedure parse_error(msg) {
	stop("At ", parse_tok_rec.coord, ": ", msg)
}

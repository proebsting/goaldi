#  parse.gd -- LL(1) parser that maps tokens to abstract syntax trees.


record parser (		# data associated with a particular input stream
	tok_stream,		# tokens coming in from lexer
	cur_tok,		# current token being inspected
)

record paired_lists(nameList, exprList)


#  overall control

procedure parse(lex) {
	local p := parser(lex, @lex)
	suspend p.parse_program()
}


#  peek_token(tag) -- succeed if the next token has the correct tag

procedure parser.peek_token(which_tag) {
	return self.cur_tok.tag === which_tag
}


#  match_token(tag) -- consume and return token, which must be of correct kind

procedure parser.match_token(which_tag) {
	local saved := self.cur_tok.str
	if self.cur_tok.tag === which_tag then {
		self.cur_tok := @self.tok_stream
		return saved
	} else {
		self.abort("Expecting "|| which_tag.str ||
				", but found " || self.cur_tok.str)
	}
}


#  eat_token() -- consume and return token

procedure parser.eat_token() {
	local saved
	saved := self.cur_tok.str
	self.cur_tok := @self.tok_stream
	return saved
}


#  abort(message) -- diagnose error and abort

procedure parser.abort(msg) {
	stop("At ", self.cur_tok.coord, ": ", msg)
}


# ----------------------------------------------------------------------------


procedure parser.parse_program() {
# {  decl }
	static program_set
	local d

	/program_set := set([lex_GLOBAL, lex_PROCEDURE, lex_INITIAL, lex_RECORD])

	if  self.peek_token(lex_PACKAGE) then {
		d := self.parse_do_package()
		suspend d
	}
	while not self.peek_token(lex_EOFX) do {
		if \d then {
			self.match_token(lex_SEMICOL)
		}
		if self.peek_token(lex_EOFX) then {
			break
		}
		if program_set.member(self.cur_tok.tag) then {
			d := self.parse_decl()
			suspend d
		} else {
			d := self.parse_expr()
			# suspend d
		}
	}
}

#  PACKAGE  name, coord
procedure parser.parse_do_package() {
	local coord
	local coord2
	local id

	coord := self.cur_tok.coord
	self.match_token(lex_PACKAGE)

	coord2 := self.cur_tok.coord
	id := a_Ident(self.match_token(lex_IDENT), nil, coord2)

	return a_Package(id, coord)
}

#  RECORD... | PROCEDURE... | GLOBAL... | INITIAL...
procedure parser.parse_decl() {
	case self.cur_tok.tag of {
		lex_RECORD    : return self.parse_do_record()
		lex_PROCEDURE : return self.parse_do_proc()
		lex_GLOBAL    : return self.parse_do_global()
		lex_INITIAL   : return self.parse_do_initial()
		default       : self.abort("Expecting declaration")
	}
}

procedure parser.parse_do_record() {
#  RECORD  IDENT  LPAREN [  idlist ]  RPAREN
	local id
	local l
	local coord
	local idcoord
	local ex
	local expkg
	local excoord
	coord := self.cur_tok.coord
	self.match_token(lex_RECORD)
	idcoord := self.cur_tok.coord
	id := a_Ident(self.match_token(lex_IDENT), nil, idcoord)

	if self.peek_token(lex_EXTENDS) then {
		self.eat_token()
		excoord := self.cur_tok.coord
		ex := a_Ident(self.match_token(lex_IDENT), nil, excoord)
		if self.peek_token(lex_DOT) then {
		self.eat_token()
		expkg := ex
		excoord := self.cur_tok.coord
		ex := a_Ident(self.match_token(lex_IDENT), nil, excoord)
		}
	}

	self.match_token(lex_LPAREN)
	l := []
	if self.peek_token(lex_IDENT) then {
		l := self.parse_idlist()
	}
	self.match_token(lex_RPAREN)
	return a_Record(id, ex, expkg, l, coord)
}

#  GLOBAL  idlist, coord
procedure parser.parse_do_global() {
	local coord
	local e
	local id

	coord := self.cur_tok.coord
	self.match_token(lex_GLOBAL)
	id := self.match_token(lex_IDENT)
	if self.peek_token(lex_ASSIGN) then {
		self.match_token(lex_ASSIGN)
		e := self.parse_expr()
		e := a_Binop(":=", a_Ident(id, nil, coord), e)
		e := a_ProcCode(e)
	}
	return a_Global(id, e, coord)
}

procedure parser.parse_do_initial() {
# INITIAL expr
	local coord
	local e

	coord := self.cur_tok.coord
	self.match_token(lex_INITIAL)
	e := self.parse_braced()
	return a_Initial(a_ProcCode(e), coord)
}

#  PROCEDURE ...
procedure parser.parse_do_proc(noident) {
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
	coord := self.cur_tok.coord
	self.match_token(lex_PROCEDURE)
	idcoord := self.cur_tok.coord
	if /noident then {
		ident := a_Ident(self.match_token(lex_IDENT), nil, idcoord)
		if self.peek_token(lex_DOT) then {	# it's a method!
			self.eat_token()
			ident.id ||:= "." || self.match_token(lex_IDENT)
		}
	} else {
		ident := a_Ident(nil, nil, idcoord)
	}
	self.match_token(lex_LPAREN)
	paramList := []
	if self.peek_token(lex_IDENT) then {
		paramList := self.parse_idlist()
		if self.peek_token(lex_LBRACK) then {
			self.eat_token()
			self.match_token(lex_RBRACK)
			accumulate := 1
		}
	}
	if !!contains(\ident.id, ".") then {	# if a method
		# add an initial "self" parameter; note NOT a reserved word
		paramList.push(a_Ident("self", nil, idcoord))
	}

	self.match_token(lex_RPAREN)
#     self.match_token(lex_SEMICOL)
	if self.peek_token(lex_SEMICOL) then {
		# horrible hack to get rid of inserted semicolon when
		# opening brace is on the next line
		self.match_token(lex_SEMICOL)
	}

	body := self.parse_braced()
#    nexprList := []
#    while do_proc_set.member(self.cur_tok.tag) do {
#        e := self.parse_nexpr()
#        put(nexprList, e)
#        self.match_token(lex_SEMICOL)
#    }
	endcoord := self.cur_tok.coord
#    self.match_token(lex_END)
	return a_ProcDecl(ident, paramList, accumulate,
		# a_ProcCode(a_ProcBody(nexprList, endcoord), endcoord),
		a_ProcCode(body, endcoord), coord, endcoord)
}

procedure parser.parse_braced() {
#  "{" compound "}"
	local e
	local coord

	coord := self.cur_tok.coord
	self.match_token(lex_LBRACE)
	e := self.parse_compound()
	self.match_token(lex_RBRACE)
	return a_Compound(e, coord)
}

procedure parser.parse_compound() {
#  nexpr {  SEMICOL  nexpr }
	local l
	local e

	l := [self.parse_nexpr()]
	while self.peek_token(lex_SEMICOL) do {
		self.eat_token()
		e := self.parse_nexpr()
		l.put(e)
	}
	return l
}

procedure parser.parse_nexpr() {
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
	if nexpr_set.member(self.cur_tok.tag) then {
		return self.parse_expr()
	}
	return nil
}

procedure parser.parse_expr() {
#  expr1x {  ANDAND  expr1x }
	local ret
	local L

	ret := self.parse_expr1x()
	if self.peek_token(lex_ANDAND) then {
		L := [ret]
		while self.peek_token(lex_ANDAND) do {
			self.eat_token()
			ret := self.parse_expr1x()
			L.put(ret)
		}
		return a_Parallel(L, ret.coord)
	} else {
		return ret
	}
}

procedure parser.parse_expr1x() {
#  expr1 {  AND  expr1 }
	local ret
	local op
	local right

	ret := self.parse_expr1()
	while self.peek_token(lex_AND) do {
		op := self.eat_token()
		right := self.parse_expr1()
		ret := a_Binop(op, ret, right, ret.coord)
	}
	return ret
}

procedure parser.parse_expr1() {
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
	ret := self.parse_expr2()
	while expr1_set.member(self.cur_tok.tag) do {
		coord := self.cur_tok.coord
		op := self.eat_token()
		right := self.parse_expr1()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parser.parse_expr2() {
#  expr3 {  TO  expr2 [  BY  expr3 ] }
	local e1
	local e2
	local e3
	local ret
	local coord

	e1 := self.parse_expr3()
	e2 := nil
	e3 := nil
	ret := e1
	while self.peek_token(lex_TO) do {
	coord := self.cur_tok.coord
		self.eat_token()
		e2 := self.parse_expr3()
		if self.peek_token(lex_BY) then {
			self.match_token(lex_BY)
			e3 := self.parse_expr3()
		}
		ret := a_ToBy(e1, e2, e3, coord)
		e1 := ret
	}
	return ret
}

procedure parser.parse_expr3() {
#  expr3a {  TILDEBAR  expr3 }
	local ret
	local a

	ret := self.parse_expr3a()
	while self.peek_token(lex_TILDEBAR) do {
		/a := a_ExcAlt([ret], self.cur_tok.coord)
		self.eat_token()
		a.eList.put(self.parse_expr3a())
	}
	ret := \a
	return ret
}

procedure parser.parse_expr3a() {
#  expr4 {  BAR  expr3 }
	local ret
	local a

	ret := self.parse_expr4()
	while self.peek_token(lex_BAR) do {
		/a := a_Alt([ret], self.cur_tok.coord)
		self.eat_token()
		a.eList.put(self.parse_expr4())
	}
	ret := \a
	return ret
}

procedure parser.parse_expr4() {
	local ret
	local op
	local right
	local coord
	static expr4_set
	/expr4_set := set([lex_EQUIV, lex_NEQUIV, lex_NMEQ, lex_NMGE,
		lex_NMGT, lex_NMLE, lex_NMLT, lex_NMNE, lex_SEQ,
		lex_SGE, lex_SGT, lex_SLE, lex_SLT, lex_SNE])
	#  expr5 {  expr4op  expr4 }
	ret := self.parse_expr5()

	while expr4_set.member(self.cur_tok.tag) do {
		coord := self.cur_tok.coord
		op := self.eat_token()
		right := self.parse_expr5()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parser.parse_expr5() {
#  expr6 {  expr5op  expr5 }
	local ret
	local right
	local op
	local coord

	ret := self.parse_expr6()
	while self.peek_token(lex_CONCAT | lex_LCONCAT) do {
		coord := self.cur_tok.coord
		op := self.eat_token()
		right := self.parse_expr6()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}


procedure parser.parse_expr6() {
#  expr7 {  expr6op  expr6 }
	local ret
	local op
	local right
	local coord
	static expr6_set

	/expr6_set := set([lex_DIFF, lex_MINUS, lex_PLUS, lex_UNION])

	ret := self.parse_expr7()
	while expr6_set.member(self.cur_tok.tag) do {
		coord := self.cur_tok.coord
		op := self.eat_token()
		right := self.parse_expr7()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parser.parse_expr7() {
#  expr8 {  expr7op  expr7 }
	local ret
	local op
	local right
	local coord
	static expr7_set

	/expr7_set := set([lex_INTER, lex_MOD, lex_SLASH, lex_STAR, lex_SLASHSLASH])

	ret := self.parse_expr8()

	while expr7_set.member(self.cur_tok.tag) do {
		coord := self.cur_tok.coord
		op := self.eat_token()
		right := self.parse_expr8()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parser.parse_expr8() {
#  expr9 {  CARET  expr8 }  (Right Associative)
	local ret
	local op
	local right
	local coord

	ret := self.parse_expr9()
	while self.peek_token(lex_CARET) do {
		coord := self.cur_tok.coord
		op := self.eat_token()
		right := self.parse_expr8()
		ret := a_Binop(op, ret, right, coord)
	}
	return ret
}

procedure parser.parse_expr9() {
#  expr10 {  expr9op  expr9 }
	local ret
	local op
	local right
	local coord

	ret := self.parse_expr10()
	while self.peek_token( lex_BACKSLASH | lex_BANG ) do {
		coord := self.cur_tok.coord
		op := self.eat_token()
		right := self.parse_expr10()
		if op == "\\" then {
			ret := a_Limitation(ret, right, coord)
		} else {
			ret := a_Binop(op, ret, right, coord)
		}
	}
	return ret
}

procedure parser.parse_expr10() {
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
	/expr10_set2 := set([lex_AT, lex_NOT, lex_BAR, lex_BANG,
		lex_PLUS, lex_STAR, lex_SLASH, 
		lex_MINUS, 
		lex_QMARK,
		lex_BACKSLASH])
	/expr10_set3 := set([lex_CONCAT, lex_LCONCAT, lex_UNION, lex_INTER,
		lex_SLASHSLASH,
		lex_DIFF])

	if expr10_set1.member(self.cur_tok.tag) then {
		return self.parse_expr11a()
	} else if expr10_set2.member(self.cur_tok.tag) then {
		coord := self.cur_tok.coord
		op := self.eat_token()
		operand := self.parse_expr10()
		case (op) of {
			"|":        return a_RepAlt(operand, coord)
			"not":      return a_Not(operand, coord)
			default:    return a_Unop(op, operand, coord)
		}
	} else if expr10_set3.member(self.cur_tok.tag) then {
		tmp_tok := self.cur_tok
		coord := self.cur_tok.coord
		op := self.eat_token()
		operand := self.parse_expr10()
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
		self.abort("\""||self.cur_tok.str||"\": Expecting expression")
	}
}

procedure parser.parse_expr11a() {
	local left
	local right
	local op
	local oprec
	static expr11_set

	/expr11_set := set([lex_DOT, lex_LBRACE, lex_LBRACK, lex_LPAREN, lex_LCOMP])
	#  expr11 {  expr11suffix }
	left := self.parse_expr11()
	while expr11_set.member(self.cur_tok.tag) do {
		left := self.parse_expr11suffix(left)
	}
	return left
}

procedure parser.parse_expr11() {
	local e
	local id
	local coord
	local ns
	case self.cur_tok.tag of {
		lex_INTLIT  |
		lex_REALLIT |
		lex_STRINGLIT   :   # literal
			return self.parse_literal()
		lex_CATCH   :
			return self.parse_do_catch()
		lex_NIL    : {
			coord := self.cur_tok.coord
			self.eat_token()
			return a_Nil(coord)
			}
		lex_FAIL    : {
			coord := self.cur_tok.coord
			self.eat_token()
			return a_Fail(coord)
			}
		lex_RETURN  |
		lex_SUSPEND :   # return
			return self.parse_do_return()
		lex_IF  :   # if
			return self.parse_do_if()
		lex_CASE    :   # case
			return self.parse_do_case()
		lex_SELECT  :   # select
			return self.parse_do_select()
		lex_WHILE   :   # while
			return self.parse_do_while()
		lex_EVERY   :   # every
			return self.parse_do_every()
		lex_REPEAT  :   # repeat
			return self.parse_do_repeat()
		lex_CREATE  :   # CREATE  expr
			{
			coord := self.cur_tok.coord
			self.eat_token()
			e := self.parse_expr()
			return a_Create(e, coord)
			}
		lex_IDENT   :   # IDENT
			{
			coord := self.cur_tok.coord
			id := self.eat_token()
			ns := nil
			if self.peek_token(lex_COLONCOLON) then {
				self.eat_token()
				ns := id
				id := self.match_token(lex_IDENT)
			}
			return a_Ident(id, ns, coord)
			}
		lex_CARET   |
		lex_LOCAL   :
			{
			coord := self.cur_tok.coord
			self.eat_token()
			id := self.match_token(lex_IDENT)
			return a_Local(id, coord)
			}
		lex_STATIC   :   # STATIC
			{
			coord := self.cur_tok.coord
			self.eat_token()
			id := self.match_token(lex_IDENT)
			return a_Static(id, coord)
			}
		lex_CONTINUE    :   # CONTINUE
			{
			coord := self.cur_tok.coord
			self.eat_token()
			if self.peek_token(lex_COLON) then {
				self.eat_token()
				id := self.match_token(lex_IDENT)
			}
			return a_Continue(id, coord)
			}
		lex_YIELD   :   # YIELD  [ : IDENT ] nexpr
			{
			coord := self.cur_tok.coord
			self.eat_token()
			if self.peek_token(lex_COLON) then {
				self.eat_token()
				id := self.match_token(lex_IDENT)
			}
			e := self.parse_nexpr()
			return a_Yield(e, id, coord)
			}
		lex_BREAK   :   # BREAK  nexpr
			{
			coord := self.cur_tok.coord
			self.eat_token()
			if self.peek_token(lex_COLON) then {
				self.eat_token()
				id := self.match_token(lex_IDENT)
			}
			return a_Break(id, coord)
			}
		lex_LPAREN  :   # LPAREN  exprlist  RPAREN
			{
			coord := self.cur_tok.coord
			self.eat_token()
			e := self.parse_exprlist()
			self.match_token(lex_RPAREN)
			return a_Mutual(e, coord)
			}
		lex_LBRACE  :   # LBRACE  compound  RBRACE
			{
			return self.parse_braced()
			}
		lex_LBRACK  :   # LBRACK  exprlist  RBRACK
			{
			coord := self.cur_tok.coord
			self.eat_token()
			e := self.parse_exprlist()
			self.match_token(lex_RBRACK)
			return a_ListConstructor(e, coord)
			}
		lex_LCOMP  :   # LCOMP  expr  RCOMP
			{
			coord := self.cur_tok.coord
			self.eat_token()
			e := self.parse_expr()
			self.match_token(lex_RCOMP)
			return a_ListComprehension(e, coord)
			}
		lex_MOD :   # MOD IDENT
			{
			coord := self.cur_tok.coord
			self.eat_token()
			return a_Key(self.match_token(lex_IDENT), coord)
			}
		lex_PROCEDURE : {
				return self.parse_do_proc("noident")
			}
		lex_LAMBDA : {
				return self.parse_do_lambda()
			}
		lex_WITH : {
				return self.parse_do_with()
			}
		default : {
			self.abort("Expecting expression")
		}
	}
}

procedure parser.parse_expr11suffix(lhs) {
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

	case self.cur_tok.tag of {
		lex_LBRACE  :   # LBRACE [  { e : e [ COMMA ] ]  RBRACE
			{
			lefts := []
			rights := []
			coord := self.cur_tok.coord
			self.eat_token()
			if self.cur_tok.tag ~=== lex_RBRACE then {
				while self.cur_tok.tag ~=== lex_RBRACE do {
					left := self.parse_expr()
					self.match_token(lex_COLON)
					right := self.parse_expr()
					lefts.put(left)
					rights.put(right)
					if self.peek_token(lex_COMMA) then {
						self.eat_token()
					} else {
						break
					}
				}
			}
			self.match_token(lex_RBRACE)
			return a_Paired(lhs, lefts, rights, coord)
			}
		lex_LPAREN  :   # LPAREN  exprlist  RPAREN
			{
			coord := self.cur_tok.coord
			self.eat_token()
			l := self.parse_named_exprlist()
			self.match_token(lex_RPAREN)
				if type(l) === paired_lists then {
					return a_Call(lhs, a_Arglist(l.exprList, l.nameList), coord)
				} else {
					return a_Call(lhs, a_Arglist(l, nil), coord)
				}
			}
		lex_DOT :   # DOT  IDENT
			{
			coord := self.cur_tok.coord
			x := self.eat_token()
			idcoord := self.cur_tok.coord
			id := self.match_token(lex_IDENT)
			return a_Field(lhs, a_Ident(id, nil, idcoord), coord)
			}
		lex_LBRACK  :   # LBRACK  expr [  sectop  expr ]  RBRACK
			{
			coord := self.cur_tok.coord
			self.eat_token()
			left := self.parse_nexpr()
			if expr11suffix_set2.member(self.cur_tok.tag) then {
				coord := self.cur_tok.coord
				op := "[" || self.eat_token() || "]"
				right := self.parse_expr()
				lhs := a_Sectionop(op, lhs, left, right, coord)
			} else {
				lhs := a_Binop("[]", lhs, left, coord)
				while self.peek_token(lex_COMMA) do {
					self.match_token(lex_COMMA)
					left := self.parse_nexpr()
					lhs := a_Binop("[]", lhs, left, coord)
				}
			}
			self.match_token(lex_RBRACK)
			return lhs
			}
		default : {
			self.abort("Malformed argument list")
		}
	}
}

procedure parser.parse_named_exprlist0(id) {
	local e
	local N
	local E
	self.match_token(lex_COLON)
	N := []
	E := []
	e := self.parse_expr()
	N.put(id.id)
	E.put(e)

	while self.peek_token(lex_COMMA) do {
		self.eat_token()
		if not self.peek_token(lex_IDENT) then {
			break
		}
		id := self.match_token(lex_IDENT)
		self.match_token(lex_COLON)
		e := self.parse_nexpr()
		N.put(id)
		E.put(e)
	}
	return paired_lists(N, E)
}

procedure parser.parse_named_exprlistX() {
	local l
	local e

	e := self.parse_nexpr()
	if \e & type(e) === a_Ident & self.peek_token(lex_COLON) then {
		return self.parse_named_exprlist0(e)
	} else if \e | (self.peek_token(lex_COMMA)) then {
		l := [ e ]
	} else {
		l := []
	}
	while self.peek_token(lex_COMMA) do {
		self.eat_token()
		e := self.parse_nexpr()
		l.put(e)
	}
	if /l[-1] then {
		l := l[1:-1]
	}
	return l
}

procedure parser.parse_named_exprlist() {
	local L
	local N
	local e

	L := []
	N := []
	repeat {
		e := self.parse_nexpr()
		if self.peek_token(lex_COMMA) then {
			L.put(e)
		} else if \e then {
			if type(e) === a_Ident & self.peek_token(lex_COLON) then {
				N.put(e.id)
				self.match_token(lex_COLON)
				e := self.parse_expr()
				L.put(e)
			} else {
				L.put(e)
			}
		} else {
		}
		if not self.peek_token(lex_COMMA) then {
			break
		}
		self.match_token(lex_COMMA)
	}

	if *N > 0 then {
		return paired_lists(N, L)
	} else {
		return L
	}
}

procedure parser.parse_exprlist() {
	local l
	local e

	e := self.parse_nexpr()
	if \e | (self.peek_token(lex_COMMA)) then {
		l := [ e ]
	} else {
		l := []
	}
	while self.peek_token(lex_COMMA) do {
		self.eat_token()
		e := self.parse_nexpr()
		l.put(e)
	}
	if /l[-1] then {
		l := l[1:-1]
	}
	return l
}


procedure parser.parse_idlist() {
#  IDENT {  COMMA  IDENT }
	local l
	local id
	local coord

	coord := self.cur_tok.coord

	l := [a_Ident(self.match_token(lex_IDENT), nil, coord)]
	while self.peek_token(lex_COMMA) do {
		self.eat_token()
		coord := self.cur_tok.coord
		if not self.peek_token(lex_IDENT) then {
			break
		}
		id := a_Ident(self.match_token(lex_IDENT), nil, coord)
		l.put(id)
	}
	return l
}

procedure parser.parse_literal() {
	local coord
	coord := self.cur_tok.coord
	case self.cur_tok.tag of {
		lex_INTLIT  :   # INTLIT
			return a_Intlit(integer(self.eat_token()), coord)
		lex_REALLIT :   # REALLIT
			return a_Reallit(number(self.eat_token()), coord)
		lex_STRINGLIT   :   # STRINGLIT
			return a_Stringlit(self.eat_token(), coord)
		default :
			self.abort("Expecting literal")
	}
}

procedure parser.parse_do_lambda() {
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
	coord := self.cur_tok.coord
	self.match_token(lex_LAMBDA)

	self.match_token(lex_LPAREN)
	paramList := []
	if self.peek_token(lex_IDENT) then {
		paramList := self.parse_idlist()
		if self.peek_token(lex_LBRACK) then {
			self.eat_token()
			self.match_token(lex_RBRACK)
			accumulate := 1
		}
	}
	self.match_token(lex_RPAREN)

#     self.match_token(lex_SEMICOL)
	if self.peek_token(lex_SEMICOL) then {
		# horrible hack to get rid of inserted semicolon when
		# opening brace is on the next line
		self.eat_token()
	}

	body := self.parse_expr()
	endcoord := self.cur_tok.coord
	return a_ProcDecl(a_Ident(nil, nil, idcoord), paramList, accumulate,
		a_ProcCode(
			a_Compound([a_Suspend(body, nil, nil, coord)], coord), endcoord),
		coord, endcoord)
}

procedure parser.parse_do_if() {
#  IF  expr  THEN  expr [  ELSE  expr ]
	local ex
	local then_expr
	local else_expr
	local coord

	coord := self.cur_tok.coord
	self.match_token(lex_IF)
	ex := self.parse_expr()
	self.match_token(lex_THEN)
	then_expr := self.parse_expr()
	else_expr := nil
	if self.peek_token(lex_ELSE) then {
		self.eat_token()
		else_expr := self.parse_expr()
	}
	return a_If(ex, then_expr, else_expr, coord)
}

procedure parser.parse_do_repeat() {
	local b
	local e
	local coord
	local id
	#  REPEAT  expr
	coord := self.cur_tok.coord
	self.match_token(lex_REPEAT)
	if self.peek_token(lex_COLON) then {
		self.eat_token()
		id := self.match_token(lex_IDENT)
	}
	b := self.parse_expr()
	if self.peek_token(lex_UNTIL) then {
		self.match_token(lex_UNTIL)
		e := self.parse_expr()
	}
	return a_Repeat(b, e, id, coord)
}

#  WHILE  expr [  DO  expr ]
procedure parser.parse_do_while() {
	local e
	local do_expr
	local coord
	local id

	coord := self.cur_tok.coord
	self.match_token(lex_WHILE)
	if self.peek_token(lex_COLON) then {
		self.eat_token()
		id := self.match_token(lex_IDENT)
	}
	e := self.parse_expr()
	do_expr := nil
	if self.peek_token(lex_DO) then {
		self.eat_token()
		do_expr := self.parse_expr()
	}
	return a_While(e, do_expr, id, coord)
}

procedure parser.parse_do_every() {
#  EVERY  expr [  DO  expr ]
	local e
	local body
	local coord
	local id
	coord := self.cur_tok.coord
	self.match_token(lex_EVERY)
	if self.peek_token(lex_COLON) then {
		self.eat_token()
		id := self.match_token(lex_IDENT)
	}
	e := self.parse_expr()
	if self.peek_token(lex_DO) then {
		self.match_token(lex_DO)
		body := self.parse_expr()
	} else {
		body := nil
	}
	return a_Every(e, body, id, coord)
}

procedure parser.parse_do_catch() {
	local e
	local coord
	#  CATCH  expr
	coord := self.cur_tok.coord
	self.match_token(lex_CATCH)
	e := self.parse_expr()
	return a_Catch(e, coord)
}

procedure parser.parse_do_return() {
	local e
	local coord
	local do_expr
	local id

	coord := self.cur_tok.coord
	case self.cur_tok.tag of {
		lex_RETURN  :   # RETURN  nexpr
			{
			self.eat_token()
			e := self.parse_nexpr()
			return a_Return(e, coord)
			}
		lex_SUSPEND :   # SUSPEND  expr [  DO  expr ]
			{
			self.eat_token()
			if self.peek_token(lex_COLON) then {
				self.eat_token()
				id := self.match_token(lex_IDENT)
			}
			e := self.parse_nexpr()
			do_expr := nil
			if self.peek_token(lex_DO) then {
				self.eat_token()
				do_expr := self.parse_expr()
			}
			return a_Suspend(e, do_expr, id, coord)
			}
		default :
			self.abort("Expecting RETURN or SUSPEND")
	}
}

#  WITH  % id [ := expr ] DO  expr
procedure parser.parse_do_with() {
	local e
	local init
	local coord
	local id
	local current
	local root
	local tmp

	coord := self.cur_tok.coord
	self.match_token(lex_WITH)

	while self.peek_token(lex_MOD) do {
		self.match_token(lex_MOD)
		id := self.match_token(lex_IDENT)
		init := nil
		if self.peek_token(lex_ASSIGN) then {
			self.eat_token()
			init := self.parse_expr()
		}
		tmp := a_With(id, init, nil, coord)
		/root:= tmp
		(\current).expr := tmp
		current := tmp
		if not self.peek_token(lex_COMMA) then {
			break
		}
		self.match_token(lex_COMMA)
	}

	self.match_token(lex_DO)
	current.expr := self.parse_braced()
	return root
}


#  CASE  expr  OF  LBRACE  cclause {  SEMICOL  cclause }  RBRACE
procedure parser.parse_do_case() {
	local e
	local body
	local element
	local dflt
	local coord

	coord := self.cur_tok.coord
	self.eat_token()
	e := self.parse_expr()
	self.match_token(lex_OF)
	self.match_token(lex_LBRACE)
	body := []
	element := self.parse_cclause()
	if element.expr === lex_DEFAULT then {
		dflt := element.body
	} else {
		body.put(element)
	}
	while self.peek_token(lex_SEMICOL) do {
		self.match_token(lex_SEMICOL)
		if self.peek_token(lex_RBRACE) then {
			break
		}
		element := self.parse_cclause()
		if element.expr === lex_DEFAULT then {
			(/dflt := element.body) | self.abort("Multiple default clauses")
		} else {
			body.put(element)
		}
	}
	self.match_token(lex_RBRACE)
	return a_Case(e, body, dflt, coord)
}

procedure parser.parse_cclause() {
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
	if self.peek_token(lex_DEFAULT) then {
		e := lex_DEFAULT
		self.eat_token()
		coord := self.cur_tok.coord
		self.match_token(lex_COLON)
		body := self.parse_expr()
	} else if cclause_set.member(self.cur_tok.tag) then {
		e := self.parse_expr()
		coord := self.cur_tok.coord
		self.match_token(lex_COLON)
		body := self.parse_expr()
	} else {
		self.abort("\""||self.cur_tok.str||"\": Invalid case clause")
	}
	return a_Cclause(e, body, coord)
}

procedure parser.parse_do_select() {
#  SELECT  LBRACE  selcase {  SEMICOL  selcase }  RBRACE
	local caseList
	local element
	local dflt
	local coord

	coord := self.cur_tok.coord
	self.eat_token()
	self.match_token(lex_LBRACE)
	caseList := []
	repeat {
		if self.peek_token(lex_RBRACE) then {
			break
		}
		element := self.parse_selcase()
		if element.kind === "default" then {
			(/dflt := element) | self.abort("More than one default clause")
		} else {
			caseList.put(element)
		}
		if self.peek_token(lex_SEMICOL) then {
			self.eat_token()
		} else {
			break
		}
	}
	self.match_token(lex_RBRACE)
	return a_Select(caseList, dflt, coord)
}

procedure parser.parse_selcase() {
# select-case:  select-condition : body
	local sc

	sc := a_SelectCase()
	self.parse_selectby(sc) | self.abort("Malformed select condition")
	sc.coord := self.cur_tok.coord
	self.match_token(lex_COLON)
	sc.body := self.parse_expr()
	return sc
}

procedure parser.parse_selectby(sc) {
# select-condition:  default  |  expr := @expr  |  expr @: expr
# sets sc.kind,left,right; fails on parse error
	local e
	local r

	if self.peek_token(lex_DEFAULT) then {
		self.eat_token()
		sc.kind := "default"
		return
	}
	e := self.parse_expr()
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

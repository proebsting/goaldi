#  lex.gd -- Goaldi tokenizer


#  a token record -- one shared/reused instance for each distinct token type

record lex_tkrec (
	str,		# actual or canonicalized form of the source token
	flags,		# beginning and/or ending flags
	coord,		# coordinate in source code (#%#% TO BE DONE)
)


#  data structures for tokenizing

global lex_lnum					# current line number
global lex_kwtab := table()		# maps keyword strings to token records
global lex_optab := table()		# maps operator strings to token records
global lex_flags := table()		# maps token records flag strings


#  report an error (and do nothing else)	#%#%#%

procedure lex_error(problem, input) {
	%stderr.write("lex error: ", problem, ": ", image(input))
}


#  generate a sequence of tokens, with coordinates, from a stream of input lines

procedure lex(src, fname) {
	fname := fname || ":"
	lex_lnum := 0
	every ^tk := lex_gentok(src) do {
		tk.coord := fname || lex_lnum
		suspend tk
	}
}


#  generate a sequence of tokens from a stream of lines

procedure lex_gentok(src) {
	while ^line := @src do {
		lex_lnum +:= 1
		^tk := nil
		while *line > 0 do {
			if ^s := match(line, lex_ws_rx) then {
				# whitespace: ignore
				line := line[1+*s:0]
			} else if s := match(line, lex_id_rx) then {
				# identifier form: check for possible keyword
				line := line[1+*s:0]
				if ^t := \lex_kwtab[s] then {
					suspend tk := t
				} else {
					lex_IDENT.str := s
					suspend tk := lex_IDENT
				}
			} else if s := match(line, lex_n1_rx | lex_n2_rx | lex_n3_rx) then {
				# number: test must precede operators to match ".123"
				line := line[1+*s:0]
				if ^n := number(s) then {
					lex_REALLIT.str := image(n)	# put in canonical form
					suspend tk := lex_REALLIT
				} else {
					lex_error("malformed number", s)
				}
			} else if s := match(line, lex_op_rx) then {
				# operator
				line := line[1+*s:0]
				suspend tk := \lex_optab[s] ~| throw("lost operator", s)
			} else if s := match(line, lex_s1_rx | lex_r1_rx) then {
				# simple string literal
				line := line[1+*s:0]
				lex_STRINGLIT.str := stringval(s)
				suspend tk := lex_STRINGLIT
			} else if s := match(line, lex_s2_rx) then {
				# unterminated string literal: error
				line := ""
				lex_error("unterminated string", s)
			} else if s := match(line, lex_r2_rx) then {
				# unterminated raw literal: may span lines, so keep reading
				repeat {
					s ||:= "\n"
					if line := @src then {
						lex_lnum +:= 1
						if ^t := match(line, lex_r3_rx) then {
							# found terminator
							line := line[1+*t:0]
							s ||:= t
							lex_STRINGLIT.str := stringval(s)
							suspend tk := lex_STRINGLIT
							break
						} else {
							s ||:= line
						}
					} else {
						s := s[1+:40] || "..."	# truncate for sane message
						lex_error("unterminated raw literal", s)
						line := ""
						break
					}
				}
			} else {
				# unrecognized
				s := line[1]
				line := line[2:0]
				lex_error("unrecognized token", s)
			}
		}
		if (\lex_flags[tk])[-1] == "e" then {
			suspend lex_SEMICOL				# semicolon insertion
		}
	}
	lex_lnum +:= 1
	suspend lex_EOFX
}


#  match(line,rx) -- return matching string if line is matched by rx, else fail
procedure match(line, rx) {
	return "" ~== rx.FindString(line)
}


#  stringval(s) -- put a string in canonical form, checking escapes
procedure stringval(s) {
	if not (s := quote(unquote(s))) then {
		lex_error("invalid string literal", s)
	}
	return s
}


#  These globals define regular expressions for tokenizing

global lex_ws_rx := regex(`^([ \t]+|\#.*)`)					# whitespace
global lex_id_rx := regex(`^[a-zA-Z_][a-zA-Z_0-9]*`)		# identifier / kwd
global lex_n1_rx := regex(`^[0-9]+r[0-9a-zA-Z]+`)			# radix prefix int
global lex_n2_rx := regex(`^[0-9][0-9a-fA-F]*[box]`)		# radix suffix int
global lex_n3_rx := regex(`^(\.[0-9]+|[0-9]+\.?[0-9]*)([eE][+-]?[0-9]+)?`) # dec
global lex_s1_rx := regex(`^"([^"\\]|\\.)*"`)				# "string"
global lex_s2_rx := regex(`^"([^"\\]|\\.)*$`)				# "unterminated...
global lex_r1_rx := regex("^`[^`]*`")						# `simple raw str`
global lex_r2_rx := regex("^`[^`]*$")						# `unterminated...
global lex_r3_rx := regex("^[^`]*`")						# ....`
global lex_op_rx			# (built by initial{} below)	# operator

#  This initial{} code builds the regular expression matching all operators.
#  It takes the form:  `^(op|op|...|op)`  (with all chars of each op escaped)
initial {				# builds lex_op_rx
	^expr := `^(`
	every ^op := (!lex_optab).key do {
		every ^c := !op do {
			expr ||:= `\` || c
		}
		expr ||:= `|`
	}
	expr := expr[1:-1] || `)`
	# this MUST use regexp (regex.CompilePOSIX) and not regex (regex.Compile)
	# because we depend on the "leftmost-longest" (greedy) matching rule
	lex_op_rx := regexp(expr)
}


#  These globals provide named handles for all the distinct token types.

global lex_IDENT         := lex_lit("identifier",      "be")
global lex_INTLIT        := lex_lit("integer-literal", "be")	# never produced
global lex_REALLIT       := lex_lit("real-literal",    "be")
global lex_STRINGLIT     := lex_lit("string-literal",  "be")
global lex_EOFX          := lex_lit("end-of-file",     "")

global lex_BREAK         := lex_kwd("break",     "be")
global lex_BY            := lex_kwd("by",        "")
global lex_CASE          := lex_kwd("case",      "b")
global lex_CATCH         := lex_kwd("catch",     "b")
global lex_CONTINUE      := lex_kwd("continue",  "be")
global lex_CREATE        := lex_kwd("create",    "b")
global lex_DEFAULT       := lex_kwd("default",   "b")
global lex_DO            := lex_kwd("do",        "")
global lex_ELSE          := lex_kwd("else",      "")
global lex_EVERY         := lex_kwd("every",     "b")
global lex_EXTENDS       := lex_kwd("extends",   "")
global lex_FAIL          := lex_kwd("fail",      "be")
global lex_GLOBAL        := lex_kwd("global",    "b")
global lex_IF            := lex_kwd("if",        "b")
global lex_INITIAL       := lex_kwd("initial",   "b")
global lex_LAMBDA        := lex_kwd("lambda",    "b")
global lex_LOCAL         := lex_kwd("local",     "b")
global lex_NOT           := lex_kwd("not",       "b")
global lex_NIL           := lex_kwd("nil",       "be")
global lex_OF            := lex_kwd("of",        "")
global lex_PACKAGE       := lex_kwd("package",   "b")
global lex_PROCEDURE     := lex_kwd("procedure", "b")
global lex_RECORD        := lex_kwd("record",    "b")
global lex_REPEAT        := lex_kwd("repeat",    "b")
global lex_RETURN        := lex_kwd("return",    "be")
global lex_SELECT        := lex_kwd("select",    "b")
global lex_STATIC        := lex_kwd("static",    "b")
global lex_SUSPEND       := lex_kwd("suspend",   "be")
global lex_THEN          := lex_kwd("then",      "")
global lex_TO            := lex_kwd("to",        "")
global lex_UNTIL         := lex_kwd("until",     "b")
global lex_WHILE         := lex_kwd("while",     "b")
global lex_WITH          := lex_kwd("with",      "b")
global lex_YIELD         := lex_kwd("yield",     "be")

global lex_ASSIGN        := lex_opr(":=",     "")
global lex_AT            := lex_opr("@",      "b")
global lex_ATCOLON       := lex_opr("@:",     "")
global lex_AUGAND        := lex_opr("&:=",    "")
global lex_AUGNMEQ       := lex_opr("=:=",    "")
global lex_AUGEQUIV      := lex_opr("===:=",  "")
global lex_AUGNMGE       := lex_opr(">=:=",   "")
global lex_AUGNMGT       := lex_opr(">:=",    "")
global lex_AUGNMLE       := lex_opr("<=:=",   "")
global lex_AUGNMLT       := lex_opr("<:=",    "")
global lex_AUGNMNE       := lex_opr("~=:=",   "")
global lex_AUGNEQUIV     := lex_opr("~===:=", "")
global lex_AUGSEQ        := lex_opr("==:=",   "")
global lex_AUGSGE        := lex_opr(">>=:=",  "")
global lex_AUGSGT        := lex_opr(">>:=",   "")
global lex_AUGSLE        := lex_opr("<<=:=",  "")
global lex_AUGSLT        := lex_opr("<<:=",   "")
global lex_AUGSNE        := lex_opr("~==:=",  "")
global lex_BACKSLASH     := lex_opr("\\",     "b")
global lex_BANG          := lex_opr("!",      "b")
global lex_BAR           := lex_opr("|",      "b")
global lex_TILDEBAR      := lex_opr("~|",     "")
global lex_CARET         := lex_opr("^",      "b")
global lex_AUGCARET      := lex_opr("^:=",    "b")
global lex_COLON         := lex_opr(":",      "")
global lex_COLONCOLON    := lex_opr("::",     "")
global lex_COMMA         := lex_opr(",",      "")
global lex_CONCAT        := lex_opr("||",     "b")
global lex_AUGCONCAT     := lex_opr("||:=",   "")
global lex_AND           := lex_opr("&",      "")
global lex_ANDAND        := lex_opr("&&",     "b")
global lex_DOT           := lex_opr(".",      "b")
global lex_DIFF          := lex_opr("--",     "b")
global lex_AUGDIFF       := lex_opr("--:=",   "")
global lex_EQUIV         := lex_opr("===",    "b")
global lex_INTER         := lex_opr("**",     "b")
global lex_AUGINTER      := lex_opr("**:=",   "")
global lex_LBRACE        := lex_opr("{",      "b")
global lex_LBRACK        := lex_opr("[",      "b")
global lex_LCOMP         := lex_opr("[:",     "b")
global lex_LCONCAT       := lex_opr("|||",    "b")
global lex_AUGLCONCAT    := lex_opr("|||:=",   "")
global lex_SEQ           := lex_opr("==",     "b")
global lex_SGE           := lex_opr(">>=",    "")
global lex_SGT           := lex_opr(">>",     "")
global lex_SLE           := lex_opr("<<=",    "")
global lex_SLT           := lex_opr("<<",     "")
global lex_SNE           := lex_opr("~==",    "b")
global lex_LPAREN        := lex_opr("(",      "b")
global lex_MCOLON        := lex_opr("-:",     "")
global lex_MINUS         := lex_opr("-",      "b")
global lex_AUGMINUS      := lex_opr("-:=",    "")
global lex_MOD           := lex_opr("%",      "b")
global lex_AUGMOD        := lex_opr("%:=",    "")
global lex_NEQUIV        := lex_opr("~===",   "b")
global lex_NMEQ          := lex_opr("=",      "b")
global lex_NMGE          := lex_opr(">=",     "")
global lex_NMGT          := lex_opr(">",      "")
global lex_NMLE          := lex_opr("<=",     "")
global lex_NMLT          := lex_opr("<",      "")
global lex_NMNE          := lex_opr("~=",     "b")
global lex_PCOLON        := lex_opr("+:",     "")
global lex_PLUS          := lex_opr("+",      "b")
global lex_AUGPLUS       := lex_opr("+:=",    "")
global lex_QMARK         := lex_opr("?",      "b")
global lex_REVASSIGN     := lex_opr("<-",     "")
global lex_REVSWAP       := lex_opr("<->",    "")
global lex_RBRACE        := lex_opr("}",      "e")
global lex_RBRACK        := lex_opr("]",      "e")
global lex_RCOMP         := lex_opr(":]",     "e")
global lex_RPAREN        := lex_opr(")",      "e")
global lex_SEMICOL       := lex_opr(";",      "")
global lex_SLASH         := lex_opr("/",      "b")
global lex_AUGSLASH      := lex_opr("/:=",    "")
global lex_SLASHSLASH    := lex_opr("//",     "b")
global lex_AUGSLASHSLASH := lex_opr("//:=",   "")
global lex_STAR          := lex_opr("*",      "b")
global lex_AUGSTAR       := lex_opr("*:=",    "")
global lex_SWAP          := lex_opr(":=:",    "")
global lex_TILDE         := lex_opr("~",      "b")
global lex_UNION         := lex_opr("++",     "b")
global lex_AUGUNION      := lex_opr("++:=",   "")

#  lex_lit defines a literal token.
procedure lex_lit(str, flags) {
	return lex_token(str, flags)
}

#  lex_kwd defines a keyword token and enters it in the lex_kwtab table.
procedure lex_kwd(str, flags) {
	return lex_kwtab[str] := lex_token(str, flags)
}

#  lex_opr defines an operator token and enters it in the lex_optab table.
procedure lex_opr(str, flags) {
	return lex_optab[str] := lex_token(str, flags)
}

#  lex_token defines a token record and puts its flags in the lex_flags table.
procedure lex_token(str, flags) {
	^r := lex_tkrec(str, flags)
	lex_flags[r] := flags
	return r
}

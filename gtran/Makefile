#  gtran Makefile

ICONT = icont -u -s

.SUFFIXES: .icn .u2
.icn.u2:	; $(ICONT) -c $<


BASESRC = \
	dump.icn

LEXSRC = \
	do_ops.icn \
	lexer.icn

PARSESRC = \
	ast.icn \
	parse.icn

IRGENSRC = \
	ir.icn \
	keyword.icn \
	irgen.icn

CODEGENSRC = \
	gen_json.icn \
	gen_symbolic.icn \
	gen_dot.icn \
	optimize.icn

TRANSRC = $(BASESRC) \
	$(LEXSRC)  \
	$(PARSESRC) \
	$(IRGENSRC) \
	$(CODEGENSRC) \
	gtran.icn

TRANOBJ = $(TRANSRC:.icn=.u2)


#  the Goaldi translator

default:  gtran

gtran:	$(TRANOBJ)
	$(ICONT) -o gtran $(TRANOBJ)
	cp gtran $${GOPATH%%:*}/bin/gtran


# programs that build programs

oplexgen: oplexgen.icn
	$(ICONT) -o oplexgen oplexgen.icn

do_ops.icn: oplexgen
	./oplexgen > do_ops.icn


#  cleanup

clean:
	rm -rf *.u[12] do_ops.icn oplexgen

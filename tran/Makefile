#  jcon/tran/Makefile -- for building the Jcon translator and linker
#
#  To build translator and linker:	make
#  To build them using Jcon:		make jj
#  To benchmark the two versions:	make bmark


ICONT = icont -u -s
JCONT = ../bin/jcont -s

.SUFFIXES: .icn .u2 .zip
.icn.u2:	; $(ICONT) -c $<
.icn.zip:	; $(JCONT) -c $<


BASESRC = \
	dump.icn

LEXSRC = \
	do_ops.icn \
	preprocessor.icn \
	lexer.icn

PARSESRC = \
	ast.icn \
	parse.icn

IRGENSRC = \
	ir.icn \
	keyword.icn \
	irgen.icn

CODEGENSRC = \
	gen_bc.icn \
	gen_symbolic.icn \
	gen_dot.icn \
	gen_ucode.icn \
	optimize.icn \
	interface.icn \
	bytecode.icn 

TRANSRC = $(BASESRC) \
	$(LEXSRC)  \
	$(PARSESRC) \
	$(IRGENSRC) \
	$(CODEGENSRC) \
	jtran_main.icn

LINKSRC = \
	linker.icn \
	bytecode.icn

TRANOBJ = $(TRANSRC:.icn=.u2)
LINKOBJ = $(LINKSRC:.icn=.u2)

TRANZIP = $(TRANSRC:.icn=.zip)
LINKZIP = $(LINKSRC:.icn=.zip)


# usual versions built by icont

all:	jtran jlink

jtran:	$(TRANOBJ)
	$(ICONT) -o jtran $(TRANOBJ)
	cp jtran ../bin

jlink:	$(LINKOBJ)
	$(ICONT) -o jlink $(LINKOBJ)
	cp jlink ../bin


# versions built by Jcon (not built by default)

jj:	jjtran jjlink

jjtran:	$(TRANZIP)
	$(JCONT) -o jjtran $(TRANZIP)
	cp jjtran ../bin

jjlink:	$(LINKZIP)
	$(JCONT) -o jjlink $(LINKZIP)
	cp jjlink ../bin


# benchmarking

BSRC = irgen.icn
ARGS = preproc $(BSRC) : yylex : parse : ast2ir : \
	bc_File -O -class:'l$$$(BSRC:.icn=)' -dir:./tmp/

bmark:	jtran jjtran
	test -d tmp || mkdir tmp
	rm -rf tmp/*
	BLKSIZE=6000000 STRSIZE=1000000 COEXPSIZE=50000 time ./jtran $(ARGS)
	rm -rf tmp/*
	time ./jjtran $(ARGS)


# programs that build programs

oplexgen: oplexgen.icn
	$(ICONT) -o oplexgen oplexgen.icn

interfacegen: interfacegen.icn
	$(ICONT) -o interfacegen interfacegen.icn

interface.icn: interfacegen
	./interfacegen > interface.icn

do_ops.icn: oplexgen
	./oplexgen > do_ops.icn


#  cleanup

clean:
	rm -rf *.u[12] *.zip tmp do_ops.icn interface.icn 
	rm -rf jtran jjtran jlink jjlink interfacegen oplexgen

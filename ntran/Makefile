##	Makefile for new translator written in Goaldi
##
##	EXPERIMENTAL -- NOT CURRENTLY IN USE

TESTFILE = test.gd

SRC = ast.gd gen_json.gd gtran.gd ir.gd irgen.gd lex.gd optimize.gd parse.gd
GIR = $(SRC:.gd=.gir)
GOBIN = ../../../bin

.SUFFIXES:	.gd .gir
.gd.gir:	; goaldi -c $<

default:	runtest

runtest:	ntran $(TESTFILE)
			gexec -D $(GIR) -- $(TESTFILE)

ntran:		$(GIR) 

$(GIR):		$(GOBIN)/goaldi $(GOBIN)/gtran

clean:
			rm -f *.gir *.json

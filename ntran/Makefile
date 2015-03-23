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
			./ntran $(TESTFILE)
			gexec $(TESTFILE:.gd=.gir)

ntran:		$(GIR)
			echo '#!/usr/bin/env gexec' | cat - $(GIR) >ntran
			chmod +x ntran
			cp ntran $(GOBIN)

$(GIR):		$(GOBIN)/gtran

clean:
			rm -f *.gir ntran

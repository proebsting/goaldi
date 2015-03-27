##	Makefile for new translator written in Goaldi
##
##	EXPERIMENTAL -- NOT CURRENTLY THE DEFAULT TRANSLATOR

SRC = gtran.gd ast.gd ir.gd \
	lex.gd parse.gd irgen.gd optimize.gd gen_json.gd gengo.gd
GIR = $(SRC:.gd=.gir)
GOBIN = ../../../bin

.SUFFIXES:	.gd .gir
.gd.gir:	; goaldi -c $<

ntran:		$(GIR)
			echo '#!/usr/bin/env gexec' | cat - $(GIR) >ntran
			gexec -l ntran			# check for link errors
			chmod +x ntran
			cp ntran $(GOBIN)

$(GIR):		$(GOBIN)/gtran

clean:
			rm -f *.gir ntran

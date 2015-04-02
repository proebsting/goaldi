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
			echo '#!/usr/bin/env gexec'				>ntran
			echo "# ntran `git rev-parse head`"		>>ntran
			echo "# $$USER `date`"					>>ntran
			echo "# `uname -n -s -m`"				>>ntran
			cat $(GIR) | bzip2 -9					>>ntran
			gexec -l ntran			# check for link errors
			chmod +x ntran
			head -4 ntran
			cp ntran $(GOBIN)

$(GIR):		$(GOBIN)/gtran

clean:
			rm -f *.gir ntran

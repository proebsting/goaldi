##	Makefile for new translator written in Goaldi

GEN=X

SRC = gtran.gd ast.gd ir.gd \
	lex.gd parse.gd irgen.gd optimize.gd gen_json.gd gengo.gd
GIR = $(SRC:.gd=.gir)
GOBIN = ../../../bin

.SUFFIXES:	.gd .gir
.gd.gir:	; goaldi -c $<

ntran.go:	ntran gobytes.sh
			./gobytes.sh ntran AppCode <ntran >ntran.go

oldbed:
			./gobytes.sh ntran AppCode <stable.gix >ntran.go

ntran:		$(GIR)
			echo '#!/usr/bin/env gexec -.!!'				>ntran
			echo "# ntran gen$(GEN) `git rev-parse HEAD`"	>>ntran
			echo "# $$USER `date`"							>>ntran
			echo "# `uname -n -s -m`"						>>ntran
			cat $(GIR) | bzip2 -9							>>ntran
			gexec -.!! -l ntran			# check for link errors
			chmod +x ntran
			head -4 ntran
			cp ntran $(GOBIN)

$(GIR):		$(GOBIN)/gexec

clean:
			rm -f *.gir ntran ntran.go

##	Makefile for Goaldi translator

GEN=X

SRC = gtran.gd ast.gd ir.gd \
	lex.gd parse.gd irgen.gd optimize.gd gen_json.gd gengo.gd
GIR = $(SRC:.gd=.gir)
GOBIN = ../../../bin

#  rule for compiling a .gd file to make a .gir file
.SUFFIXES:	.gd .gir
.gd.gir:	; goaldi -c $<

#  make Go source file for embedding the translator
ntran.go:	ntran gobytes.sh
			./gobytes.sh ntran AppCode <ntran >ntran.go

#  make bootstrap Go file from saved, stable translator version
oldbed:
			./gobytes.sh ntran AppCode <stable.gix >ntran.go

#  make translator executable from component .gir files
ntran:		$(GIR)
			echo '#!/usr/bin/env gexec -x'					>ntran
			echo "# ntran gen$(GEN) `git rev-parse HEAD`"	>>ntran
			echo "# $$USER `date`"							>>ntran
			echo "# `uname -n -s -m`"						>>ntran
			cat $(GIR) | bzip2 -9							>>ntran
			gexec -x -l ntran			# check for link errors
			chmod +x ntran
			head -4 ntran

#  install the new translator as the stable version for future builds
accept:		ntran
			@echo '------'
			@head -4 stable.gix | sed 's/^#/< /'
			@echo '------'
			@head -4 ntran      | sed 's/^#/> /'
			@echo '------'
			cp -p ntran stable.gix

#  clean up build products
clean:
			rm -f *.gir ntran ntran.go

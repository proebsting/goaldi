##	Makefile for Goaldi translator

SRC = main.gd ast.gd ir.gd \
	lex.gd parse.gd irgen.gd optimize.gd gen_json.gd gengo.gd
GIR = $(SRC:.gd=.gir)
GOBIN = ../../../bin
GEN = X

#  rule for compiling a .gd file to make a .gir file
.SUFFIXES:	.gd .gir
.gd.gir:	; goaldi -c $<

#  make Go source file for embedding the translator
gtran.go:	gtran gobytes.sh
			./gobytes.sh tran GCode <gtran >gtran.go

#  make bootstrap Go file from saved, stable translator version
boot:
			./gobytes.sh tran GCode <stable-gtran >gtran.go

#  make translator executable from component .gir files
gtran:	$(GIR)
			echo '#!/usr/bin/env goaldi -x'					>gtran
			echo "# gtran gen$(GEN) `git rev-parse HEAD`"	>>gtran
			echo "# $$USER `date`"							>>gtran
			echo "# `uname -n -s -m`"						>>gtran
			cat $(GIR) | bzip2 -9							>>gtran
			goaldi -x -l gtran			# check for link errors
			chmod +x gtran
			head -4 gtran

#  install the new translator as the stable version for future builds
accept:	gtran
			@echo '------'
			@head -4 stable-gtran | sed 's/^#/< /'
			@echo '------'
			@head -4 gtran        | sed 's/^#/> /'
			@echo '------'
			cp -p gtran stable-gtran

#  clean up build products
clean:
			rm -f *.gir gtran gtran.go

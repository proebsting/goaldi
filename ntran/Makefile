##	Makefile for Goaldi translator

GEN=X

SRC = gmain.gd ast.gd ir.gd \
	lex.gd parse.gd irgen.gd optimize.gd gen_json.gd gengo.gd
GIR = $(SRC:.gd=.gir)
GOBIN = ../../../bin

#  rule for compiling a .gd file to make a .gir file
.SUFFIXES:	.gd .gir
.gd.gir:	; goaldi -c $<

#  make Go source file for embedding the translator
gmain.go:	gmain gobytes.sh
			./gobytes.sh ntran AppCode <gmain >gmain.go

#  make bootstrap Go file from saved, stable translator version
oldbed:
			./gobytes.sh ntran AppCode <stable.gix >gmain.go

#  make translator executable from component .gir files
gmain:	$(GIR)
			echo '#!/usr/bin/env goaldi -x'					>gmain
			echo "# gmain gen$(GEN) `git rev-parse HEAD`"	>>gmain
			echo "# $$USER `date`"							>>gmain
			echo "# `uname -n -s -m`"						>>gmain
			cat $(GIR) | bzip2 -9							>>gmain
			goaldi -x -l gmain			# check for link errors
			chmod +x gmain
			head -4 gmain

#  install the new translator as the stable version for future builds
accept:	gmain
			@echo '------'
			@head -4 stable.gix | sed 's/^#/< /'
			@echo '------'
			@head -4 gmain      | sed 's/^#/> /'
			@echo '------'
			cp -p gmain stable.gix

#  clean up build products
clean:
			rm -f *.gir gmain gmain.go

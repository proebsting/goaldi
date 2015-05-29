##	Makefile for Goaldi translator

SRC = main.gd ast.gd ir.gd \
	lex.gd parse.gd irgen.gd optimize.gd gen_json.gd gengo.gd
GIR = $(SRC:.gd=.gir)
GEN = X
GOALDI = goaldi


#  rule for compiling a .gd file to make a .gir file
.SUFFIXES:	.gd .gir
.gd.gir:	; $(GOALDI) -c $<

#  make Go source file for embedding the translator
gtran.go:	gtran0 gtran gobytes.sh
			./gobytes.sh tran GCode <gtran >gtran.go

#  make translator executable from component .gir files
gtran:	$(GIR)
			echo '#!/usr/bin/env goaldi -x'					>gtran
			echo "# gtran gen$(GEN) `git rev-parse HEAD`"	>>gtran
			echo "# $$USER `date`"							>>gtran
			echo "# `uname -n -s -m`"						>>gtran
			cat $(GIR) | bzip2 -9							>>gtran
			chmod +x gtran

#  if gtran0 doesn't exist (e.g. after make clean) force full two-pass rebuild
gtran0:
			rm -f *.gir
			$(GOALDI) -c $(SRC)
			$(MAKE) gtran
			$(GOALDI) -x gtran -c $(SRC)
			mv gtran gtran0

#  install the new translator as the stable version for future builds
accept:	gtran
			@echo '------'
			@head -4 stable-gtran | sed 's/^#/< /'
			@echo '------'
			@head -4 gtran        | sed 's/^#/> /'
			@echo '------'
			cp -p gtran stable-gtran

#  make bootstrap Go file from saved, stable translator version
boot:
			./gobytes.sh tran GCode <stable-gtran >gtran.go

#  get profiling data (by running gtran to compile itself) and start pprof
profile: gtran
			$(GOALDI) -x -P -t gtran -t -c *.gd
			: suggested pprof commands: "top25" or "weblist main"
			go tool pprof `command -v $(GOALDI)` ./PROFILE

#  clean up build products
clean:
			rm -f *.gir gtran0 gtran gtran.go PROFILE

#  tests/Makefile -- Goaldi test files

GOALDI = ../goaldi

default:  runall

#  run all tests
runall:
	GOALDI=$(GOALDI) ./runtest.sh

#  run a quick sanity check:
quick:
	GOALDI=$(GOALDI) ./runtest.sh alltypes


gir:	; +make F=c X=gir derivatives   # make .gir for every .gd file
gia:	; +make F=a X=gia derivatives   # make .gia for every .gd file

F=c
X=gir
derivatives:	# run goaldi -$F on every *.gd to make *.$X
	for SRC in *.gd; do B=$${SRC%.gd}; \
		(set -x; $(GOALDI) -$F $$B.gd 2>$$B.err || rm -f FAILED: $$B.$X); \
	done
	rm -f `find * -name '*.err' -size 0`
	ls *.$X

clean:
	rm -f *.tmp *.out *.err *.gir *.gia *.dot

#  gtests/Makefile -- Goaldi test files

default:  runall

#  run all tests
runall:
	./runtest.sh

gir:	; $(MAKE) F=c X=gir derivatives   # make .gir for every .gdi file
dot:	; $(MAKE) F=d X=dot derivatives   # make .gir for every .gdi file

F=c
X=gir
derivatives:	# run goaldi -$F on every *.gdi to make *.$X
	for SRC in *.gdi; do B=$${SRC%.gdi}; \
		(set -x; goaldi -$F $$B.gdi 2>$$B.err || rm -f FAILED: $$B.$X); \
	done
	rm -f `find * -name '*.err' -size 0`
	ls *.$X

clean:
	rm -f *.gir *.out *.err *.tmp *.dot

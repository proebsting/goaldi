#  apps/Makefile -- apps directory Makefile
#
#  The default action runs all the demos interactively.

GOALDI = ../goaldi

default:	run

#  build and link without running (useful as an automated test)
link:	
	for F in *.gd; do (set -x; $(GOALDI) -l $$F); done

#  run all the demos interactively.
#  each one must be closed manually.
run:	
	for F in *.gd; do (set -x; $(GOALDI) $$F); done

#  clean up  
clean:
	# nothing to do

#  demos/Makefile -- demo area Makefile
#
#  The default action runs all the demos,
#  but inspecting and running them individually may be more useful.

GOALDI = ../goaldi

default:	run

#  build and link without running (useful as an automated test)
link:	
	for F in *.gd; do (set -x; $(GOALDI) -l $$F); done

#  run all the demos (output is nondeterministic)
run:	
	for F in *.gd; do (set -x; $(GOALDI) $$F); done

#  clean up  
clean:
	# nothing to do

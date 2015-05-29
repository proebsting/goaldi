#  Makefile for Goaldi stdlib documentation

GOALDI = ../goaldi

default:	stdlib.adoc

show:		stdlib.html
			open stdlib.html


#  extract stdlib procedure documentation.
stdlib.adoc:	libdoc.sh libdoc.hdr libdoc.gd .FORCE
	GOALDI=$(GOALDI) ./libdoc.sh >stdlib.tmp
	mv stdlib.tmp stdlib.adoc


#  make HTML documentation.
stdlib.html:	stdlib.adoc gdstyle.css
	asciidoctor -a stylesheet=gdstyle.css stdlib.adoc


#  clean up
#
#  The build product libdoc.adoc is not removed because it gets checked in.
clean:
	rm -f stdlib.tmp stdlib.html


.FORCE:

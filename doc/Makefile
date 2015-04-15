#  Makefile for Goaldi stdlib documentation

default:	stdlib.adoc

show:		stdlib.html
			open stdlib.html


#  extract stdlib procedure documentation.
#  gets data from installed "goaldi" but does not specify explicit dependency.
stdlib.adoc:	libdoc.sh libdoc.hdr libdoc.gd .FORCE
	./libdoc.sh >stdlib.adoc


#  make HTML documentation.
stdlib.html:	stdlib.adoc gdstyle.css
	asciidoctor -a stylesheet=gdstyle.css stdlib.adoc


#  clean up
#
#  The build product libdoc.adoc is not removed because it gets checked in.
clean:
	rm -f stdlib.html


.FORCE:

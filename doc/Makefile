#  Makefile for Goaldi stdlib documentation

default:	stdlib.adoc

show:		stdlib.html
			open stdlib.html


#  extract stdlib procedure documentation.
#  gets data from installed "goaldi" but does not specify explicit dependency.
stdlib.adoc:	libdoc.sh libdoc.gd 
	./libdoc.sh >stdlib.adoc


#  make HTML documentation and open in browser.
#  assumes asciidoctor is installed.
stdlib.html:	stdlib.adoc gdstyle.css
	asciidoctor -a stylesheet=gdstyle.css stdlib.adoc
	open stdlib.html


#  clean up
#
#  The build product libdoc.adoc is not removed because it gets checked in,
#  but libdoc.gd is touched to force a rebuild of stdlib.adoc next time.
clean:
	rm -f stdlib.html
	touch libdoc.gd

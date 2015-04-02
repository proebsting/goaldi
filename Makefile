#  Goaldi Makefile
#
#  Assumptions:
#	$GOPATH is set per Go documentation
#	$GOPATH/bin (first GOPATH component) is destination for built programs
#	go, icon, and $GOPATH/bin are all in search path
#
#  Additional temporary assumption:  working goaldi in search path.

PKG = goaldi
PROGS = $(PKG)/gexec
# GOBIN expands in shell to {first component of $GOPATH}/bin
GOBIN = $${GOPATH%%:*}/bin

#  default action: set up, build all, run test suite, run expt.gd if present
default:  setup build test buildx expt

#  configure Git pre-commit hook
HOOKMASTER = ./pre-commit.hook
HOOKFILE = .git/hooks/pre-commit
setup:	$(HOOKFILE)
$(HOOKFILE):	$(HOOKMASTER)
	cp $(HOOKMASTER) $(HOOKFILE)

#  build and install Goaldi
build: embed
	go install $(PROGS)
	cd gtran; $(MAKE)
	cp goaldi.sh $(GOBIN)/goaldi

#  make a sample embedded app
embed:
	goaldi -c embedded.gd
	bzip2 -9 <embedded.gir | ./gobytes.sh main appcode >gexec/embed.go

#  build and install ntran (experimental translator)
buildx:
	cd ntran; $(MAKE)

#  run Go unit tests; build and link demos; run Goaldi test suite
test:
	cd runtime; go test
	cd demo; $(MAKE) link
	cd gtests; $(MAKE)

#  run expt.gd (presumably the test of the moment) if present
#  passes $GXOPTS to interpreter if set in environment
expt:
	test -f expt.gd && $(GOBIN)/goaldi $$GXOPTS expt.gd || :

#  prepare Go source for check-in by running standard Go reformatter
format:
	go fmt *.go
	go fmt ir/*.go
	go fmt gexec/*.go
	go fmt runtime/*.go
	go fmt extensions/*.go

#  gather together source for single-file editing; requires "bundle" util
bundle:
	@bundle *.go */*.go

#  extract stdlib procedure documentation
libdoc:	libdoc.txt
	@: # don't try to build ./libdoc, it's just an alias
libdoc.txt:	libdoc.sh libdoc.gd build
	./libdoc.sh >libdoc.txt

#  remove temporary and built files from source tree
clean:
	rm -f embedded.gir gexec/embed.go
	rm -f libdoc.txt
	go clean $(PKG) $(PROGS)
	cd gtran; $(MAKE) clean
	cd ntran; $(MAKE) clean
	cd gtests; $(MAKE) clean

#  remove files placed elsewhere in $GOPATH
uninstall:
	rm -f $(GOBIN)/gtran $(GOBIN)/gexec $(GOBIN)/goaldi $(GOBIN)/ntran
	go clean -i $(PKG) $(PROGS)

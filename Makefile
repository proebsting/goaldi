#  Goaldi Makefile
#
#  Assumptions:
#	$GOPATH is set per Go documentation
#	$GOPATH/bin (first GOPATH component) is destination for built programs
#	go, icon, and $GOPATH/bin are all in search path

PKG = goaldi
PROGS = $(PKG)/gexec
# GOBIN expands in shell to {first component of $GOPATH}/bin
GOBIN = $${GOPATH%%:*}/bin

#  default action: set up, build all, run test suite, run expt.gdi if present
default:  setup build test expt

#  configure Git pre-commit hook
HOOKMASTER = ./pre-commit.hook
HOOKFILE = .git/hooks/pre-commit
setup:	$(HOOKFILE)
$(HOOKFILE):	$(HOOKMASTER)
	cp $(HOOKMASTER) $(HOOKFILE)

#  build and install Goaldi
build:
	go install $(PROGS)
	cd gtran; $(MAKE)
	cp goaldi.sh $(GOBIN)/goaldi

#  run Go unit tests and Goaldi test suite
test:
	go test
	cd gtests; $(MAKE)

#  run translate-and-link tests of old Icon programs
itest:
	cd itests; $(MAKE)

#  run expt.gdi (presumably the test of the moment) if present
#  passes $GXOPTS to interpreter if set in environment
expt:
	test -f expt.gdi && $(GOBIN)/goaldi $$GXOPTS expt.gdi || :

#  prepare Go source for check-in by running standard Go reformatter
format:
	go fmt *.go
	go fmt gexec/*.go

#  gather together source for single-file editing; requires "bundle" util
bundle:
	@bundle *.go */*.go

#  remove temporary and built files from source tree
clean:
	go clean $(PKG) $(PROGS)
	cd gtran; $(MAKE) clean
	cd gtests; $(MAKE) clean
	cd itests; $(MAKE) clean

#  remove files placed elsewhere in $GOPATH
uninstall:
	rm -f $(GOBIN)/gtran $(GOBIN)/gexec $(GOBIN)/goaldi
	go clean -i $(PKG) $(PROGS)

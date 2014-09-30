#  Goaldi Makefile
#
#  Assumptions:
#	go is in search path (typically located in /usr/local/go/bin)
#	$GOPATH is set per Go documentation
#	$GOPATH/bin (first GOPATH component) is destination for built programs

PKG = goaldi
GEXEC = $(PKG)/gexec
HTESTS = $(PKG)/test1 $(PKG)/test2 $(PKG)/test3 $(PKG)/test9
PROGS = $(GEXEC) $(HTESTS)
# GOBIN expands in shell to {first component of $GOPATH}/bin
GOBIN = $${GOPATH%%:*}/bin

#  default action: build all, run test suite, run expt if present
default:  build test expt

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

#  run early hand-coded test programs
htests:
	$(GOBIN)/test1
	$(GOBIN)/test2
	$(GOBIN)/test3
	$(GOBIN)/test9

#  prepare Go source for check-in by running standard Go reformatter
format:
	go fmt *.go
	for D in gexec test*; do go fmt $$D/*.go; done

#  gather together source for single-file editing; requires "bundle" util
bundle:
	@bundle *.go */*.go

#  remove temporary and built files from source tree
clean:
	go clean $(PKG) $(PROGS)
	go clean -i $(HTESTS)
	cd gtran; $(MAKE) clean
	cd gtests; $(MAKE) clean
	cd itests; $(MAKE) clean

#  remove files placed elsewhere in $GOPATH
uninstall:
	rm -f $(GOBIN)/gtran $(GOBIN)/gexec $(GOBIN)/goaldi
	go clean -i $(PKG) $(PROGS)

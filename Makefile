#  Goaldi Makefile
#
#  Assumptions:
#	$GOPATH is set per Go documentation
#	$GOPATH/bin (first GOPATH component) is destination for built programs
#	$GOPATH/bin are in search path, as is the Go compiler

PKG = goaldi
PROGS = $(PKG)/goaldi
# GOBIN expands in shell to {first component of $GOPATH}/bin
GOBIN = $${GOPATH%%:*}/bin

#  default action: set up, build all, run test suite, run expt.gd if present
default:  setup quick doc test expt

#  configure Git pre-commit hook
HOOKMASTER = ./pre-commit.hook
HOOKFILE = .git/hooks/pre-commit
setup:	$(HOOKFILE)
$(HOOKFILE):	$(HOOKMASTER)
	cp $(HOOKMASTER) $(HOOKFILE)

#  quick rebuild using available translator
quick:
	$(GOBIN)/goaldi -x -l /dev/null || $(MAKE) boot
	cd translator; $(MAKE)
	go install $(PROGS)
	cd tests; $(MAKE) quick

#  bootstrap build goaldi using stable translator binary
boot:
	cd translator; $(MAKE) boot
	go install $(PROGS)

#  full three-pass bootstrap process
full:
	#
	# make an executable that embeds an old version of the front end
	#
	cd translator; $(MAKE) boot
	go install $(PROGS)
	cd runtime; go test
	cd tests; $(MAKE) quick
	#
	# make an executable embedding the latest front end, built by the old one
	#
	cd translator; $(MAKE) clean; $(MAKE) GEN=1 gtran.go
	go install $(PROGS)
	cd tests; $(MAKE) quick
	#
	# make an executable embedding the latest front end as built by itself
	#
	cd translator; $(MAKE) clean; $(MAKE) GEN=2 gtran.go
	go install $(PROGS)
	cd tests; $(MAKE) quick
	#
	# looks like a keeper.
	#

#  extract stdlib documentation from the Goaldi binary
doc:	.FORCE
	cd doc; $(MAKE)

#  run Go unit tests; build and link demos; run Goaldi test suite
test:
	cd runtime; go test
	cd demos; $(MAKE) link
	cd tests; $(MAKE)

#  run expt.gd (presumably the test of the moment) if present
#  passes $GXOPTS to interpreter if set in environment
expt:
	test -f expt.gd && $(GOBIN)/goaldi $$GXOPTS expt.gd || :

#  install the new translator as the stable version for future builds
accept:
	cd translator; $(MAKE) accept

#  prepare Go source for check-in by running standard Go reformatter
format:
	go fmt *.go
	go fmt ir/*.go
	go fmt goaldi/*.go
	go fmt runtime/*.go
	go fmt extensions/*.go

#  gather together source for single-file editing; requires "bundle" util
bundle:
	@bundle `find * -name '*.go' ! -name gtran.go`

#  remove temporary and built files from source tree
#  and also subpackages built and saved in $GOPATH
clean:
	go clean $(PKG) $(PROGS)
	cd translator; $(MAKE) clean
	cd tests; $(MAKE) clean
	cd doc; $(MAKE) clean
	rm -rf $(GOBIN)/../pkg/*/goaldi

#  remove files placed elsewhere in $GOPATH
uninstall:
	go clean -i $(PKG) $(PROGS)
	rm -rf $(GOBIN)/goaldi $(GOBIN)/../pkg/*/goaldi


.FORCE:

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
default:  setup build test expt

#  configure Git pre-commit hook
HOOKMASTER = ./pre-commit.hook
HOOKFILE = .git/hooks/pre-commit
setup:	$(HOOKFILE)
$(HOOKFILE):	$(HOOKMASTER)
	cp $(HOOKMASTER) $(HOOKFILE)

#  build and install Goaldi
build:
	cp goaldi.sh $(GOBIN)/goaldi
	#
	# make an executable that embeds an old version of the front end
	#
	cd ntran; $(MAKE) oldbed
	go install $(PROGS)
	cd runtime; go test
	cd gtests; $(MAKE) quick
	#
	# make an executable embedding the latest front end, built by the old one
	#
	cd ntran; $(MAKE) clean; $(MAKE) GEN=1 ntran.go
	go install $(PROGS)
	cd gtests; $(MAKE) quick
	#
	# make an executable embedding the latest front end as built by itself
	#
	cd ntran; $(MAKE) clean; $(MAKE) GEN=2 ntran.go
	go install $(PROGS)
	cd gtests; $(MAKE) quick
	#
	# looks like a keeper.
	#

gexec/embed.go: ntran/ntran gobytes.sh
	./gobytes.sh main appcode <ntran/ntran >gexec/embed.go

ntran/ntran:
	cp ntran/stable.gix ntran/ntran

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
	rm -f libdoc.txt
	go clean $(PKG) $(PROGS)
	cd ntran; $(MAKE) clean
	cd gtests; $(MAKE) clean

#  remove files placed elsewhere in $GOPATH
uninstall:
	rm -f $(GOBIN)/gexec $(GOBIN)/goaldi $(GOBIN)/ntran
	go clean -i $(PKG) $(PROGS)

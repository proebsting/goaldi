#  Goaldi Makefile
#
#  ASSUMPTIONS:
#
#  (1) Go is installed and in the search path, and $GOPATH is set properly.
#      IF NOT: Install Go before proceeding, and set $GOPATH.
#
#  (2) A functioning “goaldi” executable is in the search path.
#      IF NOT: Remove any nonfunctional goaldi executable from the path.
#      Then run “make boot” in a clean, unedited release package to build
#      and install in $GOPATH/bin from a stable intermediate snapshot.
#
#  (3) Any intermediate build products are functional and mutually compatible.
#      IF NOT: Run “make clean” to remove them after a failed build
#      or after making incompatible changes to the intermediate representation.


PKG = goaldi
PROGS = $(PKG)/goaldi
# GOBIN expands in the shell to {first component of $GOPATH}/bin
GOBIN = $${GOPATH%%:*}/bin
# a Git pre-commit hook validates formatting before check-in
HOOKMASTER = ./pre-commit.hook
HOOKFILE = .git/hooks/pre-commit


#  -- shorthand targets --

#  default: configure, build all, run tests, run expt.gd if present
default:	setup build test expt

#  setup:  prepare for building and developing
setup:  gcommand $(HOOKFILE)

#  build:  make the goaldi executable and the extracted documentation
build:	gexe doc

#  quick rebuild and test
quick:		build qktest expt


#  -- setup targets --

gcommand:	# ensure that we have a "goaldi" command
	+command -v goaldi >/dev/null || $(MAKE) boot

boot:		# install goaldi using stable pre-built translator IR code
	cd tran; $(MAKE) boot
	go build -o gexe $(PROGS)
	$(MAKE) install
	rm -f tran/gtran.go gexe

$(HOOKFILE): $(HOOKMASTER)	# configure Git pre-commit hook
	test -d .git && cp $(HOOKMASTER) $(HOOKFILE) || :


#  -- build targets --

gexe: 
	cd tran; $(MAKE)
	go build -o gexe $(PROGS)
	./gexe -l /dev/null	# validate build

install: gexe
	./gexe -l /dev/null	# validate build
	cp ./gexe $(GOBIN)/goaldi

doc:	.FORCE
	cd doc; $(MAKE)


#  full three-pass rebuild using bootstrapping from old stable front end
#%#% this does some intermediate installs of untested code
full:
	#
	# make an executable that embeds an old version of the front end
	#
	cd tran; $(MAKE) boot
	cd runtime; go test
	go install $(PROGS)
	cd tests; $(MAKE) quick
	#
	# make an executable embedding the latest front end, built by the old one
	#
	cd tran; $(MAKE) clean; $(MAKE) GEN=1 gtran.go
	go install $(PROGS)
	cd tests; $(MAKE) quick
	#
	# make an executable embedding the latest front end as built by itself
	#
	cd tran; $(MAKE) clean; $(MAKE) GEN=2 gtran.go
	go install $(PROGS)
	$(MAKE) doc
	cd tests; $(MAKE) quick
	#
	# looks good in quick tests; now run full suite
	#
	$(MAKE) test


#  -- testing targets --

#  run Go unit tests; build and link demos; run Goaldi test suite
test:
	cd runtime; go test
	cd demos; $(MAKE) link
	cd tests; $(MAKE)

#  run a single quick test
qktest:
	cd tests; $(MAKE) quick

#  run expt.gd (presumably the test of the moment) if present
#  passes $GXOPTS to interpreter if set in environment
expt:
	test -f expt.gd && ./gexe $$GXOPTS expt.gd || :

#  run demo programs (non-automated, with output to stdout)
demos: .FORCE
	cd demos; $(MAKE)


#  -- miscellaneous targets --

#  prepare Go source for check-in by running standard Go reformatter
format:
	go fmt *.go
	go fmt ir/*.go
	go fmt goaldi/*.go
	go fmt runtime/*.go
	go fmt extensions/*.go

#  gather together Go source for single-file editing; requires "bundle" utility
bundle:
	@bundle `find * -name '*.go' ! -name gtran.go`

#  install the newly built translator as the stable version for future builds
#  (be sure this is a good one or you'll lose the ability to bootstrap)
accept:
	cd tran; $(MAKE) accept

#  a prerequesite target to force unconditional rebuild of dependent items
.FORCE:


#  -- cleanup targets --

#  remove temporary and built files from source tree
#  and also subpackages built and saved in $GOPATH
clean:
	go clean $(PKG) $(PROGS)
	cd tran; $(MAKE) clean
	cd tests; $(MAKE) clean
	cd doc; $(MAKE) clean
	rm -rf $(GOBIN)/../pkg/*/goaldi
	rm -rf gexe Goaldi-*-*

#  remove files placed elsewhere in $GOPATH
uninstall:
	go clean -i $(PKG) $(PROGS)
	rm -rf $(GOBIN)/goaldi $(GOBIN)/../pkg/*/goaldi

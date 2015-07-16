#  Goaldi Makefile
#
#  ASSUMPTIONS:
#
#  (1) Go is installed and in the search path, and $GOPATH is set properly.
#      IF NOT: Install Go before proceeding, and set $GOPATH.
#
#  (2) A functioning "goaldi" executable is in the search path.
#      IF NOT: Remove any nonfunctional goaldi executable from the path.
#      Then run "make boot" in a clean, unedited release package to build
#      and install in $GOPATH/bin from a stable intermediate snapshot.
#
#  (3) Any intermediate build products are functional and mutually compatible.
#      IF NOT: Run "make clean" to remove them after a failed build
#      or after making incompatible changes to the intermediate representation.


PKG = goaldi
PROGS = $(PKG)/interp
# GOBIN expands in the shell to {first component of $GOPATH}/bin
GOBIN = $${GOPATH%%:*}/bin
# a Git pre-commit hook validates formatting before check-in
HOOKMASTER = ./pre-commit.hook
HOOKFILE = .git/hooks/pre-commit


#  -- shorthand targets --

#  default: configure, build all, run tests, run expt.gd if present
default:	setup build test expt

#  setup:  prepare for building and developing
setup:  gcommands $(HOOKFILE)

#  build:  make the goaldi executable and the extracted documentation
build:	goaldi doc

#  quick rebuild and test
quick:		build qktest expt


#  -- setup targets --

gcommands:	# ensure that we have "go" and "goaldi" commands
	@+command -v go >/dev/null || \
		(echo 'No "go" command in $$PATH;' \
		'install Go before proceeding' && exit 1)
	@+command -v goaldi >/dev/null || \
		(echo 'No "goaldi" command in $$PATH;' \
		'run "make boot" to bootstrap into $$GOPATH/bin' && exit 1)

boot:		# install goaldi using stable pre-built translator IR code
	+cd tran; make boot
	go build -o goaldi $(PROGS)
	+make install
	rm -f tran/gtran.go goaldi
	ls -l $(GOBIN)/goaldi
	: stable version installed successfully for bootstrapping

$(HOOKFILE): $(HOOKMASTER)	# configure Git pre-commit hook
	test -d .git && cp $(HOOKMASTER) $(HOOKFILE) || :

getlibs:	# get latest versions of supplemental libraries
	go get -u -x golang.org/x/mobile/app
	go get -u -x golang.org/x/mobile/exp/gl/glutil
	go get -u -x code.google.com/p/freetype-go/freetype/truetype


#  -- build targets --

goaldi: .FORCE
	+cd tran; make
	go build -o goaldi $(PROGS)
	./goaldi -l /dev/null	# validate build

install:
	./goaldi -l /dev/null	# validate binary
	cp ./goaldi $(GOBIN)/goaldi

doc:	.FORCE
	+cd doc; make

#  self: rebuild using already-built goaldi
#  (this confirms that the latest version can build itself)
self:
	rm -f tran/*.gir
	cd tran; make GOALDI=../goaldi
	go build -o goaldi $(PROGS)
	./goaldi -l /dev/null	# validate build


#  -- testing targets --

#  run Go unit tests; build and link demos; run Goaldi test suite
test:
	cd runtime; go test
	+cd demos; make link
	+cd apps; make link
	+cd tests; make

#  run a single quick test
qktest:
	+cd tests; make quick

#  run expt.gd (presumably the test of the moment) if present
#  passes $GXOPTS to interpreter if set in environment
expt:
	test -f expt.gd && ./goaldi $$GXOPTS expt.gd || :

#  run demo programs (non-automated, with output to stdout)
demos: .FORCE
	+cd demos; make

#  run apps (non-automated, with output to stdout)
apps: .FORCE
	+cd apps; make


#  -- miscellaneous targets --

#  prepare Go source for check-in by running standard Go reformatter
format:
	go fmt *.go
	go fmt ir/*.go
	go fmt interp/*.go
	go fmt runtime/*.go
	go fmt graphics/*.go
	go fmt extensions/*.go

#  gather together Go source for single-file editing; requires "bundle" utility
bundle:
	@bundle `find * -name '*.go' ! -name gtran.go`

#  install the newly built translator as the stable version for future builds
#  (be sure this is a good one or you'll lose the ability to bootstrap)
accept:
	+cd tran; make accept

#  a prerequisite target to force unconditional rebuild of dependent items
.FORCE:


#  -- cleanup targets --

#  remove temporary and built files from source tree
#  and also subpackages built and saved in $GOPATH
clean:
	go clean $(PKG) $(PROGS)
	+cd tran; make clean
	+cd tests; make clean
	+cd doc; make clean
	rm -rf $(GOBIN)/../pkg/*/goaldi
	rm -rf goaldi Goaldi-*-*

#  remove files placed elsewhere in $GOPATH
uninstall:
	go clean -i $(PKG) $(PROGS)
	rm -rf $(GOBIN)/goaldi $(GOBIN)/../pkg/*/goaldi

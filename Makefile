#  Goaldi Makefile
#
#  Assumptions:
#	go is in search path (typically located in /usr/local/go/bin)
#	$GOPATH is set  
#	$GOPATH/bin (first GOPATH component) is destination for built programs

PKG = goaldi
MAIN = $(PKG)/gexec
TESTS = $(PKG)/test1 $(PKG)/test2 $(PKG)/test3 $(PKG)/test9
PROGS = $(MAIN) $(TESTS)
GOBIN = $${GOPATH%%:*}/bin

default:  build test

build:
	go install $(PROGS)
	cd gtran; $(MAKE)
	cp goaldi.sh $(GOBIN)/goaldi

test:
	go test
	cd itests; $(MAKE)

format:
	go fmt *.go
	for D in exec test*; do go fmt $$D/*.go; done

mains:					# early test drivers
	$(GOBIN)/test1
	$(GOBIN)/test2
	$(GOBIN)/test3
	$(GOBIN)/test9

bundle:
	@bundle *.go */*.go

clean:
	go clean $(PKG) $(PROGS)
	go clean -i $(TESTS)
	cd gtran; $(MAKE) clean
	cd itests; $(MAKE) clean

uninstall:
	rm -f $(GOBIN)/gtran $(GOBIN)/gexec $(GOBIN)/goaldi
	go clean -i $(PKG) $(PROGS)

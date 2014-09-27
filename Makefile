#  Goaldi Makefile
#  (a work in progress -- expect drastic changes and reorganizations)

PKG = goaldi
MAIN = $(PKG)/gexec
TESTS = $(PKG)/test1 $(PKG)/test2 $(PKG)/test3 $(PKG)/test9
PROGS = $(MAIN) $(TESTS)
GOBIN = $${GOPATH%%:*}/bin

default:  build test

build:
	go install $(PROGS)
	cd gtran; $(MAKE)

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
	rm -f $(GOBIN)/gtran
	go clean -i $(PKG) $(PROGS)
	cd gtran; $(MAKE) clean
	cd itests; $(MAKE) clean

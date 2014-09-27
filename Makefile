#  Goaldi Makefile
#  (a work in progress -- expect drastic changes and reorganizations)

PKG = goaldi
MAIN = $(PKG)/terp
TESTS = $(PKG)/test1 $(PKG)/test2 $(PKG)/test3 $(PKG)/test9
PROGS = $(MAIN) $(TESTS)
GOBIN = $${GOPATH%%:*}/bin

default:  build test

build:
	go install $(PROGS)
	cd tran; $(MAKE)

test:
	go test
	cd itests; $(MAKE)

format:	
	go fmt *.go
	for D in terp test*; do go fmt $$D/*.go; done

mains:					# early test drivers
	$(GOBIN)/test1
	$(GOBIN)/test2
	$(GOBIN)/test3
	$(GOBIN)/test9

bundle:
	@bundle *.go */*.go

clean:
	go clean -i $(PKG) $(PROGS)
	cd tran; $(MAKE) clean
	cd itests; $(MAKE) clean

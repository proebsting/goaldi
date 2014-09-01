
PKG = goaldi
PROGS = $(PKG)/test1 $(PKG)/test2 $(PKG)/test9
GOBIN = $$GOPATH/bin

default:  build
	$$GOPATH/bin/test1
	$$GOPATH/bin/test2
	$$GOPATH/bin/test9

build:
	go install $(PROGS)

format:	
	go fmt *.go
	for D in test*; do go fmt $$D/*.go; done

bundle:
	@bundle *.go */*.go

clean:
	go clean -i $(PKG) $(PROGS)

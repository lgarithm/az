GOPATH = $(shell pwd)/.gopath
PKG = $(shell cat .goimportpath)

main:
	mkdir -p bin
	mkdir -p $(GOPATH)/src/$(shell dirname $(PKG))
	rm -fr $(GOPATH)/src/$(PKG) && ln -s $(shell pwd) $(GOPATH)/src/$(PKG)
	rm -fr $(GOPATH)/bin && ln -s $(shell pwd)/bin $(GOPATH)/bin

	GOPATH=$(GOPATH) \
		go install -v $(PKG)/cmd/example-1

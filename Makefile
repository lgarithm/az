binaries:
	GOBIN=$(CURDIR)/bin go install -v ./cmd/...

clean:
	go clean -cache

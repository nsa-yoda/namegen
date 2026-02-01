BINDIR=bin
GO?=go

.PHONY: all build clean

all: build

build:
	cd cmd/namegen && GO111MODULE=on $(GO) build -o ../../$(BINDIR)/namegen

clean:
	rm -rf $(BINDIR)/*

check:
	go vet ./...
	test -z "$$(gofmt -l .)"
	go test ./...

distcheck: check
	go mod tidy
	git diff --exit-code
	make clean
	make build
	./bin/namegen -p
	./bin/namegen -mode english -s 123 | sha256sum
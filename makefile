BINDIR=bin
GO?=go

.PHONY: all build clean

all: build

ci: check distcheck

fuckmeup: clean test-race staticcheck vulncheck fmt lint-fmt check distcheck cross build

build:
	cd cmd/namegen && GO111MODULE=on $(GO) build -o ../../$(BINDIR)/namegen

clean:
	rm -rf $(BINDIR)/*

test-race:
	go test -race ./...

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

fmt:
	gofmt -w .

lint-fmt:
	test -z "$$(gofmt -l .)"

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
	./bin/namegen -mode english -s 123 | sha256sum | grep f8b26d75173ccc67bfc3e4d40bada8007d48b049d4caff23b850ab1adf682f00

# This catches any accidental platform-specific assumptions
cross:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINDIR)/namegen-linux-amd64 ./cmd/namegen
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINDIR)/namegen-darwin-amd64 ./cmd/namegen
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BINDIR)/namegen-darwin-arm64 ./cmd/namegen
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINDIR)/namegen-windows-amd64.exe ./cmd/namegen
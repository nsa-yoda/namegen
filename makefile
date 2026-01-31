BINDIR=bin
GO?=go

.PHONY: all build clean

all: build

build:
	cd cmd/namegen && GO111MODULE=on $(GO) build -o ../../$(BINDIR)/namegen

clean:
	rm -rf $(BINDIR)/*
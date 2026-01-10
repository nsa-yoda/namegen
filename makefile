PLUGIN_DIR=plugins
BINDIR=bin

.PHONY: all build main plugins clean

all: build

build: main plugins

main:
	cd cmd/namegen && GO111MODULE=on go build -o ../../$(BINDIR)/namegen

plugins:
	@echo "Building plugins..."
	# English
	cd $(PLUGIN_DIR)/english && GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ../../$(PLUGIN_DIR)/english.so
	# Japanese
	cd $(PLUGIN_DIR)/japanese && GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ../../$(PLUGIN_DIR)/japanese.so
	# Spanish
	cd $(PLUGIN_DIR)/spanish && GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ../../$(PLUGIN_DIR)/spanish.so

clean:
	rm -rf $(BINDIR)/*
	rm -f $(PLUGIN_DIR)/*.so
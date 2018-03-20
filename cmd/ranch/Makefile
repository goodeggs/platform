.PHONY: all test

SHELL=/bin/bash -o pipefail
export GOPATH=$(shell pwd)/.workspace
export PATH := $(GOPATH)/bin:$(PATH)

all: test

test:
	@pushd .workspace/src/github.com/goodeggs/platform/cmd/ranch > /dev/null && \
		go test -v ./cmd ./util && \
		popd > /dev/null

build:
	go build

install:
	go build -o .workspace/bin/ranch

build_deps: clean
	@mkdir -p ./.workspace/src/github.com/goodeggs/platform/cmd
	@[ -s .workspace/src/github.com/goodeggs/platform/cmd/ranch ] || ln -s $$(pwd) .workspace/src/github.com/goodeggs/platform/cmd/ranch
	go get -u github.com/mitchellh/gox
	go get -u github.com/tcnksm/ghr
	go get -u github.com/sanbornm/go-selfupdate
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/Clever/gitsem
	@pushd .workspace/src/github.com/goodeggs/platform/cmd/ranch > /dev/null && \
		dep ensure && \
		popd > /dev/null
	@echo 'Ready to build!'

clean:
	rm -fR ./vendor
	rm -fR ./vendor.orig
	rm -fR ./.workspace

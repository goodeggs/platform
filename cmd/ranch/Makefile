.PHONY: all test

SHELL=/bin/bash -o pipefail
export GOPATH=$(shell pwd)/.workspace
export PATH := $(GOPATH)/bin:$(PATH)

all: test

test: build_deps
	go test -v ./cmd ./util

build: build_deps
	go build -o .workspace/bin/ranch

build_deps:
	@mkdir -p ./.workspace/src/github.com/goodeggs/platform/cmd
	@ln -s $$(pwd) .workspace/src/github.com/goodeggs/platform/cmd/ranch
	go get -u github.com/mitchellh/gox
	go get -u github.com/tcnksm/ghr
	go get -u github.com/sanbornm/go-selfupdate
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/Clever/gitsem
	dep ensure
	@echo 'Ready to build!'

clean:
	rm -fR ./.workspace

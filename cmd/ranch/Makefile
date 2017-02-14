.PHONY: all test

all: test

test:
	go test -v ./cmd ./util

vendor:
	godep save -r ./...


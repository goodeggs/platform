.PHONY: all test

all: test

test:
	docker info >/dev/null
	go get -t ./...
	sh test.sh

vendor:
	godep save -r ./...


.PHONY: all build release

GIT_SHA := $(shell git rev-parse --short HEAD)

all: build

build:
	docker build -t goodeggs/logspout-goodeggs .

release:
	( git diff --quiet && git diff --cached --quiet ) || ( echo "checkout must be clean"; false )
	docker build -t goodeggs/logspout-goodeggs:$(GIT_SHA) .
	docker push goodeggs/logspout-goodeggs:$(GIT_SHA)
	docker tag goodeggs/logspout-goodeggs:$(GIT_SHA) goodeggs/logspout-goodeggs:latest
	docker push goodeggs/logspout-goodeggs:latest


.PHONY: all test

all: test

test:
	docker info >/dev/null
	go get -t ./...
	sh test.sh

vendor:
	godep save -r ./...

version:
	if git status --porcelain | grep '^[^?]' >/dev/null; then echo "git working directory is dirty"; false; fi
	if [ -z "$(VERSION)" ]; then echo "must provide VERSION"; false; fi
	sed -i '/^const Version =/c\const Version = "$(VERSION)"' util/version.go
	git commit -a -m 'v$(VERSION)'
	git tag -am 'v$(VERSION)' v$(VERSION)
	go build
	tar czf ranch-osx-$(VERSION).tar.gz ranch
	shasum -a 256 ranch-osx-$(VERSION).tar.gz


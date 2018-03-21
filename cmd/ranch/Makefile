.PHONY: all test

SHELL=/bin/bash -o pipefail
export GOPATH=$(shell pwd)/.workspace:$(shell pwd)/vendor
export PATH := $(GOPATH)/bin:$(PATH)

export VERSION=$(shell cat ../../VERSION)
export MINORVER=$(shell echo '$(VERSION)'| awk -F. '{print $$1 "." $$2}')

all: build_deps test

test:
	@pushd .workspace/src/github.com/goodeggs/platform/cmd/ranch > /dev/null && \
		go test -v ./cmd ./util && \
		popd > /dev/null

build:
	go build -a -ldflags '-extldflags "-static"' .

release: build_deps test build_release release_github update_bindiffs generate_brew_formula

build_release: build_deps test
	gox -osarch "darwin/amd64 linux/amd64" -ldflags "-X main.VERSION=$${VERSION} -extldflags '-static'" -output "releases/$${VERSION}/{{.OS}}_{{.Arch}}/ranch"
	rm -rf "releases/$$VERSION/dist" && mkdir -p "releases/$$VERSION/dist"
	cp "releases/$$VERSION/darwin_amd64/ranch" "releases/$$VERSION/dist/ranch-Darwin-x86_64"
	cp "releases/$$VERSION/linux_amd64/ranch" "releases/$$VERSION/dist/ranch-Linux-x86_64"

install:
	go build -a -ldflags '-extldflags "-static"' -o .workspace/bin/ranch .

build_deps:
	@rm -fR ./vendor.orig
	@mkdir -p ./.workspace/src/github.com/goodeggs/platform/cmd
	@[ -s .workspace/src/github.com/goodeggs/platform/cmd/ranch ] || ln -s $$(pwd) .workspace/src/github.com/goodeggs/platform/cmd/ranch
	go get -u github.com/mitchellh/gox
	go get -u github.com/tcnksm/ghr
	go get -u github.com/sanbornm/go-selfupdate
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/Clever/gitsem
	go get -u gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install
	@pushd .workspace/src/github.com/goodeggs/platform/cmd/ranch > /dev/null && \
		dep ensure -v && \
		popd > /dev/null
	@echo 'Ready to build!'

release_github:
	# Github Release
	echo "releasing v$${VERSION} to github ..."
	ghr -t "$${GITHUB_TOKEN}" -u goodeggs -r platform --replace "v$${VERSION}" "releases/$${VERSION}/dist/"

update_bindiffs:
	# Releasing to S3
	echo "syncing ranch-updates S3 bucket"
	mkdir -p public
	( ls -d public/*/ | grep -v "^public/${MINORVER}." | xargs rm -rf ) || true
	aws s3 sync --exclude "*" --include "${MINORVER}.*" --include "darwin-amd64.json" --include "linux-amd64.json" s3://ranch-updates.goodeggs.com/stable/ranch/ public/

	# Generating bindiffs for binary
	echo "go-selfupdate generating bindiffs"
	mkdir releases/$${VERSION}/bins
	cp releases/$${VERSION}/darwin_amd64/ranch releases/$${VERSION}/bins/darwin-amd64
	cp releases/$${VERSION}/linux_amd64/ranch releases/$${VERSION}/bins/linux-amd64
	go-selfupdate releases/$${VERSION}/bins/ $${VERSION}

	# Sync bindiffs to s3
	echo "syncing ranch-updates S3 bucket"
	aws s3 sync --acl public-read public/ s3://ranch-updates.goodeggs.com/stable/ranch/

generate_brew_formula:
	@export SHA=$$(shasum -a 256 releases/$${VERSION}/dist/ranch-Darwin-x86_64 | awk '{print $$1}'); \
	/bin/echo -e "ranch v${VERSION} released. \\n\
NOTE: you must go update the homebrew formula manually.\\n\
  source:  https://github.com/goodeggs/homebrew-delivery-eng/tree/master/Formula/ranch.rb \\n\
  version: $${VERSION} \\n\
  sha256:  $${SHA}"

clean:
	rm -fR ./vendor
	rm -fR ./vendor.orig
	rm -fR ./.workspace
	rm -fR ./releases
	rm -fR ./public
	rm -f ./ranch


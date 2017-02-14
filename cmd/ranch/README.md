ranch
=====
A CLI interface to the Good Eggs Platform.

Installation
------------

You'll need to get a working ~/.ranch.yml from someone.

```
$ brew tap goodeggs/delivery-eng
$ brew install ranch
```

Development
-----------

```
$ brew install golang direnv
$ mkdir -p platform/src/github.com/goodeggs
$ cd platform
$ echo 'layout "go"' > .envrc
$ cd src/github.com/goodeggs
$ git clone https://github.com/goodeggs/platform.git
$ cd platform/cmd/ranch
$ make
```

We use [gvt](https://github.com/FiloSottile/gvt) for dependency management, so `gvt fetch` instead of `go get`

Releasing
---------

One time, you'll need to install [goxc](https://github.com/laher/goxc) and set your Github API Token so that we can automatically create a release:

```
$ go get github.com/laher/goxc
$ goxc -wlc default publish-github -apikey=123456789012
```

Then, to create a release:

```
$ sh version.sh {major,minor,patch}
$ sh release.sh
```

Don't forget to update the [ranch homebrew formula](https://github.com/goodeggs/homebrew-delivery-eng/blob/master/Formula/ranch.rb) with the new `version` and `sha256`.


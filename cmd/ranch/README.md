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

We use [dep](https://github.com/golang/dep) for dependency management.

Releasing
---------

To create a release:

```
$ go get github.com/Clever/gitsem
$ gitsem {major,minor,patch}
$ git push
$ GITHUB_TOKEN=xxx ./release.sh
```

Don't forget to update the [ranch homebrew formula](https://github.com/goodeggs/homebrew-delivery-eng/blob/master/Formula/ranch.rb) with the new `version` and `sha256`.


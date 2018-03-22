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
$ git clone https://github.com/goodeggs/platform.git
$ cd platform/cmd/ranch
$ make build_deps && make build
```

Testing
-----------

```
$ make test
```

We use [dep](https://github.com/golang/dep) for dependency management.

Releasing
---------

To create a release:

```
$ gitsem {major,minor,patch}
$ git push
$ GITHUB_TOKEN=xxx prod make release
```

Then update the [ranch homebrew formula](https://github.com/goodeggs/homebrew-delivery-eng/blob/master/Formula/ranch.rb) with the new `version` and `sha256`.

Then update the blessed version in [goodeggs/travis-utils/install-ci-tools.sh](https://github.com/goodeggs/travis-utils/blob/master/install-ci-tools.sh)

Then update the blessed version in [goodeggs/ecru/bin/onbuild](https://github.com/goodeggs/ecru/blob/master/bin/onbuild)

Then update the blessed version in [goodeggs/mk/bin/onbuild](https://github.com/goodeggs/mk/blob/master/bin/onbuild)


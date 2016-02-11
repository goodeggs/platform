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
$ mkdir -p platform/src/github.com/goodeggs
$ cd platform/src/github.com/goodeggs
$ git clone https://github.com/goodeggs/platform.git
$ cd cmd/ranch
$ go get ...
```

Releasing
---------

```
$ VERSION=x.x.x make version
$ git push
```

And then go update the [homebrew formula](https://github.com/goodeggs/homebrew-delivery-eng/blob/master/Formula/ranch.rb).


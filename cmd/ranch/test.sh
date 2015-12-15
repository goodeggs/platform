#!/bin/bash

function die {
  echo "$1"
  exit 1
}

go get .
ranch=`which ranch`

wd=`mktemp -d -t ranch`
cd $wd

mkdir foo && cd foo

( $ranch init ) || die "expect 'ranch init' to exit 0"

( cat .ranch.yaml | grep "^name: foo$" ) || die "expected .ranch.yaml to contain 'name: foo'"

[ $($ranch version) == "v1" ] || die "expected 'ranch version' to return v1"

( $ranch version bump ) && die "expected 'ranch version bump' to fail without a git repo"

git init .
git add .ranch.yaml
git commit -m 'initial commit'

( $ranch version bump ) || die "expected 'ranch version bump' to exit 0"
[ $($ranch version) == "v2" ] || die "expected 'ranch version' to return v2"

( git show HEAD | grep v2 ) || die "expected HEAD commit to be for v2"
( git tag | grep v2 ) || die "expected 'v2' git tag"

cd - ; rm -rf $wd
exit 0

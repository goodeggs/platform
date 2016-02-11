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

git init .
git add .ranch.yaml
git commit -m 'initial commit'

cd - ; rm -rf $wd
exit 0

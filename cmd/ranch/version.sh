#!/usr/bin/env bash

( git status --porcelain | grep -v '^??' > /dev/null ) && echo "git checkout is not clean" && exit 1

if [ "$1" == "major" ]; then
  dot=0
elif [ "$1" == "minor" ]; then
  dot=1
elif [ "$1" == "patch" ]; then
  dot=2
else
  echo "usage: $(basename $0) <major,minor,patch>"
  exit 1
fi

goxc bump -dot=${dot}
git add .goxc.json
version=v$(cat .goxc.json | jq -r '.PackageVersion')
git commit -a -m ${version}
git tag -am ${version} ${version}
echo ${version}

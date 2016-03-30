#!/usr/bin/env bash
set -e
set -o pipefail

indent() {
  sed -u 's/^/       /'
}

version=$(cat .goxc.json | jq -r '.PackageVersion')

echo "releasing v${version}..."

goxc 2>&1 | indent

sha=$(shasum -a 256 releases/${version}/ranch_${version}_darwin_amd64.zip)

cat <<-EOF

ranch v${version} released.

NOTE: you must go update the homebrew formula manually.
    source:  https://github.com/goodeggs/homebrew-delivery-eng/tree/master/Formula/ranch.rb
    version: ${version}
    sha256:  ${sha}

EOF

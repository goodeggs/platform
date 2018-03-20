#!/usr/bin/env bash
set -e
set -o pipefail

make build_deps
make test

exit




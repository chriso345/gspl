#!/usr/bin/env bash
set -euo pipefail

# Go to bindings/c directory
cd "$(dirname "$0")/../../bindings/c"

echo "Building Go shared library..."
go build -buildmode=c-shared -o libgspl.so .

echo "Done. libgspl.so and libgspl.h generated in $(pwd)"

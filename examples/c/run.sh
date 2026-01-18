#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINDINGS_DIR="$SCRIPT_DIR/../../bindings/c"
EXAMPLE_SRC="$SCRIPT_DIR/simple_lp.c"
OUTPUT_BIN="$SCRIPT_DIR/simple_lp.out"

echo "Building Go shared library..."
go build -buildmode=c-shared -o "$BINDINGS_DIR/libgspl.so" "$BINDINGS_DIR"

echo "Compiling C example..."
gcc "$EXAMPLE_SRC" -I"$BINDINGS_DIR" -L"$BINDINGS_DIR" -Wl,-rpath="$BINDINGS_DIR" -lgspl -o "$OUTPUT_BIN"

echo "Running example..."
"$OUTPUT_BIN"

#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINDINGS_DIR="$SCRIPT_DIR/../../bindings/c"
EXAMPLE_DIR="$SCRIPT_DIR"
OUTPUT_DIR="$SCRIPT_DIR"

mkdir -p "$OUTPUT_DIR"

echo "Building Go shared library..."
go build -buildmode=c-shared -o "$BINDINGS_DIR/libgspl.so" "$BINDINGS_DIR"

# Loop over all C files in the example directory
for EXAMPLE_SRC in "$EXAMPLE_DIR"/*.c; do
  # Check if any C files exist
  if [ ! -e "$EXAMPLE_SRC" ]; then
    echo "No C files found in $EXAMPLE_DIR"
    exit 1
  fi

  BASENAME=$(basename "$EXAMPLE_SRC" .c)
  OUTPUT_BIN="$OUTPUT_DIR/$BASENAME.out"

  echo "Compiling $EXAMPLE_SRC..."
  gcc "$EXAMPLE_SRC" -I"$BINDINGS_DIR" -L"$BINDINGS_DIR" -Wl,-rpath="$BINDINGS_DIR" -lgspl -o "$OUTPUT_BIN"

  echo "Running $OUTPUT_BIN..."
  "$OUTPUT_BIN"
done

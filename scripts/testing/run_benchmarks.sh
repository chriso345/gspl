#!/usr/bin/env bash
set -euo pipefail

CURRENT_REF="WORKTREE"
# BENCH_DIR="./benchmarks"
BENCH_DIR=$(mktemp -d -t go_benchmarks_XXXX)

PACKAGE="./..."
BASE_REF="HEAD"
TAG=""

COUNT=10
BENCHTIME="1ms"

# Print help
usage() {
  cat <<EOF
Usage: $0 [BRANCH_OR_COMMIT] [-p PACKAGE] [-t TAG]

Compare Go benchmarks between a base commit/tag and the current working tree.

Options:
  -p, --package PACKAGE   Specify package (defaults to ./...)
  -t, --tag TAG           Specify tag to benchmark against
  -h, --help              Show this help message
EOF
}

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
  -h | --help)
    usage
    exit 0
    ;;
  -p | --package)
    PACKAGE="$2"
    shift 2
    ;;
  -t | --tag)
    TAG="$2"
    shift 2
    ;;
  *)
    BASE_REF="$1"
    shift
    ;;
  esac
done

# Override BASE_REF with tag if provided
if [[ -n "$TAG" ]]; then
  BASE_REF="$TAG"
fi

mkdir -p "$BENCH_DIR"
OLD_BENCH_FILE="$BENCH_DIR/bench_${BASE_REF//\//_}.txt"
NEW_BENCH_FILE="$BENCH_DIR/bench_current.txt"

cleanup() {
  git worktree remove -f "${WORKTREE_DIR:-}" >/dev/null 2>&1 || true
}
trap cleanup EXIT

run_benchmarks() {
  local ref="$1"
  local file="$2"
  local pkg="$3"

  if [[ "$ref" == "$CURRENT_REF" || "$ref" == "HEAD" ]]; then
    echo "Running benchmarks on $ref..."
    go test -bench=. -benchmem -count="$COUNT" -benchtime="$BENCHTIME" "$pkg" >"$file"
  else
    WORKTREE_DIR=$(mktemp -d -t bench_ref_XXXX)
    echo "Creating temporary worktree for $ref..."
    git worktree add "$WORKTREE_DIR" "$ref" >/dev/null
    echo "Running benchmarks on $ref..."
    (cd "$WORKTREE_DIR" && go test -bench=. -benchmem -count="$COUNT" -benchtime="$BENCHTIME" "$pkg") >"$file"
  fi
}

# Run benchmarks
run_benchmarks "$BASE_REF" "$OLD_BENCH_FILE" "$PACKAGE"
run_benchmarks "$CURRENT_REF" "$NEW_BENCH_FILE" "$PACKAGE"

# Generate CSV using benchstat
CSV_FILE="$BENCH_DIR/bench.csv"
benchstat -format=csv "$OLD_BENCH_FILE" "$NEW_BENCH_FILE" >"$CSV_FILE"

# Print benchmark table using Python
echo
./scripts/testing/print_bench.py "$CSV_FILE" "$BASE_REF" "$CURRENT_REF"

#!/usr/bin/env bash
set -euo pipefail

# Deterministic benchmark runner
# Enhancements: repeated runs, CPU pinning (taskset), explicit GOMAXPROCS,
# longer benchtime, and benchstat aggregation to reduce noise.

CURRENT_REF="WORKTREE"
BENCH_DIR=$(mktemp -d -t go_benchmarks_XXXX)

PACKAGE="./..."
BASE_REF="HEAD"
TAG=""

# Number of repeated runs per ref (odd number recommended for median)
REPEATS=7
# benchtime to allow stable measurement (increase for larger problems)
BENCHTIME="2s"
# Use a single OS thread for reproducibility by default; can be overridden
GOMAXPROCS=${GOMAXPROCS:-1}

# Optional CPU pinning: set TASKSET to a CPU mask like "0x1" to pin to CPU0
TASKSET_CMD=""
if command -v taskset >/dev/null 2>&1; then
  # default: pin to CPU 0
  TASKSET_CMD="taskset 0x1"
fi

usage() {
  cat <<EOF
Usage: $0 [BRANCH_OR_COMMIT] [-p PACKAGE] [-t TAG]

Runs repeated, pinned Go benchmarks for HEAD/worktree and a base ref and
aggregates results using benchstat. Outputs CSV used by print_bench.py.

Options:
  -p, --package PACKAGE   Specify package (defaults to ./...)
  -t, --tag TAG           Specify tag to benchmark against
  -h, --help              Show this help message

Environment:
  GOMAXPROCS - number of OS threads to use (default: ${GOMAXPROCS})
  TASKSET    - if available, will pin process to CPU 0 by default
EOF
}

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

if [[ -n "$TAG" ]]; then
  BASE_REF="$TAG"
fi

mkdir -p "$BENCH_DIR"
OLD_DIR="$BENCH_DIR/old"
NEW_DIR="$BENCH_DIR/new"
mkdir -p "$OLD_DIR" "$NEW_DIR"

cleanup() {
  git worktree remove -f "${WORKTREE_DIR:-}" >/dev/null 2>&1 || true
}
trap cleanup EXIT

run_one() {
  local ref="$1"
  local outdir="$2"
  local pkg="$3"
  local iter="$4"

  echo "[${ref}] Run #${iter}"
  local outfile="$outdir/run_${iter}.txt"

  if [[ "$ref" == "$CURRENT_REF" || "$ref" == "HEAD" ]]; then
    env GOMAXPROCS="$GOMAXPROCS" GODEBUG="cgocheck=0" bash -c "${TASKSET_CMD} go test -bench=. -benchmem -run=^$ -benchtime=$BENCHTIME $pkg" >"$outfile" 2>&1
  else
    WORKTREE_DIR=$(mktemp -d -t bench_ref_XXXX)
    git worktree add "$WORKTREE_DIR" "$ref" >/dev/null
    (cd "$WORKTREE_DIR" && env GOMAXPROCS="$GOMAXPROCS" GODEBUG="cgocheck=0" bash -c "${TASKSET_CMD} go test -bench=. -benchmem -run=^$ -benchtime=$BENCHTIME $pkg") >"$outfile" 2>&1
  fi
}

run_repeated() {
  local ref="$1"
  local outdir="$2"
  local pkg="$3"
  for i in $(seq 1 $REPEATS); do
    run_one "$ref" "$outdir" "$pkg" "$i"
  done
  # Concatenate runs into a single file that benchstat can consume as multiple samples
  # benchstat expects two files, so we create one combined file per ref (newline separated)
  cat "$outdir"/run_*.txt >"$outdir/combined.txt"
}

# Run benchmarks for base and worktree
run_repeated "$BASE_REF" "$OLD_DIR" "$PACKAGE"
run_repeated "$CURRENT_REF" "$NEW_DIR" "$PACKAGE"

# Use benchstat to aggregate results; it accepts multiple samples arranged vertically
OLD_BENCH_FILE="$OLD_DIR/combined.txt"
NEW_BENCH_FILE="$NEW_DIR/combined.txt"
CSV_FILE="$BENCH_DIR/bench_${BASE_REF//\//_}_vs_worktree.csv"

benchstat -format=csv "$OLD_BENCH_FILE" "$NEW_BENCH_FILE" >"$CSV_FILE"

# Print benchmark table using Python
./scripts/testing/print_bench.py "$CSV_FILE" "$BASE_REF" "$CURRENT_REF"

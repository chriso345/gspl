#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

EXCLUDE_DIRS=(
  ".github"
  "bindings"
  "cmd"
  "examples"
  "internal"
  "scripts"
)
EXCLUDE_FILES=(
  "*_test.go"
)

# Build rg/glob exclude arguments
RG_EXCLUDES=()
for d in "${EXCLUDE_DIRS[@]}"; do
  RG_EXCLUDES+=(--glob "!$d/**")
done
for f in "${EXCLUDE_FILES[@]}"; do
  RG_EXCLUDES+=(--glob "!$f")
done

# Find exported top-level functions with file context
if command -v rg >/dev/null 2>&1; then
  exported=$(rg --no-heading --line-number \
    '^func\s+[A-Z][A-Za-z0-9_]*' \
    "${RG_EXCLUDES[@]}" . |
    sed -E 's/^(.+):[0-9]+:func\s+([A-Z][A-Za-z0-9_]*).*/\1:\2/' |
    sort -u)
else
  GREP_EXCLUDES=()
  for d in "${EXCLUDE_DIRS[@]}"; do
    GREP_EXCLUDES+=(--exclude-dir="$d")
  done
  for f in "${EXCLUDE_FILES[@]}"; do
    GREP_EXCLUDES+=(--exclude="$f")
  done

  exported=$(grep -R --line-number -h \
    -E "^func\s+[A-Z][A-Za-z0-9_]*" \
    "${GREP_EXCLUDES[@]}" . |
    sed -E 's/^(.+):[0-9]+:func\s+([A-Z][A-Za-z0-9_]*).*/\1:\2/' |
    sort -u)
fi

if [ -z "$exported" ]; then
  echo "No exported functions found."
  exit 0
fi

declare -A missing_by_file
missing_count=0

while IFS=: read -r file fn; do
  if ! rg -q "\b$fn\b" bindings/c/gspl.go bindings/c/gspl.h 2>/dev/null; then
    missing_by_file["$file"]+=$'\n'"  - $fn"
    missing_count=$((missing_count + 1))
  fi
done <<<"$exported"

if [ "$missing_count" -eq 0 ]; then
  echo "Bindings cover all exported functions (heuristic)"
  exit 0
fi

echo "Missing bindings detected:"
echo

for file in "${!missing_by_file[@]}"; do
  echo "$file"
  echo "${missing_by_file[$file]}"
  echo
done

echo "────────────────────────────"
echo "Total missing bindings: $missing_count"
echo "Files affected: ${#missing_by_file[@]}"
echo "────────────────────────────"

exit 2

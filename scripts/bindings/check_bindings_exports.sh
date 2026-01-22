#!/usr/bin/env bash
set -euo pipefail

# Move to the repository root reliably
ROOT_DIR="$(git rev-parse --show-toplevel)"
cd "$ROOT_DIR"

# Directories/files to exclude from scanning for exported functions
EXCLUDE_DIRS=(
  ".github"
  "cmd"
  "examples"
  "internal"
  "scripts"
)
EXCLUDE_FILES=(
  "*_test.go"
)

# Build ripgrep/glob exclude arguments
RG_EXCLUDES=()
for d in "${EXCLUDE_DIRS[@]}"; do
  RG_EXCLUDES+=(--glob "!$d/**")
done
for f in "${EXCLUDE_FILES[@]}"; do
  RG_EXCLUDES+=(--glob "!$f")
done

# Find exported top-level functions
GREP_EXCLUDES=()
for d in "${EXCLUDE_DIRS[@]}"; do
  GREP_EXCLUDES+=(--exclude-dir="$d")
done
for f in "${EXCLUDE_FILES[@]}"; do
  GREP_EXCLUDES+=(--exclude="$f")
done

exported=$(grep -R -n "${GREP_EXCLUDES[@]}" -E '^func[[:space:]]+[A-Z][A-Za-z0-9_]*' . |
  while IFS=: read -r file _ rest; do
    fn=$(printf '%s' "$rest" | sed -E 's/^func[[:space:]]+([A-Z][A-Za-z0-9_]*).*/\1/')
    printf '%s:%s\n' "$file" "$fn"
  done | sort -u)

if [ -z "$exported" ]; then
  echo "No exported functions found."
  exit 0
fi

# Absolute paths to bindings files
BINDINGS_GO="$ROOT_DIR/bindings/c/gspl.go"
BINDINGS_H="$ROOT_DIR/bindings/c/gspl.h"

declare -A missing_by_file
missing_count=0

while IFS=: read -r file fn; do
  if ! rg -q "\b$fn\b" "$BINDINGS_GO" "$BINDINGS_H" 2>/dev/null; then
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

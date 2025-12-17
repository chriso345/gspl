#!/usr/bin/env bash

set -euo pipefail

# Directories to exclude from tests
EXCLUDE_DIRS=(
  "examples/"
)

# Coverage threshold (percent)
COVERAGE_THRESHOLD=60

# Build a list of packages to test, excluding the specified directories
mapfile -t ALL_PACKAGES < <(go list ./...)
PACKAGES=()
for pkg in "${ALL_PACKAGES[@]}"; do
  skip=false
  for exclude in "${EXCLUDE_DIRS[@]}"; do
    if [[ "$pkg" == *"$exclude"* ]]; then
      skip=true
      break
    fi
  done
  $skip || PACKAGES+=("$pkg")
done

if [[ ${#PACKAGES[@]} -eq 0 ]]; then
  echo "No packages to test after exclusions."
  exit 0
fi

# Run tests with coverage
go test -coverprofile=coverage.out "${PACKAGES[@]}"

# Extract total coverage
COVERAGE=$(go tool cover -func=coverage.out | awk '/total:/ {print $3}' | sed 's/%//')

# Compare coverage against threshold
COVERAGE_INT=${COVERAGE%.*}
if ((COVERAGE_INT < COVERAGE_THRESHOLD)); then
  echo "Coverage $COVERAGE% is below threshold of $COVERAGE_THRESHOLD%"
  exit 1
else
  echo "Coverage $COVERAGE% meets threshold of $COVERAGE_THRESHOLD%"
fi

#!/usr/bin/env bash

set -euo pipefail

# Directories to exclude from tests
EXCLUDE_DIRS=(
  "examples/"
  "bindings/"
  "cmd/"
)

# Exclude the root module/package (true/false)
EXCLUDE_ROOT=true

# Coverage thresholds (percent)
DEFAULT_PACKAGE_THRESHOLD=50
TOTAL_THRESHOLD=60
CORE_PACKAGE_THRESHOLD=90

# Core packages (require higher coverage)
CORE_PACKAGES=(
  "github.com/chriso345/gspl/internal/simplex"
  "github.com/chriso345/gspl/solver"
  "github.com/chriso345/gspl/lp"
)

# ANSI colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

mapfile -t ALL_PACKAGES < <(go list ./...)

PACKAGES=()
for pkg in "${ALL_PACKAGES[@]}"; do
  skip=false

  # Exclude root package if requested
  if $EXCLUDE_ROOT && [[ "$pkg" == "$(go list .)" ]]; then
    skip=true
  fi

  # Exclude directories in EXCLUDE_DIRS
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

echo "Packages to test:"
printf ' - %s\n' "${PACKAGES[@]}"
echo

echo "Checking per-package coverage:"
echo "--------------------------------"

package_failed=false
for pkg in "${PACKAGES[@]}"; do
  output=$(go test -cover "$pkg")
  coverage=$(echo "$output" | sed -n 's/.*coverage: \([0-9.]*\)%.*/\1/p')

  # Skip packages with no statements
  if [[ -z "$coverage" ]]; then
    echo "$pkg: no statements to cover."
    continue
  fi

  coverage_int=${coverage%.*}

  # Decide threshold: core packages get higher threshold
  threshold=$DEFAULT_PACKAGE_THRESHOLD
  for core in "${CORE_PACKAGES[@]}"; do
    if [[ "$pkg" == "$core" ]]; then
      threshold=$CORE_PACKAGE_THRESHOLD
      break
    fi
  done

  if ((coverage_int < threshold)); then
    echo -e "${RED}(FAIL)${NC} $pkg: coverage ${coverage}% < threshold ${threshold}%"
    package_failed=true
  else
    echo -e "${GREEN}(PASS)${NC} $pkg: coverage ${coverage}%"
  fi
done

[[ $package_failed == false ]] && echo "All packages meet the threshold."
echo

echo "Checking total coverage:"
echo "------------------------"

# Run tests with coverage profile
go test -coverprofile=coverage.out "${PACKAGES[@]}" >/dev/null

# Extract total coverage
total_coverage=$(go tool cover -func=coverage.out | awk '/total:/ {print $3}' | sed 's/%//')
total_int=${total_coverage%.*}

if ((total_int < TOTAL_THRESHOLD)); then
  echo -e "${RED}(FAIL)${NC} TOTAL coverage ${total_coverage}% < threshold ${TOTAL_THRESHOLD}%"
  total_failed=true
else
  echo -e "${GREEN}(PASS)${NC} TOTAL coverage ${total_coverage}%"
  total_failed=false
fi

if [[ $package_failed == true || $total_failed == true ]]; then
  echo
  echo "Coverage checks failed."
  exit 1
else
  echo
  echo "All coverage checks passed."
  exit 0
fi

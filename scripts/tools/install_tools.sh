#!/usr/bin/env bash

set -euo pipefail

TOOLS=(
  "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.1"
  "golang.org/x/perf/cmd/benchstat@4febdb082b2ce2a0f9f6f496698aa5a2273279b8"
)

# Helper to get binary name from Go module
get_bin_name() {
  local module="$1"
  # Extract the last path component as the binary name
  echo "${module##*/}" | cut -d'@' -f1
}

echo ""

for tool in "${TOOLS[@]}"; do
  bin_name=$(get_bin_name "$tool")

  if command -v "$bin_name" >/dev/null 2>&1; then
    echo "$bin_name already installed."
  else
    echo "Installing $bin_name..."
    go install "$tool"
  fi
done

echo -e "\nTool installation complete."

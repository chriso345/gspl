scripts_dir := "./scripts"

# List tasks available
default:
    @just --list --list-prefix " - "

# Run unit tests
test:
    go test ./... -count=1

# Run benchmarks
bench folder="./benchmarks" branch="HEAD":
    {{ scripts_dir }}/testing/run_benchmarks.sh {{ folder }} {{ branch }}

# Coverage report
cover:
    {{ scripts_dir }}/testing/run_coverage.sh

# Run race detector
race:
    go test ./... -count=1 -race

# Lint (if golangci-lint is installed)
lint:
    golangci-lint run || true

# Install development tools
install-tools:
    {{ scripts_dir }}/tools/install_tools.sh

# Run the command line application
run *args:
    go run ./cmd/gspl {{ args }}

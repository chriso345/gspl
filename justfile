scripts_dir := "./scripts"

# List tasks available
default:
    @just --list --list-prefix " - "

# Run unit tests
test:
    go test ./... -count=1

# Run benchmarks for solver package
bench:
    go test ./solver -bench . -benchmem

# Coverage report
cover:
    {{ scripts_dir }}/coverage/run_coverage.sh

# Lint (if golangci-lint is installed)
lint:
    golangci-lint run || true

# Install development tools
install-tools:
    {{ scripts_dir }}/tools/install_tools.sh

# Run the command line application
run *args:
    go run ./cmd/gspl {{ args }}

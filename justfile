# =========================================
# Variables
# =========================================

scripts_dir := "./scripts"

# =========================================
# DEFAULT
# =========================================

# List tasks available
[group("Default")]
default:
    @just --list --list-prefix " - "

# =========================================
# TESTING
# =========================================

# Run all unit tests
[group("Unit Tests & Benchmarks")]
test:
    go test ./... -count=1

# Run benchmarks for a folder and branch
[group("Unit Tests & Benchmarks")]
bench folder="." branch="HEAD":
    {{ scripts_dir }}/testing/run_benchmarks.sh {{ folder }} {{ branch }}

# Run benchmarks for a single package
[group("Unit Tests & Benchmarks")]
bench-pkg pkg="./benchmarks/solver":
    go test {{ pkg }} -run=^$ -bench=. -benchtime=2s -benchmem

# Generate coverage report
[group("Unit Tests & Benchmarks")]
cover:
    {{ scripts_dir }}/testing/run_coverage.sh

# Run race detector
[group("Unit Tests & Benchmarks")]
race:
    go test ./... -count=1 -race

# =========================================
# LINTING & FORMATTING
# =========================================

# Run golangci-lint
[group("Linting & Formatting")]
lint:
    golangci-lint run || true

# =========================================
# DEVELOPMENT TOOLS
# =========================================

# Install development tools
[group("Development Tools")]
install-tools:
    {{ scripts_dir }}/tools/install_tools.sh

# Generate Go/C bindings
[group("Development Tools")]
bind-gen:
    {{ scripts_dir }}/bindings/generate_bindings.sh
    {{ scripts_dir }}/bindings/check_bindings_exports.sh

# =========================================
# APPLICATION
# =========================================

# Run the command line application
[group("Run Application")]
run *args:
    go run ./cmd/gspl {{ args }}

# =========================================
# DOCUMENTATION
# =========================================

# Open the documentation in the browser
[group("Documentation")]
docs:
    pkgsite

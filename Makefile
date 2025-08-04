VERSION ?= $(shell cat version/VERSION | tr -d '[:space:]')
GOFMT_FILES ?= $$(find . -name '*.go')
PKG_NAME = vastdata
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

EXT :=
ifeq ($(GOOS),windows)
EXT := .exe
endif

BINARY_NAME := terraform-provider-$(PKG_NAME)
BUILD_DIR := build/$(GOOS)_$(GOARCH)
BINARY := $(BUILD_DIR)/$(BINARY_NAME)$(EXT)
GOFLAGS := -mod=readonly
LDFLAGS := -X main.version=$(VERSION)

.PHONY: build show-resources show-datasources test test-unit test-benchmarks test-coverage test-all

build:
	@echo "Building $(BINARY_NAME) for $(GOOS)_$(GOARCH)..."
	mkdir -p $(BUILD_DIR)
	env GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY)
	@echo "Binary created at $(BINARY)"
	@chmod +x $(BINARY)


vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -s -w $(GOFMT_FILES)


ifeq (show, $(firstword $(MAKECMDGOALS)))
  runargs := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(foreach arg,$(runargs),$(eval $(arg):;@true))
endif

#   Run the Go schema visualization tool with dynamic arguments.
#   Usage:
#     make show <type> [filter] [automation]
#
#   Arguments:
#     <type>       - Required. The type of schema to show: 'resource', 'datasource', or short forms ('r', 'd').
#     [filter]     - Optional. Case-insensitive substring to filter by.
#     [automation] - Optional. If present (any value), enables automation mode (🔸 markers).
#
#   Examples:
#     make show resource
#     make show datasource user
#     make show r quota automation
show:
	@type=$(word 1, $(runargs)); \
	arg2=$(word 2, $(runargs)); \
	arg3=$(word 3, $(runargs)); \
	filter=""; \
	auto_flag=""; \
	[ "$$arg2" = "automation" ] && auto_flag="-automation" || filter="$$arg2"; \
	[ "$$arg3" = "automation" ] && auto_flag="-automation"; \
	echo "Running show with type=$$type, filter=$$filter, automation=$$auto_flag"; \
	go run $(CURDIR)/misc/show_schemas.go -type=$$type $${filter:+-filter=$$filter} $$auto_flag

# Test targets

test:
	@echo "Running unit tests..."
	go test -v -cover ./vastdata/provider/... ./vastdata/internalstate/... ./vastdata/schema_generation/... ./vastdata/client/...
	@echo "Running error handling and validation tests..."
	go test -v -cover ./vastdata/ -run '^Test(ErrorHandling_|Validation_|Normalize|ConvertMapKeys|KeyTransform|ValidateOneOf|ValidateAllOf|ValidateNoneOf|TFState_)'

# Run unit tests with verbose output
test-unit:
	@echo "Running unit tests..."
	go test -v -race -timeout=30s ./vastdata/provider/... ./vastdata/internalstate/... ./vastdata/schema_generation/...

# Note: Integration tests removed - they required real network connections

# Note: Resource lifecycle and client integration tests removed - they required real network connections

# Run performance benchmarks
test-benchmarks:
	@echo "Running performance benchmarks..."
	go test -v -bench=. -benchmem -timeout=120s ./vastdata/schema_generation/

# Run benchmarks and save results for comparison
test-benchmarks-save:
	@echo "Running benchmarks and saving results..."
	mkdir -p benchmarks
	go test -bench=. -benchmem -timeout=120s ./vastdata/schema_generation/ | tee benchmarks/benchmark_$(shell date +%Y%m%d_%H%M%S).txt

# Compare current benchmarks with previous results
test-benchmarks-compare:
	@echo "Comparing benchmarks..."
	@if [ -f benchmarks/baseline.txt ]; then \
		go test -bench=. -benchmem ./vastdata/schema_generation/ > benchmarks/current.txt; \
		echo "=== Benchmark Comparison ==="; \
		echo "Baseline vs Current:"; \
		diff -u benchmarks/baseline.txt benchmarks/current.txt || true; \
	else \
		echo "No baseline benchmark found. Run 'make test-benchmarks-baseline' first."; \
	fi

# Set current benchmarks as baseline
test-benchmarks-baseline:
	@echo "Setting benchmark baseline..."
	mkdir -p benchmarks
	go test -bench=. -benchmem ./vastdata/schema_generation/ > benchmarks/baseline.txt
	@echo "Baseline set. Use 'make test-benchmarks-compare' to compare future runs."

# Run tests with coverage reporting
test-coverage:
	@echo "Running tests with coverage..."
	mkdir -p coverage
	go test -v -race -coverprofile=coverage/coverage.out -covermode=atomic ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	go tool cover -func=coverage/coverage.out | grep "total:" | awk '{print "Total coverage: " $$3}'
	@echo "Coverage report saved to coverage/coverage.html"

# Run coverage and open in browser (macOS/Linux)
test-coverage-open: test-coverage
	@if command -v open >/dev/null 2>&1; then \
		open coverage/coverage.html; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open coverage/coverage.html; \
	else \
		echo "Coverage report available at coverage/coverage.html"; \
	fi

# Run specific test patterns
test-pattern:
	@if [ -z "$(PATTERN)" ]; then \
		echo "Usage: make test-pattern PATTERN=<test_pattern>"; \
		echo "Example: make test-pattern PATTERN=TestProvider"; \
		exit 1; \
	fi
	go test -v -run="$(PATTERN)" ./...

# Run tests for a specific package
test-pkg:
	@if [ -z "$(PKG)" ]; then \
		echo "Usage: make test-pkg PKG=<package_path>"; \
		echo "Example: make test-pkg PKG=./vastdata/provider"; \
		exit 1; \
	fi
	go test -v -race $(PKG)

# Run all tests (unit, benchmarks, coverage)
test-all: vet fmt test-unit test-benchmarks test-coverage
	@echo "All tests completed successfully!"

# Test with different Go versions (requires Docker)
test-go-versions:
	@echo "Testing with multiple Go versions..."
	@for version in 1.21 1.22 1.23; do \
		echo "Testing with Go $$version..."; \
		docker run --rm -v $(PWD):/workspace -w /workspace golang:$$version go test ./...; \
	done

# Clean test artifacts
test-clean:
	@echo "Cleaning test artifacts..."
	rm -rf coverage/ benchmarks/ *.test
	go clean -testcache

# Run tests with verbose output and save logs
test-verbose:
	@echo "Running verbose tests..."
	mkdir -p logs
	go test -v -race ./... 2>&1 | tee logs/test_$(shell date +%Y%m%d_%H%M%S).log

# Check for test flakiness by running tests multiple times
test-flakiness:
	@echo "Checking for test flakiness (running tests 10 times)..."
	@for i in $$(seq 1 10); do \
		echo "Run $$i/10..."; \
		go test -race ./... > /dev/null 2>&1 || { echo "Test failed on run $$i"; exit 1; }; \
	done
	@echo "No flaky tests detected!"

# Run only fast tests (exclude slow benchmarks)
test-fast:
	@echo "Running fast tests..."
	go test -v -race -short -timeout=30s ./...

# Test with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -v -race -timeout=60s ./...

# Run mutation tests (requires go-mutesting)
test-mutation:
	@echo "Running mutation tests..."
	@if ! command -v go-mutesting >/dev/null 2>&1; then \
		echo "Installing go-mutesting..."; \
		go install github.com/zimmski/go-mutesting/cmd/go-mutesting@latest; \
	fi
	go-mutesting ./...

# Performance profiling
test-profile-cpu:
	@echo "Running CPU profiling..."
	mkdir -p profiles
	go test -cpuprofile=profiles/cpu.prof -bench=. ./vastdata/schema_generation/
	@echo "CPU profile saved to profiles/cpu.prof"
	@echo "View with: go tool pprof profiles/cpu.prof"

test-profile-mem:
	@echo "Running memory profiling..."
	mkdir -p profiles
	go test -memprofile=profiles/mem.prof -bench=. ./vastdata/schema_generation/
	@echo "Memory profile saved to profiles/mem.prof"
	@echo "View with: go tool pprof profiles/mem.prof"

# Generate docs and copywrite headers
generate-docs:
	@echo "Generating documentation and headers..."
	@go run $(CURDIR)/misc/generate_examples.go
	cd tools; go generate ./...; cd -


ifeq (gen-openapi-tar, $(firstword $(MAKECMDGOALS)))
  runargs := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(foreach arg,$(runargs),$(eval $(arg):;@true))
endif

ifeq (validate-api, $(firstword $(MAKECMDGOALS)))
  runargs := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(foreach arg,$(runargs),$(eval $(arg):;@true))
endif

# Enhanced OpenAPI conversion with validation and auto-fixes
#
#   This target uses the enhanced Python converter that:
#     - Validates Swagger schemas and identifies common issues
#     - Automatically fixes known problems (null properties, missing types)  
#     - Provides detailed debugging information
#     - Converts to OpenAPI v3 using the Go converter
#     - Creates final outputs: api.json and api.tar.gz
#
#   Usage:
#     make gen-openapi-tar <path> [options]
#
#   Arguments:
#     <path> - Required. Path to the Swagger/OpenAPI YAML file.
#
#   Options:
#     --debug        - Enable detailed debugging output
#     --no-auto-fix  - Disable automatic fixes (validation only)
#     --output-dir   - Custom output directory for intermediate files
#
#   Examples:
#     make gen-openapi-tar ./specs/swagger.yaml
#     make gen-openapi-tar /tmp/vast-openapi.yaml --debug
#     make gen-openapi-tar ./specs/swagger.yaml --no-auto-fix
gen-openapi-tar:
	@set -e; \
	args="$(runargs)"; \
	if [ -z "$$args" ]; then \
		echo "❌ Usage: make gen-openapi-tar <path> [options]"; \
		echo "   Example: make gen-openapi-tar /path/to/swagger.yaml --debug"; \
		echo "   Options: --debug --no-auto-fix --output-dir /custom/dir"; \
		exit 1; \
	fi; \
	echo "🚀 Running enhanced OpenAPI conversion with validation and auto-fixes..."; \
	python3 $(CURDIR)/misc/convert_swagger_with_fixes.py $$args --dest-dir $(CURDIR); \
	echo "✅ Enhanced conversion completed!"

# Validate Swagger/OpenAPI schema
#
# Usage: make validate-api <path> [options]
#
# Arguments:
#   <path>         - Path to Swagger/OpenAPI YAML or JSON file
#
# Options:
#   --debug        - Enable detailed debugging output
#   --json-output  - Export JSON report to specified file
#
# Examples:
#   make validate-api ./specs/swagger.yaml
#   make validate-api /tmp/api.yaml --debug
#   make validate-api ./api.yaml --json-output report.json
validate-api:
	@set -e; \
	args="$(runargs)"; \
	if [ -z "$$args" ]; then \
		echo "❌ Usage: make validate-api <path> [options]"; \
		echo "   Example: make validate-api /path/to/swagger.yaml --debug"; \
		echo "   Options: --debug --json-output <file>"; \
		exit 1; \
	fi; \
	echo "🔍 Running Swagger/OpenAPI schema validation..."; \
	python3 $(CURDIR)/misc/validate_swagger.py $$args; \
	echo "✅ Validation completed!"


# Help target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build                    - Build the provider binary"
	@echo "  vet                      - Run go vet"
	@echo "  fmt                      - Format Go code"
	@echo ""
	@echo "Test targets:"
	@echo "  test                     - Run basic tests with coverage"
	@echo "  test-unit               - Run unit tests only"

	@echo "  test-benchmarks         - Run performance benchmarks"
	@echo "  test-benchmarks-save    - Run benchmarks and save results"
	@echo "  test-benchmarks-compare - Compare current benchmarks with baseline"
	@echo "  test-benchmarks-baseline - Set current benchmarks as baseline"
	@echo "  test-coverage           - Run tests with coverage reporting"
	@echo "  test-coverage-open      - Run coverage and open in browser"
	@echo "  test-all                - Run all tests (unit, benchmarks, coverage)"
	@echo "  test-fast               - Run only fast tests"
	@echo "  test-race               - Run tests with race detection"
	@echo "  test-verbose            - Run tests with verbose output and save logs"
	@echo "  test-flakiness          - Check for flaky tests"
	@echo "  test-clean              - Clean test artifacts"
	@echo ""
	@echo "Profiling targets:"
	@echo "  test-profile-cpu        - Run CPU profiling"
	@echo "  test-profile-mem        - Run memory profiling"
	@echo ""
	@echo "Utility targets:"
	@echo "  show <type> [filter]    - Show resource/datasource schemas"
	@echo "  generate-docs           - Generate documentation"
	@echo "  gen-openapi-tar <path> [options] - Convert OpenAPI YAML to tarball with validation & auto-fixes"
	@echo "  validate-api <path> [options]    - Validate Swagger/OpenAPI schema with detailed diagnostics"
	@echo "  help                    - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make test-pattern PATTERN=TestProvider"
	@echo "  make test-pkg PKG=./vastdata/provider"
	@echo "  make show resource user"

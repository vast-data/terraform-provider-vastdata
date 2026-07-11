#!/bin/bash
# Copyright (c) HashiCorp, Inc.

# Simple script to run Terraform provider E2E tests

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Terraform Provider E2E Tests ===${NC}\n"

# Check environment variables
if [ -z "$VASTDATA_HOST" ]; then
    echo -e "${RED}ERROR: VASTDATA_HOST environment variable is required${NC}"
    echo "Usage: export VASTDATA_HOST=<cluster-ip> && ./run_tests.sh"
    exit 1
fi

echo -e "${YELLOW}Configuration:${NC}"
echo "  VASTDATA_HOST: $VASTDATA_HOST"
echo "  VASTDATA_USERNAME: ${VASTDATA_USERNAME:-admin}"
echo "  VASTDATA_PORT: ${VASTDATA_PORT:-443}"
echo ""

# Check if provider is built (match local `make build` platform)
GOOS="$(go env GOOS)"
GOARCH="$(go env GOARCH)"
PROVIDER_BINARY="../build/${GOOS}_${GOARCH}/terraform-provider-vastdata"
if [ ! -f "$PROVIDER_BINARY" ]; then
    echo -e "${RED}ERROR: Provider binary not found at $PROVIDER_BINARY${NC}"
    echo "Please build the provider first:"
    echo "  cd .. && make build"
    exit 1
fi

echo -e "${GREEN}✓ Provider binary found${NC}\n"

# Install dependencies if needed
if ! python3 -c "import pytest" 2>/dev/null; then
    echo -e "${YELLOW}Installing test dependencies...${NC}"
    pip install -r requirements-test.txt
fi

# Default: stop at first failure. Pass --continue to run all tests.
KEEP_GOING=false
FILTERED_ARGS=()
for arg in "$@"; do
    case "$arg" in
        --continue|--keep-going)
            KEEP_GOING=true
            ;;
        *)
            FILTERED_ARGS+=("$arg")
            ;;
    esac
done
set -- "${FILTERED_ARGS[@]}"

PYTEST_FAILFAST=(-x)
if [ "$KEEP_GOING" = true ]; then
    PYTEST_FAILFAST=()
    echo -e "${YELLOW}Continue mode: running all tests even after failures${NC}\n"
else
    echo -e "${YELLOW}Fail-fast mode: stopping at first failure (use --continue to run all)${NC}\n"
fi

# Parse command line arguments
TEST_TYPE="${1:-all}"

case "$TEST_TYPE" in
    "all")
        shift || true
        echo -e "${GREEN}Running all tests...${NC}\n"
        pytest . -v "${PYTEST_FAILFAST[@]}" "$@"
        ;;
    "e2e")
        shift || true
        echo -e "${GREEN}Running e2e resource tests (parametrized)...${NC}\n"
        pytest test_e2e_resources.py::test_individual_resource -v "${PYTEST_FAILFAST[@]}" "$@"
        ;;
    "auxiliary")
        shift || true
        echo -e "${GREEN}Running auxiliary tests...${NC}\n"
        pytest test_auxiliary.py -v "${PYTEST_FAILFAST[@]}" "$@"
        ;;
    "failfast")
        shift || true
        echo -e "${GREEN}Running tests (stop at first failure)...${NC}\n"
        pytest . -x -v "$@"
        ;;
    "resource")
        shift || true
        RESOURCE_NAME="$1"
        if [ -z "$RESOURCE_NAME" ]; then
            echo -e "${RED}ERROR: Please specify resource name${NC}"
            echo "Usage: ./run_tests.sh resource vastdata_view/1.tf [pytest-args]"
            exit 1
        fi
        shift || true
        echo -e "${GREEN}Running test for $RESOURCE_NAME...${NC}\n"
        pytest "test_e2e_resources.py::test_individual_resource[$RESOURCE_NAME]" -v "${PYTEST_FAILFAST[@]}" "$@"
        ;;
    *)
        echo -e "${RED}ERROR: Unknown test type: $TEST_TYPE${NC}"
        echo ""
        echo "Usage: ./run_tests.sh [--continue] [all|e2e|auxiliary|failfast|resource <name>] [pytest-args]"
        echo ""
        echo "By default tests stop at the first failure. Use --continue to run the full suite."
        echo ""
        echo "Examples:"
        echo "  ./run_tests.sh e2e                    # E2E tests, stop on first failure"
        echo "  ./run_tests.sh --continue e2e         # Run all E2E tests"
        echo "  ./run_tests.sh all -s                 # All tests, fail-fast, show output"
        echo "  ./run_tests.sh resource vastdata_view/1.tf"
        echo "  ./run_tests.sh --continue resource vastdata_view/1.tf -vv"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}=== Tests completed ===${NC}"

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

# Check if provider is built
PROVIDER_BINARY="../build/linux_amd64/terraform-provider-vastdata"
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

# Parse command line arguments
TEST_TYPE="${1:-all}"

case "$TEST_TYPE" in
    "all")
        shift || true  # Remove test type, keep the rest for pytest
        echo -e "${GREEN}Running all tests (parametrized)...${NC}\n"
        pytest test_e2e_resources.py::test_individual_resource -v "$@"
        ;;
    "single")
        shift || true  # Remove test type, keep the rest for pytest
        echo -e "${GREEN}Running single test (all resources)...${NC}\n"
        pytest test_e2e_resources.py::test_terraform_provider_e2e -v "$@"
        ;;
    "parallel")
        shift || true  # Remove test type, keep the rest for pytest
        echo -e "${GREEN}Running tests in parallel...${NC}\n"
        pytest test_e2e_resources.py::test_individual_resource -n auto -v "$@"
        ;;
    "failfast")
        shift || true  # Remove test type, keep the rest for pytest
        echo -e "${GREEN}Running tests (stop at first failure)...${NC}\n"
        pytest test_e2e_resources.py::test_individual_resource -x -v "$@"
        ;;
    "resource")
        shift || true  # Remove test type
        RESOURCE_NAME="$1"
        if [ -z "$RESOURCE_NAME" ]; then
            echo -e "${RED}ERROR: Please specify resource name${NC}"
            echo "Usage: ./run_tests.sh resource vastdata_view/1.tf [pytest-args]"
            exit 1
        fi
        shift || true  # Remove resource name, keep extra args
        echo -e "${GREEN}Running test for $RESOURCE_NAME...${NC}\n"
        pytest "test_e2e_resources.py::test_individual_resource[$RESOURCE_NAME]" -v "$@"
        ;;
    *)
        echo -e "${RED}ERROR: Unknown test type: $TEST_TYPE${NC}"
        echo ""
        echo "Usage: ./run_tests.sh [all|single|parallel|failfast|resource <name>] [pytest-args]"
        echo ""
        echo "Examples:"
        echo "  ./run_tests.sh all              # Run all tests (one per resource)"
        echo "  ./run_tests.sh single           # Run single test (all resources)"
        echo "  ./run_tests.sh parallel         # Run tests in parallel"
        echo "  ./run_tests.sh failfast         # Stop at first failure"
        echo "  ./run_tests.sh failfast -vv     # Stop at first failure, very verbose"
        echo "  ./run_tests.sh all -s           # Run all, show output"
        echo "  ./run_tests.sh resource vastdata_view/1.tf"
        echo "  ./run_tests.sh resource vastdata_view/1.tf -vv -s"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}=== Tests completed ===${NC}"


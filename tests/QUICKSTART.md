# Quick Start Guide

Get started with Terraform provider E2E tests in 5 minutes.

## Step 1: Build the Provider

```bash
make build
```

Verify the binary exists:
```bash
ls -lh build/linux_amd64/terraform-provider-vastdata
```

## Step 2: Set Environment Variables

```bash
export VASTDATA_HOST="10.0.0.100"        # Your VAST cluster VMS IP
export VASTDATA_USERNAME="admin"         # Optional, defaults to 'admin'
export VASTDATA_PASSWORD="123456"        # Optional, defaults to '123456'
```

## Step 3: Install Test Dependencies

```bash
cd tests
pip install -r requirements-test.txt
```

## Step 4: Run Tests

### Option A: Use the script

```bash
./run_tests.sh all
```

### Option B: Use pytest directly

```bash
pytest test_e2e_resources.py::test_individual_resource -v
```

### Option C: Run a single resource test

```bash
./run_tests.sh resource vastdata_view/1.tf
```

## Common Commands

```bash
# Run all tests (one per resource)
./run_tests.sh all

# Run all resources in one test (faster, but harder to debug)
./run_tests.sh single

# Run tests in parallel (fastest)
./run_tests.sh parallel

# Stop at first failure (useful for debugging)
./run_tests.sh failfast

# Run specific resource
./run_tests.sh resource vastdata_view/1.tf

# Direct pytest - stop at first failure
pytest test_e2e_resources.py::test_individual_resource -x -v

# Run with more verbose output
pytest test_e2e_resources.py::test_individual_resource -vv

# Run only tests matching pattern
pytest -k "vastdata_view" -v

# Show test collection without running
pytest --collect-only
```

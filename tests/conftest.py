# Copyright (c) HashiCorp, Inc.

"""
Pytest configuration and fixtures for Terraform provider E2E tests.
"""

import os
import tempfile
import zipfile
import shutil
from pathlib import Path
from plumbum import local
import pytest


# Environment variables for connection
VASTDATA_HOST = os.environ.get('VASTDATA_HOST')
VASTDATA_USERNAME = os.environ.get('VASTDATA_USERNAME', 'admin')
VASTDATA_PASSWORD = os.environ.get('VASTDATA_PASSWORD', '123456')
VASTDATA_PORT = os.environ.get('VASTDATA_PORT', '443')

# Path to provider binary (built locally)
PROVIDER_BINARY_PATH = Path(__file__).parent.parent / "build" / "linux_amd64" / "terraform-provider-vastdata"


@pytest.fixture(scope="session", autouse=True)
def validate_environment():
    """Validate required environment variables are set."""
    if not VASTDATA_HOST:
        pytest.exit("ERROR: VASTDATA_HOST environment variable is required", returncode=1)
    
    if not PROVIDER_BINARY_PATH.exists():
        pytest.exit(
            f"ERROR: Provider binary not found at {PROVIDER_BINARY_PATH}\n"
            "Please build the provider first: make build",
            returncode=1
        )
    
    print(f"\n=== Test Configuration ===")
    print(f"VASTDATA_HOST: {VASTDATA_HOST}")
    print(f"VASTDATA_USERNAME: {VASTDATA_USERNAME}")
    print(f"VASTDATA_PORT: {VASTDATA_PORT}")
    print(f"Provider binary: {PROVIDER_BINARY_PATH}")
    print(f"=========================\n")


@pytest.fixture(scope="session")
def terraform_version():
    """Terraform version to use."""
    return "1.9.5"


@pytest.fixture(scope="session", autouse=True)
def install_terraform(terraform_version):
    """Ensure Terraform is available."""
    # Check if terraform is in PATH
    try:
        result = local.cmd.which('terraform')
        print(f"Using terraform from PATH: {result}")
        return
    except:
        pass
    
    # Check common locations
    for path in ['/usr/local/bin/terraform', '/usr/bin/terraform', '/opt/terraform/terraform']:
        if local.path(path).exists():
            print(f"Found terraform at {path}")
            return
    
    # Install terraform if not found
    print(f"Terraform not found, installing version {terraform_version}...")
    artifact = f'terraform_{terraform_version}_linux_amd64.zip'
    
    with tempfile.TemporaryDirectory() as tmpdir:
        tmpdir = local.path(tmpdir)
        with local.cwd(tmpdir):
            # Download terraform
            local.cmd.wget[f'https://releases.hashicorp.com/terraform/{terraform_version}/{artifact}']()
            
            # Extract to /usr/local/bin
            local.cmd.sudo['mkdir', '-p', '/usr/local/bin']()
            local.cmd.sudo['unzip', '-od', '/tmp', artifact]()
            local.cmd.sudo['mv', '/tmp/terraform', '/usr/local/bin/terraform']()
            local.cmd.sudo['chmod', '+x', '/usr/local/bin/terraform']()
    
    print(f"Terraform {terraform_version} installed successfully")


@pytest.fixture(scope="session", autouse=True)
def install_provider(install_terraform):
    """Install the locally built provider binary and configure terraform to use it."""
    # Just use the provider binary from build directory directly
    # No need to copy it - terraform will use it via TF_CLI_CONFIG_FILE
    
    if not PROVIDER_BINARY_PATH.exists():
        pytest.exit(
            f"ERROR: Provider binary not found at {PROVIDER_BINARY_PATH}\n"
            "Please build the provider first: make build",
            returncode=1
        )
    
    print(f"Using provider binary from: {PROVIDER_BINARY_PATH}")




@pytest.fixture
def terraform_workdir():
    """Create a temporary working directory for Terraform operations."""
    with tempfile.TemporaryDirectory() as tmpdir:
        workdir = local.path(tmpdir)
        
        # Create provider configuration with explicit required_providers
        provider_config = f"""\
terraform {{
  required_providers {{
    vastdata = {{
      source = "vastdata/vastdata"
    }}
  }}
}}

provider vastdata {{
  username = "{VASTDATA_USERNAME}"
  port = "{VASTDATA_PORT}"
  password = "{VASTDATA_PASSWORD}"
  host = "{VASTDATA_HOST}"
  skip_ssl_verify = true
  version_validation_mode = "warn"
}}
"""
        
        provider_file = workdir / "provider.tf"
        provider_file.write(provider_config)
        
        # Create an empty lock file - dev overrides don't need entries in the lock file
        # An empty lock file tells Terraform "yes, we know about dependencies" without specifying them
        lock_file = workdir / ".terraform.lock.hcl"
        lock_file.write("# This file is maintained automatically by \"terraform init\".\n# Manual edits may be lost in future updates.\n")
        
        yield workdir


@pytest.fixture
def terraform_config_file(tmp_path_factory):
    """Create a temporary terraform config file for dev overrides."""
    # Create temporary .terraformrc file
    config_dir = tmp_path_factory.mktemp("terraform_config")
    config_file = config_dir / "terraformrc"
    
    # Point to the build directory where the binary is
    provider_dir = PROVIDER_BINARY_PATH.parent
    
    config = f"""\
provider_installation {{
  dev_overrides {{
    "vastdata/vastdata" = "{provider_dir}"
  }}
  direct {{}}
}}
"""
    
    config_file.write_text(config)
    print(f"Created terraform config at {config_file}")
    print(f"Provider directory: {provider_dir}")
    
    return str(config_file)


@pytest.fixture
def terraform_cmd(terraform_workdir, terraform_config_file):
    """Get terraform command with proper working directory and config file."""
    # Use TF_CLI_CONFIG_FILE to point to our temporary config (not ~/.terraformrc)
    cmd = local['terraform'].with_cwd(terraform_workdir)
    cmd = cmd.with_env(TF_CLI_CONFIG_FILE=terraform_config_file)
    return cmd


@pytest.fixture(scope="session")
def tf_examples_dir():
    """Path to examples directory."""
    return Path(__file__).parent.parent / "examples" / "resources"


def pytest_collection_modifyitems(config, items):
    """Add markers to tests based on their path."""
    for item in items:
        if "test_e2e" in item.nodeid:
            item.add_marker(pytest.mark.e2e)


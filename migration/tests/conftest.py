# Copyright (c) HashiCorp, Inc.

"""
Pytest configuration and fixtures for migration tool tests.
"""

import pytest
import tempfile
import shutil
from pathlib import Path
import sys
import os

# Add the migration script directory to the Python path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import (
    transform_file,
    transform_resource_block,
    get_group_for_key,
    parse_nested_block,
    resource_type_rename_map,
    key_groups
)


@pytest.fixture
def temp_dir():
    """Create a temporary directory for test files."""
    temp_dir = tempfile.mkdtemp()
    yield Path(temp_dir)
    shutil.rmtree(temp_dir)


@pytest.fixture
def sample_terraform_content():
    """Provide sample Terraform content for testing."""
    return {
        "basic_resource": '''resource "vastdata_example" "test" {
  name = "test-resource"
  enabled = true
}''',
        
        "resource_with_blocks": '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
  }
  
  frames {
    name = "frame1"
    ip = "10.0.1.1"
  }
  
  frames {
    name = "frame2" 
    ip = "10.0.1.2"
  }
}''',
        
        "resource_with_lists": '''resource "vastdata_example" "test" {
  name = "test-resource"
  cnode_ids = [1, 2, 3, 4]
  permissions_list = ["read", "write", "execute"]
  roles = [100, 200, 300]
}''',
        
        "resource_with_ip_ranges": '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  client_ip_ranges {
    start_ip = "192.168.1.1"
    end_ip = "192.168.1.100"
  }
  
  client_ip_ranges {
    start_ip = "10.0.0.1"
    end_ip = "10.0.0.50"
  }
}''',
        
        "resource_with_dynamic_block": '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  dynamic "client_ip_ranges" {
    for_each = var.ip_ranges
    content {
      start_ip = client_ip_ranges.value.start
      end_ip = client_ip_ranges.value.end
    }
  }
}''',
        
        "renamed_resource_types": '''resource "vastdata_administators_managers" "test1" {
  name = "admin-manager"
}

resource "vastdata_kafka_brokers" "test2" {
  broker_id = 1
}

resource "vastdata_active_directory2" "test3" {
  domain = "example.com"
}'''
    }


@pytest.fixture
def expected_outputs():
    """Provide expected output content after migration."""
    return {
        "basic_resource": '''resource "vastdata_example" "test" {
  name = "test-resource"
  enabled = true
}''',
        
        "resource_with_blocks": '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  capacity_limits = {
    soft_limit = 1000
    hard_limit = 2000
  }
  
  frames = [
    {
      name = "frame1"
      ip = "10.0.1.1"
    },
    {
      name = "frame2"
      ip = "10.0.1.2"
    },
  ]
}''',
        
        "resource_with_lists": '''resource "vastdata_example" "test" {
  name = "test-resource"
  cnode_ids = "1,2,3,4"
  permissions_list = ["read", "write", "execute"]
  roles = [100, 200, 300]
}''',
        
        "resource_with_ip_ranges": '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  client_ip_ranges = [
    ["192.168.1.1", "192.168.1.100"],
    ["10.0.0.1", "10.0.0.50"]
  ]
}''',
        
        "renamed_resource_types": '''resource "vastdata_administrator_manager" "test1" {
  name = "admin-manager"
}

resource "vastdata_kafka_broker" "test2" {
  broker_id = 1
}

resource "vastdata_active_directory" "test3" {
  domain = "example.com"
}'''
    }


@pytest.fixture
def fixture_files():
    """Provide paths to fixture files."""
    fixtures_dir = Path(__file__).parent / "fixtures"
    return {
        "input_dir": fixtures_dir / "input",
        "expected_dir": fixtures_dir / "expected"
    }


@pytest.fixture
def migration_script_path():
    """Provide path to the migration script."""
    return Path(__file__).parent.parent / "migration_script.py"

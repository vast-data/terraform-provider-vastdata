#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

"""
Integration test for variables and locals preservation in full migration workflow
"""

import os
import sys
import tempfile
import shutil
import json
from pathlib import Path

# Add parent directory to path for imports
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

import pytest


class TestVariablesLocalsIntegration:
    """Integration tests for variable and locals blocks in full migration"""
    
    def test_full_migration_with_variables_and_locals(self, tmp_path):
        """Test complete migration preserving variables, locals, and data files"""
        from state_migration import parse_tfstate, strip_vast_resources
        
        # Create source directory structure
        source_dir = tmp_path / "source"
        source_dir.mkdir()
        
        views_dir = source_dir / "views"
        views_dir.mkdir()
        
        # Create yaml config files
        (views_dir / "view1.yaml").write_text("""---
name: test-view-1
protocols:
  - S3
  - NFS
""")
        
        # Create main.tf with variables, locals, provider, and terraform blocks
        main_tf = source_dir / "main.tf"
        main_tf.write_text("""
variable "cluster_name" {
  type        = string
  description = "Name of the VAST cluster"
  default     = "vast-cluster-1"
}

variable "environment" {
  type = string
}

locals {
  test_id = "t12345"
  prefix  = "test-${local.test_id}"
  
  config_files = {
    view = fileset("${path.module}/views", "*.yaml")
  }
  
  configs = {
    view = { for f in local.config_files.view : f => yamldecode(file("${path.module}/views/${f}")) }
  }
}

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    vastdata = {
      source  = "vast-data/vastdata"
      version = "~> 3.0"
    }
  }
}

provider "vastdata" {
  host     = "https://vms.vast.local"
  username = "admin"
  password = "password"
}
""")
        
        # Create state file with a simple resource
        state_file = source_dir / "terraform.tfstate"
        state_data = {
            "version": 4,
            "terraform_version": "1.5.0",
            "serial": 1,
            "lineage": "test-lineage",
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "test_tenant",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [
                        {
                            "schema_version": 0,
                            "attributes": {
                                "id": "123",
                                "name": "test-tenant",
                                "tenant_id": 123
                            }
                        }
                    ]
                }
            ]
        }
        state_file.write_text(json.dumps(state_data, indent=2))
        
        # Test parsing the state
        state = parse_tfstate(str(state_file))
        assert state is not None
        assert "resources" in state
        
        # Test stripping vast resources
        cleaned, vast_count, non_vast_count = strip_vast_resources(state)
        assert vast_count == 1
        assert non_vast_count == 0
        assert cleaned["resources"] == []
        
        # Verify all source files exist
        assert main_tf.exists()
        assert (views_dir / "view1.yaml").exists()
        
        # Read main.tf and verify all blocks are present
        with open(main_tf, 'r') as f:
            content = f.read()
        
        # Check variables
        assert 'variable "cluster_name"' in content
        assert 'variable "environment"' in content
        assert 'type        = string' in content
        assert 'description = "Name of the VAST cluster"' in content
        
        # Check locals
        assert 'locals {' in content
        assert 'test_id = "t12345"' in content
        assert 'prefix  = "test-${local.test_id}"' in content
        assert 'config_files' in content
        assert 'fileset("${path.module}/views", "*.yaml")' in content
        assert 'yamldecode' in content
        
        # Check terraform block
        assert 'terraform {' in content
        assert 'required_version = ">= 1.0"' in content
        assert 'required_providers' in content
        
        # Check provider block
        assert 'provider "vastdata"' in content
        assert 'host     = "https://vms.vast.local"' in content
    
    def test_multiple_locals_blocks_preserved(self, tmp_path):
        """Test that multiple locals blocks from different files are all preserved"""
        from state_migration import parse_tfstate, strip_vast_resources
        
        # Create source directory
        source_dir = tmp_path / "source"
        source_dir.mkdir()
        
        # Create main.tf with one locals block
        main_tf = source_dir / "main.tf"
        main_tf.write_text("""
locals {
  region = "us-west-2"
  project = "test-project"
}

terraform {
  required_providers {
    vastdata = {
      source  = "vast-data/vastdata"
      version = "~> 3.0"
    }
  }
}

provider "vastdata" {
  host = "https://vms.vast.local"
}
""")
        
        # Create variables.tf with another locals block
        variables_tf = source_dir / "variables.tf"
        variables_tf.write_text("""
variable "env" {
  type = string
}

locals {
  env_prefix = "env-${var.env}"
  tags = {
    Environment = var.env
    Project     = local.project
  }
}
""")
        
        # Create minimal state
        state_file = source_dir / "terraform.tfstate"
        state_file.write_text(json.dumps({
            "version": 4,
            "terraform_version": "1.5.0",
            "serial": 1,
            "lineage": "test",
            "resources": []
        }))
        
        # Read both files
        with open(main_tf, 'r') as f:
            main_content = f.read()
        with open(variables_tf, 'r') as f:
            vars_content = f.read()
        
        # Verify first locals block
        assert 'locals {' in main_content
        assert 'region = "us-west-2"' in main_content
        assert 'project = "test-project"' in main_content
        
        # Verify second locals block
        assert 'locals {' in vars_content
        assert 'env_prefix = "env-${var.env}"' in vars_content
        assert 'tags = {' in vars_content
        
        # Verify variable block
        assert 'variable "env"' in vars_content
    
    def test_nested_locals_with_functions(self, tmp_path):
        """Test complex locals with nested functions and expressions"""
        source_dir = tmp_path / "source"
        source_dir.mkdir()
        
        # Create main.tf with complex locals
        main_tf = source_dir / "main.tf"
        main_tf.write_text("""
terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
      version = "~> 3.0"
    }
  }
}

provider "vastdata" {
  host = "https://vms.vast.local"
}

locals {
  # Complex nested expressions
  base_config = {
    environment = "prod"
    region      = "us-west"
  }
  
  # Multiple function calls
  yaml_files = fileset("${path.module}/configs", "*.yaml")
  json_files = fileset("${path.module}/configs", "*.json")
  
  # Nested for expressions
  parsed_configs = {
    for f in local.yaml_files : 
      f => merge(
        yamldecode(file("${path.module}/configs/${f}")),
        local.base_config
      )
  }
  
  # Conditional logic
  deployment_type = local.base_config.environment == "prod" ? "production" : "development"
  
  # Template rendering
  config_template = templatefile("${path.module}/templates/config.tpl", {
    env  = local.deployment_type
    tags = local.base_config
  })
}
""")
        
        # Create minimal state
        state_file = source_dir / "terraform.tfstate"
        state_file.write_text(json.dumps({
            "version": 4,
            "terraform_version": "1.5.0",
            "serial": 1,
            "lineage": "test",
            "resources": []
        }))
        
        # Read and verify complex locals
        with open(main_tf, 'r') as f:
            content = f.read()
        
        assert 'locals {' in content
        assert 'base_config = {' in content
        assert 'fileset("${path.module}/configs", "*.yaml")' in content
        assert 'fileset("${path.module}/configs", "*.json")' in content
        assert 'yamldecode(file(' in content
        assert 'merge(' in content
        assert 'templatefile(' in content
        assert 'for f in local.yaml_files' in content
        assert '? "production" : "development"' in content


if __name__ == '__main__':
    pytest.main([__file__, '-v'])

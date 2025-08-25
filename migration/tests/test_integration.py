# Copyright (c) HashiCorp, Inc.

"""
Integration tests for the migration tool.
"""

import pytest
import subprocess
import tempfile
import shutil
from pathlib import Path
import sys
import os

# Add the migration script directory to the Python path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import main, transform_file, VERSION


class TestIntegration:
    """Integration test cases for end-to-end migration scenarios."""
    
    def test_migration_script_version(self):
        """Test that the migration script version is correctly defined."""
        assert VERSION == "1.2.1"
    
    def test_provider_version_update(self, temp_dir):
        """Test that VastData provider version is updated from 1.x.x to 2.0.0."""
        # Create input file with provider version 1.x.x
        terraform_content = '''terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
      version = "1.7.0"
    }
  }
}

provider "vastdata" {
  host = "192.168.1.100"
}

resource "vastdata_administators_managers" "admin" {
  username = "admin1"
}
'''
        input_file = temp_dir / "main.tf"
        input_file.write_text(terraform_content)
        
        # Run transformation
        output_file = temp_dir / "main_converted.tf"
        transform_file(input_file, output_file)
        
        # Read and verify output
        result = output_file.read_text()
        
        # Verify provider version was updated
        assert 'version = "2.0.0"' in result
        assert 'version = "1.7.0"' not in result
        
        # Verify resource was also transformed
        assert 'resource "vastdata_administrator_manager"' in result
        assert 'resource "vastdata_administators_managers"' not in result
        
        # Verify provider block structure is preserved
        assert 'source = "vast-data/vastdata"' in result
        assert 'required_providers' in result
        assert 'host = "192.168.1.100"' in result
    
    def test_end_to_end_single_file_migration(self, temp_dir):
        """Test complete migration of a single file."""
        # Create input file with various transformation scenarios
        terraform_content = '''# VastData Terraform Configuration
terraform {
  required_providers {
    vastdata = {
      source = "vastdataorg/vastdata"
      version = "~> 0.9"
    }
  }
}

# Administrator manager with typo in resource name
resource "vastdata_administators_managers" "admin_mgr" {
  name = "primary-admin"
  enabled = true
  
  # Block list that should become attributes
  capacity_limits {
    soft_limit = 5000
    hard_limit = 10000
    enabled = true
  }
  
  # Block list that should become attributes list
  frames {
    name = "frame1"
    ip = "10.0.1.1"
    port = 8080
  }
  
  frames {
    name = "frame2"
    ip = "10.0.1.2"
    port = 8081
  }
  
  # List of numbers that should become string
  cnode_ids = [1, 2, 3, 4]
  
  # List of strings that should become set (no visible change)
  permissions_list = ["read", "write", "admin"]
}

# Kafka broker with plural name
resource "vastdata_kafka_brokers" "broker" {
  broker_id = 1
  name = "kafka-broker-1"
  
  # IP ranges that should become list of lists
  client_ip_ranges {
    start_ip = "192.168.1.1"
    end_ip = "192.168.1.100"
  }
  
  client_ip_ranges {
    start_ip = "10.0.0.1"
    end_ip = "10.0.0.50"
  }
  
  # Dynamic block that should be preserved
  dynamic "security_groups" {
    for_each = var.security_groups
    content {
      name = security_groups.value.name
      description = security_groups.value.description
      rules = security_groups.value.rules
    }
  }
}

# Resource with version suffix
resource "vastdata_active_directory2" "ad" {
  domain = "example.com"
  server = "ad.example.com"
  
  # Mixed static and dynamic content
  default_user_quota {
    hard_limit = 1000000
    soft_limit = 800000
  }
  
  dynamic "user_groups" {
    for_each = var.ad_groups
    content {
      name = user_groups.value.name
      dn = user_groups.value.distinguished_name
    }
  }
}'''
        
        input_file = temp_dir / "test_config.tf"
        output_file = temp_dir / "test_config_converted.tf"
        
        with open(input_file, 'w') as f:
            f.write(terraform_content)
        
        # Run the transformation
        transform_file(input_file, output_file)
        
        # Verify output file was created
        assert output_file.exists()
        
        with open(output_file, 'r') as f:
            result = f.read()
        
        # Test resource renaming
        assert "vastdata_administrator_manager" in result
        assert "vastdata_kafka_broker" in result
        assert "vastdata_active_directory" in result
        assert "vastdata_administators_managers" not in result
        assert "vastdata_kafka_brokers" not in result
        assert "vastdata_active_directory2" not in result
        
        # Test schema transformations
        assert 'capacity_limits = {' in result
        assert 'frames = [' in result
        assert 'cnode_ids = [1, 2, 3, 4]' in result  # cnode_ids should remain as list
        assert 'client_ip_ranges = [' in result
        assert '["192.168.1.1", "192.168.1.100"]' in result
        assert '["10.0.0.1", "10.0.0.50"]' in result
        
        # Test dynamic block preservation
        assert 'dynamic "security_groups" {' in result
        assert 'dynamic "user_groups" {' in result
        assert 'for_each = var.security_groups' in result
        assert 'for_each = var.ad_groups' in result
        
        # Test that comments and other content is preserved
        assert "# VastData Terraform Configuration" in result
        assert "terraform {" in result
        assert "required_providers" in result
    
    def test_migration_with_multiple_files(self, temp_dir):
        """Test migration of multiple files in a directory structure."""
        # Create directory structure
        src_dir = temp_dir / "source"
        dst_dir = temp_dir / "destination"
        src_dir.mkdir()
        
        # Create main.tf
        main_tf_content = '''resource "vastdata_administators_managers" "main_admin" {
  name = "main-admin"
  
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
  }
}'''
        
        # Create variables.tf
        variables_tf_content = '''variable "ip_ranges" {
  description = "List of IP ranges"
  type = list(object({
    start = string
    end = string
  }))
  default = []
}

variable "admin_settings" {
  description = "Administrator settings"
  type = object({
    name = string
    enabled = bool
  })
}'''
        
        # Create modules/storage/main.tf
        modules_dir = src_dir / "modules" / "storage"
        modules_dir.mkdir(parents=True)
        
        storage_tf_content = '''resource "vastdata_kafka_brokers" "storage_broker" {
  broker_id = 1
  
  client_ip_ranges {
    start_ip = "10.0.0.1"
    end_ip = "10.0.0.100"
  }
  
  cnode_ids = [1, 2, 3]
}'''
        
        # Write files
        (src_dir / "main.tf").write_text(main_tf_content)
        (src_dir / "variables.tf").write_text(variables_tf_content)
        (modules_dir / "main.tf").write_text(storage_tf_content)
        
        # Run migration using the main function
        main(str(src_dir), str(dst_dir))
        
        # Verify converted files exist
        assert (dst_dir / "main_converted.tf").exists()
        assert (dst_dir / "variables_converted.tf").exists()
        assert (dst_dir / "modules" / "storage" / "main_converted.tf").exists()
        
        # Check main.tf transformations
        main_result = (dst_dir / "main_converted.tf").read_text()
        assert "vastdata_administrator_manager" in main_result
        assert 'capacity_limits = {' in main_result
        
        # Check variables.tf is preserved (no resources to transform)
        variables_result = (dst_dir / "variables_converted.tf").read_text()
        assert "variable \"ip_ranges\"" in variables_result
        assert "variable \"admin_settings\"" in variables_result
        
        # Check storage module transformations
        storage_result = (dst_dir / "modules" / "storage" / "main_converted.tf").read_text()
        assert "vastdata_kafka_broker" in storage_result
        assert 'client_ip_ranges = [' in storage_result
        assert 'cnode_ids = [1, 2, 3]' in storage_result  # cnode_ids should remain as list
    
    def test_migration_preserves_file_structure(self, temp_dir):
        """Test that migration preserves the original file and directory structure."""
        # Create complex directory structure
        src_dir = temp_dir / "complex_source"
        dst_dir = temp_dir / "complex_destination"
        
        # Create nested directories
        (src_dir / "environments" / "prod").mkdir(parents=True)
        (src_dir / "environments" / "dev").mkdir(parents=True)
        (src_dir / "modules" / "networking").mkdir(parents=True)
        (src_dir / "modules" / "compute").mkdir(parents=True)
        
        # Create files in different locations
        files_to_create = [
            ("main.tf", 'resource "vastdata_administators_managers" "main" { name = "main" }'),
            ("environments/prod/main.tf", 'resource "vastdata_kafka_brokers" "prod" { broker_id = 1 }'),
            ("environments/dev/main.tf", 'resource "vastdata_active_directory2" "dev" { domain = "dev.com" }'),
            ("modules/networking/main.tf", 'resource "vastdata_non_local_user" "net_user" { username = "netuser" }'),
            ("modules/compute/main.tf", 'resource "vastdata_saml" "compute_saml" { provider = "okta" }'),
        ]
        
        for file_path, content in files_to_create:
            full_path = src_dir / file_path
            full_path.write_text(content)
        
        # Run migration
        main(str(src_dir), str(dst_dir))
        
        # Verify all converted files exist in correct locations
        expected_files = [
            "main_converted.tf",
            "environments/prod/main_converted.tf",
            "environments/dev/main_converted.tf", 
            "modules/networking/main_converted.tf",
            "modules/compute/main_converted.tf",
        ]
        
        for file_path in expected_files:
            converted_file = dst_dir / file_path
            assert converted_file.exists(), f"Expected file not found: {file_path}"
        
        # Verify transformations in each file
        main_result = (dst_dir / "main_converted.tf").read_text()
        assert "vastdata_administrator_manager" in main_result
        
        prod_result = (dst_dir / "environments/prod/main_converted.tf").read_text()
        assert "vastdata_kafka_broker" in prod_result
        
        dev_result = (dst_dir / "environments/dev/main_converted.tf").read_text()
        assert "vastdata_active_directory" in dev_result
        
        net_result = (dst_dir / "modules/networking/main_converted.tf").read_text()
        assert "vastdata_nonlocal_user" in net_result
        
        compute_result = (dst_dir / "modules/compute/main_converted.tf").read_text()
        assert "vastdata_saml_config" in compute_result
    
    def test_migration_handles_empty_files(self, temp_dir):
        """Test that migration handles empty files correctly."""
        src_dir = temp_dir / "empty_source"
        dst_dir = temp_dir / "empty_destination"
        src_dir.mkdir()
        
        # Create empty file
        empty_file = src_dir / "empty.tf"
        empty_file.touch()
        
        # Create file with only comments
        comments_file = src_dir / "comments.tf"
        comments_file.write_text('''# This is a comment
# Another comment

# Yet another comment''')
        
        # Run migration
        main(str(src_dir), str(dst_dir))
        
        # Verify files are created but content is preserved
        assert (dst_dir / "empty_converted.tf").exists()
        assert (dst_dir / "comments_converted.tf").exists()
        
        # Check that empty file remains empty
        empty_result = (dst_dir / "empty_converted.tf").read_text()
        assert empty_result == ""
        
        # Check that comments are preserved
        comments_result = (dst_dir / "comments_converted.tf").read_text()
        assert "# This is a comment" in comments_result
        assert "# Another comment" in comments_result
        assert "# Yet another comment" in comments_result
    
    def test_migration_handles_non_terraform_blocks(self, temp_dir):
        """Test that migration correctly handles non-resource blocks."""
        terraform_content = '''terraform {
  required_version = ">= 1.0"
  required_providers {
    vastdata = {
      source = "vastdataorg/vastdata"
      version = "~> 1.0"
    }
  }
  
  backend "s3" {
    bucket = "terraform-state"
    key = "vastdata/terraform.tfstate"
    region = "us-west-2"
  }
}

provider "vastdata" {
  endpoint = var.vastdata_endpoint
  username = var.vastdata_username
  password = var.vastdata_password
}

variable "environment" {
  description = "Environment name"
  type = string
  default = "dev"
}

locals {
  common_tags = {
    Environment = var.environment
    Project = "vastdata-migration"
  }
}

data "vastdata_cluster_info" "current" {}

output "cluster_version" {
  value = data.vastdata_cluster_info.current.version
}

resource "vastdata_administators_managers" "admin" {
  name = "test-admin"
  tags = local.common_tags
}'''
        
        src_dir = temp_dir / "mixed_source"
        dst_dir = temp_dir / "mixed_destination"
        src_dir.mkdir()
        
        input_file = src_dir / "mixed.tf"
        input_file.write_text(terraform_content)
        
        # Run migration
        main(str(src_dir), str(dst_dir))
        
        result = (dst_dir / "mixed_converted.tf").read_text()
        
        # Verify non-resource blocks are preserved
        assert "terraform {" in result
        assert "required_version" in result
        assert "required_providers" in result
        assert "backend \"s3\"" in result
        assert "provider \"vastdata\"" in result
        assert "variable \"environment\"" in result
        assert "locals {" in result
        assert "data \"vastdata_cluster_info\"" in result
        assert "output \"cluster_version\"" in result
        
        # Verify resource transformation occurred
        assert "vastdata_administrator_manager" in result
        assert "vastdata_administators_managers" not in result
        
        # Verify references are preserved
        assert "var.vastdata_endpoint" in result
        assert "local.common_tags" in result
        assert "data.vastdata_cluster_info.current.version" in result
    
    @pytest.mark.slow
    def test_large_file_migration_performance(self, temp_dir):
        """Test migration performance with a large file (marked as slow test)."""
        # Create a large file with many resources
        large_content_parts = []
        
        # Add header
        large_content_parts.append('''terraform {
  required_providers {
    vastdata = {
      source = "vastdataorg/vastdata"
    }
  }
}''')
        
        # Generate many resources
        for i in range(100):
            resource_content = f'''
resource "vastdata_administators_managers" "admin_{i}" {{
  name = "admin-{i}"
  enabled = true
  
  capacity_limits {{
    soft_limit = {1000 + i * 100}
    hard_limit = {2000 + i * 100}
  }}
  
  frames {{
    name = "frame_{i}_1"
    ip = "10.0.{i}.1"
  }}
  
  frames {{
    name = "frame_{i}_2"
    ip = "10.0.{i}.2"
  }}
  
  cnode_ids = [{i}, {i+1}, {i+2}]
  permissions_list = ["read", "write"]
}}'''
            large_content_parts.append(resource_content)
        
        large_content = '\n'.join(large_content_parts)
        
        src_dir = temp_dir / "large_source"
        dst_dir = temp_dir / "large_destination"
        src_dir.mkdir()
        
        input_file = src_dir / "large.tf"
        input_file.write_text(large_content)
        
        # Measure migration time
        import time
        start_time = time.time()
        
        main(str(src_dir), str(dst_dir))
        
        end_time = time.time()
        migration_time = end_time - start_time
        
        # Verify migration completed
        output_file = dst_dir / "large_converted.tf"
        assert output_file.exists()
        
        result = output_file.read_text()
        
        # Verify all transformations occurred
        assert "vastdata_administrator_manager" in result
        assert "vastdata_administators_managers" not in result
        assert 'capacity_limits = {' in result
        assert 'frames = [' in result
        
        # Count transformed resources (should be 100)
        admin_manager_count = result.count("vastdata_administrator_manager")
        assert admin_manager_count == 100
        
        # Performance assertion (should complete within reasonable time)
        # This is a rough benchmark - adjust based on expected performance
        assert migration_time < 30.0, f"Migration took too long: {migration_time:.2f} seconds"

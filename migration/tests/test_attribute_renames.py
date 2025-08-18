# Copyright (c) HashiCorp, Inc.

"""
Tests for attribute name changes during migration.
"""

import pytest
from pathlib import Path
import sys
import os

# Add the migration script directory to the Python path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import (
    transform_resource_block,
    transform_file
)


class TestAttributeRenames:
    """Test cases for attribute name changes during migration."""
    
    def test_type_underscore_to_type_transformation(self):
        """Test that type_ is renamed to type."""
        terraform_content = '''resource "vastdata_s3_replication_peers" "s3peer" {
  name              = "s3peer"
  bucket_name       = "s3bucket"
  http_protocol     = "https"
  type_             = "CUSTOM_S3"
  custom_bucket_url = "customs3.bucket.com"
  access_key        = "W21E6X5ZQEOODYB6J0UY"
  secret_key        = "fcESVNih9Ykb/bDSmKipQdinnHObrRyv9nre+nR1"
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that type_ is renamed to type
        assert 'type = "CUSTOM_S3"' in result
        assert 'type_' not in result
        
        # Check that resource name is also renamed
        assert "vastdata_s3_replication_peer" in result
        assert "vastdata_s3_replication_peers" not in result
    
    def test_permissions_list_to_permissions_transformation(self):
        """Test that permissions_list is renamed to permissions."""
        terraform_content = '''resource "vastdata_administators_roles" "role1" {
  name             = "role1"
  permissions_list = ["create_support", "create_settings", "create_security"]
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that permissions_list is renamed to permissions
        assert 'permissions = ["create_support", "create_settings", "create_security"]' in result
        assert 'permissions_list' not in result
        
        # Check that resource name is also renamed
        assert "vastdata_administrator_role" in result
        assert "vastdata_administators_roles" not in result
    
    def test_multiple_attribute_renames_in_single_resource(self):
        """Test multiple attribute renames in the same resource."""
        terraform_content = '''resource "vastdata_s3_replication_peers" "complex" {
  name              = "complex-peer"
  type_             = "AWS_S3"
  aws_region        = "us-west-2"
  permissions_list  = ["read", "write"]
  bucket_name       = "my-bucket"
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check both attribute renames
        assert 'type = "AWS_S3"' in result
        assert 'permissions = ["read", "write"]' in result
        assert 'type_' not in result
        assert 'permissions_list' not in result
        
        # Check resource rename
        assert "vastdata_s3_replication_peer" in result
    
    def test_attribute_renames_preserve_indentation(self):
        """Test that attribute renames preserve original indentation."""
        terraform_content = '''resource "vastdata_administators_roles" "test" {
  name = "test-role"
  
  # This should preserve indentation
  permissions_list = [
    "create_support",
    "create_settings"
  ]
  
  # Nested block
  settings {
    type_ = "ADVANCED"
    enabled = true
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that indentation is preserved
        assert '  permissions = [' in result
        assert '    type = "ADVANCED"' in result
    
    def test_attribute_renames_with_complex_values(self):
        """Test attribute renames with complex value expressions."""
        terraform_content = '''resource "vastdata_administators_managers" "manager" {
  username = "test-manager"
  permissions_list = concat(
    ["create_support"],
    var.additional_permissions
  )
  
  type_ = var.s3_type != null ? var.s3_type : "CUSTOM_S3"
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that complex expressions are preserved
        assert 'permissions = concat(' in result
        assert 'var.additional_permissions' in result
        assert 'type = var.s3_type != null ? var.s3_type : "CUSTOM_S3"' in result
    
    def test_no_change_for_correct_attribute_names(self):
        """Test that already correct attribute names are not modified."""
        terraform_content = '''resource "vastdata_administrator_manager" "manager" {
  username = "test-manager"
  permissions = ["read", "write"]
  type = "STANDARD"
  enabled = true
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that correct attributes are unchanged
        assert 'permissions = ["read", "write"]' in result
        assert 'type = "STANDARD"' in result
        # Should not add underscore
        assert 'permissions_list' not in result
        assert 'type_' not in result
    
    def test_attribute_rename_in_nested_blocks(self):
        """Test that attribute renames work within nested blocks."""
        terraform_content = '''resource "vastdata_administators_managers" "manager" {
  username = "test-manager"
  
  advanced_settings {
    type_ = "CUSTOM"
    permissions_list = ["admin"]
    enabled = true
  }
  
  backup_config {
    type_ = "INCREMENTAL"
    schedule = "daily"
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that nested attribute renames work
        assert 'type = "CUSTOM"' in result
        assert 'permissions = ["admin"]' in result
        assert 'type = "INCREMENTAL"' in result
        # Originals should be gone
        assert 'type_' not in result
        assert 'permissions_list' not in result
    
    def test_file_level_attribute_renames(self, temp_dir):
        """Test attribute renames work at the file level with multiple resources."""
        terraform_content = '''# Test configuration with attribute renames
        
resource "vastdata_administators_roles" "role1" {
  name = "role1"
  permissions_list = ["create_support", "create_settings"]
}

resource "vastdata_s3_replication_peers" "peer1" {
  name = "peer1"
  type_ = "AWS_S3"
  bucket_name = "test-bucket"
}

# Regular resource without renames
resource "vastdata_user" "user1" {
  name = "user1"
  uid = 1001
}

resource "vastdata_administators_managers" "manager1" {
  username = "manager1"
  permissions_list = ["create_monitoring"]
  roles = [1, 2, 3]
}'''
        
        input_file = temp_dir / "attribute_renames_test.tf"
        output_file = temp_dir / "attribute_renames_test_converted.tf"
        
        with open(input_file, 'w') as f:
            f.write(terraform_content)
        
        transform_file(input_file, output_file)
        
        with open(output_file, 'r') as f:
            result = f.read()
        
        # Check resource renames
        assert "vastdata_administrator_role" in result
        assert "vastdata_s3_replication_peer" in result
        assert "vastdata_administrator_manager" in result
        
        # Check attribute renames
        assert 'permissions = ["create_support", "create_settings"]' in result
        assert 'type = "AWS_S3"' in result
        assert 'permissions = ["create_monitoring"]' in result
        
        # Check originals are gone
        assert 'permissions_list' not in result
        assert 'type_' not in result
        
        # Check that unchanged resources remain unchanged
        assert 'vastdata_user' in result
        assert 'name = "user1"' in result
        
        # Check that roles list is transformed correctly (List of Number -> Set)
        assert 'roles = [1, 2, 3]' in result
    
    def test_edge_cases_for_attribute_renames(self):
        """Test edge cases for attribute renaming."""
        terraform_content = '''resource "vastdata_test" "edge_cases" {
  # Comments with type_ and permissions_list should not be changed
  
  # Test attribute names that contain but are not exactly the target
  my_type_ = "should_not_change"
  permissions_list_backup = "should_not_change"
  
  # Test exact matches
  type_ = "EXACT_MATCH"
  permissions_list = ["EXACT_MATCH"]
  
  # Test in strings (should not change)
  description = "This uses type_ and permissions_list internally"
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check that exact matches are renamed
        assert 'type = "EXACT_MATCH"' in result
        assert 'permissions = ["EXACT_MATCH"]' in result
        
        # Check that partial matches are not renamed
        assert 'my_type_ = "should_not_change"' in result
        assert 'permissions_list_backup = "should_not_change"' in result
        
        # Check that strings are not modified
        assert 'description = "This uses type_ and permissions_list internally"' in result
        
        # Check that comments are preserved
        assert '# Comments with type_ and permissions_list should not be changed' in result

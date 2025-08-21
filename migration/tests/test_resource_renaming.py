# Copyright (c) HashiCorp, Inc.

"""
Tests for resource type renaming functionality.
"""

import pytest
from pathlib import Path
import sys
import os

# Add the migration script directory to the Python path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import (
    transform_resource_block,
    resource_type_rename_map,
    transform_file
)


class TestResourceRenaming:
    """Test cases for resource type renaming."""
    
    def test_resource_type_rename_map_completeness(self):
        """Test that all expected resource type mappings are present."""
        expected_mappings = {
            # Administrators (typo fix + plural→singular)
            "vastdata_administators_managers": "vastdata_administrator_manager",
            "vastdata_administators_roles": "vastdata_administrator_role", 
            "vastdata_administators_realms": "vastdata_administrator_realm",
            # Other plural→singular renames
            "vastdata_kafka_brokers": "vastdata_kafka_broker",
            "vastdata_replication_peers": "vastdata_replication_peer",
            "vastdata_s3_replication_peers": "vastdata_s3_replication_peer",
            # Specific resource renames
            "vastdata_active_directory2": "vastdata_active_directory",
            "vastdata_non_local_user": "vastdata_nonlocal_user",
            "vastdata_non_local_user_key": "vastdata_nonlocal_user_key",
            "vastdata_non_local_group": "vastdata_nonlocal_group",
            "vastdata_saml": "vastdata_saml_config",
        }
        
        for old_name, new_name in expected_mappings.items():
            assert old_name in resource_type_rename_map
            assert resource_type_rename_map[old_name] == new_name
    
    def test_administrator_typo_fixes(self):
        """Test that administrator typos are fixed correctly."""
        terraform_content = '''resource "vastdata_administators_managers" "test" {
  name = "admin-manager"
  enabled = true
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert "vastdata_administrator_manager" in result
        assert "vastdata_administators_managers" not in result
        
    def test_plural_to_singular_conversion(self):
        """Test conversion from plural to singular resource names."""
        test_cases = [
            ("vastdata_kafka_brokers", "vastdata_kafka_broker"),
            ("vastdata_replication_peers", "vastdata_replication_peer"),
            ("vastdata_s3_replication_peers", "vastdata_s3_replication_peer")
        ]
        
        for old_type, new_type in test_cases:
            terraform_content = f'''resource "{old_type}" "test" {{
  name = "test-resource"
}}'''
            
            lines = terraform_content.split('\n')
            result, consumed = transform_resource_block(lines, 0)
            
            assert result is not None
            assert new_type in result
            assert old_type not in result
    
    def test_specific_resource_renames(self):
        """Test specific resource type renames."""
        test_cases = [
            ("vastdata_active_directory2", "vastdata_active_directory"),
            ("vastdata_non_local_user", "vastdata_nonlocal_user"),
            ("vastdata_non_local_user_key", "vastdata_nonlocal_user_key"),
            ("vastdata_non_local_group", "vastdata_nonlocal_group"),
            ("vastdata_saml", "vastdata_saml_config")
        ]
        
        for old_type, new_type in test_cases:
            terraform_content = f'''resource "{old_type}" "test" {{
  domain = "example.com"
}}'''
            
            lines = terraform_content.split('\n')
            result, consumed = transform_resource_block(lines, 0)
            
            assert result is not None
            assert new_type in result
            # Check that the old type is not used as a complete resource name
            assert f'resource "{old_type}"' not in result
    
    def test_resource_instance_name_preserved(self):
        """Test that resource instance names are preserved during renaming."""
        terraform_content = '''resource "vastdata_administators_managers" "my_custom_instance" {
  name = "admin-manager"
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert '"my_custom_instance"' in result
        assert "vastdata_administrator_manager" in result
    
    def test_no_rename_for_unknown_types(self):
        """Test that unknown resource types are not modified."""
        terraform_content = '''resource "vastdata_unknown_resource" "test" {
  name = "test-resource"
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert "vastdata_unknown_resource" in result
    
    def test_multiple_resources_in_file(self, temp_dir):
        """Test renaming multiple resources in a single file."""
        terraform_content = '''terraform {
  required_providers {
    vastdata = {
      source = "vastdataorg/vastdata"
    }
  }
}

resource "vastdata_administators_managers" "manager1" {
  name = "manager1"
}

resource "vastdata_kafka_brokers" "broker1" {
  broker_id = 1
}

resource "vastdata_active_directory2" "ad1" {
  domain = "example.com"
}

resource "vastdata_unknown_resource" "unknown1" {
  name = "should_not_change"
}'''
        
        input_file = temp_dir / "test.tf"
        output_file = temp_dir / "test_converted.tf"
        
        with open(input_file, 'w') as f:
            f.write(terraform_content)
        
        transform_file(input_file, output_file)
        
        with open(output_file, 'r') as f:
            result = f.read()
        
        # Check that renames occurred
        assert "vastdata_administrator_manager" in result
        assert "vastdata_kafka_broker" in result  
        assert "vastdata_active_directory" in result
        
        # Check that original names are gone
        assert "vastdata_administators_managers" not in result
        assert "vastdata_kafka_brokers" not in result
        assert "vastdata_active_directory2" not in result
        
        # Check that unknown resource is unchanged
        assert "vastdata_unknown_resource" in result
        
        # Check that non-resource content is preserved
        assert "terraform {" in result
        assert "required_providers" in result
        assert "vastdataorg/vastdata" in result
    
    def test_quoted_resource_types(self):
        """Test that quoted resource types are handled correctly."""
        terraform_content = '''resource "vastdata_administators_managers" "test" {
  name = "admin-manager"
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert '"vastdata_administrator_manager"' in result
        assert '"vastdata_administators_managers"' not in result
    
    def test_preserve_attributes_during_rename(self):
        """Test that resource attributes are preserved during type renaming."""
        terraform_content = '''resource "vastdata_administators_managers" "test" {
  name = "admin-manager"
  enabled = true
  tags = {
    environment = "test"
    team = "platform"
  }
  
  settings {
    auto_backup = true
    retention_days = 30
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert "vastdata_administrator_manager" in result
        
        # Check that all attributes are preserved
        assert 'name = "admin-manager"' in result
        assert 'enabled = true' in result
        assert 'environment = "test"' in result
        assert 'team = "platform"' in result
        assert 'auto_backup = true' in result
        assert 'retention_days = 30' in result

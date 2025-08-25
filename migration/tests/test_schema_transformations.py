# Copyright (c) HashiCorp, Inc.

"""
Tests for schema transformation functionality.
"""

import pytest
from pathlib import Path
import sys
import os

# Add the migration script directory to the Python path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import (
    transform_resource_block,
    get_group_for_key,
    parse_nested_block,
    key_groups,
    transform_file
)


class TestSchemaTransformations:
    """Test cases for schema transformations."""
    
    def test_key_groups_completeness(self):
        """Test that key groups contain expected attributes."""
        expected_groups = {
            "Block List --> Attributes": [
                "capacity_total_limits", "capacity_limits", "static_limits", "static_total_limits",
                "default_group_quota", "default_user_quota", "share_acl", "owner_root_snapshot", 
                "owner_tenant", "bucket_logging", "protocols_audit"
            ],
            "Block List --> Attributes List": ["frames"],
            "Block List --> List of Maps": ["addresses", "group_quotas", "user_quotas"],
            "Block List --> Attributes Set": [],
            "List of Number --> Set of Number": ["roles", "s3_policies_ids", "gids", "tenants"],
            "Block List --> List of List of String": ["client_ip_ranges", "ip_ranges"],
            "List of Number --> String": ["active_cnode_ids"],
            "List of String --> Set of String": [
                "object_types", "ldap_groups", "permissions_list", "groups", "users",
                "abac_tags", "hosts", "abe_protocols", "bucket_creators", "bucket_creators_groups",
                "nfs_all_squash", "nfs_no_squash", "nfs_read_only"
            ]
        }
        
        for group_name, attributes in expected_groups.items():
            assert group_name in key_groups
            for attr in attributes:
                assert attr in key_groups[group_name]
    
    def test_get_group_for_key_function(self):
        """Test the get_group_for_key function."""
        test_cases = [
            ("capacity_limits", "Block List --> Attributes"),
            ("frames", "Block List --> Attributes List"),
            ("group_quotas", "Block List --> List of Maps"),
            ("user_quotas", "Block List --> List of Maps"),
            ("addresses", "Block List --> List of Maps"),
            ("roles", "List of Number --> Set of Number"),
            ("client_ip_ranges", "Block List --> List of List of String"),
            ("active_cnode_ids", "List of Number --> String"),
            ("permissions_list", "List of String --> Set of String"),
            ("unknown_attribute", None)
        ]
        
        for attribute, expected_group in test_cases:
            result = get_group_for_key(attribute)
            assert result == expected_group
    
    def test_block_list_to_attributes_transformation(self):
        """Test Block List --> Attributes transformation."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
    enabled = true
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        expected_lines = [
            'capacity_limits = {',
            'soft_limit = 1000',
            'hard_limit = 2000', 
            'enabled = true',
            '}'
        ]
        
        for line in expected_lines:
            assert line.strip() in result
    
    def test_block_list_to_attributes_list_transformation(self):
        """Test Block List --> Attributes List transformation."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  frames {
    name = "frame1"
    ip = "10.0.1.1"
    enabled = true
  }
  
  frames {
    name = "frame2"
    ip = "10.0.1.2"
    enabled = false
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        expected_elements = [
            'frames = [',
            'name = "frame1"',
            'ip = "10.0.1.1"',
            'enabled = true',
            'name = "frame2"',
            'ip = "10.0.1.2"',
            'enabled = false',
            ']'
        ]
        
        for element in expected_elements:
            assert element.strip() in result
    
    def test_list_of_number_to_string_transformation(self):
        """Test List of Number --> String transformation."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  active_cnode_ids = [1, 2, 3, 4, 5]
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert 'active_cnode_ids = "1,2,3,4,5"' in result
        assert 'active_cnode_ids = [1, 2, 3, 4, 5]' not in result
    
    def test_list_of_number_to_string_with_negative_numbers(self):
        """Test List of Number --> String transformation with negative numbers."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  active_cnode_ids = [-1, 0, 1, 2, -3]
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert 'active_cnode_ids = "-1,0,1,2,-3"' in result
    
    def test_block_list_to_list_of_list_of_string_transformation(self):
        """Test Block List --> List of List of String transformation."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  client_ip_ranges {
    start_ip = "192.168.1.1"
    end_ip = "192.168.1.100"
  }
  
  client_ip_ranges {
    start_ip = "10.0.0.1"
    end_ip = "10.0.0.50"
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        expected_transformation = 'client_ip_ranges = ['
        assert expected_transformation in result
        assert '["192.168.1.1", "192.168.1.100"]' in result
        assert '["10.0.0.1", "10.0.0.50"]' in result
    
    def test_parse_nested_block_function(self):
        """Test the parse_nested_block function."""
        lines = [
            'capacity_limits {',
            '  soft_limit = 1000',
            '  hard_limit = 2000',
            '  enabled = true',
            '}'
        ]
        
        attrs, consumed = parse_nested_block(lines, 0)
        
        assert consumed == 5
        assert attrs['soft_limit'] == '1000'
        assert attrs['hard_limit'] == '2000'
        assert attrs['enabled'] == 'true'
    
    def test_nested_block_with_complex_values(self):
        """Test parsing nested blocks with complex values."""
        lines = [
            'settings {',
            '  name = "test-setting"',
            '  count = 42',
            '  enabled = true',
            '  tags = ["tag1", "tag2"]',
            '  config = { key = "value" }',
            '}'
        ]
        
        attrs, consumed = parse_nested_block(lines, 0)
        
        assert consumed == 7
        assert attrs['name'] == '"test-setting"'
        assert attrs['count'] == '42'
        assert attrs['enabled'] == 'true'
        assert attrs['tags'] == '["tag1", "tag2"]'
        assert attrs['config'] == '{ key = "value" }'
    
    def test_mixed_transformations_in_single_resource(self):
        """Test multiple transformation types in a single resource."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  # Block List --> Attributes
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
  }
  
  # Block List --> Attributes List  
  frames {
    name = "frame1"
    ip = "10.0.1.1"
  }
  
  # List of Number --> String (for active_cnode_ids only)
  active_cnode_ids = [1, 2, 3]
  
  # cnode_ids should remain as list (no transformation)
  cnode_ids = [4, 5, 6]
  
  # List of String --> Set of String (no transformation needed)
  permissions_list = ["read", "write"]
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check capacity_limits transformation
        assert 'capacity_limits = {' in result
        assert 'soft_limit = 1000' in result
        
        # Check frames transformation  
        assert 'frames = [' in result
        assert 'name = "frame1"' in result
        
        # Check active_cnode_ids transformation (should become string)
        assert 'active_cnode_ids = "1,2,3"' in result
        
        # Check cnode_ids remains as list (no transformation)
        assert 'cnode_ids = [4, 5, 6]' in result
        
        # Check permissions_list is transformed to permissions
        assert 'permissions = ["read", "write"]' in result
    
    def test_transformation_preserves_indentation(self):
        """Test that transformations preserve proper indentation."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check that the transformed block maintains proper indentation
        lines_in_result = result.split('\n')
        capacity_limits_line = None
        for i, line in enumerate(lines_in_result):
            if 'capacity_limits = {' in line:
                capacity_limits_line = i
                break
        
        assert capacity_limits_line is not None
        # The line should be indented with 2 spaces
        assert lines_in_result[capacity_limits_line].startswith('  ')
    
    def test_empty_block_transformation(self):
        """Test transformation of empty blocks."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  capacity_limits {
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert 'capacity_limits = {' in result
        assert '}' in result
    
    def test_complex_file_transformation(self, temp_dir):
        """Test transformation of a complete file with multiple resources and transformations."""
        terraform_content = '''# VastData configuration
terraform {
  required_providers {
    vastdata = {
      source = "vastdataorg/vastdata"
    }
  }
}

resource "vastdata_administators_managers" "manager1" {
  name = "manager1"
  
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
  }
  
  frames {
    name = "frame1"
    ip = "10.0.1.1"
  }
  
  cnode_ids = [1, 2, 3]  # Should remain as list
  active_cnode_ids = [4, 5, 6]  # Should become string
}

resource "vastdata_kafka_brokers" "broker1" {
  broker_id = 1
  
  client_ip_ranges {
    start_ip = "192.168.1.1"
    end_ip = "192.168.1.100"
  }
  
  permissions_list = ["read", "write", "execute"]
}'''
        
        input_file = temp_dir / "complex_test.tf"
        output_file = temp_dir / "complex_test_converted.tf"
        
        with open(input_file, 'w') as f:
            f.write(terraform_content)
        
        transform_file(input_file, output_file)
        
        with open(output_file, 'r') as f:
            result = f.read()
        
        # Check resource renaming
        assert "vastdata_administrator_manager" in result
        assert "vastdata_kafka_broker" in result
        
        # Check schema transformations
        assert 'capacity_limits = {' in result
        assert 'frames = [' in result
        assert 'cnode_ids = [1, 2, 3]' in result  # Should remain as list
        assert 'active_cnode_ids = "4,5,6"' in result  # Should become string
        assert 'client_ip_ranges = [' in result
        assert '["192.168.1.1", "192.168.1.100"]' in result
        
        # Check preserved content
        assert "# VastData configuration" in result
        assert "terraform {" in result
        assert "required_providers" in result

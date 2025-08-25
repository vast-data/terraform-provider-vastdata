# Copyright (c) HashiCorp, Inc.

"""
Tests for dynamic block handling functionality.
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


class TestDynamicBlocks:
    """Test cases for dynamic block preservation."""
    
    def test_simple_dynamic_block_preservation(self):
        """Test that simple dynamic blocks are preserved unchanged."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  dynamic "client_ip_ranges" {
    for_each = var.ip_ranges
    content {
      start_ip = client_ip_ranges.value.start
      end_ip = client_ip_ranges.value.end
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check that dynamic block is preserved exactly
        expected_lines = [
            'dynamic "client_ip_ranges" {',
            'for_each = var.ip_ranges',
            'content {',
            'start_ip = client_ip_ranges.value.start',
            'end_ip = client_ip_ranges.value.end',
            '}',
            '}'
        ]
        
        for line in expected_lines:
            assert line.strip() in result
    
    def test_multiple_dynamic_blocks_preservation(self):
        """Test that multiple dynamic blocks are preserved."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  dynamic "client_ip_ranges" {
    for_each = var.ip_ranges
    content {
      start_ip = client_ip_ranges.value.start
      end_ip = client_ip_ranges.value.end
    }
  }
  
  dynamic "permissions" {
    for_each = var.user_permissions
    content {
      user = permissions.value.user
      role = permissions.value.role
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check both dynamic blocks are preserved
        assert 'dynamic "client_ip_ranges" {' in result
        assert 'dynamic "permissions" {' in result
        assert 'for_each = var.ip_ranges' in result
        assert 'for_each = var.user_permissions' in result
    
    def test_nested_dynamic_blocks_preservation(self):
        """Test that nested dynamic blocks are preserved."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  dynamic "network_config" {
    for_each = var.networks
    content {
      name = network_config.value.name
      
      dynamic "subnet" {
        for_each = network_config.value.subnets
        content {
          cidr = subnet.value.cidr
          gateway = subnet.value.gateway
        }
      }
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check that nested dynamic blocks are preserved
        assert 'dynamic "network_config" {' in result
        assert 'dynamic "subnet" {' in result
        assert 'for_each = var.networks' in result
        assert 'for_each = network_config.value.subnets' in result
    
    def test_dynamic_block_with_complex_expressions(self):
        """Test dynamic blocks with complex expressions are preserved."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  dynamic "client_ip_ranges" {
    for_each = { for idx, range in var.ip_ranges : idx => range if range.enabled }
    iterator = ip_range
    content {
      start_ip = ip_range.value.start_ip
      end_ip = ip_range.value.end_ip
      description = "Range ${ip_range.key}: ${ip_range.value.description}"
      tags = merge(var.default_tags, ip_range.value.tags)
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check that complex expressions are preserved
        assert 'for_each = { for idx, range in var.ip_ranges : idx => range if range.enabled }' in result
        assert 'iterator = ip_range' in result
        assert '"Range ${ip_range.key}: ${ip_range.value.description}"' in result
        assert 'merge(var.default_tags, ip_range.value.tags)' in result
    
    def test_dynamic_block_with_static_attributes(self):
        """Test resources with both dynamic blocks and static attributes."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  enabled = true
  
  # Static block that should be transformed
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
  }
  
  # Dynamic block that should be preserved
  dynamic "client_ip_ranges" {
    for_each = var.ip_ranges
    content {
      start_ip = client_ip_ranges.value.start
      end_ip = client_ip_ranges.value.end
    }
  }
  
  # Static list that should be transformed
  cnode_ids = [1, 2, 3]
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check static attributes are preserved
        assert 'name = "test-resource"' in result
        assert 'enabled = true' in result
        
        # Check static block is transformed
        assert 'capacity_limits = {' in result
        assert 'soft_limit = 1000' in result
        
        # Check dynamic block is preserved
        assert 'dynamic "client_ip_ranges" {' in result
        assert 'for_each = var.ip_ranges' in result
        
        # Check static list remains as list (cnode_ids should not be converted to string)
        assert 'cnode_ids = [1, 2, 3]' in result
    
    def test_dynamic_block_with_comments(self):
        """Test that dynamic blocks with comments are preserved."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  # Configure IP ranges dynamically based on environment
  dynamic "client_ip_ranges" {
    # Iterate over all defined IP ranges
    for_each = var.ip_ranges
    content {
      # Set the start IP from the range definition
      start_ip = client_ip_ranges.value.start
      # Set the end IP from the range definition
      end_ip = client_ip_ranges.value.end
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check that dynamic block structure is preserved
        assert 'dynamic "client_ip_ranges" {' in result
        assert 'for_each = var.ip_ranges' in result
        
        # Comments should be preserved within the dynamic block
        assert '# Iterate over all defined IP ranges' in result
        assert '# Set the start IP from the range definition' in result
    
    def test_dynamic_block_vs_static_block_recognition(self):
        """Test that dynamic blocks are distinguished from static blocks with same name."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  # This static block should be transformed
  client_ip_ranges {
    start_ip = "192.168.1.1" 
    end_ip = "192.168.1.100"
  }
  
  # This dynamic block should be preserved
  dynamic "client_ip_ranges" {
    for_each = var.additional_ranges
    content {
      start_ip = client_ip_ranges.value.start
      end_ip = client_ip_ranges.value.end
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check that static block is transformed to list format
        assert 'client_ip_ranges = [' in result
        assert '["192.168.1.1", "192.168.1.100"]' in result
        
        # Check that dynamic block is preserved unchanged
        assert 'dynamic "client_ip_ranges" {' in result
        assert 'for_each = var.additional_ranges' in result
    
    def test_dynamic_block_with_conditional_logic(self):
        """Test dynamic blocks with conditional logic are preserved."""
        terraform_content = '''resource "vastdata_example" "test" {
  name = "test-resource"
  
  dynamic "network_rule" {
    for_each = var.enable_network_rules ? var.network_rules : []
    content {
      protocol = network_rule.value.protocol
      port = network_rule.value.port
      source = network_rule.value.source != null ? network_rule.value.source : "0.0.0.0/0"
      action = network_rule.value.action
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check conditional expression is preserved
        assert 'for_each = var.enable_network_rules ? var.network_rules : []' in result
        assert 'source = network_rule.value.source != null ? network_rule.value.source : "0.0.0.0/0"' in result
    
    def test_file_with_mixed_dynamic_and_static_content(self, temp_dir):
        """Test complete file transformation with mixed dynamic and static content."""
        terraform_content = '''resource "vastdata_administators_managers" "manager1" {
  name = "test-manager"
  
  # Static block - should be transformed
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
  }
  
  # Dynamic block - should be preserved
  dynamic "client_ip_ranges" {
    for_each = var.ip_ranges
    content {
      start_ip = client_ip_ranges.value.start
      end_ip = client_ip_ranges.value.end
    }
  }
  
  # Static list - should be transformed
  cnode_ids = [1, 2, 3]
}

resource "vastdata_kafka_brokers" "broker1" {
  broker_id = 1
  
  # Another dynamic block
  dynamic "security_groups" {
    for_each = var.security_groups
    content {
      name = security_groups.value.name
      rules = security_groups.value.rules
    }
  }
}'''
        
        input_file = temp_dir / "mixed_content.tf"
        output_file = temp_dir / "mixed_content_converted.tf"
        
        with open(input_file, 'w') as f:
            f.write(terraform_content)
        
        transform_file(input_file, output_file)
        
        with open(output_file, 'r') as f:
            result = f.read()
        
        # Check resource renaming occurred
        assert "vastdata_administrator_manager" in result
        assert "vastdata_kafka_broker" in result
        
        # Check static transformations occurred
        assert 'capacity_limits = {' in result
        assert 'cnode_ids = [1, 2, 3]' in result  # cnode_ids should remain as list
        
        # Check dynamic blocks were preserved
        assert 'dynamic "client_ip_ranges" {' in result
        assert 'dynamic "security_groups" {' in result
        assert 'for_each = var.ip_ranges' in result
        assert 'for_each = var.security_groups' in result

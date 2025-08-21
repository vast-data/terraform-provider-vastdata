#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

"""
Tests for specific QA issues identified in the migration script.
These tests verify that all reported transformation issues are properly handled.
"""

import pytest
from pathlib import Path
import tempfile
import sys
sys.path.append(str(Path(__file__).parent.parent))
from migration_script import transform_file, transform_resource_block


class TestQAIssues:
    """Test class for specific QA issues reported."""

    def test_blockhost_resource_rename(self):
        """Test QA Issue #1: vastdata_blockhost → vastdata_block_host rename."""
        terraform_content = '''resource vastdata_blockhost blockhost1 {
    name = var.host_name
    nqn = var.host_nqn
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert 'resource vastdata_block_host blockhost1' in result
        assert 'vastdata_blockhost' not in result
        # Verify attributes are preserved
        assert 'name = var.host_name' in result
        assert 'nqn = var.host_nqn' in result

    def test_kafka_brokers_resource_rename(self):
        """Test QA Issue #2: vastdata_kafka_brokers → vastdata_kafka_broker rename."""
        terraform_content = '''resource vastdata_kafka_brokers broker1 {
    name = var.broker_name
    addresses {
        host = var.broker_host
        port = var.broker_port
    }
    tenant_id = vastdata_tenant.broker_tenant1.id
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert 'resource vastdata_kafka_broker broker1' in result
        assert 'vastdata_kafka_brokers' not in result
        # Verify attributes are preserved
        assert 'name = var.broker_name' in result
        assert 'tenant_id = vastdata_tenant.broker_tenant1.id' in result

    def test_kafka_addresses_block_to_list_of_maps(self):
        """Test QA Issue #2: addresses block → list of maps transformation."""
        terraform_content = '''resource vastdata_kafka_broker broker1 {
    name = var.broker_name
    addresses {
        host = "10.131.21.121"
        port = 31485
    }
    tenant_id = vastdata_tenant.broker_tenant1.id
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that addresses block is converted to list of maps
        assert 'addresses = [' in result
        assert 'host = "10.131.21.121"' in result
        assert 'port = 31485' in result
        # Verify the structure is a list of maps, not just attributes
        assert '{\n      host = "10.131.21.121"' in result
        assert '    },' in result

    def test_kafka_addresses_multiple_blocks(self):
        """Test addresses block with multiple entries."""
        terraform_content = '''resource vastdata_kafka_broker broker1 {
    name = var.broker_name
    addresses {
        host = "10.131.21.121"
        port = 31485
    }
    addresses {
        host = "10.131.21.122"
        port = 31486
    }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that both addresses are converted
        assert 'addresses = [' in result
        assert 'host = "10.131.21.121"' in result
        assert 'host = "10.131.21.122"' in result
        assert 'port = 31485' in result
        assert 'port = 31486' in result

    def test_client_ip_ranges_variable_quoting(self):
        """Test QA Issue #3: Variable references should not be quoted in client_ip_ranges."""
        terraform_content = '''resource vastdata_tenant broker_tenant1 {
    name = var.tenant_name
    client_ip_ranges {
        start_ip = var.tenant_client_start_ip
        end_ip = var.tenant_client_end_ip
    }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check transformation to list of lists
        assert 'client_ip_ranges = [' in result
        # Verify variables are NOT quoted (no extra quotes around var.*)
        assert '[var.tenant_client_start_ip, var.tenant_client_end_ip]' in result
        # Ensure no double quoting
        assert '""' not in result

    def test_client_ip_ranges_mixed_variables_and_literals(self):
        """Test client_ip_ranges with mix of variables and string literals."""
        terraform_content = '''resource vastdata_tenant tenant1 {
    name = "test-tenant"
    client_ip_ranges {
        start_ip = var.start_ip_var
        end_ip = "192.168.1.100"
    }
    client_ip_ranges {
        start_ip = "10.0.0.1"
        end_ip = var.end_ip_var
    }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Variables should not be quoted, literals should be quoted
        assert '[var.start_ip_var, "192.168.1.100"]' in result
        assert '["10.0.0.1", var.end_ip_var]' in result
        # Ensure no double quoting
        assert '""' not in result

    def test_quota_blocks_to_list_of_maps(self):
        """Test QA Issue #4: user_quotas and group_quotas blocks → list of maps."""
        terraform_content = '''resource vastdata_quota quota1 {
  name = var.quota_name
  default_email = "user@example.com"
  path = vastdata_view.view1.path
  soft_limit = 100000
  hard_limit = 100000
  is_user_quota = true
  default_user_quota {
    grace_period = var.quota_default_grace_period
    hard_limit = 2000
    soft_limit = 1000
    hard_limit_inodes = 20000000
  }
  user_quotas {
    grace_period = var.quota_user_grace_period
    hard_limit = 15000
    soft_limit = 15000
    entity {
      name = var.quota_user_name
      email = "user1@example.com"
      identifier = var.quota_user_name
      identifier_type = "username"
      is_group = "false"
    }
  }
  group_quotas {
    grace_period = var.quota_group_grace_period
    hard_limit = 15000
    soft_limit = 15000
    entity {
      name = vastdata_group.group3.name
      identifier = vastdata_group.group3.name
      identifier_type = "group"
      is_group = "false"
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        
        # Check that default_user_quota becomes attributes (single block)
        assert 'default_user_quota = {' in result
        assert 'grace_period = var.quota_default_grace_period' in result
        
        # Check that user_quotas becomes list of maps
        assert 'user_quotas = [' in result
        assert 'grace_period = var.quota_user_grace_period' in result
        assert 'name = var.quota_user_name' in result
        assert 'email = "user1@example.com"' in result
        
        # Check that group_quotas becomes list of maps
        assert 'group_quotas = [' in result
        assert 'grace_period = var.quota_group_grace_period' in result
        assert 'name = vastdata_group.group3.name' in result
        assert 'identifier_type = "group"' in result

    def test_s3_lifecycle_rule_rename(self):
        """Test QA Issue #5: vastdata_s3_life_cycle_rule → vastdata_s3_lifecycle_rule rename."""
        terraform_content = '''resource vastdata_s3_life_cycle_rule s3_lifecycle_rule1 {
  name = var.s3_lifecycle_name
  max_size = 10000000
  min_size = 100000
  newer_noncurrent_versions = 3
  prefix = "prefix"
  view_id = vastdata_view.s3_view1.id
  expiration_days = 30
  enabled = var.s3_lifecycle_enabled
  noncurrent_days = 20
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        assert 'resource vastdata_s3_lifecycle_rule s3_lifecycle_rule1' in result
        assert 'vastdata_s3_life_cycle_rule' not in result
        # Verify attributes are preserved
        assert 'name = var.s3_lifecycle_name' in result
        assert 'max_size = 10000000' in result
        assert 'enabled = var.s3_lifecycle_enabled' in result

    def test_comprehensive_qa_issues_file_transformation(self, tmp_path):
        """Test all QA issues in a complete file transformation."""
        terraform_content = '''# Test file covering all QA issues
variable "host_name" {
  type = string
}

variable "tenant_client_start_ip" {
  type = string
}

variable "tenant_client_end_ip" {
  type = string
}

# QA Issue #1: blockhost rename
resource vastdata_blockhost blockhost1 {
    name = var.host_name
    nqn = var.host_nqn
}

# QA Issue #2: kafka_brokers rename and addresses transformation
resource vastdata_kafka_brokers broker1 {
    name = var.broker_name
    addresses {
        host = "10.131.21.121"
        port = 31485
    }
    tenant_id = vastdata_tenant.broker_tenant1.id
}

# QA Issue #3: client_ip_ranges variable quoting
resource vastdata_tenant broker_tenant1 {
    name = var.tenant_name
    client_ip_ranges {
        start_ip = var.tenant_client_start_ip
        end_ip = var.tenant_client_end_ip
    }
}

# QA Issue #4: quota blocks transformation
resource vastdata_quota quota1 {
  name = var.quota_name
  default_user_quota {
    grace_period = var.quota_default_grace_period
    hard_limit = 2000
  }
  user_quotas {
    grace_period = var.quota_user_grace_period
    hard_limit = 15000
    entity {
      name = var.quota_user_name
      identifier_type = "username"
    }
  }
  group_quotas {
    grace_period = var.quota_group_grace_period
    hard_limit = 15000
    entity {
      name = vastdata_group.group3.name
      identifier_type = "group"
    }
  }
}

# QA Issue #5: s3_life_cycle_rule rename
resource vastdata_s3_life_cycle_rule s3_lifecycle_rule1 {
  name = var.s3_lifecycle_name
  enabled = var.s3_lifecycle_enabled
}'''

        input_file = tmp_path / "qa_issues_test.tf"
        output_file = tmp_path / "qa_issues_test_converted.tf"
        
        input_file.write_text(terraform_content)
        
        # Run transformation
        transform_file(input_file, output_file)
        
        result = output_file.read_text()
        
        # Verify all QA issues are resolved
        
        # Issue #1: blockhost rename
        assert 'resource vastdata_block_host blockhost1' in result
        assert 'vastdata_blockhost' not in result
        
        # Issue #2: kafka_brokers rename and addresses transformation
        assert 'resource vastdata_kafka_broker broker1' in result
        assert 'vastdata_kafka_brokers' not in result
        assert 'addresses = [' in result
        assert 'host = "10.131.21.121"' in result
        assert 'port = 31485' in result
        
        # Issue #3: client_ip_ranges variable quoting
        assert 'client_ip_ranges = [' in result
        assert '[var.tenant_client_start_ip, var.tenant_client_end_ip]' in result
        assert '""' not in result
        
        # Issue #4: quota blocks transformation
        assert 'default_user_quota = {' in result
        assert 'user_quotas = [' in result
        assert 'group_quotas = [' in result
        
        # Issue #5: s3_life_cycle_rule rename
        assert 'resource vastdata_s3_lifecycle_rule s3_lifecycle_rule1' in result
        assert 'vastdata_s3_life_cycle_rule' not in result

    def test_nested_entity_blocks_in_quotas(self):
        """Test that nested entity blocks within quotas are properly flattened."""
        terraform_content = '''resource vastdata_quota quota1 {
  name = "test-quota"
  user_quotas {
    grace_period = 7
    hard_limit = 15000
    soft_limit = 15000
    entity {
      name = "testuser"
      email = "user1@example.com"
      identifier = "testuser"
      identifier_type = "username"
      is_group = "false"
    }
  }
}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None
        # Check that user_quotas becomes list of maps with flattened entity attributes
        assert 'user_quotas = [' in result
        assert 'grace_period = 7' in result
        assert 'hard_limit = 15000' in result
        assert 'name = "testuser"' in result
        assert 'email = "user1@example.com"' in result
        assert 'identifier = "testuser"' in result
        assert 'identifier_type = "username"' in result
        assert 'is_group = "false"' in result
        
        # Ensure entity block structure is flattened (no nested entity block)
        entity_block_pattern = 'entity {'
        assert entity_block_pattern not in result

    def test_resource_reference_updates_for_renamed_types(self, tmp_path):
        """Test that resource references are updated when resource types are renamed."""
        terraform_content = '''# Resources that will be renamed
resource vastdata_blockhost host1 {
    name = "test-host"
}

resource vastdata_kafka_brokers broker1 {
    name = "test-broker"
}

resource vastdata_s3_life_cycle_rule rule1 {
    name = "test-rule"
}

# Resources that reference the renamed resources
resource vastdata_example example1 {
    host_ref = vastdata_blockhost.host1.id
    broker_ref = vastdata_kafka_brokers.broker1.name
    rule_ref = vastdata_s3_life_cycle_rule.rule1.id
}

# Output that references renamed resources
output "host_name" {
    value = vastdata_blockhost.host1.name
}

output "broker_id" {
    value = vastdata_kafka_brokers.broker1.id
}

# Data source that references renamed type
data "vastdata_blockhost" "existing_host" {
    name = "existing"
}'''

        input_file = tmp_path / "reference_test.tf"
        output_file = tmp_path / "reference_test_converted.tf"
        
        input_file.write_text(terraform_content)
        
        # Run transformation
        transform_file(input_file, output_file)
        
        result = output_file.read_text()
        
        # Verify resource type renames
        assert 'resource vastdata_block_host host1' in result
        assert 'resource vastdata_kafka_broker broker1' in result
        assert 'resource vastdata_s3_lifecycle_rule rule1' in result
        
        # Verify references are updated
        assert 'host_ref = vastdata_block_host.host1.id' in result
        assert 'broker_ref = vastdata_kafka_broker.broker1.name' in result
        assert 'rule_ref = vastdata_s3_lifecycle_rule.rule1.id' in result
        
        # Verify output references are updated
        assert 'value = vastdata_block_host.host1.name' in result
        assert 'value = vastdata_kafka_broker.broker1.id' in result
        
        # Verify data source type is updated
        assert 'data "vastdata_block_host" "existing_host"' in result
        
        # Ensure old names are completely removed
        assert 'vastdata_blockhost' not in result
        assert 'vastdata_kafka_brokers' not in result
        assert 'vastdata_s3_life_cycle_rule' not in result

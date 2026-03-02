#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Comprehensive tests for v1 legacy resource name conversions.
Tests configuration migration (migration_script.py).
"""

import pytest
import json
import tempfile
import os
import sys
from pathlib import Path

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import (
    transform_resource_block,
    resource_type_rename_map,
)


# Test data: v1 resource name -> v3 resource name
V1_TO_V3_RESOURCE_MAPPINGS = {
    # Plural forms (administrators with typo)
    "vastdata_administators_managers": "vastdata_administrator_manager",
    "vastdata_administators_roles": "vastdata_administrator_role",
    "vastdata_administators_realms": "vastdata_administrator_realm",
    # Plural forms (other resources)
    "vastdata_kafka_brokers": "vastdata_kafka_broker",
    "vastdata_replication_peers": "vastdata_replication_peer",
    "vastdata_s3_replication_peers": "vastdata_s3_replication_peer",
    # Old naming
    "vastdata_active_directory2": "vastdata_active_directory",
    "vastdata_non_local_user": "vastdata_nonlocal_user",
    "vastdata_non_local_user_key": "vastdata_nonlocal_user_key",
    "vastdata_non_local_group": "vastdata_nonlocal_group",
    "vastdata_saml": "vastdata_saml_config",
    "vastdata_blockhost": "vastdata_block_host",
    # Naming variant (underscores)
    "vastdata_s3_lifecycle_rule": "vastdata_s3_life_cycle_rule",
}


class TestConfigurationMigration:
    """Tests for migration_script.py (configuration file migration)"""
    
    def test_all_v1_names_in_rename_map(self):
        """Verify all v1 legacy names are in the configuration migration rename map"""
        for v1_name in V1_TO_V3_RESOURCE_MAPPINGS.keys():
            assert v1_name in resource_type_rename_map, \
                f"Configuration migration missing v1 resource: {v1_name}"
    
    def test_v1_to_v3_mapping_correctness(self):
        """Verify all v1->v3 mappings are correct in configuration migration"""
        for v1_name, expected_v3_name in V1_TO_V3_RESOURCE_MAPPINGS.items():
            actual_v3_name = resource_type_rename_map.get(v1_name)
            assert actual_v3_name == expected_v3_name, \
                f"Configuration migration: {v1_name} should map to {expected_v3_name}, got {actual_v3_name}"
    
    @pytest.mark.parametrize("v1_name,v3_name", list(V1_TO_V3_RESOURCE_MAPPINGS.items()))
    def test_individual_resource_conversion(self, v1_name, v3_name):
        """Test each v1 resource type is correctly converted to v3 in .tf files"""
        terraform_content = f'''resource "{v1_name}" "test_resource" {{
  name = "test-{v1_name}"
  enabled = true
}}'''
        
        lines = terraform_content.split('\n')
        result, consumed = transform_resource_block(lines, 0)
        
        assert result is not None, f"Failed to transform {v1_name}"
        assert v3_name in result, f"v3 name {v3_name} not found in result for {v1_name}"
        assert f'resource "{v1_name}"' not in result, f"v1 name {v1_name} still present in result"


class TestEndToEndConversion:
    """End-to-end tests with multiple v1 resources"""
    
    def test_mixed_v1_configuration_file(self):
        """Test .tf file with multiple v1 resources"""
        terraform_content = '''
resource "vastdata_s3_replication_peers" "peer1" {
  name = "test-peer"
}

resource "vastdata_administators_managers" "admin1" {
  name = "test-admin"
}

resource "vastdata_active_directory2" "ad1" {
  domain = "example.com"
}
'''
        
        lines = terraform_content.split('\n')
        converted_blocks = []
        idx = 0
        
        while idx < len(lines):
            if 'resource "vastdata_' in lines[idx]:
                result, consumed = transform_resource_block(lines, idx)
                if result:
                    converted_blocks.append(result)
                idx += consumed
            else:
                idx += 1
        
        # Should convert all 3 resources
        assert len(converted_blocks) == 3
        
        full_result = '\n'.join(converted_blocks)
        
        # Verify all v3 names are present
        assert 'vastdata_s3_replication_peer' in full_result
        assert 'vastdata_administrator_manager' in full_result
        assert 'vastdata_active_directory' in full_result
        
        # Verify no v1 names remain
        assert 'vastdata_s3_replication_peers' not in full_result
        assert 'vastdata_administators_managers' not in full_result
        assert 'vastdata_active_directory2' not in full_result


if __name__ == '__main__':
    pytest.main([__file__, '-v', '--tb=short'])

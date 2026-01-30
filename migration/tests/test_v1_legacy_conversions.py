#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Comprehensive tests for v1 legacy resource name conversions.
Tests both configuration migration (migration_script.py) and state migration (state_migration.py).
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

from state_migration import (
    extract_resources,
    build_import_id,
    generate_import_script,
    RESOURCE_IMPORT_MAP,
    RESOURCE_NAME_TRANSLATION,
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


class TestStateMigration:
    """Tests for state_migration.py (state file migration)"""
    
    def test_all_v1_names_in_import_map(self):
        """Verify all v1 legacy names are in the state migration import map"""
        for v1_name in V1_TO_V3_RESOURCE_MAPPINGS.keys():
            assert v1_name in RESOURCE_IMPORT_MAP, \
                f"State migration RESOURCE_IMPORT_MAP missing v1 resource: {v1_name}"
    
    def test_all_v1_names_in_translation_map(self):
        """Verify all v1 legacy names are in the translation map"""
        for v1_name in V1_TO_V3_RESOURCE_MAPPINGS.keys():
            assert v1_name in RESOURCE_NAME_TRANSLATION, \
                f"State migration RESOURCE_NAME_TRANSLATION missing v1 resource: {v1_name}"
    
    def test_translation_map_correctness(self):
        """Verify all translations in RESOURCE_NAME_TRANSLATION are correct"""
        for v1_name, expected_v3_name in V1_TO_V3_RESOURCE_MAPPINGS.items():
            actual_v3_name = RESOURCE_NAME_TRANSLATION.get(v1_name)
            assert actual_v3_name == expected_v3_name, \
                f"Translation map: {v1_name} should map to {expected_v3_name}, got {actual_v3_name}"
    
    @pytest.mark.parametrize("v1_name,v3_name", list(V1_TO_V3_RESOURCE_MAPPINGS.items()))
    def test_state_extraction_for_each_v1_resource(self, v1_name, v3_name):
        """Test each v1 resource type can be extracted from state"""
        state = {
            "version": 4,
            "terraform_version": "1.5.7",
            "resources": [
                {
                    "mode": "managed",
                    "type": v1_name,
                    "name": "test_resource",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [
                        {
                            "schema_version": 0,
                            "attributes": {
                                "id": 123,
                                "name": f"test-{v1_name}"
                            }
                        }
                    ]
                }
            ]
        }
        
        resources = extract_resources(state)
        
        # Should extract the resource
        assert len(resources) == 1, f"Failed to extract {v1_name} from state"
        assert resources[0]['type'] == v1_name
        assert resources[0]['attributes']['id'] == 123
    
    @pytest.mark.parametrize("v1_name,v3_name", list(V1_TO_V3_RESOURCE_MAPPINGS.items()))
    def test_import_id_generation_for_each_v1_resource(self, v1_name, v3_name):
        """Test import ID can be generated for each v1 resource"""
        # Skip non_local resources as they have composite keys
        if "non_local" in v1_name:
            pytest.skip("Composite key resources tested separately")
        
        resource = {
            'type': v1_name,
            'name': 'test_resource',
            'address': f'{v1_name}.test_resource',
            'attributes': {'id': 123}
        }
        
        import_id = build_import_id(resource)
        
        # Should have import configuration
        assert v1_name in RESOURCE_IMPORT_MAP, f"No import config for {v1_name}"
        # Should generate import ID
        assert import_id == '123', f"Failed to generate import ID for {v1_name}"
    
    def test_non_local_resources_composite_keys(self):
        """Test non_local resources with composite keys"""
        non_local_resources = {
            "vastdata_non_local_user": {
                "username": "testuser",
                "context": "LOCAL",
                "tenant_id": 1
            },
            "vastdata_non_local_group": {
                "groupname": "testgroup",
                "context": "LOCAL",
                "tenant_id": 1
            }
        }
        
        for resource_type, attrs in non_local_resources.items():
            resource = {
                'type': resource_type,
                'name': 'test',
                'address': f'{resource_type}.test',
                'attributes': attrs
            }
            
            import_id = build_import_id(resource)
            assert import_id is not None, f"Failed to generate import ID for {resource_type}"
            # For non_local_user, should be: username=X,context=Y,tenant_id=Z
            if "user" in resource_type and "key" not in resource_type:
                assert import_id == "username=testuser,context=LOCAL,tenant_id=1"
            # For non_local_group, should be: groupname=X,context=Y,tenant_id=Z
            elif "group" in resource_type:
                assert import_id == "groupname=testgroup,context=LOCAL,tenant_id=1"


class TestStubAndImportScriptGeneration:
    """Tests for generated stub files and import scripts"""
    
    @pytest.mark.parametrize("v1_name,v3_name", list(V1_TO_V3_RESOURCE_MAPPINGS.items()))
    def test_stub_file_uses_v3_names(self, v1_name, v3_name):
        """Test that generated stub files use v3 resource names, not v1"""
        # Skip non_local resources for simplicity
        if "non_local" in v1_name:
            pytest.skip("Composite key resources tested separately")
        
        with tempfile.TemporaryDirectory() as temp_dir:
            output_dir = os.path.join(temp_dir, "output")
            terraform_dir = os.path.join(temp_dir, "terraform")
            os.makedirs(output_dir)
            os.makedirs(terraform_dir)
            
            resources = [{
                'type': v1_name,
                'name': 'test_resource',
                'address': f'{v1_name}.test_resource',
                'attributes': {'id': '123'}
            }]
            
            generate_import_script(resources, output_dir, terraform_dir)
            
            # Check stub file
            resources_tf_path = os.path.join(terraform_dir, "resources.tf")
            with open(resources_tf_path, 'r') as f:
                content = f.read()
            
            # Should use v3 name in stub
            assert f'resource "{v3_name}"' in content, \
                f"Stub file should use v3 name {v3_name} for {v1_name}"
            # Should NOT use v1 name in stub
            assert f'resource "{v1_name}"' not in content, \
                f"Stub file should not use v1 name {v1_name}"
    
    @pytest.mark.parametrize("v1_name,v3_name", list(V1_TO_V3_RESOURCE_MAPPINGS.items()))
    def test_import_script_uses_v3_names(self, v1_name, v3_name):
        """Test that import commands use v3 resource names, not v1"""
        # Skip non_local resources for simplicity
        if "non_local" in v1_name:
            pytest.skip("Composite key resources tested separately")
        
        with tempfile.TemporaryDirectory() as temp_dir:
            output_dir = os.path.join(temp_dir, "output")
            terraform_dir = os.path.join(temp_dir, "terraform")
            os.makedirs(output_dir)
            os.makedirs(terraform_dir)
            
            resources = [{
                'type': v1_name,
                'name': 'test_resource',
                'address': f'{v1_name}.test_resource',
                'attributes': {'id': '123'}
            }]
            
            script_path = generate_import_script(resources, output_dir, terraform_dir)
            
            # Check import script
            with open(script_path, 'r') as f:
                content = f.read()
            
            # Should use v3 name in import command
            assert f"terraform import '{v3_name}.test_resource'" in content, \
                f"Import script should use v3 name {v3_name} for {v1_name}"
            # Should NOT use v1 name in import command
            assert f"terraform import '{v1_name}.test_resource'" not in content, \
                f"Import script should not use v1 name {v1_name}"


class TestEndToEndConversion:
    """End-to-end tests with multiple v1 resources"""
    
    def test_mixed_v1_resources_in_state(self):
        """Test state with multiple different v1 resources"""
        state = {
            "version": 4,
            "terraform_version": "1.5.7",
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_s3_replication_peers",
                    "name": "peer1",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [{"schema_version": 0, "attributes": {"id": 1}}]
                },
                {
                    "mode": "managed",
                    "type": "vastdata_kafka_brokers",
                    "name": "broker1",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [{"schema_version": 0, "attributes": {"id": 2}}]
                },
                {
                    "mode": "managed",
                    "type": "vastdata_blockhost",
                    "name": "host1",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [{"schema_version": 0, "attributes": {"id": 3}}]
                },
            ]
        }
        
        resources = extract_resources(state)
        
        # Should extract all 3 resources
        assert len(resources) == 3
        
        # All should have import configuration
        for resource in resources:
            assert resource['type'] in RESOURCE_IMPORT_MAP
            assert resource['type'] in RESOURCE_NAME_TRANSLATION
    
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

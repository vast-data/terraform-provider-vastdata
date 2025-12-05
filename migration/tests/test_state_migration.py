#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Tests for state_migration.py
"""

import pytest
import json
import tempfile
import os
import sys
from pathlib import Path
import shutil

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from state_migration import (
    parse_tfstate,
    extract_resources,
    build_import_id,
    RESOURCE_IMPORT_MAP
)


class TestStateParsing:
    """Test parsing of Terraform state files"""
    
    @pytest.fixture
    def sample_state_v4(self, tmp_path):
        """Create a sample TF state file (version 4, TF >= 0.12)"""
        state = {
            "version": 4,
            "terraform_version": "1.0.0",
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
                                "id": 1,
                                "name": "test-tenant",
                                "encryption_crn": "crn:v1:vastdata:encryption:tenant:1"
                            }
                        }
                    ]
                },
                {
                    "mode": "managed",
                    "type": "vastdata_view",
                    "name": "test_view",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [
                        {
                            "schema_version": 0,
                            "attributes": {
                                "id": 10,
                                "path": "/test-view",
                                "policy_id": 2,
                                "tenant_id": 1
                            }
                        }
                    ]
                },
                {
                    "mode": "data",
                    "type": "vastdata_vip_pool",
                    "name": "default_pool",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [
                        {
                            "schema_version": 0,
                            "attributes": {
                                "id": 1,
                                "name": "default"
                            }
                        }
                    ]
                }
            ]
        }
        
        state_file = tmp_path / "terraform.tfstate"
        with open(state_file, 'w') as f:
            json.dump(state, f)
        
        return state_file
    
    @pytest.fixture
    def sample_state_with_composite_keys(self, tmp_path):
        """Create a state file with resources that need composite import keys"""
        state = {
            "version": 4,
            "terraform_version": "1.0.0",
            "serial": 1,
            "lineage": "test-lineage",
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_nonlocal_user",
                    "name": "ldap_user",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [
                        {
                            "schema_version": 0,
                            "attributes": {
                                "username": "jdoe",
                                "context": "ldap",
                                "tenant_id": 5,
                                "uid": 1001
                            }
                        }
                    ]
                },
                {
                    "mode": "managed",
                    "type": "vastdata_nonlocal_group",
                    "name": "ldap_group",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [
                        {
                            "schema_version": 0,
                            "attributes": {
                                "groupname": "developers",
                                "context": "ldap",
                                "tenant_id": 5,
                                "gid": 2001
                            }
                        }
                    ]
                }
            ]
        }
        
        state_file = tmp_path / "terraform.tfstate"
        with open(state_file, 'w') as f:
            json.dump(state, f)
        
        return state_file
    
    def test_parse_valid_state(self, sample_state_v4):
        """Test parsing a valid state file"""
        state = parse_tfstate(str(sample_state_v4))
        
        assert state is not None
        assert state['version'] == 4
        assert 'resources' in state
        assert len(state['resources']) == 3
    
    def test_parse_invalid_json(self, tmp_path):
        """Test parsing an invalid JSON file"""
        invalid_file = tmp_path / "invalid.tfstate"
        with open(invalid_file, 'w') as f:
            f.write("{ invalid json }")
        
        with pytest.raises(SystemExit):
            parse_tfstate(str(invalid_file))
    
    def test_parse_missing_file(self):
        """Test parsing a non-existent file"""
        with pytest.raises(SystemExit):
            parse_tfstate("/non/existent/file.tfstate")


class TestResourceExtraction:
    """Test extraction of resources from state"""
    
    @pytest.fixture
    def parsed_state(self):
        """Sample parsed state"""
        return {
            "version": 4,
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "tenant1",
                    "instances": [
                        {
                            "attributes": {
                                "id": 1,
                                "name": "tenant1"
                            }
                        }
                    ]
                },
                {
                    "mode": "managed",
                    "type": "vastdata_view",
                    "name": "view1",
                    "instances": [
                        {
                            "attributes": {
                                "id": 10,
                                "path": "/view1"
                            }
                        }
                    ]
                },
                {
                    "mode": "data",
                    "type": "vastdata_vip_pool",
                    "name": "pool1",
                    "instances": [
                        {
                            "attributes": {
                                "id": 1
                            }
                        }
                    ]
                }
            ]
        }
    
    def test_extract_resources_basic(self, parsed_state):
        """Test basic resource extraction"""
        resources = extract_resources(parsed_state)
        
        # Should extract 2 managed resources, skip data source
        assert len(resources) == 2
        
        # Check first resource
        assert resources[0]['type'] == 'vastdata_tenant'
        assert resources[0]['name'] == 'tenant1'
        assert resources[0]['address'] == 'vastdata_tenant.tenant1'
        assert resources[0]['attributes']['id'] == 1
        
        # Check second resource
        assert resources[1]['type'] == 'vastdata_view'
        assert resources[1]['name'] == 'view1'
        assert resources[1]['address'] == 'vastdata_view.view1'
    
    def test_extract_resources_with_module(self):
        """Test extraction of resources in modules"""
        state = {
            "version": 4,
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "tenant1",
                    "module": "module.networking",
                    "instances": [
                        {
                            "attributes": {
                                "id": 1
                            }
                        }
                    ]
                }
            ]
        }
        
        resources = extract_resources(state)
        assert len(resources) == 1
        assert resources[0]['address'] == 'module.networking.vastdata_tenant.tenant1'
        assert resources[0]['module'] == 'module.networking'
    
    def test_extract_resources_with_count(self):
        """Test extraction of resources with count"""
        state = {
            "version": 4,
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "tenant",
                    "instances": [
                        {
                            "index_key": 0,
                            "attributes": {
                                "id": 1
                            }
                        },
                        {
                            "index_key": 1,
                            "attributes": {
                                "id": 2
                            }
                        }
                    ]
                }
            ]
        }
        
        resources = extract_resources(state)
        assert len(resources) == 2
        assert resources[0]['address'] == 'vastdata_tenant.tenant[0]'
        assert resources[1]['address'] == 'vastdata_tenant.tenant[1]'
    
    def test_extract_resources_with_for_each(self):
        """Test extraction of resources with for_each"""
        state = {
            "version": 4,
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "tenants",
                    "instances": [
                        {
                            "index_key": "prod",
                            "attributes": {
                                "id": 1,
                                "name": "prod-tenant"
                            }
                        },
                        {
                            "index_key": "dev",
                            "attributes": {
                                "id": 2,
                                "name": "dev-tenant"
                            }
                        }
                    ]
                }
            ]
        }
        
        resources = extract_resources(state)
        assert len(resources) == 2
        assert resources[0]['address'] == 'vastdata_tenant.tenants["prod"]'
        assert resources[1]['address'] == 'vastdata_tenant.tenants["dev"]'


class TestImportIDGeneration:
    """Test generation of import IDs"""
    
    def test_simple_id_import(self):
        """Test simple ID import (single field)"""
        resource = {
            'type': 'vastdata_tenant',
            'name': 'tenant1',
            'address': 'vastdata_tenant.tenant1',
            'attributes': {
                'id': 5,
                'name': 'my-tenant'
            }
        }
        
        import_id = build_import_id(resource)
        assert import_id == '5'
    
    def test_composite_key_import(self):
        """Test composite key import"""
        resource = {
            'type': 'vastdata_nonlocal_user',
            'name': 'ldap_user',
            'address': 'vastdata_nonlocal_user.ldap_user',
            'attributes': {
                'username': 'jdoe',
                'context': 'ldap',
                'tenant_id': 10,
                'uid': 1001
            }
        }
        
        import_id = build_import_id(resource)
        assert import_id == 'username=jdoe,context=ldap,tenant_id=10'
    
    def test_missing_required_field(self):
        """Test handling of missing required field"""
        resource = {
            'type': 'vastdata_tenant',
            'name': 'tenant1',
            'address': 'vastdata_tenant.tenant1',
            'attributes': {
                'name': 'my-tenant'
                # Missing 'id' field
            }
        }
        
        import_id = build_import_id(resource)
        assert import_id is None
    
    def test_unknown_resource_type(self):
        """Test handling of unknown resource type"""
        resource = {
            'type': 'vastdata_unknown_resource',
            'name': 'unknown',
            'address': 'vastdata_unknown_resource.unknown',
            'attributes': {
                'id': 1
            }
        }
        
        import_id = build_import_id(resource)
        assert import_id is None
    
    def test_view_simple_id_import(self):
        """Test that view now uses simple ID import (not composite)"""
        resource = {
            'type': 'vastdata_view',
            'name': 'test_view',
            'address': 'vastdata_view.test_view',
            'attributes': {
                'id': 4,
                'path': '/test-view',
                'tenant_name': None,  # May be null
                'tenant_id': 1
            }
        }
        
        import_id = build_import_id(resource)
        assert import_id == '4'  # Should use simple ID, not composite
    
    def test_view_policy_simple_id_import(self):
        """Test that view_policy now uses simple ID import"""
        resource = {
            'type': 'vastdata_view_policy',
            'name': 'policy',
            'address': 'vastdata_view_policy.policy',
            'attributes': {
                'id': 6,
                'name': 'test-policy',
                'tenant_name': None
            }
        }
        
        import_id = build_import_id(resource)
        assert import_id == '6'


class TestResourceImportMapping:
    """Test the resource import mapping configuration"""
    
    def test_all_simple_id_resources(self):
        """Test that common resources use simple ID import"""
        simple_id_resources = [
            'vastdata_tenant',
            'vastdata_vip_pool',
            'vastdata_view',
            'vastdata_view_policy',
            'vastdata_user',
            'vastdata_group',
            'vastdata_quota',
            'vastdata_protection_policy',
            'vastdata_protected_path',
        ]
        
        for resource_type in simple_id_resources:
            assert resource_type in RESOURCE_IMPORT_MAP
            import_fields, id_field = RESOURCE_IMPORT_MAP[resource_type]
            assert import_fields == ['id']
            assert id_field == 'id'
    
    def test_composite_key_resources(self):
        """Test that composite key resources are configured correctly"""
        # nonlocal_user requires username, context, tenant_id
        config = RESOURCE_IMPORT_MAP.get('vastdata_nonlocal_user')
        assert config is not None
        import_fields, id_field = config
        assert 'username' in import_fields
        assert 'context' in import_fields
        assert 'tenant_id' in import_fields
        assert id_field is None  # No single ID field
        
        # nonlocal_group requires groupname, context, tenant_id
        config = RESOURCE_IMPORT_MAP.get('vastdata_nonlocal_group')
        assert config is not None
        import_fields, id_field = config
        assert 'groupname' in import_fields
        assert 'context' in import_fields
        assert 'tenant_id' in import_fields


class TestEndToEnd:
    """End-to-end integration tests"""
    
    @pytest.fixture
    def temp_workspace(self, tmp_path):
        """Create a temporary workspace with TF files and state"""
        workspace = tmp_path / "workspace"
        workspace.mkdir()
        
        # Create main.tf
        main_tf = workspace / "main.tf"
        main_tf.write_text('''
terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
      version = "2.0.0"
    }
  }
}

provider "vastdata" {
  host            = "10.0.0.1"
  username        = "admin"
  password        = "password"
  skip_ssl_verify = true
}
''')
        
        # Create terraform.tfstate
        state = {
            "version": 4,
            "terraform_version": "1.0.0",
            "serial": 1,
            "lineage": "test",
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
                                "id": 1,
                                "name": "test-tenant"
                            }
                        }
                    ]
                }
            ]
        }
        
        state_file = workspace / "terraform.tfstate"
        with open(state_file, 'w') as f:
            json.dump(state, f)
        
        return workspace
    
    def test_parse_and_extract_workflow(self, temp_workspace):
        """Test the complete workflow: parse -> extract -> build import IDs"""
        state_file = temp_workspace / "terraform.tfstate"
        
        # Parse state
        state = parse_tfstate(str(state_file))
        assert state is not None
        
        # Extract resources
        resources = extract_resources(state)
        assert len(resources) == 1
        assert resources[0]['type'] == 'vastdata_tenant'
        
        # Build import ID
        import_id = build_import_id(resources[0])
        assert import_id == '1'


if __name__ == '__main__':
    pytest.main([__file__, '-v'])


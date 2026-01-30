#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Tests for parallel import and mixed provider state migration features
"""

import pytest
import json
import tempfile
import os
import sys
from pathlib import Path
import shutil
from unittest.mock import patch, MagicMock

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from state_migration import (
    separate_vast_resources,
    extract_state_resources_raw,
    merge_non_vast_resources_to_state,
    import_single_resource,
    run_parallel_imports,
    extract_resources
)


class TestSeparateVastResources:
    """Test separation of VAST and non-VAST resources"""
    
    def test_separate_all_vast(self):
        """Test with only VAST resources"""
        resources = [
            {'type': 'vastdata_tenant', 'address': 'vastdata_tenant.test1', 'attributes': {}},
            {'type': 'vastdata_view', 'address': 'vastdata_view.test2', 'attributes': {}},
        ]
        
        vast, non_vast = separate_vast_resources(resources)
        
        assert len(vast) == 2
        assert len(non_vast) == 0
        assert all(r['type'].startswith('vastdata_') for r in vast)
    
    def test_separate_all_non_vast(self):
        """Test with only non-VAST resources"""
        resources = [
            {'type': 'aws_s3_bucket', 'address': 'aws_s3_bucket.test1', 'attributes': {}},
            {'type': 'google_storage_bucket', 'address': 'google_storage_bucket.test2', 'attributes': {}},
        ]
        
        vast, non_vast = separate_vast_resources(resources)
        
        assert len(vast) == 0
        assert len(non_vast) == 2
        assert all(not r['type'].startswith('vastdata_') for r in non_vast)
    
    def test_separate_mixed_providers(self):
        """Test with mixed VAST and non-VAST resources"""
        resources = [
            {'type': 'vastdata_tenant', 'address': 'vastdata_tenant.test1', 'attributes': {}},
            {'type': 'aws_s3_bucket', 'address': 'aws_s3_bucket.test1', 'attributes': {}},
            {'type': 'vastdata_view', 'address': 'vastdata_view.test2', 'attributes': {}},
            {'type': 'google_compute_instance', 'address': 'google_compute_instance.vm1', 'attributes': {}},
            {'type': 'vastdata_vip_pool', 'address': 'vastdata_vip_pool.pool1', 'attributes': {}},
        ]
        
        vast, non_vast = separate_vast_resources(resources)
        
        assert len(vast) == 3
        assert len(non_vast) == 2
        assert all(r['type'].startswith('vastdata_') for r in vast)
        assert all(not r['type'].startswith('vastdata_') for r in non_vast)
        
        # Verify specific resources
        vast_types = [r['type'] for r in vast]
        assert 'vastdata_tenant' in vast_types
        assert 'vastdata_view' in vast_types
        assert 'vastdata_vip_pool' in vast_types
        
        non_vast_types = [r['type'] for r in non_vast]
        assert 'aws_s3_bucket' in non_vast_types
        assert 'google_compute_instance' in non_vast_types
    
    def test_separate_empty_list(self):
        """Test with empty resource list"""
        resources = []
        
        vast, non_vast = separate_vast_resources(resources)
        
        assert len(vast) == 0
        assert len(non_vast) == 0


class TestExtractStateResourcesRaw:
    """Test extraction of raw state resources for preservation"""
    
    def test_extract_raw_resources_modern_format(self):
        """Test with TF >= 0.12 format"""
        state = {
            "version": 4,
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "test",
                    "instances": [{"attributes": {"id": 1}}]
                },
                {
                    "mode": "managed",
                    "type": "aws_s3_bucket",
                    "name": "bucket",
                    "instances": [{"attributes": {"id": "my-bucket"}}]
                }
            ]
        }
        
        raw_resources = extract_state_resources_raw(state)
        
        assert len(raw_resources) == 2
        assert raw_resources[0]['type'] == 'vastdata_tenant'
        assert raw_resources[1]['type'] == 'aws_s3_bucket'
    
    def test_extract_raw_resources_with_data_sources(self):
        """Test that data sources are included in raw extraction"""
        state = {
            "version": 4,
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "test",
                    "instances": [{"attributes": {"id": 1}}]
                },
                {
                    "mode": "data",
                    "type": "aws_ami",
                    "name": "ubuntu",
                    "instances": [{"attributes": {"id": "ami-12345"}}]
                }
            ]
        }
        
        raw_resources = extract_state_resources_raw(state)
        
        # Raw extraction should include both managed and data resources
        assert len(raw_resources) == 2


class TestMergeNonVastResources:
    """Test merging non-VAST resources into migrated state"""
    
    @pytest.fixture
    def temp_state_file(self, tmp_path):
        """Create a temporary migrated state file"""
        state = {
            "version": 4,
            "terraform_version": "1.5.0",
            "serial": 1,
            "lineage": "test-lineage",
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "migrated",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [
                        {
                            "schema_version": 0,
                            "attributes": {
                                "id": 1,
                                "name": "migrated-tenant"
                            }
                        }
                    ]
                }
            ]
        }
        
        state_file = tmp_path / "terraform.tfstate"
        with open(state_file, 'w') as f:
            json.dump(state, f)
        
        return str(state_file)
    
    def test_merge_non_vast_resources(self, temp_state_file):
        """Test merging non-VAST resources into migrated state"""
        non_vast_resources = [
            {
                "mode": "managed",
                "type": "aws_s3_bucket",
                "name": "data",
                "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
                "instances": [
                    {
                        "schema_version": 0,
                        "attributes": {
                            "id": "my-bucket",
                            "bucket": "my-bucket"
                        }
                    }
                ]
            }
        ]
        
        original_state = {"version": 4, "serial": 1}
        
        merge_non_vast_resources_to_state(temp_state_file, non_vast_resources, original_state)
        
        # Read merged state
        with open(temp_state_file, 'r') as f:
            merged_state = json.load(f)
        
        assert len(merged_state['resources']) == 2
        assert merged_state['serial'] == 2  # Incremented
        
        # Verify both resources present
        resource_types = [r['type'] for r in merged_state['resources']]
        assert 'vastdata_tenant' in resource_types
        assert 'aws_s3_bucket' in resource_types
    
    def test_merge_multiple_non_vast_resources(self, temp_state_file):
        """Test merging multiple non-VAST resources"""
        non_vast_resources = [
            {
                "mode": "managed",
                "type": "aws_s3_bucket",
                "name": "data",
                "instances": [{"attributes": {"id": "bucket1"}}]
            },
            {
                "mode": "managed",
                "type": "google_storage_bucket",
                "name": "backup",
                "instances": [{"attributes": {"id": "bucket2"}}]
            },
            {
                "mode": "data",
                "type": "aws_ami",
                "name": "ubuntu",
                "instances": [{"attributes": {"id": "ami-123"}}]
            }
        ]
        
        original_state = {"version": 4, "serial": 1}
        
        merge_non_vast_resources_to_state(temp_state_file, non_vast_resources, original_state)
        
        with open(temp_state_file, 'r') as f:
            merged_state = json.load(f)
        
        # 1 original VAST + 3 non-VAST = 4 total
        assert len(merged_state['resources']) == 4
    
    def test_merge_empty_list(self, temp_state_file):
        """Test merging empty non-VAST resource list"""
        with open(temp_state_file, 'r') as f:
            original_state_content = json.load(f)
        
        original_count = len(original_state_content['resources'])
        
        merge_non_vast_resources_to_state(temp_state_file, [], {})
        
        with open(temp_state_file, 'r') as f:
            merged_state = json.load(f)
        
        # Should remain unchanged
        assert len(merged_state['resources']) == original_count


class TestImportSingleResource:
    """Test single resource import with mocking"""
    
    @patch('subprocess.run')
    def test_successful_import(self, mock_run):
        """Test successful resource import"""
        mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")
        
        resource = {
            'type': 'vastdata_tenant',
            'address': 'vastdata_tenant.test',
            'attributes': {'id': 1}
        }
        
        import threading
        lock = threading.Lock()
        
        success, address, error = import_single_resource(resource, '/tmp/terraform', lock)
        
        assert success is True
        assert address == 'vastdata_tenant.test'
        assert error == ""
        mock_run.assert_called_once()
    
    @patch('subprocess.run')
    def test_failed_import(self, mock_run):
        """Test failed resource import"""
        mock_run.return_value = MagicMock(
            returncode=1,
            stdout="",
            stderr="Error: resource not found"
        )
        
        resource = {
            'type': 'vastdata_view',
            'address': 'vastdata_view.missing',
            'attributes': {'id': 999}
        }
        
        import threading
        lock = threading.Lock()
        
        success, address, error = import_single_resource(resource, '/tmp/terraform', lock)
        
        assert success is False
        assert address == 'vastdata_view.missing'
        assert "Error:" in error or "not found" in error.lower()
    
    @patch('subprocess.run')
    def test_import_with_resource_name_translation(self, mock_run):
        """Test import with v1 to v3 resource name translation"""
        mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")
        
        resource = {
            'type': 'vastdata_administators_managers',  # v1 name
            'address': 'vastdata_administators_managers.admin1',
            'attributes': {'id': 5}
        }
        
        import threading
        lock = threading.Lock()
        
        success, address, error = import_single_resource(resource, '/tmp/terraform', lock)
        
        assert success is True
        # Should translate to v3 name
        assert address == 'vastdata_administrator_manager.admin1'
        
        # Verify terraform import was called with translated name
        call_args = mock_run.call_args[0][0]
        assert 'vastdata_administrator_manager.admin1' in call_args


class TestRunParallelImports:
    """Test parallel import execution"""
    
    @patch('state_migration.import_single_resource')
    def test_parallel_import_all_success(self, mock_import):
        """Test parallel import with all resources succeeding"""
        # Mock successful imports
        mock_import.return_value = (True, "vastdata_tenant.test", "")
        
        resources = [
            {'type': 'vastdata_tenant', 'address': f'vastdata_tenant.test{i}', 'attributes': {'id': i}}
            for i in range(10)
        ]
        
        success_count, failed_count, failed_resources = run_parallel_imports(
            resources, '/tmp/terraform', max_workers=3
        )
        
        assert success_count == 10
        assert failed_count == 0
        assert len(failed_resources) == 0
        assert mock_import.call_count == 10
    
    @patch('state_migration.import_single_resource')
    def test_parallel_import_some_failures(self, mock_import):
        """Test parallel import with some failures"""
        # Mock: first 7 succeed, last 3 fail
        def side_effect(resource, terraform_dir, lock):
            resource_id = int(resource['address'].split('test')[1])
            if resource_id < 7:
                return (True, resource['address'], "")
            else:
                return (False, resource['address'], "Import failed")
        
        mock_import.side_effect = side_effect
        
        resources = [
            {'type': 'vastdata_view', 'address': f'vastdata_view.test{i}', 'attributes': {'id': i}}
            for i in range(10)
        ]
        
        success_count, failed_count, failed_resources = run_parallel_imports(
            resources, '/tmp/terraform', max_workers=5
        )
        
        assert success_count == 7
        assert failed_count == 3
        assert len(failed_resources) == 3
    
    @patch('state_migration.import_single_resource')
    def test_parallel_import_respects_max_workers(self, mock_import):
        """Test that parallel import respects max_workers setting"""
        import threading
        import time
        
        active_workers = []
        max_concurrent = 0
        lock = threading.Lock()
        
        def slow_import(resource, terraform_dir, state_lock):
            nonlocal max_concurrent
            with lock:
                active_workers.append(1)
                current = len(active_workers)
                if current > max_concurrent:
                    max_concurrent = current
            
            time.sleep(0.01)  # Simulate work
            
            with lock:
                active_workers.pop()
            
            return (True, resource['address'], "")
        
        mock_import.side_effect = slow_import
        
        resources = [
            {'type': 'vastdata_tenant', 'address': f'vastdata_tenant.test{i}', 'attributes': {'id': i}}
            for i in range(20)
        ]
        
        success_count, failed_count, failed_resources = run_parallel_imports(
            resources, '/tmp/terraform', max_workers=3
        )
        
        assert success_count == 20
        assert failed_count == 0
        # Max concurrent should not exceed max_workers
        assert max_concurrent <= 3


class TestEndToEndMixedProviders:
    """End-to-end tests for mixed provider migration"""
    
    @pytest.fixture
    def mixed_state_file(self, tmp_path):
        """Create a state file with mixed providers"""
        state = {
            "version": 4,
            "terraform_version": "1.5.0",
            "serial": 5,
            "lineage": "mixed-providers",
            "resources": [
                {
                    "mode": "managed",
                    "type": "vastdata_tenant",
                    "name": "vast1",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [{"schema_version": 0, "attributes": {"id": 1, "name": "vast1"}}]
                },
                {
                    "mode": "managed",
                    "type": "aws_s3_bucket",
                    "name": "data",
                    "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
                    "instances": [{"schema_version": 0, "attributes": {"id": "my-bucket"}}]
                },
                {
                    "mode": "managed",
                    "type": "vastdata_view",
                    "name": "view1",
                    "provider": "provider[\"registry.terraform.io/vast-data/vastdata\"]",
                    "instances": [{"schema_version": 0, "attributes": {"id": 10, "path": "/view1"}}]
                },
                {
                    "mode": "managed",
                    "type": "google_storage_bucket",
                    "name": "backup",
                    "provider": "provider[\"registry.terraform.io/hashicorp/google\"]",
                    "instances": [{"schema_version": 0, "attributes": {"id": "backup-bucket"}}]
                }
            ]
        }
        
        state_file = tmp_path / "terraform.tfstate"
        with open(state_file, 'w') as f:
            json.dump(state, f)
        
        return state_file
    
    def test_extract_and_separate_mixed_state(self, mixed_state_file):
        """Test extraction and separation of mixed provider state"""
        from state_migration import parse_tfstate
        
        state = parse_tfstate(str(mixed_state_file))
        resources = extract_resources(state)
        
        assert len(resources) == 4
        
        vast, non_vast = separate_vast_resources(resources)
        
        assert len(vast) == 2  # vastdata_tenant and vastdata_view
        assert len(non_vast) == 2  # aws_s3_bucket and google_storage_bucket
        
        vast_types = {r['type'] for r in vast}
        assert vast_types == {'vastdata_tenant', 'vastdata_view'}
        
        non_vast_types = {r['type'] for r in non_vast}
        assert non_vast_types == {'aws_s3_bucket', 'google_storage_bucket'}


if __name__ == '__main__':
    pytest.main([__file__, '-v'])

#!/usr/bin/env python3
"""
Test for_each and count resource handling in state migration
Tests that resources created with for_each/count are properly handled
"""

import sys
import os
import json
import tempfile
import shutil

# Add parent directory to path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from state_migration import extract_resources, separate_vast_resources, RESOURCE_NAME_TRANSLATION


def test_foreach_resource_extraction():
    """Test that for_each resources are correctly extracted from state"""
    state = {
        "version": 4,
        "resources": [
            {
                "mode": "managed",
                "type": "vastdata_view_policy",
                "name": "nfs_bu",
                "instances": [
                    {
                        "index_key": "ahl",
                        "attributes": {
                            "id": "100",
                            "name": "nfs_bu_ahl"
                        }
                    },
                    {
                        "index_key": "frm",
                        "attributes": {
                            "id": "101",
                            "name": "nfs_bu_frm"
                        }
                    }
                ]
            }
        ]
    }
    
    resources = extract_resources(state)
    
    assert len(resources) == 2, f"Expected 2 resources, got {len(resources)}"
    
    # Check addresses include the for_each key
    addresses = [r['address'] for r in resources]
    assert 'vastdata_view_policy.nfs_bu["ahl"]' in addresses, f"Missing ahl instance in {addresses}"
    assert 'vastdata_view_policy.nfs_bu["frm"]' in addresses, f"Missing frm instance in {addresses}"
    
    # Check resource names (without index)
    assert all(r['name'] == 'nfs_bu' for r in resources), "Resource name should be 'nfs_bu'"
    
    print("✓ for_each resource extraction test passed")


def test_count_resource_extraction():
    """Test that count resources are correctly extracted from state"""
    state = {
        "version": 4,
        "resources": [
            {
                "mode": "managed",
                "type": "vastdata_view",
                "name": "test_views",
                "instances": [
                    {
                        "index_key": 0,
                        "attributes": {
                            "id": "1",
                            "name": "view-0"
                        }
                    },
                    {
                        "index_key": 1,
                        "attributes": {
                            "id": "2",
                            "name": "view-1"
                        }
                    }
                ]
            }
        ]
    }
    
    resources = extract_resources(state)
    
    assert len(resources) == 2, f"Expected 2 resources, got {len(resources)}"
    
    # Check addresses include the count index
    addresses = [r['address'] for r in resources]
    assert 'vastdata_view.test_views[0]' in addresses, f"Missing index 0 in {addresses}"
    assert 'vastdata_view.test_views[1]' in addresses, f"Missing index 1 in {addresses}"
    
    print("✓ count resource extraction test passed")


def test_stub_resource_generation():
    """Test that stub resource blocks are correctly generated for for_each resources"""
    resources = [
        {
            'address': 'vastdata_view_policy.nfs_bu["ahl"]',
            'type': 'vastdata_view_policy',
            'name': 'nfs_bu',
            'module': '',
            'attributes': {'id': '100'}
        },
        {
            'address': 'vastdata_view_policy.nfs_bu["frm"]',
            'type': 'vastdata_view_policy',
            'name': 'nfs_bu',
            'module': '',
            'attributes': {'id': '101'}
        },
        {
            'address': 'vastdata_view.test_views[0]',
            'type': 'vastdata_view',
            'name': 'test_views',
            'module': '',
            'attributes': {'id': '1'}
        },
        {
            'address': 'vastdata_view.test_views[1]',
            'type': 'vastdata_view',
            'name': 'test_views',
            'module': '',
            'attributes': {'id': '2'}
        }
    ]
    
    # Simulate stub resource generation logic (new simplified approach)
    seen_resources = set()
    generated_blocks = []
    
    for resource in resources:
        resource_type_from_state = resource['type']
        resource_type = RESOURCE_NAME_TRANSLATION.get(resource_type_from_state, resource_type_from_state)
        
        # Use the base name directly from state structure (no parsing!)
        # State already has the base name without for_each/count indices
        base_name = resource['name']
        
        # Only write each unique resource block once
        resource_key = f"{resource_type}.{base_name}"
        if resource_key not in seen_resources:
            seen_resources.add(resource_key)
            generated_blocks.append(f'resource "{resource_type}" "{base_name}"')
    
    # Should generate exactly 2 unique resource blocks
    assert len(generated_blocks) == 2, f"Expected 2 unique blocks, got {len(generated_blocks)}"
    
    # Check the generated blocks
    assert 'resource "vastdata_view_policy" "nfs_bu"' in generated_blocks, \
        "Should generate resource block with base name 'nfs_bu'"
    assert 'resource "vastdata_view" "test_views"' in generated_blocks, \
        "Should generate resource block with base name 'test_views'"
    
    # Should NOT contain any brackets
    for block in generated_blocks:
        assert '[' not in block, f"Resource block should not contain brackets: {block}"
        assert '"[' not in block, f"Resource block should not contain brackets: {block}"
    
    print("✓ stub resource generation test passed")


def test_base_name_from_state():
    """Test that we use the base name directly from state structure (no parsing!)"""
    # The state file already provides the base name without indices
    state_resources = [
        {'type': 'vastdata_view_policy', 'name': 'nfs_bu'},      # Base name, no index
        {'type': 'vastdata_view', 'name': 'test_views'},          # Base name, no index
        {'type': 'vastdata_view', 'name': 'simple'},              # Base name, no index
        {'type': 'vastdata_view', 'name': 'with_underscore'},     # Base name, no index
        {'type': 'vastdata_view', 'name': 'numbered_key'},        # Base name, no index
    ]
    
    # We simply use resource['name'] - no parsing needed!
    for resource in state_resources:
        base_name = resource['name']
        # The name is already clean, no indices to strip
        assert '[' not in base_name, f"State name should not have indices: {base_name}"
        assert ']' not in base_name, f"State name should not have indices: {base_name}"
    
    print("✓ base name from state test passed")


def test_mixed_resources():
    """Test handling of mixed for_each and regular resources"""
    state = {
        "version": 4,
        "resources": [
            {
                "mode": "managed",
                "type": "vastdata_view_policy",
                "name": "foreach_policy",
                "instances": [
                    {
                        "index_key": "key1",
                        "attributes": {"id": "100"}
                    }
                ]
            },
            {
                "mode": "managed",
                "type": "vastdata_view_policy",
                "name": "simple_policy",
                "instances": [
                    {
                        "attributes": {"id": "200"}
                    }
                ]
            },
            {
                "mode": "managed",
                "type": "aws_s3_bucket",
                "name": "test_bucket",
                "instances": [
                    {
                        "attributes": {"id": "my-bucket"}
                    }
                ]
            }
        ]
    }
    
    resources = extract_resources(state)
    vast_resources, non_vast_resources = separate_vast_resources(resources)
    
    assert len(vast_resources) == 2, f"Expected 2 VAST resources, got {len(vast_resources)}"
    assert len(non_vast_resources) == 1, f"Expected 1 non-VAST resource, got {len(non_vast_resources)}"
    
    # Check VAST resources
    vast_addresses = [r['address'] for r in vast_resources]
    assert 'vastdata_view_policy.foreach_policy["key1"]' in vast_addresses, f"Missing foreach_policy in {vast_addresses}"
    assert 'vastdata_view_policy.simple_policy' in vast_addresses, f"Missing simple_policy in {vast_addresses}"
    
    # Check non-VAST resources
    assert non_vast_resources[0]['type'] == 'aws_s3_bucket'
    
    print("✓ mixed resources test passed")


if __name__ == '__main__':
    print("Running for_each/count resource tests...")
    print()
    
    test_foreach_resource_extraction()
    test_count_resource_extraction()
    test_stub_resource_generation()
    test_base_name_from_state()
    test_mixed_resources()
    
    print()
    print("=" * 80)
    print("All tests passed! ✓")
    print("=" * 80)

#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

"""
Test for_each and count resource handling in state migration.

The new state_migration.py works at the raw state-resource level
(not per-instance), so these tests verify that resources with
for_each/count instances are correctly separated (vast vs non-vast)
and stripped.
"""

import sys
import os
import json

# Add parent directory to path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from state_migration import separate_vast_resources, strip_vast_resources


def _make_state(resources, serial=5):
    return {
        "version": 4,
        "terraform_version": "1.5.0",
        "serial": serial,
        "lineage": "foreach-test",
        "outputs": {},
        "resources": resources,
    }


def test_foreach_resource_separation():
    """Test that for_each resources are correctly classified (vast vs non-vast)"""
    resources = [
        {
            "mode": "managed",
            "type": "vastdata_view_policy",
            "name": "nfs_bu",
            "instances": [
                {
                    "index_key": "ahl",
                    "attributes": {"id": "100", "name": "nfs_bu_ahl"}
                },
                {
                    "index_key": "frm",
                    "attributes": {"id": "101", "name": "nfs_bu_frm"}
                }
            ]
        },
        {
            "mode": "managed",
            "type": "aws_s3_bucket",
            "name": "buckets",
            "instances": [
                {
                    "index_key": "logs",
                    "attributes": {"id": "logs-bucket"}
                },
                {
                    "index_key": "data",
                    "attributes": {"id": "data-bucket"}
                }
            ]
        }
    ]

    vast, non_vast = separate_vast_resources(resources)

    assert len(vast) == 1, f"Expected 1 vast resource, got {len(vast)}"
    assert len(non_vast) == 1, f"Expected 1 non-vast resource, got {len(non_vast)}"

    # Vast resource keeps all its for_each instances
    assert vast[0]["type"] == "vastdata_view_policy"
    assert len(vast[0]["instances"]) == 2

    # Non-vast resource keeps all its for_each instances
    assert non_vast[0]["type"] == "aws_s3_bucket"
    assert len(non_vast[0]["instances"]) == 2

    print("✓ for_each resource separation test passed")


def test_count_resource_separation():
    """Test that count resources are correctly classified"""
    resources = [
        {
            "mode": "managed",
            "type": "vastdata_view",
            "name": "test_views",
            "instances": [
                {"index_key": 0, "attributes": {"id": "1", "name": "view-0"}},
                {"index_key": 1, "attributes": {"id": "2", "name": "view-1"}},
            ]
        },
        {
            "mode": "managed",
            "type": "aws_instance",
            "name": "servers",
            "instances": [
                {"index_key": 0, "attributes": {"id": "i-001"}},
                {"index_key": 1, "attributes": {"id": "i-002"}},
                {"index_key": 2, "attributes": {"id": "i-003"}},
            ]
        }
    ]

    vast, non_vast = separate_vast_resources(resources)

    assert len(vast) == 1
    assert len(non_vast) == 1

    # Vast resource instances preserved
    assert len(vast[0]["instances"]) == 2

    # Non-vast resource instances preserved
    assert len(non_vast[0]["instances"]) == 3

    print("✓ count resource separation test passed")


def test_strip_foreach_resources():
    """Test that strip_vast_resources removes for_each vast resources while keeping non-vast."""
    state = _make_state([
        {
            "mode": "managed",
            "type": "vastdata_view_policy",
            "name": "foreach_policy",
            "instances": [
                {"index_key": "key1", "attributes": {"id": "100"}},
                {"index_key": "key2", "attributes": {"id": "101"}},
            ]
        },
        {
            "mode": "managed",
            "type": "vastdata_view_policy",
            "name": "simple_policy",
            "instances": [
                {"attributes": {"id": "200"}}
            ]
        },
        {
            "mode": "managed",
            "type": "aws_s3_bucket",
            "name": "test_bucket",
            "instances": [
                {"attributes": {"id": "my-bucket"}}
            ]
        }
    ])

    cleaned, vc, nvc = strip_vast_resources(state)

    assert vc == 2, f"Expected 2 vast resources stripped, got {vc}"
    assert nvc == 1, f"Expected 1 non-vast resource preserved, got {nvc}"
    assert len(cleaned["resources"]) == 1
    assert cleaned["resources"][0]["type"] == "aws_s3_bucket"

    print("✓ strip for_each resources test passed")


def test_mixed_foreach_count_and_simple():
    """Test state with a mix of for_each, count, and simple resources."""
    state = _make_state([
        # for_each vastdata
        {
            "mode": "managed",
            "type": "vastdata_view_policy",
            "name": "foreach_policy",
            "instances": [
                {"index_key": "key1", "attributes": {"id": "100"}},
            ]
        },
        # count vastdata
        {
            "mode": "managed",
            "type": "vastdata_view",
            "name": "counted_views",
            "instances": [
                {"index_key": 0, "attributes": {"id": "1"}},
                {"index_key": 1, "attributes": {"id": "2"}},
            ]
        },
        # simple vastdata
        {
            "mode": "managed",
            "type": "vastdata_tenant",
            "name": "main_tenant",
            "instances": [
                {"attributes": {"id": "10"}}
            ]
        },
        # for_each AWS
        {
            "mode": "managed",
            "type": "aws_s3_bucket",
            "name": "buckets",
            "instances": [
                {"index_key": "logs", "attributes": {"id": "logs-bucket"}},
                {"index_key": "data", "attributes": {"id": "data-bucket"}},
            ]
        },
        # simple AWS
        {
            "mode": "managed",
            "type": "aws_iam_role",
            "name": "role",
            "instances": [
                {"attributes": {"id": "my-role"}}
            ]
        },
    ])

    cleaned, vc, nvc = strip_vast_resources(state)

    assert vc == 3, f"Expected 3 vast resources, got {vc}"
    assert nvc == 2, f"Expected 2 non-vast resources, got {nvc}"
    assert len(cleaned["resources"]) == 2

    remaining_types = {r["type"] for r in cleaned["resources"]}
    assert remaining_types == {"aws_s3_bucket", "aws_iam_role"}

    # Verify for_each instances are preserved on AWS bucket
    bucket = [r for r in cleaned["resources"] if r["type"] == "aws_s3_bucket"][0]
    assert len(bucket["instances"]) == 2

    print("✓ mixed for_each/count/simple test passed")


def test_foreach_instances_preserved_after_strip():
    """Ensure non-vast for_each instances are completely preserved, including index_key."""
    state = _make_state([
        {
            "mode": "managed",
            "type": "vastdata_tenant",
            "name": "t",
            "instances": [{"attributes": {"id": 1}}],
        },
        {
            "mode": "managed",
            "type": "google_compute_instance",
            "name": "vms",
            "instances": [
                {"index_key": "web", "schema_version": 0, "attributes": {"id": "vm-web", "zone": "us-east1-b"}},
                {"index_key": "api", "schema_version": 0, "attributes": {"id": "vm-api", "zone": "us-east1-c"}},
            ],
        },
    ])

    cleaned, _, _ = strip_vast_resources(state)

    assert len(cleaned["resources"]) == 1
    instances = cleaned["resources"][0]["instances"]
    assert len(instances) == 2
    assert instances[0]["index_key"] == "web"
    assert instances[1]["index_key"] == "api"
    assert instances[0]["attributes"]["zone"] == "us-east1-b"

    print("✓ for_each instances preserved after strip test passed")


if __name__ == '__main__':
    print("Running for_each/count resource tests...")
    print()

    test_foreach_resource_separation()
    test_count_resource_separation()
    test_strip_foreach_resources()
    test_mixed_foreach_count_and_simple()
    test_foreach_instances_preserved_after_strip()

    print()
    print("=" * 80)
    print("All tests passed! ✓")
    print("=" * 80)

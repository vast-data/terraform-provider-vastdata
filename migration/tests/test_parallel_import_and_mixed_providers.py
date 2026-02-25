#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Tests for mixed-provider state migration scenarios.

Focuses on realistic end-to-end scenarios where the customer's .tfstate
contains resources from VastData alongside AWS, GCP, Azure, etc.  The new
state_migration.py must strip only the VastData resources and keep the rest.
"""

import json
import os
import sys
from pathlib import Path
from unittest.mock import patch, MagicMock

import pytest

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from state_migration import (
    parse_tfstate,
    separate_vast_resources,
    strip_vast_resources,
    write_state,
    main,
)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def _make_state(resources, *, version=4, serial=5, lineage="mixed-test"):
    return {
        "version": version,
        "terraform_version": "1.5.0",
        "serial": serial,
        "lineage": lineage,
        "outputs": {},
        "resources": resources,
    }


def _managed(rtype, name, attrs, *, provider=None, module=None):
    entry = {
        "mode": "managed",
        "type": rtype,
        "name": name,
        "instances": [{"schema_version": 0, "attributes": attrs}],
    }
    if provider:
        entry["provider"] = provider
    if module:
        entry["module"] = module
    return entry


def _data(rtype, name, attrs, *, provider=None):
    entry = {
        "mode": "data",
        "type": rtype,
        "name": name,
        "instances": [{"schema_version": 0, "attributes": attrs}],
    }
    if provider:
        entry["provider"] = provider
    return entry


def _write(path, obj):
    with open(path, "w") as fh:
        json.dump(obj, fh, indent=2)


# ---------------------------------------------------------------------------
# Separation tests
# ---------------------------------------------------------------------------

class TestSeparateVastResources:
    """Test the resource-splitting logic with various provider mixes."""

    def test_all_vast(self):
        resources = [
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("vastdata_view", "v1", {"id": 2}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 2
        assert len(non_vast) == 0
        assert all(r["type"].startswith("vastdata_") for r in vast)

    def test_all_non_vast(self):
        resources = [
            _managed("aws_s3_bucket", "b1", {"id": "bucket-1"}),
            _managed("google_storage_bucket", "gs1", {"id": "gs-1"}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 0
        assert len(non_vast) == 2

    def test_mixed_aws_vastdata(self):
        resources = [
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
            _managed("vastdata_view", "v1", {"id": 2}),
            _managed("aws_iam_role", "role1", {"id": "role-1"}),
            _managed("vastdata_vip_pool", "pool1", {"id": 3}),
        ]
        vast, non_vast = separate_vast_resources(resources)

        assert len(vast) == 3
        assert len(non_vast) == 2

        assert {r["type"] for r in vast} == {
            "vastdata_tenant", "vastdata_view", "vastdata_vip_pool",
        }
        assert {r["type"] for r in non_vast} == {"aws_s3_bucket", "aws_iam_role"}

    def test_mixed_gcp_vastdata(self):
        resources = [
            _managed("vastdata_quota", "q1", {"id": 1}),
            _managed("google_compute_instance", "vm1", {"id": "vm-1"}),
            _managed("google_storage_bucket", "gs1", {"id": "gs-1"}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 1
        assert len(non_vast) == 2

    def test_mixed_azure_vastdata(self):
        resources = [
            _managed("vastdata_snapshot", "snap1", {"id": 1}),
            _managed("azurerm_resource_group", "rg1", {"id": "rg-1"}),
            _managed("azurerm_virtual_network", "vnet1", {"id": "vnet-1"}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 1
        assert len(non_vast) == 2

    def test_mixed_multiple_clouds(self):
        """Realistic: AWS + GCP + Azure + VastData."""
        resources = [
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("vastdata_view", "v1", {"id": 2}),
            _managed("vastdata_view_policy", "vp1", {"id": 3}),
            _managed("aws_s3_bucket", "b1", {"id": "b-1"}),
            _managed("aws_iam_role", "role1", {"id": "r-1"}),
            _managed("google_compute_instance", "vm1", {"id": "vm-1"}),
            _managed("azurerm_resource_group", "rg1", {"id": "rg-1"}),
            _managed("kubernetes_deployment", "k1", {"id": "deploy-1"}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 3
        assert len(non_vast) == 5

    def test_empty_list(self):
        vast, non_vast = separate_vast_resources([])
        assert vast == []
        assert non_vast == []

    def test_data_sources_mixed(self):
        resources = [
            _data("vastdata_vip_pool", "pool", {"id": 1}),
            _data("aws_ami", "ubuntu", {"id": "ami-123"}),
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 2   # data + managed
        assert len(non_vast) == 2

    def test_v1_legacy_names_still_classified_as_vast(self):
        """Old v1/v2 resource names like vastdata_administators_managers."""
        resources = [
            _managed("vastdata_administators_managers", "admin", {"id": 1}),
            _managed("vastdata_kafka_brokers", "broker", {"id": 2}),
            _managed("vastdata_non_local_user", "user", {"username": "x", "context": "ldap", "tenant_id": 1}),
            _managed("vastdata_blockhost", "host", {"id": 3}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 4
        assert len(non_vast) == 1


# ---------------------------------------------------------------------------
# Strip tests with mixed providers
# ---------------------------------------------------------------------------

class TestStripMixedProviders:
    """Test strip_vast_resources with multi-provider state."""

    def test_strips_vast_preserves_aws(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("vastdata_view", "v1", {"id": 2}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket-1"}),
            _managed("aws_iam_role", "role1", {"id": "role-1"}),
        ], serial=10)

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 2
        assert nvc == 2
        assert len(cleaned["resources"]) == 2
        assert all(r["type"].startswith("aws_") for r in cleaned["resources"])
        assert cleaned["serial"] == 11

    def test_strips_vast_preserves_gcp_and_azure(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("google_compute_instance", "vm", {"id": "vm-1"}),
            _managed("azurerm_resource_group", "rg", {"id": "rg-1"}),
        ])

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 1
        assert nvc == 2
        types = {r["type"] for r in cleaned["resources"]}
        assert types == {"google_compute_instance", "azurerm_resource_group"}

    def test_preserves_data_sources_of_other_providers(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _data("aws_ami", "ubuntu", {"id": "ami-123"}),
            _data("vastdata_vip_pool", "pool", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ])

        cleaned, vc, nvc = strip_vast_resources(state)

        # vast: tenant (managed) + vip_pool (data) = 2
        assert vc == 2
        # non-vast: ami (data) + s3 (managed) = 2
        assert nvc == 2
        types = {r["type"] for r in cleaned["resources"]}
        assert types == {"aws_ami", "aws_s3_bucket"}

    def test_preserves_resource_attributes_intact(self):
        """Non-vast resource attributes must survive stripping untouched."""
        aws_attrs = {
            "id": "my-bucket",
            "bucket": "my-bucket",
            "acl": "private",
            "tags": {"env": "prod", "team": "infra"},
            "versioning": [{"enabled": True}],
        }
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "data_bucket", aws_attrs,
                     provider='provider["registry.terraform.io/hashicorp/aws"]'),
        ])

        cleaned, _, _ = strip_vast_resources(state)

        preserved = cleaned["resources"][0]
        assert preserved["type"] == "aws_s3_bucket"
        assert preserved["name"] == "data_bucket"
        assert preserved["instances"][0]["attributes"] == aws_attrs
        assert preserved["provider"] == 'provider["registry.terraform.io/hashicorp/aws"]'

    def test_preserves_module_non_vast_resources(self):
        state = _make_state([
            _managed("vastdata_view", "v1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "b"}, module="module.storage"),
            _managed("vastdata_tenant", "t1", {"id": 2}, module="module.vast"),
        ])

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 2
        assert nvc == 1
        assert cleaned["resources"][0]["module"] == "module.storage"

    def test_preserves_outputs(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
        ])
        state["outputs"] = {
            "bucket_arn": {"value": "arn:aws:s3:::my-bucket", "type": "string"},
        }

        cleaned, _, _ = strip_vast_resources(state)

        assert cleaned["outputs"]["bucket_arn"]["value"] == "arn:aws:s3:::my-bucket"


# ---------------------------------------------------------------------------
# End-to-end mixed-provider CLI tests
# ---------------------------------------------------------------------------

class TestEndToEndMixedProviders:
    """Integration tests using main() with multi-provider states."""

    def _create_workspace(self, tmp_path, state_resources, *, serial=5):
        workdir = tmp_path / "workdir"
        workdir.mkdir()
        (workdir / "main.tf").write_text(
            'resource "vastdata_tenant" "t" { name = "x" }\n'
            'resource "aws_s3_bucket" "b" { bucket = "b" }\n'
        )
        state = _make_state(state_resources, serial=serial)
        state_file = tmp_path / "original.tfstate"
        _write(state_file, state)
        return workdir, state_file

    @patch("state_migration.run_terraform_apply_migrate")
    @patch("state_migration.verify_migrate_mode_support")
    @patch("state_migration.run_terraform_init")
    def test_mixed_state_strips_vast_keeps_aws(self, mock_init, mock_verify, mock_apply, tmp_path, monkeypatch):
        workdir, state_file = self._create_workspace(tmp_path, [
            _managed("vastdata_tenant", "t1", {"id": 1, "name": "prod"}),
            _managed("vastdata_view", "v1", {"id": 2, "path": "/data"}),
            _managed("aws_s3_bucket", "data", {"id": "data-bucket", "bucket": "data-bucket"}),
            _managed("aws_iam_role", "role", {"id": "role-1", "name": "terraform-role"}),
        ], serial=15)

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file)])

        main()

        # Verify cleaned state
        dest = workdir / "terraform.tfstate"
        with open(dest) as fh:
            cleaned = json.load(fh)

        assert len(cleaned["resources"]) == 2
        types = {r["type"] for r in cleaned["resources"]}
        assert types == {"aws_s3_bucket", "aws_iam_role"}
        assert cleaned["serial"] == 16

        mock_init.assert_called_once()
        mock_verify.assert_called_once()
        mock_apply.assert_called_once()

    @patch("state_migration.run_terraform_apply_migrate")
    @patch("state_migration.verify_migrate_mode_support")
    @patch("state_migration.run_terraform_init")
    def test_all_vast_state_produces_empty_resources(self, mock_init, mock_verify, mock_apply, tmp_path, monkeypatch):
        workdir, state_file = self._create_workspace(tmp_path, [
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("vastdata_view", "v1", {"id": 2}),
            _managed("vastdata_quota", "q1", {"id": 3}),
        ])

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file)])

        main()

        dest = workdir / "terraform.tfstate"
        with open(dest) as fh:
            cleaned = json.load(fh)

        assert cleaned["resources"] == []

    @patch("state_migration.run_terraform_apply_migrate")
    @patch("state_migration.run_terraform_init")
    def test_dry_run_mixed_state(self, mock_init, mock_apply, tmp_path, monkeypatch):
        workdir, state_file = self._create_workspace(tmp_path, [
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ])

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file), "--dry-run"])

        with pytest.raises(SystemExit) as exc_info:
            main()

        assert exc_info.value.code == 0
        mock_init.assert_not_called()
        mock_apply.assert_not_called()

        # State should still be written
        dest = workdir / "terraform.tfstate"
        with open(dest) as fh:
            cleaned = json.load(fh)
        assert len(cleaned["resources"]) == 1
        assert cleaned["resources"][0]["type"] == "aws_s3_bucket"

    @patch("state_migration.run_terraform_apply_migrate")
    @patch("state_migration.run_terraform_init")
    def test_no_vast_in_mixed_state(self, mock_init, mock_apply, tmp_path, monkeypatch):
        """State has only non-vast resources → nothing to migrate."""
        workdir, state_file = self._create_workspace(tmp_path, [
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
            _managed("google_compute_instance", "vm1", {"id": "vm"}),
        ])

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file)])

        with pytest.raises(SystemExit) as exc_info:
            main()

        assert exc_info.value.code == 0
        mock_init.assert_not_called()
        mock_apply.assert_not_called()


# ---------------------------------------------------------------------------
# Realistic production-like scenarios
# ---------------------------------------------------------------------------

class TestRealisticScenarios:
    """Simulate real customer environments."""

    def test_customer_with_aws_gcp_vastdata(self):
        """
        Customer has:
          - 10 VastData resources (tenants, views, quotas, policies)
          - 15 AWS resources (S3 buckets, IAM roles, EC2 instances)
          - 5  GCP resources (storage buckets, compute instances)
        """
        resources = []

        # VastData
        for i in range(3):
            resources.append(_managed("vastdata_tenant", f"tenant_{i}", {"id": i}))
        for i in range(4):
            resources.append(_managed("vastdata_view", f"view_{i}", {"id": 10 + i}))
        for i in range(2):
            resources.append(_managed("vastdata_quota", f"quota_{i}", {"id": 20 + i}))
        resources.append(_managed("vastdata_view_policy", "policy_0", {"id": 30}))

        # AWS
        for i in range(5):
            resources.append(_managed("aws_s3_bucket", f"bucket_{i}", {"id": f"b-{i}"}))
        for i in range(5):
            resources.append(_managed("aws_iam_role", f"role_{i}", {"id": f"r-{i}"}))
        for i in range(5):
            resources.append(_managed("aws_instance", f"ec2_{i}", {"id": f"i-{i}"}))

        # GCP
        for i in range(3):
            resources.append(_managed("google_storage_bucket", f"gs_{i}", {"id": f"gs-{i}"}))
        for i in range(2):
            resources.append(_managed("google_compute_instance", f"gce_{i}", {"id": f"gce-{i}"}))

        state = _make_state(resources, serial=100)

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 10
        assert nvc == 20
        assert len(cleaned["resources"]) == 20
        assert cleaned["serial"] == 101

        # None of the remaining should be vastdata
        assert all(
            not r["type"].startswith("vastdata_")
            for r in cleaned["resources"]
        )

    def test_customer_with_for_each_vastdata_and_aws(self):
        """Resources with for_each instances."""
        vast_resource = {
            "mode": "managed",
            "type": "vastdata_view",
            "name": "views",
            "instances": [
                {"index_key": "prod", "schema_version": 0, "attributes": {"id": 1}},
                {"index_key": "staging", "schema_version": 0, "attributes": {"id": 2}},
                {"index_key": "dev", "schema_version": 0, "attributes": {"id": 3}},
            ],
        }
        aws_resource = {
            "mode": "managed",
            "type": "aws_s3_bucket",
            "name": "buckets",
            "instances": [
                {"index_key": "logs", "schema_version": 0, "attributes": {"id": "logs-bucket"}},
                {"index_key": "data", "schema_version": 0, "attributes": {"id": "data-bucket"}},
            ],
        }

        state = _make_state([vast_resource, aws_resource])

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 1   # 1 resource entry (with 3 instances)
        assert nvc == 1  # 1 resource entry (with 2 instances)
        assert len(cleaned["resources"]) == 1
        # AWS instances should be intact
        assert len(cleaned["resources"][0]["instances"]) == 2

    def test_customer_with_modules_mixed(self):
        """Resources inside modules from different providers."""
        resources = [
            _managed("vastdata_tenant", "t", {"id": 1}, module="module.vast_infra"),
            _managed("vastdata_view", "v", {"id": 2}, module="module.vast_infra"),
            _managed("aws_s3_bucket", "b", {"id": "bucket"}, module="module.aws_storage"),
            _managed("aws_iam_role", "r", {"id": "role"}, module="module.aws_iam"),
            _managed("google_compute_instance", "vm", {"id": "vm"}, module="module.gcp_compute"),
        ]

        state = _make_state(resources)

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 2
        assert nvc == 3

        # Verify modules are preserved
        modules = {r["module"] for r in cleaned["resources"]}
        assert "module.aws_storage" in modules
        assert "module.aws_iam" in modules
        assert "module.gcp_compute" in modules

    def test_customer_state_with_sensitive_attrs(self):
        """Non-vast resources may have sensitive_attributes — must be preserved."""
        aws_resource = {
            "mode": "managed",
            "type": "aws_iam_access_key",
            "name": "key",
            "instances": [
                {
                    "schema_version": 0,
                    "attributes": {
                        "id": "AKIA...",
                        "secret": "s3cr3t",
                    },
                    "sensitive_attributes": [
                        [{"type": "get_attr", "value": "secret"}],
                    ],
                    "private": "base64data==",
                }
            ],
        }

        state = _make_state([
            _managed("vastdata_tenant", "t", {"id": 1}),
            aws_resource,
        ])

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 1
        assert nvc == 1

        preserved = cleaned["resources"][0]
        assert preserved["type"] == "aws_iam_access_key"
        inst = preserved["instances"][0]
        assert inst["sensitive_attributes"] == [
            [{"type": "get_attr", "value": "secret"}],
        ]
        assert inst["private"] == "base64data=="


if __name__ == "__main__":
    pytest.main([__file__, "-v"])

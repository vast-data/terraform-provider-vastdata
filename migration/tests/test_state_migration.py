#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Tests for state_migration.py — the new VastData state migration script.

Covers:
  - Parsing .tfstate files
  - Separating vast vs non-vast resources
  - Stripping VastData resources from state
  - Writing cleaned state
  - CLI / main() behaviour (dry-run, missing files, etc.)
"""

import json
import os
import sys
import shutil
import subprocess
from pathlib import Path
from unittest.mock import patch, MagicMock

import pytest

# Add parent directory to path so we can import the migration script
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from state_migration import (
    parse_tfstate,
    separate_vast_resources,
    strip_vast_resources,
    extract_resource_ids,
    extract_v1_attributes,
    patch_imported_attributes,
    prettify_json_fields,
    write_state,
    run_terraform_init,
    run_terraform_apply_migrate,
    verify_migrate_mode_support,
    get_vastdata_provider_version,
    main,
)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def _make_state(resources, *, version=4, serial=5, lineage="test-lineage"):
    """Build a minimal TF state dict."""
    return {
        "version": version,
        "terraform_version": "1.5.0",
        "serial": serial,
        "lineage": lineage,
        "outputs": {},
        "resources": resources,
    }


def _managed(rtype, name, attrs, *, provider=None, module=None):
    """Build a managed-resource entry."""
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
    """Build a data-source entry."""
    entry = {
        "mode": "data",
        "type": rtype,
        "name": name,
        "instances": [{"schema_version": 0, "attributes": attrs}],
    }
    if provider:
        entry["provider"] = provider
    return entry


def _write_state_file(path, state_dict):
    with open(path, "w") as fh:
        json.dump(state_dict, fh, indent=2)


# ---------------------------------------------------------------------------
# TestParseTfstate
# ---------------------------------------------------------------------------

class TestParseTfstate:
    """Tests for parse_tfstate()."""

    def test_valid_state(self, tmp_path):
        state = _make_state([_managed("vastdata_tenant", "t", {"id": 1})])
        f = tmp_path / "ok.tfstate"
        _write_state_file(f, state)

        result = parse_tfstate(str(f))

        assert result["version"] == 4
        assert len(result["resources"]) == 1

    def test_invalid_json(self, tmp_path):
        f = tmp_path / "bad.tfstate"
        f.write_text("{ not valid json !!!")

        with pytest.raises(SystemExit):
            parse_tfstate(str(f))

    def test_missing_file(self):
        with pytest.raises(SystemExit):
            parse_tfstate("/does/not/exist.tfstate")

    def test_empty_resources(self, tmp_path):
        state = _make_state([])
        f = tmp_path / "empty.tfstate"
        _write_state_file(f, state)

        result = parse_tfstate(str(f))
        assert result["resources"] == []


# ---------------------------------------------------------------------------
# TestSeparateVastResources
# ---------------------------------------------------------------------------

class TestSeparateVastResources:
    """Tests for separate_vast_resources()."""

    def test_all_vast(self):
        resources = [
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("vastdata_view", "v1", {"id": 2}),
            _managed("vastdata_vip_pool", "p1", {"id": 3}),
        ]

        vast, non_vast = separate_vast_resources(resources)

        assert len(vast) == 3
        assert len(non_vast) == 0

    def test_all_non_vast(self):
        resources = [
            _managed("aws_s3_bucket", "b1", {"id": "bucket-1"}),
            _managed("google_compute_instance", "vm1", {"id": "vm-1"}),
            _managed("azurerm_resource_group", "rg1", {"id": "rg-1"}),
        ]

        vast, non_vast = separate_vast_resources(resources)

        assert len(vast) == 0
        assert len(non_vast) == 3

    def test_mixed_providers(self):
        resources = [
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
            _managed("vastdata_view", "v1", {"id": 2}),
            _managed("google_storage_bucket", "gs1", {"id": "gs-bucket"}),
            _managed("vastdata_quota", "q1", {"id": 3}),
        ]

        vast, non_vast = separate_vast_resources(resources)

        assert len(vast) == 3
        assert len(non_vast) == 2

        vast_types = {r["type"] for r in vast}
        assert vast_types == {"vastdata_tenant", "vastdata_view", "vastdata_quota"}

        non_vast_types = {r["type"] for r in non_vast}
        assert non_vast_types == {"aws_s3_bucket", "google_storage_bucket"}

    def test_empty_list(self):
        vast, non_vast = separate_vast_resources([])
        assert vast == []
        assert non_vast == []

    def test_data_sources_are_preserved_not_imported(self):
        """VastData data sources stay in state and are never terraform import-ed."""
        resources = [
            _data("vastdata_vip_pool", "pool", {"id": 1}),
            _data("aws_ami", "ubuntu", {"id": "ami-123"}),
        ]

        importable, preserved, non_vast = separate_vast_resources(resources)

        assert len(importable) == 0
        assert len(preserved) == 1
        assert preserved[0]["type"] == "vastdata_vip_pool"
        assert len(non_vast) == 1
        assert non_vast[0]["type"] == "aws_ami"

    def test_v1_legacy_resource_names(self):
        """v1/v2 legacy vastdata names still start with 'vastdata_'."""
        resources = [
            _managed("vastdata_administators_managers", "admin", {"id": 1}),
            _managed("vastdata_kafka_brokers", "broker", {"id": 2}),
            _managed("vastdata_active_directory2", "ad", {"id": 3}),
            _managed("vastdata_non_local_user", "user", {"username": "x", "context": "ldap", "tenant_id": 1}),
            _managed("aws_iam_role", "role", {"id": "role-1"}),
        ]

        vast, non_vast = separate_vast_resources(resources)

        assert len(vast) == 4
        assert len(non_vast) == 1

    def test_resource_with_missing_type(self):
        """A resource dict without 'type' key should land in non-vast."""
        resources = [{"name": "orphan", "instances": []}]

        vast, non_vast = separate_vast_resources(resources)

        assert len(vast) == 0
        assert len(non_vast) == 1


# ---------------------------------------------------------------------------
# TestExtractResourceIds
# ---------------------------------------------------------------------------

class TestExtractResourceIds:
    """Import addresses must use v3 resource types from migration_script."""

    def test_renamed_resource_type_uses_v3_address(self):
        state = _make_state([
            _managed("vastdata_administators_roles", "read_only", {"id": 18}),
        ])

        result = extract_resource_ids(state)

        assert result == {"vastdata_administrator_role.read_only": "18"}

    def test_unchanged_resource_type_keeps_address(self):
        state = _make_state([
            _managed("vastdata_tenant", "tenant1", {"id": 3}),
        ])

        result = extract_resource_ids(state)

        assert result == {"vastdata_tenant.tenant1": "3"}

    def test_renamed_resource_with_for_each_index(self):
        state = _make_state([{
            "mode": "managed",
            "type": "vastdata_replication_peers",
            "name": "peer",
            "instances": [
                {
                    "index_key": "a",
                    "schema_version": 0,
                    "attributes": {"id": 7},
                },
            ],
        }])

        result = extract_resource_ids(state)

        assert result == {'vastdata_replication_peer.peer["a"]': "7"}

    def test_module_resource_includes_module_path(self):
        state = _make_state([
            _managed(
                "vastdata_view",
                "this",
                {"id": 20, "path": "/tfmod-env1", "tenant_name": "default"},
                module='module.views["env1"]',
            ),
        ])

        result = extract_resource_ids(state)

        assert result == {
            'module.views["env1"].vastdata_view.this': "/tfmod-env1|default",
        }

    def test_module_for_each_instances_get_unique_addresses(self):
        state = _make_state([
            _managed(
                "vastdata_view_policy",
                "nfs",
                {"id": 5, "name": "tfmod-policy", "tenant_name": "default"},
            ),
            _managed(
                "vastdata_view",
                "this",
                {"id": 20, "path": "/tfmod-env2", "tenant_name": "default"},
                module='module.views["env2"]',
            ),
            _managed(
                "vastdata_view",
                "this",
                {"id": 21, "path": "/tfmod-env1", "tenant_name": "default"},
                module='module.views["env1"]',
            ),
        ])

        result = extract_resource_ids(state)

        assert len(result) == 3
        assert result['vastdata_view_policy.nfs'] == "tfmod-policy|default"
        assert result['module.views["env1"].vastdata_view.this'] == "/tfmod-env1|default"
        assert result['module.views["env2"].vastdata_view.this'] == "/tfmod-env2|default"

    def test_data_sources_are_excluded_from_import(self):
        state = _make_state([
            _data("vastdata_tenant", "default", {"id": 1, "name": "default"}),
            _managed(
                "vastdata_view_policy",
                "nfs",
                {"id": 5, "name": "tfds-policy", "tenant_name": "default"},
            ),
        ])

        result = extract_resource_ids(state)

        assert result == {"vastdata_view_policy.nfs": "5"}


class TestPreserveCreateOnlyAttributes:
    """create-only attributes must survive import via v1 state patching."""

    def test_extract_v1_attributes_maps_view_create_dir(self):
        state = _make_state([
            _managed("vastdata_view", "global_mgmt", {
                "id": 42,
                "path": "/global_mgmt",
                "create_dir": True,
            }),
        ])

        result = extract_v1_attributes(state)

        assert result == {
            "vastdata_view.global_mgmt": {
                "id": 42,
                "path": "/global_mgmt",
                "create_dir": True,
            },
        }

    def test_patch_imported_attributes_restores_null_create_dir(self, tmp_path):
        v1_state = _make_state([
            _managed("vastdata_view", "global_mgmt", {
                "id": 42,
                "path": "/global_mgmt",
                "create_dir": True,
            }),
        ])
        v1_attrs = extract_v1_attributes(v1_state)
        imported = {"vastdata_view.global_mgmt": "42"}

        post_import = _make_state([
            _managed("vastdata_view", "global_mgmt", {
                "id": 42,
                "path": "/global_mgmt",
                "create_dir": None,
            }),
        ])
        state_path = tmp_path / "terraform.tfstate"
        write_state(post_import, str(state_path))

        patched = patch_imported_attributes(str(tmp_path), v1_attrs, imported)

        assert patched == 1
        with open(state_path) as fh:
            loaded = json.load(fh)
        attrs = loaded["resources"][0]["instances"][0]["attributes"]
        assert attrs["create_dir"] is True

    def test_prettify_json_fields_restores_v1_policy_format(self, tmp_path):
        pretty_policy = (
            '{\n  "Version": "2012-10-17",\n  "Statement": []\n}'
        )
        compact_policy = '{"Statement":[],"Version":"2012-10-17"}'

        v1_state = _make_state([
            _managed("vastdata_s3_policy", "s3policy_ro", {
                "id": 3,
                "name": "s3_codex_ro",
                "policy": pretty_policy,
            }),
        ])
        v1_attrs = extract_v1_attributes(v1_state)
        imported = {"vastdata_s3_policy.s3policy_ro": "3"}

        post_import = _make_state([
            _managed("vastdata_s3_policy", "s3policy_ro", {
                "id": 3,
                "name": "s3_codex_ro",
                "policy": pretty_policy,
            }),
        ])
        state_path = tmp_path / "terraform.tfstate"
        write_state(post_import, str(state_path))

        count = prettify_json_fields(str(tmp_path), v1_attrs, imported)

        assert count == 1
        with open(state_path) as fh:
            loaded = json.load(fh)
        attrs = loaded["resources"][0]["instances"][0]["attributes"]
        assert attrs["policy"] == compact_policy

    def test_prettify_json_fields_prettifies_without_v1(self, tmp_path):
        pretty_policy = (
            '{\n  "Version": "2012-10-17",\n  "Statement": [{"Effect": "Allow"}]\n}'
        )
        compact_policy = '{"Statement":[{"Effect":"Allow"}],"Version":"2012-10-17"}'

        post_import = _make_state([
            _managed("vastdata_s3_policy", "s3policy_ro", {
                "id": 3,
                "policy": pretty_policy,
            }),
        ])
        state_path = tmp_path / "terraform.tfstate"
        write_state(post_import, str(state_path))

        count = prettify_json_fields(str(tmp_path), {}, {"vastdata_s3_policy.s3policy_ro": "3"})

        assert count == 1
        with open(state_path) as fh:
            loaded = json.load(fh)
        attrs = loaded["resources"][0]["instances"][0]["attributes"]
        assert attrs["policy"] == compact_policy


# ---------------------------------------------------------------------------
# TestStripVastResources
# ---------------------------------------------------------------------------

class TestStripVastResources:
    """Tests for strip_vast_resources()."""

    def test_strips_only_vast(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
            _managed("vastdata_view", "v1", {"id": 2}),
        ], serial=10)

        cleaned, vast_count, non_vast_count = strip_vast_resources(state)

        assert vast_count == 2
        assert non_vast_count == 1
        assert len(cleaned["resources"]) == 1
        assert cleaned["resources"][0]["type"] == "aws_s3_bucket"

    def test_serial_is_bumped(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
        ], serial=42)

        cleaned, _, _ = strip_vast_resources(state)

        assert cleaned["serial"] == 43

    def test_no_vast_resources(self):
        state = _make_state([
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ], serial=7)

        cleaned, vast_count, non_vast_count = strip_vast_resources(state)

        assert vast_count == 0
        assert non_vast_count == 1
        assert len(cleaned["resources"]) == 1
        assert cleaned["serial"] == 8  # still bumped

    def test_all_vast_resources(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("vastdata_view", "v1", {"id": 2}),
        ])

        cleaned, vast_count, non_vast_count = strip_vast_resources(state)

        assert vast_count == 2
        assert non_vast_count == 0
        assert cleaned["resources"] == []

    def test_empty_resources(self):
        state = _make_state([])

        cleaned, vast_count, non_vast_count = strip_vast_resources(state)

        assert vast_count == 0
        assert non_vast_count == 0

    def test_no_resources_key(self):
        state = {"version": 4, "serial": 1}

        cleaned, vast_count, non_vast_count = strip_vast_resources(state)

        assert vast_count == 0
        assert non_vast_count == 0
        # Original state returned as-is
        assert cleaned is state

    def test_preserves_non_resource_fields(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ], serial=5, lineage="my-lineage")
        state["outputs"] = {"url": {"value": "https://example.com"}}

        cleaned, _, _ = strip_vast_resources(state)

        assert cleaned["lineage"] == "my-lineage"
        assert cleaned["outputs"] == {"url": {"value": "https://example.com"}}
        assert cleaned["version"] == 4

    def test_preserves_data_sources_of_non_vast(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _data("aws_ami", "ubuntu", {"id": "ami-123"}),
            _data("vastdata_vip_pool", "pool", {"id": 1}),
        ])

        cleaned, importable_count, preserved_count, non_vast_count = strip_vast_resources(state)

        assert importable_count == 1
        assert preserved_count == 1
        assert non_vast_count == 1
        assert len(cleaned["resources"]) == 2
        types = {r["type"] for r in cleaned["resources"]}
        assert types == {"aws_ami", "vastdata_vip_pool"}

    def test_large_mixed_state(self):
        """State with many resources from multiple providers."""
        resources = []
        for i in range(50):
            resources.append(_managed("vastdata_view", f"view_{i}", {"id": i}))
        for i in range(30):
            resources.append(_managed("aws_s3_bucket", f"bucket_{i}", {"id": f"b-{i}"}))
        for i in range(20):
            resources.append(_managed("google_compute_instance", f"vm_{i}", {"id": f"vm-{i}"}))

        state = _make_state(resources)

        cleaned, vast_count, non_vast_count = strip_vast_resources(state)

        assert vast_count == 50
        assert non_vast_count == 50
        assert len(cleaned["resources"]) == 50
        assert all(not r["type"].startswith("vastdata_") for r in cleaned["resources"])


# ---------------------------------------------------------------------------
# TestWriteState
# ---------------------------------------------------------------------------

class TestWriteState:
    """Tests for write_state()."""

    def test_writes_valid_json(self, tmp_path):
        state = _make_state([_managed("aws_s3_bucket", "b1", {"id": "x"})])
        out = str(tmp_path / "out.tfstate")

        write_state(state, out)

        with open(out) as fh:
            loaded = json.load(fh)

        assert loaded == state

    def test_overwrites_existing_file(self, tmp_path):
        out = tmp_path / "out.tfstate"
        out.write_text('{"old": true}')

        new_state = _make_state([])
        write_state(new_state, str(out))

        with open(out) as fh:
            loaded = json.load(fh)

        assert loaded == new_state
        assert "old" not in loaded


# ---------------------------------------------------------------------------
# TestRunTerraformInit
# ---------------------------------------------------------------------------

class TestRunTerraformInit:
    """Tests for run_terraform_init()."""

    @patch("subprocess.run")
    def test_success(self, mock_run):
        mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")

        run_terraform_init("/some/dir")

        mock_run.assert_called_once_with(
            ["terraform", "init"],
            cwd="/some/dir",
            capture_output=True,
            text=True,
        )

    @patch("subprocess.run")
    def test_failure_exits(self, mock_run):
        mock_run.return_value = MagicMock(returncode=1, stdout="", stderr="some error")

        with pytest.raises(SystemExit):
            run_terraform_init("/some/dir")

    @patch("subprocess.run")
    def test_dev_overrides_failure_continues(self, mock_run):
        """When dev_overrides are active, init may fail — we should continue."""
        mock_run.return_value = MagicMock(
            returncode=1,
            stdout="",
            stderr=(
                "Warning: Provider development overrides are in effect\n"
                "Error: Failed to query available provider packages\n"
            ),
        )

        # Should NOT raise SystemExit
        run_terraform_init("/some/dir")

    @patch("subprocess.run")
    def test_genuine_failure_still_exits(self, mock_run):
        """A real init failure (not dev overrides) should still exit."""
        mock_run.return_value = MagicMock(
            returncode=1,
            stdout="",
            stderr="Error: Could not load plugin\n",
        )

        with pytest.raises(SystemExit):
            run_terraform_init("/some/dir")


# ---------------------------------------------------------------------------
# TestRunTerraformApplyMigrate
# ---------------------------------------------------------------------------

class TestRunTerraformApplyMigrate:
    """Tests for run_terraform_apply_migrate()."""

    @patch("subprocess.run")
    def test_success(self, mock_run):
        mock_run.return_value = MagicMock(returncode=0)

        run_terraform_apply_migrate("/some/dir")

        args, kwargs = mock_run.call_args
        assert args[0] == ["terraform", "apply", "-auto-approve"]
        assert kwargs["cwd"] == "/some/dir"
        assert kwargs["env"]["VASTDATA_MIGRATE_MODE"] == "1"

    @patch("subprocess.run")
    def test_failure_exits(self, mock_run):
        mock_run.return_value = MagicMock(returncode=1)

        with pytest.raises(SystemExit):
            run_terraform_apply_migrate("/some/dir")

    @patch("subprocess.run")
    def test_inherits_environment(self, mock_run):
        """VASTDATA_MIGRATE_MODE should be added on top of the current env."""
        mock_run.return_value = MagicMock(returncode=0)

        with patch.dict(os.environ, {"MY_CUSTOM_VAR": "hello"}):
            run_terraform_apply_migrate("/some/dir")

        env_passed = mock_run.call_args[1]["env"]
        assert env_passed["MY_CUSTOM_VAR"] == "hello"
        assert env_passed["VASTDATA_MIGRATE_MODE"] == "1"


# ---------------------------------------------------------------------------
# TestGetVastdataProviderVersion
# ---------------------------------------------------------------------------

class TestGetVastdataProviderVersion:
    """Tests for get_vastdata_provider_version()."""

    @patch("subprocess.run")
    def test_returns_version_from_provider_selections(self, mock_run):
        mock_run.return_value = MagicMock(
            returncode=0,
            stdout=json.dumps({
                "terraform_version": "1.14.7",
                "provider_selections": {
                    "registry.terraform.io/vast-data/vastdata": "3.1.1",
                    "registry.terraform.io/hashicorp/null": "3.2.4",
                },
            }),
        )
        assert get_vastdata_provider_version("/some/dir") == "3.1.1"

    @patch("subprocess.run")
    def test_returns_empty_string_when_no_vastdata_provider(self, mock_run):
        mock_run.return_value = MagicMock(
            returncode=0,
            stdout=json.dumps({
                "terraform_version": "1.14.7",
                "provider_selections": {
                    "registry.terraform.io/hashicorp/null": "3.2.4",
                },
            }),
        )
        assert get_vastdata_provider_version("/some/dir") == ""

    @patch("subprocess.run")
    def test_returns_empty_string_on_command_failure(self, mock_run):
        mock_run.return_value = MagicMock(returncode=1, stdout="", stderr="error")
        assert get_vastdata_provider_version("/some/dir") == ""

    @patch("subprocess.run")
    def test_returns_empty_string_on_invalid_json(self, mock_run):
        mock_run.return_value = MagicMock(returncode=0, stdout="not-json")
        assert get_vastdata_provider_version("/some/dir") == ""

    @patch("subprocess.run")
    def test_matches_vastdata_case_insensitively(self, mock_run):
        mock_run.return_value = MagicMock(
            returncode=0,
            stdout=json.dumps({
                "provider_selections": {
                    "registry.terraform.io/VAST-DATA/VastData": "3.0.0",
                },
            }),
        )
        assert get_vastdata_provider_version("/some/dir") == "3.0.0"


# ---------------------------------------------------------------------------
# TestVerifyMigrateModeSupport
# ---------------------------------------------------------------------------

def _version_mock(version_str: str):
    """Return a mock for `terraform version -json` with the given vastdata version."""
    return MagicMock(
        returncode=0,
        stdout=json.dumps({
            "terraform_version": "1.14.7",
            "provider_selections": {
                "registry.terraform.io/vast-data/vastdata": version_str,
            },
        }),
    )


def _version_mock_empty():
    """Return a mock for `terraform version -json` that reports no vastdata provider."""
    return MagicMock(
        returncode=0,
        stdout=json.dumps({
            "terraform_version": "1.14.7",
            "provider_selections": {},
        }),
    )


class TestVerifyMigrateModeSupport:
    """Tests for verify_migrate_mode_support()."""

    @patch("subprocess.run")
    def test_v3_provider_passes_without_plan(self, mock_run):
        """v3.x detected via terraform version -json → passes immediately, no plan run."""
        mock_run.return_value = _version_mock("3.1.1")

        verify_migrate_mode_support("/some/dir")

        # Only one subprocess call (terraform version -json), no terraform plan
        mock_run.assert_called_once()
        args, _ = mock_run.call_args
        assert args[0] == ["terraform", "version", "-json"]

    @patch("subprocess.run")
    def test_v3_0_0_passes(self, mock_run):
        """v3.0.0 is the minimum supported version."""
        mock_run.return_value = _version_mock("3.0.0")
        verify_migrate_mode_support("/some/dir")
        mock_run.assert_called_once()

    @patch("subprocess.run")
    def test_old_provider_v2_exits(self, mock_run):
        """v2.x detected via terraform version -json → exits immediately, no plan run."""
        mock_run.return_value = _version_mock("2.9.9")

        with pytest.raises(SystemExit):
            verify_migrate_mode_support("/some/dir")

        mock_run.assert_called_once()

    @patch("subprocess.run")
    def test_version_not_available_falls_back_to_banner(self, mock_run):
        """When terraform version -json returns no vastdata entry, fall back to plan banner."""
        plan_mock = MagicMock(
            returncode=0,
            stdout="",
            stderr="VASTDATA MIGRATE MODE ENABLED\n",
        )
        mock_run.side_effect = [_version_mock_empty(), plan_mock]

        verify_migrate_mode_support("/some/dir")

        assert mock_run.call_count == 2
        plan_call_args, plan_call_kwargs = mock_run.call_args
        assert plan_call_args[0] == ["terraform", "plan", "-input=false"]
        assert plan_call_kwargs["env"]["VASTDATA_MIGRATE_MODE"] == "1"
        assert plan_call_kwargs["env"]["TF_LOG"] == "WARN"

    @patch("subprocess.run")
    def test_version_command_fails_falls_back_to_banner(self, mock_run):
        """When terraform version -json fails, fall back to plan banner check."""
        version_fail = MagicMock(returncode=1, stdout="", stderr="error")
        plan_mock = MagicMock(
            returncode=0,
            stdout="VASTDATA MIGRATE MODE ENABLED",
            stderr="",
        )
        mock_run.side_effect = [version_fail, plan_mock]

        verify_migrate_mode_support("/some/dir")

        assert mock_run.call_count == 2

    @patch("subprocess.run")
    def test_fallback_banner_absent_exits(self, mock_run):
        """Fallback: no banner in plan output → abort."""
        plan_mock = MagicMock(
            returncode=0,
            stdout="No changes.\n",
            stderr="",
        )
        mock_run.side_effect = [_version_mock_empty(), plan_mock]

        with pytest.raises(SystemExit):
            verify_migrate_mode_support("/some/dir")

    @patch("subprocess.run")
    def test_fallback_plan_fails_but_banner_present(self, mock_run):
        """Fallback: plan exits non-zero but banner is present → provider is OK."""
        plan_mock = MagicMock(
            returncode=1,
            stdout="",
            stderr="VASTDATA MIGRATE MODE ENABLED\nError: something else\n",
        )
        mock_run.side_effect = [_version_mock_empty(), plan_mock]

        verify_migrate_mode_support("/some/dir")

    @patch("subprocess.run")
    def test_fallback_plan_fails_no_banner_exits(self, mock_run):
        """Fallback: plan fails and no banner → old provider → abort."""
        plan_mock = MagicMock(
            returncode=1,
            stdout="",
            stderr="Error: something went wrong\n",
        )
        mock_run.side_effect = [_version_mock_empty(), plan_mock]

        with pytest.raises(SystemExit):
            verify_migrate_mode_support("/some/dir")


# ---------------------------------------------------------------------------
# TestMainCLI
# ---------------------------------------------------------------------------

class TestMainCLI:
    """Tests for the main() CLI entry-point."""

    def _create_workspace(self, tmp_path, state_resources, *, serial=5):
        """Create a workspace with .tf files and a state file."""
        workdir = tmp_path / "workdir"
        workdir.mkdir()
        (workdir / "main.tf").write_text('resource "vastdata_tenant" "t" { name = "x" }\n')

        state = _make_state(state_resources, serial=serial)
        state_file = tmp_path / "original.tfstate"
        _write_state_file(state_file, state)

        return workdir, state_file

    @patch("state_migration.run_terraform_apply_migrate")
    @patch("state_migration.verify_migrate_mode_support")
    @patch("state_migration.run_terraform_init")
    def test_full_run(self, mock_init, mock_verify, mock_apply, tmp_path, monkeypatch):
        workdir, state_file = self._create_workspace(tmp_path, [
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ])

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file)])

        main()

        mock_init.assert_called_once_with(str(workdir))
        mock_verify.assert_called_once_with(str(workdir))
        mock_apply.assert_called_once_with(str(workdir))

        # Verify cleaned state was written
        dest = workdir / "terraform.tfstate"
        assert dest.exists()
        with open(dest) as fh:
            cleaned = json.load(fh)
        assert len(cleaned["resources"]) == 1
        assert cleaned["resources"][0]["type"] == "aws_s3_bucket"

    @patch("state_migration.run_terraform_apply_migrate")
    @patch("state_migration.run_terraform_init")
    def test_dry_run(self, mock_init, mock_apply, tmp_path, monkeypatch):
        workdir, state_file = self._create_workspace(tmp_path, [
            _managed("vastdata_tenant", "t1", {"id": 1}),
        ])

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file), "--dry-run"])

        with pytest.raises(SystemExit) as exc_info:
            main()

        # Dry-run calls sys.exit(0)
        assert exc_info.value.code == 0

        # terraform should NOT have been called
        mock_init.assert_not_called()
        mock_apply.assert_not_called()

        # But the cleaned state should have been written
        dest = workdir / "terraform.tfstate"
        assert dest.exists()

    @patch("state_migration.run_terraform_apply_migrate")
    @patch("state_migration.run_terraform_init")
    def test_no_vast_resources_exits_zero(self, mock_init, mock_apply, tmp_path, monkeypatch):
        workdir, state_file = self._create_workspace(tmp_path, [
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ])

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file)])

        with pytest.raises(SystemExit) as exc_info:
            main()

        assert exc_info.value.code == 0
        mock_init.assert_not_called()
        mock_apply.assert_not_called()

    def test_missing_state_file_exits(self, tmp_path, monkeypatch):
        workdir = tmp_path / "workdir"
        workdir.mkdir()
        (workdir / "main.tf").write_text('resource "vastdata_tenant" "t" {}\n')

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", "/no/such/file.tfstate"])

        with pytest.raises(SystemExit):
            main()

    def test_no_tf_files_exits(self, tmp_path, monkeypatch):
        workdir = tmp_path / "empty_workdir"
        workdir.mkdir()

        state = _make_state([_managed("vastdata_tenant", "t1", {"id": 1})])
        state_file = tmp_path / "s.tfstate"
        _write_state_file(state_file, state)

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file)])

        with pytest.raises(SystemExit):
            main()

    @patch("state_migration.run_terraform_apply_migrate")
    @patch("state_migration.verify_migrate_mode_support")
    @patch("state_migration.run_terraform_init")
    def test_existing_state_in_workdir_is_backed_up(self, mock_init, mock_verify, mock_apply, tmp_path, monkeypatch):
        workdir, state_file = self._create_workspace(tmp_path, [
            _managed("vastdata_tenant", "t1", {"id": 1}),
        ])
        # Pre-existing state in workdir
        existing = workdir / "terraform.tfstate"
        existing.write_text('{"pre_existing": true}')

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file)])

        main()

        backup = workdir / "terraform.tfstate.backup-pre-migration"
        assert backup.exists()
        assert json.loads(backup.read_text()) == {"pre_existing": True}


# ---------------------------------------------------------------------------
# TestFixtureFiles
# ---------------------------------------------------------------------------

class TestFixtureFiles:
    """Tests using the checked-in fixture state files."""

    FIXTURES_DIR = Path(__file__).parent / "fixtures" / "state"

    def test_empty_state_fixture(self):
        f = self.FIXTURES_DIR / "empty_state.tfstate"
        if not f.exists():
            pytest.skip("Fixture not found")

        state = parse_tfstate(str(f))
        _, vast_count, non_vast_count = strip_vast_resources(state)

        assert vast_count == 0
        assert non_vast_count == 0

    def test_production_state_fixture(self):
        f = self.FIXTURES_DIR / "sample_production_state.tfstate"
        if not f.exists():
            pytest.skip("Fixture not found")

        state = parse_tfstate(str(f))
        cleaned, vast_count, non_vast_count = strip_vast_resources(state)

        # The sample has 6 managed vastdata resources + 1 data source = 7 vast
        assert vast_count >= 5
        assert non_vast_count == 0  # no non-vast in this fixture
        assert cleaned["resources"] == []


if __name__ == "__main__":
    pytest.main([__file__, "-v"])

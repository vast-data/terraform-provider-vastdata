# Copyright (c) HashiCorp, Inc.

"""TERF-277: state_migration must not terraform import vastdata_* data sources."""

import json
import os
import sys
from unittest.mock import patch

import pytest

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from state_migration import (
    extract_resource_ids,
    main,
    separate_vast_resources,
    strip_vast_resources,
)


def _make_state(resources, *, serial=1):
    return {
        "version": 4,
        "terraform_version": "1.15.8",
        "serial": serial,
        "lineage": "test-lineage",
        "outputs": {},
        "resources": resources,
    }


def _managed(rtype, name, attrs, *, module=None):
    entry = {
        "mode": "managed",
        "type": rtype,
        "name": name,
        "instances": [{"schema_version": 0, "attributes": attrs}],
    }
    if module:
        entry["module"] = module
    return entry


def _data(rtype, name, attrs, *, module=None):
    entry = {
        "mode": "data",
        "type": rtype,
        "name": name,
        "instances": [{"schema_version": 0, "attributes": attrs}],
    }
    if module:
        entry["module"] = module
    return entry


def _write_state(path, state):
    with open(path, "w") as fh:
        json.dump(state, fh, indent=2)


class TestDataSourceStateMigration:
    """Regression tests for TERF-277."""

    def test_extract_resource_ids_skips_vastdata_data_sources(self):
        """Repro from TERF-277: data.vastdata_tenant must not be imported."""
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
        assert "vastdata_tenant.default" not in result
        assert "data.vastdata_tenant.default" not in result

    def test_strip_vast_resources_preserves_vastdata_data_sources(self):
        state = _make_state([
            _data("vastdata_tenant", "default", {"id": 1, "name": "default"}),
            _managed(
                "vastdata_view_policy",
                "nfs",
                {"id": 5, "name": "tfds-policy", "tenant_name": "default"},
            ),
        ])

        cleaned, importable_count, preserved_count, non_vast_count = strip_vast_resources(state)

        assert importable_count == 1
        assert preserved_count == 1
        assert non_vast_count == 0
        assert len(cleaned["resources"]) == 1
        assert cleaned["resources"][0]["mode"] == "data"
        assert cleaned["resources"][0]["type"] == "vastdata_tenant"
        assert cleaned["resources"][0]["name"] == "default"

    def test_separate_vast_resources_buckets_data_sources_as_preserved(self):
        resources = [
            _data("vastdata_tenant", "default", {"id": 1}),
            _managed("vastdata_view_policy", "nfs", {"id": 5}),
            _data("aws_ami", "ubuntu", {"id": "ami-123"}),
        ]

        importable, preserved, non_vast = separate_vast_resources(resources)

        assert len(importable) == 1
        assert importable[0]["name"] == "nfs"
        assert len(preserved) == 1
        assert preserved[0]["type"] == "vastdata_tenant"
        assert len(non_vast) == 1
        assert non_vast[0]["type"] == "aws_ami"

    def test_many_vastdata_view_data_sources_are_preserved_not_imported(self):
        """MANGROUP configs may contain many data.vastdata_view blocks."""
        data_views = [
            _data("vastdata_view", f"view_{i}", {"id": i, "path": f"/path/{i}"})
            for i in range(1, 14)
        ]
        state = _make_state(data_views + [
            _managed("vastdata_view_policy", "nfs", {"id": 99, "name": "policy"}),
        ])

        result = extract_resource_ids(state)
        cleaned, importable_count, preserved_count, _ = strip_vast_resources(state)

        assert result == {"vastdata_view_policy.nfs": "99"}
        assert importable_count == 1
        assert preserved_count == 13
        assert len(cleaned["resources"]) == 13
        assert all(r["mode"] == "data" for r in cleaned["resources"])

    @patch("state_migration.prettify_json_fields")
    @patch("state_migration.patch_imported_attributes")
    @patch("state_migration.run_terraform_import_all")
    @patch("state_migration.run_terraform_init")
    def test_main_imports_only_managed_resources(
        self,
        mock_init,
        mock_import,
        mock_patch,
        mock_prettify,
        tmp_path,
        monkeypatch,
    ):
        workdir = tmp_path / "workdir"
        workdir.mkdir()
        (workdir / "main.tf").write_text(
            'data "vastdata_tenant" "default" { name = "default" }\n'
            'resource "vastdata_view_policy" "nfs" { name = "tfds-policy" }\n'
        )

        state = _make_state([
            _data("vastdata_tenant", "default", {"id": 1, "name": "default"}),
            _managed(
                "vastdata_view_policy",
                "nfs",
                {"id": 5, "name": "tfds-policy", "tenant_name": "default"},
            ),
        ])
        state_file = tmp_path / "terraform.tfstate.v1"
        _write_state(state_file, state)

        monkeypatch.chdir(workdir)
        monkeypatch.setattr("sys.argv", ["state_migration.py", str(state_file)])

        main()

        mock_init.assert_called_once_with(str(workdir))
        mock_import.assert_called_once()
        imported = mock_import.call_args[0][1]
        assert imported == {"vastdata_view_policy.nfs": "5"}

        with open(workdir / "terraform.tfstate") as fh:
            cleaned = json.load(fh)
        assert len(cleaned["resources"]) == 1
        assert cleaned["resources"][0]["mode"] == "data"

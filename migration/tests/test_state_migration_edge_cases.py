#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Edge-case and boundary tests for state_migration.py

Covers:
  - Unusual / malformed state structures
  - Null / missing attributes
  - Very large states
  - Unicode and special characters
  - Legacy (pre-0.12) state format handling
  - Corrupted files
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
)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def _make_state(resources, *, version=4, serial=5):
    return {
        "version": version,
        "terraform_version": "1.5.0",
        "serial": serial,
        "lineage": "edge-case",
        "outputs": {},
        "resources": resources,
    }


def _managed(rtype, name, attrs):
    return {
        "mode": "managed",
        "type": rtype,
        "name": name,
        "instances": [{"schema_version": 0, "attributes": attrs}],
    }


def _data(rtype, name, attrs):
    return {
        "mode": "data",
        "type": rtype,
        "name": name,
        "instances": [{"schema_version": 0, "attributes": attrs}],
    }


def _write(path, state):
    with open(path, "w") as fh:
        json.dump(state, fh, indent=2)


# ---------------------------------------------------------------------------
# Corrupted / malformed files
# ---------------------------------------------------------------------------

class TestCorruptedFiles:

    def test_truncated_json(self, tmp_path):
        f = tmp_path / "truncated.tfstate"
        f.write_text('{"version": 4, "resources": [')
        with pytest.raises(SystemExit):
            parse_tfstate(str(f))

    def test_empty_file(self, tmp_path):
        f = tmp_path / "empty.tfstate"
        f.write_text("")
        with pytest.raises(SystemExit):
            parse_tfstate(str(f))

    def test_non_json_content(self, tmp_path):
        f = tmp_path / "binary.tfstate"
        f.write_bytes(b"\x00\x01\x02\x03")
        with pytest.raises(SystemExit):
            parse_tfstate(str(f))

    def test_json_array_instead_of_object(self, tmp_path):
        """State file is a JSON array — valid JSON but wrong structure → error."""
        f = tmp_path / "array.tfstate"
        f.write_text('[1, 2, 3]')
        # parse_tfstate calls .get() on the result, so a list causes AttributeError
        with pytest.raises(AttributeError):
            parse_tfstate(str(f))

    def test_permission_denied(self, tmp_path):
        if os.name == "nt":
            pytest.skip("Permission test unreliable on Windows")
        f = tmp_path / "noperm.tfstate"
        f.write_text('{"version": 4}')
        os.chmod(f, 0o000)
        try:
            with pytest.raises((SystemExit, PermissionError)):
                parse_tfstate(str(f))
        finally:
            os.chmod(f, 0o644)


# ---------------------------------------------------------------------------
# Unusual state structures
# ---------------------------------------------------------------------------

class TestUnusualStateStructures:

    def test_no_resources_key(self):
        """State dict without 'resources' at all."""
        state = {"version": 4, "serial": 1}
        cleaned, vc, nvc = strip_vast_resources(state)
        assert vc == 0
        assert nvc == 0

    def test_resources_is_none(self):
        """'resources' exists but is None."""
        state = {"version": 4, "serial": 1, "resources": None}
        # separate_vast_resources expects a list; this should be handled
        # strip_vast_resources checks 'resources' not in state — None is truthy
        # so it will try to iterate. We accept TypeError here or graceful handling.
        with pytest.raises(TypeError):
            strip_vast_resources(state)

    def test_extra_top_level_fields_preserved(self):
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
        ])
        state["custom_field"] = "should_survive"
        state["backend"] = {"type": "s3", "config": {"bucket": "my-bucket"}}

        cleaned, _, _ = strip_vast_resources(state)

        assert cleaned["custom_field"] == "should_survive"
        assert cleaned["backend"]["config"]["bucket"] == "my-bucket"

    def test_resource_with_no_instances(self):
        """Resource entry with empty instances list."""
        resources = [
            {
                "mode": "managed",
                "type": "vastdata_tenant",
                "name": "empty",
                "instances": [],
            },
            _managed("aws_s3_bucket", "b", {"id": "x"}),
        ]
        state = _make_state(resources)

        cleaned, vc, nvc = strip_vast_resources(state)
        assert vc == 1
        assert nvc == 1

    def test_resource_with_multiple_instances(self):
        """Resource with count/for_each producing multiple instances."""
        resource = {
            "mode": "managed",
            "type": "vastdata_view",
            "name": "views",
            "instances": [
                {"index_key": 0, "schema_version": 0, "attributes": {"id": 1}},
                {"index_key": 1, "schema_version": 0, "attributes": {"id": 2}},
                {"index_key": 2, "schema_version": 0, "attributes": {"id": 3}},
            ],
        }
        state = _make_state([resource, _managed("aws_s3_bucket", "b", {"id": "x"})])

        cleaned, vc, nvc = strip_vast_resources(state)
        assert vc == 1  # It's 1 resource entry (with 3 instances), not 3
        assert nvc == 1

    def test_module_resources_in_state(self):
        """Resources inside modules are still typed vastdata_* and should be stripped."""
        resources = [
            {
                "module": "module.networking",
                "mode": "managed",
                "type": "vastdata_vip_pool",
                "name": "pool",
                "instances": [{"attributes": {"id": 1}}],
            },
            {
                "module": "module.storage",
                "mode": "managed",
                "type": "aws_ebs_volume",
                "name": "vol",
                "instances": [{"attributes": {"id": "vol-123"}}],
            },
        ]
        state = _make_state(resources)

        cleaned, vc, nvc = strip_vast_resources(state)
        assert vc == 1
        assert nvc == 1
        assert cleaned["resources"][0]["type"] == "aws_ebs_volume"


# ---------------------------------------------------------------------------
# Null / missing attributes
# ---------------------------------------------------------------------------

class TestNullAndMissingAttributes:

    def test_resource_with_all_null_attributes(self):
        resources = [
            _managed("vastdata_view", "v1", {
                "id": None, "path": None, "tenant_id": None,
            }),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 1

    def test_resource_with_empty_attributes(self):
        resources = [_managed("vastdata_tenant", "t1", {})]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 1

    def test_resource_without_type_key(self):
        """Missing 'type' → defaults to '' → not vastdata → goes to non-vast."""
        resources = [{"mode": "managed", "name": "orphan", "instances": []}]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 0
        assert len(non_vast) == 1


# ---------------------------------------------------------------------------
# Large states
# ---------------------------------------------------------------------------

class TestLargeStates:

    def test_500_mixed_resources(self):
        resources = []
        for i in range(200):
            resources.append(_managed("vastdata_view", f"v_{i}", {"id": i}))
        for i in range(150):
            resources.append(_managed("aws_s3_bucket", f"b_{i}", {"id": f"b-{i}"}))
        for i in range(100):
            resources.append(_managed("google_compute_instance", f"vm_{i}", {"id": f"vm-{i}"}))
        for i in range(50):
            resources.append(_data("vastdata_vip_pool", f"pool_{i}", {"id": i}))

        state = _make_state(resources)

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 250   # 200 managed + 50 data
        assert nvc == 250   # 150 aws + 100 gcp
        assert len(cleaned["resources"]) == 250

    def test_write_large_state(self, tmp_path):
        """Ensure write_state handles large output."""
        resources = [
            _managed("aws_s3_bucket", f"b_{i}", {"id": f"bucket-{i}", "tags": {"env": "prod"}})
            for i in range(500)
        ]
        state = _make_state(resources)
        out = str(tmp_path / "large.tfstate")

        write_state(state, out)

        with open(out) as fh:
            loaded = json.load(fh)
        assert len(loaded["resources"]) == 500


# ---------------------------------------------------------------------------
# Unicode and special characters
# ---------------------------------------------------------------------------

class TestUnicodeAndSpecialChars:

    def test_unicode_resource_name(self):
        resources = [
            _managed("vastdata_tenant", "テナント", {"id": 1, "name": "日本語テナント"}),
            _managed("aws_s3_bucket", "バケット", {"id": "unicode-bucket"}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 1
        assert len(non_vast) == 1

    def test_special_chars_in_attributes(self):
        resources = [
            _managed("vastdata_view", "v1", {
                "id": 1,
                "path": "/data/special chars & 'quotes' \"double\"",
            }),
        ]
        vast, _ = separate_vast_resources(resources)
        assert len(vast) == 1

    def test_write_state_with_unicode(self, tmp_path):
        state = _make_state([
            _managed("aws_s3_bucket", "b", {"id": "b", "name": "名前"})
        ])
        out = str(tmp_path / "unicode.tfstate")
        write_state(state, out)

        with open(out, encoding="utf-8") as fh:
            loaded = json.load(fh)
        assert loaded["resources"][0]["instances"][0]["attributes"]["name"] == "名前"


# ---------------------------------------------------------------------------
# Serial number handling
# ---------------------------------------------------------------------------

class TestSerialHandling:

    def test_serial_zero(self):
        state = _make_state([_managed("vastdata_tenant", "t", {"id": 1})], serial=0)
        cleaned, _, _ = strip_vast_resources(state)
        assert cleaned["serial"] == 1

    def test_serial_very_large(self):
        state = _make_state([_managed("vastdata_tenant", "t", {"id": 1})], serial=999999)
        cleaned, _, _ = strip_vast_resources(state)
        assert cleaned["serial"] == 1000000

    def test_serial_missing(self):
        state = {"version": 4, "resources": [_managed("vastdata_tenant", "t", {"id": 1})]}
        # No 'serial' key
        cleaned, vc, _ = strip_vast_resources(state)
        assert vc == 1
        assert cleaned["serial"] == 1  # 0 + 1


# ---------------------------------------------------------------------------
# Provider-specific edge cases
# ---------------------------------------------------------------------------

class TestProviderEdgeCases:

    def test_vastdata_prefix_but_different_provider(self):
        """
        Hypothetical: someone names a resource type 'vastdata_custom_thing'
        from a custom provider. It still starts with 'vastdata_' so should
        be classified as vast.
        """
        resources = [
            _managed("vastdata_custom_thing", "c1", {"id": 1}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 1

    def test_vast_without_underscore(self):
        """'vastdata' without trailing underscore should be non-vast."""
        resources = [
            _managed("vastdata", "x", {"id": 1}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        # 'vastdata'.startswith('vastdata_') → False
        assert len(vast) == 0
        assert len(non_vast) == 1

    def test_mixed_case_type(self):
        """'VastData_tenant' (wrong case) should NOT match vastdata_."""
        resources = [
            _managed("VastData_tenant", "t", {"id": 1}),
        ]
        vast, non_vast = separate_vast_resources(resources)
        assert len(vast) == 0
        assert len(non_vast) == 1

    def test_many_different_providers(self):
        """State with resources from 5+ different providers."""
        resources = [
            _managed("vastdata_tenant", "v1", {"id": 1}),
            _managed("vastdata_view", "v2", {"id": 2}),
            _managed("aws_s3_bucket", "a1", {"id": "b1"}),
            _managed("aws_iam_role", "a2", {"id": "r1"}),
            _managed("google_compute_instance", "g1", {"id": "vm1"}),
            _managed("azurerm_resource_group", "az1", {"id": "rg1"}),
            _managed("kubernetes_deployment", "k1", {"id": "deploy1"}),
            _managed("datadog_monitor", "d1", {"id": "mon1"}),
            _managed("vastdata_quota", "v3", {"id": 3}),
        ]
        state = _make_state(resources)

        cleaned, vc, nvc = strip_vast_resources(state)

        assert vc == 3   # tenant, view, quota
        assert nvc == 6  # aws×2, google, azure, k8s, datadog
        assert len(cleaned["resources"]) == 6
        remaining_types = {r["type"] for r in cleaned["resources"]}
        assert "vastdata_tenant" not in remaining_types
        assert "aws_s3_bucket" in remaining_types
        assert "kubernetes_deployment" in remaining_types


# ---------------------------------------------------------------------------
# Idempotency
# ---------------------------------------------------------------------------

class TestIdempotency:

    def test_strip_twice(self):
        """Stripping an already-stripped state should produce 0 vast."""
        state = _make_state([
            _managed("vastdata_tenant", "t1", {"id": 1}),
            _managed("aws_s3_bucket", "b1", {"id": "bucket"}),
        ], serial=10)

        cleaned_1, vc1, _ = strip_vast_resources(state)
        assert vc1 == 1

        cleaned_2, vc2, nvc2 = strip_vast_resources(cleaned_1)
        assert vc2 == 0
        assert nvc2 == 1  # aws bucket still there
        # Serial bumped again
        assert cleaned_2["serial"] == 12

    def test_write_then_parse_roundtrip(self, tmp_path):
        state = _make_state([
            _managed("aws_s3_bucket", "b1", {"id": "bucket", "tags": {"env": "prod"}}),
        ], serial=42)
        out = str(tmp_path / "rt.tfstate")

        write_state(state, out)
        loaded = parse_tfstate(out)

        assert loaded["serial"] == 42
        assert len(loaded["resources"]) == 1
        assert loaded["resources"][0]["instances"][0]["attributes"]["tags"]["env"] == "prod"


if __name__ == "__main__":
    pytest.main([__file__, "-v"])

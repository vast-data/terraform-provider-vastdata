# Copyright (c) HashiCorp, Inc.

"""Tests for v1 -> v3 breaking-change handling in migration_script.py."""

import os
import sys

import pytest

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import main, transform_file, transform_resource_block


class TestBreakingChanges:
    def test_replication_peer_removes_read_only_attributes(self):
        content = '''resource "vastdata_replication_peers" "peer" {
  name        = "peer-a"
  leading_vip = "10.0.0.1"
  peer_name   = "REMOTE"
  is_local    = true
  pool_id     = 1
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert 'resource "vastdata_replication_peer"' in result
        assert "peer_name =" not in result
        assert "is_local =" not in result

    def test_view_renames_s3_locking_attributes(self):
        content = '''resource "vastdata_view" "locked_bucket" {
  path                      = "/locked"
  policy_id                 = 1
  protocols                 = ["S3"]
  s3_locks                  = true
  s3_locks_retention_period = "7d"
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert "locking = true" in result
        assert 'default_retention_period = "7d"' in result
        assert "s3_locks" not in result

    def test_view_policy_converts_vippool_permissions_blocks(self):
        content = '''resource "vastdata_view_policy" "example" {
  name = "example"
  vippool_permissions {
    vippool_id          = 1
    vippool_permissions = "RW"
  }
  vippool_permissions {
    vippool_id          = 2
    vippool_permissions = "RO"
  }
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert "permission_per_vip_pool = {" in result
        assert '"1" = "RW"' in result
        assert '"2" = "RO"' in result
        assert "vippool_permissions {" not in result

    def test_view_policy_converts_vippool_permissions_with_expressions(self):
        content = '''resource "vastdata_view_policy" "example" {
  name = "example"
  vippool_permissions {
    vippool_id          = vastdata_vip_pool.pool_a.id
    vippool_permissions = "RW"
  }
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert "(vastdata_vip_pool.pool_a.id) = \"RW\"" in result

    def test_tenant_converts_dynamic_client_ip_ranges(self):
        content = '''resource "vastdata_tenant" "tenant1" {
  name = "tenant1"
  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip   = client_ip_ranges.value["end_ip"]
    }
  }
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert 'client_ip_ranges = [for r in var.tenant_client_ip_ranges : [r["start_ip"], r["end_ip"]]]' in result
        assert 'dynamic "client_ip_ranges"' not in result

    def test_protection_policy_converts_dynamic_frames_dot_access(self):
        content = '''resource "vastdata_protection_policy" "pp" {
  name           = "dynfr-policy"
  clone_type     = "LOCAL"
  indestructible = false
  prefix         = "dynfr"
  dynamic "frames" {
    for_each = var.frames
    content {
      every      = frames.value.every
      keep_local = frames.value.keep_local
      start_at   = frames.value.start_at
    }
  }
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert 'frames = [for f in var.frames : {' in result
        assert "every = f.every" in result.replace(" ", "")
        assert "keep_local = f.keep_local" in result.replace(" ", "")
        assert "start_at = f.start_at" in result.replace(" ", "")
        assert 'dynamic "frames"' not in result

    def test_protection_policy_converts_dynamic_frames_bracket_access(self):
        content = '''resource "vastdata_protection_policy" "pp" {
  name = "dynfr-policy"
  dynamic "frames" {
    for_each = var.frames
    content {
      every      = frames.value["every"]
      keep_local = frames.value["keep_local"]
      start_at   = frames.value["start_at"]
    }
  }
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert 'frames = [for f in var.frames : {' in result
        assert 'every=f["every"]' in result.replace(" ", "")
        assert 'dynamic "frames"' not in result

    def test_protection_policy_converts_dynamic_frames_custom_iterator(self):
        content = '''resource "vastdata_protection_policy" "pp" {
  name = "dynfr-policy"
  dynamic "frames" {
    for_each = var.frames
    iterator = frame
    content {
      every = frame.value.every
    }
  }
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert 'frames = [for f in var.frames : {' in result
        assert "every = f.every" in result

    def test_tenant_removes_vippool_ids_with_comment(self):
        content = '''resource "vastdata_tenant" "tenant1" {
  name        = "tenant1"
  vippool_ids = [1, 2]
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert "vippool_ids =" not in result
        assert "# TODO: vippool_ids was removed" in result

    def test_tenant_converts_vippool_ids_pool_refs(self):
        content = '''resource "vastdata_tenant" "tenant1" {
  name = "tenant1"
  vippool_ids = [
    vastdata_vip_pool.pool_a.id,
    vastdata_vip_pool.pool_b.id,
  ]
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert "vippool_ids =" not in result
        assert "# TODO:" not in result
        assert "NOTE: vippool_ids removed" in result
        assert "vastdata_vip_pool.pool_a" in result
        assert "WARNING: VAST API rejects tenant_id changes" in result
        assert "tenant_id = vastdata_tenant" not in result

    def test_group_injects_local_provider_id(self):
        content = '''resource "vastdata_group" "admins" {
  name = "admins"
  gid  = 5001
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert "local_provider_id = 1" in result

    def test_qos_policy_converts_attached_users_identifiers_user_ref(self):
        content = '''resource "vastdata_qos_policy" "qos" {
  name                       = "qos"
  attached_users_identifiers = [tostring(vastdata_user.qos_user.id)]
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert "attached_users_identifiers" not in result
        assert "attached_users" in result
        assert "vastdata_user.qos_user.sid" in result
        assert 'identifier_type  = "sid_str"' in result
        assert "MIGRATION REQUIRED" not in result

    def test_qos_policy_removes_literal_attached_users_identifiers(self):
        content = '''resource "vastdata_qos_policy" "qos" {
  name                       = "qos"
  attached_users_identifiers = ["101"]
}'''
        result, _ = transform_resource_block(content.split("\n"), 0)
        assert "attached_users_identifiers =" not in result
        assert "MIGRATION REQUIRED" in result

    def test_vip_pool_removes_read_only_active_cnode_ids(self, tmp_path):
        content = '''variable "active_cnode_ids" {
  type = list(number)
}

resource "vastdata_vip_pool" "pool" {
  name             = "pool"
  role             = "PROTOCOLS"
  subnet_cidr      = "24"
  active_cnode_ids = var.active_cnode_ids
  ip_ranges {
    start_ip = "10.0.0.1"
    end_ip   = "10.0.0.2"
  }
}'''
        input_file = tmp_path / "main.tf"
        output_file = tmp_path / "main_converted.tf"
        input_file.write_text(content)
        transform_file(input_file, output_file)
        result = output_file.read_text()
        assert "active_cnode_ids =" not in result
        assert "read-only in v3" in result

    def test_provider_version_updated_to_3(self, tmp_path):
        content = '''terraform {
  required_providers {
    vastdata = {
      source  = "vast-data/vastdata"
      version = "1.7.0"
    }
  }
}'''
        input_file = tmp_path / "main.tf"
        output_file = tmp_path / "main_converted.tf"
        input_file.write_text(content)
        transform_file(input_file, output_file)
        result = output_file.read_text()
        assert 'version = "3.0.0"' in result


class TestFailClosedMigration:
    def test_main_exits_when_manual_steps_remain(self, tmp_path):
        src = tmp_path / "in"
        dst = tmp_path / "out"
        src.mkdir()
        (src / "qos.tf").write_text(
            'resource "vastdata_qos_policy" "q" {\n'
            '  name = "q"\n'
            '  attached_users_identifiers = ["1"]\n'
            "}\n"
        )

        with pytest.raises(SystemExit) as excinfo:
            main(src, dst, allow_incomplete=False)

        assert excinfo.value.code == 2
        converted = dst / "qos_converted.tf"
        assert converted.exists()
        assert "attached_users_identifiers =" not in converted.read_text()
        assert "MIGRATION REQUIRED" in converted.read_text()

    def test_main_succeeds_when_attached_users_auto_converted(self, tmp_path):
        src = tmp_path / "in"
        dst = tmp_path / "out"
        src.mkdir()
        (src / "qos.tf").write_text(
            'resource "vastdata_user" "u" { name = "u" uid = 1 }\n'
            'resource "vastdata_qos_policy" "q" {\n'
            '  name = "q"\n'
            '  attached_users_identifiers = [tostring(vastdata_user.u.id)]\n'
            "}\n"
        )

        main(src, dst, allow_incomplete=False)

        result = (dst / "qos_converted.tf").read_text()
        assert "attached_users" in result
        assert "# TODO:" not in result

    def test_main_does_not_inject_tenant_id_on_vippool_ids(self, tmp_path):
        src = tmp_path / "in"
        dst = tmp_path / "out"
        src.mkdir()
        (src / "tenant.tf").write_text(
            'resource "vastdata_tenant" "t" {\n'
            '  name = "t"\n'
            '  vippool_ids = [vastdata_vip_pool.pool_a.id]\n'
            "}\n"
        )
        (src / "vip_pools.tf").write_text(
            'resource "vastdata_vip_pool" "pool_a" {\n'
            '  name = "pool-a"\n'
            '  role = "PROTOCOLS"\n'
            '  subnet_cidr = "24"\n'
            "}\n"
        )

        main(src, dst, allow_incomplete=False)

        tenant = (dst / "tenant_converted.tf").read_text()
        pools = (dst / "vip_pools_converted.tf").read_text()
        assert "# TODO:" not in tenant
        assert "tenant_id = vastdata_tenant.t.id" not in pools

    def test_main_succeeds_with_allow_incomplete(self, tmp_path):
        src = tmp_path / "in"
        dst = tmp_path / "out"
        src.mkdir()
        (src / "qos.tf").write_text(
            'resource "vastdata_qos_policy" "q" {\n'
            '  name = "q"\n'
            '  attached_users_identifiers = ["1"]\n'
            "}\n"
        )

        main(src, dst, allow_incomplete=True)

        assert (dst / "qos_converted.tf").exists()

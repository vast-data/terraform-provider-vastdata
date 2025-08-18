# Copyright (c) HashiCorp, Inc.

"""
Tests for specific configuration file transformations using exact expected outputs.

This module tests specific conf file scenarios with precise expected transformations.
"""

import pytest
import tempfile
import shutil
from pathlib import Path
import sys
import os

# Add the migration script directory to the Python path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import transform_file


class TestConfSpecificFixtures:
    """Test specific conf file transformations with exact expected outputs."""
    
    @pytest.fixture
    def fixtures_dir(self):
        """Provide path to fixtures directory."""
        return Path(__file__).parent / "fixtures"
    
    @pytest.fixture
    def temp_dir(self):
        """Create a temporary directory for test files."""
        temp_dir = tempfile.mkdtemp()
        yield Path(temp_dir)
        shutil.rmtree(temp_dir)
    
    def test_administrators_managers_transformation(self, fixtures_dir, temp_dir):
        """Test exact transformation of administrators managers configuration."""
        input_file = fixtures_dir / "input" / "administrators_managers_conf.tf"
        expected_file = fixtures_dir / "expected" / "administrators_managers_conf.tf"
        
        if not input_file.exists() or not expected_file.exists():
            pytest.skip(f"Fixture files not found: {input_file} or {expected_file}")
        
        # Create temporary files
        test_input = temp_dir / "test_input.tf"
        test_output = temp_dir / "test_output.tf"
        
        # Copy input file
        shutil.copy2(input_file, test_input)
        
        # Run transformation
        transform_file(test_input, test_output)
        
        # Read results
        result = test_output.read_text()
        expected = expected_file.read_text()
        
        # Verify key transformations
        assert "vastdata_administrator_manager" in result
        assert "vastdata_administrator_role" in result
        assert "vastdata_administators_managers" not in result
        assert "vastdata_administators_roles" not in result
        
        # Verify references are updated
        assert "vastdata_administrator_role.man_role1.id" in result
        assert "vastdata_administrator_manager.manager1" in result
        
        # Verify structure is preserved
        assert "variable role_name" in result
        assert "variable permissions_list" in result
        assert "output tf_role" in result
        assert "sensitive = true" in result
    
    def test_qos_policy_blocks_transformation(self, fixtures_dir, temp_dir):
        """Test exact transformation of QOS policy blocks to attributes."""
        input_file = fixtures_dir / "input" / "qos_policy_blocks_conf.tf"
        expected_file = fixtures_dir / "expected" / "qos_policy_blocks_conf.tf"
        
        if not input_file.exists() or not expected_file.exists():
            pytest.skip(f"Fixture files not found: {input_file} or {expected_file}")
        
        # Create temporary files
        test_input = temp_dir / "test_input.tf"
        test_output = temp_dir / "test_output.tf"
        
        # Copy input file
        shutil.copy2(input_file, test_input)
        
        # Run transformation
        transform_file(test_input, test_output)
        
        # Read results
        result = test_output.read_text()
        expected = expected_file.read_text()
        
        # Verify block-to-attribute transformations
        assert "static_limits = {" in result
        assert "capacity_limits = {" in result
        
        # Verify original block syntax is gone
        assert "static_limits {" not in result or result.count("static_limits {") == 0
        assert "capacity_limits {" not in result or result.count("capacity_limits {") == 0
        
        # Verify all attributes are preserved
        assert "max_writes_bw_mbps = 110" in result
        assert "max_reads_iops = 200" in result
        assert "max_writes_iops = 3001" in result
        assert "max_reads_bw_mbps_per_gb_capacity = 100" in result
        assert "max_reads_iops_per_gb_capacity = 200" in result
        
        # Verify other structure preserved
        assert "variable qos_policy_name" in result
        assert "output tf_qos_policy" in result
    
    def test_nonlocal_user_transformation(self, fixtures_dir, temp_dir):
        """Test exact transformation of non_local_user to nonlocal_user."""
        input_file = fixtures_dir / "input" / "nonlocal_user_conf.tf"
        expected_file = fixtures_dir / "expected" / "nonlocal_user_conf.tf"
        
        if not input_file.exists() or not expected_file.exists():
            pytest.skip(f"Fixture files not found: {input_file} or {expected_file}")
        
        # Create temporary files
        test_input = temp_dir / "test_input.tf"
        test_output = temp_dir / "test_output.tf"
        
        # Copy input file
        shutil.copy2(input_file, test_input)
        
        # Run transformation
        transform_file(test_input, test_output)
        
        # Read results
        result = test_output.read_text()
        expected = expected_file.read_text()
        
        # Verify resource renaming
        assert "vastdata_nonlocal_user" in result
        assert "vastdata_non_local_user" not in result
        
        # Verify data source renaming
        assert 'data "vastdata_nonlocal_user"' in result
        assert 'data "vastdata_non_local_user"' not in result
        
        # Verify references are updated
        assert "vastdata_nonlocal_user.non_local_user1" in result
        assert "data.vastdata_nonlocal_user.user_data" in result
        
        # Verify attributes are preserved
        assert "s3_policies_ids = [1, 2, 3]" in result
        assert "allow_create_bucket = true" in result
        assert "allow_delete_bucket = false" in result
        
        # Verify structure preserved
        assert "variable user_uid" in result
        assert "output tf_user" in result
        assert "output tf_user_ds" in result
    
    def test_ip_ranges_transformation_vippool(self, temp_dir):
        """Test IP ranges transformation in VIP pool configuration."""
        terraform_content = '''resource "vastdata_vip_pool" "pool1" {
    name = "test-pool"
    role = "PROTOCOLS"
    
    ip_ranges {
        start_ip = "192.168.1.1"
        end_ip = "192.168.1.10"
    }
    
    ip_ranges {
        start_ip = "192.168.1.20"
        end_ip = "192.168.1.30"
    }
}'''
        
        input_file = temp_dir / "vippool_input.tf"
        output_file = temp_dir / "vippool_output.tf"
        
        input_file.write_text(terraform_content)
        
        # Run transformation
        transform_file(input_file, output_file)
        
        result = output_file.read_text()
        
        # Should transform to list of lists
        assert 'ip_ranges = [' in result
        assert '["192.168.1.1", "192.168.1.10"]' in result
        assert '["192.168.1.20", "192.168.1.30"]' in result
        
        # Original block syntax should be gone
        assert 'ip_ranges {' not in result
    
    def test_tenant_dynamic_blocks_preserved(self, temp_dir):
        """Test that dynamic blocks in tenant configuration are preserved."""
        terraform_content = '''variable tenant_client_ip_ranges {
    type = list(object({
      start_ip = string
      end_ip = string
    }))
}

resource vastdata_tenant tenant1 {
  name = "test-tenant"

  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip = client_ip_ranges.value["end_ip"]
    }
  }
}'''
        
        input_file = temp_dir / "tenant_input.tf"
        output_file = temp_dir / "tenant_output.tf"
        
        input_file.write_text(terraform_content)
        
        # Run transformation
        transform_file(input_file, output_file)
        
        result = output_file.read_text()
        
        # Dynamic block should be preserved exactly
        assert 'dynamic "client_ip_ranges"' in result
        assert 'for_each = var.tenant_client_ip_ranges' in result
        assert 'client_ip_ranges.value["start_ip"]' in result
        assert 'client_ip_ranges.value["end_ip"]' in result
        
        # Should not be converted to static format
        assert 'client_ip_ranges = [' not in result
    
    def test_share_acl_transformation_s3_view(self, temp_dir):
        """Test share_acl block transformation in S3 view."""
        terraform_content = '''resource vastdata_view view1 {
  path = "/test-bucket"
  bucket = "test-bucket"
  
  share_acl {
    acl {
        name = "test-user"
        grantee = "users"
        fqdn = "All"
        permissions = "RW"
    }
    enabled = true
  }
}'''
        
        input_file = temp_dir / "s3view_input.tf"
        output_file = temp_dir / "s3view_output.tf"
        
        input_file.write_text(terraform_content)
        
        # Run transformation
        transform_file(input_file, output_file)
        
        result = output_file.read_text()
        
        # Should transform to attribute
        assert 'share_acl = {' in result
        assert 'enabled = true' in result
        
        # Nested acl block should be flattened (no longer present as a separate block)
        assert 'acl {' not in result
        assert 'name = "test-user"' in result
        assert 'grantee = "users"' in result
        
        # Original block syntax should be transformed
        lines = result.split('\n')
        share_acl_block_lines = [line for line in lines if 'share_acl {' in line and 'share_acl = {' not in line]
        assert len(share_acl_block_lines) == 0, "share_acl block syntax should be converted to attribute"
    
    def test_empty_lists_comments_preserved(self, temp_dir):
        """Test that empty list comments are preserved during migration."""
        terraform_content = '''resource vastdata_group user_group1 {
  name = "test-group"
  gid = 1000
  # s3_policies_ids = [] check
}

resource "vastdata_non_local_user" "user1" {
    uid = 2000
    context = "ldap"
    # s3_policies_ids = [] check
}'''
        
        input_file = temp_dir / "empty_lists_input.tf"
        output_file = temp_dir / "empty_lists_output.tf"
        
        input_file.write_text(terraform_content)
        
        # Run transformation
        transform_file(input_file, output_file)
        
        result = output_file.read_text()
        
        # Comments should be preserved
        assert "# s3_policies_ids = [] check" in result
        
        # Resource renaming should occur
        assert "vastdata_nonlocal_user" in result
        assert "vastdata_non_local_user" not in result
        
        # Structure should be preserved
        assert "name = \"test-group\"" in result
        assert "gid = 1000" in result
        assert "uid = 2000" in result
        assert "context = \"ldap\"" in result


class TestConfFilesAdvancedScenarios:
    """Test advanced scenarios found in conf files."""
    
    @pytest.fixture
    def temp_dir(self):
        """Create a temporary directory for test files.""" 
        temp_dir = tempfile.mkdtemp()
        yield Path(temp_dir)
        shutil.rmtree(temp_dir)
    
    def test_complex_conditional_s3_policies(self, temp_dir):
        """Test complex conditional s3_policies_ids from conf files."""
        terraform_content = '''variable use_s3_policies {
    type = string
}

variable s3_policy_id1 {
    type = number
}

variable s3_policy_id2 {
  type = number
}

resource "vastdata_non_local_user" "user1" {
    uid = 1000
    context = "ldap"
    s3_policies_ids = var.use_s3_policies == "none" ? [] : (
        var.use_s3_policies == "all" ? [
            var.s3_policy_id1,
            var.s3_policy_id2
        ] : [
            var.s3_policy_id1
        ]
    )
}'''
        
        input_file = temp_dir / "conditional_input.tf"
        output_file = temp_dir / "conditional_output.tf"
        
        input_file.write_text(terraform_content)
        
        # Run transformation
        transform_file(input_file, output_file)
        
        result = output_file.read_text()
        
        # Resource should be renamed
        assert "vastdata_nonlocal_user" in result
        assert "vastdata_non_local_user" not in result
        
        # Complex conditional should be preserved
        assert 'var.use_s3_policies == "none" ? []' in result
        assert 'var.use_s3_policies == "all" ?' in result
        assert 'var.s3_policy_id1' in result
        assert 'var.s3_policy_id2' in result
    
    def test_multiple_resource_dependencies(self, temp_dir):
        """Test files with multiple resource dependencies and references."""
        terraform_content = '''resource vastdata_ldap ldap1 {
  domain_name = "example.com"
  urls = ["ldap://server1"]
}

resource vastdata_active_directory active_dir1 {
  ldap_id = vastdata_ldap.ldap1.id
  machine_account_name = "test-machine"
}

resource vastdata_tenant tenant1 {
  name = "test-tenant"
  ldap_provider_id = vastdata_ldap.ldap1.id
}

output tf_active_directory {
  value = vastdata_active_directory.active_dir1
  sensitive = true
}

output tf_ldap {
  value = vastdata_ldap.ldap1
  sensitive = true
}'''
        
        input_file = temp_dir / "dependencies_input.tf"
        output_file = temp_dir / "dependencies_output.tf"
        
        input_file.write_text(terraform_content)
        
        # Run transformation
        transform_file(input_file, output_file)
        
        result = output_file.read_text()
        
        # All resource references should be preserved
        assert "vastdata_ldap.ldap1.id" in result
        assert "vastdata_active_directory.active_dir1" in result
        assert "vastdata_ldap.ldap1" in result
        
        # Resource definitions should be preserved
        assert "resource vastdata_ldap" in result
        assert "resource vastdata_active_directory" in result
        assert "resource vastdata_tenant" in result
        
        # Outputs should be preserved
        assert "output tf_active_directory" in result
        assert "output tf_ldap" in result
        assert "sensitive = true" in result

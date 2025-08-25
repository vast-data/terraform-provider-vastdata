# Copyright (c) HashiCorp, Inc.

"""
Tests based on real customer configuration scenarios.

This test module validates the migration script against real-world customer
configurations found in the vast-terraform directory.
"""

import pytest
import tempfile
import shutil
from pathlib import Path
import sys
import os

# Add the migration script directory to the Python path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from migration_script import transform_file, main


class TestRealCustomerScenarios:
    """Test migration script against real customer configuration scenarios."""
    
    @pytest.fixture
    def vast_terraform_dir(self):
        """Provide path to real customer vast-terraform directory."""
        return Path(__file__).parent / "vast-terraform"
    
    @pytest.fixture
    def temp_migration_dirs(self):
        """Create temporary directories for migration testing."""
        src_temp = tempfile.mkdtemp(prefix="migration_customer_src_")
        dst_temp = tempfile.mkdtemp(prefix="migration_customer_dst_") 
        
        yield Path(src_temp), Path(dst_temp)
        
        # Cleanup
        shutil.rmtree(src_temp)
        shutil.rmtree(dst_temp)
    
    def test_administrator_resources_and_references(self, vast_terraform_dir, temp_migration_dirs):
        """Test that administrator resources are renamed and references updated."""
        src_dir, dst_dir = temp_migration_dirs
        
        main_tf = vast_terraform_dir / "main.tf"
        if main_tf.exists():
            # Copy file to source directory
            shutil.copy2(main_tf, src_dir)
            
            # Run migration
            main(str(src_dir), str(dst_dir))
            
            # Check transformation
            converted_file = dst_dir / "main_converted.tf"
            assert converted_file.exists(), "main.tf was not converted"
            
            content = converted_file.read_text()
            
            # Check resource type renaming
            assert 'resource "vastdata_administrator_role"' in content
            assert 'resource "vastdata_administrator_manager"' in content
            assert 'resource "vastdata_administators_roles"' not in content
            assert 'resource "vastdata_administators_managers"' not in content
            
            # Check resource references are updated
            assert 'vastdata_administrator_role.read_only.id' in content
            assert 'vastdata_administrator_role.csi.id' in content
            assert 'vastdata_administators_roles.read_only.id' not in content
            assert 'vastdata_administators_roles.csi.id' not in content
            
            # Check attribute renaming
            assert 'permissions = [' in content
            assert 'permissions_list = [' not in content
    
    def test_vip_pool_cnode_ids_transformation(self, vast_terraform_dir, temp_migration_dirs):
        """Test that cnode_ids lists remain as lists (no conversion to strings)."""
        src_dir, dst_dir = temp_migration_dirs
        
        vip_pools_tf = vast_terraform_dir / "vip_pools.tf"
        if vip_pools_tf.exists():
            # Copy file to source directory
            shutil.copy2(vip_pools_tf, src_dir)
            
            # Run migration
            main(str(src_dir), str(dst_dir))
            
            # Check transformation
            converted_file = dst_dir / "vip_pools_converted.tf"
            assert converted_file.exists(), "vip_pools.tf was not converted"
            
            content = converted_file.read_text()
            
            # Check cnode_ids remain as lists (not converted to strings)
            assert 'cnode_ids         = [' in content or 'cnode_ids = [' in content
            # Verify the lists contain the expected values
            assert '1,' in content and '2,' in content  # from first pool
            assert '3,' in content and '4,' in content  # from second pool
            
            # Check IP ranges are still transformed correctly
            assert 'ip_ranges = [' in content
            assert '["10.66.20.201", "10.66.20.208"]' in content
            assert '["10.66.20.141", "10.66.20.148"]' in content
            assert '["192.168.55.20", "192.168.55.30"]' in content
            
            # Note: cnode_ids should remain as lists (multi-line format preserved)
            # Only active_cnode_ids should be converted to strings if present
    
    def test_view_policy_list_attributes(self, vast_terraform_dir, temp_migration_dirs):
        """Test that view policy list attributes are preserved correctly."""
        src_dir, dst_dir = temp_migration_dirs
        
        view_policies_tf = vast_terraform_dir / "view_policies.tf"
        if view_policies_tf.exists():
            # Copy file to source directory
            shutil.copy2(view_policies_tf, src_dir)
            
            # Run migration
            main(str(src_dir), str(dst_dir))
            
            # Check transformation
            converted_file = dst_dir / "view_policies_converted.tf"
            assert converted_file.exists(), "view_policies.tf was not converted"
            
            content = converted_file.read_text()
            
            # Check that list attributes are preserved (these are likely sets in new provider)
            assert 'nfs_no_squash    = ["10.66.4.67"]' in content
            assert 'nfs_read_write   = ["10.66.4.67"]' in content
            assert 'nfs_no_squash       = ["10.72.34.0/24"]' in content
            assert 'nfs_read_write   = ["10.72.34.0/24"]' in content
            assert 'nfs_read_write   = ["10.72.33.0/24","10.72.34.0/24"]' in content
            
            # Check that vip_pools references are preserved
            assert 'vip_pools     = [vastdata_vip_pool.prod.id]' in content
            assert 'vip_pools     = [vastdata_vip_pool.dev.id]' in content
            assert 'vip_pools     = [vastdata_vip_pool.dev.id,vastdata_vip_pool.prod.id]' in content
    
    def test_all_customer_files_migrate_successfully(self, vast_terraform_dir, temp_migration_dirs):
        """Test that all customer files can be migrated without syntax errors."""
        src_dir, dst_dir = temp_migration_dirs
        
        # Copy all terraform files from customer directory
        for tf_file in vast_terraform_dir.glob("*.tf"):
            shutil.copy2(tf_file, src_dir)
        
        # Run migration 
        main(str(src_dir), str(dst_dir))
        
        # Verify all files were converted
        original_files = list(vast_terraform_dir.glob("*.tf"))
        converted_files = list(dst_dir.glob("*_converted.tf"))
        
        assert len(converted_files) == len(original_files), \
            f"Expected {len(original_files)} converted files, got {len(converted_files)}"
        
        # Check each converted file has content and basic syntax
        for converted_file in converted_files:
            content = converted_file.read_text()
            assert len(content.strip()) > 0, f"Converted file {converted_file.name} is empty"
            
            # Basic syntax validation - balanced braces
            open_braces = content.count('{')
            close_braces = content.count('}')
            assert open_braces == close_braces, \
                f"Mismatched braces in {converted_file.name}: {open_braces} open, {close_braces} close"
                
            # Basic syntax validation - balanced brackets
            open_brackets = content.count('[')
            close_brackets = content.count(']')
            assert open_brackets == close_brackets, \
                f"Mismatched brackets in {converted_file.name}: {open_brackets} open, {close_brackets} close"
    
    def test_provider_configuration_unchanged(self, vast_terraform_dir, temp_migration_dirs):
        """Test that provider configuration block is preserved correctly."""
        src_dir, dst_dir = temp_migration_dirs
        
        main_tf = vast_terraform_dir / "main.tf"
        if main_tf.exists():
            # Copy file to source directory
            shutil.copy2(main_tf, src_dir)
            
            # Run migration
            main(str(src_dir), str(dst_dir))
            
            # Check transformation
            converted_file = dst_dir / "main_converted.tf"
            assert converted_file.exists(), "main.tf was not converted"
            
            content = converted_file.read_text()
            
            # Provider configuration should be unchanged
            assert 'provider "vault"' in content
            assert 'provider vastdata' in content
            assert 'data "vault_generic_secret"' in content
            assert 'terraform {' in content
            assert 'required_providers' in content
            
            # Backend configuration should be preserved
            assert 'backend "pg"' in content
    
    def test_lifecycle_blocks_preserved(self, vast_terraform_dir, temp_migration_dirs):
        """Test that lifecycle blocks in resources are preserved."""
        src_dir, dst_dir = temp_migration_dirs
        
        view_policies_tf = vast_terraform_dir / "view_policies.tf"
        if view_policies_tf.exists():
            # Copy file to source directory
            shutil.copy2(view_policies_tf, src_dir)
            
            # Run migration
            main(str(src_dir), str(dst_dir))
            
            # Check transformation
            converted_file = dst_dir / "view_policies_converted.tf"
            assert converted_file.exists(), "view_policies.tf was not converted"
            
            content = converted_file.read_text()
            
            # Lifecycle blocks should be preserved
            assert 'lifecycle {' in content
            assert 'ignore_changes = [' in content
            assert 'count_views,' in content
    
    def test_comments_and_formatting_preserved(self, vast_terraform_dir, temp_migration_dirs):
        """Test that comments and basic formatting are preserved during migration."""
        src_dir, dst_dir = temp_migration_dirs
        
        view_policies_tf = vast_terraform_dir / "view_policies.tf"
        if view_policies_tf.exists():
            # Copy file to source directory
            shutil.copy2(view_policies_tf, src_dir)
            
            # Get original content
            original_content = view_policies_tf.read_text()
            
            # Run migration
            main(str(src_dir), str(dst_dir))
            
            # Check transformation
            converted_file = dst_dir / "view_policies_converted.tf"
            assert converted_file.exists(), "view_policies.tf was not converted"
            
            content = converted_file.read_text()
            
            # Comments should be preserved
            assert '## holds shared view policies' in content
            assert '// only allow root admin access' in content
            assert '// do not allow access over' in content
            assert '// nfs_no_squash is needed' in content
            assert '// read write from k8s' in content
    
    @pytest.mark.parametrize("customer_file", [
        "main.tf",
        "vip_pools.tf", 
        "view_policies.tf",
        "views.tf",
        "view_coredata.tf",
        "view_services.tf",
    ])
    def test_individual_customer_file_migration(self, vast_terraform_dir, temp_migration_dirs, customer_file):
        """Test migration of individual customer files."""
        src_dir, dst_dir = temp_migration_dirs
        
        source_file = vast_terraform_dir / customer_file
        if source_file.exists():
            # Copy file to source directory
            shutil.copy2(source_file, src_dir)
            
            # Run migration
            main(str(src_dir), str(dst_dir))
            
            # Check transformation
            converted_file = dst_dir / f"{source_file.stem}_converted.tf"
            assert converted_file.exists(), f"Customer file {customer_file} was not converted"
            
            content = converted_file.read_text()
            assert len(content.strip()) > 0, f"Converted {customer_file} is empty"
            
            # Basic structure validation
            original_content = source_file.read_text()
            
            # Count major elements to ensure nothing was lost
            original_resources = original_content.count('resource ')
            original_variables = original_content.count('variable ')
            original_outputs = original_content.count('output ')
            original_data = original_content.count('data ')
            
            converted_resources = content.count('resource ')
            converted_variables = content.count('variable ')
            converted_outputs = content.count('output ')
            converted_data = content.count('data ')
            
            assert converted_resources == original_resources, \
                f"Resource count mismatch in {customer_file}"
            assert converted_variables == original_variables, \
                f"Variable count mismatch in {customer_file}"
            assert converted_outputs == original_outputs, \
                f"Output count mismatch in {customer_file}"
            assert converted_data == original_data, \
                f"Data source count mismatch in {customer_file}"

# Copyright (c) HashiCorp, Inc.

"""
Auxiliary tests for Terraform provider edge cases and specific scenarios.

These tests cover scenarios that aren't part of the standard e2e test suite,
such as state recovery after cluster restore.
"""

import json
import pytest
from pathlib import Path


@pytest.mark.e2e
def test_id_refresh_after_state_corruption(terraform_cmd, terraform_workdir):
    """
    Test that terraform refresh correctly restores resource IDs when state is corrupted.
    
    This simulates a scenario where:
    1. A resource is created (e.g., tenant with id=3)
    2. The state file is manually corrupted (id changed to 99999)
    3. terraform refresh is run
    4. The id should be restored to the correct value (3)
    """

    # Create tenant configuration
    tenant_config = """\
resource "vastdata_tenant" "test_tenant" {
  name                            = "tf-test-refresh-tenant"
  client_ip_ranges                = [["193.168.99.1", "193.168.99.2"]]
  local_provider_id               = 1
  default_others_share_level_perm = "FULL"
}
"""
    
    test_tf_file = terraform_workdir / "test.tf"
    test_tf_file.write(tenant_config)
    
    original_id = None
    
    try:
        # Apply configuration to create tenant
        terraform_cmd['apply', '-auto-approve', '-input=false']()

        # Extract the original ID from state
        tfstate_path = terraform_workdir / "terraform.tfstate"
        
        with open(tfstate_path, 'r') as f:
            tfstate = json.load(f)
        
        # Find the tenant resource and get its ID
        for resource in tfstate.get('resources', []):
            if resource.get('type') == 'vastdata_tenant' and resource.get('name') == 'test_tenant':
                instances = resource.get('instances', [])
                if instances:
                    original_id = instances[0].get('attributes', {}).get('id')
                    break
        
        if original_id is None:
            raise AssertionError("Could not find tenant ID in state")
        
        print(f"Original ID: {original_id}")
        
        # Step 3: Corrupt the state by changing the ID to 99999
        print("Corrupting state - changing ID to 99999")
        
        for resource in tfstate.get('resources', []):
            if resource.get('type') == 'vastdata_tenant' and resource.get('name') == 'test_tenant':
                instances = resource.get('instances', [])
                if instances:
                    instances[0]['attributes']['id'] = 99999
                    break
        
        with open(tfstate_path, 'w') as f:
            json.dump(tfstate, f, indent=2)
        
        # Verify corruption
        with open(tfstate_path, 'r') as f:
            corrupted_state = json.load(f)
        
        corrupted_id = None
        for resource in corrupted_state.get('resources', []):
            if resource.get('type') == 'vastdata_tenant' and resource.get('name') == 'test_tenant':
                instances = resource.get('instances', [])
                if instances:
                    corrupted_id = instances[0].get('attributes', {}).get('id')
                    break
        
        assert corrupted_id == 99999, f"State corruption failed: ID is {corrupted_id}, expected 99999"

        # Run terraform refresh to recover correct ID
        terraform_cmd['refresh', '-input=false']()

        # Verify the ID was restored
        with open(tfstate_path, 'r') as f:
            refreshed_state = json.load(f)
        
        refreshed_id = None
        for resource in refreshed_state.get('resources', []):
            if resource.get('type') == 'vastdata_tenant' and resource.get('name') == 'test_tenant':
                instances = resource.get('instances', [])
                if instances:
                    refreshed_id = instances[0].get('attributes', {}).get('id')
                    break

        # The refreshed ID should match the original ID
        assert str(refreshed_id) == str(original_id), \
            f"ID was not restored correctly! Expected {original_id}, got {refreshed_id}"

        # Destroy the tenant
        terraform_cmd['destroy', '-auto-approve', '-input=false', '-lock=false']()

    except Exception as e:
        # Clean up on failure
        print(f"\n  ERROR: {e}")
        print("  Cleaning up after error...")
        try:
            terraform_cmd['destroy', '-auto-approve', '-input=false', '-lock=false']()
        except:
            pass
        
        # Delete state file
        state_file = terraform_workdir / "terraform.tfstate"
        if state_file.exists():
            state_file.unlink()
        
        raise
    
    finally:
        # Always clean up test file
        if test_tf_file.exists():
            test_tf_file.unlink()


@pytest.mark.e2e
def test_optional_field_preserved_after_refresh(terraform_cmd, terraform_workdir):
    """
    Test that user-declared optional fields are preserved after terraform refresh
    when the API returns null for those fields.
    
    This tests the fix for the drift issue where optional+computed fields
    were being overwritten by null values from the API.
    """

    # Create a protected path with target_exported_dir (optional+computed field)
    # Note: This test requires a view and tenant to exist
    # For simplicity, we'll use tenant's default_others_share_level_perm as the test field
    
    tenant_config = """\
resource "vastdata_tenant" "test_tenant" {
  name                            = "tf-test-optional-field"
  client_ip_ranges                = [["193.168.88.1", "193.168.88.2"]]
  local_provider_id               = 1
  default_others_share_level_perm = "FULL"
}
"""
    
    test_tf_file = terraform_workdir / "test.tf"
    test_tf_file.write(tenant_config)
    
    try:
        # Apply configuration
        terraform_cmd['apply', '-auto-approve', '-input=false']()

        # Run terraform plan to check for drift
        result = terraform_cmd['plan', '-input=false', '-detailed-exitcode', '-lock=false'].run(retcode=None)
        returncode, stdout, stderr = result
        
        if returncode == 0:
            print("No drift detected - optional field preserved correctly")
        elif returncode == 2:
            print(f"Drift detected!")
            print(f"Plan output:\n{stdout}")
            raise AssertionError("Drift detected - optional field was not preserved")
        else:
            print(f"  ERROR: terraform plan failed: {stderr}")
            raise Exception(f"terraform plan failed")
        
        # Run refresh and check again
        terraform_cmd['refresh', '-input=false']()
        
        # Run plan again to verify no drift after refresh
        result = terraform_cmd['plan', '-input=false', '-detailed-exitcode', '-lock=false'].run(retcode=None)
        returncode, stdout, stderr = result
        
        if returncode == 0:
            print("No drift after refresh - optional field still preserved")
        elif returncode == 2:
            print(f"Drift detected after refresh!")
            print(f"Plan output:\n{stdout}")
            raise AssertionError("Drift detected after refresh - optional field was overwritten")
        else:
            print(f"  ERROR: terraform plan failed: {stderr}")
            raise Exception(f"terraform plan failed")
        
        # Destroy
        terraform_cmd['destroy', '-auto-approve', '-input=false', '-lock=false']()

    except Exception as e:
        print(f"\n  ERROR: {e}")
        print("  Cleaning up after error...")
        try:
            terraform_cmd['destroy', '-auto-approve', '-input=false', '-lock=false']()
        except:
            pass
        
        state_file = terraform_workdir / "terraform.tfstate"
        if state_file.exists():
            state_file.unlink()
        
        raise
    
    finally:
        if test_tf_file.exists():
            test_tf_file.unlink()

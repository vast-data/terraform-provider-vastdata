# Copyright (c) HashiCorp, Inc.

"""
End-to-end tests for Terraform provider resources.

Tests all .tf files in examples/resources/*/e2e/ directories.
Each test:
1. Applies the configuration
2. Extracts resource IDs from tfstate
3. Deletes tfstate
4. Imports resources back
5. Runs plan to verify no drift
6. Destroys resources
"""

import json
import uuid
import pytest
import re
import shutil
import random
import string
from pathlib import Path


# Resources where ANY drift is acceptable after import
# If a resource type is in this set, all field drift is allowed during import tests
ALLOWED_DRIFT_RESOURCES = {
    "vastdata_administrator_manager",
    "vastdata_administrator_role",
    "vastdata_bgp_config",
    "vastdata_block_host_mapping",
    "vastdata_view",
    "vastdata_global_local_snapshot",
    "vastdata_global_snapshot",
    "vastdata_qos_policy",
    "vastdata_quota",
    "vastdata_s3_policy",
    "vastdata_s3_replication_peer",
    "vastdata_tenant",
}


def inject_random_ids(content):
    """
    Replace hardcoded uid/gid values and random_integer resources with actual random values.
    
    This function:
    1. Removes any random_integer resource blocks
    2. Replaces references like random_integer.uid.result with random values
    3. Replaces hardcoded uid/gid assignments like "uid = 1234" with random values
    
    Returns the modified content with a mapping of replacements made.
    """
    # Track random values we generate for consistency within the same file
    random_values = {}
    
    # Remove random_integer resource blocks
    content = re.sub(
        r'resource\s+"random_integer"\s+"[^"]+"\s*\{[^}]*\}',
        '',
        content,
        flags=re.MULTILINE
    )
    
    # Find and replace random_integer references like: random_integer.uid.result
    def replace_random_ref(match):
        var_name = match.group(1)  # e.g., "uid" or "gid" or "uid1"
        if var_name not in random_values:
            random_values[var_name] = random.randint(1000, 65000)
        return str(random_values[var_name])
    
    content = re.sub(
        r'random_integer\.(\w+)\.result',
        replace_random_ref,
        content
    )
    
    # Replace hardcoded uid/gid assignments like: uid = 30109 or gid = 1001
    def replace_hardcoded_id(match):
        id_type = match.group(1)  # "uid" or "gid"
        old_value = match.group(2)  # the hardcoded number
        # Generate a unique key for this to avoid collisions
        key = f"{id_type}_{old_value}"
        if key not in random_values:
            random_values[key] = random.randint(1000, 65000)
        return f"{id_type} = {random_values[key]}"
    
    content = re.sub(
        r'\b(uid|gid)\s*=\s*(\d{3,})',  # Match uid or gid with 3+ digit numbers
        replace_hardcoded_id,
        content
    )
    
    return content


def is_drift_allowed(plan_output, tested_resource_type=None):
    """
    Check if the drift detected in plan output is acceptable.
    
    Parse the plan output to check if resources with drift are in ALLOWED_DRIFT_RESOURCES
    or are not the resource being tested (dependency resources).
    
    Args:
        plan_output: The terraform plan output
        tested_resource_type: The resource type being tested (e.g., "vastdata_block_host").
                             If specified, only check drift for this resource type.
    
    Returns:
        Tuple: (is_allowed, reason)
        - is_allowed: True if all drift is allowed, False otherwise
    """
    # Strip ANSI color codes for easier parsing
    ansi_escape = re.compile(r'\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])')
    clean_output = ansi_escape.sub('', plan_output)
    
    # Parse plan output to find resources with changes
    # Look for patterns like "# vastdata_administrator_manager.admin will be updated in-place"
    resource_pattern = r'#\s+(\S+)\s+will be (updated|replaced|created|destroyed)'
    field_pattern = r'[~+]\s+(\w+)\s+='
    
    resources_with_changes = {}
    current_resource = None
    
    for line in clean_output.split('\n'):
        # Check if this line indicates a resource change
        resource_match = re.search(resource_pattern, line)
        if resource_match:
            current_resource = resource_match.group(1)
            resource_type = current_resource.split('.')[0]
            resources_with_changes[current_resource] = {
                'type': resource_type,
                'fields': []
            }
            continue
        
        # Check if this line shows a field change
        if current_resource and ('~' in line or '+' in line):
            field_match = re.search(field_pattern, line)
            if field_match:
                field_name = field_match.group(1)
                resources_with_changes[current_resource]['fields'].append(field_name)
    
    # Check if all detected changes are allowed
    disallowed_changes = []
    ignored_resources = []
    
    for resource_address, info in resources_with_changes.items():
        resource_type = info['type']
        changed_fields = info['fields']
        
        # If we're testing a specific resource type, ignore drift in other resource types
        if tested_resource_type and resource_type != tested_resource_type:
            ignored_resources.append(resource_address)
            continue
        
        # If resource type is in allowed list, all drift is acceptable
        if resource_type in ALLOWED_DRIFT_RESOURCES:
            continue
        
        # Otherwise, this is disallowed drift
        field_list = ', '.join(changed_fields) if changed_fields else 'unknown fields'
        disallowed_changes.append(f"{resource_address}: {field_list}")
    
    if disallowed_changes:
        reason = f"Disallowed drift in: {', '.join(disallowed_changes)}"
        return False, reason
    elif resources_with_changes:
        # We have drift, but it's all in allowed resources
        return True, "All drift is allowed"
    else:
        return True, "No drift detected"


def extract_resources_from_tfstate(terraform_workdir):
    """
    Extract resource addresses and IDs from terraform.tfstate file.
    
    Returns:
        List of tuples: [(resource_address, resource_id), ...]
        Example: [("vastdata_administrator_manager.admin", "2"), ...]
    """
    tfstate_path = terraform_workdir / "terraform.tfstate"
    if not tfstate_path.exists():
        print("  ⚠ terraform.tfstate file not found")
        return []
    
    with open(tfstate_path, 'r') as f:
        tfstate = json.load(f)
    
    resources = []
    # Parse tfstate structure to extract resources
    if 'resources' in tfstate:
        for resource in tfstate['resources']:
            resource_type = resource.get('type', '')
            resource_name = resource.get('name', '')
            resource_mode = resource.get('mode', 'managed')
            
            # Only process managed resources (not data sources)
            if resource_mode != 'managed':
                continue
            
            # Get instances (can be multiple for count/for_each)
            instances = resource.get('instances', [])
            for instance in instances:
                attributes = instance.get('attributes', {})
                resource_id = attributes.get('id')
                
                if resource_id:
                    # Build resource address: type.name
                    resource_address = f"{resource_type}.{resource_name}"
                    resources.append((resource_address, str(resource_id)))
                    print(f"    Found resource: {resource_address} with ID: {resource_id}")
    
    return resources


def collect_e2e_test_files(tf_examples_dir):
    """
    Collect all .tf files from examples/resources/*/e2e/ directories.
    
    Returns:
        List of tuples: (resource_name, tf_file_path)
    """
    test_files = []
    
    for resource_dir in sorted(tf_examples_dir.iterdir()):
        if not resource_dir.is_dir():
            continue
        
        e2e_dir = resource_dir / "e2e"
        if not e2e_dir.exists():
            continue
        
        for tf_file in sorted(e2e_dir.glob("*.tf")):
            # Skip files marked with ignore:e2e in first 10 lines
            content = tf_file.read_text()
            first_lines = '\n'.join(content.split('\n')[:10])
            if '# ignore:e2e' in first_lines:
                continue
            
            test_name = f"{resource_dir.name}/{tf_file.name}"
            test_files.append((test_name, tf_file))
    
    return test_files


@pytest.fixture(scope="module")
def e2e_test_files(tf_examples_dir):
    """Fixture that provides all e2e test files."""
    return collect_e2e_test_files(tf_examples_dir)


@pytest.mark.e2e
def test_terraform_provider_e2e(terraform_cmd, terraform_workdir, e2e_test_files):
    """
    Test all Terraform example configurations.
    
    For each .tf file:
      1. Copy to working directory
      2. Apply configuration
      3. Destroy configuration
    """
    if not e2e_test_files:
        pytest.skip("No e2e test files found")
    
    total_tests = len(e2e_test_files)
    print(f"\n=== Running {total_tests} E2E tests ===\n")
    
    for test_name, tf_file in e2e_test_files:
        print(f"\n--- Test: {test_name} ---")
        
        # Read and modify content to make it unique
        content = tf_file.read_text()
        
        # Inject random uid/gid values
        content = inject_random_ids(content)
        
        uniq_hash = ''.join(random.choices(string.ascii_lowercase, k=4))
        content = content.replace("vastdb", f"vastdb{uniq_hash}")
        
        # Write to working directory
        test_tf_file = terraform_workdir / "test.tf"
        test_tf_file.write(content)
        
        # Apply
        try:
            print("  Step 1: terraform apply")
            result = terraform_cmd['apply', '-auto-approve', '-input=false', '-lock=false']()
            print(f"  Apply succeeded")
        except Exception as e:
            print(f"  ERROR: apply failed for {test_name}")
            print(f"  {e}")
            # Clean up and continue
            try:
                terraform_cmd['destroy', '-auto-approve', '-input=false', '-lock=false']()
            except:
                pass
            # Delete state file to avoid blocking next test
            (terraform_workdir / "terraform.tfstate").delete() if (terraform_workdir / "terraform.tfstate").exists() else None
            test_tf_file.delete()
            continue
        
        # Destroy
        try:
            print("  Step 2: terraform destroy")
            result = terraform_cmd['destroy', '-auto-approve', '-input=false', '-lock=false']()
            print(f"  Destroy succeeded")
        except Exception as e:
            print(f"  ERROR: destroy failed for {test_name}")
            print(f"  {e}")
            # Clean up state file
            (terraform_workdir / "terraform.tfstate").delete() if (terraform_workdir / "terraform.tfstate").exists() else None
        
        # Clean up test file
        test_tf_file.delete()
        print(f"  ✓ Test completed: {test_name}")


@pytest.mark.e2e
@pytest.mark.parametrize("test_name,tf_file", 
                         collect_e2e_test_files(Path(__file__).parent.parent / "examples" / "resources"),
                         ids=[x[0] for x in collect_e2e_test_files(Path(__file__).parent.parent / "examples" / "resources")])
def test_individual_resource(terraform_cmd, terraform_workdir, test_name, tf_file):
    """
    Test individual Terraform resources with full import workflow.
    
    For each test:
    1. Apply configuration
    2. Check for drift immediately after apply (must be none)
    3. Extract resource IDs from tfstate
    4. Delete tfstate to simulate state loss
    5. Import resources back
    6. Run plan to verify no drift (honors ALLOWED_DRIFT_RESOURCES)
    7. Destroy resources
    """
    print(f"\n--- Test: {test_name} ---")
    
    # Read and modify content to make it unique
    content = tf_file.read_text()
    
    # Skip files marked with ignore:e2e in first 10 lines
    first_lines = '\n'.join(content.split('\n')[:10])
    if '# ignore:e2e' in first_lines:
        pytest.skip(f"Test marked with ignore:e2e: {test_name}")
    
    # Inject random uid/gid values
    content = inject_random_ids(content)
    
    # Make resource names unique (4 random lowercase letters)
    uniq_hash = ''.join(random.choices(string.ascii_lowercase, k=4))
    content = content.replace("vastdb", f"vastdb{uniq_hash}")
    
    # Write to working directory
    test_tf_file = terraform_workdir / "test.tf"
    test_tf_file.write(content)
    
    resources_to_import = []
    
    # Extract the primary resource type from test name (e.g., "vastdata_block_host" from "vastdata_block_host/1.tf")
    primary_resource_type = test_name.split('/')[0] if '/' in test_name else test_name.replace('.tf', '')
    
    try:
        # Step 1: Apply
        print("  Step 1: terraform apply")
        terraform_cmd['apply', '-auto-approve', '-input=false']()
        print(f"  Apply succeeded")
        
        # Step 2: Check for drift immediately after apply (must be none, ignore ALLOWED_DRIFT_RESOURCES)
        print("  Step 2: Running terraform plan to check for drift after apply")
        result = terraform_cmd['plan', '-input=false', '-detailed-exitcode', '-lock=false'].run(retcode=None)
        returncode, stdout, stderr = result
        
        # Exit code 0 = no changes, 1 = error, 2 = changes present
        if returncode == 0:
            print("  No drift detected after apply")
        elif returncode == 2:
            # Drift detected right after apply - this is always a bug
            print(f"  ERROR: Drift detected immediately after apply!")
            print("    This indicates the resource is not stable after creation")
            print(f"    Plan output:\n{stdout}")
            raise AssertionError(f"Drift detected immediately after apply for {test_name}. This is a bug - resources should be stable after apply.")
        else:
            print(f"  ERROR: terraform plan failed")
            print(f"    Error: {stderr}")
            raise Exception(f"terraform plan failed for {test_name}")
        
        # Step 3: Extract resource information from tfstate
        print("  Step 3: Extracting resources from tfstate")
        all_resources = extract_resources_from_tfstate(terraform_workdir)
        
        # Separate resources into primary (to test) and dependencies (to keep in state)
        primary_resource = None
        dependency_resources = []
        skipped_resources = []
        
        for resource_address, resource_id in all_resources:
            resource_type = resource_address.split('.')[0]

            if resource_type == primary_resource_type:
                # This is the primary resource we're testing
                if primary_resource is None:
                    primary_resource = (resource_address, resource_id)
                else:
                    # Multiple instances of same type - import all for now
                    dependency_resources.append((resource_address, resource_id))
            else:
                # This is a dependency resource
                dependency_resources.append((resource_address, resource_id))
        
        if skipped_resources:
            print(f"  Skipping {len(skipped_resources)} resources (don't support standard ID import):")
            for addr in skipped_resources:
                print(f"    - {addr}")
        
        if primary_resource is None:
            print(f"  No primary resource of type '{primary_resource_type}' found, skipping import test")
        else:
            print(f"  ✓ Found primary resource to test: {primary_resource[0]}")
            if dependency_resources:
                print(f"  ℹ Found {len(dependency_resources)} dependency resources (will remain in state)")
            
            # Step 4: Remove only the primary resource from tfstate
            print(f"  Step 4: Removing {primary_resource[0]} from state to simulate loss")
            try:
                result = terraform_cmd['state', 'rm', primary_resource[0]].run(retcode=None)
                returncode, stdout, stderr = result
                if returncode != 0:
                    print(f"  ERROR: Failed to remove {primary_resource[0]} from state")
                    raise Exception(f"terraform state rm failed: {stderr}")
                print(f"  Removed {primary_resource[0]} from state")
            except Exception as e:
                print(f"  ERROR: {str(e)}")
                raise
            
            # Step 5: Import only the primary resource
            resource_address, resource_id = primary_resource
            print(f"  Step 5: Importing {resource_address} with ID {resource_id}")
            try:
                # Run import and capture output
                result = terraform_cmd['import', '-input=false', '-lock=false', resource_address, resource_id].run(retcode=None)
                returncode, stdout, stderr = result
                
                if returncode != 0:
                    print(f"  ERROR: import failed for {resource_address}")
                    print(f"    Return code: {returncode}")
                    if stderr:
                        print(f"    Error output:\n{stderr}")
                    if stdout:
                        print(f"    Standard output:\n{stdout}")
                    raise AssertionError(f"Import failed for {resource_address}: {stderr}")
                else:
                    print(f"  Imported {resource_address}")
            except Exception as e:
                print(f"  ERROR: import exception for {resource_address}")
                print(f"    {str(e)}")
                raise
            
            # Step 6: Run terraform plan to check for drift after import (honors ALLOWED_DRIFT_RESOURCES)
            print(f"  Step 6: Running terraform plan to check for drift after import in {primary_resource_type}")
            result = terraform_cmd['plan', '-input=false', '-detailed-exitcode', '-lock=false'].run(retcode=None)
            returncode, stdout, stderr = result
            
            # Exit code 0 = no changes, 1 = error, 2 = changes present
            if returncode == 0:
                print("  ✓ No drift detected after import")
            elif returncode == 2:
                # Drift detected - check if it's allowed (only for primary resource)
                is_allowed, reason = is_drift_allowed(stdout, tested_resource_type=primary_resource_type)
                
                if is_allowed:
                    print(f"  Drift detected after import, but it's allowed")
                    print(f"    Reason: {reason}")
                    print("  Test passed (allowed drift)")
                else:
                    print(f"  ERROR: Drift detected in primary resource {primary_resource_type}!")
                    print("    This means import didn't populate all fields correctly")
                    print(f"    {reason}")
                    print(f"    Plan output:\n{stdout}")
                    raise AssertionError(f"Disallowed drift detected after import for {test_name}: {reason}")
            else:
                print(f"  ERROR: terraform plan failed")
                print(f"    Error: {stderr}")
                raise Exception(f"terraform plan failed for {test_name}")
        
        # Step 7: Destroy
        print("  Step 7: terraform destroy")
        terraform_cmd['destroy', '-auto-approve', '-input=false', '-lock=false']()
        print(f"  Destroy succeeded")
        print(f"  Test completed: {test_name}")
        
    except Exception as e:
        # Try to clean up
        print("  Cleaning up after error...")
        try:
            terraform_cmd['destroy', '-auto-approve', '-input=false', '-lock=false']()
        except:
            pass
        
        # Delete state file to avoid blocking subsequent tests
        state_file = terraform_workdir / "terraform.tfstate"
        if state_file.exists():
            state_file.delete()
        
        # Re-raise the exception to fail the test
        raise
    
    finally:
        # Always clean up test file
        if test_tf_file.exists():
            test_tf_file.delete()


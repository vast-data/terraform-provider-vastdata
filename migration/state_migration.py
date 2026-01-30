#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Terraform State Migration Tool for VastData Provider
Migrates existing tfstate from v1.6.7 to v2.x by re-importing resources.
"""

import os
import sys
import json
import argparse
import subprocess
import tempfile
import shutil
import threading
from pathlib import Path
from typing import Dict, List, Tuple, Optional
from concurrent.futures import ThreadPoolExecutor, as_completed

VERSION = "1.1.0"

# Resource import field mappings
# Format: resource_type -> (import_fields, id_field_in_state)
RESOURCE_IMPORT_MAP = {
    # Resources that import by simple ID
    "vastdata_tenant": (["id"], "id"),
    "vastdata_vip_pool": (["id"], "id"),
    "vastdata_view_policy": (["id"], "id"),
    "vastdata_view": (["id"], "id"),
    "vastdata_user": (["id"], "id"),
    "vastdata_group": (["id"], "id"),
    "vastdata_nonlocal_user": (["username", "context", "tenant_id"], None),
    "vastdata_nonlocal_group": (["groupname", "context", "tenant_id"], None),
    "vastdata_quota": (["id"], "id"),
    "vastdata_protection_policy": (["id"], "id"),
    "vastdata_protected_path": (["id"], "id"),
    "vastdata_s3_policy": (["id"], "id"),
    "vastdata_s3_policy_attachment": (["id"], "id"),
    "vastdata_replication_peer": (["id"], "id"),
    "vastdata_s3_replication_peer": (["id"], "id"),
    "vastdata_active_directory": (["id"], "id"),
    "vastdata_ldap": (["id"], "id"),
    "vastdata_dns": (["id"], "id"),
    "vastdata_nis": (["id"], "id"),
    "vastdata_local_provider": (["id"], "id"),
    "vastdata_snapshot": (["id"], "id"),
    "vastdata_global_snapshot": (["id"], "id"),
    "vastdata_global_local_snapshot": (["id"], "id"),
    "vastdata_qos_policy": (["id"], "id"),
    "vastdata_administrator_manager": (["id"], "id"),
    "vastdata_administrator_role": (["id"], "id"),
    "vastdata_administrator_realm": (["id"], "id"),
    "vastdata_bgp_config": (["id"], "id"),
    "vastdata_block_host": (["id"], "id"),
    "vastdata_block_host_mapping": (["id"], "id"),
    "vastdata_encryption_group": (["id"], "id"),
    "vastdata_encryption_group_control": (["id"], "id"),
    "vastdata_event_definition": (["id"], "id"),
    "vastdata_event_definition_config": (["id"], "id"),
    "vastdata_folder_readonly": (["id"], "id"),
    "vastdata_kafka_broker": (["id"], "id"),
    "vastdata_saml_config": (["id"], "id"),
    "vastdata_vms": (["id"], "id"),
    "vastdata_volume": (["id"], "id"),
    "vastdata_user_copy": (["id"], "id"),
    "vastdata_user_key": (["id"], "id"),
    "vastdata_nonlocal_user_key": (["id"], "id"),
    "vastdata_apitoken": (["id"], "id"),
    "vastdata_s3_lifecycle_rule": (["id"], "id"),
    "vastdata_s3_life_cycle_rule": (["id"], "id"),
    # v1.x legacy resource names (plural forms and old naming)
    "vastdata_administators_managers": (["id"], "id"),
    "vastdata_administators_roles": (["id"], "id"),
    "vastdata_administators_realms": (["id"], "id"),
    "vastdata_kafka_brokers": (["id"], "id"),
    "vastdata_replication_peers": (["id"], "id"),
    "vastdata_s3_replication_peers": (["id"], "id"),
    "vastdata_active_directory2": (["id"], "id"),
    "vastdata_non_local_user": (["username", "context", "tenant_id"], None),
    "vastdata_non_local_user_key": (["id"], "id"),
    "vastdata_non_local_group": (["groupname", "context", "tenant_id"], None),
    "vastdata_saml": (["id"], "id"),
    "vastdata_blockhost": (["id"], "id"),
    # Note: v3.0 new resources are not included as they don't need migration from v1.6.7
}

# Map v1/v2 legacy resource names to v3 resource names
# Used when generating stub resource definitions for terraform import
RESOURCE_NAME_TRANSLATION = {
    # v1.x plural forms → v3.x singular forms
    "vastdata_administators_managers": "vastdata_administrator_manager",
    "vastdata_administators_roles": "vastdata_administrator_role",
    "vastdata_administators_realms": "vastdata_administrator_realm",
    "vastdata_kafka_brokers": "vastdata_kafka_broker",
    "vastdata_replication_peers": "vastdata_replication_peer",
    "vastdata_s3_replication_peers": "vastdata_s3_replication_peer",
    # v1.x/v2.x old naming → v3.x new naming
    "vastdata_active_directory2": "vastdata_active_directory",
    "vastdata_non_local_user": "vastdata_nonlocal_user",
    "vastdata_non_local_user_key": "vastdata_nonlocal_user_key",
    "vastdata_non_local_group": "vastdata_nonlocal_group",
    "vastdata_saml": "vastdata_saml_config",
    "vastdata_blockhost": "vastdata_block_host",
    # v2.x naming variant → v3.x naming
    "vastdata_s3_lifecycle_rule": "vastdata_s3_life_cycle_rule",
}


class Colors:
    """ANSI color codes for terminal output"""
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    NC = '\033[0m'  # No Color


def log_info(msg):
    print(f"{Colors.BLUE}[INFO]{Colors.NC} {msg}")


def log_success(msg):
    print(f"{Colors.GREEN}[SUCCESS]{Colors.NC} {msg}")


def log_warning(msg):
    print(f"{Colors.YELLOW}[WARNING]{Colors.NC} {msg}")


def log_error(msg):
    print(f"{Colors.RED}[ERROR]{Colors.NC} {msg}")


def download_state_from_s3(bucket: str, key: str, access_key: str, secret_key: str, 
                           endpoint: Optional[str] = None) -> str:
    """Download tfstate file from S3"""
    import boto3
    from botocore.exceptions import ClientError
    
    log_info(f"Downloading state from S3: s3://{bucket}/{key}")
    
    # Configure S3 client
    s3_config = {
        'aws_access_key_id': access_key,
        'aws_secret_access_key': secret_key,
    }
    
    if endpoint:
        s3_config['endpoint_url'] = endpoint
    
    try:
        s3_client = boto3.client('s3', **s3_config)
        
        # Create temp file
        temp_file = tempfile.NamedTemporaryFile(mode='w', suffix='.tfstate', delete=False)
        temp_file.close()
        
        # Download file
        s3_client.download_file(bucket, key, temp_file.name)
        log_success(f"Downloaded state to {temp_file.name}")
        
        return temp_file.name
    except ClientError as e:
        log_error(f"Failed to download from S3: {e}")
        sys.exit(1)
    except Exception as e:
        log_error(f"Unexpected error downloading from S3: {e}")
        sys.exit(1)


def parse_tfstate(state_file: str) -> Dict:
    """Parse Terraform state file"""
    log_info(f"Parsing state file: {state_file}")
    
    try:
        with open(state_file, 'r') as f:
            state = json.load(f)
        
        log_success(f"Successfully parsed state file (version: {state.get('version', 'unknown')})")
        return state
    except json.JSONDecodeError as e:
        log_error(f"Invalid JSON in state file: {e}")
        sys.exit(1)
    except FileNotFoundError:
        log_error(f"State file not found: {state_file}")
        sys.exit(1)


def extract_resources(state: Dict) -> List[Dict]:
    """Extract all resources from tfstate"""
    resources = []
    
    # Handle different tfstate versions
    if 'resources' in state:
        # TF >= 0.12 format
        for resource in state['resources']:
            mode = resource.get('mode', 'managed')
            if mode != 'managed':
                continue  # Skip data sources
            
            resource_type = resource['type']
            name = resource['name']
            module = resource.get('module', '')
            
            # Handle instances
            instances = resource.get('instances', [])
            for idx, instance in enumerate(instances):
                attributes = instance.get('attributes', {})
                
                # Build resource address
                if module:
                    address = f"{module}.{resource_type}.{name}"
                else:
                    address = f"{resource_type}.{name}"
                
                # Handle indexed resources
                if len(instances) > 1:
                    index_key = instance.get('index_key')
                    if index_key is not None:
                        if isinstance(index_key, str):
                            address = f"{address}[\"{index_key}\"]"
                        else:
                            address = f"{address}[{index_key}]"
                
                resources.append({
                    'address': address,
                    'type': resource_type,
                    'name': name,
                    'module': module,
                    'attributes': attributes
                })
    else:
        # Legacy format (TF < 0.12)
        log_warning("Legacy state format detected (TF < 0.12), attempting to parse...")
        modules = state.get('modules', [])
        for module in modules:
            module_path = '.'.join(module.get('path', []))
            for res_key, res_val in module.get('resources', {}).items():
                if not res_key.startswith('data.'):
                    resource_type = res_val.get('type', '')
                    attributes = res_val.get('primary', {}).get('attributes', {})
                    
                    resources.append({
                        'address': f"{module_path}.{res_key}" if module_path and module_path != 'root' else res_key,
                        'type': resource_type,
                        'name': res_key.split('.', 1)[1] if '.' in res_key else res_key,
                        'module': module_path if module_path != 'root' else '',
                        'attributes': attributes
                    })
    
    log_success(f"Extracted {len(resources)} resources from state")
    return resources


def separate_vast_resources(resources: List[Dict]) -> Tuple[List[Dict], List[Dict]]:
    """Separate VAST resources from other provider resources
    
    Returns:
        Tuple of (vast_resources, non_vast_resources)
    """
    vast_resources = []
    non_vast_resources = []
    
    for resource in resources:
        resource_type = resource.get('type', '')
        if resource_type.startswith('vastdata_'):
            vast_resources.append(resource)
        else:
            non_vast_resources.append(resource)
    
    log_info(f"Separated resources: {len(vast_resources)} VAST, {len(non_vast_resources)} non-VAST")
    return vast_resources, non_vast_resources


def extract_state_resources_raw(state: Dict) -> List[Dict]:
    """Extract raw resource entries from state for preservation
    
    This preserves the original state structure for non-VAST resources
    """
    raw_resources = []
    
    if 'resources' in state:
        # TF >= 0.12 format
        for resource in state['resources']:
            raw_resources.append(resource)
    else:
        # Legacy format
        log_warning("Legacy state format - non-VAST resources may not be preserved correctly")
    
    return raw_resources


def merge_non_vast_resources_to_state(state_file: str, non_vast_resources: List[Dict], original_state: Dict):
    """Merge non-VAST resources into the migrated state file
    
    Args:
        state_file: Path to the migrated state file
        non_vast_resources: List of raw resource objects from original state
        original_state: Original state dict for metadata
    """
    if not non_vast_resources:
        log_info("No non-VAST resources to merge")
        return
    
    log_info(f"Merging {len(non_vast_resources)} non-VAST resources into migrated state")
    
    try:
        # Read the migrated state
        with open(state_file, 'r') as f:
            migrated_state = json.load(f)
        
        # Add non-VAST resources to the migrated state
        if 'resources' not in migrated_state:
            migrated_state['resources'] = []
        
        migrated_state['resources'].extend(non_vast_resources)
        
        # Increment serial number
        migrated_state['serial'] = migrated_state.get('serial', 0) + 1
        
        # Write back
        with open(state_file, 'w') as f:
            json.dump(migrated_state, f, indent=2)
        
        log_success(f"Successfully merged {len(non_vast_resources)} non-VAST resources")
    except Exception as e:
        log_error(f"Failed to merge non-VAST resources: {e}")
        raise


def build_import_id(resource: Dict) -> Optional[str]:
    """Build import ID string for a resource"""
    resource_type = resource['type']
    attributes = resource['attributes']
    
    # Get import configuration for this resource type
    import_config = RESOURCE_IMPORT_MAP.get(resource_type)
    if not import_config:
        log_warning(f"No import configuration for resource type: {resource_type}")
        return None
    
    import_fields, id_field = import_config
    
    # Simple ID import
    if len(import_fields) == 1 and import_fields[0] == 'id' and id_field:
        resource_id = attributes.get(id_field)
        if resource_id:
            return str(resource_id)
        else:
            log_warning(f"Missing {id_field} for {resource['address']}")
            return None
    
    # Composite key import (use key=value format)
    import_values = []
    for field_spec in import_fields:
        # Handle field specifications: can be a string or tuple of (field, fallback_field, ...)
        if isinstance(field_spec, tuple):
            # Try each field in order until we find a non-null value
            value = None
            used_field = None
            for field in field_spec:
                value = attributes.get(field)
                if value is not None:
                    used_field = field
                    break
            
            if value is None:
                log_warning(f"Missing all fallback fields {field_spec} for {resource['address']}")
                return None
            
            import_values.append(f"{used_field}={value}")
        else:
            # Simple field (string)
            field = field_spec
            value = attributes.get(field)
            if value is None:
                log_warning(f"Missing field '{field}' for {resource['address']}")
                return None
            import_values.append(f"{field}={value}")
    
    return ','.join(import_values)


def run_import_script(script_path: str, output_dir: str) -> Tuple[int, int]:
    """Execute the import script and return success/failure counts"""
    log_info("Executing import script...")
    
    try:
        # Make script executable
        os.chmod(script_path, 0o755)
        
        # Run the script
        result = subprocess.run(
            [script_path],
            cwd=output_dir,
            capture_output=False,
            text=True
        )
        
        return result.returncode
    except Exception as e:
        log_error(f"Failed to execute import script: {e}")
        return 1


def import_single_resource(resource: Dict, terraform_dir: str, lock: threading.Lock) -> Tuple[bool, str, str]:
    """Import a single resource using terraform import
    
    Args:
        resource: Resource dict with type, address, attributes
        terraform_dir: Directory containing terraform configuration
        lock: Threading lock for terraform state operations
        
    Returns:
        Tuple of (success, resource_address, error_message)
    """
    import_id = build_import_id(resource)
    if not import_id:
        return False, resource['address'], "No import ID could be built"
    
    original_address = resource['address']
    
    # Translate resource type to v3 name for import command
    parts = original_address.split('.', 1)
    if len(parts) == 2:
        resource_type_from_state = parts[0]
        resource_name = parts[1]
        resource_type_v3 = RESOURCE_NAME_TRANSLATION.get(resource_type_from_state, resource_type_from_state)
        translated_address = f"{resource_type_v3}.{resource_name}"
    else:
        translated_address = original_address
    
    # Use lock to ensure terraform operations are serialized (terraform doesn't support concurrent state modifications)
    with lock:
        try:
            result = subprocess.run(
                ['terraform', 'import', translated_address, import_id],
                cwd=terraform_dir,
                capture_output=True,
                text=True,
                timeout=300  # 5 minute timeout per resource
            )
            
            if result.returncode == 0:
                return True, translated_address, ""
            else:
                error_msg = result.stderr if result.stderr else result.stdout
                return False, translated_address, error_msg
        except subprocess.TimeoutExpired:
            return False, translated_address, "Import timeout (>5 minutes)"
        except Exception as e:
            return False, translated_address, str(e)


def run_parallel_imports(resources: List[Dict], terraform_dir: str, max_workers: int = 5) -> Tuple[int, int, List[str]]:
    """Import resources in parallel using ThreadPoolExecutor
    
    Args:
        resources: List of resource dicts to import
        terraform_dir: Directory containing terraform configuration
        max_workers: Number of parallel workers (default: 5)
        
    Returns:
        Tuple of (success_count, failed_count, failed_resources)
    """
    log_info(f"Starting parallel import with {max_workers} workers for {len(resources)} resources")
    
    # Terraform state operations need to be serialized
    state_lock = threading.Lock()
    
    success_count = 0
    failed_count = 0
    failed_resources = []
    
    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        # Submit all import tasks
        future_to_resource = {
            executor.submit(import_single_resource, resource, terraform_dir, state_lock): resource
            for resource in resources
        }
        
        # Process completed imports
        for future in as_completed(future_to_resource):
            resource = future_to_resource[future]
            try:
                success, address, error = future.result()
                
                if success:
                    success_count += 1
                    log_success(f"[{success_count + failed_count}/{len(resources)}] ✓ Imported {address}")
                else:
                    failed_count += 1
                    failed_resources.append(address)
                    log_error(f"[{success_count + failed_count}/{len(resources)}] ✗ Failed {address}: {error[:100]}")
            except Exception as e:
                failed_count += 1
                failed_resources.append(resource['address'])
                log_error(f"✗ Exception importing {resource['address']}: {e}")
    
    log_info("=" * 80)
    log_info("Import Summary:")
    log_info(f"  Total resources: {len(resources)}")
    log_info(f"  Successfully imported: {success_count}")
    log_info(f"  Failed: {failed_count}")
    log_info("=" * 80)
    
    return success_count, failed_count, failed_resources


def generate_import_script(resources: List[Dict], output_dir: str, terraform_dir: str) -> str:
    """Generate shell script with terraform import commands"""
    script_path = os.path.join(output_dir, "import_resources.sh")
    
    log_info(f"Generating import script: {script_path}")
    
    # Generate stub resource definitions
    resources_tf_path = os.path.join(terraform_dir, "resources.tf")
    with open(resources_tf_path, 'w') as f:
        f.write("# Auto-generated stub resource definitions for import\n")
        f.write("# These are minimal definitions required by terraform import\n\n")
        for resource in resources:
            resource_type_from_state = resource['type']
            # Translate v1/v2 resource names to v3 names
            resource_type = RESOURCE_NAME_TRANSLATION.get(resource_type_from_state, resource_type_from_state)
            resource_name = resource['address'].split('.', 1)[1]  # Extract name from address
            f.write(f'resource "{resource_type}" "{resource_name}" {{\n')
            f.write('  # Configuration will be populated from import\n')
            f.write('}\n\n')
    
    with open(script_path, 'w') as f:
        f.write("#!/bin/bash\n")
        f.write("# Auto-generated Terraform import script\n")
        f.write("# This script will create a new state file with all imported resources\n\n")
        f.write("set -e  # Exit on error\n\n")
        f.write(f"cd \"{terraform_dir}\"\n\n")
        f.write("# Initialize Terraform (may be skipped if dev_overrides are active)\n")
        f.write("echo \"Initializing Terraform...\"\n")
        f.write("terraform init || true  # Don't fail if dev_overrides are in effect\n\n")
        f.write("# Remove existing state if present\n")
        f.write("if [ -f terraform.tfstate ]; then\n")
        f.write("    echo \"Removing existing state...\"\n")
        f.write("    rm terraform.tfstate\n")
        f.write("fi\n\n")
        f.write("# Import resources\n")
        f.write("echo \"Starting resource import...\"\n")
        f.write("IMPORT_COUNT=0\n")
        f.write("FAILED_COUNT=0\n\n")
        
        for idx, resource in enumerate(resources, 1):
            import_id = build_import_id(resource)
            if not import_id:
                f.write(f"# SKIPPED: {resource['address']} (no import ID)\n")
                continue
            
            # Original address from state (may use v1/v2 resource name)
            original_address = resource['address']
            
            # Translate resource type to v3 name for import command
            # Address format: resource_type.resource_name
            parts = original_address.split('.', 1)
            if len(parts) == 2:
                resource_type_from_state = parts[0]
                resource_name = parts[1]
                # Translate to v3 resource name
                resource_type_v3 = RESOURCE_NAME_TRANSLATION.get(resource_type_from_state, resource_type_from_state)
                translated_address = f"{resource_type_v3}.{resource_name}"
            else:
                translated_address = original_address
            
            f.write(f"# Import {idx}/{len(resources)}: {original_address} (as {translated_address})\n")
            f.write(f"echo \"[{idx}/{len(resources)}] Importing {translated_address}...\"\n")
            f.write(f"if terraform import '{translated_address}' '{import_id}'; then\n")
            f.write(f"    echo \"  ✓ Successfully imported {translated_address}\"\n")
            f.write(f"    IMPORT_COUNT=$((IMPORT_COUNT + 1))\n")
            f.write(f"else\n")
            f.write(f"    echo \"  ✗ Failed to import {translated_address}\"\n")
            f.write(f"    FAILED_COUNT=$((FAILED_COUNT + 1))\n")
            f.write(f"fi\n\n")
        
        f.write("# Summary\n")
        f.write("echo \"\"\n")
        f.write("echo \"===========================================\"\n")
        f.write("echo \"Import Summary:\"\n")
        f.write("echo \"  Total resources: " + str(len(resources)) + "\"\n")
        f.write("echo \"  Successfully imported: $IMPORT_COUNT\"\n")
        f.write("echo \"  Failed: $FAILED_COUNT\"\n")
        f.write("echo \"===========================================\"\n")
        f.write("\n")
        f.write("if [ $FAILED_COUNT -gt 0 ]; then\n")
        f.write("    echo \"WARNING: Some imports failed. Please review the output above.\"\n")
        f.write("    exit 1\n")
        f.write("fi\n\n")
        f.write("echo \"\"\n")
        f.write("echo \"Migration complete! New state file created: terraform.tfstate\"\n")
        f.write("echo \"You can now run 'terraform plan' to verify the migration.\"\n")
    
    # Make script executable
    os.chmod(script_path, 0o755)
    log_success(f"Generated import script with {len(resources)} resources")
    
    return script_path


def create_import_summary(resources: List[Dict], output_dir: str):
    """Create a summary file of resources to be imported"""
    summary_path = os.path.join(output_dir, "import_summary.txt")
    
    log_info(f"Creating import summary: {summary_path}")
    
    with open(summary_path, 'w') as f:
        f.write("=" * 80 + "\n")
        f.write("Terraform State Migration Summary\n")
        f.write("=" * 80 + "\n\n")
        f.write(f"Total resources to import: {len(resources)}\n\n")
        
        # Group by resource type
        by_type = {}
        for resource in resources:
            res_type = resource['type']
            if res_type not in by_type:
                by_type[res_type] = []
            by_type[res_type].append(resource)
        
        f.write("Resources by type:\n")
        f.write("-" * 80 + "\n")
        for res_type in sorted(by_type.keys()):
            f.write(f"\n{res_type}: {len(by_type[res_type])} resource(s)\n")
            for resource in by_type[res_type]:
                import_id = build_import_id(resource)
                f.write(f"  - {resource['address']}")
                if import_id:
                    f.write(f" (ID: {import_id})")
                else:
                    f.write(" (SKIP: No import ID)")
                f.write("\n")
        
        f.write("\n" + "=" * 80 + "\n")
        f.write("Next Steps:\n")
        f.write("=" * 80 + "\n")
        f.write("1. Review this summary and the generated import script\n")
        f.write("2. Ensure your .tf files are updated to v2.x format\n")
        f.write("3. Run: ./import_resources.sh\n")
        f.write("4. After successful import, run: terraform plan\n")
        f.write("5. The migrated state will be created as terraform.tfstate\n")
        f.write("6. If using S3 backend, upload terraform.tfstate to S3\n")
        f.write("\n")
    
    log_success(f"Created import summary")


def main():
    parser = argparse.ArgumentParser(
        description="Migrate Terraform state from VastData provider v1.6.7 to v2.x",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Migrate local state file
  ./state_migration.py /path/to/source/dir /path/to/output/dir

  # With S3 state
  ./state_migration.py /path/to/source/dir /path/to/output/dir \\
    --s3-bucket my-bucket \\
    --s3-key terraform.tfstate \\
    --s3-access-key AKIAIOSFODNN7EXAMPLE \\
    --s3-secret-key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

  # With custom S3 endpoint (for VAST S3)
  ./state_migration.py /path/to/source/dir /path/to/output/dir \\
    --s3-bucket my-bucket \\
    --s3-key terraform.tfstate \\
    --s3-access-key KEY \\
    --s3-secret-key SECRET \\
    --s3-endpoint https://s3.vast.local
        """
    )
    
    parser.add_argument('--version', action='version', version=f'%(prog)s {VERSION}')
    
    # Positional arguments (like run_migration.sh)
    parser.add_argument('source_dir', 
                       help='Source directory containing Terraform configuration and state')
    parser.add_argument('output_dir', 
                       help='Output directory for migrated configuration and state')
    
    # S3 options (optional)
    parser.add_argument('--s3-bucket', help='S3 bucket name containing tfstate (optional, uses local state if not provided)')
    parser.add_argument('--s3-key', help='S3 object key (path) to tfstate file')
    parser.add_argument('--s3-access-key', help='AWS/S3 access key ID')
    parser.add_argument('--s3-secret-key', help='AWS/S3 secret access key')
    parser.add_argument('--s3-endpoint', help='Custom S3 endpoint URL (for VAST S3 or S3-compatible storage)')
    parser.add_argument('--max-workers', type=int, default=5, 
                       help='Maximum number of parallel workers for resource imports (default: 5)')
    
    args = parser.parse_args()
    
    # Validate S3 arguments
    if args.s3_bucket:
        if not all([args.s3_key, args.s3_access_key, args.s3_secret_key]):
            log_error("When using --s3-bucket, you must also provide --s3-key, --s3-access-key, and --s3-secret-key")
            sys.exit(1)
        
        # Check boto3 availability
        try:
            import boto3
        except ImportError:
            log_error("boto3 library is required for S3 support. Install it with: pip install boto3")
            sys.exit(1)
    
    # Validate directories
    source_dir = os.path.abspath(args.source_dir)
    output_dir = os.path.abspath(args.output_dir)
    
    if not os.path.exists(source_dir):
        log_error(f"Source directory does not exist: {source_dir}")
        sys.exit(1)
    
    # Create temp directory for import operations
    temp_dir = tempfile.mkdtemp(prefix='terraform_migration_')
    
    log_info(f"VastData Terraform State Migration Tool v{VERSION}")
    log_info("=" * 80)
    log_info(f"Source directory: {source_dir}")
    log_info(f"Output directory: {output_dir}")
    
    # Copy .tf files from source to temp directory for import
    log_info("Copying Terraform configuration files to temporary directory...")
    
    # Extract provider and terraform blocks from source .tf files
    provider_blocks = []
    terraform_blocks = []
    
    for file in os.listdir(source_dir):
        if file.endswith('.tf'):
            src_path = os.path.join(source_dir, file)
            with open(src_path, 'r') as f:
                content = f.read()
                
                # Simple extraction of provider and terraform blocks
                # This is a basic approach - we capture the blocks we need
                if 'provider "' in content:
                    # Extract provider block(s)
                    import re
                    providers = re.findall(r'provider\s+"[^"]+"\s+\{[^}]*\}', content, re.DOTALL)
                    provider_blocks.extend(providers)
                
                if 'terraform {' in content or 'required_providers' in content:
                    # Extract terraform block
                    start = content.find('terraform {')
                    if start != -1:
                        brace_count = 0
                        in_block = False
                        block_content = ""
                        for i, char in enumerate(content[start:]):
                            block_content += char
                            if char == '{':
                                brace_count += 1
                                in_block = True
                            elif char == '}':
                                brace_count -= 1
                                if brace_count == 0 and in_block:
                                    terraform_blocks.append(block_content)
                                    break
    
    # Create a minimal provider.tf file in temp directory
    provider_tf_path = os.path.join(temp_dir, 'provider.tf')
    with open(provider_tf_path, 'w') as f:
        f.write("# Auto-generated provider configuration for import\n\n")
        if terraform_blocks:
            f.write(terraform_blocks[0] + "\n\n")
        if provider_blocks:
            for provider in provider_blocks:
                f.write(provider + "\n\n")
    
    log_info("Created minimal provider configuration (resource definitions excluded)")
    
    # Get state file
    if args.s3_bucket:
        state_file = download_state_from_s3(
            args.s3_bucket,
            args.s3_key,
            args.s3_access_key,
            args.s3_secret_key,
            args.s3_endpoint
        )
    else:
        # Look for state file in source directory
        state_file = os.path.join(source_dir, 'terraform.tfstate')
        if not os.path.exists(state_file):
            log_error(f"State file not found: {state_file}")
            sys.exit(1)
    
    # Parse state
    state = parse_tfstate(state_file)
    
    # Extract raw resources for preservation
    raw_resources = extract_state_resources_raw(state)
    
    # Extract resources with structured format
    resources = extract_resources(state)
    
    if not resources:
        log_warning("No resources found in state file")
        shutil.rmtree(temp_dir)
        return
    
    # Separate VAST resources from non-VAST resources
    vast_resources, non_vast_resource_structs = separate_vast_resources(resources)
    
    # Get raw non-VAST resources for preservation
    non_vast_raw_resources = []
    if non_vast_resource_structs:
        non_vast_types = {r['type'] for r in non_vast_resource_structs}
        for raw_res in raw_resources:
            if raw_res.get('type') in non_vast_types:
                non_vast_raw_resources.append(raw_res)
        log_info(f"Preserving {len(non_vast_raw_resources)} non-VAST provider resources")
    
    if not vast_resources:
        log_warning("No VAST resources found in state file")
        if non_vast_raw_resources:
            log_info("Only non-VAST resources found - copying state as-is")
            os.makedirs(output_dir, exist_ok=True)
            shutil.copy(state_file, os.path.join(output_dir, 'terraform.tfstate'))
            log_success("State file copied successfully")
        shutil.rmtree(temp_dir)
        return
    
    # Generate stub resource definitions for VAST resources only
    resources_tf_path = os.path.join(temp_dir, "resources.tf")
    with open(resources_tf_path, 'w') as f:
        f.write("# Auto-generated stub resource definitions for VAST resources import\n")
        f.write("# Non-VAST resources will be preserved from original state\n\n")
        for resource in vast_resources:
            resource_type_from_state = resource['type']
            resource_type = RESOURCE_NAME_TRANSLATION.get(resource_type_from_state, resource_type_from_state)
            resource_name = resource['address'].split('.', 1)[1]
            f.write(f'resource "{resource_type}" "{resource_name}" {{\n')
            f.write('  # Configuration will be populated from import\n')
            f.write('}\n\n')
    
    # Create summary in temp directory (for reference during import)
    create_import_summary(vast_resources, temp_dir)
    
    # Clean up temp file if downloaded from S3
    if args.s3_bucket and state_file:
        try:
            os.unlink(state_file)
        except:
            pass
    
    log_success("=" * 80)
    log_info(f"Running parallel import with {args.max_workers} workers for {len(vast_resources)} VAST resources...")
    log_success("=" * 80)
    
    # Initialize terraform in temp directory
    log_info("Initializing Terraform...")
    init_result = subprocess.run(
        ['terraform', 'init'],
        cwd=temp_dir,
        capture_output=True,
        text=True
    )
    if init_result.returncode != 0 and 'dev_overrides' not in init_result.stdout:
        log_warning(f"Terraform init had issues: {init_result.stderr}")
    
    # Remove existing state if present
    state_path = os.path.join(temp_dir, 'terraform.tfstate')
    if os.path.exists(state_path):
        os.remove(state_path)
    
    # Execute parallel imports
    success_count, failed_count, failed_resources = run_parallel_imports(vast_resources, temp_dir, max_workers=args.max_workers)
    
    returncode = 0 if failed_count == 0 else 1
    
    log_success("=" * 80)
    if returncode == 0 or (failed_count > 0 and success_count > 0):
        # Create output directory
        os.makedirs(output_dir, exist_ok=True)
        
        # Get migrated state file path
        migrated_state = os.path.join(temp_dir, 'terraform.tfstate')
        output_state = os.path.join(output_dir, 'terraform.tfstate')
        
        if os.path.exists(migrated_state):
            # Merge non-VAST resources into the migrated state
            if non_vast_raw_resources:
                log_info("=" * 80)
                merge_non_vast_resources_to_state(migrated_state, non_vast_raw_resources, state)
            
            # Move the final merged state to output directory
            shutil.move(migrated_state, output_state)
            
            if returncode == 0:
                log_success("State migration complete!")
                log_success("=" * 80)
                print(f"\nMigration successful!")
                print(f"  - VAST resources imported: {success_count}")
                if non_vast_raw_resources:
                    print(f"  - Non-VAST resources preserved: {len(non_vast_raw_resources)}")
            else:
                log_warning("State migration completed with some failures")
                log_warning("=" * 80)
                print(f"\nMigration partially successful!")
                print(f"  - VAST resources imported: {success_count}")
                print(f"  - VAST resources failed: {failed_count}")
                if non_vast_raw_resources:
                    print(f"  - Non-VAST resources preserved: {len(non_vast_raw_resources)}")
                if failed_resources:
                    print(f"\nFailed resources:")
                    for res in failed_resources[:10]:  # Show first 10
                        print(f"    - {res}")
                    if len(failed_resources) > 10:
                        print(f"    ... and {len(failed_resources) - 10} more")
            
            print(f"\n  - Source directory (unchanged): {source_dir}")
            print(f"  - Migrated state file: {output_state}")
            print(f"\nNext steps:")
            print(f"  1. Update .tf files from {source_dir} to v2.x format")
            print(f"     Run: cd {os.path.dirname(os.path.abspath(__file__))}")
            print(f"          ./run_migration.sh {source_dir} <config_output_dir>")
            print(f"  2. Copy the migrated state to your config directory")
            print(f"     cp {output_state} <config_output_dir>/terraform.tfstate")
            print(f"  3. cd <config_output_dir>")
            print(f"  4. Run: terraform plan")
            if args.s3_bucket:
                print(f"  5. Upload terraform.tfstate to S3")
        else:
            log_error("Migrated state file not found!")
            sys.exit(1)
            
        # Don't exit with error if some resources succeeded
        if returncode != 0 and success_count > 0:
            returncode = 0
    else:
        log_error("State migration failed!")
        log_error("=" * 80)
        print(f"\nAll imports failed. Check the logs in: {temp_dir}")
        print(f"  - Resources file: {os.path.join(temp_dir, 'resources.tf')}")
        print(f"  - Summary: {temp_dir}/import_summary.txt")
        print(f"\nTemp directory preserved for debugging: {temp_dir}")
        sys.exit(1)
    
    # Clean up temp directory (only if migration succeeded)
    if returncode == 0:
        try:
            shutil.rmtree(temp_dir)
            log_info("Cleaned up temporary files")
        except Exception as e:
            log_warning(f"Could not clean up temp directory {temp_dir}: {e}")


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\n\nInterrupted by user")
        sys.exit(130)
    except Exception as e:
        log_error(f"Unexpected error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


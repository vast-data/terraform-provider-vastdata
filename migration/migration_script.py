#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

import os
import sys
import re
from pathlib import Path
import argparse

VERSION = "1.2.1"

# Resource type rename map (old → new)
resource_type_rename_map = {
    # Administrators (typo fix + plural→singular)
    "vastdata_administators_managers": "vastdata_administrator_manager",
    "vastdata_administators_roles": "vastdata_administrator_role",
    "vastdata_administators_realms": "vastdata_administrator_realm",
    # Other plural→singular renames
    "vastdata_kafka_brokers": "vastdata_kafka_broker",
    "vastdata_replication_peers": "vastdata_replication_peer",
    "vastdata_s3_replication_peers": "vastdata_s3_replication_peer",
    # Specific resource renames
    "vastdata_active_directory2": "vastdata_active_directory",
    "vastdata_non_local_user": "vastdata_nonlocal_user",
    "vastdata_non_local_user_key": "vastdata_nonlocal_user_key",
    "vastdata_non_local_group": "vastdata_nonlocal_group",
    "vastdata_saml": "vastdata_saml_config",
    # Additional resource renames from QA
    "vastdata_blockhost": "vastdata_block_host",
    "vastdata_s3_life_cycle_rule": "vastdata_s3_lifecycle_rule",
}

# Key groups for different Terraform block types
key_groups = {
    "Block List --> Attributes": [
        "capacity_total_limits", "capacity_limits", "static_limits", "static_total_limits",
        "default_group_quota", "default_user_quota", "share_acl", "owner_root_snapshot",
        "owner_tenant", "bucket_logging", "protocols_audit"
    ],
    "Block List --> Attributes List": ["frames"],
    "Block List --> Attributes Set": [],
    "Block List --> List of Maps": ["addresses", "group_quotas", "user_quotas"],
    "List of Number --> Set of Number": ["roles", "s3_policies_ids", "gids", "tenants"],
    "Block List --> List of List of String": ["client_ip_ranges", "ip_ranges"],
    "List of Number --> String": ["active_cnode_ids"],
    "List of String --> Set of String": [
        "object_types", "ldap_groups", "permissions_list", "groups", "users",
        "abac_tags", "hosts", "abe_protocols", "bucket_creators", "bucket_creators_groups",
        "nfs_all_squash", "nfs_no_squash", "nfs_read_only"
    ]
}

# Reverse lookup for key group by key (attribute-focused)
key_to_group = {k: g for g, keys in key_groups.items() for k in keys}

# Attributes to remove entirely (no longer supported or read-only in v2.0)
attributes_to_remove = {
    # Unsupported arguments
    "s3_bucket_full_control",
    
    # Read-only attributes that should not be set
    "tenant_name",
    "smb_directory_mode_padded", 
    "smb_file_mode_padded",
    "log_username",
    "log_hostname",
    "log_full_path", 
    "log_deleted",
    "enable_snapshot_lookup",
    "enable_listing_of_snapshot_dir",
    "data_modify",
    "data_create_delete",
    "data_read", 
    "cluster",
    "count_views",
}

def get_group_for_key(key):
    # Exact attribute key match
    for k in key_to_group:
        if k == key:
            return key_to_group[k]
    return None

def parse_nested_block(lines, start_index):
    attrs = {}
    i = start_index
    brace_level = 0
    block_lines = []
    while i < len(lines):
        line = lines[i]
        brace_level += line.count("{") - line.count("}")
        block_lines.append(line.strip())
        i += 1
        if brace_level == 0:
            break

    for line in block_lines:
        if "=" in line and not line.strip().endswith("{"):
            key, val = map(str.strip, line.split("=", 1))
            attrs[key] = val
    return attrs, i - start_index

def transform_resource_block(lines, i):
    if not lines[i].strip().startswith("resource "):
        return None, 0

    block = []
    brace_level = 0
    start = i
    while i < len(lines):
        block.append(lines[i])
        brace_level += lines[i].count("{") - lines[i].count("}")
        i += 1
        if brace_level == 0:
            break

    body_lines = block[1:-1]
    transformed = [block[0].rstrip()]
    j = 0

    ##########
    # Split the first line to get the resource type
    parts = transformed[0].split()
    current_resource_type = None
    if len(parts) > 1:
        resource_type = parts[1]  # This includes quotes, e.g., "vastdata_resource"
        resource_type_clean = resource_type.strip('"')  # Remove quotes for lookup
        current_resource_type = resource_type_clean
        # Check if the resource type is in the resource_type_rename_map
        new_resource_type = resource_type_rename_map.get(resource_type_clean)
        if new_resource_type:
            transformed[0] = transformed[0].replace(resource_type_clean, new_resource_type)
            current_resource_type = new_resource_type
    while j < len(body_lines):
        line = body_lines[j]
        stripped = line.strip()

        # Detect dynamic block, e.g. dynamic "client_ip_ranges" {
        dyn_match = re.match(r'dynamic\s+"(\w+)"\s*{', stripped)
        if dyn_match:
            # Preserve dynamic blocks as-is; auto-converting is error-prone without context
            brace = 0
            while j < len(body_lines):
                transformed.append(body_lines[j].rstrip())
                brace += body_lines[j].count("{") - body_lines[j].count("}")
                j += 1
                if brace == 0:
                    break
            continue  # Continue outer loop

        # Check normal block pattern e.g. client_ip_ranges { ... }
        m = re.match(r'(\w+)\s*{', stripped)
        if m:
            key = m.group(1)
            group = get_group_for_key(key)

            if group == "Block List --> List of List of String":
                pairs = []
                while j < len(body_lines) and body_lines[j].strip().startswith(key + " {"):
                    attrs, consumed = parse_nested_block(body_lines, j)
                    j += consumed
                    start_ip = attrs.get("start_ip", '""').strip()
                    end_ip = attrs.get("end_ip", '""').strip()
                    
                    # Remove existing quotes to check the actual value
                    start_ip_clean = start_ip.strip('"')
                    end_ip_clean = end_ip.strip('"')
                    
                    # Don't quote variable references or resource references
                    if start_ip_clean.startswith("var.") or start_ip_clean.startswith("vastdata_") or start_ip_clean.startswith("local.") or start_ip_clean.startswith("data."):
                        start_ip_formatted = start_ip_clean
                    else:
                        start_ip_formatted = f'"{start_ip_clean}"'
                    if end_ip_clean.startswith("var.") or end_ip_clean.startswith("vastdata_") or end_ip_clean.startswith("local.") or end_ip_clean.startswith("data."):
                        end_ip_formatted = end_ip_clean
                    else:
                        end_ip_formatted = f'"{end_ip_clean}"'
                    pairs.append(f'[{start_ip_formatted}, {end_ip_formatted}]')
                pairs_str = ',\n    '.join(pairs)
                transformed.append(f"  {key} = [\n    {pairs_str}\n  ]")
                continue

            if group == "Block List --> List of Maps":
                items = []
                while j < len(body_lines) and body_lines[j].strip().startswith(key + " {"):
                    attrs, consumed = parse_nested_block(body_lines, j)
                    j += consumed
                    items.append(attrs)
                transformed.append(f"  {key} = [")
                for item in items:
                    transformed.append("    {")
                    for k, v in item.items():
                        transformed.append(f"      {k} = {v}")
                    transformed.append("    },")
                transformed.append("  ]")
                continue

            if group == "Block List --> Attributes List":
                items = []
                while j < len(body_lines) and body_lines[j].strip().startswith(key + " {"):
                    attrs, consumed = parse_nested_block(body_lines, j)
                    j += consumed
                    items.append(attrs)
                transformed.append(f"  {key} = [")
                for item in items:
                    transformed.append("    {")
                    for k, v in item.items():
                        transformed.append(f"      {k} = {v}")
                    transformed.append("    },")
                transformed.append("  ]")
                continue

            if group == "Block List --> Attributes":
                attrs, consumed = parse_nested_block(body_lines, j)
                j += consumed
                transformed.append(f"  {key} = {{")
                for k, v in attrs.items():
                    transformed.append(f"    {k} = {v}")
                transformed.append("  }")
                continue

            if group == "Block List --> Attributes Set":
                items = []
                while j < len(body_lines) and body_lines[j].strip().startswith(key + " {"):
                    attrs, consumed = parse_nested_block(body_lines, j)
                    j += consumed
                    items.append(attrs)
                transformed.append(f"  {key} = [")
                for item in items:
                    transformed.append("    {")
                    for k, v in item.items():
                        transformed.append(f"      {k} = {v}")
                    transformed.append("    },")
                transformed.append("  ]")
                continue

        # Attribute-line transforms (e.g., List of Number --> String for specific keys)
        assign = re.match(r"(\s*)(\w+)\s*=\s*(.+)$", line)
        if assign:
            indent, attr_key, value_expr = assign.groups()
            original_attr_key = attr_key
            
            # Skip attributes that should be removed entirely
            if attr_key in attributes_to_remove:
                j += 1
                continue

            # Handle attribute name changes
            if attr_key == "type_":
                attr_key = "type"
            elif attr_key == "use32bit_fileid":
                attr_key = "use_32bit_fileid"
            elif attr_key == "permissions_list" and current_resource_type != "vastdata_administrator_manager":
                # For most resources: permissions_list -> permissions
                attr_key = "permissions"
            elif attr_key == "permissions" and current_resource_type == "vastdata_administrator_manager":
                # For administrator_manager specifically: permissions -> permissions_list
                attr_key = "permissions_list"

            group = get_group_for_key(attr_key)
            if group == "List of Number --> String":
                # Handle multi-line list by collecting all content until closing bracket
                if value_expr.strip() == "[":
                    # Multi-line list - collect all lines until closing bracket
                    all_content = ""
                    temp_j = j + 1
                    brace_count = 1
                    while temp_j < len(body_lines) and brace_count > 0:
                        next_line = body_lines[temp_j].strip()
                        all_content += " " + next_line
                        brace_count += next_line.count("[") - next_line.count("]")
                        temp_j += 1
                    
                    # Extract numbers from the collected content
                    numbers = re.findall(r"-?\d+", all_content)
                    joined = ",".join(numbers)
                    transformed.append(f"{indent}{attr_key} = \"{joined}\"")
                    j = temp_j
                    continue
                else:
                    # Single-line list
                    numbers = re.findall(r"-?\d+", value_expr)
                    joined = ",".join(numbers)
                    transformed.append(f"{indent}{attr_key} = \"{joined}\"")
                    j += 1
                    continue
            
            # If attribute name was changed but no group transformation needed
            if attr_key != original_attr_key:
                transformed.append(f"{indent}{attr_key} = {value_expr}")
                j += 1
                continue

        # No special case — just append line
        transformed.append(line.rstrip())
        j += 1

    transformed.append(block[-1].rstrip())
    return ("\n".join(transformed) + "\n"), i - start

def transform_variable_block(lines, i):
    # Optional: Implement variable block transformation if needed
    pass

def update_resource_references(content):
    """Update resource references when resource types have been renamed."""
    for old_type, new_type in resource_type_rename_map.items():
        # Pattern to match resource references like: vastdata_old_type.resource_name.attribute
        # This handles references in resource blocks, outputs, locals, etc.
        old_ref_pattern = rf'\b{old_type}\.([a-zA-Z_][a-zA-Z0-9_]*)'
        new_ref = rf'{new_type}.\1'
        content = re.sub(old_ref_pattern, new_ref, content)
        
        # Also handle data source references like: data.vastdata_old_type.resource_name.attribute
        old_data_pattern = rf'\bdata\.{old_type}\.([a-zA-Z_][a-zA-Z0-9_]*)'
        new_data_ref = rf'data.{new_type}.\1'
        content = re.sub(old_data_pattern, new_data_ref, content)
        
        # Also handle resource type names in comments and strings
        # This ensures consistency even in commented-out code
        old_resource_pattern = rf'\b{old_type}\b'
        content = re.sub(old_resource_pattern, new_type, content)
    return content

def transform_data_block(lines, i):
    """Transform data source blocks to update resource type names."""
    if not lines[i].strip().startswith("data "):
        return None, 0

    # Find the data block and extract the data source type
    data_line = lines[i].strip()
    # Pattern: data "vastdata_resource_type" "name" {
    data_match = re.match(r'data\s+"([^"]+)"\s+"([^"]+)"\s*{', data_line)
    if not data_match:
        return None, 0
    
    old_data_type = data_match.group(1)
    data_name = data_match.group(2)
    
    # Check if this data source type needs renaming
    new_data_type = resource_type_rename_map.get(old_data_type, old_data_type)
    
    # If no renaming needed, return original
    if new_data_type == old_data_type:
        return None, 0
    
    # Collect the entire data block
    block = []
    brace_level = 0
    j = i
    while j < len(lines):
        block.append(lines[j])
        brace_level += lines[j].count("{") - lines[j].count("}")
        j += 1
        if brace_level == 0:
            break
    
    # Replace the data source type in the first line
    block[0] = block[0].replace(f'"{old_data_type}"', f'"{new_data_type}"')
    
    return "".join(block), j - i

def transform_terraform_block(lines, i):
    """Transform terraform blocks to update VastData provider version from 1.x.x to 2.0.0."""
    if not lines[i].strip().startswith("terraform "):
        return None, 0

    # Collect the entire terraform block
    block = []
    brace_level = 0
    j = i
    while j < len(lines):
        block.append(lines[j])
        brace_level += lines[j].count("{") - lines[j].count("}")
        j += 1
        if brace_level == 0:
            break

    transformed_block = []
    provider_version_updated = False
    
    for line in block:
        # Look for vastdata provider version lines
        # Match patterns like: version = "1.7.0" or version = "1.x.x"
        version_match = re.search(r'(\s*version\s*=\s*")(1\.\d+\.\d+)(")', line)
        if version_match:
            # Check if we're within a vastdata provider block by looking at surrounding context
            block_str = "".join(block)
            # Look for vastdata provider context
            if 'vastdata' in block_str and 'vast-data/vastdata' in block_str:
                # Update version from 1.x.x to 2.0.0
                updated_line = line.replace(version_match.group(2), "2.0.0")
                transformed_block.append(updated_line)
                provider_version_updated = True
                continue
        
        transformed_block.append(line)
    
    # Only return transformed content if we actually updated something
    if provider_version_updated:
        return "".join(transformed_block), j - i
    else:
        return None, 0

def transform_file(input_path: Path, output_path: Path):
    with open(input_path, "r") as f:
        lines = f.readlines()

    transformed_lines = []
    i = 0
    while i < len(lines):
        line = lines[i]
        if line.strip().startswith("resource "):
            tr, consumed = transform_resource_block(lines, i)
            if consumed:
                transformed_lines.append(tr)
                i += consumed
                continue
        elif line.strip().startswith("data "):
            tr, consumed = transform_data_block(lines, i)
            if consumed:
                transformed_lines.append(tr)
                i += consumed
                continue
        elif line.strip().startswith("terraform "):
            tr, consumed = transform_terraform_block(lines, i)
            if consumed:
                transformed_lines.append(tr)
                i += consumed
                continue
        transformed_lines.append(line)
        i += 1

    # Join all transformed lines into a single string for reference updates
    content = "".join(transformed_lines)
    
    # Update resource references to use new resource type names
    content = update_resource_references(content)

    with open(output_path, "w") as f:
        f.write(content)



def main(src_folder, dst_folder):
    src_folder = Path(src_folder)
    dst_folder = Path(dst_folder)

    if not src_folder.is_dir():
        print(f"❌ Source folder '{src_folder}' does not exist or is not a directory")
        sys.exit(1)

    dst_folder.mkdir(parents=True, exist_ok=True)
    
    print("🚀 VastData Terraform File Converter")
    print("=" * 50)
    print(f"📂 Source: {src_folder}")
    print(f"📁 Output: {dst_folder}")
    print("=" * 50)

    converted_files = []

    for tf_file in src_folder.rglob("*.tf"):
        relative_path = tf_file.relative_to(src_folder)
        output_file = dst_folder / relative_path.with_name(tf_file.stem + "_converted.tf")
        output_file.parent.mkdir(parents=True, exist_ok=True)

        print(f"🔄 Converting {tf_file} -> {output_file}")
        transform_file(tf_file, output_file)
        converted_files.append(output_file)

    # Print summary
    print("\n" + "=" * 60)
    print("📋 CONVERSION COMPLETE")
    print("=" * 60)
    
    print(f"✅ Files converted: {len(converted_files)}")
    print(f"📁 Output location: {dst_folder}")
    
    print("\n" + "⚠️ " * 20)
    print("📋 NEXT STEPS - USER ACTION REQUIRED")
    print("⚠️ " * 20)
    print("1. 🔍 Review each converted file carefully")
    print("2. 💾 Backup your terraform state files") 
    print("3. ✅ Run 'terraform validate' to check syntax")
    print("4. 📋 Run 'terraform plan' to preview changes")
    print("5. 🧪 Test in non-production environment first")
    print("6. 🚀 Run 'terraform apply' only after thorough review")
    print("⚠️ " * 20)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Convert Terraform files for VastData provider v2.0 migration. "
                   "Performs file conversion only - user must validate and apply manually."
    )
    parser.add_argument("src_folder", nargs="?", help="Source folder containing .tf files")
    parser.add_argument("dst_folder", nargs="?", help="Destination folder for converted files")
    parser.add_argument("--version", action="version", version=f"%(prog)s {VERSION}")
    args = parser.parse_args()

    if not args.src_folder or not args.dst_folder:
        parser.print_help()
        sys.exit(1)

    main(args.src_folder, args.dst_folder)
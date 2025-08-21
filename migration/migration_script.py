#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

import os
import sys
import re
import subprocess
from pathlib import Path
import argparse

VERSION = "1.1.0"

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
    "List of Number --> String": ["cnode_ids", "active_cnode_ids"],
    "List of String --> Set of String": [
        "object_types", "ldap_groups", "permissions_list", "groups", "users",
        "abac_tags", "hosts", "abe_protocols", "bucket_creators", "bucket_creators_groups",
        "nfs_all_squash", "nfs_no_squash", "nfs_read_only"
    ]
}

# Reverse lookup for key group by key (attribute-focused)
key_to_group = {k: g for g, keys in key_groups.items() for k in keys}

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
    if len(parts) > 1:
        resource_type = parts[1]  # This includes quotes, e.g., "vastdata_resource"
        resource_type_clean = resource_type.strip('"')  # Remove quotes for lookup
        # Check if the resource type is in the resource_type_rename_map
        new_resource_type = resource_type_rename_map.get(resource_type_clean)
        if new_resource_type:
            transformed[0] = transformed[0].replace(resource_type_clean, new_resource_type)
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
            
            # Handle attribute name changes
            if attr_key == "type_":
                attr_key = "type"
            elif attr_key == "permissions_list":
                attr_key = "permissions"
            
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
        transformed_lines.append(line)
        i += 1

    # Join all transformed lines into a single string for reference updates
    content = "".join(transformed_lines)
    
    # Update resource references to use new resource type names
    content = update_resource_references(content)

    with open(output_path, "w") as f:
        f.write(content)

def terraform_apply(file_path: Path):
    work_dir = file_path.parent

    print(f"\nRunning 'terraform init' in {work_dir}")
    result_init = subprocess.run(["terraform", "init", "-input=false"], cwd=work_dir, capture_output=True, text=True)
    if result_init.returncode != 0:
        print(f"terraform init failed:\n{result_init.stderr}")
        return False
    print(result_init.stdout)

    print(f"Running 'terraform apply' on {file_path.name}")
    result_apply = subprocess.run(["terraform", "apply", "-auto-approve", file_path.name], cwd=work_dir, capture_output=True, text=True)
    if result_apply.returncode != 0:
        print(f"terraform apply failed:\n{result_apply.stderr}")
        return False
    print(result_apply.stdout)
    return True

def main(src_folder, dst_folder):
    src_folder = Path(src_folder)
    dst_folder = Path(dst_folder)

    if not src_folder.is_dir():
        print(f"Source folder '{src_folder}' does not exist or is not a directory")
        sys.exit(1)

    dst_folder.mkdir(parents=True, exist_ok=True)

    for tf_file in src_folder.rglob("*.tf"):
        relative_path = tf_file.relative_to(src_folder)
        output_file = dst_folder / relative_path.with_name(tf_file.stem + "_converted.tf")
        output_file.parent.mkdir(parents=True, exist_ok=True)

        print(f"\nConverting {tf_file} -> {output_file}")
        transform_file(tf_file, output_file)

        # Run terraform apply to verify
        success = terraform_apply(output_file)
        if not success:
            print(f"Apply failed for {output_file}")
        else:
            print(f"Apply succeeded for {output_file}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Convert and apply Terraform files in a folder.")
    parser.add_argument("src_folder", nargs="?", help="Source folder containing .tf files")
    parser.add_argument("dst_folder", nargs="?", help="Destination folder for converted files")
    parser.add_argument("--version", action="version", version=f"%(prog)s {VERSION}")
    args = parser.parse_args()

    if not args.src_folder or not args.dst_folder:
        parser.print_help()
        sys.exit(1)

    main(args.src_folder, args.dst_folder)
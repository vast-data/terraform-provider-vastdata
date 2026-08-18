#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

import os
import sys
import re
from pathlib import Path
import argparse

VERSION = "1.2.8"

# Marker inserted when a v1 attribute cannot be converted automatically.
MANUAL_STEP_PREFIX = "# TODO:"

VIP_POOL_REF_PATTERN = re.compile(r"vastdata_vip_pool\.(\w+)\.id")

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
    # Correct the naming: add underscores
    "vastdata_s3_lifecycle_rule": "vastdata_s3_life_cycle_rule",
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

# Per-resource-type attribute renames: {resource_type: {old_attr: new_attr}}
# Applied in addition to (and after) the global attribute rename logic.
resource_specific_renames = {
    # target_id was renamed to remote_target_id in v3
    "vastdata_protected_path": {
        "target_id": "remote_target_id",
    },
    # S3 object locking attribute renames in v3
    "vastdata_view": {
        "s3_locks": "locking",
        "s3_locks_retention_period": "default_retention_period",
    },
}

# Per-resource-type attributes to remove: {resource_type: set(attr_names)}
# vip_pools was removed from vastdata_view_policy in v3 and replaced by
# permission_per_vip_pool (Map of String: {pool_id: permission}).
# The migration cannot infer permissions from pool IDs alone, so the attribute
# is removed and a TODO comment is inserted so the user can add it manually.
resource_specific_removals = {
    "vastdata_view_policy": {"vip_pools"},
    "vastdata_replication_peer": {"peer_name", "is_local"},
    "vastdata_tenant": {"vippool_ids"},
    "vastdata_qos_policy": {"attached_users_identifiers"},
    "vastdata_vip_pool": {"active_cnode_ids"},
}

# Replacement comments inserted when a resource-specific attribute is removed.
resource_specific_removal_comments = {
    "vastdata_view_policy": {
        "vip_pools": (
            "# TODO: vip_pools was replaced by permission_per_vip_pool in v3.\n"
            '# Add: permission_per_vip_pool = { "<pool_id>" = "RW" }'
        ),
    },
    "vastdata_replication_peer": {
        "peer_name": (
            "# NOTE: peer_name is read-only in v3 and is populated by the provider."
        ),
        "is_local": (
            "# NOTE: is_local is read-only in v3 and is populated by the provider."
        ),
    },
    "vastdata_tenant": {
        "vippool_ids": (
            "# TODO: vippool_ids was removed from vastdata_tenant in v3.\n"
            "# MIGRATION REQUIRED: in v3 set vastdata_vip_pool.tenant_id on each pool.\n"
            "# VAST API rejects tenant_id changes when views exist on the pool's current tenant."
        ),
    },
    "vastdata_qos_policy": {
        "attached_users_identifiers": (
            "# TODO: attached_users_identifiers was removed in v3.\n"
            "# MIGRATION REQUIRED: add attached_users manually (literal user IDs cannot be\n"
            "# converted automatically — use vastdata_user.<name>.id in v1 config, or add\n"
            "# attached_users with name/fqdn/identifier_type/identifier_value by hand)."
        ),
    },
    "vastdata_vip_pool": {
        "active_cnode_ids": (
            "# NOTE: active_cnode_ids is read-only in v3 and is populated by the provider."
        ),
    },
}

# Attributes to inject for specific resource types if not already present.
# {resource_type: {attr_name: attr_value_as_hcl_string}}
# local_provider_id became required in v3 for vastdata_user. Default to 1
# (the built-in local LDAP provider that every VAST cluster ships with).
resource_specific_additions = {
    "vastdata_user": {
        "local_provider_id": "1",
    },
    "vastdata_group": {
        "local_provider_id": "1",
    },
}

def get_group_for_key(key):
    # Exact attribute key match
    for k in key_to_group:
        if k == key:
            return key_to_group[k]
    return None

def try_convert_dynamic_client_ip_ranges(body_lines, start_index):
    """Convert a standard dynamic client_ip_ranges block to a list attribute."""
    block_lines = []
    brace = 0
    j = start_index
    while j < len(body_lines):
        line = body_lines[j]
        block_lines.append(line)
        brace += line.count("{") - line.count("}")
        j += 1
        if brace == 0:
            break

    block_str = "\n".join(block_lines)
    for_each_match = re.search(r"for_each\s*=\s*(.+)", block_str)
    if not for_each_match:
        return None, 0

    for_each_expr = for_each_match.group(1).strip()
    has_bracket_access = (
        '["start_ip"]' in block_str
        or "['start_ip']" in block_str
        or '["end_ip"]' in block_str
        or "['end_ip']" in block_str
    )
    has_dot_access = (
        ".start_ip" in block_str
        or ".end_ip" in block_str
    )
    if not has_bracket_access and not has_dot_access:
        return None, 0

    if has_bracket_access:
        comprehension = (
            f'client_ip_ranges = [for r in {for_each_expr} : [r["start_ip"], r["end_ip"]]]'
        )
    else:
        comprehension = (
            f'client_ip_ranges = [for r in {for_each_expr} : [r.start_ip, r.end_ip]]'
        )

    return f"  {comprehension}", j - start_index

def convert_vippool_permissions_blocks(body_lines, start_index):
    """Convert vippool_permissions blocks to permission_per_vip_pool map."""
    permissions = {}
    j = start_index
    while j < len(body_lines) and body_lines[j].strip().startswith("vippool_permissions {"):
        attrs, consumed = parse_nested_block(body_lines, j)
        j += consumed
        pool_id = attrs.get("vippool_id", "").strip().strip('"')
        permission = attrs.get("vippool_permissions", "").strip().strip('"')
        if pool_id and permission:
            permissions[pool_id] = permission

    if not permissions:
        return None, 0

    lines = ["  permission_per_vip_pool = {"]
    for pool_id, permission in permissions.items():
        pool_id = pool_id.strip()
        if re.match(r'^\d+$', pool_id):
            key = f'"{pool_id}"'
        else:
            key = f"({pool_id})"
        lines.append(f'    {key} = "{permission}"')
    lines.append("  }")
    return "\n".join(lines), j - start_index


def try_convert_attached_users_identifiers(indent: str, value_expr: str):
    """Convert attached_users_identifiers to attached_users when it references a user resource.

    v1 configs commonly use [tostring(vastdata_user.foo.id)]. The v3 API stores the
    attachment as sid_str + user SID — reference vastdata_user.foo.sid (not invented values).
    Returns converted HCL or None when the expression cannot be converted safely.
    """
    match = re.search(r"vastdata_user\.(\w+)\.id", value_expr)
    if not match:
        return None
    user_ref = f"vastdata_user.{match.group(1)}"
    return (
        f"{indent}attached_users = [{{\n"
        f"{indent}  name             = {user_ref}.name\n"
        f'{indent}  fqdn             = ""\n'
        f'{indent}  identifier_type  = "sid_str"\n'
        f"{indent}  identifier_value = {user_ref}.sid\n"
        f"{indent}}}]"
    )


def read_assignment_value(body_lines, start_index):
    """Return (indent, attr_key, value_expr, end_index) for a full assignment."""
    assign = re.match(r"(\s*)(\w+)\s*=\s*(.*)$", body_lines[start_index])
    if not assign:
        return None
    indent, attr_key, remainder = assign.groups()
    end = skip_attribute_assignment(body_lines, start_index)
    parts = [remainder.strip()]
    for idx in range(start_index + 1, end):
        parts.append(body_lines[idx].strip())
    return indent, attr_key, " ".join(parts), end


def is_empty_vippool_ids_list(value_expr: str) -> bool:
    cleaned = re.sub(r"\s+", "", value_expr)
    return cleaned in ("[]",)


def try_convert_vippool_ids(indent: str, value_expr: str, tenant_resource_name: str):
    """Remove tenant.vippool_ids and document the v3 manual association model.

    v1: vastdata_tenant.vippool_ids = [vastdata_vip_pool.a.id, ...]
    v3: vastdata_vip_pool.a.tenant_id = vastdata_tenant.<tenant>.id

    tenant_id is NOT auto-injected: the VAST API rejects moving a pool to another
    tenant when views already exist on the pool's current tenant.
    """
    if is_empty_vippool_ids_list(value_expr):
        return "", True
    pool_names = VIP_POOL_REF_PATTERN.findall(value_expr)
    if not pool_names:
        return None, False
    pools = ", ".join(f"vastdata_vip_pool.{name}" for name in pool_names)
    note_lines = [
        f"{indent}# NOTE: vippool_ids removed from vastdata_tenant in v3.",
        (
            f"{indent}# In v3, set tenant_id on each pool if appropriate: {pools} "
            f"(vastdata_tenant.{tenant_resource_name})."
        ),
        (
            f"{indent}# WARNING: VAST API rejects tenant_id changes when views exist "
            "on the pool's current tenant."
        ),
    ]
    return "\n".join(note_lines), True


def skip_attribute_assignment(body_lines, start_index):
    """Return the line index after a complete attribute assignment."""
    line = body_lines[start_index]
    assign = re.match(r'(\s*)(\w+)\s*=\s*(.*)$', line)
    if not assign:
        return start_index + 1

    remainder = assign.group(3)
    square = remainder.count('[') - remainder.count(']')
    curly = remainder.count('{') - remainder.count('}')
    j = start_index + 1
    while square > 0 or curly > 0:
        if j >= len(body_lines):
            break
        next_line = body_lines[j]
        square += next_line.count('[') - next_line.count(']')
        curly += next_line.count('{') - next_line.count('}')
        j += 1
    return j

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
    current_resource_name = None
    header_match = re.match(
        r'resource\s+"([^"]+)"\s+"([^"]+)"',
        transformed[0].strip(),
    )
    if header_match:
        resource_type_clean = header_match.group(1)
        current_resource_name = header_match.group(2)
        current_resource_type = resource_type_clean
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
            dyn_key = dyn_match.group(1)
            if dyn_key == "client_ip_ranges":
                converted, consumed = try_convert_dynamic_client_ip_ranges(body_lines, j)
                if converted:
                    transformed.append(converted)
                    j += consumed
                    continue
            # Preserve other dynamic blocks as-is; auto-converting is error-prone without context
            brace = 0
            while j < len(body_lines):
                transformed.append(body_lines[j].rstrip())
                brace += body_lines[j].count("{") - body_lines[j].count("}")
                j += 1
                if brace == 0:
                    break
            continue  # Continue outer loop

        if stripped.startswith("vippool_permissions {"):
            converted, consumed = convert_vippool_permissions_blocks(body_lines, j)
            if converted:
                transformed.append(converted)
                j += consumed
                continue

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
            
            # Skip attributes that should be removed entirely (global)
            if attr_key in attributes_to_remove:
                j = skip_attribute_assignment(body_lines, j)
                continue

            # vastdata_qos_policy: auto-convert attached_users_identifiers when possible
            if (
                current_resource_type == "vastdata_qos_policy"
                and attr_key == "attached_users_identifiers"
            ):
                converted = try_convert_attached_users_identifiers(indent, value_expr)
                if converted:
                    transformed.append(converted)
                    j = skip_attribute_assignment(body_lines, j)
                    continue

            # vastdata_tenant: map vippool_ids -> vastdata_vip_pool.tenant_id
            if current_resource_type == "vastdata_tenant" and attr_key == "vippool_ids":
                assignment = read_assignment_value(body_lines, j)
                if assignment:
                    _, _, full_value, end_j = assignment
                    note, handled = try_convert_vippool_ids(
                        indent, full_value, current_resource_name or ""
                    )
                    if handled:
                        if note:
                            transformed.append(note)
                        j = end_j
                        continue

            # Skip attributes removed for this specific resource type, inserting a TODO comment
            if current_resource_type in resource_specific_removals and \
                    attr_key in resource_specific_removals[current_resource_type]:
                comment = resource_specific_removal_comments.get(
                    current_resource_type, {}
                ).get(attr_key)
                if comment:
                    for comment_line in comment.splitlines():
                        transformed.append(f"{indent}{comment_line}")
                j = skip_attribute_assignment(body_lines, j)
                continue

            # Handle attribute name changes
            if attr_key == "type_":
                attr_key = "type"
            elif attr_key == "use32bit_fileid":
                attr_key = "use_32bit_fileid"
            elif attr_key == "permissions" and current_resource_type == "vastdata_administrator_manager":
                # For administrator_manager: permissions -> permissions_list
                attr_key = "permissions_list"
            elif attr_key == "permissions_list" and current_resource_type != "vastdata_administrator_manager":
                # For administrator_role and other resources: permissions_list -> permissions
                # Exception: administrator_manager keeps permissions_list
                attr_key = "permissions"

            # Apply resource-specific attribute renames
            resource_renames = resource_specific_renames.get(current_resource_type, {})
            if attr_key in resource_renames:
                attr_key = resource_renames[attr_key]

            group = get_group_for_key(attr_key)

            # If attribute name was changed but no group transformation needed
            if attr_key != original_attr_key:
                transformed.append(f"{indent}{attr_key} = {value_expr}")
                j += 1
                continue

        # No special case — just append line
        transformed.append(line.rstrip())
        j += 1

    # Inject resource-specific additions that are not yet present.
    additions = resource_specific_additions.get(current_resource_type, {})
    for add_key, add_value in additions.items():
        # Only inject if the attribute doesn't already appear in the transformed block.
        already_present = any(
            re.match(rf"\s*{re.escape(add_key)}\s*=", line)
            for line in transformed
        )
        if not already_present:
            transformed.append(f"  {add_key} = {add_value}")

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
    """Transform terraform blocks to update VastData provider version from 1.x.x to 3.0.0."""
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
                # Update version from 1.x.x to 3.0.0
                updated_line = line.replace(version_match.group(2), "3.0.0")
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



def find_manual_steps(directory: Path):
    """Return (file, line_no, line_text) for every unresolved migration TODO."""
    findings = []
    for tf_file in sorted(directory.rglob("*_converted.tf")):
        try:
            lines = tf_file.read_text(encoding="utf-8").splitlines()
        except OSError:
            continue
        for line_no, line in enumerate(lines, start=1):
            if MANUAL_STEP_PREFIX in line:
                findings.append((tf_file, line_no, line.strip()))
    return findings


def report_manual_steps(findings, dst_folder: Path) -> None:
    print("\n" + "=" * 60)
    print("❌ MIGRATION INCOMPLETE — MANUAL STEPS REQUIRED")
    print("=" * 60)
    print(
        "The converter removed v1 attributes that cannot be mapped automatically.\n"
        "Fix every item below in the converted files, then re-run the migration script.\n"
        "Do not run terraform apply until this script exits successfully."
    )
    print(f"\nOutput directory: {dst_folder}\n")
    current_file = None
    for tf_file, line_no, line_text in findings:
        if tf_file != current_file:
            current_file = tf_file
            print(f"  {tf_file}:")
        print(f"    line {line_no}: {line_text}")
    print("\n" + "=" * 60)


def main(src_folder, dst_folder, allow_incomplete: bool = False):
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

    # Fail closed: do not leave the user with a silently incomplete stack.
    manual_steps = find_manual_steps(dst_folder)
    if manual_steps and not allow_incomplete:
        report_manual_steps(manual_steps, dst_folder)
        sys.exit(2)

    # Print summary
    print("\n" + "=" * 60)
    print("📋 CONVERSION COMPLETE")
    print("=" * 60)
    
    print(f"✅ Files converted: {len(converted_files)}")
    print(f"📁 Output location: {dst_folder}")
    
    if manual_steps:
        print(f"\n⚠️  {len(manual_steps)} manual step(s) remain (--allow-incomplete was set).")
    else:
        print("\n✅ No manual migration steps detected in converted files.")

    print("\n" + "⚠️ " * 20)
    print("📋 NEXT STEPS")
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
    parser.add_argument(
        "--allow-incomplete",
        action="store_true",
        help="Write converted files even when manual migration steps remain (default: exit with error)",
    )
    parser.add_argument("--version", action="version", version=f"%(prog)s {VERSION}")
    args = parser.parse_args()

    if not args.src_folder or not args.dst_folder:
        parser.print_help()
        sys.exit(1)

    main(args.src_folder, args.dst_folder, allow_incomplete=args.allow_incomplete)
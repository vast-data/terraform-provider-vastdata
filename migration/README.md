# VastData Terraform Provider Migration Tool

Migrate Terraform configurations from VastData provider **1.x** to **2.0**.

## Why This Tool is Needed

VastData provider **2.0** uses the new Terraform Plugin Framework and includes breaking changes:
- Resource type renames (e.g., `vastdata_administators_managers` → `vastdata_administrator_manager`)
- Attribute name changes (e.g., `type_` → `type`, `permissions_list` → `permissions`)
- Schema structure transformations (block lists to attributes, IP range formats, etc.)

## Quick Start

### Prerequisites
- Python 3.9 or higher
- Terraform CLI (for validation)

### Usage

```bash
# Run migration
./run_migration.sh /path/to/old/configs /path/to/migrated/configs

# Show help
./run_migration.sh --help

# Show version
./run_migration.sh --version

# Run tests
./run_migration.sh --test

# Clean up environment
./run_migration.sh --clean
```

### What It Does

1. **Automatically detects** and sets up Python 3.9+ environment
2. **Transforms resource types**: `vastdata_administators_managers` → `vastdata_administrator_manager`
3. **Updates attribute names**: `type_` → `type`, `permissions_list` → `permissions`
4. **Converts schema structures**: Block lists to attributes, IP ranges, etc.
5. **Preserves dynamic blocks** (requires manual review)
6. **Validates results** with `terraform init` and `terraform apply`

### Example Transformation

**Before (Provider 1.x):**
```hcl
resource "vastdata_administators_managers" "admin" {
  username         = "admin1"
  permissions_list = ["create_support", "create_settings"]
  
  capacity_limits {
    soft_limit = 1000
    hard_limit = 2000
  }
}
```

**After (Provider 2.0):**
```hcl
resource "vastdata_administrator_manager" "admin" {
  username    = "admin1"
  permissions = ["create_support", "create_settings"]
  
  capacity_limits = {
    soft_limit = 1000
    hard_limit = 2000
  }
}
```

## Important Notes

- **Backup** your configurations before migration
- **Review** converted files carefully
- **Test** with provider 2.0 before production use
- **Dynamic blocks** are preserved unchanged and may need manual review

## Support

For issues or questions, refer to the VastData Terraform provider documentation or create an issue in the repository.

---

**Migration Tool Entry Points: 5**
- Main migration: `./run_migration.sh source_dir dest_dir`
- Help: `./run_migration.sh --help`
- Version: `./run_migration.sh --version`
- Tests: `./run_migration.sh --test`
- Clean: `./run_migration.sh --clean`


# Manual Update Instructions for VAST Terraform Plugin v2.0.0

Version 2.0.0 of the VAST Terraform plugin introduces breaking changes, requiring manual updates to your `.tf` configuration files. Below are step-by-step instructions to update resource attributes for compatibility. Follow these carefully to avoid issues.

---

## Prerequisites
- *Backup* your `.tf` files before making changes.
- Identify resources using the VAST plugin (e.g., `vastdata_resource`).
- Review the attributes below to determine which need updates in your configs.

## Attribute Transformations
Update the following attributes in your `.tf` files based on their group. Use a text editor to locate and modify these attributes.

1. *Block List → Attributes*  
   _Attributes_: `capacity_total_limits`, `capacity_limits`, `static_limits`, `static_total_limits`, `default_group_quota`, `default_user_quota`, `share_acl`, `owner_root_snapshot`, `owner_tenant`, `bucket_logging`, `protocols_audit`  
   _Action_: Convert blocks to a single map.  
   _Example_:
   ```hcl
   # Before
   capacity_limits {
     key1 = value1
     key2 = value2
   }
   # After
   capacity_limits = {
     key1 = value1
     key2 = value2
   }
   ```

2. *Block List → Attributes List*  
   _Attribute_: `frames`  
   _Action_: Convert multiple blocks to a list of maps.  
   _Example_:
   ```hcl
   # Before
   frames {
     key1 = value1
   }
   frames {
     key1 = value2
   }
   # After
   frames = [
     { key1 = value1 },
     { key1 = value2 }
   ]
   ```

3. *Block List → Attributes Set*  
   _Attributes_: `group_quotas`, `user_quotas`  
   _Action_: Convert blocks to a list of maps (treated as a set for uniqueness).  
   _Example_:
   ```hcl
   # Before
   group_quotas {
     key1 = value1
   }
   # After
   group_quotas = [{ key1 = value1 }]
   ```

4. *List of Number → Set of Number*  
   _Attributes_: `roles`, `s3_policies_ids`, `gids`, `tenants`, `kafka_vip_pools`, `vip_pools`  
   _Action_: Ensure lists are unique (Terraform treats as sets). No structural change needed unless duplicates exist.  
   _Example_:
   ```hcl
   roles = [1, 2, 3] # Remains unchanged, ensure no duplicates
   ```

5. *Block List → List of List of String*  
   _Attributes_: `client_ip_ranges`, `ip_ranges`  
   _Action_: Convert blocks to a list of `[start_ip, end_ip]` lists.  
   _Example_:
   ```hcl
   # Before
   client_ip_ranges {
     start_ip = "192.168.1.1"
     end_ip = "192.168.1.10"
   }
   # After
   client_ip_ranges = [["192.168.1.1", "192.168.1.10"]]
   ```

6. *List of Number → String*  
   _Attribute_: `cnode_ids`  
   _Action_: Convert list to a comma-separated string.  
   _Example_:
   ```hcl
   # Before
   cnode_ids = [1, 2, 3]
   # After
   cnode_ids = "1,2,3" # Verify format in VAST docs
   ```

7. *List of String → Set of String*  
   _Attributes_: `object_types`, `ldap_groups`, `permissions_list`, `groups`, `users`, `abac_tags`, `hosts`, `abe_protocols`, `bucket_creators`, `bucket_creators_groups`, `nfs_all_squash`, `nfs_no_squash`, `nfs_read_only`, `nfs_read_write`, `nfs_root_squash`, `protocols`, `read_only`, `read_write`, `s3_read_only`, `s3_read_write`, `s3_visibility`, `s3_visibility_groups`, `smb_read_only`, `smb_read_write`, `trash_access`  
   _Action_: Ensure lists are unique (no structural change needed unless duplicates exist).  
   _Example_:
   ```hcl
   users = ["user1", "user2"] # Remains unchanged, ensure no duplicates
   ```

8. *Set of String → List of String*  
   _Attribute_: `urls`  
   _Action_: No change needed unless duplicates are allowed.  
   _Example_:
   ```hcl
   urls = ["url1", "url2"] # Remains unchanged
   ```

9. *Block → List of Maps*
   _Attribute_: `addresses`  
   _Action_: Convert block to list of maps with `host` and `port`.
   ```hcl
   # Before
   addresses {
     host = "10.131.21.121"
     port = 31485
   }
   # After
   addresses = [{ host = "10.131.21.121", port = 31485 }]
   ```

10. *Block → List of Maps*  
    _Attribute_: `share_acl`  
    _Action_: Convert `share_acl.acl` block to list of maps, rename `permissions` to `perm`.
   ```hcl
   # Before
   share_acl {
     acl {
       name = "user1"
       grantee = "users"
       fqdn = "All"
       permissions = "FULL"
     }
     enabled = true
   }
   # After
   share_acl = {
     acl = [{ name = "user1", grantee = "users", fqdn = "All", perm = "FULL" }]
     enabled = true
   }
   ```

11. *List to List of Maps*
    _Attribute_: `attached_users_identifiers`  
    _Action_: Convert list of strings to list of maps with `name`, `fqdn`, `identifier_type`, `identifier_value`.
   ```hcl
   # Before
   attached_users_identifiers = [tostring(vastdata_user.qos_user1.id)]
   # After
   attached_users = [{ name = "user1", fqdn = "user1.vastdb.local", identifier_type = "username", identifier_value = "user1" }]
   ```

## Resource Renames
Rename resources in `.tf` files using find-and-replace:

- `vastdata_blockhost` → `vastdata_block_host`
- `vastdata_non_local_user` → `vastdata_nonlocal_user` (including `data` resources)
- `vastdata_non_local_group` → `vastdata_nonlocal_group` (including `data` resources)
- `vastdata_s3_life_cycle_rule` → `vastdata_s3_lifecycle_rule`
- `vastdata_s3_replication_peers` → `vastdata_s3_replication_peer`
- `vastdata_replication_peers` → `vastdata_replication_peer`
- `vastdata_saml` → `vastdata_saml_config` (including `data` resources, move `idp_entityid`, `idp_metadata_url`, `encrypt_assertion`, `want_assertions_or_response_signed` to `saml_settings`)
- `vastdata_user_key` → `vastdata_nonlocal_user_key` (for non-local users)
- `vastdata_administators_managers` → `vastdata_administrator_manager`
- `vastdata_administators_roles` → `vastdata_administrator_role`
- `vastdata_administators_realms` → `vastdata_administrator_realm`
- `vastdata_active_directory2` → `vastdata_active_directory`
- `vastdata_kafka_brokers` → `vastdata_kafka_broker`

## Validation and Testing
1. Run `terraform validate` to check for syntax errors.
2. Run `terraform plan` to preview changes and ensure correctness.
3. Run `terraform apply` to apply the updated configurations.
4. Monitor for errors and verify resources in your VAST infrastructure.

## Notes
- *Backup*: Always keep a backup of your original `.tf` files.
- *Documentation*: Check the VAST plugin v2.0.0 docs for specific attribute formats (e.g., `cnode_ids` string format).
- *Testing*: Test in a non-production environment if possible.
- *Support*: Contact VAST support for complex cases or issues.

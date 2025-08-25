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

After conversion, **you** need to:
1. Review the converted files
2. Run `terraform validate` 
3. Run `terraform plan`
4. Run `terraform apply` (if you're satisfied with the changes)

### What It Does

This tool **ONLY** converts your Terraform configuration files:

1. **Updates provider version**: `version = "1.7.0"` → `version = "2.0.0"`
2. **Transforms resource types**: `vastdata_administators_managers` → `vastdata_administrator_manager`
3. **Updates attribute names**: `type_` → `type`, `permissions_list` → `permissions`
4. **Converts schema structures**: Block lists to attributes, IP ranges, etc.
5. **Preserves dynamic blocks** (may need manual review)
6. **Updates resource references** throughout your files

**This tool handles all the tedious file conversion work for you, so you can focus on reviewing and applying the changes.**

### Example Transformation

**Before (Provider 1.x):**
```hcl
terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
      version = "1.7.0"
    }
  }
}

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
terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
      version = "2.0.0"
    }
  }
}

resource "vastdata_administrator_manager" "admin" {
  username    = "admin1"
  permissions = ["create_support", "create_settings"]
  
  capacity_limits = {
    soft_limit = 1000
    hard_limit = 2000
  }
}
```

## 🔧 FILE CONVERSION CAPABILITIES

**Smart Terraform configuration converter for VastData provider v2.0 migration.**

The conversion tool:
- ✅ **Updates** VastData provider version from 1.x.x to 2.0.0
- ✅ **Converts** your `.tf` configuration files to the new format
- ✅ **Saves** converted files with new names (keeps originals safe)  
- ✅ **Updates** all resource references automatically
- ✅ **Handles** complex schema transformations intelligently
- ✅ **Preserves** your file structure and formatting

**You maintain complete control over when and how to apply the converted configurations.**

## Streamlined Conversion Workflow

### 1. Prepare for Migration
- **💾 Backup** your Terraform files and state
- **📊 Document** your current infrastructure

### 2. Convert Your Files
```bash
./run_migration.sh /path/to/old/configs /path/to/converted/configs
```

**What happens:**
- All `.tf` files are automatically converted to v2.0 format
- Converted files saved with `_converted.tf` suffix  
- Original files remain untouched
- Resource references updated throughout your configurations

### 3. Review and Deploy
- **🔍 Review** converted files (see exactly what changed)
- **📊 Validate** syntax with `terraform validate`
- **👀 Preview** changes with `terraform plan`
- **🧪 Test** in staging environment
- **🚀 Deploy** with `terraform apply`

## Best Practices

- **🧪 Test first** in staging before production deployment
- **💾 Keep backups** of your Terraform state and configuration files
- **🔍 Review changes** to understand what was converted
- **🚀 Use provider 2.0** for all new configurations after migration
- **🔧 Check dynamic blocks** - these are preserved as-is and may need adjustment

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

### After File Conversion (Your Responsibility)
After the converter finishes, you must:
1. **Review converted files** - Check each file for correctness
2. **Test in staging** - Apply in a non-production environment first
3. **Run `terraform plan`** - Verify planned changes match expectations
4. **Manual `terraform apply`** - Only after thorough review and testing

### Common Validation Issues and Solutions

**Error: "Provider configuration not found"**
- Ensure your `terraform` block includes the new VastData provider v2.0 configuration
- Update your provider source to the correct version

**Error: "Resource not found" during plan**
- This may be expected if you're migrating to new resource types
- Review the resource rename mappings in the migration summary

**Error: "Invalid attribute name"** 
- Check for dynamic blocks that may need manual adjustment
- Verify attribute transformations are correct for your use case

**Error: "Type mismatch"**
- Some attributes change types (e.g., list of numbers → comma-separated string)
- Review the transformed values to ensure they match expected format

## 🚀 Migration Success Tips

- **💾 Prepare thoroughly**: Backup files, state, and document your current setup
- **📖 Stay informed**: Review VastData provider v2.0.0 documentation for new features
- **🧪 Test confidently**: Use staging environments to validate conversions
- **🔍 Review systematically**: Understand each change before applying
- **📞 Leverage support**: VastData support team is available for complex scenarios
- **⏱️ Plan strategically**: Schedule appropriate maintenance windows
- **📋 Be prepared**: Have rollback procedures ready for peace of mind

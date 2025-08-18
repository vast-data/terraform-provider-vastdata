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
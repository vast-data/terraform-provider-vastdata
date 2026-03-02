# VastData Terraform Provider Migration Guide

Migrate your Terraform configurations and state files from VastData provider v1.x/v2.x to v3.0.

---

## Table of Contents

- [Overview](#overview)
- [Migration Tools](#migration-tools)
  - [1. Configuration Migration (.tf files)](#1-configuration-migration-tf-files)
  - [2. State Migration (.tfstate files)](#2-state-migration-tfstate-files)
- [Complete Migration Workflow](#complete-migration-workflow)
- [How Migrate Mode Works (Under the Hood)](#how-migrate-mode-works-under-the-hood)
- [Troubleshooting](#troubleshooting)
- [Support](#support)

---

## Overview

When upgrading to VastData Terraform provider v3.0, you need to migrate two things:

1. **Configuration files** (`.tf`) — resource type names and attribute formats have changed.
2. **State files** (`.tfstate`) — resource schemas have changed and the state must be rebuilt.

> **Important:** Your existing `.tfstate` file may contain resources from other providers
>  alongside VastData resources. The migration process preserves all
> non-VastData resources and only re-imports the VastData ones.

This guide covers both steps.

---

## Migration Tools

### 1. Configuration Migration (.tf files)

**Tool:** `migration_script.py` (wrapped by `run_migration.sh`)  
**Purpose:** Converts your `.tf` configuration files from v1.x/v2.x format to v3.0 format.

#### What It Does

- Renames resource types (e.g., `vastdata_administators_managers` → `vastdata_administrator_manager`)
- Updates attribute names and formats
- Converts block structures to attributes
- Preserves your original files (creates `*_converted.tf` output files)

#### Usage

```bash
# Using the wrapper script
./run_migration.sh /path/to/source/configs /path/to/output/configs

# Or directly with Python
python3 migration_script.py /path/to/source/configs /path/to/output/configs

# Show help
python3 migration_script.py --help
```

**Input:** Your existing `.tf` files (v1.x/v2.x format)  
**Output:** New `*_converted.tf` files with v3.0 syntax

---

### 2. State Migration (.tfstate files)

**Tool:** `state_migration.py`  
**Purpose:** Migrates your `.tfstate` file to work with the v3.0 provider.

#### Usage

```bash
# From the directory containing your converted v3.0 .tf files:
python3 /path/to/migration/state_migration.py /path/to/terraform.tfstate
```

The script will:
1. Parse the provided `.tfstate` file.
2. Remove all `vastdata_*` resources, keeping everything else.
3. Write the cleaned state as `terraform.tfstate` in the current directory.
4. Run `terraform init`.
5. Run `VASTDATA_MIGRATE_MODE=1 terraform apply -auto-approve`.

After completion, the resulting `terraform.tfstate` contains:
- All non-VastData resources exactly as they were in the original state.
- All VastData resources re-imported from the cluster under the v3.0 schema.

#### Options

| Flag | Description |
|------|-------------|
| `--dry-run` | Only strip the state and write the cleaned file; do **not** run Terraform |
| `--version` | Show version |
| `--help` | Show help |

#### Example: Dry Run

```bash
# Preview what the script will do without running Terraform
python3 state_migration.py /path/to/terraform.tfstate --dry-run
```

This writes the cleaned state but does not execute `terraform init` or `terraform apply`.
You can then inspect the state and run Terraform manually:

```bash
terraform init
VASTDATA_MIGRATE_MODE=1 terraform apply -auto-approve
```

---

## Complete Migration Workflow

### Prerequisites

```bash
# 1. Back up your existing configuration and state
cp -r /path/to/configs /path/to/configs.backup

# 2. Download your current state file
#    For S3 backend:
aws s3 cp s3://your-bucket/path/to/terraform.tfstate ./terraform.tfstate.original
#    For local state: just note the path to your existing terraform.tfstate
```

### Step 1: Migrate Configuration Files

```bash
cd /path/to/terraform-provider-vastdata/migration

# Convert .tf files from v1/v2 format to v3 format
./run_migration.sh /path/to/source/configs /path/to/output/configs
```

Review the generated `*_converted.tf` files to ensure correctness.

### Step 2: Prepare the Working Directory

The output directory from Step 1 already contains only the converted `.tf` files and
no Terraform state — so you can continue working directly in it.

```bash
cd /path/to/output/configs
```

Make sure your provider block specifies v3.0:

```hcl
terraform {
  required_providers {
    vastdata = {
      source  = "vast-data/vastdata"
      version = "~> 3.0"
    }
  }
}
```

### Step 3: Migrate the State

Run the state migration script, pointing it at your **existing** `.tfstate` file:

```bash
python3 /path/to/migration/state_migration.py /path/to/terraform.tfstate.original
```

The script will:
- Strip all `vastdata_*` resources from the state (preserving AWS/GCP/etc.).
- Run `terraform init` to set up the v3.0 provider.
- Run `VASTDATA_MIGRATE_MODE=1 terraform apply -auto-approve` to re-read every
  VastData resource from the cluster and populate the state.

### Step 4: Verify — No Drift

```bash
# Run plan WITHOUT migrate mode — should show zero changes or only expected changes
terraform plan

# Expected:
# No changes. Your infrastructure matches the configuration.
```

If `terraform plan` shows unexpected differences, review and adjust your `.tf` files to
match the actual cluster state, then re-run the migration from Step 3.

### Step 5: Replace Your Old State

Once verification passes, replace your old state with the newly generated one:

```bash
# For S3 backend: upload the new state to replace the old one
aws s3 cp terraform.tfstate s3://your-bucket/path/to/terraform.tfstate

# For local state: copy to your production directory
cp terraform.tfstate /path/to/production/configs/terraform.tfstate
```

From this point on, use the v3.0 provider normally — `VASTDATA_MIGRATE_MODE` is no
longer needed.

---

## Support

For issues or questions:
- Check the [CHANGELOG.md](../CHANGELOG.md) for breaking changes
- Create an issue in the repository
- Contact VastData support

---

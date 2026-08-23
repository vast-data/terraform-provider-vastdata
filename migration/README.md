# VastData Terraform Provider Migration Guide

Migrate your Terraform configurations and state files from the VastData provider 1.x to a newer provider version.

---

## Table of Contents

- [Overview](#overview)
- [Migration Tools](#migration-tools)
  - [1. Configuration Migration (.tf files)](#1-configuration-migration-tf-files)
  - [2. State Migration (.tfstate files)](#2-state-migration-tfstate-files)
- [Complete Migration Workflow](#complete-migration-workflow)
- [Manual Migration Steps](#manual-migration-steps)
- [How Migrate Mode Works (Under the Hood)](#how-migrate-mode-works-under-the-hood)
- [Troubleshooting](#troubleshooting)
- [Support](#support)

---

## Overview

When upgrading from the VastData Terraform provider 1.x to a newer version, you need to migrate two things:

1. **Configuration files** (`.tf`) — resource types and attribute schemas changed between provider 1.x and newer provider versions.
2. **State files** (`.tfstate`) — the state must be rebuilt to match the new resource schemas.

> **Important:** Your existing `.tfstate` file may contain resources from other providers
>  alongside VastData resources. The migration process preserves all
> non-VastData resources and only re-imports the VastData ones.

This guide covers both steps.

---

## Migration Tools

### 1. Configuration Migration (.tf files)

**Tool:** `migration_script.py` (wrapped by `run_migration.sh`)  
**Purpose:** Updates your `.tf` configuration files so resource definitions match the schemas used by newer provider versions (they differ from provider 1.x).

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

**Input:** Your existing `.tf` files written for provider 1.x  
**Output:** New `*_converted.tf` files updated for newer provider versions

---

### 2. State Migration (.tfstate files)

**Tool:** `state_migration.py`  
**Purpose:** Rebuilds your `.tfstate` file so VastData resources match the schemas used by newer provider versions.

#### Usage

```bash
# From the directory containing your converted .tf files:
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
- All VastData resources re-imported from the cluster under the newer provider schema.

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

# Convert .tf files from provider 1.x schemas to newer provider schemas
./run_migration.sh /path/to/source/configs /path/to/output/configs
```

Review the generated `*_converted.tf` files to ensure correctness, then complete any
[manual migration steps](#manual-migration-steps) before continuing.

### Step 2: Prepare the Working Directory

The output directory from Step 1 already contains only the converted `.tf` files and
no Terraform state — so you can continue working directly in it.

```bash
cd /path/to/output/configs
```

Make sure your provider block pins a release compatible with your VAST cluster version:

```hcl
terraform {
  required_providers {
    vastdata = {
      source  = "vast-data/vastdata"
      version = "~> 3.2.2"  # example — pick the release that matches your cluster
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
- Run `terraform init` to set up the newer provider.
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

From this point on, use the provider normally — `VASTDATA_MIGRATE_MODE` is no
longer needed.

---

## Manual Migration Steps

The configuration migration script handles most schema differences between provider 1.x
and newer provider versions automatically, but some changes cannot be inferred from your
existing files.

### Provider version

The script updates the `vastdata` provider pin to `3.0.0`. Adjust this to the release
that matches your VAST cluster version (for example, `~> 3.2.2` on release-5.4.x
clusters). See the [provider releases](https://registry.terraform.io/providers/vast-data/vastdata/latest)
and [CHANGELOG.md](../CHANGELOG.md) for compatibility details.

### `vastdata_protected_path`: add `tenant_id`

In newer provider versions, `tenant_id` is **required** on every `vastdata_protected_path`
resource. In provider 1.x it was optional, so older configurations often omit it. The
migration script does **not** add `tenant_id` automatically and will not flag the
omission — you must set it by hand in each `*_converted.tf` file that defines a
protected path.


## Support

For issues or questions:
- Check the [CHANGELOG.md](../CHANGELOG.md) for breaking changes
- Create an issue in the repository
- Contact VastData support

---

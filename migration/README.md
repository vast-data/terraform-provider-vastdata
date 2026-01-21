# VastData Terraform Provider Migration Guide

Migrate your Terraform configurations and state files from VastData provider v1.x/v2.x to v3.0.

---

## Table of Contents

- [Overview](#overview)
- [Migration Tools](#migration-tools)
  - [1. Configuration Migration (.tf files)](#1-configuration-migration-tf-files)
  - [2. State Migration (.tfstate files)](#2-state-migration-tfstate-files)
- [Quick Start](#quick-start)
- [Complete Migration Workflow](#complete-migration-workflow)
- [Troubleshooting](#troubleshooting)

---

## Overview

When upgrading to VastData Terraform provider v3.0, you need to migrate:

1. **Configuration files** (`.tf`) - resource names and attribute formats have changed
2. **State files** (`.tfstate`) - resource schemas have changed

This directory provides automated tools for both migrations.

---

## Migration Tools

### 1. Configuration Migration (.tf files)

**Tool:** `migration_script.py`  
**Purpose:** Converts your `.tf` configuration files to v3.0 format

#### What It Does

- Renames resources (e.g., `vastdata_administators_managers` → `vastdata_administrator_manager`)
- Updates attribute names and formats
- Converts block structures to attributes
- Preserves your original files (creates `*_converted.tf` files)

#### Usage

```bash
# Basic usage
./run_migration.sh /path/to/source/configs /path/to/output/configs

# Show help
./run_migration.sh --help
```

**Input:** Your existing `.tf` files  
**Output:** New `*_converted.tf` files with v3.0 syntax

---

### 2. State Migration (.tfstate files)

**Tool:** `state_migration.py`  
**Purpose:** Migrates your Terraform state file to v3.0 by re-importing all resources

#### What It Does

- Extracts all resources from your existing state file
- Generates `terraform import` commands for each resource
- Creates a migration script you can review before executing
- Supports both local and S3-backed state files

#### Usage

**Local state file:**

```bash
python3 state_migration.py /path/to/terraform/config /path/to/output
```

**S3-backed state file:**

```bash
python3 state_migration.py \
  /path/to/terraform/config \
  /path/to/output \
  --s3-bucket my-bucket \
  --s3-key terraform.tfstate \
  --s3-access-key YOUR_ACCESS_KEY \
  --s3-secret-key YOUR_SECRET_KEY

# With custom S3 endpoint (VAST S3)
python3 state_migration.py \
  /path/to/terraform/config \
  /path/to/output \
  --s3-bucket my-bucket \
  --s3-key terraform.tfstate \
  --s3-access-key YOUR_ACCESS_KEY \
  --s3-secret-key YOUR_SECRET_KEY \
  --s3-endpoint https://s3.vast.example.com
```

**Input:** Existing `terraform.tfstate`  
**Output:** Directory where new migrated state file will be created

---

## Quick Start

### Prerequisites

```bash
# Install Python dependencies (for state migration)
pip install -r requirements-state.txt

# Backup everything
cp -r /path/to/configs /path/to/configs.backup
cp terraform.tfstate terraform.tfstate.backup
```

### Step 1: Migrate Configuration Files

```bash
cd /path/to/terraform-provider-vastdata/migration
./run_migration.sh /path/to/configs /path/to/output
```

### Step 2: Update Provider Version

Edit your `versions.tf`:

```hcl
terraform {
  required_providers {
    vastdata = {
      source  = "vastdata/vastdata"
      version = "~> 3.0"
    }
  }
}
```

### Step 3: Migrate State File

```bash
# For local state
python3 state_migration.py /path/to/configs /path/to/output

# For S3 state
python3 state_migration.py /path/to/configs /path/to/output \
  --s3-bucket my-bucket --s3-key terraform.tfstate \
  --s3-access-key KEY --s3-secret-key SECRET
```

### Step 4: Verify Migration

```bash
cd /path/to/output
terraform validate
```

## Support

For issues or questions:
- Check the [CHANGELOG.md](../CHANGELOG.md) for breaking changes
- Create an issue in the repository
- Contact VastData support

---

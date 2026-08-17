# Copyright (c) HashiCorp, Inc.
#
# StackAV customer migration reproduction scenario.
# Covers all five issues reported during their v1 → v3 migration:
#
#  1. bucket_logging stored as array in old state (state_migration.py handles this)
#  2. target_id not recognised in vastdata_protected_path (renamed → remote_target_id)
#  3. s3_life_cycle_rule ambiguous lookup when same name spans multiple views
#  4. hard_limit_inodes = 0 treated as null after apply in vastdata_quota
#  5. vip_pools removed from vastdata_view_policy (replaced by permission_per_vip_pool)

# ── Variables ─────────────────────────────────────────────────────────────────

variable stackav_bucket_name {
  type = string
}

variable stackav_bucket_owner {
  type = string
}

variable stackav_policy_name {
  type = string
}

variable stackav_protected_path_name {
  type = string
}

variable stackav_target_object_id {
  type = number
}

variable stackav_hard_limit_inodes {
  type    = number
  default = 0
}

variable stackav_pool_name {
  type = string
}

variable stackav_pool_start_ip {
  type = string
}

variable stackav_pool_end_ip {
  type = string
}

# ── Issue 5: vip_pools on vastdata_view_policy ────────────────────────────────
# Old v1 usage: vip_pools = [vastdata_vip_pool.stackav_pool.id]
# Removed in v3 → must be replaced with permission_per_vip_pool manually.

resource vastdata_vip_pool stackav_pool {
  name       = var.stackav_pool_name
  role       = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    start_ip = var.stackav_pool_start_ip
    end_ip   = var.stackav_pool_end_ip
  }
}

resource vastdata_view_policy stackav_policy {
  name      = var.stackav_policy_name
  flavor    = "S3_NATIVE"
  vip_pools = [vastdata_vip_pool.stackav_pool.id]
}

# ── Issue 1: bucket_logging as old block-list on vastdata_view ────────────────
# Old v1 block syntax is transformed to attribute object by migration_script.py.
# The old *state* storing it as an array is handled by state_migration.py.

resource vastdata_view stackav_view {
  path         = "/${var.stackav_bucket_name}"
  bucket       = var.stackav_bucket_name
  bucket_owner = var.stackav_bucket_owner
  policy_id    = vastdata_view_policy.stackav_policy.id
  create_dir   = true
  protocols    = ["S3"]
  bucket_logging {
    destination_id = 42
    prefix         = "logs/"
  }
}

# ── Issue 3: s3_life_cycle_rule ambiguous name across views ───────────────────
# Two views use a rule called "abort-incomplete-mpu-1-day". MigrateMode was
# searching by name only → "too many records found". Fixed by adding view_id
# to SearchableFields in s3_lifecycle_rule.go.

resource vastdata_s3_life_cycle_rule stackav_abort_mpu_1d {
  name                            = "abort-incomplete-mpu-1-day"
  prefix                          = ""
  view_id                         = vastdata_view.stackav_view.id
  abort_mpu_days_after_initiation = 1
  enabled                         = true
}

resource vastdata_s3_life_cycle_rule stackav_expiration_1d {
  name            = "expiration-1-day"
  prefix          = ""
  view_id         = vastdata_view.stackav_view.id
  expiration_days = 1
  noncurrent_days = 1
  enabled         = true
}

# ── Issue 4: hard_limit_inodes = 0 on vastdata_quota ─────────────────────────
# The VAST API omits zero-value inode fields from its response. Without
# PreserveUserValueFieldsWhenApiReturnsNull the provider would mark
# hard_limit_inodes as null after apply, triggering "inconsistent result" error.

resource vastdata_quota stackav_quota {
  name             = "sdr:bucket:stackav:app:dev:${var.stackav_bucket_name}"
  default_email    = "owner@example.com"
  path             = "/${var.stackav_bucket_name}"
  hard_limit       = 0
  hard_limit_inodes = var.stackav_hard_limit_inodes
  soft_limit       = 10000000000
  soft_limit_inodes = 0
  is_user_quota    = false
}

# ── Issue 2: target_id on vastdata_protected_path ─────────────────────────────
# Old v1 attribute name was target_id; renamed to remote_target_id in v3.
# migration_script.py now applies this rename automatically for this resource.

resource vastdata_protected_path stackav_protected_path {
  name                = "${var.stackav_protected_path_name}-to-venus"
  source_dir          = vastdata_view.stackav_view.path
  tenant_id           = vastdata_view.stackav_view.tenant_id
  target_exported_dir = "/stackav-vast-backups/${var.stackav_bucket_name}"
  protection_policy_id = 1
  target_id           = var.stackav_target_object_id
  capabilities        = "ASYNC_REPLICATION"
  enabled             = true
}

output stackav_view {
  value = vastdata_view.stackav_view
}

output stackav_quota {
  value = vastdata_quota.stackav_quota
}

output stackav_protected_path {
  value = vastdata_protected_path.stackav_protected_path
}

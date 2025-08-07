resource "vastdata_block_host" "vastdb_block_host" {
  name      = "vastdb-block-host"
  tenant_id = 1
  nqn       = "nqn.2014-08.org.nvmexpress:uuid:12345678-1234-1234-1234-123456789012"
}

# ---------------------
# Complete examples
# ---------------------

data "vastdata_tenant" "vastdb_default_tenant" {
  name = "default"
}

data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path                 = "/vastdb_view/subsystem"
  name                 = "vastdb-subsystem"
  create_dir           = true
  is_default_subsystem = true
  policy_id            = data.vastdata_view_policy.vastdb_view_policy_default.id
  protocols            = ["BLOCK"]
}

resource "vastdata_volume" "vastdb_volume" {
  name    = "vastdb-volume"
  size    = 10737418240
  view_id = vastdata_view.vastdb_view.id
}


resource "vastdata_block_host" "vastdb_block_host" {
  name      = "vastdb-block-host"
  tenant_id = data.vastdata_tenant.vastdb_default_tenant.id
  nqn       = "nqn.2014-08.org.nvmexpress:uuid:12345678-1234-1234-1234-123456789012"
}

resource "vastdata_block_host_mapping" "vastdb_block_host_mapping" {
  host_id   = vastdata_block_host.vastdb_block_host.id
  volume_id = vastdata_volume.vastdb_volume.id
}

# --------------------



resource "vastdata_view_policy" "vastdb_view_policy" {
  name          = "vastdb_view_policy"
  tenant_id     = 1
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
  permission_per_vip_pool = {
    "1" = "RW"
    "2" = "RW"
    "3" = "RW"
  }
}

# ---------------------
# Complete examples
# ---------------------

data "vastdata_tenant" "vastdb_default_tenant" {
  name = "default"
}

resource "vastdata_view_policy" "vastdb_view_policy" {
  name          = "vastdb_view_policy"
  tenant_id     = data.vastdata_tenant.vastdb_default_tenant.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

# --------------------


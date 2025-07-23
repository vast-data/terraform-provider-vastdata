
resource "vastdata_view_policy" "vastdb_view_policy" {
  name          = "vastdb_view_policy"
  vip_pools     = [1, 2, 3]
  tenant_id     = 1
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

# ---------------------
# Complete examples
# ---------------------


data "vastdata_vip_pool" "vastdb_vippool" {
  name = "vippool-1"
}

data "vastdata_tenant" "vastdb_default_tenant" {
  name = "default"
}

resource "vastdata_view_policy" "vastdb_view_policy" {
  name          = "vastdb_view_policy"
  vip_pools     = [data.vastdata_vip_pool.vastdb_vippool.id]
  tenant_id     = data.vastdata_tenant.vastdb_default_tenant.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

# --------------------


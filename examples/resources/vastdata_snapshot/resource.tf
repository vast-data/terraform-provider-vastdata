resource "vastdata_tenant" "tenant" {
  name = "tenant1"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
}

resource "vastdata_vip_pool" "pool" {
  name        = "protocols-pool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.80"
    start_ip = "11.0.0.70"
  }

}

resource "vastdata_view_policy" "view-policy" {
  name          = "view-policy1"
  vip_pools     = [vastdata_vip_pool.pool.id]
  tenant_id     = vastdata_tenant.tenant.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

resource "vastdata_view" "view" {
  path       = "/view1"
  policy_id  = vastdata_view_policy.view-policy.id
  tenant_id  = vastdata_tenant.tenant.id
  create_dir = "true"
}

resource "vastdata_snapshot" "snapshot" {
  name            = "snapshot1"
  path            = vastdata_view.view.path
  tenant_id       = vastdata_tenant.tenant.id
  indestructible  = false
  expiration_time = "2023-06-02T12:22:32Z"
}

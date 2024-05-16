#Create view policy with the name vpolicy1 
resource "vastdata_vip_pool" "pool3" {
  name        = "pool3"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.80"
    start_ip = "11.0.0.70"
  }

}


resource "vastresource" "vastdata_tenant" "tenant1" {
  name = "tenant01"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
}


resource "vastdata_view_policy" "vpolicy1" {
  name          = "vpolicy1"
  vip_pools     = [vastdata_vip_pool.pool3.id]
  tenant_id     = vastdata_tenant.tenant1.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

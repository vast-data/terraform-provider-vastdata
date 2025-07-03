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

#Note that the `vip_pools` attribute is supported with VAST Cluster version 5.0.0 and later.
#When using a VAST Cluster version later than 5.1.0, is is advisable to migrate using `vippool_permissions` in case you used the `vip_pools` attribute when the cluster was running version 5.0.0.
#When using `vip_pools` with version 5.1.0 and later, the virtual IP pool permissions will always be set to W.

#Example of using `vip_pools` with version 5.0.0 and later:

resource "vastdata_vip_pool" "pool4" {
  name        = "pool4"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.80"
    start_ip = "11.0.0.70"
  }

}

resource "vastdata_vip_pool" "pool5" {
  name        = "pool5"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.110"
    start_ip = "11.0.0.90"
  }

}


resource "vastresource" "vastdata_tenant" "tenant2" {
  name = "tenant02"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
  vippool_ids = [vastdata_vip_pool.pool4.id,vastdata_vip_pool.pool5.id]

}


resource "vastdata_view_policy" "vpolicy1" {
  name          = "vpolicy1"
  vip_pools     = [vastdata_vip_pool.pool3.id]
  tenant_id     = vastdata_tenant.tenant1.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
  vip_pools =  [vastdata_vip_pool.pool4.id,vastdata_vip_pool.pool5.id]
}

#Example of using `vippool_permissions` for version 5.1.0 and later, where one virtual IP pool is set to read/write and the other is set to read-only:

resource "vastdata_vip_pool" "pool4" {
  name        = "pool4"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.80"
    start_ip = "11.0.0.70"
  }

}

resource "vastdata_vip_pool" "pool5" {
  name        = "pool5"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.110"
    start_ip = "11.0.0.90"
  }

}


resource "vastresource" "vastdata_tenant" "tenant2"{
  name = "tenant02"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
  vippool_ids = [vastdata_vip_pool.pool4.id,vastdata_vip_pool.pool5.id]

}


resource "vastdata_view_policy" "vpolicy1" {
  name          = "vpolicy1"
  vip_pools     = [vastdata_vip_pool.pool3.id]
  tenant_id     = vastdata_tenant.tenant1.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
  vippool_permissions {
    vippool_id = vastdata_vip_pool.pool4.id
    vippool_permissions  = "RW"
  }
  vippool_permissions {
    vippool_id = vastdata_vip_pool.pool5.id
    vippool_permissions  = "RO"
  }  
}



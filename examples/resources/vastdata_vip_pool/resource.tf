#Basic definition of virtual IP pool 
resource "vastdata_vip_pool" "pool1" {
  name        = "pool1"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.40"
    start_ip = "11.0.0.20"
  }

  ip_ranges {
    start_ip = "11.0.0.5"
    end_ip   = "11.0.0.10"
  }
}

#A virtual IP pool can be associated with a tenant in two ways: by defining the `vastdata_tenant` resource, or by setting the `tenant_id` attribute on the virtual IP pool.
#It is recommended to use only one of these methods per tenant and pool combination.

#Define a virtual IP pool and associate it with a tenant using the `vastdata_tenant` resource
resource "vastdata_vip_pool" "pool1" {
  name        = "pool1"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.40"
    start_ip = "11.0.0.20"
  }

  ip_ranges {
    start_ip = "11.0.0.5"
    end_ip   = "11.0.0.10"
  }
}

resource "vastdata_tenant" "tenant1" {
  name        = "tenant01"
  vippool_ids = [vastdata_vip_pool.pool1.id]
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
}


#Define a virtual IP pool and associate it with a tenant using the `tenant_id` attribute of the pool
resource "vastdata_vip_pool" "pool1" {
  name        = "pool1"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  tenant_id   = vastdata_tenant.tenant1.id
  ip_ranges {
    end_ip   = "11.0.0.40"
    start_ip = "11.0.0.20"
  }

  ip_ranges {
    start_ip = "11.0.0.5"
    end_ip   = "11.0.0.10"
  }
}

resource "vastdata_tenant" "tenant1" {
  name = "tenant01"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
}

#Define a virtual IP pool for all tenants by setting `tenant_id = 0`
resource "vastdata_vip_pool" "pool1" {
  name        = "pool1"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  tenant_id   = 0
  ip_ranges {
    end_ip   = "11.0.0.40"
    start_ip = "11.0.0.20"
  }

  ip_ranges {
    start_ip = "11.0.0.5"
    end_ip   = "11.0.0.10"
  }
}

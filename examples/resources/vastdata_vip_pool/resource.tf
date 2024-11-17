#Basic defenition of VIP pool 
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

#Setting up VIP pool related tenant can be done in 2 ways.
#It is advisable to select only one method per tenant,vippool combination.

#Define VIP pool setting up the tenant_id of it using vastdata_tenant resource 
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


#Define VIP pool setting up the tenant_id using the tenent_id attribute.
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

#Define a VIP pool for all tenants by setting tenant_id = 0
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

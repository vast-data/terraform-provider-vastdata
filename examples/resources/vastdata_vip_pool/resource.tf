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
  name = "tenant01"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
}

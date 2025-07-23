# Copyright (c) HashiCorp, Inc.

resource "vastdata_vip_pool" "vastdb_vippool" {
  name        = "vastdb_vippool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"

  ip_ranges = [
    ["11.0.0.6", "11.0.0.10"],
    ["11.0.0.20", "11.0.0.40"]
  ]
}
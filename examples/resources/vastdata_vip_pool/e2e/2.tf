# Copyright (c) HashiCorp, Inc.

resource "vastdata_vip_pool" "vastdb_vippool" {
  name                      = "vastdb_vippool"
  role                      = "PROTOCOLS"
  subnet_cidr               = "24"
  enable_weighted_balancing = true
  ip_ranges = [
    ["11.0.0.50", "11.0.0.80"],
  ]
}
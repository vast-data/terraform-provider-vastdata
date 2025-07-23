# Copyright (c) HashiCorp, Inc.

data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_vip_pool" "vastdb_vippool" {
  name          = "vastdb_vippool"
  role          = "PROTOCOLS"
  tenant_id     = data.vastdata_tenant.vastdb_tenant.id
  domain_name   = "vastdb.example.com"
  vms_preferred = true
  subnet_cidr   = "24"
  ip_ranges = [
    ["11.0.0.50", "11.0.0.80"],
  ]
}
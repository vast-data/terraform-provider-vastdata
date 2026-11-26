resource "vastdata_vip_pool" "dataengine" {
  name        = local.vip_pool_name
  role        = "PROTOCOLS"
  tenant_id   = vastdata_tenant.dataengine.id
  subnet_cidr = var.vip_pool.subnet_cidr
  ip_ranges   = var.vip_pool.ip_ranges
  gw_ip       = try(var.vip_pool.gw_ip, null)
  vlan        = try(var.vip_pool.vlan, null)
  domain_name = try(var.vip_pool.domain_name, null)
}

#A view when there is only one view with that name on the entire cluster
data "vastdata_view" "view1" {
  path = "/path"
}

#When there is more than one view with the same path at differant tenants
#If a tenant_id is not specfied, error is returned

data "vastdata_tenant" "tenants1" {
  name = "tenant01"
}

data "vastdata_view" "view1" {
  path      = "/path2"
  tenant_id = data.vastdata_tenant.tenants1.id
}

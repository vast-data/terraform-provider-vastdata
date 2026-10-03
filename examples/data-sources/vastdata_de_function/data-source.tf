# Lookup by name (requires tenant_name for cluster-admin auth).
data "vastdata_de_function" "by_name" {
  name        = "de-test-function"
  tenant_name = "de-lab"
}

# Lookup by guid.
data "vastdata_de_function" "by_guid" {
  guid        = "00000000-0000-0000-0000-000000000001"
  tenant_name = "de-lab"
}

# Lookup by vrn.
data "vastdata_de_function" "by_vrn" {
  vrn         = "vast:dataengine:functions:de-test-function"
  tenant_name = "de-lab"
}

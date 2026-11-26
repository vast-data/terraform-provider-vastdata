# Lookup by name (requires tenant_name for cluster-admin auth).
data "vastdata_de_mtls_authentication_credential" "by_name" {
  name        = "de-lab-mtls"
  tenant_name = "de-lab"
}

# Lookup by guid.
data "vastdata_de_mtls_authentication_credential" "by_guid" {
  guid        = "00000000-0000-0000-0000-000000000001"
  tenant_name = "de-lab"
}

# Lookup by vrn.
data "vastdata_de_mtls_authentication_credential" "by_vrn" {
  vrn         = "vast:dataengine:mtls-authentication-credentials:de-lab-mtls"
  tenant_name = "de-lab"
}

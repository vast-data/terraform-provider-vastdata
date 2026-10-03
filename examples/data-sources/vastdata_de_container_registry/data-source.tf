# Lookup by name (requires tenant_name for cluster-admin auth).
data "vastdata_de_container_registry" "by_name" {
  name        = "dockerhub"
  tenant_name = "de-lab"
}

# Lookup by guid.
data "vastdata_de_container_registry" "by_guid" {
  guid        = "00000000-0000-0000-0000-000000000001"
  tenant_name = "de-lab"
}

# Lookup by vrn.
data "vastdata_de_container_registry" "by_vrn" {
  vrn         = "vast:dataengine:container-registries:dockerhub"
  tenant_name = "de-lab"
}

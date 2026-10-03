# Lookup by name (requires tenant_name for cluster-admin auth).
data "vastdata_de_pipeline" "by_name" {
  name        = "de-test-pipeline"
  tenant_name = "de-lab"
}

# Lookup by guid.
data "vastdata_de_pipeline" "by_guid" {
  guid        = "00000000-0000-0000-0000-000000000001"
  tenant_name = "de-lab"
}

# Lookup by vrn.
data "vastdata_de_pipeline" "by_vrn" {
  vrn         = "vast:dataengine:pipelines:de-test-pipeline"
  tenant_name = "de-lab"
}

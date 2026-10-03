# Lookup by name (requires tenant_name for cluster-admin auth).
data "vastdata_de_kubernetes_cluster" "by_name" {
  name        = "de-lab-k8s"
  tenant_name = "de-lab"
}

# Lookup by guid.
data "vastdata_de_kubernetes_cluster" "by_guid" {
  guid        = "00000000-0000-0000-0000-000000000001"
  tenant_name = "de-lab"
}

# Lookup by vrn.
data "vastdata_de_kubernetes_cluster" "by_vrn" {
  vrn         = "vast:dataengine:kubernetes-clusters:de-lab-k8s"
  tenant_name = "de-lab"
}

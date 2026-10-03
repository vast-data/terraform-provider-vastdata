# Attach an external Kubernetes cluster to Data Engine using mTLS credentials.
# Requires vastdata_de_engine and vastdata_de_mtls_authentication_credential.

provider "vastdata" {
  host     = "vast-cluster"
  username = "admin"
  password = "******"
  # tenant = "de-lab"
}

resource "vastdata_de_mtls_authentication_credential" "example" {
  name            = "de-mtls"
  tenant_name     = "de-lab" # required for cluster-admin auth
  kubeconfig_data = file("${path.module}/kubeconfig")
  context_name    = "my-cluster-context"
  # or PEMs: certificate_authority_b64 / client_certificate_b64 / client_key_b64
}

resource "vastdata_de_kubernetes_cluster" "example" {
  name                  = "de-test"
  tenant_name           = "de-lab"
  kube_api_url          = "https://k8s-api.example.com:6443"
  mtls_credentials_guid = vastdata_de_mtls_authentication_credential.example.guid
  namespaces = [
    "default",
    "vast-dataengine",
    "vast-pipelines",
  ]
}

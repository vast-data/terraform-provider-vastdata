# Data Engine mTLS credentials used when attaching a Kubernetes cluster.
# Requires vastdata_de_engine to be completed first.
#
# Provide exactly one of:
#   A) kubeconfig_data (+ optional context_name) — write-only, not stored in state
#   B) certificate_authority_b64 + client_certificate_b64 + client_key_b64

provider "vastdata" {
  host     = "vast-cluster"
  username = "admin"
  password = "******"
}

# --- Option A: kubeconfig + context -------------------------------------------
resource "vastdata_de_mtls_authentication_credential" "from_kubeconfig" {
  name            = "de-mtls"
  tenant_name     = "de-lab" # required for cluster-admin auth
  kubeconfig_data = file("${path.module}/kubeconfig")
  context_name    = "my-cluster-context" # optional; defaults to current-context
}

# --- Option B: explicit base64 PEMs -------------------------------------------
resource "vastdata_de_mtls_authentication_credential" "from_pems" {
  name                      = "de-mtls-pems"
  tenant_name               = "de-lab"
  certificate_authority_b64 = filebase64("${path.module}/ca.crt")
  client_certificate_b64    = filebase64("${path.module}/client.crt")
  client_key_b64            = filebase64("${path.module}/client.key")
}

# Register mTLS credentials + an external Kubernetes cluster with Data Engine.
# mTLS material: kubeconfig (preferred) or ca_cert/client_cert/client_key (*_b64 / artifacts dir).

check "kubernetes_cluster_requires_data_engine" {
  assert {
    condition     = !var.enable_kubernetes_cluster || var.enable_data_engine
    error_message = "enable_kubernetes_cluster requires enable_data_engine=true."
  }
}

check "kubernetes_cluster_requires_api_url" {
  assert {
    condition     = !var.enable_kubernetes_cluster || try(length(trimspace(var.kubernetes_kube_api_url)) > 0, false)
    error_message = "enable_kubernetes_cluster requires kubernetes_kube_api_url (e.g. https://<cp-ip>:6443)."
  }
}

check "kubernetes_cluster_requires_mtls_material" {
  assert {
    condition = !var.enable_kubernetes_cluster || (
      (local.k8s_use_kubeconfig && !local.k8s_has_pems) ||
      (!local.k8s_use_kubeconfig && local.k8s_has_pems)
    )
    error_message = <<-EOT
      enable_kubernetes_cluster requires exactly one mTLS source:
        • kubernetes_kubeconfig or kubernetes_kubeconfig_path, or
        • kubernetes_mtls_artifacts_dir (ca_cert/client_cert/client_key), or
        • kubernetes_certificate_authority_b64 + kubernetes_client_certificate_b64 + kubernetes_client_key_b64
      Do not mix kubeconfig with PEM inputs.
    EOT
  }
}

# kubeconfig_data is write-only (not in state), so Terraform cannot detect file refreshes.
# Hash the material and replace the credential when it changes.
resource "terraform_data" "mtls_material_hash" {
  count = var.enable_kubernetes_cluster ? 1 : 0

  input = local.k8s_use_kubeconfig ? sha256(local.k8s_kubeconfig_data) : sha256(join("|", [
    coalesce(local.k8s_ca_b64, ""),
    coalesce(local.k8s_client_cert_b64, ""),
    coalesce(local.k8s_client_key_b64, ""),
  ]))
}

resource "vastdata_de_mtls_authentication_credential" "cluster" {
  count = var.enable_kubernetes_cluster ? 1 : 0

  name        = local.k8s_mtls_name
  tenant_name = vastdata_tenant.dataengine.name

  # Option A: kubeconfig (provider parses CA + client cert/key; kubeconfig_data is write-only).
  kubeconfig_data = local.k8s_use_kubeconfig ? local.k8s_kubeconfig_data : null
  context_name    = local.k8s_use_kubeconfig ? var.kubernetes_kubeconfig_context : null

  # Option B: explicit base64 PEMs (from *_b64 vars or artifacts dir files).
  certificate_authority_b64 = local.k8s_use_kubeconfig ? null : local.k8s_ca_b64
  client_certificate_b64    = local.k8s_use_kubeconfig ? null : local.k8s_client_cert_b64
  client_key_b64            = local.k8s_use_kubeconfig ? null : local.k8s_client_key_b64

  depends_on = [vastdata_de_engine.this]

  lifecycle {
    replace_triggered_by = [terraform_data.mtls_material_hash[count.index]]
  }
}

resource "vastdata_de_kubernetes_cluster" "this" {
  count = var.enable_kubernetes_cluster ? 1 : 0

  name                  = local.k8s_cluster_name
  tenant_name           = vastdata_tenant.dataengine.name
  kube_api_url          = var.kubernetes_kube_api_url
  mtls_credentials_guid = vastdata_de_mtls_authentication_credential.cluster[0].guid
  namespaces            = var.kubernetes_namespaces
  is_internal           = var.kubernetes_is_internal

  depends_on = [vastdata_de_engine.this]
}

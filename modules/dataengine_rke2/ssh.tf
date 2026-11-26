resource "local_sensitive_file" "ssh_key" {
  count = local.ssh_private_key_path == null ? 1 : 0

  content         = var.ssh_private_key
  filename        = "${path.module}/.terraform/ssh_id"
  file_permission = "0600"
}

locals {
  ssh_identity_file = coalesce(
    local.ssh_private_key_path,
    try(local_sensitive_file.ssh_key[0].filename, null)
  )

  ssh_env = {
    SSH_USER                     = var.ssh_user
    SSH_PORT                     = tostring(var.ssh_port)
    SSH_STRICT_HOST_KEY_CHECKING = var.ssh_strict_host_key_checking
    SSH_IDENTITY_FILE            = local.ssh_identity_file
    RKE2_CHANNEL                 = var.rke2_channel
    RKE2_TOKEN                   = var.cluster_token
    RKE2_WRITE_KUBECONFIG        = var.write_kubeconfig_mode
    RKE2_EXTRA_TLS_SANS          = join(" ", var.extra_tls_sans)
  }
}

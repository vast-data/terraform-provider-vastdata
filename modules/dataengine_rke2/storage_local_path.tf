resource "null_resource" "local_path_storage" {
  count = local.local_path_storage_enabled ? 1 : 0

  depends_on = [null_resource.finalize]

  triggers = {
    # Match prep/first_server/finalize so a cluster rebuild reinstalls the SC.
    fingerprint = local.install_fingerprint
    kubeconfig  = local.kubeconfig_abs
    version     = var.local_path_version
    default     = tostring(local.local_path_is_default)
  }

  provisioner "local-exec" {
    command     = abspath("${path.module}/scripts/install_local_path.sh")
    interpreter = ["bash"]
    environment = {
      KUBECONFIG         = local.kubeconfig_abs
      LOCAL_PATH_VERSION = var.local_path_version
      SET_DEFAULT        = tostring(local.local_path_is_default)
    }
  }
}

resource "null_resource" "podman_registry" {
  count = local.registry_enabled ? 1 : 0

  depends_on = [null_resource.finalize]

  triggers = {
    fingerprint = local.registry_fingerprint
    smoke       = tostring(var.registry_run_smoke_test)
  }

  provisioner "local-exec" {
    command     = abspath("${path.module}/scripts/run_remote_script.sh")
    interpreter = ["bash"]
    environment = merge(local.ssh_env, {
      REMOTE_HOST    = local.registry_ssh_host
      LOCAL_SCRIPT   = abspath("${path.module}/scripts/setup_podman_registry.sh")
      REMOTE_SCRIPT  = "/tmp/setup_podman_registry.sh"
      REGISTRY_HOST  = var.registry_host
      REGISTRY_PORT  = tostring(var.registry_port)
      NFS_EXPORT     = var.registry_nfs_export != null ? var.registry_nfs_export : ""
      MOUNTPOINT     = var.registry_mountpoint
      CERT_DAYS      = tostring(var.registry_cert_days)
      FORCE_RECERT   = tostring(var.registry_force_recert)
      RUN_SMOKE_TEST = tostring(var.registry_run_smoke_test)
      FORWARD_VARS   = "REGISTRY_HOST REGISTRY_PORT NFS_EXPORT MOUNTPOINT CERT_DAYS FORCE_RECERT RUN_SMOKE_TEST CONTAINER_NAME"
    })
  }
}

resource "null_resource" "registry_trust" {
  count = local.registry_enabled ? 1 : 0

  depends_on = [null_resource.podman_registry]

  triggers = {
    fingerprint = local.registry_fingerprint
    nodes       = join(",", local.rke2_node_mgmt_ips)
    restart     = tostring(var.registry_restart_rke2)
  }

  provisioner "local-exec" {
    command     = abspath("${path.module}/scripts/configure_registry_trust.sh")
    interpreter = ["bash"]
    environment = merge(local.ssh_env, {
      REGISTRY_HOST        = var.registry_host
      REGISTRY_PORT        = tostring(var.registry_port)
      REGISTRY_SSH_HOST    = local.registry_ssh_host
      CA_FILE              = var.registry_ca_file
      MASTER_IP            = local.first_cp.management_ip
      RKE2_HOSTS           = join(" ", local.rke2_node_mgmt_ips)
      CONTROL_HOSTS        = join(" ", [for n in var.control_plane_nodes : n.management_ip])
      WORKER_HOSTS         = join(" ", [for n in var.worker_nodes : n.management_ip])
      INSECURE_SKIP_VERIFY = tostring(var.registry_insecure_skip_verify)
      RESTART_RKE2         = tostring(var.registry_restart_rke2)
      MOUNTPOINT           = var.registry_mountpoint
    })
  }
}

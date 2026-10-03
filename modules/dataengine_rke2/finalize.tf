resource "null_resource" "finalize" {
  depends_on = [
    null_resource.first_server,
    null_resource.additional_servers,
    null_resource.agents,
  ]

  triggers = {
    fingerprint = local.install_fingerprint
    label       = tostring(var.label_workers)
    artifacts   = local.kubeconfig_artifacts_dir
    workers     = join(",", [for n in var.worker_nodes : n.name])
  }

  provisioner "local-exec" {
    command     = abspath("${path.module}/scripts/finalize_cluster.sh")
    interpreter = ["bash"]
    environment = merge(local.ssh_env, {
      READY_TIMEOUT = tostring(var.ready_timeout_seconds)
      EXPECTED      = tostring(local.expected_node_count)
      LABEL_WORKERS = tostring(var.label_workers)
      WORKER_NAMES  = join(" ", [for n in var.worker_nodes : n.name])
      KUBE_OUT      = local.kubeconfig_abs
      FIRST_HOST    = local.first_cp.management_ip
      FIRST_USER    = var.ssh_user
    })
  }
}

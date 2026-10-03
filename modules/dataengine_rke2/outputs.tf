output "first_control_plane" {
  description = "First control-plane node (manager)."
  value       = local.first_cp
}

output "join_address" {
  description = "Host used for RKE2 registration (server: https://HOST:9345)."
  value       = local.join_address
}

output "control_plane_nodes" {
  value = var.control_plane_nodes
}

output "worker_nodes" {
  value = var.worker_nodes
}

output "rke2_channel" {
  value = var.rke2_channel
}

output "kubeconfig_path" {
  description = "Local kubeconfig path ({artifacts_dir}/kubeconfig)."
  value       = local.kubeconfig_abs
}

output "kubernetes_config_path" {
  description = <<-EOT
    Value for the root module's kubernetes/helm provider config_path.
    Equals kubeconfig_path once finalize has written it; null before that, so
    the providers stay unconfigured and a greenfield apply can bootstrap RKE2.
    Depends only on inputs, so it is safe to reference from provider blocks.
  EOT
  value       = fileexists(local.kubeconfig_abs) ? local.kubeconfig_abs : null
}

output "kubeconfig_artifacts_dir" {
  description = "Directory with kubeconfig plus VMS paste files (ca_cert, client_cert, client_key)."
  value       = local.kubeconfig_artifacts_abs
}

output "kubectl_hint" {
  description = "Example kubectl invocation using the written kubeconfig."
  value       = "KUBECONFIG=${local.kubeconfig_abs} kubectl get nodes"
}

output "podman_registry_enabled" {
  value = local.registry_enabled
}

output "registry_endpoint" {
  description = "Registry host:port when enable_podman_registry=true; otherwise null."
  value       = local.registry_enabled ? "${var.registry_host}:${var.registry_port}" : null
}

output "registry_ca_file" {
  description = "CA path installed on RKE2 nodes when registry is enabled."
  value       = local.registry_enabled ? var.registry_ca_file : null
}

output "local_path_storage_enabled" {
  value = local.local_path_storage_enabled
}

output "vast_csi_enabled" {
  value = local.vast_csi_enabled
}

output "storage_class_name" {
  description = "Preferred StorageClass name (default SC when configured; otherwise first enabled backend; null if none)."
  value       = local.storage_class_name
}

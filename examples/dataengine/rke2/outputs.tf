output "join_address" {
  value = module.rke2.join_address
}

output "kubeconfig_path" {
  value = module.rke2.kubeconfig_path
}

output "kubeconfig_artifacts_dir" {
  value = module.rke2.kubeconfig_artifacts_dir
}

output "kubectl_hint" {
  value = module.rke2.kubectl_hint
}

output "first_control_plane" {
  value = module.rke2.first_control_plane
}

output "registry_endpoint" {
  value = module.rke2.registry_endpoint
}

output "storage_class_name" {
  value = module.rke2.storage_class_name
}

terraform {
  required_version = ">= 1.5.0"

  required_providers {
    null = {
      source  = "hashicorp/null"
      version = ">= 3.2.0"
    }
    local = {
      source  = "hashicorp/local"
      version = ">= 2.4.0"
    }
    external = {
      source  = "hashicorp/external"
      version = ">= 2.3.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.30.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 3.0.0"
    }
  }
}

module "rke2" {
  source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/dataengine_rke2?ref=85ba8e8f53fafcae17b4fafa7422348e06069c57"

  ssh_user             = var.ssh_user
  ssh_private_key_path = var.ssh_private_key_path
  ssh_private_key      = var.ssh_private_key
  ssh_port             = var.ssh_port

  rke2_channel   = var.rke2_channel
  cluster_token  = var.cluster_token
  join_address   = var.join_address
  extra_tls_sans = var.extra_tls_sans

  control_plane_nodes = var.control_plane_nodes
  worker_nodes        = var.worker_nodes

  prepare_nodes = var.prepare_nodes
  apt_upgrade   = var.apt_upgrade
  ping_vip      = var.ping_vip

  kubeconfig_artifacts_dir = var.kubeconfig_artifacts_dir
  label_workers            = var.label_workers
  ready_timeout_seconds    = var.ready_timeout_seconds

  enable_podman_registry        = var.enable_podman_registry
  registry_host                 = var.registry_host
  registry_management_ip        = var.registry_management_ip
  registry_port                 = var.registry_port
  registry_nfs_export           = var.registry_nfs_export
  registry_mountpoint           = var.registry_mountpoint
  registry_cert_days            = var.registry_cert_days
  registry_ca_file              = var.registry_ca_file
  registry_insecure_skip_verify = var.registry_insecure_skip_verify
  registry_force_recert         = var.registry_force_recert
  registry_run_smoke_test       = var.registry_run_smoke_test
  registry_restart_rke2         = var.registry_restart_rke2

  enable_local_path_storage = var.enable_local_path_storage
  local_path_version        = var.local_path_version
  local_path_set_default    = var.local_path_set_default

  enable_vast_csi              = var.enable_vast_csi
  vast_csi_endpoint            = var.vast_csi_endpoint
  vast_csi_username            = var.vast_csi_username
  vast_csi_password            = var.vast_csi_password
  vast_csi_tenant              = var.vast_csi_tenant
  vast_csi_vip_pool            = var.vast_csi_vip_pool
  vast_csi_view_policy         = var.vast_csi_view_policy
  vast_csi_storage_path        = var.vast_csi_storage_path
  vast_csi_namespace           = var.vast_csi_namespace
  vast_csi_secret_name         = var.vast_csi_secret_name
  vast_csi_release_name        = var.vast_csi_release_name
  vast_csi_chart_version       = var.vast_csi_chart_version
  vast_csi_storage_class_name  = var.vast_csi_storage_class_name
  vast_csi_set_default         = var.vast_csi_set_default
  vast_csi_snapshotter_version = var.vast_csi_snapshotter_version
  vast_csi_install_snapshotter = var.vast_csi_install_snapshotter
}

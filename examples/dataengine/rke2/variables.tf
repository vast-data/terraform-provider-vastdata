variable "ssh_user" {
  type    = string
  default = "vastdata"
}

variable "ssh_private_key_path" {
  type     = string
  default  = "~/.ssh/id_ed25519"
  nullable = true
}

variable "ssh_private_key" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "ssh_port" {
  type    = number
  default = 22
}

variable "rke2_channel" {
  type    = string
  default = "v1.35"
}

variable "cluster_token" {
  type      = string
  default   = "RKE2CLUSTER"
  sensitive = true
}

variable "join_address" {
  type        = string
  description = "Defaults to first control-plane storage_ip when null."
  default     = null
  nullable    = true
}

variable "extra_tls_sans" {
  type    = list(string)
  default = []
}

variable "control_plane_nodes" {
  type = list(object({
    name          = string
    management_ip = string
    storage_ip    = string
  }))
}

variable "worker_nodes" {
  type = list(object({
    name          = string
    management_ip = string
    storage_ip    = string
  }))
  default = []
}

variable "prepare_nodes" {
  type    = bool
  default = true
}

variable "apt_upgrade" {
  type    = bool
  default = false
}

variable "ping_vip" {
  type        = string
  description = "Optional CSI/VIP address to ping from all nodes before install."
  default     = null
  nullable    = true
}

variable "kubeconfig_artifacts_dir" {
  type    = string
  default = "out"
}

variable "label_workers" {
  type    = bool
  default = true
}

variable "ready_timeout_seconds" {
  type    = number
  default = 900
}

variable "enable_podman_registry" {
  type    = bool
  default = false
}

variable "registry_host" {
  type     = string
  default  = null
  nullable = true
}

variable "registry_management_ip" {
  type     = string
  default  = null
  nullable = true
}

variable "registry_port" {
  type    = number
  default = 5000
}

variable "registry_nfs_export" {
  type     = string
  default  = null
  nullable = true
}

variable "registry_mountpoint" {
  type    = string
  default = "/var/lib/dataengine-registry"
}

variable "registry_cert_days" {
  type    = number
  default = 3650
}

variable "registry_ca_file" {
  type    = string
  default = "/etc/rancher/rke2/certs/registry-ca.crt"
}

variable "registry_insecure_skip_verify" {
  type    = bool
  default = false
}

variable "registry_force_recert" {
  type    = bool
  default = false
}

variable "registry_run_smoke_test" {
  type    = bool
  default = false
}

variable "registry_restart_rke2" {
  type    = bool
  default = true
}

variable "enable_local_path_storage" {
  type     = bool
  default  = null
  nullable = true
}

variable "local_path_version" {
  type    = string
  default = "v0.0.31"
}

variable "local_path_set_default" {
  type    = bool
  default = true
}

variable "enable_vast_csi" {
  type    = bool
  default = false
}

variable "vast_csi_endpoint" {
  type     = string
  default  = null
  nullable = true
}

variable "vast_csi_username" {
  type     = string
  default  = null
  nullable = true
}

variable "vast_csi_password" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "vast_csi_tenant" {
  type     = string
  default  = null
  nullable = true
}

variable "vast_csi_vip_pool" {
  type     = string
  default  = null
  nullable = true
}

variable "vast_csi_view_policy" {
  type     = string
  default  = null
  nullable = true
}

variable "vast_csi_storage_path" {
  type    = string
  default = "/k8s/csi/pvcs"
}

variable "vast_csi_namespace" {
  type    = string
  default = "vast-csi"
}

variable "vast_csi_secret_name" {
  type    = string
  default = "de-lab-csi"
}

variable "vast_csi_release_name" {
  type    = string
  default = "csi-driver"
}

variable "vast_csi_chart_version" {
  type     = string
  default  = null
  nullable = true
}

variable "vast_csi_storage_class_name" {
  type    = string
  default = "vast-nfs"
}

variable "vast_csi_set_default" {
  type    = bool
  default = true
}

variable "vast_csi_snapshotter_version" {
  type    = string
  default = "v7.0.1"
}

variable "vast_csi_install_snapshotter" {
  type    = bool
  default = true
}

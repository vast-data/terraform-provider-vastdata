# -----------------------------------------------------------------------------
# SSH
# -----------------------------------------------------------------------------

variable "ssh_user" {
  type        = string
  description = "SSH user on all cluster nodes (must have passwordless sudo)."
  default     = "vastdata"
}

variable "ssh_private_key_path" {
  type        = string
  description = "Path to the SSH private key used to reach cluster nodes."
  default     = null
  nullable    = true
}

variable "ssh_private_key" {
  type        = string
  description = "SSH private key contents (alternative to ssh_private_key_path)."
  default     = null
  sensitive   = true
  nullable    = true
}

variable "ssh_port" {
  type        = number
  description = "SSH port on cluster nodes."
  default     = 22
}

variable "ssh_strict_host_key_checking" {
  type        = string
  description = "OpenSSH StrictHostKeyChecking value (accept-new is typical for labs)."
  default     = "accept-new"
}

# -----------------------------------------------------------------------------
# Nodes
# -----------------------------------------------------------------------------

variable "control_plane_nodes" {
  type = list(object({
    name          = string
    management_ip = string
    storage_ip    = string
  }))
  description = <<-EOT
    Control-plane / etcd nodes. Index 0 is the first server (manager).
    Terraform SSHes to management_ip; cluster join uses storage_ip of the first node.
  EOT

  validation {
    condition     = length(var.control_plane_nodes) >= 1
    error_message = "At least one control-plane node is required."
  }
}

variable "worker_nodes" {
  type = list(object({
    name          = string
    management_ip = string
    storage_ip    = string
  }))
  description = "Worker (agent) nodes. Join uses the first control-plane storage_ip."
  default     = []
}

# -----------------------------------------------------------------------------
# RKE2
# -----------------------------------------------------------------------------

variable "rke2_channel" {
  type        = string
  description = "INSTALL_RKE2_CHANNEL passed to get.rke2.io (e.g. v1.35, stable, latest)."
  default     = "v1.35"
}

variable "cluster_token" {
  type        = string
  description = "Shared RKE2 cluster token (config token:)."
  default     = "RKE2CLUSTER"
  sensitive   = true
}

variable "join_address" {
  type        = string
  description = <<-EOT
    API/registration address for additional servers and agents (host only, no scheme/port).
    Defaults to the first control-plane node's storage_ip (backend/data network).
  EOT
  default     = null
  nullable    = true
}

variable "extra_tls_sans" {
  type        = list(string)
  description = "Extra TLS SANs for control-plane nodes (in addition to each node's mgmt/storage IPs)."
  default     = []
}

variable "write_kubeconfig_mode" {
  type        = string
  description = "write-kubeconfig-mode for server nodes."
  default     = "0644"
}

# -----------------------------------------------------------------------------
# Node prep (no ClusterShell; SSH per node)
# -----------------------------------------------------------------------------

variable "prepare_nodes" {
  type        = bool
  description = "If true, apt-get update and install nfs-common on all nodes before RKE2."
  default     = true
}

variable "apt_upgrade" {
  type        = bool
  description = "If true (and prepare_nodes), run apt-get upgrade -y on all nodes. Slow; off by default."
  default     = false
}

variable "ping_vip" {
  type        = string
  description = "Optional VIP to ping from every node before install (e.g. CSI VIP pool address). Null skips."
  default     = null
  nullable    = true
}

# -----------------------------------------------------------------------------
# Local outputs / kubectl helpers
# -----------------------------------------------------------------------------

variable "kubeconfig_artifacts_dir" {
  type        = string
  description = <<-EOT
    Directory for local kubeconfig and VMS paste files. Finalize writes:
      kubeconfig, ca_cert, client_cert, client_key
    Default: $${path.root}/out
  EOT
  default     = null
  nullable    = true
}

variable "label_workers" {
  type        = bool
  description = "If true, set node-role.kubernetes.io/worker=\"\" on each worker node."
  default     = true
}

variable "ready_timeout_seconds" {
  type        = number
  description = "How long to wait for all expected nodes to report Ready."
  default     = 900
}

# -----------------------------------------------------------------------------
# Podman registry (optional; DataEngine “Create a Podman Registry” guide)
# -----------------------------------------------------------------------------

variable "enable_podman_registry" {
  type        = bool
  description = "If true, install a TLS Podman registry and configure RKE2 nodes to trust it."
  default     = false
}

variable "registry_host" {
  type        = string
  description = "Registry IP used in cert SAN, registries.yaml, and image names (no DNS)."
  default     = null
  nullable    = true
}

variable "registry_management_ip" {
  type        = string
  description = "SSH target for the registry VM. Defaults to registry_host."
  default     = null
  nullable    = true
}

variable "registry_port" {
  type        = number
  description = "Host port published to the registry container (maps to container 5000)."
  default     = 5000
}

variable "registry_nfs_export" {
  type        = string
  description = <<-EOT
    Optional NFS export for registry data (VAST lab, e.g. 172.200.205.19:/de-lab-docker/images).
    When null, the module uses a local directory on the registry VM (created automatically).
  EOT
  default     = null
  nullable    = true
}

variable "registry_mountpoint" {
  type        = string
  description = <<-EOT
    Directory on the registry VM for image data and certs.
    When registry_nfs_export is set, this is the NFS mount path; otherwise a local directory.
  EOT
  default     = "/var/lib/dataengine-registry"
}

variable "registry_cert_days" {
  type        = number
  description = "Self-signed certificate validity in days."
  default     = 3650
}

variable "registry_ca_file" {
  type        = string
  description = "Path where the registry CA is installed on every RKE2 node."
  default     = "/etc/rancher/rke2/certs/registry-ca.crt"
}

variable "registry_insecure_skip_verify" {
  type        = bool
  description = <<-EOT
    registries.yaml tls.insecure_skip_verify. Prefer false when ca_file is installed.
    The published lab paste sets true; troubleshooting says not to — default is false.
  EOT
  default     = false
}

variable "registry_force_recert" {
  type        = bool
  description = "If true, regenerate domain.crt/key even when they already exist on the NFS mount."
  default     = false
}

variable "registry_run_smoke_test" {
  type        = bool
  description = "If true, push alpine into the new registry as a check. Off by default; setup already curls /v2/_catalog."
  default     = false
}

variable "registry_restart_rke2" {
  type        = bool
  description = "If true, restart rke2-server/agent after writing registries.yaml. RKE2 reads that file only at process start."
  default     = true
}

# -----------------------------------------------------------------------------
# Storage: local-path (lab default) and/or VAST CSI
# -----------------------------------------------------------------------------

variable "enable_local_path_storage" {
  type = bool
  description = <<-EOT
    Install Rancher local-path provisioner and StorageClass local-path after the
    cluster is Ready. Default null means auto: install when enable_vast_csi is
    false (lab default); skip when CSI alone is enabled. Set true to
    keep local-path alongside CSI, or false to opt out entirely.
  EOT
  default  = null
  nullable = true
}

variable "local_path_version" {
  type        = string
  description = "Rancher local-path-provisioner release tag for the upstream deploy YAML."
  default     = "v0.0.31"
}

variable "local_path_set_default" {
  type        = bool
  description = <<-EOT
    Annotate local-path as the default StorageClass when it is installed.
    Ignored (treated as false) when VAST CSI is also installed with
    vast_csi_set_default=true so only one default SC is configured.
  EOT
  default = true
}

variable "enable_vast_csi" {
  type        = bool
  description = "If true, install VAST CSI Driver (Helm) and create StorageClass after the cluster is Ready."
  default     = false
}

variable "vast_csi_endpoint" {
  type        = string
  description = "VMS endpoint IP/hostname for the CSI driver secret. Required when enable_vast_csi=true."
  default     = null
  nullable    = true
}

variable "vast_csi_username" {
  type        = string
  description = "CSI driver service-account username. Required when enable_vast_csi=true."
  default     = null
  nullable    = true
}

variable "vast_csi_password" {
  type        = string
  description = "CSI driver service-account password. Required when enable_vast_csi=true."
  default     = null
  sensitive   = true
  nullable    = true
}

variable "vast_csi_tenant" {
  type        = string
  description = "VAST tenant name for CSI. Required when enable_vast_csi=true."
  default     = null
  nullable    = true
}

variable "vast_csi_vip_pool" {
  type        = string
  description = "VIP pool name used by the CSI StorageClass. Required when enable_vast_csi=true."
  default     = null
  nullable    = true
}

variable "vast_csi_view_policy" {
  type        = string
  description = "View policy name used by the CSI StorageClass. Required when enable_vast_csi=true."
  default     = null
  nullable    = true
}

variable "vast_csi_storage_path" {
  type        = string
  description = "Root path on VAST for dynamically provisioned volumes."
  default     = "/k8s/csi/pvcs"
}

variable "vast_csi_namespace" {
  type        = string
  description = "Kubernetes namespace for the VAST CSI Driver release."
  default     = "vast-csi"
}

variable "vast_csi_secret_name" {
  type        = string
  description = "Name of the Opaque secret holding VMS CSI credentials."
  default     = "de-lab-csi"
}

variable "vast_csi_release_name" {
  type        = string
  description = "Helm release name for vastcsi."
  default     = "csi-driver"
}

variable "vast_csi_chart_version" {
  type        = string
  description = "Optional pin of the vastcsi chart version. Null uses the latest from the repo."
  default     = null
  nullable    = true
}

variable "vast_csi_storage_class_name" {
  type        = string
  description = "StorageClass name created by the VAST CSI chart (storageClasses map key)."
  default     = "vast-nfs"
}

variable "vast_csi_set_default" {
  type        = bool
  description = "If true, set the VAST CSI StorageClass as the cluster default."
  default     = true
}

variable "vast_csi_snapshotter_version" {
  type        = string
  description = "kubernetes-csi/external-snapshotter tag for CRDs and snapshot-controller manifests."
  default     = "v7.0.1"
}

variable "vast_csi_install_snapshotter" {
  type        = bool
  description = "If true, apply external-snapshotter CRDs/RBAC/controller before the CSI Helm release."
  default     = true
}

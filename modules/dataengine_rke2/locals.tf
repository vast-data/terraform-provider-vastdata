locals {
  ssh_private_key_path = var.ssh_private_key_path != null ? pathexpand(var.ssh_private_key_path) : null

  ssh_private_key = coalesce(
    var.ssh_private_key,
    try(file(local.ssh_private_key_path), null)
  )

  join_address = coalesce(var.join_address, var.control_plane_nodes[0].storage_ip)

  first_cp = var.control_plane_nodes[0]

  additional_control_planes = slice(var.control_plane_nodes, 1, length(var.control_plane_nodes))

  all_nodes = concat(var.control_plane_nodes, var.worker_nodes)

  all_nodes_map = {
    for n in local.all_nodes : n.name => n
  }

  expected_node_count = length(local.all_nodes)

  kubeconfig_artifacts_dir = coalesce(
    var.kubeconfig_artifacts_dir,
    "${path.root}/out"
  )
  kubeconfig_local_path = "${local.kubeconfig_artifacts_dir}/kubeconfig"

  ssh_common = {
    user                     = var.ssh_user
    port                     = var.ssh_port
    strict_host_key_checking = var.ssh_strict_host_key_checking
  }

  # Fingerprint of inputs that should force re-bootstrap when changed.
  install_fingerprint = sha256(jsonencode({
    channel   = var.rke2_channel
    token     = var.cluster_token
    join      = local.join_address
    cps       = var.control_plane_nodes
    workers   = var.worker_nodes
    tls_extra = var.extra_tls_sans
  }))

  registry_enabled = var.enable_podman_registry

  registry_ssh_host = try(coalesce(var.registry_management_ip, var.registry_host), null)

  registry_fingerprint = sha256(jsonencode({
    host        = var.registry_host
    ssh         = local.registry_ssh_host
    port        = var.registry_port
    nfs         = var.registry_nfs_export
    mount       = var.registry_mountpoint
    days        = var.registry_cert_days
    ca_file     = var.registry_ca_file
    skip_verify = var.registry_insecure_skip_verify
    force_cert  = var.registry_force_recert
  }))

  rke2_node_mgmt_ips = [for n in local.all_nodes : n.management_ip]

  kubeconfig_abs           = abspath(local.kubeconfig_local_path)
  kubeconfig_artifacts_abs = abspath(local.kubeconfig_artifacts_dir)

  # Auto: local-path when CSI is off; explicit true/false overrides.
  local_path_storage_enabled = coalesce(var.enable_local_path_storage, !var.enable_vast_csi)
  vast_csi_enabled           = var.enable_vast_csi

  # Only one default SC: CSI wins when both ask to be default.
  local_path_is_default = (
    local.local_path_storage_enabled &&
    var.local_path_set_default &&
    !(local.vast_csi_enabled && var.vast_csi_set_default)
  )
  vast_csi_is_default = local.vast_csi_enabled && var.vast_csi_set_default

  # Preferred default SC name for consumers (zarf, docs). Null if neither is default.
  storage_class_name = (
    local.vast_csi_is_default ? var.vast_csi_storage_class_name :
    local.local_path_is_default ? "local-path" :
    local.vast_csi_enabled ? var.vast_csi_storage_class_name :
    local.local_path_storage_enabled ? "local-path" :
    null
  )
}

check "ssh_key_provided" {
  assert {
    condition     = var.ssh_private_key_path != null || var.ssh_private_key != null
    error_message = "Set ssh_private_key_path or ssh_private_key."
  }
}

check "registry_inputs" {
  assert {
    condition     = !var.enable_podman_registry || var.registry_host != null
    error_message = "enable_podman_registry requires registry_host."
  }
}

check "vast_csi_kubeconfig_present" {
  assert {
    condition     = !var.enable_vast_csi || fileexists(local.kubeconfig_abs)
    error_message = "enable_vast_csi needs ${local.kubeconfig_abs} to exist: kubernetes/helm providers read it at plan time. Apply once with enable_vast_csi=false to bootstrap RKE2, then enable it."
  }
}

check "vast_csi_inputs" {
  assert {
    condition = !var.enable_vast_csi || (
      var.vast_csi_endpoint != null && var.vast_csi_endpoint != "" &&
      var.vast_csi_username != null && var.vast_csi_username != "" &&
      var.vast_csi_password != null && var.vast_csi_password != "" &&
      var.vast_csi_tenant != null && var.vast_csi_tenant != "" &&
      var.vast_csi_vip_pool != null && var.vast_csi_vip_pool != "" &&
      var.vast_csi_view_policy != null && var.vast_csi_view_policy != ""
    )
    error_message = "enable_vast_csi requires vast_csi_endpoint, username, password, tenant, vip_pool, and view_policy."
  }
}

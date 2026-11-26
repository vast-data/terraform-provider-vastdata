# Modules

Terraform modules in this directory. Pin a commit or tag when you consume them:

```hcl
source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/<name>?ref=<commit-or-tag>"
```

| Module | What it does |
| --- | --- |
| [dataengine_vms_setup](#dataengine_vms_setup) | VMS objects for a DataEngine tenant (identity, VIP pool, broker, topics, optional CSI) |
| [dataengine_rke2](#dataengine_rke2) | RKE2 on existing Ubuntu VMs, optional Podman registry and StorageClass |

Pipeline workload resources (`vastdata_de_function`, `vastdata_de_schedule_trigger`, `vastdata_de_pipeline`) are **not** part of these modules. See [`examples/dataengine/hello_world`](../examples/dataengine/hello_world) after the image is built and pushed to the registry.

**Important** in the tables means you must set the value for a successful apply (or must set it when the named feature is on). Other inputs have defaults.

---

## dataengine_vms_setup

`modules/dataengine_vms_setup`

Creates the VMS side of a DataEngine lab: local (or existing LDAP/AD) identity, tenant, app group, broker owner, VIP pool, S3 view policy, event-broker view, main and dead-letter topics, and an optional API token. With `enable_data_engine` it enables Data Engine; with `enable_kubernetes_cluster` it registers mTLS credentials and an external Kubernetes cluster (e.g. from `dataengine_rke2` `out/` certs); with `enable_container_registry` it registers a container registry bound to that cluster. With `enable_csi` it also creates the CSI VIP pool, NFS view policy, docker-repo view, roles, and manager users.

Resources are split by concern: `identity.tf`, `tenant.tf`, `vip_pool.tf`, `broker.tf`, `topics.tf`, `api_token.tf`, `de_engine.tf`, `kubernetes_cluster.tf`, `container_registry.tf`, `tenant_admin.tf`, `csi.tf` (plus `variables.tf` / `locals.tf` / `outputs.tf` / `versions.tf`).

The root module must configure the `vastdata` provider. This module does not open a VMS connection itself.

### Inputs

| Parameter | Important | Default | Description |
| --- | --- | --- | --- |
| `name_prefix` | no | `de-lab` | Prefix used when a name override is null |
| `tenant_name` | no | `name_prefix` | DataEngine tenant name |
| `identity_provider_type` | no | `local` | `local`, `ldap`, or `active_directory` |
| `create_local_provider` | no | `true` | Create a local provider, or reuse `local_provider_id` |
| `local_provider_name` | no | `<prefix>-provider` | Name of the created local provider |
| `local_provider_id` | when reusing local | `1` | Existing local provider ID when `create_local_provider=false` |
| `ldap_provider_id` | when `ldap` | null | Existing LDAP provider ID |
| `ad_provider_id` | when `active_directory` | null | Existing AD provider ID |
| `app_group_name` | when LDAP/AD | `<prefix>-app-group` | Created name (local) or existing group name (LDAP/AD) |
| `app_group_gid` | yes, when `local` | null | POSIX GID for the created app group |
| `broker_owner_name` | when LDAP/AD | `<prefix>-broker-owner` | Created name (local) or existing username (LDAP/AD) |
| `broker_owner_uid` | yes, when `local` | null | POSIX UID for the created broker owner |
| `broker_owner_password` | no | null | Password for the created local broker owner |
| `vip_pool_name` | no | `<prefix>-vippool` | DataEngine VIP pool name |
| `vip_pool` | yes | none | `subnet_cidr`, `ip_ranges`, optional `gw_ip`, `vlan`, `domain_name` |
| `view_policy_name` | no | `<prefix>-s3-policy` | S3_NATIVE view policy name |
| `s3_view_policy_auth_source` | no | `RPC_AND_PROVIDERS` | Group membership source for the S3 policy |
| `broker_bucket` | no | `<prefix>-broker` | S3 bucket / Kafka database name |
| `broker_view_path` | no | `/engine-event-broker` | Filesystem path of the broker view |
| `delete_dir` | no | `true` | Delete view directories on destroy (needs Trash API) |
| `force_delete_tenant` | no | `false` | Force-delete tenants on destroy |
| `enable_data_engine` | no | `false` | Create `vastdata_de_engine` (POST `/data-engine`) after topics exist; passes `tenant_name` on API queries |
| `enable_kubernetes_cluster` | no | `false` | After DE enable: create mTLS credentials + register an external K8s cluster |
| `kubernetes_kube_api_url` | when enabling K8s | null | API server URL (e.g. `https://<cp-ip>:6443` from RKE2) |
| `kubernetes_kubeconfig_path` | preferred mTLS source | null | Path to kubeconfig (e.g. RKE2 `out/kubeconfig`); provider derives certs |
| `kubernetes_kubeconfig` | alt to path | null | Inline kubeconfig YAML (sensitive; write-only on provider) |
| `kubernetes_kubeconfig_context` | no | null | Kubeconfig context (default: `current-context`) |
| `kubernetes_mtls_artifacts_dir` | alt mTLS source | null | Dir with `ca_cert` / `client_cert` / `client_key` (RKE2 `out/`) |
| `kubernetes_certificate_authority_b64` | alt to artifacts dir | null | Base64 CA (overrides file) |
| `kubernetes_client_certificate_b64` | alt to artifacts dir | null | Base64 client cert (overrides file) |
| `kubernetes_client_key_b64` | alt to artifacts dir | null | Base64 client key (overrides file) |
| `kubernetes_cluster_name` | no | `<prefix>-k8s` | Data Engine Kubernetes cluster name |
| `kubernetes_mtls_name` | no | `<prefix>-mtls` | mTLS credential resource name |
| `kubernetes_namespaces` | no | default, vast-dataengine, vast-pipelines | Namespaces registered on the cluster |
| `kubernetes_is_internal` | no | `false` | Internal-to-VAST cluster flag |
| `enable_container_registry` | no | `false` | Register a container registry bound to the K8s cluster |
| `container_registry_url` | when enabling registry | null | Registry URL (`https://host:5000` or ngrok HTTPS) |
| `container_registry_name` | no | `<prefix>-registry` | Data Engine container registry name |
| `container_registry_auth_type` | no | `none` | `none`, `password`, or `secret` |
| `container_registry_primary_namespace` | no | `default` | Primary namespace for registry binding |
| `container_registry_is_accessible_to_data_engine` | no | `true` | Periodic registry reachability checks |
| `create_tenant_admin` | no | `false` | Create a TENANT_ADMIN manager on the DataEngine tenant |
| `tenant_admin_username` | no | `<prefix>-admin` | Tenant-admin manager username |
| `tenant_admin_password` | when `create_tenant_admin` | null | Password for that manager |
| `create_tenant_admin_role` | no | `true` | `true` = create role; `false` = look up existing role by name |
| `tenant_admin_role_name` | when using existing role | `<prefix>-tenant-admin` | Role to create, or existing role name to assign |
| `tenant_admin_permissions_list` | no | `[]` | Fine-grained codenames when creating the role; empty uses permission types |
| `tenant_admin_extra_permissions_list` | no | `[]` | Extra codenames merged into the admin list (create-role mode) |
| `tenant_admin_permission_types` | no | create, view, edit, delete | Used when creating the role and the permissions list is empty |
| `main_topic_name` | no | `engine-main` | Default DataEngine topic |
| `deadletter_topic_name` | no | `engine-deadletter` | Dead-letter topic |
| `topic_partitions` | no | `3` | Kafka partition count |
| `topic_retention_ms` | no | `21600000` | Topic retention (6h) |
| `kafka_unencrypted_auth_mechanism` | no | `NONE` | Unencrypted Kafka auth mechanism |
| `kafka_is_authorization_required` | no | `false` | Require Kafka authorization on the broker view |
| `create_api_token` | no | `true` | Create an API token for later `vastde` use |
| `api_token_name` | no | `<prefix>-vastde-admin` | API token name |
| `api_token_owner` | no | null | Token owner; null uses the requesting user |
| `api_token_expiry` | no | `1Y` | Token expiry (`1Y`, `30D`, …) |
| `enable_csi` | no | `false` | Also create CSI tenant helpers |
| `create_csi_tenant` | no | `false` | Dedicated CSI tenant instead of the DataEngine tenant |
| `csi_tenant_name` | no | `<prefix>-csi-tenant` | CSI tenant name |
| `trash_gid` | no | `0` | Trash folder GID on the CSI tenant |
| `csi_vip_pool_name` | no | `<prefix>-csi-vippool` | CSI VIP pool name |
| `csi_vip_pool` | when `enable_csi` | null | Same shape as `vip_pool` |
| `csi_view_policy_name` | no | `<prefix>-csi-policy` | NFS CSI view policy name |
| `csi_view_policy` | no | RPC, bind VIP, no-squash `*` | `auth_source`, `bind_vip_pool`, `nfs_no_squash`, `nfs_root_squash` |
| `docker_repo_path` | no | `/docker/images` | Path of the NFS docker-repo view |
| `docker_repo_view_name` | no | `<prefix>-docker-repo` | Docker repo view name |
| `csi_tenant_admin_role_name` | no | `<prefix>-csi-tenant-admin` | Tenant-admin role name |
| `csi_service_account_role_name` | no | `<prefix>-csi-service-account` | CSI driver role name |
| `csi_tenant_admin_username` | no | `<prefix>-csi-admin` | Tenant-admin manager username |
| `csi_driver_username` | no | `<prefix>-csi-driver` | CSI driver manager username |
| `csi_manager_password` | when `enable_csi` | null | Password for both CSI manager accounts |
| `csi_tenant_admin_permissions_list` | no | `[]` | Fine-grained codenames; empty uses permission types |
| `csi_tenant_admin_extra_permissions_list` | no | `[]` | Extra codenames merged into the admin list |
| `csi_tenant_admin_permission_types` | no | create, view, edit, delete | Used when the permissions list is empty |
| `csi_service_account_permissions_list` | no | `[]` | Fine-grained codenames for the driver role |
| `csi_service_account_extra_permissions_list` | no | `[]` | Extra codenames merged into the driver list |
| `csi_service_account_permission_types` | no | create, view, edit, delete | Used when the driver permissions list is empty |

### Examples

Basic: local identity, one VIP pool.

```hcl
provider "vastdata" {
  host            = "10.143.11.201"
  username        = "admin"
  password        = "..."
  skip_ssl_verify = true
}

module "dataengine" {
  source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/dataengine_vms_setup?ref=<commit-or-tag>"

  app_group_gid    = 854321
  broker_owner_uid = 123456

  vip_pool = {
    subnet_cidr = 16
    ip_ranges   = [["172.200.205.13", "172.200.205.18"]]
  }
}
```

Existing LDAP provider (the provider must already exist, with VMS authentication enabled).

```hcl
module "dataengine" {
  source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/dataengine_vms_setup?ref=<commit-or-tag>"

  identity_provider_type = "ldap"
  ldap_provider_id       = 2
  app_group_name         = "de-app-group"
  broker_owner_name      = "de-broker-owner"

  vip_pool = {
    subnet_cidr = 16
    ip_ranges   = [["172.200.205.13", "172.200.205.18"]]
  }
}
```

CSI helpers on a dedicated tenant (VIP pool must not overlap the DataEngine pool).

```hcl
module "dataengine" {
  source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/dataengine_vms_setup?ref=<commit-or-tag>"

  app_group_gid    = 854321
  broker_owner_uid = 123456

  vip_pool = {
    subnet_cidr = 16
    ip_ranges   = [["172.200.205.13", "172.200.205.18"]]
  }

  enable_csi           = true
  create_csi_tenant    = true
  csi_manager_password = "..."
  csi_vip_pool = {
    subnet_cidr = 16
    ip_ranges   = [["172.200.205.19", "172.200.205.24"]]
  }
}
```

A complete root module (provider, every input, `terraform.tfvars.example`) is in [examples/dataengine/vms_setup](../examples/dataengine/vms_setup).

---

## dataengine_rke2

`modules/dataengine_rke2`

Installs RKE2 over SSH on existing Ubuntu VMs (no ClusterShell, no VM provisioning). Optional: TLS Podman registry trusted by every node, Rancher local-path StorageClass, and the VAST CSI Driver.

The root module must declare `kubernetes` and `helm` and point them at `kubernetes_config_path`. That output is null until finalize writes the kubeconfig, so the first apply can bootstrap the cluster. Enable `enable_vast_csi` only on a later apply, after the kubeconfig file exists.

`terraform destroy` does not uninstall RKE2 or the registry container.

### Inputs

| Parameter | Important | Default | Description |
| --- | --- | --- | --- |
| `ssh_user` | no | `vastdata` | SSH user; needs passwordless sudo |
| `ssh_private_key_path` | yes (or key) | null | Path to the SSH private key |
| `ssh_private_key` | yes (or path) | null | Key contents, instead of a path |
| `ssh_port` | no | `22` | SSH port |
| `ssh_strict_host_key_checking` | no | `accept-new` | OpenSSH `StrictHostKeyChecking` |
| `control_plane_nodes` | yes | none | List of `{name, management_ip, storage_ip}`; index 0 is the first server |
| `worker_nodes` | no | `[]` | Agent nodes, same object shape |
| `rke2_channel` | no | `v1.35` | `INSTALL_RKE2_CHANNEL` (`v1.35`, `stable`, `latest`) |
| `cluster_token` | no | `RKE2CLUSTER` | Shared RKE2 token |
| `join_address` | no | first CP `storage_ip` | Registration host for extra servers and agents |
| `extra_tls_sans` | no | `[]` | Extra API TLS SANs on control-plane nodes |
| `write_kubeconfig_mode` | no | `0644` | `write-kubeconfig-mode` on servers |
| `prepare_nodes` | no | `true` | `apt-get update` and install `nfs-common` before RKE2 |
| `apt_upgrade` | no | `false` | Also `apt-get upgrade` during prep |
| `ping_vip` | no | null | Address to ping from every node before install |
| `kubeconfig_artifacts_dir` | no | `<path.root>/out` | Where finalize writes kubeconfig and VMS paste files |
| `label_workers` | no | `true` | Set `node-role.kubernetes.io/worker` on workers |
| `ready_timeout_seconds` | no | `900` | Wait for every node to be Ready |
| `enable_podman_registry` | no | `false` | Install a TLS Podman registry and trust it on RKE2 nodes |
| `registry_host` | when registry on | null | Registry IP (cert SAN and image names; no DNS) |
| `registry_management_ip` | no | `registry_host` | SSH target for the registry VM |
| `registry_port` | no | `5000` | Host port for the registry container |
| `registry_nfs_export` | no | null | NFS export for registry data; null uses a local directory |
| `registry_mountpoint` | no | `/var/lib/dataengine-registry` | Data and cert directory on the registry VM |
| `registry_cert_days` | no | `3650` | Self-signed cert lifetime |
| `registry_ca_file` | no | `/etc/rancher/rke2/certs/registry-ca.crt` | CA path installed on every RKE2 node |
| `registry_insecure_skip_verify` | no | `false` | `registries.yaml` `tls.insecure_skip_verify` |
| `registry_force_recert` | no | `false` | Regenerate the registry cert even if one exists |
| `registry_run_smoke_test` | no | `false` | Optional check that pushes `alpine` into the new registry. Setup already curls `/v2/_catalog` |
| `registry_restart_rke2` | no | `true` | Restart RKE2 after writing `registries.yaml`. The file is read only when the process starts |
| `enable_local_path_storage` | no | `null` (auto) | Install local-path. Auto: on when CSI is off, off when CSI alone is on. Set `true`/`false` to override |
| `local_path_version` | no | `v0.0.31` | local-path-provisioner release tag |
| `local_path_set_default` | no | `true` | Make `local-path` the default SC (ignored when CSI is also default) |
| `enable_vast_csi` | no | `false` | Install the VAST CSI Driver (needs kubeconfig already on disk) |
| `vast_csi_endpoint` | when CSI on | null | VMS address for the CSI secret |
| `vast_csi_username` | when CSI on | null | CSI driver service-account username |
| `vast_csi_password` | when CSI on | null | CSI driver service-account password |
| `vast_csi_tenant` | when CSI on | null | VAST tenant name for CSI |
| `vast_csi_vip_pool` | when CSI on | null | VIP pool name on the StorageClass |
| `vast_csi_view_policy` | when CSI on | null | View policy name on the StorageClass |
| `vast_csi_storage_path` | no | `/k8s/csi/pvcs` | Root path for dynamic volumes |
| `vast_csi_namespace` | no | `vast-csi` | Namespace for the Helm release |
| `vast_csi_secret_name` | no | `de-lab-csi` | Secret holding VMS CSI credentials |
| `vast_csi_release_name` | no | `csi-driver` | Helm release name |
| `vast_csi_chart_version` | no | null | Chart version; null uses the repo latest |
| `vast_csi_storage_class_name` | no | `vast-nfs` | StorageClass name |
| `vast_csi_set_default` | no | `true` | Make the VAST StorageClass the default |
| `vast_csi_snapshotter_version` | no | `v7.0.1` | external-snapshotter tag |
| `vast_csi_install_snapshotter` | no | `true` | Install snapshot CRDs and controller before Helm |

By default (neither flag set) the module installs **local-path** and marks it default.

### Examples

The provider blocks are the same for every example. `kubernetes_config_path` is null on the first apply and the real kubeconfig path after finalize.

```hcl
provider "kubernetes" {
  config_path = module.rke2.kubernetes_config_path
}

provider "helm" {
  kubernetes = {
    config_path = module.rke2.kubernetes_config_path
  }
}
```

Basic: one control plane and two workers.

```hcl
module "rke2" {
  source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/dataengine_rke2?ref=<commit-or-tag>"

  ssh_private_key_path = "~/.ssh/id_ed25519"

  control_plane_nodes = [
    { name = "de-cp1", management_ip = "10.178.218.17", storage_ip = "10.178.218.17" },
  ]
  worker_nodes = [
    { name = "de-w1", management_ip = "10.178.218.226", storage_ip = "10.178.218.226" },
    { name = "de-w2", management_ip = "10.178.218.238", storage_ip = "10.178.218.238" },
  ]
}
```

Lab without VAST CSI: local-path StorageClass and a TLS Podman registry on its own VM.

```hcl
module "rke2" {
  source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/dataengine_rke2?ref=<commit-or-tag>"

  ssh_private_key_path = "~/.ssh/id_ed25519"

  control_plane_nodes = [
    { name = "de-cp1", management_ip = "10.178.218.17", storage_ip = "10.178.218.17" },
  ]
  worker_nodes = [
    { name = "de-w1", management_ip = "10.178.218.226", storage_ip = "10.178.218.226" },
    { name = "de-w2", management_ip = "10.178.218.238", storage_ip = "10.178.218.238" },
  ]


  enable_podman_registry = true
  registry_host          = "10.178.218.164"
}
```

VAST CSI after `dataengine_vms_setup` with `enable_csi=true`. Apply once without CSI so `out/kubeconfig` exists, then apply again with `enable_vast_csi=true`. Names below match that module's defaults for `name_prefix=de-lab`.

```hcl
module "rke2" {
  source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/dataengine_rke2?ref=<commit-or-tag>"

  ssh_private_key_path = "~/.ssh/id_ed25519"
  kubeconfig_artifacts_dir = "out"

  control_plane_nodes = [
    { name = "de-cp1", management_ip = "10.178.218.17", storage_ip = "10.178.218.17" },
  ]
  worker_nodes = [
    { name = "de-w1", management_ip = "10.178.218.226", storage_ip = "10.178.218.226" },
    { name = "de-w2", management_ip = "10.178.218.238", storage_ip = "10.178.218.238" },
  ]

  enable_vast_csi             = true
  vast_csi_endpoint           = "10.143.11.201"
  vast_csi_username           = "de-lab-csi-driver"
  vast_csi_password           = "..."
  vast_csi_tenant             = "de-lab-csi-tenant"
  vast_csi_vip_pool           = "de-lab-csi-vippool"
  vast_csi_view_policy        = "de-lab-csi-policy"
  local_path_set_default      = false
}
```

A complete root module (providers, every input, `terraform.tfvars.example`) is in [examples/dataengine/rke2](../examples/dataengine/rke2).

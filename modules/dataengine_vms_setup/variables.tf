# -----------------------------------------------------------------------------
# Naming
# -----------------------------------------------------------------------------

variable "name_prefix" {
  type        = string
  description = "Default prefix for resource names when individual name overrides are null."
  default     = "de-lab"
}

variable "tenant_name" {
  type        = string
  description = "DataEngine tenant name. Defaults to name_prefix."
  default     = null
  nullable    = true
}

# -----------------------------------------------------------------------------
# Identity provider
# Default: local (create dedicated provider). Alternatives: existing LDAP or AD.
# -----------------------------------------------------------------------------

variable "identity_provider_type" {
  type        = string
  description = <<-EOT
    How tenants get their identity provider:
      local              — vastdata_local_provider (default). Creates or reuses a local provider.
      ldap               — attach an existing LDAP provider (ldap_provider_id). Does not create LDAP.
      active_directory   — attach an existing AD provider (ad_provider_id). Does not create AD.
    For ldap/active_directory, enable "VMS Authentication Provider" on that provider first
    (GUI or separate terraform for vastdata_ldap / vastdata_active_directory).
  EOT
  default     = "local"

  validation {
    condition     = contains(["local", "ldap", "active_directory"], var.identity_provider_type)
    error_message = "identity_provider_type must be local, ldap, or active_directory."
  }
}

variable "create_local_provider" {
  type        = bool
  description = "When identity_provider_type=local: true creates a dedicated local provider; false reuses local_provider_id."
  default     = true
}

variable "local_provider_name" {
  type        = string
  description = "Name for the created local provider when create_local_provider=true. Defaults to \"$${name_prefix}-provider\"."
  default     = null
  nullable    = true
}

variable "local_provider_id" {
  type        = number
  description = "Existing local provider ID when create_local_provider=false (often 1)."
  default     = 1
}

variable "ldap_provider_id" {
  type        = number
  description = "Existing LDAP provider ID when identity_provider_type=ldap."
  default     = null
  nullable    = true
}

variable "ad_provider_id" {
  type        = number
  description = "Existing Active Directory provider ID when identity_provider_type=active_directory."
  default     = null
  nullable    = true
}

# -----------------------------------------------------------------------------
# App group + broker owner user
# Local mode creates them. LDAP/AD mode references existing names (required).
# -----------------------------------------------------------------------------

variable "app_group_name" {
  type        = string
  description = "App group name. Local: created (default $${name_prefix}-app-group). LDAP/AD: existing group name (required)."
  default     = null
  nullable    = true
}

variable "app_group_gid" {
  type        = number
  description = "POSIX GID when creating a local app group (identity_provider_type=local)."
  default     = null
  nullable    = true
}

variable "broker_owner_name" {
  type        = string
  description = "Broker owner username. Local: created (default $${name_prefix}-broker-owner). LDAP/AD: existing username (required)."
  default     = null
  nullable    = true
}

variable "broker_owner_uid" {
  type        = number
  description = "POSIX UID when creating a local broker owner (identity_provider_type=local)."
  default     = null
  nullable    = true
}

variable "broker_owner_password" {
  type        = string
  description = "Optional password for a created local broker owner."
  default     = null
  sensitive   = true
  nullable    = true
}

# -----------------------------------------------------------------------------
# DataEngine VIP pool
# -----------------------------------------------------------------------------

variable "vip_pool_name" {
  type        = string
  description = "DataEngine / broker VIP pool name. Defaults to \"$${name_prefix}-vippool\"."
  default     = null
  nullable    = true
}

variable "vip_pool" {
  type = object({
    subnet_cidr = number
    ip_ranges   = list(list(string))
    gw_ip       = optional(string)
    vlan        = optional(number)
    domain_name = optional(string)
  })
  description = "DataEngine VIP pool network settings (required). Optional gw_ip, vlan, domain_name."
}

# -----------------------------------------------------------------------------
# Broker view + S3 policy + topics
# -----------------------------------------------------------------------------

variable "view_policy_name" {
  type        = string
  description = "S3_NATIVE view policy name. Defaults to \"$${name_prefix}-s3-policy\"."
  default     = null
  nullable    = true
}

variable "s3_view_policy_auth_source" {
  type        = string
  description = "Group membership source for the S3 policy (GUI: NFS RPC and Providers → RPC_AND_PROVIDERS)."
  default     = "RPC_AND_PROVIDERS"
}

variable "broker_bucket" {
  type        = string
  description = "S3 bucket / Kafka database name. Defaults to \"$${name_prefix}-broker\"."
  default     = null
  nullable    = true
}

variable "broker_view_path" {
  type        = string
  description = "Filesystem path for the event broker view."
  default     = "/engine-event-broker"
}

variable "delete_dir" {
  type        = bool
  description = "If true, delete the underlying directory when destroying any vastdata_view. Requires Trash API."
  default     = true
}

variable "force_delete_tenant" {
  type        = bool
  description = "Force-delete tenants on destroy. Prefer delete_dir; enable only if tenant delete still fails."
  default     = false
}

variable "enable_data_engine" {
  type        = bool
  description = <<-EOT
    When true, create vastdata_de_engine (POST /data-engine) after topics exist.
    Auth may be tenant-admin (provider tenant) or cluster admin with tenant_name
    passed on Data Engine API queries.
  EOT
  default     = false
}

# -----------------------------------------------------------------------------
# Optional external Kubernetes cluster (after Data Engine enable)
# -----------------------------------------------------------------------------

variable "enable_kubernetes_cluster" {
  type        = bool
  description = <<-EOT
    When true, create mTLS credentials and register an external Kubernetes cluster
    with Data Engine. Requires enable_data_engine=true, kubernetes_kube_api_url, and
    mTLS material via one of:
      • kubernetes_kubeconfig / kubernetes_kubeconfig_path (preferred), or
      • kubernetes_mtls_artifacts_dir (ca_cert/client_cert/client_key), or
      • the three kubernetes_*_b64 certificate inputs.
  EOT
  default     = false
}

variable "kubernetes_cluster_name" {
  type        = string
  description = "Data Engine Kubernetes cluster resource name. Defaults to \"$${name_prefix}-k8s\"."
  default     = null
  nullable    = true
}

variable "kubernetes_kube_api_url" {
  type        = string
  description = "Kubernetes API server URL (e.g. https://<first-cp>:6443 from RKE2 finalize)."
  default     = null
  nullable    = true
}

variable "kubernetes_kubeconfig" {
  type        = string
  description = <<-EOT
    Raw kubeconfig YAML for vastdata_de_mtls_authentication_credential.kubeconfig_data.
    Alternative to PEM artifacts / *_b64. Write-only on the provider resource (not stored
    in state). Mutually exclusive with kubernetes_mtls_artifacts_dir and *_b64 inputs.
  EOT
  default     = null
  sensitive   = true
  nullable    = true
}

variable "kubernetes_kubeconfig_path" {
  type        = string
  description = <<-EOT
    Path to a kubeconfig file (e.g. examples/dataengine/rke2/out/kubeconfig). Read with
    file() into kubernetes_kubeconfig when that var is null. Preferred over separate
    ca_cert/client_cert/client_key files in kubernetes_mtls_artifacts_dir.
  EOT
  default     = null
  nullable    = true
}

variable "kubernetes_kubeconfig_context" {
  type        = string
  description = "Kubeconfig context name. Defaults to current-context when null."
  default     = null
  nullable    = true
}

variable "kubernetes_mtls_artifacts_dir" {
  type        = string
  description = <<-EOT
    Directory containing ca_cert, client_cert, and client_key as written by
    modules/dataengine_rke2 finalize (e.g. examples/dataengine/rke2/out). Values are
    already base64 (kubeconfig *-data fields). Ignored when kubeconfig or the three
    *_b64 vars are set.
  EOT
  default     = null
  nullable    = true
}

variable "kubernetes_certificate_authority_b64" {
  type        = string
  description = "Base64 CA cert (overrides file from kubernetes_mtls_artifacts_dir)."
  default     = null
  sensitive   = true
  nullable    = true
}

variable "kubernetes_client_certificate_b64" {
  type        = string
  description = "Base64 client cert (overrides file from kubernetes_mtls_artifacts_dir)."
  default     = null
  sensitive   = true
  nullable    = true
}

variable "kubernetes_client_key_b64" {
  type        = string
  description = "Base64 client key (overrides file from kubernetes_mtls_artifacts_dir)."
  default     = null
  sensitive   = true
  nullable    = true
}

variable "kubernetes_mtls_name" {
  type        = string
  description = "Name for vastdata_de_mtls_authentication_credential. Defaults to \"$${name_prefix}-mtls\"."
  default     = null
  nullable    = true
}

variable "kubernetes_namespaces" {
  type        = list(string)
  description = "Namespaces registered on the Data Engine Kubernetes cluster resource."
  default = [
    "default",
    "vast-dataengine",
    "vast-pipelines",
  ]
}

variable "kubernetes_is_internal" {
  type        = bool
  description = "Whether the Kubernetes cluster is internal to VAST."
  default     = false
}

# -----------------------------------------------------------------------------
# Optional container registry (after Kubernetes cluster register)
# -----------------------------------------------------------------------------

variable "enable_container_registry" {
  type        = bool
  description = <<-EOT
    When true, create vastdata_de_container_registry bound to the registered Kubernetes cluster.
    Requires enable_kubernetes_cluster=true and container_registry_url.
  EOT
  default     = false
}

variable "container_registry_name" {
  type        = string
  description = "Data Engine container registry name. Defaults to \"$${name_prefix}-registry\"."
  default     = null
  nullable    = true
}

variable "container_registry_url" {
  type        = string
  description = "Registry URL (e.g. https://10.178.218.55:5000 or an ngrok HTTPS URL)."
  default     = null
  nullable    = true
}

variable "container_registry_auth_type" {
  type        = string
  description = "Registry auth: none, password, or secret."
  default     = "none"
}

variable "container_registry_username" {
  type        = string
  description = "Registry username (e.g. Docker Hub user). Used when auth_type=password."
  default     = null
  nullable    = true
}

variable "container_registry_password" {
  type        = string
  description = "Registry password or PAT. Used when auth_type=password."
  default     = null
  nullable    = true
  sensitive   = true
}

variable "container_registry_primary_namespace" {
  type        = string
  description = "Primary Kubernetes namespace for the registry binding."
  default     = "default"
}

variable "container_registry_is_accessible_to_data_engine" {
  type        = bool
  description = "When true, Data Engine periodically validates registry reachability."
  default     = true
}

# -----------------------------------------------------------------------------
# Optional DataEngine tenant-admin manager
# -----------------------------------------------------------------------------

variable "create_tenant_admin" {
  type        = bool
  description = "Create a TENANT_ADMIN manager on the DataEngine tenant."
  default     = false
}

variable "tenant_admin_username" {
  type        = string
  description = "Tenant-admin manager username. Defaults to \"$${name_prefix}-admin\"."
  default     = null
  nullable    = true
}

variable "tenant_admin_password" {
  type        = string
  description = "Password for the DataEngine tenant-admin manager. Required when create_tenant_admin=true."
  default     = null
  sensitive   = true
  nullable    = true
}

variable "create_tenant_admin_role" {
  type        = bool
  description = <<-EOT
    When create_tenant_admin=true: if true (default), create a new administrator role and assign it;
    if false, look up an existing role by tenant_admin_role_name and assign that.
  EOT
  default     = true
}

variable "tenant_admin_role_name" {
  type        = string
  description = <<-EOT
    Role name. When create_tenant_admin_role=true, defaults to \"$${name_prefix}-tenant-admin\".
    When create_tenant_admin_role=false, must be set to an existing role name.
  EOT
  default     = null
  nullable    = true
}

variable "tenant_admin_permissions_list" {
  type        = set(string)
  description = "Fine-grained permission codenames when creating the role. Empty → use permission types. Ignored when create_tenant_admin_role=false."
  default     = []
}

variable "tenant_admin_extra_permissions_list" {
  type        = set(string)
  description = "Additional codenames merged into tenant_admin_permissions_list (create-role mode only)."
  default     = []
}

variable "tenant_admin_permission_types" {
  type        = set(string)
  description = "Legacy permission types when creating the role and permissions_list is empty."
  default     = ["create", "view", "edit", "delete"]
}

variable "main_topic_name" {
  type        = string
  description = "Default DataEngine topic name."
  default     = "engine-main"
}

variable "deadletter_topic_name" {
  type        = string
  description = "Dead-letter topic name."
  default     = "engine-deadletter"
}

variable "topic_partitions" {
  type        = number
  description = "Kafka topic partition count."
  default     = 3
}

variable "topic_retention_ms" {
  type        = number
  description = "Kafka topic retention in milliseconds (default 6h)."
  default     = 21600000
}

variable "kafka_unencrypted_auth_mechanism" {
  type        = string
  description = "Kafka unencrypted auth mechanism (lab often NONE)."
  default     = "NONE"
}

variable "kafka_is_authorization_required" {
  type        = bool
  description = "Whether Kafka authorization is required on the broker view."
  default     = false
}

# -----------------------------------------------------------------------------
# API token (for vastde later)
# -----------------------------------------------------------------------------

variable "create_api_token" {
  type        = bool
  description = "Create a vastdata_api_token for later vastde use."
  default     = true
}

variable "api_token_name" {
  type        = string
  description = "API token name. Defaults to \"$${name_prefix}-vastde-admin\"."
  default     = null
  nullable    = true
}

variable "api_token_owner" {
  type        = string
  description = "Token owner username. Null = requesting user."
  default     = null
  nullable    = true
}

variable "api_token_expiry" {
  type        = string
  description = "Token expiry (e.g. 1Y, 30D)."
  default     = "1Y"
}

# -----------------------------------------------------------------------------
# Optional CSI helpers
# -----------------------------------------------------------------------------

variable "enable_csi" {
  type        = bool
  description = "Create CSI-oriented tenant helpers (NFS policy/view, roles, managers, optional CSI VIP/tenant)."
  default     = false
}

variable "create_csi_tenant" {
  type        = bool
  description = "When enable_csi=true, create a dedicated CSI tenant instead of attaching CSI objects to the DataEngine tenant."
  default     = false
}

variable "csi_tenant_name" {
  type        = string
  description = "CSI tenant name when create_csi_tenant=true. Defaults to \"$${name_prefix}-csi-tenant\"."
  default     = null
  nullable    = true
}

variable "trash_gid" {
  type        = number
  description = "Trash folder GID on the CSI (or DE) tenant when enable_csi=true."
  default     = 0
}

variable "csi_vip_pool_name" {
  type        = string
  description = "CSI VIP pool name. Defaults to \"$${name_prefix}-csi-vippool\"."
  default     = null
  nullable    = true
}

variable "csi_vip_pool" {
  type = object({
    subnet_cidr = number
    ip_ranges   = list(list(string))
    gw_ip       = optional(string)
    vlan        = optional(number)
    domain_name = optional(string)
  })
  description = "CSI VIP pool settings. Required when enable_csi=true."
  default     = null
  nullable    = true
}

variable "csi_view_policy_name" {
  type        = string
  description = "NFS CSI view policy name. Defaults to \"$${name_prefix}-csi-policy\"."
  default     = null
  nullable    = true
}

variable "csi_view_policy" {
  type = object({
    # Guide: NFS RPC → "RPC". Alternatives: RPC_AND_PROVIDERS, PROVIDERS, etc.
    auth_source = optional(string, "RPC")
    # Bind the CSI VIP pool to this policy (guide selects VIP on the policy).
    bind_vip_pool   = optional(bool, true)
    nfs_no_squash   = optional(set(string), ["*"])
    nfs_root_squash = optional(set(string), [])
  })
  description = "CSI NFS view policy options. Defaults match current lab behavior (RPC, bind VIP, no-squash *)."
  default     = {}
}

variable "docker_repo_path" {
  type        = string
  description = "Path for the NFS docker-repo view."
  default     = "/docker/images"
}

variable "docker_repo_view_name" {
  type        = string
  description = "Docker repo view name. Defaults to \"$${name_prefix}-docker-repo\"."
  default     = null
  nullable    = true
}

variable "csi_tenant_admin_role_name" {
  type        = string
  description = "CSI tenant-admin role name. Defaults to \"$${name_prefix}-csi-tenant-admin\"."
  default     = null
  nullable    = true
}

variable "csi_service_account_role_name" {
  type        = string
  description = "CSI service-account role name. Defaults to \"$${name_prefix}-csi-service-account\"."
  default     = null
  nullable    = true
}

variable "csi_tenant_admin_username" {
  type        = string
  description = "CSI tenant-admin manager username. Defaults to \"$${name_prefix}-csi-admin\"."
  default     = null
  nullable    = true
}

variable "csi_driver_username" {
  type        = string
  description = "CSI driver manager username. Defaults to \"$${name_prefix}-csi-driver\"."
  default     = null
  nullable    = true
}

variable "csi_manager_password" {
  type        = string
  description = "Password for CSI tenant-admin and driver manager accounts. Required when enable_csi=true."
  default     = null
  sensitive   = true
  nullable    = true
}

# Fine-grained permissions (codenames from /permissions/get). When the base
# list is empty, fall back to permission *types* (create/view/edit/delete) —
# same as today. Use extra_* to extend without replacing the base list.

variable "csi_tenant_admin_permissions_list" {
  type        = set(string)
  description = "Fine-grained permission codenames for CSI tenant-admin role. Empty → use permission types."
  default     = []
}

variable "csi_tenant_admin_extra_permissions_list" {
  type        = set(string)
  description = "Additional codenames merged into csi_tenant_admin_permissions_list."
  default     = []
}

variable "csi_tenant_admin_permission_types" {
  type        = set(string)
  description = "Legacy permission types when permissions_list is empty (GUI Logical Full Control)."
  default     = ["create", "view", "edit", "delete"]
}

variable "csi_service_account_permissions_list" {
  type        = set(string)
  description = "Fine-grained permission codenames for CSI service-account role. Empty → use permission types."
  default     = []
}

variable "csi_service_account_extra_permissions_list" {
  type        = set(string)
  description = "Additional codenames merged into csi_service_account_permissions_list."
  default     = []
}

variable "csi_service_account_permission_types" {
  type        = set(string)
  description = "Legacy permission types when permissions_list is empty."
  default     = ["create", "view", "edit", "delete"]
}

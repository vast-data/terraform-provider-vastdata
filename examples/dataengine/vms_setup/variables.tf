variable "vms_host" {
  type        = string
  description = "VMS address (hostname or IP, no scheme/port)."
}

variable "vms_port" {
  type    = number
  default = 443
}

variable "vms_username" {
  type     = string
  default  = null
  nullable = true
}

variable "vms_password" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "vms_api_token" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "skip_ssl_verify" {
  type    = bool
  default = true
}

variable "name_prefix" {
  type    = string
  default = "de-lab"
}

variable "identity_provider_type" {
  type    = string
  default = "local"
}

variable "create_local_provider" {
  type    = bool
  default = true
}

variable "local_provider_name" {
  type     = string
  default  = null
  nullable = true
}

variable "local_provider_id" {
  type    = number
  default = 1
}

variable "ldap_provider_id" {
  type     = number
  default  = null
  nullable = true
}

variable "ad_provider_id" {
  type     = number
  default  = null
  nullable = true
}

variable "app_group_name" {
  type     = string
  default  = null
  nullable = true
}

variable "app_group_gid" {
  type     = number
  default  = null
  nullable = true
}

variable "broker_owner_name" {
  type     = string
  default  = null
  nullable = true
}

variable "broker_owner_uid" {
  type     = number
  default  = null
  nullable = true
}

variable "broker_owner_password" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "vip_pool" {
  type = object({
    subnet_cidr = number
    ip_ranges   = list(list(string))
    gw_ip       = optional(string)
    vlan        = optional(number)
    domain_name = optional(string)
  })
}

variable "broker_view_path" {
  type    = string
  default = "/engine-event-broker"
}

variable "s3_view_policy_auth_source" {
  type    = string
  default = "RPC_AND_PROVIDERS"
}

variable "delete_dir" {
  type    = bool
  default = true
}

variable "force_delete_tenant" {
  type    = bool
  default = false
}

variable "enable_data_engine" {
  type    = bool
  default = false
}

variable "enable_kubernetes_cluster" {
  type    = bool
  default = false
}

variable "kubernetes_cluster_name" {
  type     = string
  default  = null
  nullable = true
}

variable "kubernetes_kube_api_url" {
  type     = string
  default  = null
  nullable = true
}

variable "kubernetes_kubeconfig" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "kubernetes_kubeconfig_path" {
  type     = string
  default  = null
  nullable = true
}

variable "kubernetes_kubeconfig_context" {
  type     = string
  default  = null
  nullable = true
}

variable "kubernetes_mtls_artifacts_dir" {
  type     = string
  default  = null
  nullable = true
}

variable "kubernetes_certificate_authority_b64" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "kubernetes_client_certificate_b64" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "kubernetes_client_key_b64" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "kubernetes_mtls_name" {
  type     = string
  default  = null
  nullable = true
}

variable "kubernetes_namespaces" {
  type = list(string)
  default = [
    "default",
    "vast-dataengine",
    "vast-pipelines",
  ]
}

variable "kubernetes_is_internal" {
  type    = bool
  default = false
}

variable "enable_container_registry" {
  type    = bool
  default = false
}

variable "container_registry_name" {
  type     = string
  default  = null
  nullable = true
}

variable "container_registry_url" {
  type     = string
  default  = null
  nullable = true
}

variable "container_registry_auth_type" {
  type    = string
  default = "none"
}

variable "container_registry_username" {
  type     = string
  default  = null
  nullable = true
}

variable "container_registry_password" {
  type      = string
  default   = null
  nullable  = true
  sensitive = true
}

variable "container_registry_primary_namespace" {
  type    = string
  default = "default"
}

variable "container_registry_is_accessible_to_data_engine" {
  type    = bool
  default = true
}

variable "create_tenant_admin" {
  type    = bool
  default = false
}

variable "tenant_admin_username" {
  type     = string
  default  = null
  nullable = true
}

variable "tenant_admin_password" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "create_tenant_admin_role" {
  type    = bool
  default = true
}

variable "tenant_admin_role_name" {
  type     = string
  default  = null
  nullable = true
}

variable "tenant_admin_permissions_list" {
  type    = set(string)
  default = []
}

variable "tenant_admin_extra_permissions_list" {
  type    = set(string)
  default = []
}

variable "tenant_admin_permission_types" {
  type    = set(string)
  default = ["create", "view", "edit", "delete"]
}

variable "main_topic_name" {
  type    = string
  default = "engine-main"
}

variable "deadletter_topic_name" {
  type    = string
  default = "engine-deadletter"
}

variable "topic_partitions" {
  type    = number
  default = 3
}

variable "create_api_token" {
  type    = bool
  default = true
}

variable "api_token_name" {
  type        = string
  description = "API token name. Defaults to \"$${name_prefix}-vastde-admin\" in the module."
  default     = null
  nullable    = true
}

variable "api_token_owner" {
  type     = string
  default  = null
  nullable = true
}

variable "api_token_expiry" {
  type    = string
  default = "1Y"
}

variable "enable_csi" {
  type    = bool
  default = false
}

variable "create_csi_tenant" {
  type    = bool
  default = false
}

variable "csi_vip_pool" {
  type = object({
    subnet_cidr = number
    ip_ranges   = list(list(string))
    gw_ip       = optional(string)
    vlan        = optional(number)
    domain_name = optional(string)
  })
  default  = null
  nullable = true
}

variable "csi_view_policy" {
  type = object({
    auth_source     = optional(string, "RPC")
    bind_vip_pool   = optional(bool, true)
    nfs_no_squash   = optional(set(string), ["*"])
    nfs_root_squash = optional(set(string), [])
  })
  default = {}
}

variable "docker_repo_path" {
  type    = string
  default = "/docker/images"
}

variable "csi_manager_password" {
  type      = string
  default   = null
  sensitive = true
  nullable  = true
}

variable "trash_gid" {
  type    = number
  default = 0
}

variable "csi_tenant_admin_permissions_list" {
  type    = set(string)
  default = []
}

variable "csi_tenant_admin_extra_permissions_list" {
  type    = set(string)
  default = []
}

variable "csi_service_account_permissions_list" {
  type    = set(string)
  default = []
}

variable "csi_service_account_extra_permissions_list" {
  type    = set(string)
  default = []
}

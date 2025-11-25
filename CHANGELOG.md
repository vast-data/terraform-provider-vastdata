
## 3.0.2

ENHANCEMENTS:

* **vastdata_view_policy**: fixed drift due to VMS ips optimization (ORION-288457)

## 3.0.1

ENHANCEMENTS:

* **Documentation**: Added missing breaking changes to v3.0.0 release notes.
* **Schema Descriptions**: Updated OpenAPI schema processing to restore missing property descriptions.


## 3.0.0

BREAKING CHANGES:

* **vastdata_user**: `local_provider_id` is now a **required** field (was optional in v2.x). All user resources must specify a local provider.
* **vastdata_group**: `local_provider_id` is now a **required** field (was optional in v2.x). All group resources must specify a local provider.
* **vastdata_view_policy**: `vip_pools` attribute (Set of Number) has been replaced with `permission_per_vip_pool` (Map of String). Users must migrate their configurations to use the new attribute format.
* **vastdata_view_policy**: `protocols_audit` attribute is now **read-only** (was optional/configurable in v2.x). This attribute is now computed from the cluster settings and cannot be set directly on the view policy. Use the `protocols` attribute instead for specifying which protocols to audit.
* **vastdata_protected_path**: `tenant_id` is now a **required** field (was optional in v2.x). All protected path resources must specify a tenant ID.
* **vastdata_protected_path**: `target_exported_dir` is now a **required** field (was optional in v2.x). All protected path resources must specify a target export directory.

NEW RESOURCES:

* **vastdata_cluster**: Manage VAST cluster configuration
* **vastdata_cluster_ekm**: Manage cluster external key management
* **vastdata_cnode**: Manage VAST cluster nodes
* **vastdata_cnode_bgp_config**: Manage BGP configuration for cluster nodes
* **vastdata_iam_role**: Manage IAM roles
* **vastdata_kerberos**: Manage Kerberos configuration
* **vastdata_kerberos_keytab**: Manage Kerberos keytab files
* **vastdata_manager_password**: Manage VMS manager passwords
* **vastdata_oidc**: Manage OIDC authentication configuration
* **vastdata_rack**: Manage VAST rack configuration
* **vastdata_rack_bgp_config**: Manage BGP configuration for racks
* **vastdata_supported_drives**: Manage supported drive configurations
* **vastdata_tenant_nfs4_delegation**: Manage NFS4 delegation settings for tenants
* **vastdata_topic**: Manage Kafka topics

ENHANCEMENTS:

* **vastdata_nonlocal_user**: Added support for searching users by `sid` field
* **vastdata_nonlocal_user_key**: Added support for searching users by `sid` field; Improved documentation to correctly show usage with non-local users
* **vastdata_s3_policy_attachment**: Added support for searching users/groups by `sid` field; Improved documentation and examples to correctly reference non-local users/groups instead of local users
* **vastdata_active_directory**: Enhanced Active Directory integration capabilities
* **vastdata_block_host_mapping**: Improved block host mapping configuration
* **Provider HTTP Client**: Added support for `HTTPS_PROXY`, `HTTP_PROXY`, and `NO_PROXY` environment variables for proxy configuration
* **Import Operations**: Enhanced import logic to prevent drift by only populating computed and required fields (optional fields remain null after import)


## 2.1.1

ENHANCEMENTS:

* **vastdata_s3_policy_attachment**: Added support for `username` and `groupname` attributes as alternatives to `uid` and `gid` for identifying users and groups

## 2.1.0

ENHANCEMENTS:

* **vastdata_vip_pool**: Added client monitoring support for VIP pool connectivity monitoring
* **vastdata_volume**: Added volume monitoring capabilities
* **vastdata_vms**: Added management data interface configuration for L3 networks
* **vastdata_s3_policy_attachment**: get s3_policy_attachment by `s3_policy_guid`; added import support

## 2.0.1

BUG FIXES:

* **vastdata_nonlocal_user**: Fixed `uid` and `vid` attributes being incorrectly marked as non computed for nonlocal_user resource
* 
## 2.0.0

NOTES:

* The provider has been **fully rewritten** using the new [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework), replacing the legacy SDKv2 implementation.

* The provider now supports **dynamic schema generation** by parsing the VAST OpenAPI specification at runtime. This enables easier exposure of VAST API endpoints as Terraform resources and data sources, significantly improving coverage and maintainability.

BREAKING CHANGES:

* Legacy resources defined using SDKv2 have been removed or migrated to the new Plugin Framework format. You may need to review their configurations for compatibility, especially if relying on custom behaviors or non-standard attributes.
* We've prepared a migration script that automatically converts your existing configuration files into the format supported by the new Terraform provider. You can find the usage instructions [here](migration/README.md)

ENHANCEMENTS:

* Improved maintainability and extensibility via the Plugin Framework.
* Consistent handling of resource plans, state, and diagnostics.
* Enhanced API errors handling and validation using [go-vast-client](https://github.com/vast-data/go-vast-client) library.
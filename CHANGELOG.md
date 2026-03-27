## 3.2

BUG FIXES:

* **vastdata_topic**: Fixed compatibility issue with older VAST versions where the Topics API returns 404 instead of 400 for missing topics. (TERF-224)

NEW RESOURCES:

* **vastdata_webhook**: Manage VAST cluster configuration


## 3.1.1

BUG FIXES:

* **vastdata_s3_policy**: Fixed `ImportState` and `Read` failing with "Duplicate Set Element" error when the cluster API returns empty-string principals. (TERF-223)

## 3.1.0

BUG FIXES:

* **vastdata_cnode_bgp_config**: Fixed PATCH request failing with 400 "This field may not be blank" when BGP address fields (`port1_self_address`, `port1_peer_address`, `port2_self_address`, `port2_peer_address`, `self_asn`) are empty strings. (TERF-218)
* **vastdata_protected_path**: Fixed update failing with 400 "Can't modify exported path for protected path" when a linked `protection_policy` is modified. The API returns `target_exported_dir: null` after creation, causing Terraform to detect drift and attempt to PATCH the immutable field. (TERF-220)
* **vastdata_protected_path**: Fixed `remote_tenant_name` being read-only, preventing edits on the Protected Path resource. (TERF-213)
* **vastdata_qos_policy**: Fixed perpetual drift after creating a `VIEW`-type QoS policy. The API returns `is_default: null` instead of the submitted `false`, causing Terraform to detect a change on every plan. (TERF-221)
* **vastdata_view**: Fixed Terraform not detecting `share_acl` changes made outside of Terraform (e.g., via the GUI). When `share_acl` permissions were modified through the VAST management UI, `terraform plan` incorrectly reported "No changes" instead of detecting the drift and reverting to the declared configuration. (TERF-210)
* **vastdata_administrator_manager**: Fixed false drift detection on the `password` field. `terraform plan` reported changes on the `password` attribute even when no modifications were made, because the sensitive field was being compared against the API response. (TERF-214)
* **vastdata_view**: Fixed validation error when importing a BLOCK view whose `alias` does not start with `/`. The "must start with '/'" validator was incorrectly applied to the `alias` field, which does not require a leading slash for BLOCK protocol views. (TERF-211)

ENHANCEMENTS:

* Added `UseStateForUnknown` plan modifier for all computed fields to reduce noise during field updates.
* Added support for `X-Tenant-Name` header in login flow (TERF-86).
* Embedded migration mode in the Terraform provider via the `VASTDATA_MIGRATE_MODE` environment variable. See the [How Migrate Mode Works](migration/README.md).
* **vastdata_kafkabroker**: Added `guid` field.
* **vastdata_ldap**: Added `netgroup_searchbase` field.
* **vastdata_s3_lifecycle_rule**: Added `tags` field.
* **vastdata_tenant**: Added `list_open_handles_task` field.
* **vastdata_view_policy**: Added `smb_recursive_change_notify` field.

KNOWN ISSUES:

* **vastdata_view**: `share_acl.acl[].fqdn` may reset from `"All"` to `""` (empty string) after view creation for the bucket owner's ACL entry (TERF-222). The VAST API returns a modified value for `fqdn` on subsequent GET requests, causing `terraform plan` to show perpetual drift.

## 3.0.7

BUG FIXES:

* **vastdata_quota**: Fixed `default_user_quota` and `default_group_quota` fields incorrectly marked as read-only.
* **vastdata_view_policy**: Fixed `protocols_audit` field incorrectly marked as read-only.
*  Fixed drift detection not reverting changes made through VAST UI. Resources now properly detect and revert manual changes made outside of Terraform, ensuring infrastructure state consistency.

ENHANCEMENTS:

* **vastdata_s3_policy_attachment**: Added `group_sid` field to support groups that may not have valid POSIX GIDs (e.g., `gid = -1`).

## 3.0.6

ENHANCEMENTS:

* Use bearer token auth instead of basic auth

## 3.0.5

BUG FIXES:

* **vastdata_nonlocal_user**: Fixed `access_keys` field schema generation. The field was missing from the resource schema because PATCH request body schema was used instead of PATCH response schema. Now correctly uses the PATCH response schema which includes computed fields like `access_keys`.
*  Fixed GUID fallback mechanism not updating `id` in state.

## 3.0.4

BUG FIXES:

* Populate required fields during import
* **vastdata_protected_path**: Fixed drift during import

## 3.0.3

BUG FIXES:

* **vastdata_protected_path**: Fixed `estimated_read_only_time` field type inconsistency. The field changed from String to Number between v2.1.1 and v3.0.0, which prevented Terraform from automatically converting existing state. The field is now forced to always be a string type, ensuring compatibility and preventing import/state migration issues.
* **vastdata_administrator_manager**: Fixed import error where the `roles` field was returning objects with `{id, name}` structure instead of just role IDs. Import now properly converts roles to a slice of IDs, resolving "unexpected type map[string]interface {} for int" errors during resource import operations.

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
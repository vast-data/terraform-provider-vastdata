---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_ldap Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_ldap (Resource)



## Example Usage

```terraform
resource "vastdata_ldap" "ldap1" {
  domain_name        = "VastEng.lab"
  urls               = ["ldap://10.27.252.30"]
  binddn             = "cn=admin,dc=qa,dc=vastdata,dc=com"
  searchbase         = "dc=qa,dc=vastdata,dc=com"
  bindpw             = "<password>"
  use_auto_discovery = "false"
  use_ldaps          = "false"
  port               = "389"
  method             = "simple"
  query_groups_mode  = "COMPATIBLE"
  use_tls            = "false"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) FQDN of the domain.

### Optional

- `active_directory` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `binddn` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Distinguished name of the LDAP superuser.
- `bindpw` (String, Sensitive) (Valid for versions: 5.0.0,5.1.0,5.2.0) Password for the LDAP superuser.
- `gid_number` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for gid number.
- `group_login_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The attribute used to query the provider for the group login name in NFS ID mapping.
- `group_searchbase` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Base DN for group queries within the joined domain only. When auto-discovery is enabled, group queries outside the joined domain use automatically discovered base DNs.
- `is_vms_auth_provider` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies whether the LDAP is to be used for VMS authentication. There can be only two LDAP configurations that can be used for VMS authentication: one with Active Directory and the other without Active Directory.
- `mail_property_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `match_user` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for user matching.
- `method` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Bind authentication method. Allowed Values are [simple sasl anonymous]
- `port` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) LDAP server port. 389 (LDAP)  636 (LDAPS)
- `posix_account` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for posix account.
- `posix_attributes_source` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `posix_group` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for posix group.
- `posix_primary_provider` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) POSIX primary provider.
- `query_groups_mode` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Query group mode. Allowed Values are [COMPATIBLE RFC2307BIS RFC2307 NONE]
- `query_posix_attributes_from_gc` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) If 'true', users/groups from non-joined domain POSIX attributes are supported. If 'false', POSIX attributes of users/groups from non-joined domain are not supported.
- `reverse_lookup` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `searchbase` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The base DN is the starting point that the LDAP provider uses when searching for users and groups. If a group base DN is configured, it will be used instead of the base DN, for groups only.
- `tls_certificate` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `uid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for uid.
- `uid_member` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for uid member.
- `uid_member_value_property_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `uid_number` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for uid number.
- `urls` (List of String) (Valid for versions: 5.0.0,5.1.0,5.2.0) A list of URIs of LDAP servers (Domain Controllers in Active Directory), in priority order. The URI with highest priority that has a good health status is used. Specify each URI in the format '<scheme>://<address>'. '<address>' can be either a DNS name or an IP address, for example: 'ldap://ldap.company.com, ldaps://ldaps.company.com, ldap://192.0.2.2'.
- `use_auto_discovery` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) If 'true', Active Directory Domain Controllers (DCs) and Active Directory domains are automatically discovered. Queries extend beyond the joined domain to all domains in the forest. If 'false', queries are restricted to the joined domain and URIs must be provided in 'urls'.
- `use_ldaps` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies whether to use LDAPS for auto-discovery.
- `use_tls` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies whether to configure LDAP with TLS.
- `user_login_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The attribute used to query the provider for the user login name in NFS ID mapping.
- `username_property_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Username property name.

### Read-Only

- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import vastdata_ldap.example <guid>
terraform import vastdata_ldap.example <Domain name>
```

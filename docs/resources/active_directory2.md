---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_active_directory2 Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_active_directory2 (Resource)



## Example Usage

```terraform
resource "vastdata_active_directory2" "ad1" {
  machine_account_name = "vast-cluster01"
  organizational_unit  = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
  use_auto_discovery   = false
  binddn               = "cn=admin,dc=qa,dc=vastdata,dc=com"
  searchbase           = "dc=qa,dc=vastdata,dc=com"
  bindpw               = "<password>"
  use_ldaps            = "false"
  domain_name          = "VastEng.lab"
  method               = "simple"
  query_groups_mode    = "COMPATIBLE"
  use_tls              = "false"
  urls                 = ["ldap://198.51.100.3"]

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `machine_account_name` (String) (Valid for versions: 5.1.0,5.2.0) Name of the computer object/machine account to add. Recommended to be the name of the cluster

### Optional

- `binddn` (String) (Valid for versions: 5.1.0,5.2.0) Distinguished name of AD superuser
- `bindpw` (String, Sensitive) (Valid for versions: 5.1.0,5.2.0) The password used with the Bind DN to authenticate to the AD server.
- `domain_name` (String) (Valid for versions: 5.1.0,5.2.0) FQDN of the domain.
- `gid_number` (String) (Valid for versions: 5.1.0,5.2.0) The attribute of a user entry on the AD server that contains the UID number, if different from ‘uidNumber’. Often when binding VAST Cluster to AD this does not need to be set.
- `group_login_name` (String) (Valid for versions: 5.1.0,5.2.0) Specifies the attribute used to query AD for the group login name in NFS ID mapping. Applicable only with AD and NFSv4.1.
- `is_vms_auth_provider` (Boolean) (Valid for versions: 5.1.0,5.2.0) Enables use of the AD for VMS authentication. Two AD configurations per cluster can be used for VMS authentication: one with AD and one without.
- `mail_property_name` (String) (Valid for versions: 5.1.0,5.2.0) Specifies the attribute to use for the user’s email address.
- `match_user` (String) (Valid for versions: 5.1.0,5.2.0) The attribute to use when querying a provider for a user that matches a user that was already retrieved from another provider. A user entry that contains a matching value in this attribute will be considered the same user as the user previously retrieved.
- `method` (String) (Valid for versions: 5.1.0,5.2.0) Bind Authentication Method Allowed Values are [simple anonymous sasl]
- `ntlm_enabled` (Boolean) (Valid for versions: 5.1.0,5.2.0) Manages support of NTLM authentication method for SMB protocol.
- `organizational_unit` (String) (Valid for versions: 5.1.0,5.2.0) Organizational Unit within AD where the Cluster Machine account will be created. If left empty, it will go into default Computers OU
- `port` (Number) (Valid for versions: 5.1.0,5.2.0) Which port to use
- `posix_account` (String) (Valid for versions: 5.1.0,5.2.0) The object class that defines a user entry on the AD server, if different from ‘posixAccount’. When binding VAST Cluster to AD, set this parameter to ‘user’ in order for authorization to work properly.
- `posix_attributes_source` (String) (Valid for versions: 5.1.0,5.2.0) Defines which domains POSIX attributes will be supported from.
- `posix_group` (String) (Valid for versions: 5.1.0,5.2.0)  The object class that defines a group entry on the AD server, if different from ‘posixGroup’. When binding VAST Cluster to AD, set this parameter to ‘group’ in order for authorization to work properly.
- `query_groups_mode` (String) (Valid for versions: 5.1.0,5.2.0) Query group mode Allowed Values are [COMPATIBLE NONE RFC2307BIS RFC2307]
- `reverse_lookup` (Boolean) (Valid for versions: 5.1.0,5.2.0) Resolve AD netgroups into hostnames
- `searchbase` (String) (Valid for versions: 5.1.0,5.2.0) The Base DN is the starting point the AD provider uses when searching for users and groups. If the Group Base DN is configured it will be used instead of the Base DN, for groups only
- `smb_allowed` (Boolean) (Valid for versions: 5.1.0,5.2.0) Indicates if AD is allowed for SMB.
- `tls_certificate` (String) (Valid for versions: 5.1.0,5.2.0) TLS certificate to use for verifying the remote AD server’s TLS certificate.
- `uid` (String) (Valid for versions: 5.1.0,5.2.0) The attribute of a user entry on the AD server that contains the user name, if different from ‘uid’ When binding VAST Cluster to AD, you may need to set this to ‘sAMAccountname’.
- `uid_member` (String) (Valid for versions: 5.1.0,5.2.0) The attribute of a group entry on the AD server that contains names of group members, if different from ‘memberUid’. When binding VAST Cluster to AD, you may need to set this to ‘memberUID’.
- `uid_member_value_property_name` (String) (Valid for versions: 5.1.0,5.2.0) Specifies the attribute which represents the value of the AD group’s member property.
- `uid_number` (String) (Valid for versions: 5.1.0,5.2.0) The attribute of a user entry on the AD server that contains the UID number, if different from ‘uidNumber’. Often when binding VAST Cluster to AD this does not need to be set.
- `urls` (List of String) (Valid for versions: 5.1.0,5.2.0) Comma separated list of URIs of AD servers in the format SCHEME://ADDRESS. The order of listing defines the priority order. The URI with highest priority that has a good health status is used.
- `use_auto_discovery` (Boolean) (Valid for versions: 5.1.0,5.2.0) When enabled, Active Directory Domain Controllers (DCs) and Active Directory domains are auto discovered. Queries extend beyond the joined domain to all domains in the forest. When disabled, queries are restricted to the joined domain and DCs must be provided in the URLs field.
- `use_ldaps` (Boolean) (Valid for versions: 5.1.0,5.2.0) Use LDAPS for auto-Discovery. To activate, set use_auto_discovery to true also.
- `use_multi_forest` (Boolean) (Valid for versions: 5.1.0,5.2.0) Allow access for users from trusted domains on other forests.
- `use_tls` (Boolean) (Valid for versions: 5.1.0,5.2.0) Set to true to enable use of TLS to secure communication between VAST Cluster and the AD server.
- `user_login_name` (String) (Valid for versions: 5.1.0,5.2.0) Specifies the attribute used to query AD for the user login name in NFS ID mapping. Applicable only with AD and NFSv4.1.
- `username_property_name` (String) (Valid for versions: 5.1.0,5.2.0) The attribute to use for querying users in VMS user-initated user queries. Default is ‘name’. Sometimes set to ‘cn’

### Read-Only

- `guid` (String) (Valid for versions: 5.1.0,5.2.0) A uniqe ID given to this resource
- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import vastdata_active_directory2.example <guid>
terraform import vastdata_active_directory2.example <Machine Account Name>
```

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_non_local_group Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_non_local_group (Resource)



## Example Usage

```terraform
resource "vastdata_non_local_group" "ExternalGroup" {
    gid                 = 10000
    tenant_id           = 1
    s3_policies_ids     = [
        1
    ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `context` (String) (Valid for versions: 5.1.0,5.2.0) Context from which the user originates. Available: 'ad', 'nis' and 'ldap'
- `gid` (Number) (Valid for versions: 5.1.0,5.2.0) Group GID
- `tenant_id` (Number) (Valid for versions: 5.1.0,5.2.0) Tenant ID

### Optional

- `groupname` (String) (Valid for versions: 5.1.0,5.2.0) Groupname
- `s3_policies_ids` (List of Number) (Valid for versions: 5.1.0,5.2.0) List S3 policies IDs
- `sid` (String) (Valid for versions: 5.1.0,5.2.0) Group SID

### Read-Only

- `id` (String) (Valid for versions: 5.1.0,5.2.0) The NonLocalGroup identifier

## Import

Import is supported using the following syntax:

```shell
terraform import vastdata_non_local_group.example <Groupname>|<Context>|<Tenant ID>
```

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_user Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_user (Data Source)



## Example Usage

```terraform
data "vastdata_user" "user1" {
  name = "user1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) A uniq name given to the user

### Read-Only

- `allow_create_bucket` (Boolean) Allow create bucket
- `allow_delete_bucket` (Boolean) Allow delete bucket
- `gids` (List of Number) List of supplementary GID list
- `group_count` (Number) Group Count
- `groups` (List of String) List of supplementary Group list
- `guid` (String) A uniq guid given to the user
- `id` (Number) A uniq id given to user
- `leading_gid` (Number) The user leading unix GID
- `leading_group_gid` (Number) Leading Group GID
- `leading_group_name` (String) Leading Group Name
- `local` (Boolean) IS this a local user
- `primary_group_sid` (String) The user primary group SID
- `s3_policies_ids` (List of Number) List S3 policies IDs
- `s3_superuser` (Boolean) Is S3 superuser
- `sid` (String) The user SID
- `sids` (List of String) supplementary SID list
- `uid` (Number) The user unix UID

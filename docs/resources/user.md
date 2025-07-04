---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_user Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_user (Resource)



## Example Usage

```terraform
#Create a user named example with UID of 9000
resource "vastdata_user" "example-user" {
  name = "example"
  uid  = 9000
}

#Create a user named user1 with leading group and supplementary groups
resource "vastdata_group" "group2" {
  name = "group2"
  gid  = 2000
}

resource "vastdata_group" "group4" {
  name = "group4"
  gid  = 4000
}


resource "vastdata_user" "user1" {
  name        = "user1"
  uid         = 3000
  leading_gid = resource.vastdata_group.group1.gid
  gids = [
    resource.vastdata_group.group2.gid,
    resource.vastdata_group.group4.gid
  ]

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The unique name of the user.

### Optional

- `allow_create_bucket` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Allows or prohibits bucket creation by the user.
- `allow_delete_bucket` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Allows or prohibits bucket deletion by the user.
- `gids` (List of Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) A list of supplementary GIDs.
- `groups` (List of String) (Valid for versions: 5.0.0,5.1.0,5.2.0) A list of supplementary groups.
- `leading_gid` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The user's leading Unix GID.
- `primary_group_sid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The user's primary group SID.
- `s3_policies_ids` (List of Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) A list of identity policy IDs.
- `s3_superuser` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) If 'true', the user is an S3 superuser.
- `uid` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The user's Unix UID.

### Read-Only

- `group_count` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) Group count.
- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The unique GUID of the user.
- `id` (String) The ID of this resource.
- `leading_group_gid` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The GID of the leading group.
- `leading_group_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the leading group.
- `local` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) If 'true', the user is a local user.
- `sid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The user's SID.
- `sids` (List of String) (Valid for versions: 5.0.0,5.1.0,5.2.0) A list of supplementary SIDs.

## Import

Import is supported using the following syntax:

```shell
terraform import vastdata_user.example <guid>
terraform import vastdata_user.example <Name>
```

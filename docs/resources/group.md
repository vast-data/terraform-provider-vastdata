---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_group Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_group (Resource)



## Example Usage

```terraform
# Create a group with the name group1 ang GID 1000
resource "vastdata_group" "group1" {
  name = "group1"
  gid  = 1000
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `gid` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The Unix GID of the group.
- `name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The unique name of the group.

### Optional

- `s3_policies_ids` (List of Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) A list of identity policy IDs.

### Read-Only

- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The unique GUID of the group.
- `id` (String) The ID of this resource.
- `sid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The SID of the group.

## Import

Import is supported using the following syntax:

```shell
terraform import vastdata_group.example <guid>
terraform import vastdata_group.example <Name>
```

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_protection_policy Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_protection_policy (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `clone_type` (String)
- `name` (String)
- `prefix` (String)

### Optional

- `frames` (Block List) List of snapshots schedules (see [below for nested schema](#nestedblock--frames))
- `indestructible` (Boolean) Is the snapshot indestructable
- `target_name` (String) The target peer name
- `target_object_id` (Number) The id of the target peer
- `url` (String) Direct link to the replication policy

### Read-Only

- `guid` (String) A unique guid given to the  replication peer configuration
- `id` (String) The ID of this resource.

<a id="nestedblock--frames"></a>
### Nested Schema for `frames`

Optional:

- `every` (String) How often to make a snapshot
- `keep_local` (String)
- `keep_remote` (String)
- `start_at` (String)
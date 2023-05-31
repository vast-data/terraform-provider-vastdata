---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_protection_policy Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_protection_policy (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the replication peer configuration

### Optional

- `clone_type` (String) The type the replication
- `frames` (Block List) List of snapshots schedules (see [below for nested schema](#nestedblock--frames))
- `id` (Number) A unique id given to the replication peer configuration
- `indestructible` (Boolean) Is the snapshot indestructable
- `prefix` (String) The prefix to be given to the replicated data
- `target_name` (String) The target peer name
- `target_object_id` (Number) The id of the target peer
- `url` (String) Direct link to the replication policy

### Read-Only

- `guid` (String) A unique guid given to the  replication peer configuration

<a id="nestedblock--frames"></a>
### Nested Schema for `frames`

Optional:

- `every` (String) How often to make a snapshot
- `keep_local` (String)
- `keep_remote` (String)
- `start_at` (String)
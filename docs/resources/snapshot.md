---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_snapshot Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_snapshot (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `expiration_time` (String) When will this sanpshot expire
- `locked` (Boolean) Is it locked (indestructable)
- `path` (String) The path to make snapshot from
- `tenant_id` (Number) The tenant id to use

### Read-Only

- `guid` (String) A unique guid given to the snapshot
- `id` (String) The ID of this resource.
---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_replication_peers Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_replication_peers (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the replication peer configuration

### Optional

- `id` (Number) A unique id given to the replication peer configuration
- `is_local` (Boolean) Is the source of the replication local (this host is the source)
- `leading_vip` (String) The vip provided for the replication peer configuration
- `peer_name` (String) The name of the peer cluster
- `remote_version` (String) The version of the remote peer
- `remote_vip_range` (String) The vip range which were reported by the peer
- `secure_mode` (String) Is the connection secure
- `url` (String) Direct url of the replication peer configurations
- `version` (String) The version of the source

### Read-Only

- `guid` (String) A unique guid given to the  replication peer configuration
# Copyright (c) HashiCorp, Inc.

variable policy_name {
    type = string
}

variable view_name {
    type = string
}

variable host_name {
    type = string
}

variable host_nqn {
    type = string
}

variable host2_name {
    type = string
}

variable host2_nqn {
    type = string
}

variable volume_name {
    type = string
}

variable volume_size {
    type = number
}

variable volume_tags {
    type = list(string)
}

variable mapping_hosts {
    type = string
}

resource vastdata_view_policy block_policy {
    name = var.policy_name
    flavor = "NFS"
    nfs_no_squash = ["10.0.0.1","10.0.0.2"]
}

resource vastdata_view block_view {
    path = "/${var.view_name}"
    name = var.view_name
    is_default_subsystem = false
    policy_id = vastdata_view_policy.block_policy.id
    create_dir = true
    protocols = ["BLOCK"]
}

resource vastdata_blockhost blockhost1 {
    name = var.host_name
    nqn = var.host_nqn
}

resource vastdata_blockhost blockhost2 {
    name = var.host2_name
    nqn = var.host2_nqn
}

resource vastdata_volume volume1 {
    name = "/${var.volume_name}"
    size = var.volume_size
    view_id = vastdata_view.block_view.id
    volume_tags = var.volume_tags
}

resource vastdata_block_mapping volume1_blockhost1_mapping {
    volume_id = vastdata_volume.volume1.id
    hosts_ids = var.mapping_hosts == "all" ? [vastdata_blockhost.blockhost1.id, vastdata_blockhost.blockhost2.id] : [vastdata_blockhost.blockhost1.id]
}

output tf_blockhost {
  value = vastdata_blockhost.blockhost1
}

output tf_volume {
  value = vastdata_volume.volume1
}

output tf_block_mapping {
  value = vastdata_block_mapping.volume1_blockhost1_mapping
}

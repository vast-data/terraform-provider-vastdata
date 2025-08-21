# Copyright (c) HashiCorp, Inc.

# To refer to a specific system, need to add `provider = vastdata.system<IDX>` where
# IDX is the system's index (in the order of appearance in Comet's commandline)

variable group_name {
    type = string
}

variable group_gid {
    type = number
}

resource vastdata_group group1 {
  name = var.group_name
  gid = var.group_gid
}

output tf_group {
  value = vastdata_group.group1
}

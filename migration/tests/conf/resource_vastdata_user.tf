# Copyright (c) HashiCorp, Inc.

variable user_name {
    type = string
}

variable user_uid {
    type = number
}

variable group_name {
    type = string
}

variable group_gid {
    type = number
}

resource vastdata_group user_group1 {
  name = var.group_name
  gid = var.group_gid
}

resource vastdata_user user1 {
  name = var.user_name
  uid = var.user_uid
  leading_gid = vastdata_group.user_group1.gid
}

output tf_user {
  value = vastdata_user.user1
}
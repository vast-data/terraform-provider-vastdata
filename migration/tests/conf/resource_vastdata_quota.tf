# Copyright (c) HashiCorp, Inc.

variable quota_name {
    type = string
}

variable quota_user_name {
    type = string
}

variable quota_view_policy_name {
    type = string
}

variable quota_view_path {
    type = string
}

variable quota_user_uid {
    type = number
}

variable quota_group_name {
    type = string
}

variable quota_group_gid {
    type = number
}

variable quota_default_grace_period {
    type = string
}

variable quota_user_grace_period {
    type = string
}

variable quota_group_grace_period {
    type = string
}

resource vastdata_group group3 {
  name = var.quota_group_name
  gid = var.quota_group_gid
}

resource vastdata_user quota_user {
  name = var.quota_user_name
  uid = var.quota_user_uid
}

#Create a quota with user & group quota
resource vastdata_view_policy view_policy2 {
   name = var.quota_view_policy_name
   flavor = "NFS"
   nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

resource vastdata_view view1 {
  path = var.quota_view_path
  policy_id = vastdata_view_policy.view_policy2.id
  create_dir = "true"
  protocols = ["NFS","NFS4"]
}

resource vastdata_quota quota1 {
  name = var.quota_name
  default_email = "user@example.com"
  path = vastdata_view.view1.path
  soft_limit = 100000
  hard_limit = 100000
  is_user_quota = true
  default_user_quota {
    grace_period = var.quota_default_grace_period
    hard_limit = 2000
    soft_limit = 1000
    hard_limit_inodes = 20000000
  }
  user_quotas {
    grace_period = var.quota_user_grace_period
    hard_limit = 15000
    soft_limit = 15000
    entity {
      name = var.quota_user_name
      email = "user1@example.com"
      identifier = var.quota_user_name
      identifier_type = "username"
      is_group = "false"
    }
  }
  group_quotas {
    grace_period = var.quota_group_grace_period
    hard_limit = 15000
    soft_limit = 15000
    entity {
      name = vastdata_group.group3.name
      identifier = vastdata_group.group3.name
      identifier_type = "group"
      is_group = "false"
    }
  }
}

output tf_quota {
  value = vastdata_quota.quota1
}

output tf_quota_user {
  value = vastdata_user.quota_user
}

output tf_quota_view_policy {
  value = vastdata_view_policy.view_policy2
}

output tf_quota_view {
  value = vastdata_view.view1
}
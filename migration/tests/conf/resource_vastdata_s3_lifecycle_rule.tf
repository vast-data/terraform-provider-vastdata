# Copyright (c) HashiCorp, Inc.

variable s3_user_uid {
    type = number
}

variable s3_lifecycle_name {
    type = string
}

variable s3_lifecycle_user_name {
    type = string
}

variable s3_lifecycle_enabled {
    type = bool
}

resource vastdata_user s3_user1 {
  name = var.s3_lifecycle_user_name
  uid = var.s3_user_uid
}

resource vastdata_view_policy s3_viewpolicy1 {
   name = "tf_s3_viewpolicy1"
   flavor = "S3_NATIVE"
   nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

resource vastdata_view s3_view1 {
  policy_id = vastdata_view_policy.s3_viewpolicy1.id
  path = "/tf_s3view1"
  bucket = "tf-s3view1"
  protocols = ["S3"]
  bucket_owner = vastdata_user.s3_user1.name
  create_dir = "true"
}

resource vastdata_s3_life_cycle_rule s3_lifecycle_rule1 {
  name = var.s3_lifecycle_name
  max_size = 10000000
  min_size = 100000
  newer_noncurrent_versions = 3
  prefix = "prefix"
  view_id = vastdata_view.s3_view1.id
  expiration_days = 30
  enabled = var.s3_lifecycle_enabled
  noncurrent_days  = 20
}

output tf_s3_lifecycle_rule {
  value = vastdata_s3_life_cycle_rule.s3_lifecycle_rule1
}

output tf_s3_lifecycle_rule_user {
  value = vastdata_user.s3_user1
}

output tf_s3_lifecycle_rule_view_policy {
  value = vastdata_view_policy.s3_viewpolicy1
}

output tf_s3_lifecycle_rule_view {
  value = vastdata_view.s3_view1
}

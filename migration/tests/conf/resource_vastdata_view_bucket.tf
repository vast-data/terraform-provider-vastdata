# Copyright (c) HashiCorp, Inc.

variable bucket_user_name {
    type = string
}

variable bucket_user_uid {
    type = number
}

variable view_policy_name {
    type = string
}

variable view_path {
    type = string
}

variable bucket_log_prefix {
    type = string
}

variable bucket_log_key_format {
    type = string
}

resource vastdata_user bucket_user1 {
  name = var.bucket_user_name
  uid = var.bucket_user_uid
}

resource vastdata_view_policy bucket_vpolicy1 {
    name = var.view_policy_name
    flavor = "S3_NATIVE"
    nfs_no_squash = ["10.0.0.1","10.0.0.2"]
    allowed_characters = "LCD"
    auth_source = "RPC"
}

resource vastdata_view bucket_view1 {
  path = "/${var.view_path}_1"
  bucket = "${var.view_path}1"
  bucket_owner = vastdata_user.bucket_user1.name
  policy_id = vastdata_view_policy.bucket_vpolicy1.id
  create_dir = "true"
  protocols = ["S3"]
}

resource vastdata_view bucket_view2 {
  path = "/${var.view_path}_2"
  bucket = "${var.view_path}2"
  bucket_owner = vastdata_user.bucket_user1.name
  policy_id = vastdata_view_policy.bucket_vpolicy1.id
  create_dir = "true"
  protocols = ["S3"]
  bucket_logging {
    destination_id = vastdata_view.bucket_view1.id
    prefix = var.bucket_log_prefix
    key_format = var.bucket_log_key_format
  }
}

output tf_view1 {
  value = vastdata_view.bucket_view1
}

output tf_view2 {
  value = vastdata_view.bucket_view2
}

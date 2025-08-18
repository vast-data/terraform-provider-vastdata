# Copyright (c) HashiCorp, Inc.

# To refer to a specific system, need to add `provider = vastdata.system<IDX>` where
# IDX is the system's index (in the order of appearance in Comet's commandline)

variable view_path {
    type = string
}

variable view_protocols {
    type = list(string)
}

variable view_policy_name {
    type = string
}

variable user_name {
    type = string
}

variable user_uid {
    type = number
}

variable is_seamless {
    type = bool
}

variable s3_object_ownership_rule {
    type = string
}

variable locking {
    type = bool
}

resource vastdata_user s3_view_user1 {
  name = var.user_name
  uid = var.user_uid
}

# Create a view with NFS & NFSv4 protocols
resource vastdata_view_policy s3_viewpolicy1 {
   name = var.view_policy_name
   flavor = "S3_NATIVE"
   nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

resource vastdata_view s3_view2 {
  path = var.view_path
  policy_id = vastdata_view_policy.s3_viewpolicy1.id
  create_dir = "true"
  protocols = var.view_protocols
  bucket_creators = ["${var.user_name}@VastENG.lab"]

  is_seamless = var.is_seamless
  s3_object_ownership_rule = var.s3_object_ownership_rule
  locking = var.locking
}

output tf_view {
  value = vastdata_view.s3_view2
}

output tf_view_view_policy {
  value = vastdata_view_policy.s3_viewpolicy1
}
# Copyright (c) HashiCorp, Inc.

variable view_bucket_name {
    type = string
}

variable view_protocols {
    type = list(string)
}

variable user_name {
    type = string
}

variable user_uid {
    type = number
}

variable acl_perm_level {
    type = string
}

resource vastdata_user s3_view_user1 {
  name = var.user_name
  uid = var.user_uid
}

# Create a view with NFS & NFSv4 protocols
resource vastdata_view_policy viewpolicy1 {
   name = "tf_viewpolicy_${var.view_bucket_name}"
   flavor = "S3_NATIVE"
   nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

resource vastdata_view view1 {
  path = "/${var.view_bucket_name}"
  policy_id = vastdata_view_policy.viewpolicy1.id
  create_dir = "true"
  protocols = var.view_protocols
  bucket = "${var.view_bucket_name}"
  bucket_owner = "10950303@VastENG.lab"
  share_acl {
    acl {
        name = vastdata_user.s3_view_user1.name
        grantee="users"
        fqdn="All"
        permissions = "${var.acl_perm_level}"
    }
    enabled = true
  }
}

output tf_view {
  value = vastdata_view.view1
}

output tf_view_view_policy {
  value = vastdata_view_policy.viewpolicy1
}

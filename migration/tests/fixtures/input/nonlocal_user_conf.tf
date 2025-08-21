# Copyright (c) HashiCorp, Inc.

variable user_uid {
    type = number
}

variable context {
    type = string
}

variable tenant_name {
    type = string
}

resource "vastdata_non_local_user" "non_local_user1" {
    uid = var.user_uid
    context = var.context
    allow_create_bucket = true
    allow_delete_bucket = false
    s3_policies_ids = [1, 2, 3]
}

data "vastdata_non_local_user" "user_data" {
    uid = var.user_uid
    context = var.context
}

output tf_user {
  value = vastdata_non_local_user.non_local_user1
}

output tf_user_ds {
    value = data.vastdata_non_local_user.user_data
}

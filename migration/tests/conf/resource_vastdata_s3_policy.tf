# Copyright (c) HashiCorp, Inc.

variable s3_policy_name {
    type = string
}

variable s3_policy_enabled {
    type = bool
}

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

resource vastdata_s3_policy s3policy {
        name = var.s3_policy_name
        policy = <<EOT
        {
   "Version":"2012-10-17",
   "Statement":[
      {
         "Effect":"Allow",
         "Action": "s3:ListAllMyBuckets",
         "Resource":"*"
      },
      {
         "Effect":"Allow",
         "Action":["s3:ListObjects","s3:GetBucketLocation"],
         "Resource":"arn:aws:s3:::DOC-EXAMPLE-BUCKET1"
      },
      {
         "Effect":"Allow",
         "Action":[
            "s3:PutObject",
            "s3:PutObjectAcl",
            "s3:GetObject",
            "s3:GetObjectAcl",
            "s3:DeleteObject"
         ],
         "Resource":"arn:aws:s3:::DOC-EXAMPLE-BUCKET1/*"
      }
   ]
}
        EOT
        enabled = var.s3_policy_enabled
}

resource vastdata_group s3policy_group1 {
  name = var.group_name
  gid = var.group_gid
  s3_policies_ids = [vastdata_s3_policy.s3policy.id]
}

resource vastdata_user s3policy_user1 {
  name = var.user_name
  uid = var.user_uid
  s3_policies_ids = [vastdata_s3_policy.s3policy.id]
}

output tf_s3_policy {
  value = vastdata_s3_policy.s3policy
}

output tf_user {
  value = vastdata_user.s3policy_user1
}

output tf_group {
  value = vastdata_group.s3policy_group1
}
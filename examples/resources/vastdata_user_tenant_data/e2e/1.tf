data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
}

resource "vastdata_s3_policy" "vastdb_s3policy1" {
  name      = "vastdb_s3policy1"
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  policy    = <<EOT
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
}

resource "vastdata_s3_policy" "vastdb_s3policy2" {
  name      = "vastdb_s3policy2"
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  policy    = <<EOT
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
                 "Resource":"arn:aws:s3:::DOC-EXAMPLE-BUCKET2"
              },
              {
                 "Effect":"Allow",
                 "Action":[
                    "s3:PutObject",
                    "s3:PutObjectAcl"
                 ],
                 "Resource":"arn:aws:s3:::DOC-EXAMPLE-BUCKET2/*"
              }
           ]
        }
        EOT
}

resource "vastdata_user_tenant_data" "vastdb_user_tenant_data" {
  user_id             = vastdata_user.vastdb_user.id
  allow_create_bucket = true
  allow_delete_bucket = true
  s3_policies_ids = [
    vastdata_s3_policy.vastdb_s3policy1.id,
    vastdata_s3_policy.vastdb_s3policy2.id
  ]
}


data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_s3_policy" "vastdb_s3policy" {
  name      = "vastdb_s3policy"
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


resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
}


resource "vastdata_nonlocal_user" "vastdb_nonlocal_user" {
  uid                 = vastdata_user.vastdb_user.uid
  tenant_id           = data.vastdata_tenant.vastdb_tenant.id
  allow_create_bucket = false
  allow_delete_bucket = true
  s3_superuser        = false
  s3_policies_ids = [
    vastdata_s3_policy.vastdb_s3policy.id
  ]
}
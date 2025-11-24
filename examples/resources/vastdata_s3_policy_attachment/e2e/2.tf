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
  name              = "vastdb_user"
  uid               = 5001
  local_provider_id = 1
}

resource "vastdata_group" "vastdb_group" {
  name              = "vastdb_group"
  gid               = 5002
  local_provider_id = 1
}

resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment1" {
  s3_policy_id = vastdata_s3_policy.vastdb_s3policy.id
  groupname    = vastdata_group.vastdb_group.name
}

resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment2" {
  s3_policy_id   = vastdata_s3_policy.vastdb_s3policy.id
  sid            = vastdata_user.vastdb_user.sid
  ignore_present = true
}

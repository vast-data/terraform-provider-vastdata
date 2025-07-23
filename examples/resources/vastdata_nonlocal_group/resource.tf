
resource "vastdata_nonlocal_group" "vastdb_nonlocal_group" {
  gid = 1000
}

# ---------------------
# Complete examples
# ---------------------


data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_group" "vastdb_group" {
  name = "vastdb_group"
  gid  = 1001
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

resource "vastdata_nonlocal_group" "vastdb_nonlocal_group" {
  gid       = vastdata_group.vastdb_group.gid
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  s3_policies_ids = [
    vastdata_s3_policy.vastdb_s3policy1.id,
    vastdata_s3_policy.vastdb_s3policy2.id
  ]
}

# --------------------


resource "vastdata_group" "vastdb_group" {
  name = "vastdb_group"
  gid  = 1001
}

resource "vastdata_nonlocal_group" "vastdb_nonlocal_group" {
  gid = vastdata_group.vastdb_group.gid
}

# --------------------


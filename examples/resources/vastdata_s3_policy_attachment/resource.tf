resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment1" {
  s3_policy_id      = 1
  gid               = 5001
  local_provider_id = 1
}

# Create another S3 policy attachment with ignore_present set to true.
resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment2" {
  s3_policy_id   = 1
  uid            = 5002
  ignore_present = true
}

# ---------------------
# Complete examples
# ---------------------

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
  uid               = 30109
  local_provider_id = 1
}

resource "vastdata_group" "vastdb_group" {
  name              = "vastdb_group"
  gid               = 30097
  local_provider_id = 1
}

resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment1" {
  s3_policy_id = vastdata_s3_policy.vastdb_s3policy.id
  gid          = vastdata_group.vastdb_group.gid
}

resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment2" {
  s3_policy_id   = vastdata_s3_policy.vastdb_s3policy.id
  uid            = vastdata_user.vastdb_user.uid
  ignore_present = true
}

# --------------------

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

# --------------------

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
  uid               = 5003
  local_provider_id = 1
}

resource "vastdata_group" "vastdb_group" {
  name              = "vastdb_group"
  gid               = 5004
  local_provider_id = 1
}

# Test group_sid attachment (recommended for AD groups with invalid GID)
resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment1" {
  s3_policy_id = vastdata_s3_policy.vastdb_s3policy.id
  group_sid    = vastdata_group.vastdb_group.sid
}

# Test username attachment
resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment2" {
  s3_policy_id   = vastdata_s3_policy.vastdb_s3policy.id
  username       = vastdata_user.vastdb_user.name
  ignore_present = true
}

# --------------------


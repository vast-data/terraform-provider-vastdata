
# Create a user with a specific UID.
resource "vastdata_user" "example-user" {
  name = "example"
  uid  = 9000
}


resource "vastdata_local_provider" "vastdb_local_provider1" {
  name = "vastdb_local_provider1"
}

// Create a user with local provider.
resource "vastdata_user" "vastdb_user1" {
  name              = "vastdb_user1"
  uid               = 30117
  local_provider_id = vastdata_local_provider.vastdb_local_provider1.id
}

# ---------------------
# Complete examples
# ---------------------


resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
  gids = [
    1001
  ]
}

# --------------------


resource "vastdata_group" "vastdb_group" {
  name = "vastdb_group"
  gid  = 30097
}

resource "vastdata_user" "vastdb_user" {
  name                = "vastdb_user"
  uid                 = 30109
  local               = true
  allow_create_bucket = true
  allow_delete_bucket = true
  s3_superuser        = false
  leading_gid         = vastdata_group.vastdb_group.gid
  gids = [
    1001,
    vastdata_group.vastdb_group.gid
  ]
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
  name = "vastdb_user"
  uid  = 30109
  s3_policies_ids = [
    vastdata_s3_policy.vastdb_s3policy.id
  ]
}

# --------------------


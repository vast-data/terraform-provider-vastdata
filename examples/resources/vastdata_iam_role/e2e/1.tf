resource "vastdata_tenant" "vastdb_tenant" {
  name = "vastdbtenant"
}


resource "vastdata_iam_role" "vastdb_iam_role" {
  name         = "vastdb_iam_role"
  tenant_id    = vastdata_tenant.vastdb_tenant.id
  trust_policy = <<-EOT
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sts:*"
    }
  ]
}
EOT
}

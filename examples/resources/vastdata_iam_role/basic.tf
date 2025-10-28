resource "vastdata_iam_role" "vastdb_iam_role" {
  name         = "vastdb_iam_role"
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

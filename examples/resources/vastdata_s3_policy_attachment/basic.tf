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

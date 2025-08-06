resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment1" {
  s3_policy_id = 1
  gid          = 1000
}

# Create another S3 policy attachment with ignore_present set to true.
# It will not fail if the user with uid=777 already has s3_policy_id=1 attached.
resource "vastdata_s3_policy_attachment" "vastdb_policy_attachment2" {
  s3_policy_id   = 1
  uid            = 777
  ignore_present = true
}

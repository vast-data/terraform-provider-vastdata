# Example: Check authorization status for S3 GetObject action
data "vastdata_manager_authorized_status" "test" {
  action      = "GetObject"
  bucket_name = "test-bucket"
  leading_vid = 1
  object_path = "/test/object.txt"
  owner_vid   = 1
  tenant_id   = 1
}

# Example: Check different S3 action
data "vastdata_manager_authorized_status" "put_object" {
  action      = "PutObject"
  bucket_name = "upload-bucket"
  leading_vid = 1
  object_path = "/uploads/file.dat"
  owner_vid   = 1
  tenant_id   = 1
}

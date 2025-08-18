# Copyright (c) HashiCorp, Inc.

# Test file for attribute renaming transformations

# S3 replication peer with type_ attribute
resource "vastdata_s3_replication_peers" "aws_peer" {
  name          = "aws-peer"
  bucket_name   = "my-aws-s3-bucket"
  http_protocol = "https"
  type_         = "AWS_S3"
  aws_region    = "eu-west-1"
  access_key    = "W21E6X5ZQEOODYB6J0UY"
  secret_key    = "fcESVNih9Ykb/bDSmKipQdinnHObrRyv9nre+nR1"
}

# Administrator role with permissions_list
resource "vastdata_administators_roles" "test_role" {
  name             = "test-role"
  permissions_list = ["create_support", "create_settings", "create_security"]
}

# Administrator manager with both attribute types
resource "vastdata_administators_managers" "test_manager" {
  username         = "test-manager"
  password         = "SecurePassword123"
  permissions_list = ["create_monitoring", "view_logs"]
  roles            = [1, 2, 3]
}

# Custom S3 peer with type_
resource "vastdata_s3_replication_peers" "custom_peer" {
  name              = "custom-peer"
  bucket_name       = "custom-bucket"
  http_protocol     = "https"
  type_             = "CUSTOM_S3"
  custom_bucket_url = "customs3.bucket.com"
  access_key        = "CUSTOM_ACCESS_KEY"
  secret_key        = "custom_secret_key"
}

# Complex example with nested attributes
resource "vastdata_administators_managers" "complex_manager" {
  username = "complex-manager"
  password = "ComplexPassword456"
  
  # This should be transformed
  permissions_list = [
    "create_support",
    "create_monitoring", 
    "view_logs"
  ]
  
  # Dynamic content that should be preserved
  dynamic "access_controls" {
    for_each = var.access_controls
    content {
      type_ = access_controls.value.type
      permissions_list = access_controls.value.permissions
    }
  }
}

# Copyright (c) HashiCorp, Inc.

# Test file for attribute renaming transformations

# S3 replication peer with type attribute (renamed from type_)
resource "vastdata_s3_replication_peer" "aws_peer" {
  name          = "aws-peer"
  bucket_name   = "my-aws-s3-bucket"
  http_protocol = "https"
  type         = "AWS_S3"
  aws_region    = "eu-west-1"
  access_key    = "W21E6X5ZQEOODYB6J0UY"
  secret_key    = "fcESVNih9Ykb/bDSmKipQdinnHObrRyv9nre+nR1"
}

# Administrator role with permissions (renamed from permissions_list)
resource "vastdata_administrator_role" "test_role" {
  name             = "test-role"
  permissions = ["create_support", "create_settings", "create_security"]
}

# Administrator manager with both attribute renames
resource "vastdata_administrator_manager" "test_manager" {
  username         = "test-manager"
  password         = "SecurePassword123"
  permissions = ["create_monitoring", "view_logs"]
  roles            = [1, 2, 3]
}

# Custom S3 peer with type (renamed from type_)
resource "vastdata_s3_replication_peer" "custom_peer" {
  name              = "custom-peer"
  bucket_name       = "custom-bucket"
  http_protocol     = "https"
  type             = "CUSTOM_S3"
  custom_bucket_url = "customs3.bucket.com"
  access_key        = "CUSTOM_ACCESS_KEY"
  secret_key        = "custom_secret_key"
}

# Complex example with nested attributes
resource "vastdata_administrator_manager" "complex_manager" {
  username = "complex-manager"
  password = "ComplexPassword456"
  
  # This should be transformed
  permissions = [
    "create_support",
    "create_monitoring", 
    "view_logs"
  ]
  
  # Dynamic content that should be preserved
  dynamic "access_controls" {
    for_each = var.access_controls
    content {
      type = access_controls.value.type
      permissions = access_controls.value.permissions
    }
  }
}

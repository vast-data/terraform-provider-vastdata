# Copyright (c) HashiCorp, Inc.

# Create a user with a specific UID.
resource "vastdata_user" "example-user" {
  name = "example"
  uid  = 9000
}

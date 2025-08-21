# Copyright (c) HashiCorp, Inc.

# Simple resource for basic testing
resource "vastdata_administrator_manager" "simple" {
  name = "simple-admin"
  enabled = true
}

# Copyright (c) HashiCorp, Inc.

# Simple resource for basic testing
resource "vastdata_administators_managers" "simple" {
  name = "simple-admin"
  enabled = true
}

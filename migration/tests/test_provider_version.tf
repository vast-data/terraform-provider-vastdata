# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
      version = "1.7.0"
    }
  }
}

provider "vastdata" {
  host     = "192.168.1.100"
  username = "admin"
  password = "password"
}

resource "vastdata_administators_managers" "admin" {
  username         = "admin1"
  permissions_list = ["create_support", "create_settings"]
}

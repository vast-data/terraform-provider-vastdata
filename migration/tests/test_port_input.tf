# Copyright (c) HashiCorp, Inc.

resource "vastdata_vip_pool" "prod" {
    name              = "Prod"
    cnode_ids         = [3, 4]
    port_membership   = "ALL"
    enabled           = true
}

resource "vastdata_vip_pool" "tests" {
    name              = "tests"
    cnode_ids         = [1, 2, 3, 4]
    port_membership   = "ALL"
    domain_name       = "testing"
}

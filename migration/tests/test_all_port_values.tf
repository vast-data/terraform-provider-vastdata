# Copyright (c) HashiCorp, Inc.

resource "vastdata_vip_pool" "pool_all" {
    name              = "Pool_All"
    port_membership   = "ALL"
    enabled           = true
}

resource "vastdata_vip_pool" "pool_right" {
    name              = "Pool_Right"
    port_membership   = "RIGHT"
    enabled           = true
}

resource "vastdata_vip_pool" "pool_left" {
    name              = "Pool_Left"
    port_membership   = "LEFT"
    enabled           = true
}

resource "vastdata_vip_pool" "pool_already_lowercase" {
    name              = "Pool_Already"
    port_membership   = "all"
    enabled           = true
}

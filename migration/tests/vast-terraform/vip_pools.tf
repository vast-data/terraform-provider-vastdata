# Copyright (c) HashiCorp, Inc.


# Basic definition of VIP pools
# Both have been imported and is managed in terraform state

resource "vastdata_vip_pool" "dev" {
    name              = "Dev"
    cnode_ids         = [
        1,
        2,
    ]
    domain_name       = "dev"
    enable_l3         = false
    enabled           = true
    gw_ip             = "10.66.20.193"
    gw_ipv6           = null
    port_membership   = "ALL"
    role              = "PROTOCOLS"
    subnet_cidr       = 26
    vlan              = 0
    vms_preferred     = false

  ip_ranges {
      start_ip = "10.66.20.201"
      end_ip   = "10.66.20.208"

    }

}

resource "vastdata_vip_pool" "prod" {
    name              = "Prod"
    cnode_ids         = [
        3,
        4,
    ]
    domain_name       = "prod"
    enable_l3         = false
    enabled           = true
    gw_ip             = "10.66.20.129"
    gw_ipv6           = null
    port_membership   = "ALL"
    role              = "PROTOCOLS"
    subnet_cidr       = 26
    vlan              = 0
    vms_preferred     = false

    ip_ranges {
        start_ip = "10.66.20.141"
        end_ip   = "10.66.20.148"
    }
}

resource "vastdata_vip_pool" "tests" {
    name              = "tests"
    cnode_ids         = [
        1,
        2,
        3,
        4,
    ]
    domain_name       = "testing"
    enable_l3         = false
    enabled           = true
    gw_ip             = "192.168.55.1"
    gw_ipv6           = null
    port_membership   = "ALL"
    role              = "PROTOCOLS"
    subnet_cidr       = 24
    vlan              = 0
    vms_preferred     = false

    ip_ranges {
        start_ip = "192.168.55.20"
        end_ip   = "192.168.55.30"
    }
}



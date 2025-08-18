# Copyright (c) HashiCorp, Inc.

variable vippool_name {
    type = string
}

resource vastdata_vip_pool pool1v6 {
    name = var.vippool_name
    role = "PROTOCOLS"
    subnet_cidr_ipv6 = 64
    ip_ranges {
        start_ip = "fec0:10::11"
        end_ip = "fec0:10::18"
    }
}

output tf_vippool {
  value = vastdata_vip_pool.pool1v6
}
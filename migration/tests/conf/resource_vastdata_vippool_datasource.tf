# Copyright (c) HashiCorp, Inc.

variable vippool_name {
    type = string
}

data "vastdata_vip_pool" "vippool_ds1" {
    name = var.vippool_name
}

output tf_vippool_ds {
    value = data.vastdata_vip_pool.vippool_ds1
}

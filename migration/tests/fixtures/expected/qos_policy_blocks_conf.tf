# Copyright (c) HashiCorp, Inc.

variable qos_policy_name {
    type = string
}

resource vastdata_qos_policy qos1 {
  name = var.qos_policy_name
  static_limits = {
    max_writes_bw_mbps = 110
    max_reads_iops = 200
    max_writes_iops = 3001
  }
  capacity_limits = {
    max_reads_bw_mbps_per_gb_capacity=100
    max_reads_iops_per_gb_capacity = 200
  }
}

output tf_qos_policy {
  value = vastdata_qos_policy.qos1
}

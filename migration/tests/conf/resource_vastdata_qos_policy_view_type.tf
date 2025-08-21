# Copyright (c) HashiCorp, Inc.

variable qos_policy_name {
    type = string
}

variable qos_policy_limit_by {
    type = string
}

variable qos_policy_is_default {
    type = bool
}

variable "max_bw_mbps" {
  type = number
}

variable "burst_bw_mb" {
  type = number
}

variable "burst_loan_mb" {
  type = number
}

variable "max_iops" {
  type = number
}

variable "burst_iops" {
  type = number
}

variable "burst_loan_iops" {
  type = number
}

variable "max_bw_mbps_per_gb_capacity" {
  type = number
}

variable "max_iops_per_gb_capacity" {
  type = number
}


resource vastdata_qos_policy qos2 {
  name = var.qos_policy_name
  policy_type = "VIEW"
  mode = "USED_CAPACITY"

  limit_by = var.qos_policy_limit_by
  is_default = var.qos_policy_is_default

  static_total_limits {
    max_bw_mbps = var.max_bw_mbps
    burst_bw_mb = var.burst_bw_mb
    burst_loan_mb = var.burst_loan_mb
    max_iops = var.max_iops
    burst_iops = var.burst_iops
    burst_loan_iops = var.burst_loan_iops
  }
  capacity_total_limits {
    max_bw_mbps_per_gb_capacity = var.max_bw_mbps_per_gb_capacity
    max_iops_per_gb_capacity = var.max_iops_per_gb_capacity
  }
}

output tf_qos_policy {
  value = vastdata_qos_policy.qos2
}
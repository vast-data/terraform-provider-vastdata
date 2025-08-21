# Copyright (c) HashiCorp, Inc.

variable qos_user_name {
    type = string
}

variable qos_user_uid {
    type = number
}

variable qos_policy_name {
    type = string
}

variable qos_policy_limit_by {
    type = string
}

variable qos_policy_is_default {
    type = bool
}

variable "max_reads_bw_mbps" {
  type = number
}

variable "max_writes_bw_mbps" {
  type = number
}

variable "min_writes_bw_mbps" {
  type = number
}

variable "min_reads_bw_mbps" {
  type = number
}

variable "max_reads_iops" {
  type = number
}

variable "max_writes_iops" {
  type = number
}

variable "min_reads_iops" {
  type = number
}

variable "min_writes_iops" {
  type = number
}

variable "burst_reads_bw_mb" {
  type = number
}

variable "burst_reads_loan_mb" {
  type = number
}

variable "burst_writes_bw_mb" {
  type = number
}

variable "burst_writes_loan_mb" {
  type = number
}

variable "burst_reads_iops" {
  type = number
}

variable "burst_reads_loan_iops" {
  type = number
}

variable "burst_writes_iops" {
  type = number
}

variable "burst_writes_loan_iops" {
  type = number
}

resource vastdata_user qos_user1 {
  name = var.qos_user_name
  uid = var.qos_user_uid
}


resource vastdata_qos_policy qos2 {
  name = var.qos_policy_name
  policy_type = "USER"

  attached_users_identifiers = tolist([tostring(vastdata_user.qos_user1.id)])

  limit_by = var.qos_policy_limit_by
  is_default = var.qos_policy_is_default

  static_limits {
    min_reads_bw_mbps  = var.min_reads_bw_mbps
    min_writes_bw_mbps = var.min_writes_bw_mbps
    max_reads_bw_mbps = var.max_reads_bw_mbps
    max_writes_bw_mbps = var.max_writes_bw_mbps
    max_reads_iops     = var.max_reads_iops
    max_writes_iops    = var.max_writes_iops
    min_reads_iops     = var.min_reads_iops
    min_writes_iops    = var.min_writes_iops
    burst_reads_bw_mb  = var.burst_reads_bw_mb
    burst_reads_loan_mb = var.burst_reads_loan_mb
    burst_writes_bw_mb = var.burst_writes_bw_mb
    burst_writes_loan_mb = var.burst_writes_loan_mb
    burst_reads_iops   = var.burst_reads_iops
    burst_reads_loan_iops = var.burst_reads_loan_iops
    burst_writes_iops  = var.burst_writes_iops
    burst_writes_loan_iops = var.burst_writes_loan_iops
  }

}

output tf_qos_policy {
  value = vastdata_qos_policy.qos2
}
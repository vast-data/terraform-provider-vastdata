
resource "vastdata_qos_policy" "vastdb_qos_policy" {
  name = "vastdb_qos_policy"

  static_limits = {
    max_writes_bw_mbps     = 1000
    max_reads_iops         = 2000
    max_writes_iops        = 3000
    burst_reads_bw_mb      = 2500
    burst_reads_iops       = 3000
    burst_reads_loan_iops  = 3500
    burst_reads_loan_mb    = 3000
    burst_writes_bw_mb     = 3500
    burst_writes_iops      = 4000
    burst_writes_loan_iops = 4500
    max_reads_bw_mbps      = 2000
    min_reads_bw_mbps      = 1500
    min_reads_iops         = 1000
    burst_writes_loan_mb   = 4000
    min_writes_bw_mbps     = 800
    min_writes_iops        = 1000
  }

  capacity_limits = {
    max_reads_bw_mbps_per_gb_capacity  = 100
    max_reads_iops_per_gb_capacity     = 200
    max_writes_bw_mbps_per_gb_capacity = 300
    max_writes_iops_per_gb_capacity    = 400
  }
}


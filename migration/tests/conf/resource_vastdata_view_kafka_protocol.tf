# Copyright (c) HashiCorp, Inc.

variable view_bucket_name {
    type = string
}

variable view_protocols {
    type = list(string)
}

variable user_name {
    type = string
}

variable user_uid {
    type = number
}

variable view_policy_name {
    type = string
}

variable tenant_name {
    type = string
}

variable tenant_client_ip_ranges {
    type = list(object({
        start_ip = string
        end_ip = string
    }))
}

variable vippool_name {
    type = string
}

variable vippool_range_start {
    type = string
}

variable vippool_range_end {
    type = string
}

variable kafka_rejoin_group_timeout_sec {
    type = number
}

variable kafka_first_join_group_timeout_sec {
    type = number
}

# Create a user for bucket ownership
resource vastdata_user user_for_kafka_view1 {
  name = var.user_name
  uid = var.user_uid
}

# Create a tenant
resource vastdata_tenant tenant_for_kafka_view1 {
  name = var.tenant_name
  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip = client_ip_ranges.value["end_ip"]
    }
  }
}

# Create VIP pool
resource vastdata_vip_pool vippool_for_kafka_view1 {
  name = var.vippool_name
  role = "PROTOCOLS"
  subnet_cidr = 24
  tenant_id = vastdata_tenant.tenant_for_kafka_view1.id
  ip_ranges {
    start_ip = var.vippool_range_start
    end_ip = var.vippool_range_end
  }
}

# Create a view policy with S3_NATIVE flavor
resource vastdata_view_policy view_policy_for_kafka1 {
  name = var.view_policy_name
  flavor = "S3_NATIVE"
  tenant_id = vastdata_tenant.tenant_for_kafka_view1.id
  nfs_no_squash = ["10.0.0.1","10.0.0.2"]
}

# Create a view with Kafka protocol
resource vastdata_view view_for_kafka1 {
  path = "/${var.view_bucket_name}"
  policy_id = vastdata_view_policy.view_policy_for_kafka1.id
  tenant_id = vastdata_tenant.tenant_for_kafka_view1.id
  create_dir = true
  protocols = var.view_protocols
  bucket = var.view_bucket_name
  bucket_owner = vastdata_user.user_for_kafka_view1.name
  kafka_vip_pools = [vastdata_vip_pool.vippool_for_kafka_view1.id]
  kafka_rejoin_group_timeout_sec = var.kafka_rejoin_group_timeout_sec
  kafka_first_join_group_timeout_sec = var.kafka_first_join_group_timeout_sec
  share_acl {
    enabled = false
    acl {
      name = vastdata_user.user_for_kafka_view1.name
      grantee = "users"
      fqdn = "All"
      permissions = "FULL"
    }
  }
}

output tf_view {
  value = vastdata_view.view_for_kafka1
}

output tf_view_policy {
  value = vastdata_view_policy.view_policy_for_kafka1
}

output tf_tenant {
  value = vastdata_tenant.tenant_for_kafka_view1
}

output tf_vippool {
  value = vastdata_vip_pool.vippool_for_kafka_view1
}
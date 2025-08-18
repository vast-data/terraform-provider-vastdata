# Copyright (c) HashiCorp, Inc.

# Example Terraform configuration using SDK v2 format
# This file demonstrates various patterns that need migration

terraform {
  required_providers {
    vastdata = {
      source = "vastdataorg/vastdata"
      version = "~> 0.9"
    }
  }
}

# Administrator manager with typo in resource name
resource "vastdata_administators_managers" "primary_admin" {
  name = "primary-administrator"
  enabled = true
  description = "Primary system administrator"
  
  # Block list that should become attributes
  capacity_limits {
    soft_limit = 5000
    hard_limit = 10000
    enabled = true
    alert_threshold = 80
  }
  
  # Block list that should become attributes list
  frames {
    name = "frame-01"
    ip = "10.0.1.1"
    port = 8080
    ssl_enabled = true
  }
  
  frames {
    name = "frame-02"
    ip = "10.0.1.2"
    port = 8081
    ssl_enabled = true
  }
  
  frames {
    name = "frame-03"
    ip = "10.0.1.3"
    port = 8082
    ssl_enabled = false
  }
  
  # List of numbers that should become string
  cnode_ids = [1, 2, 3, 4, 5]
  
  # List of strings that should become set (no visible change)
  permissions_list = ["read", "write", "admin", "backup"]
  
  # List of numbers that should become set
  roles = [100, 200, 300]
}

# Kafka broker with plural name
resource "vastdata_kafka_brokers" "primary_broker" {
  broker_id = 1
  name = "kafka-broker-primary"
  port = 9092
  
  # IP ranges that should become list of lists
  client_ip_ranges {
    start_ip = "192.168.1.1"
    end_ip = "192.168.1.100"
  }
  
  client_ip_ranges {
    start_ip = "10.0.0.1"
    end_ip = "10.0.0.50"
  }
  
  client_ip_ranges {
    start_ip = "172.16.0.1"
    end_ip = "172.16.0.25"
  }
  
  # Dynamic block that should be preserved
  dynamic "security_groups" {
    for_each = var.kafka_security_groups
    content {
      name = security_groups.value.name
      description = security_groups.value.description
      rules = security_groups.value.rules
      priority = security_groups.value.priority
    }
  }
  
  # List of strings that should become set
  object_types = ["message", "topic", "partition"]
}

# Resource with version suffix
resource "vastdata_active_directory2" "corporate_ad" {
  domain = "corp.example.com"
  server = "ad-primary.corp.example.com"
  backup_server = "ad-backup.corp.example.com"
  port = 389
  ssl_enabled = true
  
  # Mixed static and dynamic content
  default_user_quota {
    hard_limit = 1000000
    soft_limit = 800000
    grace_period = 7
  }
  
  default_group_quota {
    hard_limit = 5000000
    soft_limit = 4000000
    grace_period = 14
  }
  
  dynamic "user_groups" {
    for_each = var.active_directory_groups
    content {
      name = user_groups.value.name
      dn = user_groups.value.distinguished_name
      description = user_groups.value.description
      members = user_groups.value.members
    }
  }
  
  # List that needs transformation
  ldap_groups = ["domain_users", "domain_admins", "backup_operators"]
}

# Non-local user with underscores
resource "vastdata_non_local_user" "service_account" {
  username = "service-account-01"
  uid = 10001
  primary_group = "services"
  
  # List of numbers to set
  gids = [1000, 1001, 1002]
  
  # List of strings to set
  groups = ["services", "backup", "monitoring"]
}

# SAML configuration with short name
resource "vastdata_saml" "corporate_sso" {
  provider_name = "Corporate SSO"
  entity_id = "https://sso.corp.example.com"
  sso_url = "https://sso.corp.example.com/saml/sso"
  certificate = file("${path.module}/sso-cert.pem")
  
  # Dynamic attribute mapping
  dynamic "attribute_mappings" {
    for_each = var.saml_attribute_mappings
    content {
      saml_attribute = attribute_mappings.value.saml_name
      vastdata_attribute = attribute_mappings.value.vastdata_name
      required = attribute_mappings.value.required
    }
  }
}

# Multiple resources with different transformation needs
resource "vastdata_replication_peers" "backup_cluster" {
  name = "backup-cluster-west"
  hostname = "backup.west.example.com"
  port = 443
  
  ip_ranges {
    start_ip = "203.0.113.1"
    end_ip = "203.0.113.50"
  }
  
  ip_ranges {
    start_ip = "198.51.100.1"
    end_ip = "198.51.100.25"
  }
}

resource "vastdata_s3_replication_peers" "s3_backup" {
  name = "s3-backup-peer"
  endpoint = "https://s3.us-west-2.amazonaws.com"
  bucket = "vastdata-backup-bucket"
  
  s3_policies_ids = [101, 102, 103]
  
  # Users and groups for access control
  users = ["backup-user", "admin-user", "monitoring-user"]
  bucket_creators_groups = ["backup-team", "admin-team"]
}

# Variables and outputs that should not be transformed
variable "kafka_security_groups" {
  description = "Security groups for Kafka broker"
  type = list(object({
    name = string
    description = string
    rules = list(string)
    priority = number
  }))
  default = []
}

variable "active_directory_groups" {
  description = "Active Directory groups to configure"
  type = list(object({
    name = string
    distinguished_name = string
    description = string
    members = list(string)
  }))
  default = []
}

output "admin_manager_id" {
  description = "ID of the primary administrator manager"
  value = vastdata_administators_managers.primary_admin.id
}

output "kafka_broker_endpoint" {
  description = "Kafka broker endpoint"
  value = "${vastdata_kafka_brokers.primary_broker.hostname}:${vastdata_kafka_brokers.primary_broker.port}"
}
